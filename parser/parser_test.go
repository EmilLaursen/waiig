package parser

import (
	"fmt"
	"testing"

	"github.com/EmilLaursen/wiig/ast"
	"github.com/EmilLaursen/wiig/lexer"
	"github.com/EmilLaursen/wiig/testutils"
	"github.com/EmilLaursen/wiig/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
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
	t.Helper()
	checkParserErrors(t, p)
	require.NotNil(t, program)
	require.Equal(t, numStatements, len(program.Statements), "statements: %+v", program.Statements)
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) {
	t.Helper()
	intLit := testutils.IsType[*ast.IntegerLiteral](t, exp)
	require.Equal(t, val, intLit.Value)
	require.Equal(t, fmt.Sprintf("%d", val), intLit.TokenLiteral())
}

func testIdentifier(t *testing.T, want ast.Expression, got ast.Expression) {
	t.Helper()
	wident := testutils.IsType[*ast.Identifier](t, want)
	ident := testutils.IsType[*ast.Identifier](t, got)
	require.Equal(t, wident.Value, ident.Value)
	require.Equal(t, wident.TokenLiteral(), ident.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, want bool, got ast.Expression) {
	t.Helper()
	bexp := testutils.IsType[*ast.Boolean](t, got)
	require.Equal(t, want, bexp.Value)
	require.Equal(t, fmt.Sprintf("%t", want), bexp.TokenLiteral())
}

func testLetStatement(t *testing.T, wantID, wantVal, got any) {
	t.Helper()
	lexp := testutils.IsType[*ast.LetStatement](t, got)
	testLiteralExpression(t, wantID, lexp.Name)
	testLiteralExpression(t, wantVal, lexp.Value)
}

func testReturnStatement(t *testing.T, wantVal, got any) {
	t.Helper()
	exp := testutils.IsType[*ast.ReturnStatement](t, got)
	testLiteralExpression(t, wantVal, exp)
}

func testLiteralExpression(
	t *testing.T,
	want any,
	got any,
) {
	t.Helper()
	var exp ast.Expression
	switch gv := got.(type) {
	case *ast.ExpressionStatement:
		exp = gv.Expression
	case ast.Expression:
		exp = gv
	case *ast.ReturnStatement:
		exp = gv.ReturnValue
	case nil:
		exp = nil
	default:
		t.Errorf("type of exp not handled. got=%+v type=%T", got, got)
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
	case nil:
		require.Equal(t, v, exp)
	default:
		t.Errorf("type of want exp not handled. want=%T", want)
	}
}

func testInfixExpD(
	t *testing.T,
	want *ast.InfixExpression,
	got any,
) {
	t.Helper()
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
	t.Helper()
	var exp ast.Expression
	switch v := got.(type) {
	case ast.Expression:
		exp = v
	case *ast.ExpressionStatement:
		exp = v.Expression
	default:
		t.Errorf("type of exp not handled. got=%T", got)
	}
	opExp := testutils.IsType[*ast.InfixExpression](t, exp)
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
	pexp := testutils.IsType[*ast.PrefixExpression](t, exp)
	require.Equal(t, wop, pexp.Operator)
	testLiteralExpression(t, wright, pexp.Right)
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input   string
		wantID  string
		wantVal any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		testLetStatement(t, tt.wantID, tt.wantVal, program.Statements[0])
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"return 5;", 5},
		{"return 10", 10},
		{"return true;", true},
		{"return false;", false},
	}

	for _, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		testReturnStatement(t, tt.want, program.Statements[0])
	}
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
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
			Value: &ast.IntegerLiteral{
				Token: token.Token{token.INT, "5"},
				Value: 5,
			},
		},

		&ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "y"}, Value: "y"},
			Value: &ast.IntegerLiteral{
				Token: token.Token{token.INT, "10"},
				Value: 10,
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
			Value: &ast.IntegerLiteral{
				Token: token.Token{token.INT, "838383"},
				Value: 838383,
			},
		},
	}

	for i := range want {
		w := want[i].(*ast.LetStatement)
		actual := program.Statements[i]
		require.Equal(t, w.TokenLiteral(), actual.TokenLiteral())
		letStmt := testutils.IsType[*ast.LetStatement](t, actual)
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
		testutils.IsType[*ast.ReturnStatement](t, actual)
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
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
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

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := testutils.IsType[*ast.IfExpression](t, stmt.Expression)
	testInfixExpression(t, "x", "<", "y", exp.Condition)
	require.Len(t, exp.Consequence.Statements, 1)
	cons := testutils.IsType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	testLiteralExpression(t, "x", cons.Expression)
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := testutils.IsType[*ast.IfExpression](t, stmt.Expression)
	testInfixExpression(t, "x", "<", "y", exp.Condition)
	require.Len(t, exp.Consequence.Statements, 1)
	cons := testutils.IsType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	testLiteralExpression(t, "x", cons.Expression)
	require.Len(t, exp.Alternative.Statements, 1)
	alt := testutils.IsType[*ast.ExpressionStatement](t, exp.Alternative.Statements[0])
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
		stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
		fn := testutils.IsType[*ast.FunctionLiteral](t, stmt.Expression)
		gotTokens := []token.Token{}
		for _, p := range fn.Params {
			gotTokens = append(gotTokens, p.Token)
		}
		require.Equal(t, tt.want, gotTokens)
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2*3,4+5);"

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	cexp := testutils.IsType[*ast.CallExpression](t, stmt.Expression)

	testLiteralExpression(t, "add", cexp.Function)

	require.Len(t, cexp.Arguments, 3)

	testLiteralExpression(t, 1, cexp.Arguments[0])
	testInfixExpression(t, 2, "*", 3, cexp.Arguments[1])
	testInfixExpression(t, 4, "+", 5, cexp.Arguments[2])
}

