package ast

import "monkey/token"

// Represents a Node in the AST.
type Node interface {
	TokenLiteral() string
}

// Represents a single statement in the series included in the Monkey program.
type Statement interface {
	Node
	statementNode()
}

// Represents an expression (produces a value) included in the Monkey program.
type Expression interface {
	Node
	expressionNode()
}

// Represents the root node of every Monkey AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Represents a let statement in the Monkey programming language, consisting of (1) the LET token, (2) the
// name of the identifier in the binding, and (3) the expression that produces the value for the binding.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Represents a return statement in the Monkey programming language, consisting of the RETURN token
// and the expression providing the value that is being returned.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// Represents an identifier, consisting of the IDENT token and the value (name of the identifier).
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
