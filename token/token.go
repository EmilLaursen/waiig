package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// OPERATORS
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COLON    = ":"

	// KEYWORDS
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
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

func Ident(ident string) Token {
	if kw, ok := keywords[ident]; ok {
		return Token{Type: kw, Literal: ident}
	}
	return Token{Type: IDENT, Literal: ident}
}

func Num(num string) Token {
	return Token{Type: INT, Literal: num}
}

func Str(str string) Token {
	return Token{
		Type:    STRING,
		Literal: str,
	}
}

func Ch(ch string) Token {
	var tp TokenType
	lit := ch

	switch ch {
	case "=":
		tp = ASSIGN
	case ";":
		tp = SEMICOLON
	case "(":
		tp = LPAREN
	case ")":
		tp = RPAREN
	case ",":
		tp = COMMA
	case "+":
		tp = PLUS
	case "-":
		tp = MINUS
	case "!":
		tp = BANG
	case "/":
		tp = SLASH
	case "*":
		tp = ASTERISK
	case "<":
		tp = LT
	case ">":
		tp = GT
	case "{":
		tp = LBRACE
	case "}":
		tp = RBRACE
	case string(byte(0)):
		tp = EOF
		lit = ""
	default:
		tp = ILLEGAL
	}
	return Token{Type: tp, Literal: lit}
}
