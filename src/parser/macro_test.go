package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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
