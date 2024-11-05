package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/bytecode"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedInstructions []bytecode.Instructions
	expectedConstants    []interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "1 + 2",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
			},
			expectedConstants: []interface{}{1, 2},
		},
	}

	runCompilerTests(t, tests)
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.ParseProgram()
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, test := range tests {
		program := parse(test.input)

		compiler := NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(test.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(test.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func testInstructions(expected []bytecode.Instructions, actual bytecode.Instructions) error {
	concatted := concatInstructions(expected)

	if len(concatted) != len(actual) {
		return fmt.Errorf("instructions are the wrong length.\nexpected=%q\ngot=%q", concatted, actual)
	}

	for i, expInstr := range concatted {
		if expInstr != actual[i] {
			return fmt.Errorf("wrong instruction at position %d.\nexpected=%q\ngot=%q", i, concatted, actual)
		}
	}

	return nil
}

func testConstants(expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("constants are the wrong length. expected=%d, got=%d", len(expected), len(actual))
	}

	for i, expConst := range expected {
		switch expConst := expConst.(type) {
		case int:
			err := testIntegerObject(int64(expConst), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func concatInstructions(s []bytecode.Instructions) bytecode.Instructions {
	result := bytecode.Instructions{}
	for _, instr := range s {
		result = append(result, instr...)
	}
	return result
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not an Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong integer value. expected=%d, got=%d", expected, result.Value)
	}

	return nil
}
