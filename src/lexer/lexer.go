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
			tok = l.readTwoCharacterToken(token.EQ)
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.readTwoCharacterToken(token.AND)
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	case '|':
		if l.peekChar() == '|' {
			tok = l.readTwoCharacterToken(token.OR)
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.MUL, l.char)
	case '/':
		if l.peekChar() == '/' {
			tok = l.readTwoCharacterToken(token.INTEGER_DIV)
		} else {
			tok = newToken(token.DIV, l.char)
		}
	case '%':
		tok = newToken(token.MODULO, l.char)
	case '!':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(char1) + string(l.char)}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '<':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(char1) + string(l.char)}
		} else {
			tok = newToken(token.LT, l.char)
		}
	case '>':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(char1) + string(l.char)}
		} else {
			tok = newToken(token.GT, l.char)
		}
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case ':':
		tok = newToken(token.COLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '[':
		tok = newToken(token.LBRACKET, l.char)
	case ']':
		tok = newToken(token.RBRACKET, l.char)
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
			tok = l.readNumber()
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

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readTwoCharacterToken(tokenType token.TokenType) token.Token {
	char1 := l.char
	l.readChar()
	tok := token.Token{Type: tokenType, Literal: string(char1) + string(l.char)}
	return tok
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) readNumber() token.Token {
	startPosition := l.position
	var tokenType token.TokenType
	tokenType = token.INT
	for isDigit(l.char) || l.char == '.' {
		if l.char == '.' {
			tokenType = token.FLOAT
		}
		l.readChar()
	}

	return token.Token{Type: tokenType, Literal: l.input[startPosition:l.position]}
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
