package parser

import (
	"fmt"
	"reflect"
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

func baseParseCheck(
	t *testing.T,
	p *Parser,
	program *ast.Program,
	numStatements int,
) {
	checkParserErrors(t, p)
	require.NotNil(t, program)
	require.Equal(t, numStatements, len(program.Statements), "statements: %+v", program.Statements)
}

func isType[T any](t *testing.T, obj any) T {
	var x T
	r, ok := obj.(T)
	require.True(t, ok, "type of obj=%+v is not type=%T but %s", obj, x, reflect.TypeOf(obj))
	return r
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) {
	intLit := isType[*ast.IntegerLiteral](t, exp)
	require.Equal(t, val, intLit.Value)
	require.Equal(t, intLit.TokenLiteral(), fmt.Sprintf("%d", val))
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
	baseParseCheck(t, p, program, 3)

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
		letStmt := isType[*ast.LetStatement](t, actual)
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
	baseParseCheck(t, p, program, 3)

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
		isType[*ast.ReturnStatement](t, actual)
	}
}

func TestIdentifiers(t *testing.T) {
	input := "foobar;"

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	want := []ast.Statement{
		&ast.ExpressionStatement{
			Token: token.Token{token.IDENT, "foobar"},
			Expression: &ast.Identifier{
				Token: token.Token{token.IDENT, "foobar"},
				Value: "foobar",
			},
		},
	}

	for i := range want {
		w := want[i].(*ast.ExpressionStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		stmt := isType[*ast.ExpressionStatement](t, actual)
		wident := w.Expression.(*ast.Identifier)
		ident := isType[*ast.Identifier](t, stmt.Expression)
		require.Equal(t, wident.Value, ident.Value)
		require.Equal(t, wident.TokenLiteral(), ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	want := []ast.Statement{
		&ast.ExpressionStatement{
			Token: token.Token{token.INT, "5"},
			Expression: &ast.IntegerLiteral{
				Token: token.Token{token.INT, "5"},
				Value: 5,
			},
		},
	}

	for i := range want {
		w := want[i].(*ast.ExpressionStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		stmt := isType[*ast.ExpressionStatement](t, actual)
		wlit := w.Expression.(*ast.IntegerLiteral)
		lit := isType[*ast.IntegerLiteral](t, stmt.Expression)
		require.Equal(t, wlit.Value, lit.Value)
		require.Equal(t, wlit.TokenLiteral(), lit.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input  string
		op     string
		intVal int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		actual := program.Statements[0]
		stmt := isType[*ast.ExpressionStatement](t, actual)
		exp := isType[*ast.PrefixExpression](t, stmt.Expression)
		require.Equal(t, tt.op, exp.Operator)
		testIntegerLiteral(t, exp.Right, tt.intVal)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  int64
		op       string
		rightVal int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		actual := program.Statements[0]
		stmt := isType[*ast.ExpressionStatement](t, actual)
		exp := isType[*ast.InfixExpression](t, stmt.Expression)
		require.Equal(t, tt.op, exp.Operator)
		testIntegerLiteral(t, exp.Left, tt.leftVal)
		testIntegerLiteral(t, exp.Right, tt.rightVal)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		// TODO: add remaining test cases p 87
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		require.Equal(t, tt.want, program.String())
	}
}
