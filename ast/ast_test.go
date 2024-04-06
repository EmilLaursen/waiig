package ast

import (
	"testing"

	"github.com/EmilLaursen/wiig/token"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{token.LET, "let"},
				Name: &Identifier{
					Token: token.Token{token.IDENT, "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{token.IDENT, "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	require.Equal(t, "let myVar = anotherVar;", program.String())
}
