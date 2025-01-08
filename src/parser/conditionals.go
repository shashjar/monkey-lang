package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.currToken}
	ieClauses := []ast.ConditionalClause{}

	ifCondition, ok := p.parseCondition()
	if !ok {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifConsequence := p.parseBlockStatement()

	ieClauses = append(ieClauses, ast.ConditionalClause{Condition: ifCondition, Consequence: ifConsequence})

	for p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) { // Parsing `else if` clause
			p.nextToken()

			elseIfCondition, ok := p.parseCondition()
			if !ok {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			elseIfConsequence := p.parseBlockStatement()

			ieClauses = append(ieClauses, ast.ConditionalClause{Condition: elseIfCondition, Consequence: elseIfConsequence})
		} else { // Parsing `else` clause
			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			ie.Alternative = p.parseBlockStatement()

			break
		}
	}

	ie.Clauses = ieClauses
	return ie
}

func (p *Parser) parseSwitchStatement() ast.Expression {
	ss := &ast.SwitchStatement{Token: p.currToken}
	switchCases := []ast.SwitchCase{}

	p.nextToken()

	switchExpression := p.parseExpression(LOWEST)
	ss.SwitchExpression = switchExpression

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	if !p.expectPeek(token.CASE) {
		return nil
	}

	p.nextToken()

	firstCaseExpression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.COLON) {
		return nil
	}

	firstCaseConsequence := p.parseBlock([]token.TokenType{token.CASE, token.DEFAULT, token.RBRACE})

	firstCase := ast.SwitchCase{Expression: firstCaseExpression, Consequence: firstCaseConsequence}
	switchCases = append(switchCases, firstCase)

	for p.currTokenIs(token.CASE) {
		p.nextToken()

		caseExpression := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		caseConsequence := p.parseBlock([]token.TokenType{token.CASE, token.DEFAULT, token.RBRACE})

		thisCase := ast.SwitchCase{Expression: caseExpression, Consequence: caseConsequence}
		switchCases = append(switchCases, thisCase)
	}

	if p.currTokenIs(token.DEFAULT) {
		if !p.expectPeek(token.COLON) {
			return nil
		}

		defaultConsequence := p.parseBlock([]token.TokenType{token.RBRACE})
		ss.Default = defaultConsequence
	}

	ss.Cases = switchCases
	return ss
}

func (p *Parser) parseCondition() (ast.Expression, bool) {
	if !p.expectPeek(token.LPAREN) {
		return nil, false
	}

	p.nextToken()

	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil, false
	}

	return condition, true
}
