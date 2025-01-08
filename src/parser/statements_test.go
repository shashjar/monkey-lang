package parser

import (
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

func TestOperatorAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		left               interface{}
		operator           string
		right              interface{}
	}{
		{"x += 5;", "x", "x", "+", 5},
		{"y -= 3;", "y", "y", "-", 3},
		{"a *= 4;", "a", "a", "*", 4},
		{"b /= 2;", "b", "b", "/", 2},
		{"z //= 7;", "z", "z", "//", 7},
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
		if !testInfixExpression(t, val, test.left, test.operator, test.right) {
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
