package token

// Represents a type of token in the Monkey programming language.
type TokenType string

// Represents a token in the Monkey programming language.
type Token struct {
	Type         TokenType
	Literal      string
	LineNumber   int
	ColumnNumber int
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

	PLUS_ASSIGN        = "+="
	MINUS_ASSIGN       = "-="
	MUL_ASSIGN         = "*="
	DIV_ASSIGN         = "/="
	INTEGER_DIV_ASSIGN = "//="

	EXP    = "**"
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
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	WHILE    = "WHILE"
	FOR      = "FOR"
	RETURN   = "RETURN"
	MACRO    = "MACRO"
)

var OPERATOR_ASSIGNMENTS = []TokenType{
	PLUS_ASSIGN,
	MINUS_ASSIGN,
	MUL_ASSIGN,
	DIV_ASSIGN,
	INTEGER_DIV_ASSIGN,
}

var keywords = map[string]TokenType{
	"fn":      FUNCTION,
	"let":     LET,
	"const":   CONST,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"switch":  SWITCH,
	"case":    CASE,
	"default": DEFAULT,
	"while":   WHILE,
	"for":     FOR,
	"return":  RETURN,
	"macro":   MACRO,
}

func LookupIdent(ident string) TokenType {
	tokType, ok := keywords[ident]
	if ok {
		return tokType
	} else {
		return IDENT
	}
}
