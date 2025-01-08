package parser

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is not an *ast.IntegerLiteral. got=%T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value is not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral is not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestFloatExpression(t *testing.T) {
	input := "8.946;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	float, ok := statement.Expression.(*ast.Float)
	if !ok {
		t.Fatalf("expression is not an *ast.Float. got=%T", statement.Expression)
	}

	if math.Abs(float.Value-8.946) > ast.FLOAT_64_EQUALITY_THRESHOLD {
		t.Errorf("float.Value is not %f. got=%f", 8.946, float.Value)
	}
	if float.TokenLiteral() != "8.946" {
		t.Errorf("float.TokenLiteral is not %s. got=%s", "8.946", float.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
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

		boolean, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp is not an *ast.Boolean. got=%T", statement.Expression)
		}

		if boolean.Value != test.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", test.expectedBoolean, boolean.Value)
		}
	}
}

func TestStringExpression(t *testing.T) {
	input := `"Hello world!"`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expression is not an *ast.String. got=%T", statement.Expression)
	}

	if literal.Value != "Hello world!" {
		t.Errorf("literal.Value is not %q. got=%q", "Hello world!", literal.Value)
	}
	if literal.TokenLiteral() != "Hello world!" {
		t.Errorf("literal.TokenLiteral is not %q. got=%q", "Hello world!", literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not an *ast.Identifier. got=%T", statement.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloat(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not an *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral is not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testFloat(t *testing.T, f ast.Expression, value float64) bool {
	float, ok := f.(*ast.Float)
	if !ok {
		t.Errorf("f is not an *ast.Float. got=%T", f)
		return false
	}

	if math.Abs(float.Value-value) > ast.FLOAT_64_EQUALITY_THRESHOLD {
		t.Errorf("float.Value is not %f. got=%f", value, float.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not an *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value is not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral is not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	sl, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp is not an *ast.StringLiteral. got=%T", sl)
		return false
	}

	if sl.Value != value {
		t.Errorf("sl.Value is not %q. got=%q", value, sl.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not an *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}
