package eval

import (
	"fmt"
	"testing"

	"github.com/EmilLaursen/wiig/object"
	"github.com/EmilLaursen/wiig/parser"
	"github.com/EmilLaursen/wiig/testutils"
	"github.com/stretchr/testify/require"
)

func testEval(input string) object.Object {
	p := parser.FromInput(input)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObj(t *testing.T, want int64, got object.Object) {
	t.Helper()
	o := testutils.IsType[*object.Integer](t, got)
	require.Equal(t, want, o.Value)
}

func testBooleanObj(t *testing.T, want bool, got object.Object) {
	t.Helper()
	o := testutils.IsType[*object.Boolean](t, got)
	require.Equal(t, want, o.Value)
}

func TestEvalIntegerExp(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		got := testEval(tt.input)
		testIntegerObj(t, tt.want, got)
	}
}

func TestEvalBooleanExp(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		got := testEval(tt.input)
		testBooleanObj(t, tt.want, got)
	}
}

func TestIfElseExp(t *testing.T) {
	// TODO: add me
	tests := []struct {
		input string
		want  any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for i, tt := range tests {
		got := testEval(tt.input)
		msg := fmt.Sprintf("case %d, input=%s, got=%s", i, tt.input, got.Inspect())
		switch w := tt.want.(type) {
		case int:
			testIntegerObj(t, int64(w), got)
		case nil:
			require.Equal(t, NULL, got, msg)
		default:
			t.Errorf("unexpected want: %+v, type=%T", tt.want, tt.want)
		}
	}
}
