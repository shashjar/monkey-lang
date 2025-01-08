package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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
