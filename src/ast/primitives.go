package ast

import "monkey/token"

// Represents an integer, consisting of the INT token and the value of the integer.
type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// Represents a 64-bit float number, consisting of the FLOAT token and the value of the float.
type Float struct {
	Token token.Token // the token.FLOAT token
	Value float64
}

const FLOAT_64_EQUALITY_THRESHOLD = 1e-9

func (f *Float) expressionNode() {}

func (f *Float) TokenLiteral() string {
	return f.Token.Literal
}

func (f *Float) String() string {
	return f.Token.Literal
}

// Represents a boolean, consisting of the TRUE or FALSE token and the value of the boolean.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

// Represents a string, consisting of the token.STRING token and the value of the string.
type StringLiteral struct {
	Token token.Token // the token.STRING token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
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

func (i *Identifier) String() string {
	return i.Value
}
