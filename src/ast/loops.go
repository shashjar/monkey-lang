package ast

import (
	"bytes"
	"monkey/token"
)

// Represents a while loop in the Monkey programming language, which executes some body as long as some
// provided condition is truthy.
type WhileLoop struct {
	Token     token.Token // the token.WHILE token
	Condition Expression
	Body      *BlockStatement
}

func (wl *WhileLoop) expressionNode() {}

func (wl *WhileLoop) TokenLiteral() string {
	return wl.Token.Literal
}

func (wl *WhileLoop) String() string {
	var out bytes.Buffer

	out.WriteString("while (")
	out.WriteString(wl.Condition.String())
	out.WriteString(") { ")
	out.WriteString(wl.Body.String())
	out.WriteString(" } ")

	return out.String()
}

// Represents a for loop in the Monkey programming language, consisting of an initialization expression, a condition,
// and an expression which is evaluated at the end of each loop iteration.
type ForLoop struct {
	Token        token.Token // the token.FOR token
	Init         Statement
	Condition    Expression
	Afterthought Statement
	Body         *BlockStatement
}

func (fl *ForLoop) expressionNode() {}

func (fl *ForLoop) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *ForLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for (")
	out.WriteString(fl.Init.String())
	out.WriteString(" ")
	out.WriteString(fl.Condition.String())
	out.WriteString("; ")
	out.WriteString(fl.Afterthought.String())
	out.WriteString(") { ")
	out.WriteString(fl.Body.String())
	out.WriteString(" } ")

	return out.String()
}
