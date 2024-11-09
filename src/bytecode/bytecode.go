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

	OpEqual
	OpNotEqual
	OpGreaterThan

	OpMinus
	OpBang
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

	OpAdd: {"OpAdd", []int{}},
	OpSub: {"OpSub", []int{}},
	OpMul: {"OpMul", []int{}},
	OpDiv: {"OpDiv", []int{}},

	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},

	OpMinus: {"OpMinus", []int{}},
	OpBang:  {"OpBang", []int{}},
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
