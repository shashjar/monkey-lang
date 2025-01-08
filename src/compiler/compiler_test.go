package compiler

import (
	"fmt"
	"math"
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

type compilerErrorTestCase struct {
	input         string
	expectedError string
}

func TestArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "-2",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{2},
		},
		{
			input: "-8.64",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{8.64},
		},
		{
			input: "1; 2;",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2},
		},
		{
			input: "1 + 2",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2},
		},
		{
			input: "3 - 4",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{3, 4},
		},
		{
			input: "2 * 8",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{2, 8},
		},
		{
			input: "15 / 5.45",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{15, 5.45},
		},
		{
			input: "22 // 7",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpIntegerDiv),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{22, 7},
		},
		{
			input: "17.213 // 6.182",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpIntegerDiv),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{17.213, 6.182},
		},
		{
			input: "4 ** 6",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpExp),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{4, 6},
		},
		{
			input: "14 % 5",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMod),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{14, 5},
		},
		{
			input: "3.1415 * 2.718 - 1",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{3.1415, 2.718, 1},
		},
	}

	runCompilerTests(t, tests)
}

func TestPostfixOperators(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "let x = 0; x++;",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),

				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpSetGlobal, 0),
			},
			expectedConstants: []interface{}{0, 1},
		},
		{
			input: "let x = 0; x--;",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),

				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpSetGlobal, 0),
			},
			expectedConstants: []interface{}{0, 1},
		},
	}

	runCompilerTests(t, tests)
}

