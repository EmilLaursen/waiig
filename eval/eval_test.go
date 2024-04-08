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
	env := object.NewEnv()
	return Eval(program, env)
}

func testIntegerObj(t *testing.T, want int64, got object.Object, msgAndArgs ...any) {
	t.Helper()
	o := testutils.IsType[*object.Integer](t, got, msgAndArgs...)
	require.Equal(t, want, o.Value, msgAndArgs...)
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
		if (10 > 1) {
		  if (10 > 1) {
		    return 10;
		  }

		  return 1;
		}
		`, 10,
		},
		// 		{
		// 			`
		// let f = fn(x) {
		//   return x;
		//   x + 10;
		// };
		// f(10);`,
		// 			10,
		// 		},
		// 		{
		// 			`
		// let f = fn(x) {
		//    let result = x + 10;
		//    return result;
		//    return 10;
		// };
		// f(10);`,
		// 			20,
		// 		},
	}

	for i, tt := range tests {
		got := testEval(tt.input)
		var msg string
		if got != nil {
			msg = fmt.Sprintf("case %d, input=%s, got=%s", i, tt.input, got.Inspect())
		}
		switch w := tt.want.(type) {
		case int:
			testIntegerObj(t, int64(w), got, msg)
		case nil:
			require.Equal(t, NULL, got, msg)
		default:
			t.Errorf("unexpected want: %+v, type=%T", tt.want, tt.want)
		}
	}
}

func TestEvalErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		// {
		// 	"true + false + true + false;",
		// 	"unknown operator: BOOLEAN + BOOLEAN",
		// },
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for i, tt := range tests {
		got := testEval(tt.input)
		var gotmsg string
		if got != nil {
			gotmsg = got.Inspect()
		}
		msg := fmt.Sprintf("case %d, input=%s, got=%s", i, tt.input, gotmsg)
		erro := testutils.IsType[*object.Error](t, got, msg)
		require.Equal(t, tt.want, erro.Msg, msg)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObj(t, tt.expected, testEval(tt.input))
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) {x + 2; };`

	wantBody := "(x + 2)"

	o := testEval(input)
	fn := testutils.IsType[*object.Function](t, o)
	require.Len(t, fn.Params, 1)
	require.Equal(t, fn.Params[0].String(), "x")
	require.Equal(t, wantBody, fn.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObj(t, tt.expected, testEval(tt.input))
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
   fn(y) { x + y;}
};

let addTwo = newAdder(2);
addTwo(2);
`
	testIntegerObj(t, 4, testEval(input))
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	want := "Hello World!"

	r := testEval(input)
	str := testutils.IsType[*object.String](t, r)
	require.Equal(t, want, str.Value)
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	want := "Hello World!"

	r := testEval(input)
	str := testutils.IsType[*object.String](t, r)
	require.Equal(t, want, str.Value)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		// {`len([1, 2, 3])`, 3},
		// {`len([])`, 0},
		// {`puts("hello", "world!")`, nil},
		// {`first([1, 2, 3])`, 1},
		// {`first([])`, nil},
		// {`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		// {`last([1, 2, 3])`, 3},
		// {`last([])`, nil},
		// {`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		// {`rest([1, 2, 3])`, []int{2, 3}},
		// {`rest([])`, nil},
		// {`push([], 1)`, []int{1}},
		// {`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
	}

	for i, tt := range tests {
		got := testEval(tt.input)
		msg := fmt.Sprintf("case %d, input=%s, got=%s", i, tt.input, got.Inspect())
		switch w := tt.want.(type) {
		case int:
			testIntegerObj(t, int64(w), got, msg)
		case nil:
			require.Equal(t, NULL, got, msg)
		case string:
			errObj := testutils.IsType[*object.Error](t, got, msg)
			require.Equal(t, w, errObj.Msg, msg)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1,2*2,3+3]"
	res := testutils.IsType[*object.Array](t, testEval(input))
	require.Equal(t, 3, len(res.Elems))
	testIntegerObj(t, 1, res.Elems[0])
	testIntegerObj(t, 4, res.Elems[1])
	testIntegerObj(t, 6, res.Elems[2])
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"[1, 2, 3][4]",
			2,
		},
		{
			"[1, 2, 3][-5]",
			2,
		},
		{
			"[][0]",
			nil,
		},
		{
			"[][-1]",
			nil,
		},
		{
			"[][123]",
			nil,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			1,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"[1, 2, 3][:-1]",
			[]int64{1, 2},
		},
		{
			"[1, 2, 3][1:2]",
			[]int64{2},
		},
		{
			"[1, 2, 3][:]",
			[]int64{1, 2, 3},
		},
		{
			"[1, 2, 3][1:]",
			[]int64{2, 3},
		},
		{
			`[1, 2, 3,4,5,6,7,8,9,10][-6:-1]`,
			[]int64{5, 6, 7, 8, 9},
		},
		{
			`[1, 2, 3][-1:]`,
			[]int64{3},
		},
		{
			`[1, 2, 3][0:]`,
			[]int64{1, 2, 3},
		},
		{
			`[1, 2, 3][:0]`,
			[]int64{},
		},
		{
			`[1, 2, 3][0:0]`,
			[]int64{},
		},
		{
			`[1, 2, 3][-100:-100]`,
			[]int64{},
		},
		{
			`[1, 2, 3][88:88]`,
			[]int64{},
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[-1:]; i[0]",
			3,
		},
	}

	for i, tt := range tests {
		got := testEval(tt.input)
		msg := fmt.Sprintf("case: %d input=%s want=%+v", i, tt.input, tt.want)
		switch w := tt.want.(type) {
		case int:
			testIntegerObj(t, int64(w), got, msg)
		case []int64:
			arr := testutils.IsType[*object.Array](t, got)
			gots := []int64{}
			for idx, o := range arr.Elems {
				ii := testutils.IsType[*object.Integer](t, o, "case %d: arr.Elems[%d] is not integer", i, idx)
				gots = append(gots, ii.Value)
			}
			require.Equal(t, w, gots, msg)
		default:
			require.Equal(t, NULL, got, msg)
		}
	}
}
