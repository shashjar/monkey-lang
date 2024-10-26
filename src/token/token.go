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
	IDENT  = "IDENT"
	INT    = "INT"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRING = "STRING"

	// Operators
	ASSIGN = "="

	PLUS  = "+"
	MINUS = "-"
	MUL   = "*"
	DIV   = "/"

	BANG = "!"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Parentheses
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	tokType, ok := keywords[ident]
	if ok {
		return tokType
	} else {
		return IDENT
	}
}
