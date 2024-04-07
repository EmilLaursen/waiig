package eval

import (
	"github.com/EmilLaursen/wiig/ast"
	"github.com/EmilLaursen/wiig/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.IfExpression:
		return evalIfExp(node)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBoolObj(node.Value)
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExp(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExp(node.Operator, left, right)
	}
	return nil
}

func evalIfExp(n *ast.IfExpression) object.Object {
	cond := Eval(n.Condition)

	if isTruthy(cond) {
		return Eval(n.Consequence)
	} else if n.Alternative != nil {
		return Eval(n.Alternative)
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

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt)
	}
	return res
}

func evalPrefixExp(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExp(right)
	case "-":
		return evalMinusOpExp(right)
	default:
		return NULL
	}
}

func evalInfixExp(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntInfixExp(op, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBoolInfixExp(op, left, right)

	default:
		return NULL
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
		return NULL
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
		return NULL
	}
}

func evalMinusOpExp(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
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

func nativeBoolToBoolObj(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}