func TestPostfixOperatorErrors(t *testing.T) {
	tests := []compilerErrorTestCase{
		{
			input:         "x++",
			expectedError: "line 1, column 0: attempting to assign value to identifier 'x' prior to declaration",
		},
		{
			input:         "x--",
			expectedError: "line 1, column 0: attempting to assign value to identifier 'x' prior to declaration",
		},
		{
			input:         "const x = 5; x++;",
			expectedError: "line 1, column 13: attempting to assign value to constant variable 'x'",
		},
	}

	runCompilerErrorTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "true",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "false",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "!true",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "!false",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "3 == 3",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{3, 3},
		},
		{
			input: "3 == 7",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{3, 7},
		},
		{
			input: "18 != 16",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{18, 16},
		},
		{
			input: "1 != 1",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 1},
		},
		{
			input: "6 < 8",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpLessThan),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{6, 8},
		},
		{
			input: "1 > 2",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2},
		},
		{
			input: "4 > 4",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{4, 4},
		},
		{
			input: "10 <= 12",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpLessThanOrEqualTo),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{10, 12},
		},
		{
			input: "3 >= 4",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThanOrEqualTo),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{3, 4},
		},
		{
			input: "true == false",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "true != false",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "true && false",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpAnd),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "false || true",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpOr),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) { 10 }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 11),
				// 0010
				bytecode.Make(bytecode.OpNull),
				// 0011
				bytecode.Make(bytecode.OpPop),
				// 0012
				bytecode.Make(bytecode.OpConstant, 1),
				// 0015
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{10, 3333},
		},
		{
			input: `
			if (true) { 10 } else if (false) { 20 }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 21),
				// 0010
				bytecode.Make(bytecode.OpFalse),
				// 0011
				bytecode.Make(bytecode.OpJumpNotTruthy, 20),
				// 0014
				bytecode.Make(bytecode.OpConstant, 1),
				// 0017
				bytecode.Make(bytecode.OpJump, 21),
				// 0020
				bytecode.Make(bytecode.OpNull),
				// 0021
				bytecode.Make(bytecode.OpPop),
				// 0022
				bytecode.Make(bytecode.OpConstant, 2),
				// 0025
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{10, 20, 3333},
		},
		{
			input: `
			if (true) { 10 } else { 20 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 13),
				// 0010
				bytecode.Make(bytecode.OpConstant, 1),
				// 0013
				bytecode.Make(bytecode.OpPop),
				// 0014
				bytecode.Make(bytecode.OpConstant, 2),
				// 0017
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			if (true) { 10 } else if (false) { 20 } else if (true) { 30 } else { 40 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 30, 40, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 33),
				// 0010
				bytecode.Make(bytecode.OpFalse),
				// 0011
				bytecode.Make(bytecode.OpJumpNotTruthy, 20),
				// 0014
				bytecode.Make(bytecode.OpConstant, 1),
				// 0017
				bytecode.Make(bytecode.OpJump, 33),
				// 0020
				bytecode.Make(bytecode.OpTrue),
				// 0021
				bytecode.Make(bytecode.OpJumpNotTruthy, 30),
				// 0024
				bytecode.Make(bytecode.OpConstant, 2),
				// 0027
				bytecode.Make(bytecode.OpJump, 33),
				// 0030
				bytecode.Make(bytecode.OpConstant, 3),
				// 0033
				bytecode.Make(bytecode.OpPop),
				// 0034
				bytecode.Make(bytecode.OpConstant, 4),
				// 0037
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestSwitchStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			switch "hello" {
			case "hello":
				10;
			}
			3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpConstant, 1),
				// 0006
				bytecode.Make(bytecode.OpEqual),
				// 0007
				bytecode.Make(bytecode.OpJumpNotTruthy, 16),
				// 0010
				bytecode.Make(bytecode.OpConstant, 2),
				// 0013
				bytecode.Make(bytecode.OpJump, 17),
				// 0016
				bytecode.Make(bytecode.OpNull),
				// 0017
				bytecode.Make(bytecode.OpPop),
				// 0018
				bytecode.Make(bytecode.OpConstant, 3),
				// 0021
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{"hello", "hello", 10, 3333},
		},
		{
			input: `
			switch "hello" {
			case "hello":
				10;
			case "world":
				20;
			}
			3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpConstant, 1),
				// 0006
				bytecode.Make(bytecode.OpEqual),
				// 0007
				bytecode.Make(bytecode.OpJumpNotTruthy, 16),
				// 0010
				bytecode.Make(bytecode.OpConstant, 2),
				// 0013
				bytecode.Make(bytecode.OpJump, 33),
				// 0016
				bytecode.Make(bytecode.OpConstant, 3),
				// 0019
				bytecode.Make(bytecode.OpConstant, 4),
				// 0022
				bytecode.Make(bytecode.OpEqual),
				// 0023
				bytecode.Make(bytecode.OpJumpNotTruthy, 32),
				// 0026
				bytecode.Make(bytecode.OpConstant, 5),
				// 0029
				bytecode.Make(bytecode.OpJump, 33),
				// 0032
				bytecode.Make(bytecode.OpNull),
				// 0033
				bytecode.Make(bytecode.OpPop),
				// 0034
				bytecode.Make(bytecode.OpConstant, 6),
				// 0037
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{"hello", "hello", 10, "hello", "world", 20, 3333},
		},
		{
			input: `
			switch "hello" {
			case "hello":
				10;
			case "world":
				20;
			default:
				30;
			}
			3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpConstant, 1),
				// 0006
				bytecode.Make(bytecode.OpEqual),
				// 0007
				bytecode.Make(bytecode.OpJumpNotTruthy, 16),
				// 0010
				bytecode.Make(bytecode.OpConstant, 2),
				// 0013
				bytecode.Make(bytecode.OpJump, 35),
				// 0016
				bytecode.Make(bytecode.OpConstant, 3),
				// 0019
				bytecode.Make(bytecode.OpConstant, 4),
				// 0022
				bytecode.Make(bytecode.OpEqual),
				// 0023
				bytecode.Make(bytecode.OpJumpNotTruthy, 32),
				// 0026
				bytecode.Make(bytecode.OpConstant, 5),
				// 0029
				bytecode.Make(bytecode.OpJump, 35),
				// 0032
				bytecode.Make(bytecode.OpConstant, 6),
				// 0035
				bytecode.Make(bytecode.OpPop),
				// 0036
				bytecode.Make(bytecode.OpConstant, 7),
				// 0039
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{"hello", "hello", 10, "hello", "world", 20, 30, 3333},
		},
	}

	runCompilerTests(t, tests)
}

