package parser

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}

		val := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, test.expectedValue) {
			return
		}
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"const X = 5;", "X", 5},
		{"const Y = true;", "Y", true},
		{"const FOOBAR = Y;", "FOOBAR", "Y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		statement := program.Statements[0]
		if !testConstStatement(t, statement, test.expectedIdentifier) {
			return
		}

		val := statement.(*ast.ConstStatement).Value
		if !testLiteralExpression(t, val, test.expectedValue) {
			return
		}
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		statement := program.Statements[0]
		if !testAssignStatement(t, statement, test.expectedIdentifier) {
			return
		}

		val := statement.(*ast.AssignStatement).Value
		if !testLiteralExpression(t, val, test.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement is not an *ast.ReturnStatement. got=%T", statement)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Fatalf("returnStatement.TokenLiteral is not 'return', got %q", returnStatement.TokenLiteral())
		}
		if testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedValue) {
			return
		}
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range prefixTests {
		l := lexer.NewLexer(test.input)
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

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression is not an *ast.PrefixExpression. got=%T", statement.Expression)
		}

		if exp.Operator != test.operator {
			t.Fatalf("exp.Operator is not '%s'. got='%s'", test.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, test.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"10 - 4;", 10, "-", 4},
		{"3 * 9;", 3, "*", 9},
		{"12 / 6;", 12, "/", 6},
		{"22 // 7;", 22, "//", 7},
		{"6 > 0;", 6, ">", 0},
		{"3 < 7;", 3, "<", 7},
		{"3 >= 2;", 3, ">=", 2},
		{"15 <= 21;", 15, "<=", 21},
		{"5 == 5;", 5, "==", 5},
		{"9 != 10;", 9, "!=", 10},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"true && true;", true, "&&", true},
		{"true && false;", true, "&&", false},
		{"false && true;", false, "&&", true},
		{"false && false;", false, "&&", false},
		{"true || true;", true, "||", true},
		{"true || false;", true, "||", false},
		{"false || true;", false, "||", true},
		{"false || false;", false, "||", false},
		{"5 % 2;", 5, "%", 2},
	}

	for _, test := range infixTests {
		l := lexer.NewLexer(test.input)
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

		if !testInfixExpression(t, statement.Expression, test.leftValue, test.operator, test.rightValue) {
			return
		}
	}
}

