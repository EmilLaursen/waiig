package eval

import (
	"fmt"

	"github.com/EmilLaursen/wiig/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErr("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elems))}
			default:
				return newErr("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErr("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newErr("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			n := len(arr.Elems)
			nels := make([]object.Object, n+1, n+1)
			copy(nels, arr.Elems)
			nels[n] = args[1]
			arobj := &object.Array{Elems: nels}
			return arobj
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
