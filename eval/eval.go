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
		if node.Value {
			return TRUE
		}
		return FALSE

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
	function, ok := fn.(*object.Function)
	if !ok {
		return newErr("not a function: %s", fn.Type())
	}

	scope := object.NewScope(function.Env)
	for i, p := range function.Params {
		scope.Set(p.Value, args[i])
	}
	ret := Eval(function.Body, scope)
	return unwrapReturn(ret)
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
	if !ok {
		return newErr("identifier not found: %s", node.Value)
	}
	return v
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
