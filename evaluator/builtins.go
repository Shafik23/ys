package evaluator

import "github.com/shafik23/ys/object"

var builtins = map[string]*object.Builtin{
	// Go compiler can infer the type of the struct literal from the the decalaration above.
	"len": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		// case *object.Array:
		// 	return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported, got type %s", args[0].Type())
		}
	}},
}
