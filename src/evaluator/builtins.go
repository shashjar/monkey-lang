package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.BuiltIn{
	"puts":   puts,
	"len":    length,
	"first":  first,
	"last":   last,
	"rest":   rest,
	"append": appendBI,
}

var puts = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return newError("need at least one argument provided to puts")
		}

		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}

		return NULL
	},
}

var length = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` is not supported, got %s", arg.Type())
		}
	},
}

var first = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			if len(arg.Elements) == 0 {
				return newError("array is empty; no first element")
			}

			return arg.Elements[0]
		default:
			return newError("argument to `first` is not supported, got %s", arg.Type())
		}
	},
}

var last = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			if len(arg.Elements) == 0 {
				return newError("array is empty; no last element")
			}

			return arg.Elements[len(arg.Elements)-1]
		default:
			return newError("argument to `last` is not supported, got %s", arg.Type())
		}
	},
}

var rest = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			if len(arg.Elements) == 0 {
				return NULL
			}

			restElements := make([]object.Object, len(arg.Elements)-1)
			copy(restElements, arg.Elements[1:len(arg.Elements)])
			return &object.Array{Elements: restElements}
		default:
			return newError("argument to `rest` is not supported, got %s", arg.Type())
		}
	},
}

var appendBI = &object.BuiltIn{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. expected=2, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			newElements := make([]object.Object, len(arg.Elements), len(arg.Elements)+1)
			copy(newElements, arg.Elements)
			newElements = append(newElements, args[1])
			return &object.Array{Elements: newElements}
		default:
			return newError("argument to `append` is not supported, got %s", arg.Type())
		}
	},
}
