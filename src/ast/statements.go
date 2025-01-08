package ast

import (
	"bytes"
	"monkey/token"
)

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

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// Represents a const declaration statement in the Monkey programming language, consisting of (1) the CONST token,
// (2) the name of the identifier in the binding, and (3) the expression that produces the value for the binding.
// Variables declared as consts cannot be reassigned other values later.
type ConstStatement struct {
	Token token.Token // the token.CONST token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode() {}

func (cs *ConstStatement) TokenLiteral() string {
	return cs.Token.Literal
}

func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")
	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// Represents an assign statement in the Monkey programming language, consisting of (1) the name of
// the identifier in the binding and (2) the expression that produces the value for the binding. Note
// that a variable binding must first be declared with 'let' before it can be assigned a new value
// with an assign statement.
type AssignStatement struct {
	Token token.Token // the token.IDENT token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {}

func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")

	return out.String()
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

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// Represents a series of statements, e.g. those contained within if-else expressions.
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
