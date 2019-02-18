package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENTIFIER" // cos, lod, x, y, ...

	NUM   = "NUMBER"    // 1343456.78

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	POWER    = "^"

	LPAREN = "("
	RPAREN = ")"

	// Keywords
	SIN  = "sin"
	COS  = "cos"
	TAN  = "tan"
	SEC  = "sec"
	LOG  = "log"
	LN   = "ln"
	ROOT = "root"
	E    = "e"
)

var keywords = map[string]TokenType{
	"sin":  SIN,
	"cos":  COS,
	"tan":  TAN,
	"sec":  SEC,
	"log":  LOG,
	"ln":   LN,
	"root": ROOT,
	"e":    E,
}

func Lookup(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT // this is probably a variable
}
