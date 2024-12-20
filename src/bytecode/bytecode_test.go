package bytecode

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpEqual, []int{}, []byte{byte(OpEqual)}},
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpCall, []int{2}, []byte{byte(OpCall), 2}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		if len(instruction) != len(test.expected) {
			t.Errorf("instruction has the wrong length. expected=%d, got=%d", len(test.expected), len(instruction))
		}

		for i, b := range test.expected {
			if instruction[i] != test.expected[i] {
				t.Errorf("wrong byte at position %d. expected=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpTrue),
		Make(OpGetLocal, 255),
		Make(OpClosure, 65535, 255),
	}

	expected := "0000 OpAdd\n0001 OpConstant 1\n0004 OpConstant 2\n0007 OpConstant 65535\n0010 OpTrue\n0011 OpGetLocal 255\n0013 OpClosure 65535 255\n"

	concatted := Instructions{}
	for _, instr := range instructions {
		concatted = append(concatted, instr...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nexpected=%q\ngot=%q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		bytesRead int
		operands  []int
	}{
		{OpConstant, 2, []int{65535}},
		{OpTrue, 0, []int{}},
		{OpFalse, 0, []int{}},
		{OpAdd, 0, []int{}},
		{OpGreaterThan, 0, []int{}},
		{OpBang, 0, []int{}},
		{OpGetLocal, 1, []int{255}},
		{OpClosure, 3, []int{65535, 255}},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		def, err := LookUp(byte(test.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != test.bytesRead {
			t.Fatalf("number of bytes read wrong. expected=%d, got=%d", test.bytesRead, n)
		}

		for i, expected := range test.operands {
			if operandsRead[i] != expected {
				t.Fatalf("operand at position %d wrong. expected=%d, got=%d", i, expected, operandsRead[i])
			}
		}
	}
}
