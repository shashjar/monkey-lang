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
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
