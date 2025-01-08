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
	lineNumber   int  // Line number (starting from 1) in the input that the lexer is processing
	columnNumber int  // Index of the column (starting from 0) in the current line that the lexer is processing
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, lineNumber: 1, columnNumber: -1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tokLineNumber := l.lineNumber
	tokColumnNumber := l.columnNumber

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			tok = l.readTwoCharacterToken(token.EQ)
		} else {
			tok = l.newToken(token.ASSIGN, l.char)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.readTwoCharacterToken(token.AND)
		} else {
			tok = l.newToken(token.ILLEGAL, l.char)
		}
	case '|':
		if l.peekChar() == '|' {
			tok = l.readTwoCharacterToken(token.OR)
		} else {
			tok = l.newToken(token.ILLEGAL, l.char)
		}
	case '+':
		if l.peekChar() == '+' {
			tok = l.readTwoCharacterToken(token.INCREMENT)
		} else if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.PLUS_ASSIGN, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.PLUS, l.char)
		}
	case '-':
		if l.peekChar() == '-' {
			tok = l.readTwoCharacterToken(token.DECREMENT)
		} else if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.MINUS_ASSIGN, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.MINUS, l.char)
		}
	case '*':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.MUL_ASSIGN, string(char1)+string(l.char))
		} else if l.peekChar() == '*' {
			tok = l.readTwoCharacterToken(token.EXP)
		} else {
			tok = l.newToken(token.MUL, l.char)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			if l.peekChar() == '=' {
				char1 := l.char
				l.readChar()
				tok = l.makeToken(token.INTEGER_DIV_ASSIGN, string(char1)+string(char1)+string(l.char))
			} else {
				tok = l.makeToken(token.INTEGER_DIV, string(l.char)+string(l.char))
			}
		} else if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.DIV_ASSIGN, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.DIV, l.char)
		}
	case '%':
		tok = l.newToken(token.MODULO, l.char)
	case '!':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.NOT_EQ, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.BANG, l.char)
		}
	case '<':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.LTE, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.LT, l.char)
		}
	case '>':
		if l.peekChar() == '=' {
			char1 := l.char
			l.readChar()
			tok = l.makeToken(token.GTE, string(char1)+string(l.char))
		} else {
			tok = l.newToken(token.GT, l.char)
		}
	case ',':
		tok = l.newToken(token.COMMA, l.char)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.char)
	case ':':
		tok = l.newToken(token.COLON, l.char)
	case '(':
		tok = l.newToken(token.LPAREN, l.char)
	case ')':
		tok = l.newToken(token.RPAREN, l.char)
	case '{':
		tok = l.newToken(token.LBRACE, l.char)
	case '}':
		tok = l.newToken(token.RBRACE, l.char)
	case '[':
		tok = l.newToken(token.LBRACKET, l.char)
	case ']':
		tok = l.newToken(token.RBRACKET, l.char)
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
			tok.LineNumber = tokLineNumber
			tok.ColumnNumber = tokColumnNumber
			return tok
		} else if isDigit(l.char) {
			tok = l.readNumber()
			tok.LineNumber = tokLineNumber
			tok.ColumnNumber = tokColumnNumber
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.char)
		}
	}

	tok.LineNumber = tokLineNumber
	tok.ColumnNumber = tokColumnNumber

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

	if l.char == '\n' {
		l.lineNumber += 1
		l.columnNumber = -1
	} else {
		l.columnNumber += 1
	}
}

func (l *Lexer) readTwoCharacterToken(tokenType token.TokenType) token.Token {
	char1 := l.char
	l.readChar()
	tok := l.makeToken(tokenType, string(char1)+string(l.char))
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

	return l.makeToken(tokenType, l.input[startPosition:l.position])
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

func (l *Lexer) newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char), LineNumber: l.lineNumber, ColumnNumber: l.columnNumber}
}

func (l *Lexer) makeToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, LineNumber: l.lineNumber, ColumnNumber: l.columnNumber}
}