func TestWhileLoops(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			while (true) { 10; }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 0),
				// 0010
				bytecode.Make(bytecode.OpNull),
				// 0011
				bytecode.Make(bytecode.OpPop),
				// 0012
				bytecode.Make(bytecode.OpConstant, 1),
				// 0015
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{10, 3333},
		},
		{
			input: `
			let x = 0; while (x < 10) { x = x + 1; }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0006
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0009
				bytecode.Make(bytecode.OpConstant, 1),
				// 0012
				bytecode.Make(bytecode.OpLessThan),
				// 0013
				bytecode.Make(bytecode.OpJumpNotTruthy, 29),
				// 0016
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0019
				bytecode.Make(bytecode.OpConstant, 2),
				// 0022
				bytecode.Make(bytecode.OpAdd),
				// 0023
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0026
				bytecode.Make(bytecode.OpJump, 6),
				// 0029
				bytecode.Make(bytecode.OpNull),
				// 0030
				bytecode.Make(bytecode.OpPop),
				// 0031
				bytecode.Make(bytecode.OpConstant, 3),
				// 0034
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{0, 10, 1, 3333},
		},
		{
			input: `
			let x = 0; while (x < 10) { x++ }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0006
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0009
				bytecode.Make(bytecode.OpConstant, 1),
				// 0012
				bytecode.Make(bytecode.OpLessThan),
				// 0013
				bytecode.Make(bytecode.OpJumpNotTruthy, 29),
				// 0016
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0019
				bytecode.Make(bytecode.OpConstant, 2),
				// 0022
				bytecode.Make(bytecode.OpAdd),
				// 0023
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0026
				bytecode.Make(bytecode.OpJump, 6),
				// 0029
				bytecode.Make(bytecode.OpNull),
				// 0030
				bytecode.Make(bytecode.OpPop),
				// 0031
				bytecode.Make(bytecode.OpConstant, 3),
				// 0034
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{0, 10, 1, 3333},
		},
		{
			input: `
			let x = 0; while (x < 10) { x++; }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0006
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0009
				bytecode.Make(bytecode.OpConstant, 1),
				// 0012
				bytecode.Make(bytecode.OpLessThan),
				// 0013
				bytecode.Make(bytecode.OpJumpNotTruthy, 29),
				// 0016
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0019
				bytecode.Make(bytecode.OpConstant, 2),
				// 0022
				bytecode.Make(bytecode.OpAdd),
				// 0023
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0026
				bytecode.Make(bytecode.OpJump, 6),
				// 0029
				bytecode.Make(bytecode.OpNull),
				// 0030
				bytecode.Make(bytecode.OpPop),
				// 0031
				bytecode.Make(bytecode.OpConstant, 3),
				// 0034
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{0, 10, 1, 3333},
		},
	}

	runCompilerTests(t, tests)
}

func TestForLoops(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			for (let i = 0; true; i = i + 1) { puts(i); }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpSetGlobal, 0),

				// 0006
				bytecode.Make(bytecode.OpTrue),
				// 0007
				bytecode.Make(bytecode.OpJumpNotTruthy, 30),

				// 0010
				bytecode.Make(bytecode.OpGetBuiltIn, 0),
				// 0012
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0015
				bytecode.Make(bytecode.OpCall, 1),

				// 0017
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0020
				bytecode.Make(bytecode.OpConstant, 1),
				// 0023
				bytecode.Make(bytecode.OpAdd),
				// 0024
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0027
				bytecode.Make(bytecode.OpJump, 6),

				// 0030
				bytecode.Make(bytecode.OpNull),
				// 0031
				bytecode.Make(bytecode.OpPop),

				// 0032
				bytecode.Make(bytecode.OpConstant, 2),
				// 0035
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{0, 1, 3333},
		},
		{
			input: `
			for (let i = 0; true; i++) { puts(i); }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpSetGlobal, 0),

				// 0006
				bytecode.Make(bytecode.OpTrue),
				// 0007
				bytecode.Make(bytecode.OpJumpNotTruthy, 30),

				// 0010
				bytecode.Make(bytecode.OpGetBuiltIn, 0),
				// 0012
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0015
				bytecode.Make(bytecode.OpCall, 1),

				// 0017
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0020
				bytecode.Make(bytecode.OpConstant, 1),
				// 0023
				bytecode.Make(bytecode.OpAdd),
				// 0024
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// 0027
				bytecode.Make(bytecode.OpJump, 6),

				// 0030
				bytecode.Make(bytecode.OpNull),
				// 0031
				bytecode.Make(bytecode.OpPop),

				// 0032
				bytecode.Make(bytecode.OpConstant, 2),
				// 0035
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{0, 1, 3333},
		},
		{
			input: `
			let arr = [1, 2, 3]; for (let i = 0; i < len(arr); i = i + 1) { puts(arr[i]); }; 3333;
			`,
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpConstant, 0),
				// 0003
				bytecode.Make(bytecode.OpConstant, 1),
				// 0006
				bytecode.Make(bytecode.OpConstant, 2),
				// 0009
				bytecode.Make(bytecode.OpArray, 3),
				// 0012
				bytecode.Make(bytecode.OpSetGlobal, 0),

				// 0015
				bytecode.Make(bytecode.OpConstant, 3),
				// 0018
				bytecode.Make(bytecode.OpSetGlobal, 1),

				// 0021
				bytecode.Make(bytecode.OpGetGlobal, 1),
				// 0024
				bytecode.Make(bytecode.OpGetBuiltIn, 1),
				// 0026
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0029
				bytecode.Make(bytecode.OpCall, 1),
				// 0031
				bytecode.Make(bytecode.OpLessThan),
				// 0032
				bytecode.Make(bytecode.OpJumpNotTruthy, 59),

				// 0035
				bytecode.Make(bytecode.OpGetBuiltIn, 0),
				// 0037
				bytecode.Make(bytecode.OpGetGlobal, 0),
				// 0040
				bytecode.Make(bytecode.OpGetGlobal, 1),
				// 0043
				bytecode.Make(bytecode.OpIndex),
				// 0044
				bytecode.Make(bytecode.OpCall, 1),

				// 0046
				bytecode.Make(bytecode.OpGetGlobal, 1),
				// 0049
				bytecode.Make(bytecode.OpConstant, 4),
				// 0052
				bytecode.Make(bytecode.OpAdd),
				// 0053
				bytecode.Make(bytecode.OpSetGlobal, 1),
				// 0056
				bytecode.Make(bytecode.OpJump, 21),

				// 0059
				bytecode.Make(bytecode.OpNull),
				// 0060
				bytecode.Make(bytecode.OpPop),

				// 0061
				bytecode.Make(bytecode.OpConstant, 5),
				// 0064
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3, 0, 1, 3333},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
			},
			expectedConstants: []interface{}{1, 2},
		},
		{
			input: `
			let one = 1;
			one;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalConstStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			const one = 1;
			const two = 2;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
			},
			expectedConstants: []interface{}{1, 2},
		},
		{
			input: `
			const one = 1;
			one;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1},
		},
		{
			input: `
			const one = 1;
			const two = one;
			two;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1},
		},
	}

	runCompilerTests(t, tests)
}

func TestConstStatementsErrors(t *testing.T) {
	tests := []compilerErrorTestCase{
		{
			input: `
			const num = 10;
			num = 20;
			`,
			expectedError: `line 3, column 3: attempting to assign value to constant variable 'num'`,
		},
		{
			input: `
			const num = 10;
			fn() {
				num = 20;
				num;
			}
			num;
			`,
			expectedError: `line 4, column 4: attempting to assign value to identifier 'num' prior to declaration`,
		},
		{
			input: `
			const num = 10;
			fn() {
				const num = 20;
				num = 30;
				return num;
			}
			num;
			`,
			expectedError: `line 5, column 4: attempting to assign value to constant variable 'num'`,
		},
		{
			input: `
			const num = 10;
			let num = 20;
			`,
			expectedError: `line 3, column 3: identifier 'num' has already been declared`,
		},
		{
			input: `
			let num = 10;
			const num = 20;
			`,
			expectedError: `line 3, column 3: identifier 'num' has already been declared`,
		},
		{
			input: `
			let num = 10;
			let num = 20;
			`,
			expectedError: `line 3, column 3: identifier 'num' has already been declared`,
		},
		{
			input: `
			const num = 10;
			const num = 20;
			`,
			expectedError: `line 3, column 3: identifier 'num' has already been declared`,
		},
	}

	runCompilerErrorTests(t, tests)
}

func TestGlobalAssignStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			one = "one";
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 0),
			},
			expectedConstants: []interface{}{1, "one"},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two = 2;
			two;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2},
		},
	}

	runCompilerTests(t, tests)
}

func TestAssignStatementsErrors(t *testing.T) {
	tests := []compilerErrorTestCase{
		{
			input: `
			one = 1;
			`,
			expectedError: `line 2, column 3: attempting to assign value to identifier 'one' prior to declaration`,
		},
		{
			input: `
			one += 1;
			`,
			expectedError: `line 2, column 3: attempting to assign value to identifier 'one' prior to declaration`,
		},
		{
			input: `
			let num = 10;
			fn() {
				num = 20;
				num;
			}
			num;
			`,
			expectedError: `line 4, column 4: attempting to assign value to identifier 'num' prior to declaration`,
		},
		{
			input: `
			let num = 10;
			let f = fn() {
				let num = 20;
				let g = fn() {
					num = 30;
					return num
				};
				return g()
			};
			num + f();
			`,
			expectedError: `line 6, column 5: attempting to assign value to identifier 'num' prior to declaration`,
		},
	}

	runCompilerErrorTests(t, tests)
}

func TestOperatorAssignStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let num = 55;
			num += 5;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpSetGlobal, 0),
			},
			expectedConstants: []interface{}{55, 5},
		},
		{
			input: `
			let num = 82;
			num //= 7;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpIntegerDiv),
				bytecode.Make(bytecode.OpSetGlobal, 0),
			},
			expectedConstants: []interface{}{82, 7},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `"monkey"`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{"monkey"},
		},
		{
			input: `"mon" + "key"`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{"mon", "key"},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `[]`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: `[1, 2, 3]`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3},
		},
		{
			input: `[1 + 2, 3 - 4, 5 * 6]`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashMapLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "{}",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpHashMap, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{},
		},
		{
			input: "{1: 2, 3: 4, 5: 6}",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpHashMap, 6),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
		},
		{
			input: "{1: 2 + 3, 4: 5 * 6}",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpHashMap, 4),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "[1, 2, 3][1 + 1]",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
		},
		{
			input: "{1: 2}[2 - 1]",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpHashMap, 2),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{1, 2, 2, 1},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := NewCompiler()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex is wrong. expected=%d, got=%d", 0, compiler.scopeIndex)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(bytecode.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex is wrong. expected=%d, got=%d", 1, compiler.scopeIndex)
	}

	compiler.emit(bytecode.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("length of instructions is wrong. expected=%d, got=%d", 1, len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != bytecode.OpSub {
		t.Errorf("lastInstruction.Opcode is wrong. expected=%d, got=%d", bytecode.OpSub, last.Opcode)
	}

	if compiler.symbolTable.outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbol table when entering scope")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex is wrong. expected=%d, got=%d", 0, compiler.scopeIndex)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table after leaving scope")
	}
	if compiler.symbolTable.outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly - is enclosed")
	}

	compiler.emit(bytecode.OpAdd)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("length of instructions is wrong. expected=%d, got=%d", 2, len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != bytecode.OpAdd {
		t.Errorf("lastInstruction.Opcode is wrong. expected=%d, got=%d", bytecode.OpAdd, last.Opcode)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousLastInstruction
	if previous.Opcode != bytecode.OpMul {
		t.Errorf("lastInstruction.Opcode is wrong. expected=%d, got=%d", bytecode.OpMul, last.Opcode)
	}
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "fn() { }",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpReturn),
				},
			},
		},
		{
			input: "fn() { return 5 + 10; }",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				5,
				10,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: "fn() { 5 + 10; }",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				5,
				10,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: "fn() { 1; 2 }",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				1,
				2,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "fn() { 42; }()",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				42,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			let noArg = fn() { 42 };
			noArg();
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				42,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			let oneArg = fn(a) { a; };
			oneArg(24);
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
				24,
			},
		},
		{
			input: `
			let manyArg = fn(a, b, c) { a; b; c; };
			manyArg(24, 25, 26);
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpCall, 3),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				24,
				25,
				26,
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let num = 55;
			fn() { num; }
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			fn() {
				let num = 55;
				num;
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			fn() {
				let a = 55;
				let b = 77;
				a + b
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				77,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 1),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConstStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			const num = 55;
			fn() { num; }
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			fn() {
				const num = 55;
				num;
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			fn() {
				const a = 55;
				const b = 77;
				a + b
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				77,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 1),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestAssignStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			fn() {
				let num = 55;
				num = num + 5;
				num;
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				5,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			let num = 10;
			fn() {
				let num = 20;
				return num;
			}
			num;
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				10,
				20,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBuiltIns(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			len([]);
			append([], 1);
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpGetBuiltIn, 1),
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetBuiltIn, 5),
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpCall, 2),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				1,
			},
		},
		{
			input: `fn() { len([]); };`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetBuiltIn, 1),
					bytecode.Make(bytecode.OpArray, 0),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			split("hello");
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpGetBuiltIn, 7),
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				"hello",
			},
		},
		{
			input: `
			split("hello world", " ");
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpGetBuiltIn, 7),
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall, 2),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				"hello world",
				" ",
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			fn(a) {
				fn(b) {
					a + b;
				}
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFreeVar, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 0, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			fn(a) {
				fn(b) {
					fn(c) {
						a + b + c;
					}
				}
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFreeVar, 0),
					bytecode.Make(bytecode.OpGetFreeVar, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFreeVar, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 0, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 1, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
		{
			input: `
			let global = 55;

			fn() {
				let a = 66;

				fn() {
					let b = 77;

					fn() {
						let c = 88;
						
						return global + a + b + c;
					}
				}
			}
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 6, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				55,
				66,
				77,
				88,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 3),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpGetFreeVar, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetFreeVar, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 2),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetFreeVar, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 4, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 5, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let countDown = fn(x) { countDown(x - 1); };
			countDown(1);
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpCurrentClosure),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSub),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
				1,
			},
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) { countDown(x - 1); };
				countDown(1);
			};
			wrapper();
			`,
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 3, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpCurrentClosure),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSub),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpClosure, 1, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 2),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
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

func runCompilerErrorTests(t *testing.T, tests []compilerErrorTestCase) {
	t.Helper()

	for _, test := range tests {
		program := parse(test.input)

		compiler := NewCompiler()
		err := compiler.Compile(program)
		if err == nil {
			t.Errorf("no error was detected in compiler in compiler error test case")
		}

		if err.Error() != test.expectedError {
			t.Errorf("compiler error test case: error message was not correct. expected=%q, got=%q", test.expectedError, err.Error())
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
		case float64:
			err := testFloatObject(float64(expConst), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testFloatObject failed: %s", i, err)
			}
		case string:
			err := testStringObject(expConst, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		case []bytecode.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T", i, actual[i])
			}

			err := testInstructions(expConst, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}
		default:
			return fmt.Errorf("unknown expected constant type: %T", expConst)
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

func testFloatObject(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not a Float. got=%T (%+v)", actual, actual)
	}

	if math.Abs(result.Value-expected) > ast.FLOAT_64_EQUALITY_THRESHOLD {
		return fmt.Errorf("object has wrong float value. expected=%f, got=%f", expected, result.Value)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not a String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong string value. expected=%q, got=%q", expected, result.Value)
	}

	return nil
}
