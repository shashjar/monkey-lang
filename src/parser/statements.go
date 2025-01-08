package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch {
	case p.currToken.Type == token.LET:
		return p.parseBindingDeclarationStatement(false)
	case p.currToken.Type == token.CONST:
		return p.parseBindingDeclarationStatement(true)
	case p.currToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN):
		return p.parseAssignStatement()
	case p.currToken.Type == token.IDENT && p.peekTokenIn(token.OPERATOR_ASSIGNMENTS):
		return p.parseOperatorAssignStatement()
	case p.peekTokenIs(token.INCREMENT) || p.peekTokenIs(token.DECREMENT):
		return p.parsePostfixStatement()
	case p.currToken.Type == token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBindingDeclarationStatement(isConst bool) ast.Statement {
	var statement ast.Statement
	if isConst {
		statement = &ast.ConstStatement{Token: p.currToken}
	} else {
		statement = &ast.LetStatement{Token: p.currToken}
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	name := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	value := p.parseExpression(LOWEST)

	if fl, ok := value.(*ast.FunctionLiteral); ok {
		fl.Name = name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	switch statement := statement.(type) {
	case *ast.LetStatement:
		statement.Name = name
		statement.Value = value
	case *ast.ConstStatement:
		statement.Name = name
		statement.Value = value
	}

	return statement
}

func (p *Parser) parseAssignStatement() ast.Statement {
	if !p.currTokenIs(token.IDENT) {
		return nil
	}

	assignStatement := &ast.AssignStatement{Token: p.currToken}
	assignStatement.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	assignStatement.Value = p.parseExpression(LOWEST)

	if fl, ok := assignStatement.Value.(*ast.FunctionLiteral); ok {
		fl.Name = assignStatement.Name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return assignStatement
}

func (p *Parser) parseOperatorAssignStatement() ast.Statement {
	if !p.currTokenIs(token.IDENT) {
		return nil
	}

	assignStatement := &ast.AssignStatement{Token: p.currToken}
	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	assignStatement.Name = identifier

	p.nextToken()

	operatorAssignmentToken := p.currToken

	p.nextToken()

	rightExpression := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	var operatorTok token.Token
	switch operatorAssignmentToken.Type {
	case token.PLUS_ASSIGN:
		operatorTok = token.Token{Type: token.PLUS, Literal: "+"}
	case token.MINUS_ASSIGN:
		operatorTok = token.Token{Type: token.MINUS, Literal: "-"}
	case token.MUL_ASSIGN:
		operatorTok = token.Token{Type: token.MUL, Literal: "*"}
	case token.DIV_ASSIGN:
		operatorTok = token.Token{Type: token.DIV, Literal: "/"}
	case token.INTEGER_DIV_ASSIGN:
		operatorTok = token.Token{Type: token.INTEGER_DIV, Literal: "//"}
	default:
		msg := fmt.Sprintf("received invalid token for operator assignment statement: %s", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	valueExpression := ast.InfixExpression{Token: operatorTok, Left: identifier, Operator: operatorTok.Literal, Right: rightExpression}
	assignStatement.Value = &valueExpression

	return assignStatement
}

func (p *Parser) parsePostfixStatement() ast.Statement {
	if !p.currTokenIs(token.IDENT) {
		msg := fmt.Sprintf("expected postfix operator '%s' to be applied to an identifier. got %s instead", p.peekToken.Literal, p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	postfixStatement := &ast.AssignStatement{Token: p.currToken}
	postfixStatement.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	p.nextToken()

	switch p.currToken.Type {
	case token.INCREMENT:
		postfixStatement.Value = &ast.InfixExpression{
			Token:    token.Token{Type: token.PLUS, Literal: "+"},
			Left:     postfixStatement.Name,
			Operator: "+",
			Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
		}
	case token.DECREMENT:
		postfixStatement.Value = &ast.InfixExpression{
			Token:    token.Token{Type: token.MINUS, Literal: "-"},
			Left:     postfixStatement.Name,
			Operator: "-",
			Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
		}
	default:
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return postfixStatement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	returnStatement := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	returnStatement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStatement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	return p.parseBlock([]token.TokenType{token.RBRACE})
}

func (p *Parser) parseBlock(endingTokenTypes []token.TokenType) *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{Token: p.currToken}
	blockStatement.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currTokenIn(endingTokenTypes) && !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		p.nextToken()
	}

	return blockStatement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{Token: p.currToken}

	statement.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
