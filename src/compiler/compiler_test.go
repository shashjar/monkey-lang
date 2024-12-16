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
			input: "-2",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{2},
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
			input: "15 / 5",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{15, 5},
		},
	}

	runCompilerTests(t, tests)
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
			input: "6 < 8",
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
			expectedConstants: []interface{}{8, 6},
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
				bytecode.Make(bytecode.OpConstant, 0),
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
				bytecode.Make(bytecode.OpConstant, 2),
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
				bytecode.Make(bytecode.OpConstant, 2),
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
				bytecode.Make(bytecode.OpConstant, 2),
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
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall),
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
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall),
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
				bytecode.Make(bytecode.OpConstant, 1),
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
				bytecode.Make(bytecode.OpConstant, 1),
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
				bytecode.Make(bytecode.OpConstant, 2),
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
