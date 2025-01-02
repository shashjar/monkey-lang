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
	FLOAT  = "FLOAT"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRING = "STRING"

	// Operators
	ASSIGN = "="

	PLUS        = "+"
	MINUS       = "-"
	MUL         = "*"
	DIV         = "/"
	INTEGER_DIV = "//"

	MODULO = "%"

	BANG = "!"

	AND = "&&"
	OR  = "||"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LTE    = "<="
	GTE    = ">="

	INCREMENT = "++"
	DECREMENT = "--"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

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
	CONST    = "CONST"
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
	FOR      = "FOR"
	RETURN   = "RETURN"
	MACRO    = "MACRO"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"const":  CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,
	"macro":  MACRO,
}

func LookupIdent(ident string) TokenType {
	tokType, ok := keywords[ident]
	if ok {
		return tokType
	} else {
		return IDENT
	}
}
