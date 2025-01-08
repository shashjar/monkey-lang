package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseWhileLoop() ast.Expression {
	wl := &ast.WhileLoop{Token: p.currToken}

	whileCondition, ok := p.parseCondition()
	if !ok {
		return nil
	}
	wl.Condition = whileCondition

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	whileBody := p.parseBlockStatement()
	wl.Body = whileBody

	return wl
}

func (p *Parser) parseForLoop() ast.Expression {
	fl := &ast.ForLoop{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	init := p.parseStatement()
	if !p.currTokenIs(token.SEMICOLON) {
		return nil
	}
	p.nextToken()

	condition := p.parseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()

	afterthought := p.parseStatement()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	forBody := p.parseBlockStatement()

	fl.Init = init
	fl.Condition = condition
	fl.Afterthought = afterthought
	fl.Body = forBody
	return fl
}
