package ast

import (
	"bytes"
)

// Represents a Node in the AST.
type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
