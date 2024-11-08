package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/bytecode"
	"monkey/object"
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

// Represents a compiler for the Monkey programming language, generating bytecode instructions to execute.
type Compiler struct {
	instructions bytecode.Instructions
	constants    []object.Object

	lastInstruction         EmittedInstruction // The latest instruction emitted by the compiler.
	previousLastInstruction EmittedInstruction // The second-to-latest instruction emitted by the compiler.
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []object.Object{},

		lastInstruction:         EmittedInstruction{},
		previousLastInstruction: EmittedInstruction{},
	}
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

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		if node.Alternative == nil {
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// Emit an `OpJump` with a bogus offset to be updated below with the position following the alternative
			jumpPos := c.emit(bytecode.OpJump, 9999)

			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

			afterAlternativePos := len(c.instructions)
			c.changeOperand(jumpPos, afterAlternativePos)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	instr := bytecode.Make(op, operands...)
	pos := c.addInstruction(instr)

	c.previousLastInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{Opcode: op, Position: pos}

	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(instr []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instr...)
	return posNewInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstr []byte) {
	for i := 0; i < len(newInstr); i++ {
		c.instructions[pos+i] = newInstr[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	fmt.Printf("Changing operand at pos %d to %d\n", opPos, operand)
	fmt.Printf("Instructions length: %d\n", len(c.instructions))

	op := bytecode.Opcode(c.instructions[opPos])
	newInstruction := bytecode.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == bytecode.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousLastInstruction
}
