package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Represents an opcode of size 1 byte, indicating some operation with some number of operands.
type Opcode byte

const (
	OpConstant Opcode = iota
	OpTrue
	OpFalse
	OpNull

	OpPop
	OpJumpNotTruthy
	OpJump

	OpAdd
	OpSub
	OpMul
	OpDiv
	OpIntegerDiv
	OpExp
	OpMod

	OpAnd
	OpOr

	OpEqual
	OpNotEqual
	OpLessThan
	OpGreaterThan
	OpLessThanOrEqualTo
	OpGreaterThanOrEqualTo

	OpMinus
	OpBang

	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal

	OpArray
	OpHashMap
	OpIndex

	OpCall
	OpReturnValue
	OpReturn
	OpGetBuiltIn
	OpClosure
	OpGetFreeVar
	OpCurrentClosure
)

// Represents a set of instructions as a slice of bytes.
type Instructions []byte

func (instr Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(instr) {
		def, err := LookUp(instr[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, instr[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, instr.fmtInstruction(def, operands))
		i += 1 + read
	}

	return out.String()
}

func (instr Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand length %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Represents the definition for an Opcode, with some readable name and the number of
// bytes that each operand takes up.
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpTrue:     {"OpTrue", []int{}},
	OpFalse:    {"OpFalse", []int{}},
	OpNull:     {"OpNull", []int{}},

	OpPop:           {"OpPop", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},

	OpAdd:        {"OpAdd", []int{}},
	OpSub:        {"OpSub", []int{}},
	OpMul:        {"OpMul", []int{}},
	OpDiv:        {"OpDiv", []int{}},
	OpIntegerDiv: {"OpIntegerDiv", []int{}},
	OpExp:        {"OpExp", []int{}},
	OpMod:        {"OpMod", []int{}},

	OpAnd: {"OpAnd", []int{}},
	OpOr:  {"OpOr", []int{}},

	OpEqual:                {"OpEqual", []int{}},
	OpNotEqual:             {"OpNotEqual", []int{}},
	OpLessThan:             {"OpLessThan", []int{}},
	OpGreaterThan:          {"OpGreaterThan", []int{}},
	OpLessThanOrEqualTo:    {"OpLessThanOrEqualTo", []int{}},
	OpGreaterThanOrEqualTo: {"OpGreaterThanOrEqualTo", []int{}},

	OpMinus: {"OpMinus", []int{}},
	OpBang:  {"OpBang", []int{}},

	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpSetGlobal: {"OpSetGlobal", []int{2}},
	OpGetLocal:  {"OpGetLocal", []int{1}},
	OpSetLocal:  {"OpSetLocal", []int{1}},

	OpArray:   {"OpArray", []int{2}},
	OpHashMap: {"OpHashMap", []int{2}},
	OpIndex:   {"OpIndex", []int{}},

	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetBuiltIn:     {"OpGetBuiltIn", []int{1}},
	OpClosure:        {"OpClosure", []int{2, 1}}, // First operand: constant index of *object.CompiledFunction. Second operand: number of free variables in the closure.
	OpGetFreeVar:     {"OpGetFreeVar", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

func LookUp(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d is undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, instr Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(instr[offset:]))
		case 2:
			operands[i] = int(ReadUint16(instr[offset:]))
		}
		offset += width
	}

	return operands, offset
}

func ReadUint16(instr Instructions) uint16 {
	return binary.BigEndian.Uint16(instr)
}

func ReadUint8(instr Instructions) uint8 {
	return uint8(instr[0])
}
