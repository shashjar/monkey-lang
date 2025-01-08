package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextTokenBasic(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - token type is wrong. expected=%q, got=%q", i, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - token literal is wrong. expected=%q, got=%q", i, test.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 > 10) {
		return 1;
	} else if (5 == 10) {
	 	return 0;
	} else {
		return -1;
	}

	10 == 10;
	10 != 9;

	""
	"foobar"
	"foo bar"
	let s = "Hello world!";

	[1, 2, "3"];

	{"hello": "world", 1: "me", true: 56, "false": false};

	macro(x, y) { x + y; };

	true && false;
	false || true;

	5 % 2;

	2 <= 6;
	9 >= 8;

	4.57 + 9;

	22 // 7;

	let a = 5;
	a = 6;

	const FIFTY = 50;

	let i = 0;
	while true {
		i = i + 1;
	};

	let arr = [1, 2, 3];
	for (let i = 0; i < len(arr); i = i + 1) {
		puts(arr[i]);
	}

	let i = 0;
	i++;
	i--;

	switch x {
	case "hello":
		puts("hello");
	case "world":
		puts("world");
	default:
		puts("hello world");
	}

	i += 1;
	i -= 3;
	i *= 5;
	i /= 7;
	i //= 9;

	2**3;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.DIV, "/"},
		{token.MUL, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.GT, ">"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.MINUS, "-"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, ""},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LET, "let"},
		{token.IDENT, "s"},
		{token.ASSIGN, "="},
		{token.STRING, "Hello world!"},
		{token.SEMICOLON, ";"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.STRING, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "hello"},
		{token.COLON, ":"},
		{token.STRING, "world"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.COLON, ":"},
		{token.STRING, "me"},
		{token.COMMA, ","},
		{token.TRUE, "true"},
		{token.COLON, ":"},
		{token.INT, "56"},
		{token.COMMA, ","},
		{token.STRING, "false"},
		{token.COLON, ":"},
		{token.FALSE, "false"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.MACRO, "macro"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.FALSE, "false"},
		{token.OR, "||"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.MODULO, "%"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.INT, "2"},
		{token.LTE, "<="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.INT, "9"},
		{token.GTE, ">="},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},

		{token.FLOAT, "4.57"},
		{token.PLUS, "+"},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.INT, "22"},
		{token.INTEGER_DIV, "//"},
		{token.INT, "7"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},

		{token.CONST, "const"},
		{token.IDENT, "FIFTY"},
		{token.ASSIGN, "="},
		{token.INT, "50"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "arr"},
		{token.ASSIGN, "="},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.LET, "let"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.IDENT, "len"},
		{token.LPAREN, "("},
		{token.IDENT, "arr"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.IDENT, "arr"},
		{token.LBRACKET, "["},
		{token.IDENT, "i"},
		{token.RBRACKET, "]"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.INCREMENT, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.DECREMENT, "--"},
		{token.SEMICOLON, ";"},

		{token.SWITCH, "switch"},
		{token.IDENT, "x"},
		{token.LBRACE, "{"},
		{token.CASE, "case"},
		{token.STRING, "hello"},
		{token.COLON, ":"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.STRING, "hello"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.CASE, "case"},
		{token.STRING, "world"},
		{token.COLON, ":"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.STRING, "world"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.DEFAULT, "default"},
		{token.COLON, ":"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.STRING, "hello world"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IDENT, "i"},
		{token.PLUS_ASSIGN, "+="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.MINUS_ASSIGN, "-="},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.MUL_ASSIGN, "*="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.DIV_ASSIGN, "/="},
		{token.INT, "7"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.INTEGER_DIV_ASSIGN, "//="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.INT, "2"},
		{token.EXP, "**"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - token type is wrong. expected=%q, got=%q", i, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - token literal is wrong. expected=%q, got=%q", i, test.expectedLiteral, tok.Literal)
		}
	}
}
