package object

import "fmt"

var BuiltIns = []struct {
	Name    string
	BuiltIn *BuiltIn
}{
	{
		"len",
		length,
	},
	{
		"puts",
		puts,
	},
	{
		"first",
		first,
	},
	{
		"last",
		last,
	},
	{
		"rest",
		rest,
	},
	{
		"append",
		appendBI,
	},
}

func GetBuiltInByName(name string) *BuiltIn {
	for _, def := range BuiltIns {
		if def.Name == name {
			return def.BuiltIn
		}
	}
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

var length = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *String:
			return &Integer{Value: int64(len(arg.Value))}
		case *Array:
			return &Integer{Value: int64(len(arg.Elements))}
		case *HashMap:
			return &Integer{Value: int64(len(arg.KVPairs))}
		default:
			return newError("argument to `len` is not supported, got %s", arg.Type())
		}
	},
}

var puts = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) == 0 {
			return newError("need at least one argument provided to puts")
		}

		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}

		return nil
	},
}

var first = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) == 0 {
				return newError("array is empty; no first element")
			}

			return arg.Elements[0]
		default:
			return newError("argument to `first` is not supported, got %s", arg.Type())
		}
	},
}

var last = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) == 0 {
				return newError("array is empty; no last element")
			}

			return arg.Elements[len(arg.Elements)-1]
		default:
			return newError("argument to `last` is not supported, got %s", arg.Type())
		}
	},
}

var rest = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) == 0 {
				return nil
			}

			restElements := make([]Object, len(arg.Elements)-1)
			copy(restElements, arg.Elements[1:len(arg.Elements)])
			return &Array{Elements: restElements}
		default:
			return newError("argument to `rest` is not supported, got %s", arg.Type())
		}
	},
}

var appendBI = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. expected=2, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *Array:
			newElements := make([]Object, len(arg.Elements), len(arg.Elements)+1)
			copy(newElements, arg.Elements)
			newElements = append(newElements, args[1])
			return &Array{Elements: newElements}
		default:
			return newError("argument to `append` is not supported, got %s", arg.Type())
		}
	},
}