func TestStringLiteralExp(t *testing.T) {
	input := `"hello world";`
	want := "hello world"

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	str := testutils.IsType[*ast.StringLiteral](t, stmt.Expression)
	require.Equal(t, want, str.Value)
}

func TestParsingArrayLiteral(t *testing.T) {
	input := `[1, 2*2,3+3];`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	arr := testutils.IsType[*ast.ArrayLiteral](t, stmt.Expression)
	assert.Equal(t, 3, len(arr.Elems))

	testIntegerLiteral(t, arr.Elems[0], 1)
	testInfixExpression(t, 2, "*", 2, arr.Elems[1])
	testInfixExpression(t, 3, "+", 3, arr.Elems[2])
}

func TestParsingIndexExpressions(t *testing.T) {
	input := `myArray[1 + 1];`

	p := FromInput(input)
	program := p.ParseProgram()
	baseParseCheck(t, p, program, 1)

	stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0])
	iexp := testutils.IsType[*ast.IndexExpression](t, stmt.Expression)

	testLiteralExpression(t, "myArray", iexp.Left)
	testInfixExpression(t, 1, "+", 1, iexp.Index)
}

func TestParsingSliceExpressions(t *testing.T) {
	tests := []struct {
		input string
		ident string
		left  any
		right any
	}{
		{"arr[1:3]", "arr", 1, 3},
		{"arr[:3]", "arr", nil, 3},
		{"arr[1:]", "arr", 1, nil},
		{"arr[:]", "arr", nil, nil},
	}

	for i, tt := range tests {
		p := FromInput(tt.input)
		program := p.ParseProgram()
		baseParseCheck(t, p, program, 1)
		msg := fmt.Sprintf("case=%d input=%s ident=%s got=%+v", i, tt.input, tt.ident, program.Statements[0])

		stmt := testutils.IsType[*ast.ExpressionStatement](t, program.Statements[0], msg)
		iexp := testutils.IsType[*ast.SliceExpression](t, stmt.Expression, msg)

		fmt.Println(msg)
		testLiteralExpression(t, tt.ident, iexp.Left)
		testLiteralExpression(t, tt.left, iexp.IndexLeft)
		testLiteralExpression(t, tt.right, iexp.IndexRight)
	}
}
