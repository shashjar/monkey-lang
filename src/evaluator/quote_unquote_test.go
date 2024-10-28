package evaluator

import (
	"monkey/object"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote("hello world!")`,
			`hello world!`,
		},
		{
			`quote(true)`,
			`true`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(foobar)`,
			`foobar`,
		},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != test.expected {
			t.Errorf("quoted node string is wrong. got=%q, want=%q", quote.Node.String(), test.expected)
		}
	}
}

func TestQuoteErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote()`,
			`wrong number of arguments provided to 'quote'. expected=1, got=0`,
		},
		{
			`quote(5, 6)`,
			`wrong number of arguments provided to 'quote'. expected=1, got=2`,
		},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		testErrorObject(t, evaluated, test.expected)
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(unquote(4))`,
			`4`,
		},
		{
			`quote(unquote(4 + 4))`,
			`8`,
		},
		{
			`quote(7 + unquote(4 + 4))`,
			`(7 + 8)`,
		},
		{
			`quote(unquote(4 + 4) + 7)`,
			`(8 + 7)`,
		},
		{
			`let foobar = 8;
			quote(foobar)`,
			`foobar`,
		},
		{
			`let foobar = 8; quote(unquote(foobar))`,
			`8`,
		},
		{
			`quote(unquote(true))`,
			`true`,
		},
		{
			`quote(unquote(true == false))`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`(4 + 4)`,
		},
		{
			`let quotedInfixExpression = quote(4 + 4);
			quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != test.expected {
			t.Errorf("quoted node string is wrong. got=%q, want=%q", quote.Node.String(), test.expected)
		}
	}
}
