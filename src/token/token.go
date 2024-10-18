package token

// Represents a type of token in the Monkey programming language.
type TokenType string

// Represents a token in the Monkey programming language.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers & Literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Parentheses
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	tokType, ok := keywords[ident]
	if ok {
		return tokType
	} else {
		return IDENT
	}
}
