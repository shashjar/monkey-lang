package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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
