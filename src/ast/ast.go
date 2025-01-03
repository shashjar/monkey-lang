package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
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

// Represents a statement that consists solely of one expression in the Monkey program.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Represents a prefix expression that has some operator in the prefix position acting on some expression to the right.
type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// Represents an infix expression that has some operator in the infix position acting on two
// expressions, one to the left and one to the right.
type InfixExpression struct {
	Token    token.Token // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
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

// Represents a function literal in the form "fn <parameters> <block statement>".
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", fl.Name))
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// Represents a call expression (calling a function) in the form "<expression>(<comma-separated expressions>)".
type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // can be Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

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

// Represents a macro literal.
type MacroLiteral struct {
	Token      token.Token // the token.MACRO token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode() {}

func (ml *MacroLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

func (ml *MacroLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(ml.Body.String())

	return out.String()
}
