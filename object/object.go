package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/EmilLaursen/wiig/ast"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	STRING_OBJ       ObjectType = "STRING"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	HASH_OBJ         ObjectType = "HASH"
)

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func IsTypeOrNULL(one Object, ot ObjectType) bool {
	if one.Type() == ot {
		return true
	}
	return one.Type() == NULL_OBJ
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (i *Boolean) HashKey() HashKey {
	var v uint64
	if i.Value {
		v = 1
	} else {
		v = 0
	}
	return HashKey{Type: i.Type(), Value: v}
}

type String struct {
	Value string
}

func (i *String) Inspect() string  { return i.Value }
func (i *String) Type() ObjectType { return STRING_OBJ }
func (i *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(i.Value))
	return HashKey{Type: i.Type(), Value: h.Sum64()}
}

type Null struct{}

func (i *Null) Inspect() string  { return "null" }
func (i *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (i *ReturnValue) Inspect() string  { return i.Value.Inspect() }
func (i *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Array struct {
	Elems []Object
}

func (i *Array) Type() ObjectType { return ARRAY_OBJ }
func (i *Array) Inspect() string {
	var out bytes.Buffer
	elems := []string{}
	for _, e := range i.Elems {
		elems = append(elems, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

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
