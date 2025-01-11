package object

import (
	"fmt"
	"strings"
)

var BuiltIns = []struct {
	Name    string
	BuiltIn *BuiltIn
}{
	{
		"puts",
		puts,
	},
	{
		"len",
		length,
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
	{
		"join",
		join,
	},
	{
		"split",
		split,
	},
	{
		"sum",
		sum,
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

var puts = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) == 0 {
			return newError("need at least one argument provided to `puts`")
		}

		for _, arg := range args {
			fmt.Print(arg.Inspect())
		}
		fmt.Println()

		return nil
	},
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

var join = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 && len(args) != 2 {
			return newError("wrong number of arguments. expected 1 or 2, got=%d", len(args))
		}

		var delimiter string
		if len(args) == 2 {
			delimStr, ok := args[1].(*String)
			if !ok {
				return newError("expected delimiter passed to `join` to be a string, got %s", args[1].Type())
			}
			delimiter = delimStr.Value
		}

		switch arr := args[0].(type) {
		case *Array:
			var joined string

			for i, elem := range arr.Elements {
				elemStr, ok := elem.(*String)
				if !ok {
					return newError("elements of array passed to `join` must be strings, got %s", elem.Type())
				}

				joined += elemStr.Value
				if i < len(arr.Elements)-1 {
					joined += delimiter
				}
			}

			return &String{Value: joined}
		default:
			return newError("first argument to `join` must be an array, got %s", arr.Type())
		}
	},
}

var split = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 && len(args) != 2 {
			return newError("wrong number of arguments. expected 1 or 2, got=%d", len(args))
		}

		var separator string
		if len(args) == 2 {
			sepStr, ok := args[1].(*String)
			if !ok {
				return newError("expected separator passed to `split` to be a string, got %s", args[1].Type())
			}
			separator = sepStr.Value
		}

		switch s := args[0].(type) {
		case *String:
			resultStrings := strings.Split(s.Value, separator)

			resultObjects := []Object{}
			for _, resultStr := range resultStrings {
				resultObjects = append(resultObjects, &String{Value: resultStr})
			}

			return &Array{Elements: resultObjects}
		default:
			return newError("first argument to `split` must be a string, got %s", s.Type())
		}
	},
}

var sum = &BuiltIn{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. expected=1, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *Array:
			var arraySum float64
			isFloat := false

			for _, elem := range arg.Elements {
				elemValue, elemIsFloat, err := GetNumericalValue(elem)
				if err != nil {
					return &Error{Message: err.Error()}
				}

				arraySum += elemValue

				if elemIsFloat {
					isFloat = true
				}
			}

			if isFloat {
				return &Float{Value: arraySum}
			} else {
				return &Integer{Value: int64(arraySum)}
			}

		default:
			return newError("argument to `sum` is not supported, got %s", arg.Type())
		}
	},
}
