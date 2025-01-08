package ast

import (
	"bytes"
	"monkey/token"
)

// Represents a clause (`if` or `else if`) in a conditional expression.
type ConditionalClause struct {
	Condition   Expression
	Consequence *BlockStatement
}

// Represents an if expression in the form "if (<condition>) <consequence> else if (<condition>) <consequence> ... else <alternative>"
// `else if` and `else` clauses are optional
type IfExpression struct {
	Token       token.Token         // the token.IF token
	Clauses     []ConditionalClause // will always contain at least one clause (the `if` clause itself), and additional if any `else if`s are used
	Alternative *BlockStatement     // the `else` clause body, if `else` was used
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Clauses[0].Condition.String())
	out.WriteString(") { ")
	out.WriteString(ie.Clauses[0].Consequence.String())
	out.WriteString(" } ")

	for i := 1; i < len(ie.Clauses); i++ {
		out.WriteString("else if (")
		out.WriteString(ie.Clauses[i].Condition.String())
		out.WriteString(") { ")
		out.WriteString(ie.Clauses[i].Consequence.String())
		out.WriteString(" } ")
	}

	if ie.Alternative != nil {
		out.WriteString("else { ")
		out.WriteString(ie.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

// Represents a case in a switch statement.
type SwitchCase struct {
	Expression  Expression
	Consequence *BlockStatement
}

// Represents a switch statement in the form "switch <expression> { case <expression>: <consequence> ... default: <default-consequence> }"
// At least one `case` must be provided, and `default` is optional
type SwitchStatement struct {
	Token            token.Token     // the token.SWITCH token
	SwitchExpression Expression      // the expression we are switching on
	Cases            []SwitchCase    // at least one `case` is required
	Default          *BlockStatement // the `default` case, if one was provided
}

func (ss *SwitchStatement) expressionNode() {}

func (ss *SwitchStatement) TokenLiteral() string {
	return ss.Token.Literal
}

func (ss *SwitchStatement) String() string {
	var out bytes.Buffer

	out.WriteString("switch ")
	out.WriteString(ss.SwitchExpression.String())
	out.WriteString(" { ")
	out.WriteString("case " + ss.Cases[0].Expression.String() + ": ")
	out.WriteString(ss.Cases[0].Consequence.String())

	for i := 1; i < len(ss.Cases); i++ {
		out.WriteString(" case " + ss.Cases[i].Expression.String() + ": ")
		out.WriteString(ss.Cases[i].Consequence.String())
	}

	if ss.Default != nil {
		out.WriteString(" default: " + ss.Default.String())
	}

	out.WriteString(" }")

	return out.String()
}
