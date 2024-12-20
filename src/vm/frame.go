package vm

import (
	"monkey/bytecode"
	"monkey/object"
)

// Represents a call frame (i.e. stack frame) used to track function call information on the stack.
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.cl.Fn.Instructions
}