func TestParsingPostfixExpressions(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
	}{
		{"x++", "x"},
		{"y--", "y"},
		{"x++;", "x"},
		{"y--;", "y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		statement := program.Statements[0]
		if !testAssignStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

func TestPostfixOperatorErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			input:         "3++",
			expectedError: "expected postfix operator '++' to be applied to an identifier. got 3 instead",
		},
		{
			input:         "true--",
			expectedError: "expected postfix operator '--' to be applied to an identifier. got true instead",
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		p.ParseProgram()

		errorCreated := false
		for _, err := range p.errors {
			if err == test.expectedError {
				errorCreated = true
			}
		}

		if !errorCreated {
			t.Fatalf("parser did not create expected error. expected=%q, got=%q", test.expectedError, p.errors)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"-6 + 7 * 8 - 4 % 9 / 10",
			"(((-6) + (7 * 8)) - ((4 % 9) / 10))",
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

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

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, exp.Clauses[0].Condition, "x", "<", "y") {
		return
	}

	if len(exp.Clauses[0].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Clauses[0].Consequence.Statements))
	}

	consequence, ok := exp.Clauses[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Clauses[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

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

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, exp.Clauses[0].Condition, "x", "<", "y") {
		return
	}

	if len(exp.Clauses[0].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Clauses[0].Consequence.Statements))
	}

	consequence, ok := exp.Clauses[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Clauses[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfElseIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else if (x == y) { 0 } else { y }`

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

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.IfExpression. got=%T", statement.Expression)
	}

	if len(exp.Clauses) != 2 {
		t.Fatalf("exp.Clauses is the wrong length. expected=%d, got=%d", 2, len(exp.Clauses))
	}

	// First clause (`if`)
	if !testInfixExpression(t, exp.Clauses[0].Condition, "x", "<", "y") {
		return
	}

	if len(exp.Clauses[0].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Clauses[0].Consequence.Statements))
	}

	consequence, ok := exp.Clauses[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Clauses[0].Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Clauses[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Second clause (`else if`)
	if !testInfixExpression(t, exp.Clauses[1].Condition, "x", "==", "y") {
		return
	}

	if len(exp.Clauses[1].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Clauses[1].Consequence.Statements))
	}

	consequence, ok = exp.Clauses[1].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Clauses[1].Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Clauses[1].Consequence.Statements[0])
	}

	if !testIntegerLiteral(t, consequence.Expression, 0) {
		return
	}

	// Alternative (`else`)
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestBasicSwitchStatement(t *testing.T) {
	input := `
	switch x {
	case "hello":
		x;
	}
	`

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

	exp, ok := statement.Expression.(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.SwitchStatement. got=%T", statement.Expression)
	}

	// Switch expression
	if !testLiteralExpression(t, exp.SwitchExpression, "x") {
		return
	}

	if len(exp.Cases) != 1 {
		t.Fatalf("exp.Cases is the wrong length. expected=%d, got=%d", 1, len(exp.Cases))
	}

	// First case ("hello")
	if !testStringLiteral(t, exp.Cases[0].Expression, "hello") {
		return
	}

	if len(exp.Cases[0].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Cases[0].Consequence.Statements))
	}

	consequence, ok := exp.Cases[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Cases[0].Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Cases[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}

func TestSwitchStatementWithDefault(t *testing.T) {
	input := `
	switch x {
	case "hello":
		x;
	case "world":
		3;
	default:
		y;
	}
	`

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

	exp, ok := statement.Expression.(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.SwitchStatement. got=%T", statement.Expression)
	}

	// Switch expression
	if !testLiteralExpression(t, exp.SwitchExpression, "x") {
		return
	}

	if len(exp.Cases) != 2 {
		t.Fatalf("exp.Cases is the wrong length. expected=%d, got=%d", 2, len(exp.Cases))
	}

	// First case ("hello")
	if !testStringLiteral(t, exp.Cases[0].Expression, "hello") {
		return
	}

	if len(exp.Cases[0].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Cases[0].Consequence.Statements))
	}

	consequence, ok := exp.Cases[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Cases[0].Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Cases[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Second case ("world")
	if !testStringLiteral(t, exp.Cases[1].Expression, "world") {
		return
	}

	if len(exp.Cases[1].Consequence.Statements) != 1 {
		t.Errorf("consequence contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Cases[1].Consequence.Statements))
	}

	consequence, ok = exp.Cases[1].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Cases[1].Consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Cases[1].Consequence.Statements[0])
	}

	if !testIntegerLiteral(t, consequence.Expression, 3) {
		return
	}

	// Default case
	if len(exp.Default.Statements) != 1 {
		t.Errorf("default contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Default.Statements))
	}

	defaultCase, ok := exp.Default.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Default.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.Default.Statements[0])
	}

	if !testIdentifier(t, defaultCase.Expression, "y") {
		return
	}
}

func TestWhileLoop(t *testing.T) {
	input := `while (x < y) { x = x + 1; }`

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

	exp, ok := statement.Expression.(*ast.WhileLoop)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.WhileLoop. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("body contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Body.Statements))
	}

	body, ok := exp.Body.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("exp.Body.Statements[0] is not an ast.AssignStatement. got=%T", exp.Body.Statements[0])
	}

	if !testAssignStatement(t, body, "x") {
		return
	}
}

func TestForLoop(t *testing.T) {
	input := `for (let i = 0; i < len(arr); i = i + 1) { puts(arr[i]); }`

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

	exp, ok := statement.Expression.(*ast.ForLoop)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.ForLoop. got=%T", statement.Expression)
	}

	if !testLetStatement(t, exp.Init, "i") {
		return
	}

	condition, ok := exp.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp.Condition is not an ast.InfixExpression")
	}

	testLiteralExpression(t, condition.Left, "i")
	lenCall, ok := condition.Right.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Condition.Right is not an ast.CallExpression")
	}

	if !testIdentifier(t, lenCall.Function, "len") {
		return
	}

	if len(lenCall.Arguments) != 1 {
		t.Fatalf("call expression has wrong number of arguments. expected=%d, got=%d", 1, len(lenCall.Arguments))
	}

	if !testAssignStatement(t, exp.Afterthought, "i") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("body contains wrong number of statements. expected=%d, got=%d", 1, len(exp.Body.Statements))
	}

	body, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Body.Statements[0] is not an *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ce, ok := body.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("body.Expression is not an ast.CallExpression. got=%T", statement.Expression)
	}

	if !testIdentifier(t, ce.Function, "puts") {
		return
	}

	if len(ce.Arguments) != 1 {
		t.Fatalf("call expression has wrong number of arguments. expected=%d, got=%d", 1, len(ce.Arguments))
	}

	indexExp, ok := ce.Arguments[0].(*ast.IndexExpression)
	if !ok {
		t.Fatalf("ce.Arguments[0] is not an ast.IndexExpression. got=%T", statement.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "arr") {
		return
	}

	if !testIdentifier(t, indexExp.Index, "i") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

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

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.FunctionLiteral. got=%T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. expected 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body contains wrong number of statements. expected=%d, got=%d\n", 1, len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not an ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn() { }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.LetStatement. got=%T", program.Statements[0])
	}

	function, ok := statement.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Value is not an ast.FunctionLiteral. got=%T", statement.Value)
	}

	if function.Name != "myFunction" {
		t.Fatalf("function literal name was wrong. expected='myFunction', got=%q", function.Name)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParams) {
			t.Errorf("function parameters are of the wrong length. expected=%d, got=%d\n", len(test.expectedParams), len(function.Parameters))
		}

		for i, ident := range test.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

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

	exp, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.CallExpression. got=%T", statement.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("call expression has wrong number of arguments. expected=%d, got=%d", 3, len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := statement.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("statement.Expression is not an ast.CallExpression. got=%T", statement.Expression)
		}

		if !testIdentifier(t, exp.Function, test.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(test.expectedArgs) {
			t.Fatalf("call expression has wrong number of arguments. expected=%d, got=%d", len(test.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range test.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. expected=%q, got=%q", i, arg, exp.Arguments[i].String())
			}
		}
	}
}

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

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; }`
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

	macro, ok := statement.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not an ast.MacroLiteral. got=%T", statement.Expression)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("number of macro literal parameters is wrong. expected=2, got=%d", len(macro.Parameters))
	}
	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")

	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro.Body.Statements has wrong number of statements. expected=1, got=%d", len(macro.Body.Statements))
	}
	bodyStatement, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body statement is not an ast.ExpressionStatement. got=%T", macro.Body.Statements[0])
	}
	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not an *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value is not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() is not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func testConstStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "const" {
		t.Errorf("statement.TokenLiteral not 'const'. got=%q", statement.TokenLiteral())
	}

	constStatement, ok := statement.(*ast.ConstStatement)
	if !ok {
		t.Errorf("statement is not an *ast.constStatement. got=%T", statement)
		return false
	}

	if constStatement.Name.Value != name {
		t.Errorf("constStatement.Name.Value is not '%s'. got=%s", name, constStatement.Name.Value)
		return false
	}

	if constStatement.Name.TokenLiteral() != name {
		t.Errorf("constStatement.Name.TokenLiteral() is not '%s'. got=%s", name, constStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func testAssignStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != name {
		t.Errorf("statement.TokenLiteral not %q. got=%q", name, statement.TokenLiteral())
	}

	assignStatement, ok := statement.(*ast.AssignStatement)
	if !ok {
		t.Errorf("statement is not an *ast.AssignStatement. got=%T", statement)
		return false
	}

	if assignStatement.Name.Value != name {
		t.Errorf("assignStatement.Name.Value is not '%s'. got=%s", name, assignStatement.Name.Value)
		return false
	}

	if assignStatement.Name.TokenLiteral() != name {
		t.Errorf("assignStatement.Name.TokenLiteral() is not '%s'. got=%s", name, assignStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not an ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloat(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
