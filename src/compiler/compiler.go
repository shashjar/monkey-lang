package compiler

import (
	"monkey/ast"
	"monkey/bytecode"
	"monkey/object"
)

// Represents bytecode generated and constants evaluated by the compiler.
type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []object.Object
}

// Represents a compiler for the Monkey programming language, generating bytecode instructions to execute.
type Compiler struct {
	instructions bytecode.Instructions
	constants    []object.Object
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []object.Object{},
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

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
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

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))
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
	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(instr []byte) int {
	c.instructions = append(c.instructions, instr...)
	return len(c.instructions) - 1
}
