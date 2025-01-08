package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseMacroLiteral() ast.Expression {
	macro := &ast.MacroLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	macro.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	macro.Body = p.parseBlockStatement()

	return macro
}
