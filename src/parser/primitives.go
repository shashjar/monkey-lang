package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as an integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	literal.Value = value

	return literal
}

func (p *Parser) parseFloat() ast.Expression {
	float := &ast.Float{Token: p.currToken}

	value, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as a float", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	float.Value = value

	return float
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}
