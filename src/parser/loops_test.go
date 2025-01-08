package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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
