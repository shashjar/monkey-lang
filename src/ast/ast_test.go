package ast

import (
	"monkey/token"
	"reflect"
	"testing"
)

func TestStringLetStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() is wrong. got=%q", program.String())
	}
}

func TestStringIfExpression(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1} }
	two := func() Expression { return &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2} }
	three := func() Expression { return &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "3"}, Value: 3} }

	ie := &IfExpression{
		Clauses: []ConditionalClause{
			{
				Condition:   &Boolean{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true},
				Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
			},
			{
				Condition:   &Boolean{Token: token.Token{Type: token.FALSE, Literal: "false"}, Value: false},
				Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
			},
		},
		Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: three()}}},
	}
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: ie,
			},
		},
	}

	if program.String() != "if (true) { 1 } else if (false) { 2 } else { 3 }" {
		t.Errorf("program.String() is wrong. got=%q", program.String())
	}
}

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpression{Left: two(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&PrefixExpression{Operator: "!", Right: one()},
			&PrefixExpression{Operator: "!", Right: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&ReturnStatement{ReturnValue: one()},
			&ReturnStatement{ReturnValue: two()},
		},
		{
			&BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}, &ExpressionStatement{Expression: two()}, &ExpressionStatement{Expression: one()}}},
			&BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}, &ExpressionStatement{Expression: two()}, &ExpressionStatement{Expression: two()}}},
		},
		{
			&IfExpression{
				Clauses: []ConditionalClause{
					{
						Condition:   one(),
						Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
					},
					{
						Condition:   one(),
						Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
					},
				},
				Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
			},
			&IfExpression{
				Clauses: []ConditionalClause{
					{
						Condition:   two(),
						Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
					},
					{
						Condition:   two(),
						Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
					},
				},
				Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
			},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&ArrayLiteral{Elements: []Expression{one(), one()}},
			&ArrayLiteral{Elements: []Expression{two(), two()}},
		},
		{
			&IndexExpression{Left: one(), Index: one()},
			&IndexExpression{Left: two(), Index: two()},
		},
	}

	for _, test := range tests {
		modified := Modify(test.input, turnOneIntoTwo)
		equal := reflect.DeepEqual(modified, test.expected)
		if !equal {
			t.Errorf("not equal. expected=%#v, got=%#v", test.expected, modified)
		}
	}

	hashmapLiteral := &HashMapLiteral{
		KVPairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}
	Modify(hashmapLiteral, turnOneIntoTwo)
	for key, val := range hashmapLiteral.KVPairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}
		val, _ := val.(*IntegerLiteral)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
