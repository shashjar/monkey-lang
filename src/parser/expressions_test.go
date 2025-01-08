package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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
		{"2**3;", 2, "**", 3},
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
