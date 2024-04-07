package lexer

import (
	"fmt"
	"testing"

	"github.com/EmilLaursen/wiig/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextSimple(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		wantType token.TokenType
		wantLit  string
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

	l := New(input)

	for i, tc := range tests {
		var c byte = 0
		if i < len(input) {
			c = input[i]
		}
		t.Run(fmt.Sprintf("input[i]=%c wantType=%s wantLit=%s\n", c, tc.wantType, tc.wantLit), func(t *testing.T) {
			tok := l.NextToken()

			assert.Equal(t, tc.wantType, tok.Type)
			assert.Equal(t, tc.wantLit, tok.Literal)
		})
	}
}

func TestNextToken(t *testing.T) {
	tests := []token.Token{
		token.Ident("let"),
		token.Ident("five"),
		token.Ch("="),
		token.Num("5"),
		token.Ch(";"),
		token.Ident("let"),
		token.Ident("ten"),
		// {token.LET, "let"},
		// {token.IDENT, "five"},
		// {token.ASSIGN, "="},
		// {token.INT, "5"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "ten"},
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
		token.Ch("!"),
		token.Ch("-"),
		token.Ch("/"),
		token.Ch("*"),
		token.Num("5"),
		token.Ch(";"),
		token.Num("5"),
		token.Ch("<"),
		token.Num("10"),
		token.Ch(">"),
		token.Num("5"),
		token.Ch(";"),
		token.Ident("if"),
		token.Ch("("),
		token.Num("0005"),
		token.Ch("<"),
		token.Num("0010"),
		token.Ch(")"),
		token.Ch("{"),
		token.Ident("return"),
		token.Ident("true"),
		token.Ch(";"),
		token.Ch("}"),
		token.Ident("else"),
		token.Ch("{"),
		token.Ident("return"),
		{token.FALSE, "false"},
		token.Ch(";"),
		token.Ch("}"),
		token.Num("10"),
		{token.EQ, "=="},
		token.Num("10"),
		token.Ch(";"),
		token.Num("10"),
		{token.NOT_EQ, "!="},
		token.Num("9"),
		token.Ch(";"),
		{token.STRING, "foobar"},
		token.Ch(";"),
		{token.STRING, "trololo"},
		token.Ch(";"),
		{token.STRING, "foo bar"},
		token.Ch(";"),
		{token.STRING, ""},
		token.Ch(";"),
		{token.EOF, ""},
	}

	input := `let five = 5;
let ten = 10;

let add = fn(x,y) {
   x+y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (0005 < 0010) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar";
"trololo";
"foo bar";
"";
`

	l := New(input)

	n := len(l.input) - 1
	window := 10

	for i, tc := range tests {

		tok := l.NextToken()
		start := min(n, l.position)
		end := min(n, l.position+window)
		require.Equal(t, tc, tok,
			fmt.Sprintf("pos=%d want=%+v window=%+v\n", i, tc, l.input[start:end]))
	}
}
