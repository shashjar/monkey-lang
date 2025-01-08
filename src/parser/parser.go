package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// Represents a parser for the Monkey programming language, consisting of (1) a lexer for tokenization,
// (2) any errors encountered during parsing, (3) the current token under examination, (4) the next token,
// and (5) the parsing functions for operators in different (prefix, infix) positions.
type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Defines the operator precedences of the Monkey programming language.
const (
	_ int = iota
	LOWEST
	EQUALS       // == or !=
	AND_OR       // && or ||
	LESS_GREATER // > or < or <= or >=
	SUM          // + or -
	PRODUCT      // * or / or // or %
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.AND:         AND_OR,
	token.OR:          AND_OR,
	token.LT:          LESS_GREATER,
	token.GT:          LESS_GREATER,
	token.LTE:         LESS_GREATER,
	token.GTE:         LESS_GREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.MUL:         PRODUCT,
	token.DIV:         PRODUCT,
	token.INTEGER_DIV: PRODUCT,
	token.MODULO:      PRODUCT,
	token.LPAREN:      CALL,
	token.LBRACKET:    INDEX,
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloat)
	p.registerPrefix(token.STRING, p.parseString)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.SWITCH, p.parseSwitchStatement)
	p.registerPrefix(token.WHILE, p.parseWhileLoop)
	p.registerPrefix(token.FOR, p.parseForLoop)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashMapLiteral)
	p.registerPrefix(token.MACRO, p.parseMacroLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.INTEGER_DIV, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Read two tokens so that both currToken & peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.createNoPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

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

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	function := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	function.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		params = append(params, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: p.currToken, Function: function}
	callExpression.Arguments = p.parseExpressionList(token.RPAREN)
	return callExpression
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	expList := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return expList
	}

	p.nextToken()
	expList = append(expList, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expList = append(expList, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return expList
}

func (p *Parser) parseHashMapLiteral() ast.Expression {
	hashmap := &ast.HashMapLiteral{Token: p.currToken}
	hashmap.KVPairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}

		hashmap.KVPairs[key] = value
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hashmap
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.createPeekError(t)
		return false
	}
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) currTokenIn(ts []token.TokenType) bool {
	for _, t := range ts {
		if p.currTokenIs(t) {
			return true
		}
	}
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenIn(ts []token.TokenType) bool {
	for _, t := range ts {
		if p.peekTokenIs(t) {
			return true
		}
	}
	return false
}

func (p *Parser) createPeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) createNoPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
