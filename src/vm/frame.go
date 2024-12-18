package vm

import (
	"monkey/bytecode"
	"monkey/object"
)

// Represents a call frame (i.e. stack frame) used to track function call information on the stack.
type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.fn.Instructions
}
