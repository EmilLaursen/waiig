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
	require.Equal(t, fmt.Sprintf("%d", val), intLit.TokenLiteral())
}

func testIdentifier(t *testing.T, want ast.Expression, got ast.Expression) {
	wident := isType[*ast.Identifier](t, want)
	ident := isType[*ast.Identifier](t, got)
	require.Equal(t, wident.Value, ident.Value)
	require.Equal(t, wident.TokenLiteral(), ident.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, want bool, got ast.Expression) {
	bexp := isType[*ast.Boolean](t, got)
	require.Equal(t, want, bexp.Value)
	require.Equal(t, fmt.Sprintf("%t", want), bexp.TokenLiteral())
}

func testLiteralExpression(
	t *testing.T,
	want any,
	got any,
) {
	var exp ast.Expression
	switch gv := got.(type) {
	case *ast.ExpressionStatement:
		exp = gv.Expression
	case ast.Expression:
		exp = gv
	default:
		t.Errorf("type of exp not handled. got=%T", got)
	}

	switch v := want.(type) {
	case *ast.ExpressionStatement:
		testLiteralExpression(t, v.Expression, exp)
	case *ast.IntegerLiteral:
		testIntegerLiteral(t, exp, v.Value)
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, &ast.Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: v,
			},
			Value: v,
		}, exp)
	case *ast.Identifier:
		testIdentifier(t, v, exp)
	case *ast.Boolean:
		testBooleanLiteral(t, v.Value, exp)
	case bool:
		testBooleanLiteral(t, v, exp)
	default:
		t.Errorf("type of want exp not handled. want=%T", want)
	}
}

func testInfixExpD(
	t *testing.T,
	want *ast.InfixExpression,
	got any,
) {
	// want token unused
	testInfixExpression(t, want.Left, want.Operator, want.Right, got)
}

func testInfixExpression(
	t *testing.T,
	wleft any,
	wop string,
	wright any,
	got any,
) {
	var exp ast.Expression
	switch v := got.(type) {
	case ast.Expression:
		exp = v
	case *ast.ExpressionStatement:
		exp = v.Expression
	default:
		t.Errorf("type of exp not handled. got=%T", got)
	}
	opExp := isType[*ast.InfixExpression](t, exp)
	testLiteralExpression(t, wleft, opExp.Left)
	require.Equal(t, wop, opExp.Operator)
	testLiteralExpression(t, wright, opExp.Right)
}

func testPrefixExpression(
	t *testing.T,
	wop string,
	wright any,
	got any,
) {
	var exp ast.Expression
	switch v := got.(type) {
	case ast.Expression:
		exp = v
	case *ast.ExpressionStatement:
		exp = v.Expression
	default:
		t.Errorf("type of exp not handled. got=%T", got)
	}
	pexp := isType[*ast.PrefixExpression](t, exp)
	require.Equal(t, wop, pexp.Operator)
	testLiteralExpression(t, wright, pexp.Right)
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
		// TODO: refactor this, together with return stmt
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
		// TODO: need to test more
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
		testLiteralExpression(t, w, actual)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `
true;
false;
`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 2)

	want := []ast.Statement{
		&ast.ExpressionStatement{
			Token: token.Token{
				Type:    token.TRUE,
				Literal: "true",
			},
			Expression: &ast.Boolean{
				Token: token.Token{
					Type:    token.TRUE,
					Literal: "true",
				},
				Value: true,
			},
		},
		&ast.ExpressionStatement{
			Token: token.Token{
				Type:    token.FALSE,
				Literal: "false",
			},
			Expression: &ast.Boolean{
				Token: token.Token{
					Type:    token.FALSE,
					Literal: "false",
				},
				Value: false,
			},
		},
	}

	for i := range want {
		w := want[i].(*ast.ExpressionStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		testLiteralExpression(t, w, actual)
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
		testLiteralExpression(t, w, actual)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		op    string
		val   any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		actual := program.Statements[0]
		testPrefixExpression(t, tt.op, tt.val, actual)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  any
		op       string
		rightVal any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		actual := program.Statements[0]
		testInfixExpression(t, tt.leftVal, tt.op, tt.rightVal, actual)
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

		{"a*b*c", "((a * b) * c)"},
		{"a*b/c", "((a * b) / c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e -f", "(((a + (b * c)) + (d / e)) - f)"},

		{"3+4;-5*5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},

		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		require.Equal(t, tt.want, program.String())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := isType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := isType[*ast.IfExpression](t, stmt.Expression)
	testInfixExpression(t, "x", "<", "y", exp.Condition)
	require.Len(t, exp.Consequence.Statements, 1)
	cons := isType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	testLiteralExpression(t, "x", cons.Expression)
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := isType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := isType[*ast.IfExpression](t, stmt.Expression)
	testInfixExpression(t, "x", "<", "y", exp.Condition)
	require.Len(t, exp.Consequence.Statements, 1)
	cons := isType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	testLiteralExpression(t, "x", cons.Expression)
	require.Len(t, exp.Alternative.Statements, 1)
	alt := isType[*ast.ExpressionStatement](t, exp.Alternative.Statements[0])
	testLiteralExpression(t, "y", alt.Expression)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x,y) {x+y;}`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	want := []ast.Statement{
		&ast.ExpressionStatement{
			Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
			Expression: &ast.FunctionLiteral{
				Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
				Params: []*ast.Identifier{
					{Token: token.Token{token.IDENT, "x"}, Value: "x"},
					{Token: token.Token{token.IDENT, "y"}, Value: "y"},
				},
				Body: &ast.BlockStatement{
					Token: token.Token{token.LBRACE, "{"},
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Token: token.Token{token.IDENT, "x"},
							Expression: &ast.InfixExpression{
								Token:    token.Token{Type: token.PLUS, Literal: "+"},
								Operator: "+",
								Left:     &ast.Identifier{Token: token.Token{token.IDENT, "x"}, Value: "x"},
								Right:    &ast.Identifier{Token: token.Token{token.IDENT, "y"}, Value: "y"},
							},
						},
					},
				},
			},
		},
	}
	require.Equal(t, want[0], program.Statements[0])
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input string
		want  []token.Token
	}{
		{"fn() {};", []token.Token{}},
		{"fn(x) {};", []token.Token{token.Ident("x")}},
		{"fn(x,y,z) {};", []token.Token{token.Ident("x"), token.Ident("y"), token.Ident("z")}},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		stmt := isType[*ast.ExpressionStatement](t, program.Statements[0])
		fn := isType[*ast.FunctionLiteral](t, stmt.Expression)
		gotTokens := []token.Token{}
		for _, p := range fn.Params {
			gotTokens = append(gotTokens, p.Token)
		}
		require.Equal(t, tt.want, gotTokens)
	}
}
