package vm

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1.49", 1.49},
		{"-7801.23189123", -7801.23189123},
		{"-4", -4},
		{"1 + 2", 3},
		{"4 - 9", -5},
		{"0 * 6", 0},
		{"0 * -6", 0},
		{"-2 * 4", -8},
		{"5 * 12", 60},
		{"21 / 7", 3},
		{"22 / 7", 3.1428571429},
		{"0 % 10", 0},
		{"36 % 6", 0},
		{"7 % 4", 3},
		{"9 % -4", 1},
		{"-9 % 4", -1},
		{"-9 % -4", -1},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"3 * -(2 + 10) + 4", -32},
		{"3 * -(2 + 10) + 4 % 3", -35},
		{"7 * 8.4178923", 58.9252461},
		{"6.84 - 10 * 9.8 + 42.17987", -48.98013},
		{"5 * (6.87 - 10.892137) / 4.21 + (10 - 7) * 11.6", 30.0231152019},
	}

	runVMTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
		{"1 > 2", false},
		{"1 < 2", true},
		{"1 < 1", false},
		{"3 <= 2", false},
		{"3 <= 3", true},
		{"3 <= 4", true},
		{"5 >= 6", false},
		{"5 >= 5", true},
		{"3 == 3", true},
		{"4 != 3", true},
		{"2 != 2", false},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != true", false},
		{"true != false", true},
		{"false != false", false},
		{`"" == ""`, true},
		{`"d" == ""`, false},
		{`"hello" == "hello"`, true},
		{`"hello" != "world"`, true},
		{`"ten" != "ten"`, false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!(if (false) { 5; })", true},
		{"true && false", false},
		{"false || true", true},
		{"false || false", false},
		{"(6 < 8) && true", true},
		{"(8 < 6) || (9 + 9 < 20)", true},
		{`"world" + "s" == "worlds"`, true},
	}

	runVMTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (3 > 5) { 10 } else if (3 >= 3) { 20 } else { 30 }", 20},
		{"if (3 > 5) { 10 } else if (3 <= 2) { 20 }", Null},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
		{"if (true) { 10 } else if (1 < 2) { 20 } else { 30 }", 10},
		{"if (3 > 4) { 10 } else if (1 < 2) { 20 } else { 30 }", 20},
		{"if (1 > 2) { 10 } else if (1 == 2) { 20 } else { 30 }", 30},
		{"if (1 > 2) { 10 } else if (1 == 2) { 20 }", Null},
		{"if (1 > 2) { 10 } else if (1 < 2) { 20 }", 20},
		{"if ((if (false) { false } else { true })) { 5 } else { 7 }", 5},
		{"if (3) { 10 } else if (4) { 20 } else { 30 }", 10},
		{"if (81.917) { 10 } else if (4) { 20 } else { 30 }", 10},
	}

	runVMTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let three = one + one + 1; one + three;", 4},
	}

	runVMTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + " banana"`, "monkey banana"},
	}

	runVMTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 - 6]", []int{3, 12, -1}},
		{"[] + []", []int{}},
		{"[1] + []", []int{1}},
		{"[] + [-7]", []int{-7}},
		{"[4, 0, -5] + [6, 11]", []int{4, 0, -5, 6, 11}},
	}

	runVMTests(t, tests)
}

func TestHashMapLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1: 2, 2: 3}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 2}).HashKey(): 3,
		}},
		{"{1 + 1: 2 * 2, 3 + 3: 4 * 4}", map[object.HashKey]int64{
			(&object.Integer{Value: 2}).HashKey(): 4,
			(&object.Integer{Value: 6}).HashKey(): 16,
		}},
	}

	runVMTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[4, 5, 6]][0][0]", 4},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() { };
			noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() { };
			let noReturnTwo = fn() { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
			expected: Null,
		},
		{
			input: `
			let fivePlusTen = fn() { 5 + 10; };
			fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
			let one = fn() { 1; };
			let two = fn() { 2; };
			one() + two();
			`,
			expected: 3,
		},
		{
			input: `
			let a = fn() { 1; };
			let b = fn() { a() + 1; };
			let c = fn() { b() + 1 };
			c();
			`,
			expected: 3,
		},
		{
			input: `
			let earlyExit = fn() { return 99; 100; };
			earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit = fn() { return 99; return 100; };
			earlyExit();
			`,
			expected: 99,
		},
	}

	runVMTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() { 1; };
			let returnsOneReturner = fn() { returnsOne; };
			returnsOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;
			};
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num;
			}
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let identity = fn(a) { a; };
			identity(42);
			`,
			expected: 42,
		},
		{
			input: `
			let sum = fn(a, b) { a + b; };
			sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				return c;
			};
			sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				return c;
			};
			sum(1, 2) + sum(3, 4);
			`,
			expected: 10,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				c;
			};
			let outer = fn() {
				sum(1, 2) + sum(3, 4);
			};
			outer();
			`,
			expected: 10,
		},
		{
			input: `
			let globalNum = 10;

			let sum = fn(a, b) {
				let c = a + b;
				return c + globalNum;
			};

			let outer = fn() {
				sum(1, 2) + sum(3, 4) + globalNum;
			};

			outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let identity = fn(a) { a; };
			identity(42, 5);
			`,
			expected: `wrong number of arguments: expected=1, got=2`,
		},
		{
			input:    `fn() { 1; }(1)`,
			expected: `wrong number of arguments: expected=0, got=1`,
		},
		{
			input:    `fn(a) { a; }()`,
			expected: `wrong number of arguments: expected=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1)`,
			expected: `wrong number of arguments: expected=2, got=1`,
		},
	}

	for _, test := range tests {
		program := parse(test.input)

		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := NewVM(compiler.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but didn't receive one")
		}

		if err.Error() != test.expected {
			t.Fatalf("wrong VM error: expected=%q, got=%q", test.expected, err)
		}
	}
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
				fn() { a; };
			};
			let closure = newClosure(99);
			closure()
			`,
			expected: 99,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				fn(c) { a + b + c };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				let c = a + b;
				fn(d) { c + d; };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
			let newAdderOuter = fn(a, b) {
				let c = a + b;
				fn(d) {
					let e = d + c;
					fn(f) { e + f; };
				};
			};
			let newAdderInner = newAdderOuter(1, 2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let a = 1;
			let newAdderOuter = fn(b) {
				fn(c) {
					fn(d) { a + b + c + d };
				};
			};
			let newAdderInner = newAdderOuter(2)
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let newClosure = fn(a, b) {
				let one = fn() { a; };
				let two = fn() { b; };
				fn() { one() + two(); };
			};
			let closure = newClosure(9, 90);
			closure();
			`,
			expected: 99,
		},
	}

	runVMTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			let wrapper = fn() {
				countDown(1);
			};
			wrapper();
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(1);
			};
			wrapper();
			`,
			expected: 0,
		},
	}

	runVMTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fibonacci = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					if (x == 1) {
						return 1;
					} else {
						fibonacci(x - 1) + fibonacci(x - 2);
					}
				}
			};
			fibonacci(15);
			`,
			expected: 610,
		},
	}

	runVMTests(t, tests)
}

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, test := range tests {
		program := parse(test.input)

		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := NewVM(compiler.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("VM error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, test.expected, stackElem)
	}
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case float64:
		err := testFloatObject(float64(expected), actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not an Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("array has wrong number of elements. expected=%d, got=%d", len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hashmap, ok := actual.(*object.HashMap)
		if !ok {
			t.Errorf("object is not a HashMap. got=%T (%+v)", actual, actual)
			return
		}

		if len(hashmap.KVPairs) != len(expected) {
			t.Errorf("hashmap has wrong number of elements. expected=%d, got=%d", len(expected), len(hashmap.KVPairs))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hashmap.KVPairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs: %q", expectedKey)
			}

			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not an Error. got=%T (%+v)", actual, actual)
			return
		}

		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q", expected.Message, errObj.Message)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	default:
		t.Errorf("unknown expected object type: %T", expected)
	}
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

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not a Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong boolean value. expected=%t, got=%t", expected, result.Value)
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
