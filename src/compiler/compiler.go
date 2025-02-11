package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/bytecode"
	"monkey/object"
	"sort"
)

// Represents bytecode generated and constants evaluated by the compiler.
type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []object.Object
}

// Represents an instruction that was emitted by the compiler.
type EmittedInstruction struct {
	Opcode   bytecode.Opcode
	Position int
}

// Represents the scope of compilation.
type CompilationScope struct {
	instructions            bytecode.Instructions
	lastInstruction         EmittedInstruction // The latest instruction emitted by the compiler.
	previousLastInstruction EmittedInstruction // The second-to-latest instruction emitted by the compiler.
}

// Represents a compiler for the Monkey programming language, generating bytecode instructions to execute.
type Compiler struct {
	constants []object.Object

	scopes     []CompilationScope
	scopeIndex int

	symbolTable *SymbolTable // The symbol table for the compiler to use for identifier associations (bindings).
}

func NewCompiler() *Compiler {
	mainScope := CompilationScope{
		instructions:            bytecode.Instructions{},
		lastInstruction:         EmittedInstruction{},
		previousLastInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	for i, builtIn := range object.BuiltIns {
		symbolTable.DefineBuiltIn(i, builtIn.Name)
	}

	return &Compiler{
		constants: []object.Object{},

		scopes:     []CompilationScope{mainScope},
		scopeIndex: 0,

		symbolTable: symbolTable,
	}
}

func NewCompilerWithState(st *SymbolTable, constants []object.Object) *Compiler {
	compiler := NewCompiler()
	compiler.symbolTable = st
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpPop)

	case *ast.LetStatement:
		symbol, ok := c.symbolTable.store[node.Name.Value] // Only able to declare this variable if it hasn't already been declared
		if ok && (symbol.Scope == GlobalScope || symbol.Scope == LocalScope) {
			return fmt.Errorf("line %d, column %d: identifier '%s' has already been declared", node.Token.LineNumber, node.Token.ColumnNumber, node.Name.Value)
		}

		symbol = c.symbolTable.Define(node.Name.Value)

		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpSetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpSetLocal, symbol.Index)
		}

	case *ast.ConstStatement:
		symbol, ok := c.symbolTable.store[node.Name.Value] // Only able to declare this variable if it hasn't already been declared
		if ok && (symbol.Scope == GlobalScope || symbol.Scope == LocalScope) {
			return fmt.Errorf("line %d, column %d: identifier '%s' has already been declared", node.Token.LineNumber, node.Token.ColumnNumber, node.Name.Value)
		}

		symbol = c.symbolTable.DefineConst(node.Name.Value)

		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpSetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpSetLocal, symbol.Index)
		}

	case *ast.AssignStatement:
		symbol, ok := c.symbolTable.store[node.Name.Value] // Only able to reassign value if the variable was declared in the same scope we're currently in
		if !ok {
			return fmt.Errorf("line %d, column %d: attempting to assign value to identifier '%s' prior to declaration", node.Token.LineNumber, node.Token.ColumnNumber, node.Name.Value)
		}
		if symbol.Const {
			return fmt.Errorf("line %d, column %d: attempting to assign value to constant variable '%s'", node.Token.LineNumber, node.Token.ColumnNumber, node.Name.Value)
		}

		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpSetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("line %d, column %d: undefined variable: %s", node.Token.LineNumber, node.Token.ColumnNumber, node.Value) // Compile-time error
		}

		c.loadSymbol(symbol)

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(bytecode.OpMinus)
		case "!":
			c.emit(bytecode.OpBang)
		default:
			return fmt.Errorf("line %d, column %d: unknown operator: %s", node.Token.LineNumber, node.Token.ColumnNumber, node.Operator)
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		case "-":
			c.emit(bytecode.OpSub)
		case "*":
			c.emit(bytecode.OpMul)
		case "/":
			c.emit(bytecode.OpDiv)
		case "//":
			c.emit(bytecode.OpIntegerDiv)
		case "**":
			c.emit(bytecode.OpExp)
		case "%":
			c.emit(bytecode.OpMod)

		case "&&":
			c.emit(bytecode.OpAnd)
		case "||":
			c.emit(bytecode.OpOr)

		case "==":
			c.emit(bytecode.OpEqual)
		case "!=":
			c.emit(bytecode.OpNotEqual)
		case "<":
			c.emit(bytecode.OpLessThan)
		case ">":
			c.emit(bytecode.OpGreaterThan)
		case "<=":
			c.emit(bytecode.OpLessThanOrEqualTo)
		case ">=":
			c.emit(bytecode.OpGreaterThanOrEqualTo)
		default:
			return fmt.Errorf("line %d, column %d: unknown operator: %s", node.Token.LineNumber, node.Token.ColumnNumber, node.Operator)
		}

	case *ast.IfExpression:
		jumpPositions := []int{}

		for _, clause := range node.Clauses {
			err := c.Compile(clause.Condition)
			if err != nil {
				return err
			}

			// Emit an `OpJumpNotTruthy` with a bogus offset to be updated below with the position following this clause's consequence
			jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

			err = c.Compile(clause.Consequence)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}

			// Emit an `OpJump` with a bogus offset to be updated below with the position following the end of the entire if expression
			jumpPos := c.emit(bytecode.OpJump, 9999)
			jumpPositions = append(jumpPositions, jumpPos)

			afterClauseConsequencePos := len(c.currentInstructions())
			c.changeOperand(jumpNotTruthyPos, afterClauseConsequencePos)
		}

		if node.Alternative == nil {
			c.emit(bytecode.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}
		}

		afterIfExpressionPos := len(c.currentInstructions())
		for _, jumpPos := range jumpPositions {
			c.changeOperand(jumpPos, afterIfExpressionPos)
		}

	case *ast.SwitchStatement:
		jumpPositions := []int{}

		for _, switchCase := range node.Cases {
			err := c.Compile(node.SwitchExpression)
			if err != nil {
				return err
			}

			err = c.Compile(switchCase.Expression)
			if err != nil {
				return err
			}

			c.emit(bytecode.OpEqual)

			// Emit an `OpJumpNotTruthy` with a bogus offset to be updated below with the position following this case's consequence
			jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

			err = c.Compile(switchCase.Consequence)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}

			// Emit an `OpJump` with a bogus offset to be updated below with the position following the end of the entire switch statement
			jumpPos := c.emit(bytecode.OpJump, 9999)
			jumpPositions = append(jumpPositions, jumpPos)

			afterCaseConsequencePos := len(c.currentInstructions())
			c.changeOperand(jumpNotTruthyPos, afterCaseConsequencePos)
		}

		if node.Default == nil {
			c.emit(bytecode.OpNull)
		} else {
			err := c.Compile(node.Default)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}
		}

		afterSwitchStatementPos := len(c.currentInstructions())
		for _, jumpPos := range jumpPositions {
			c.changeOperand(jumpPos, afterSwitchStatementPos)
		}

	case *ast.WhileLoop:
		whileLoopStartPos := len(c.currentInstructions())

		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus offset to be updated below with the position following the loop body
		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` to go back to the start of the loop
		c.emit(bytecode.OpJump, whileLoopStartPos)

		afterWhileLoopPos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterWhileLoopPos)

		// Emit an OpNull so that the OpPop emitted after this while loop is compiled doesn't change anything
		c.emit(bytecode.OpNull)

	case *ast.ForLoop:
		err := c.Compile(node.Init)
		if err != nil {
			return err
		}

		forLoopConditionStartPos := len(c.currentInstructions())

		err = c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus offset to be updated below with the position following the loop body
		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		err = c.Compile(node.Afterthought)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` to go back to the start of the loop
		c.emit(bytecode.OpJump, forLoopConditionStartPos)

		afterForLoopPos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterForLoopPos)

		// Emit an OpNull so that the OpPop emitted after this for loop is compiled doesn't change anything
		c.emit(bytecode.OpNull)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))

	case *ast.Float:
		float := &object.Float{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(float))

	case *ast.Boolean:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, element := range node.Elements {
			err := c.Compile(element)
			if err != nil {
				return err
			}
		}
		c.emit(bytecode.OpArray, len(node.Elements))

	case *ast.HashMapLiteral:
		keys := []ast.Expression{}
		for k := range node.KVPairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i int, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}

			err = c.Compile(node.KVPairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(bytecode.OpHashMap, len(node.KVPairs)*2)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(bytecode.OpIndex)

	case *ast.FunctionLiteral:
		c.enterScope()

		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name)
		}

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(bytecode.OpReturnValue) {
			c.emit(bytecode.OpReturn)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, fs := range freeSymbols {
			c.loadSymbol(fs)
		}

		compiledFunction := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}
		fnIndex := c.addConstant(compiledFunction)
		c.emit(bytecode.OpClosure, fnIndex, len(freeSymbols))

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(bytecode.OpReturnValue)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, argExp := range node.Arguments {
			err = c.Compile(argExp)
			if err != nil {
				return err
			}
		}

		c.emit(bytecode.OpCall, len(node.Arguments))
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	instr := bytecode.Make(op, operands...)
	pos := c.addInstruction(instr)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) {
	newPrevious := c.scopes[c.scopeIndex].lastInstruction
	newLast := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousLastInstruction = newPrevious
	c.scopes[c.scopeIndex].lastInstruction = newLast
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(instr []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), instr...)
	c.scopes[c.scopeIndex].instructions = updatedInstructions
	return posNewInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstr []byte) {
	instructions := c.currentInstructions()
	for i := 0; i < len(newInstr); i++ {
		instructions[pos+i] = newInstr[i]
	}
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, bytecode.Make(bytecode.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = bytecode.OpReturnValue
}

func (c *Compiler) currentInstructions() bytecode.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := bytecode.Opcode(c.currentInstructions()[opPos])
	newInstruction := bytecode.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIs(op bytecode.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousLastInstruction

	oldInstructions := c.currentInstructions()
	newInstructions := oldInstructions[:last.Position]

	c.scopes[c.scopeIndex].instructions = newInstructions
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) enterScope() {
	newScope := CompilationScope{
		instructions:            bytecode.Instructions{},
		lastInstruction:         EmittedInstruction{},
		previousLastInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, newScope)
	c.scopeIndex += 1

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() bytecode.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex -= 1

	c.symbolTable = c.symbolTable.outer

	return instructions
}

func (c *Compiler) loadSymbol(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(bytecode.OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(bytecode.OpGetLocal, symbol.Index)
	case FreeScope:
		c.emit(bytecode.OpGetFreeVar, symbol.Index)
	case FunctionScope:
		c.emit(bytecode.OpCurrentClosure)
	case BuiltInScope:
		c.emit(bytecode.OpGetBuiltIn, symbol.Index)
	}
}
