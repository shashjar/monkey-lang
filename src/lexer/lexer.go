package lexer

import (
	"monkey/token"
)

// Represents a lexer that processes (tokenizes) source code.
type Lexer struct {
	input        string
	position     int  // Current index position in input string (points to current character)
	readPosition int  // Current reading position in input string (after current character)
	char         byte // Current character under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(char1) + string(l.char)}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.MUL, l.char)
	case '/':
		tok = newToken(token.DIV, l.char)
	case '!':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(char1) + string(l.char)}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\r' || l.char == '\n' {
		l.readChar()
	}
}

// TODO: currently only supports ASCII. Maybe extend to Unicode later (see p. 19-20).
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) readNumber() string {
	startPosition := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) readString() string {
	startPosition := l.position + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}
