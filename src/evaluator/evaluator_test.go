package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"5", 5},
		{"10", 10},
		{"16;", 16},
		{"-3", -3},
		{"-7;", -7},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 3 * 2", 24},
		{"-50 - 50 + 100", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 - 4) * 10", -30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"true;", true},
		{"false", false},
		{"false;", false},

		{"true == true", true},
		{"true == false", false},
		{"false == true", false},
		{"false == false", true},

		{"true != true", false},
		{"true != false", true},
		{"false != true", true},
		{"false != false", false},

		{"(2 < 5) == true;", true},
		{"(3 == -1) != false;", false},

		{"1 == 1", true},
		{"1 == 2", false},
		{"-4 == 0", false},
		{"7 != 7;", false},
		{"6 != -1", true},
		{"1 < 2", true},
		{"16 < 10;", false},
		{"4 > 3;", true},
		{"-9 > 72", false},

		{`"" == ""`, true},
		{`"foo" == "foo"`, true},
		{`"foo" == "bar"`, false},
		{`"bar" != "bar"`, false},
		{`"baz" != "baz "`, true},
		{`"baz" != "foo"`, true},
		{`let s = "hello " + "there" + "!"; s == "hello there!"`, true},
		{`let s = "hello " + "there" + "!"; s == "hi there!;"`, false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world!"`, "hello world!"},
		{`"hello world!";`, "hello world!"},
		{`let s = "3 + 9;"; s;`, "3 + 9;"},
		{`let s = "foo" + "bar"; s;`, "foobar"},
		{`let s = "foo " + "bar " + "baz"; s;`, "foo bar baz"},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testStringObject(t, evaluated, test.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}

			return 1;
		}
		`, 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
		{
			"let a = 5; let b = a + 7; c;",
			"identifier not found: c",
		},
		{
			"fn(x, y) { return 10 * x + y - 3; }(10, 6, 7)",
			"wrong number of arguments provided to function. expected=2, received=3",
		},
		{
			"fn(x, y, z) { return 10 * x + y - 3; }(10, 6)",
			"wrong number of arguments provided to function. expected=3, received=2",
		},
		{
			`5; "true" + "false"; "true" + 5;`,
			"type mismatch: STRING + INTEGER",
		},
		{
			`"hello" - "there";`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testErrorObject(t, evaluated, test.expectedMessage)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a  = 5; a;", 5},
		{"let a  = 5 * 5; a;", 25},
		{"let a  = 5; let b = a; b;", 5},
		{"let a  = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters. expected=%d, got=%d", 1, len(fn.Parameters))
	}

	if fn.Parameters[0].Value != "x" {
		t.Fatalf("function parameter is not 'x'. got=%q", fn.Parameters[0].Value)
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("function body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let add = fn(x, y) { return x + y; }; add(4, 3);", 7},
		{"let double = fn(x) { x * 2; }; double(8);", 16},
		{"let sub = fn(x, y) { return x - y; }; sub(5, -3);", 8},
		{"let sub = fn(x, y) { return x - y; }; sub(sub(4, 1), sub(9, 20));", 14},
		{"fn(x) { x; }(5)", 5},
		{"fn(x, y) { return 10 * x + y - 3; }(10, 6)", 103},
	}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y; }
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("hello")`, 5},
		{`len("Hello world!")`, 12},
		{`let a = "hi"; let b = " there"; len(a + b)`, 8},

		{`len([])`, 0},
		{`len([1, 2])`, 2},
		{`len([-4, "hello world!", 65, fn(x) { 2 * x }])`, 4},

		{`len(1)`, "argument to `len` is not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. expected=1, got=2"},

		{`first([])`, "array is empty; no first element"},
		{`first([3])`, 3},
		{`first([-8, 7, -4, 20])`, -8},
		{`first("hello world")`, "argument to `first` is not supported, got STRING"},

		{`last([])`, "array is empty; no last element"},
		{`last([3])`, 3},
		{`last([-8, 7, -4, 20])`, 20},
		{`last("hello world")`, "argument to `last` is not supported, got STRING"},

		{`rest([])`, nil},
		{`rest([3])`, []int{}},
		{`rest([-8, 7, -4, 20])`, []int{7, -4, 20}},
		{`rest("hello world")`, "argument to `rest` is not supported, got STRING"},

		{`append([], 3)`, []int{3}},
		{`append([3], 21)`, []int{3, 21}},
		{`append([4, -10])`, "wrong number of arguments. expected=2, got=1"},
		{`append("hello world", "hi")`, "argument to `append` is not supported, got STRING"},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)

		switch expected := test.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testErrorObject(t, evaluated, expected)
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not an Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong number of elements. expected=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElement := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElement))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3, "hi there"]`
	evaluated := testEval(input)
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not an Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 4 {
		t.Fatalf("array has wrong number of elements. expected=%d, got=%d", 4, len(array.Elements))
	}

	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testIntegerObject(t, array.Elements[2], 6)
	testStringObject(t, array.Elements[3], "hi there")
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"3[3]",
			"index operator not supported: INTEGER[INTEGER]",
		},
		{
			`[1, 2, 3][fn(x) { 2 * x }]`,
			"index operator not supported: ARRAY[FUNCTION]",
		},
		{
			`let a = "s"; [1, 2, 3][a]`,
			"index operator not supported: ARRAY[STRING]",
		},
		{
			"[1, 2, 3][3]",
			"index is out-of-bounds for array",
		},
		{
			"[1, 2, 3][-1]",
			"index is out-of-bounds for array",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}
}

func TestHashMapLiterals(t *testing.T) {
	input := `
	let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: fn(x) { x + 4 }(1),
		false: 6,
	}
	`
	evaluated := testEval(input)
	hashmap, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("object is not a HashMap. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(hashmap.KVPairs) != len(expected) {
		t.Fatalf("hashmap has wrong number of pairs. expected=%d, got=%d", len(expected), len(hashmap.KVPairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hashmap.KVPairs[expectedKey]
		if !ok {
			t.Errorf("no pair found for key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashMapIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{1: "o" + "n" + "e", "two": 2, true: 3}[1]`,
			"one",
		},
		{
			`{1: "one", "two": 3 * 5 - 13, true: 3}["two"]`,
			2,
		},
		{
			`{1: "one", "two": 2, true: 3}[true != false]`,
			3,
		},
		{
			`let m = {"hi": "there", true: false, false: true, 30 * 10: 40 - 20}; m[10 == 10]`,
			false,
		},
		{
			`{-4: 4}[-4]`,
			4,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestHashMapIndexExpressionErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`{fn(x) { x }: "Monkey"}[true]`,
			"key is unusable as hash key: FUNCTION",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }]`,
			"index operator not supported: HASHMAP[FUNCTION]",
		},
		{
			`{}["foo"]`,
			"key not found in hashmap",
		},
		{
			`{1: "one", "two": 2, true: 3}[0]`,
			"key not found in hashmap",
		},
		{
			`{1: "one", "two": 2, true: 3}[false]`,
			"key not found in hashmap",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testErrorObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not a Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%t, got=%t", expected, result.Value)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not a String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%q, got=%q", expected, result.Value)
		return false
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	errorObj, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not an Error. got=%T (%+v)", obj, obj)
		return false
	}
	if errorObj.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q", expected, errorObj.Message)
		return false
	}

	return true
}
