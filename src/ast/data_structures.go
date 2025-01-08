package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// Represents an array in the form [element1, ...].
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Represents a hashmap in the form {<expression>: <expression>, ...}.
type HashMapLiteral struct {
	Token   token.Token // the '{' token
	KVPairs map[Expression]Expression
}

func (hml *HashMapLiteral) expressionNode() {}

func (hml *HashMapLiteral) TokenLiteral() string {
	return hml.Token.Literal
}

func (hml *HashMapLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range hml.KVPairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Represents an index expression in the form "<expression>[<expression>]".
type IndexExpression struct {
	Token token.Token // the '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}
