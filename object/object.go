package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/EmilLaursen/wiig/ast"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ                    = "NULL"
	RETURN_VALUE_OBJ            = "RETURN_VALUE"
	ERROR_OBJ                   = "ERROR"
	FUNCTION_OBJ                = "FUNCTION"
	STRING_OBJ                  = "STRING"
	BUILTIN_OBJ                 = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type String struct {
	Value string
}

func (i *String) Inspect() string  { return i.Value }
func (i *String) Type() ObjectType { return STRING_OBJ }

type Null struct{}

func (i *Null) Inspect() string  { return "null" }
func (i *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (i *ReturnValue) Inspect() string  { return i.Value.Inspect() }
func (i *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Builtin struct {
	Fn BuiltinFunction
}

func (i *Builtin) Inspect() string  { return "builtin function" }
func (i *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Error struct {
	Msg string
}

func (i *Error) Type() ObjectType { return ERROR_OBJ }
func (i *Error) Inspect() string  { return "ERROR: " + i.Msg }

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (*Function) Type() ObjectType { return FUNCTION_OBJ }
func (n *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range n.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(n.Body.String())
	out.WriteString("\n")
	return out.String()
}
