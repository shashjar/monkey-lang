package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestParsingArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, "hello", 6.97]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	array, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.ArrayLiteral. got=%T", statement.Expression)
	}

	if len(array.Elements) != 4 {
		t.Fatalf("len(array.Elements) is wrong. expected=4, got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testStringLiteral(t, array.Elements[2], "hello")
	testFloat(t, array.Elements[3], 6.97)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := `myArray[1 * 3]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	indexExp, ok := statement.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.IndexExpression. got=%T", statement.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "*", 3) {
		return
	}
}

func TestParsingEmptyHashMapLiteral(t *testing.T) {
	input := `{}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hashmap, ok := statement.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.HashMapLiteral. got=%T", statement.Expression)
	}

	if len(hashmap.KVPairs) != 0 {
		t.Fatalf("len(hashmap.KVPairs) is wrong. expected=0, got=%d", len(hashmap.KVPairs))
	}
}

func TestParsingHashMapLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hashmap, ok := statement.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.HashMapLiteral. got=%T", statement.Expression)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hashmap.KVPairs) != len(expected) {
		t.Fatalf("len(hashmap.KVPairs) is wrong. expected=%d, got=%d", len(expected), len(hashmap.KVPairs))
	}

	for k, v := range hashmap.KVPairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not an ast.StringLiteral. got=%T", k)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hashmap, ok := statement.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.HashMapLiteral. got=%T", statement.Expression)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hashmap.KVPairs) != len(expected) {
		t.Fatalf("len(hashmap.KVPairs) is wrong. expected=%d, got=%d", len(expected), len(hashmap.KVPairs))
	}

	for k, v := range hashmap.KVPairs {
		boolean, ok := k.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not an ast.BooleanLiteral. got=%T", k)
			continue
		}

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hashmap, ok := statement.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.HashMapLiteral. got=%T", statement.Expression)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hashmap.KVPairs) != len(expected) {
		t.Fatalf("len(hashmap.KVPairs) is wrong. expected=%d, got=%d", len(expected), len(hashmap.KVPairs))
	}

	for k, v := range hashmap.KVPairs {
		integer, ok := k.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not an ast.IntegerLiteral. got=%T", k)
			continue
		}

		expectedValue := expected[integer.String()]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	hashmap, ok := statement.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.HashMapLiteral. got=%T", statement.Expression)
	}

	if len(hashmap.KVPairs) != 3 {
		t.Fatalf("len(hashmap.KVPairs) is wrong. expected=%d, got=%d", 3, len(hashmap.KVPairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for k, v := range hashmap.KVPairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not an ast.StringLiteral. got=%T", k)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(v)
	}
}
