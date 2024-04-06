package parser

import (
	"testing"

	"github.com/EmilLaursen/wiig/ast"
	"github.com/EmilLaursen/wiig/lexer"
	"github.com/EmilLaursen/wiig/token"
	"github.com/stretchr/testify/require"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program)
	require.Equal(t, 3, len(program.Statements), program.Statements)

	want := []ast.Statement{
		&ast.LetStatement{
			Token: token.Token{
				Type:    token.LET,
				Literal: "let",
			},
			Name: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
		},

		&ast.LetStatement{
			Token: token.Token{
				Type:    token.LET,
				Literal: "let",
			},
			Name: &ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "y",
				}, Value: "y",
			},
		},

		&ast.LetStatement{
			Token: token.Token{
				Type:    token.LET,
				Literal: "let",
			},
			Name: &ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foobar",
				}, Value: "foobar",
			},
		},
	}

	for i := range want {
		w := want[i].(*ast.LetStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		letStmt, ok := actual.(*ast.LetStatement)
		require.True(t, ok)
		require.Equal(t, w.Name.Value, letStmt.Name.Value)
		require.Equal(t, w.Name.TokenLiteral(), letStmt.Name.TokenLiteral())
	}

	require.Equal(t, want, program.Statements)
}

func TestReturnStatemens(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program)
	require.Equal(t, 3, len(program.Statements), program.Statements)

	want := []ast.Statement{
		&ast.ReturnStatement{
			Token: token.Token{
				Type:    token.RETURN,
				Literal: "return",
			},
			ReturnValue: nil,
		},

		&ast.ReturnStatement{
			Token: token.Token{
				Type:    token.RETURN,
				Literal: "return",
			},
			ReturnValue: nil,
		},

		&ast.ReturnStatement{
			Token: token.Token{
				Type:    token.RETURN,
				Literal: "return",
			},
			ReturnValue: nil,
		},
	}

	for i := range want {
		w := want[i].(*ast.ReturnStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		_, ok := actual.(*ast.ReturnStatement)
		require.True(t, ok)
	}
}
