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

	return &Compiler{
		constants: []object.Object{},

		scopes:     []CompilationScope{mainScope},
		scopeIndex: 0,

		symbolTable: NewSymbolTable(),
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
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpSetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Value) // Compile-time error
		}

		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpGetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpGetLocal, symbol.Index)
		}

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
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(bytecode.OpGreaterThan)
			return nil
		}

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

		case "==":
			c.emit(bytecode.OpEqual)
		case "!=":
			c.emit(bytecode.OpNotEqual)
		case ">":
			c.emit(bytecode.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus offset to be updated below with the position following the consequence
		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus offset to be updated below with the position following the alternative
		jumpPos := c.emit(bytecode.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(bytecode.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))

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

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		compiledFunction := &object.CompiledFunction{Instructions: instructions, NumLocals: numLocals}
		c.emit(bytecode.OpConstant, c.addConstant(compiledFunction))

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

		c.emit(bytecode.OpCall)
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
