package eval

import (
	"fmt"

	"github.com/EmilLaursen/wiig/ast"
	"github.com/EmilLaursen/wiig/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func isError(o object.Object) bool {
	return o != nil && o.Type() == object.ERROR_OBJ
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExp(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBoolObj(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ReturnStatement:
		v := Eval(node.ReturnValue, env)
		if isError(v) {
			return v
		}
		return &object.ReturnValue{Value: v}

	case *ast.LetStatement:
		v := Eval(node.Value, env)
		if isError(v) {
			return v
		}

		env.Set(node.Name.Value, v)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{
			Params: node.Params,
			Body:   node.Body,
			Env:    env,
		}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.SliceExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		var ileft object.Object = NULL
		var iright object.Object = NULL

		if node.IndexLeft != nil {
			ileft = Eval(node.IndexLeft, env)
			if isError(ileft) {
				return ileft
			}
		}

		if node.IndexRight != nil {
			iright = Eval(node.IndexRight, env)
			if isError(iright) {
				return iright
			}
		}

		return evalSliceExpression(left, ileft, iright)

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyfunction(fn, args)

	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elems, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elems: elems}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExp(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExp(node.Operator, left, right)

	}

	return nil
}

func applyfunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		scope := object.NewScope(fn.Env)
		for i, p := range fn.Params {
			scope.Set(p.Value, args[i])
		}
		ret := Eval(fn.Body, scope)
		return unwrapReturn(ret)

	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newErr("not a function: %s", fn.Type())
	}
}

func unwrapReturn(o object.Object) object.Object {
	if r, ok := o.(*object.ReturnValue); ok {
		return r.Value
	}
	return o
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var r []object.Object
	for _, e := range exps {
		evald := Eval(e, env)
		if isError(evald) {
			return []object.Object{evald}
		}
		r = append(r, evald)
	}
	return r
}

func evalIfExp(n *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(n.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(n.Consequence, env)
	} else if n.Alternative != nil {
		return Eval(n.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(o object.Object) bool {
	switch o {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for kn, vn := range node.Pairs {
		key := Eval(kn, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newErr("key is not hashable: %s", key.Type())
		}

		value := Eval(vn, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return &object.Hash{Pairs: pairs}
}

func evalSliceExpression(left, ileft, rleft object.Object) object.Object {
	ileftGood := object.IsTypeOrNULL(ileft, object.INTEGER_OBJ)
	irightGood := object.IsTypeOrNULL(rleft, object.INTEGER_OBJ)

	switch {
	case left.Type() == object.ARRAY_OBJ && ileftGood && irightGood:
		return evalArraySliceExpression(left, ileft, rleft)
	default:
		return newErr("slice operator not supported: %s[%s:%s]", left.Type(), ileft.Type(), rleft.Type())
	}
}

func evalArraySliceExpression(array, ileft, iright object.Object) object.Object {
	arrobj := array.(*object.Array)
	n := len(arrobj.Elems)

	if n == 0 {
		return &object.Array{Elems: []object.Object{}}
	}

	var leftIdx int64
	var rightIdx int64 = int64(n)

	if ileft.Type() != object.NULL_OBJ {
		leftIdx = getIndex(ileft.(*object.Integer).Value, n)
	}
	if iright.Type() != object.NULL_OBJ {
		rightIdx = getIndex(iright.(*object.Integer).Value, n)
	}

	if leftIdx < 0 || rightIdx < 0 || leftIdx >= rightIdx {
		return &object.Array{Elems: []object.Object{}}
	}

	elems := make([]object.Object, rightIdx-leftIdx)
	// copying pointers to integers, but these should never change?
	// TODO: test this
	copy(elems, arrobj.Elems[leftIdx:rightIdx])
	return &object.Array{Elems: elems}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newErr("index operator not supported: %s", left.Type())
	}
}

// fits idx inside arrSize similar to python array indexing: x[-1] == x[len(x)-1]
// Return -1 sentinel for index out of bounds.
func getIndex(idx int64, arrSize int) int64 {
	m := int64(arrSize)
	if m == 0 {
		return -1
	}
	d := idx
	if idx < 0 {
		d = idx + m
	}
	if !(0 <= d && d < m) {
		return -1
	}
	return d
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrobj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	d := getIndex(idx, len(arrobj.Elems))
	if d < 0 {
		return NULL
	}
	return arrobj.Elems[d]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hsh := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newErr("unusable as hash key: %s", index.Type())
	}
	pair, ok := hsh.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt, env)

		switch r := res.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var r object.Object
	for _, statement := range block.Statements {
		r = Eval(statement, env)
		if r != nil {
			rt := r.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return r
			}
		}
	}
	return r
}

func evalPrefixExp(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExp(right)
	case "-":
		return evalMinusOpExp(right)
	default:
		return newErr("unknown operator: %s%s", op, right.Type())
	}
}

func evalInfixExp(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newErr("type mismatch: %s %s %s", left.Type(), op, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntInfixExp(op, left, right)

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBoolInfixExp(op, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExp(op, left, right)

	default:
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntInfixExp(op string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	var res int64
	switch op {
	case "+":
		res = l + r
	case "-":
		res = l - r
	case "*":
		res = l * r
	case "/":
		res = l / r
	case ">":
		return nativeBoolToBoolObj(l > r)
	case "<":
		return nativeBoolToBoolObj(l < r)
	case "==":
		return nativeBoolToBoolObj(l == r)
	case "!=":
		return nativeBoolToBoolObj(l != r)
	default:
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	return &object.Integer{Value: res}
}

func evalBoolInfixExp(op string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Boolean).Value
	r := right.(*object.Boolean).Value
	switch op {
	case "==":
		return nativeBoolToBoolObj(l == r)
	case "!=":
		return nativeBoolToBoolObj(l != r)
	default:
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExp(op string, left object.Object, right object.Object) object.Object {
	if op != "+" {
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	return &object.String{Value: l + r}
}

func evalMinusOpExp(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newErr("unknown operator: -%s", right.Type())
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalBangOperatorExp(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	v, ok := env.Get(node.Value)
	if ok {
		return v
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErr("identifier not found: %s", node.Value)
}

func nativeBoolToBoolObj(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}

func newErr(msg string, a ...any) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(msg, a...)}
}
