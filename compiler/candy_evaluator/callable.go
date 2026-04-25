package candy_evaluator

import "strings"

// InvokeCallable executes either a builtin function (by name string value)
// or a user-defined function value with the provided arguments.
func InvokeCallable(fn *Value, args []*Value) (*Value, error) {
	if fn == nil {
		return &Value{Kind: ValNull}, nil
	}
	if fn.Kind == ValString {
		if b, ok := Builtins[strings.ToLower(fn.Str)]; ok {
			return b(args)
		}
	}
	if fn.Kind == ValFunction {
		return evalUserFunction(fn, args)
	}
	return &Value{Kind: ValNull}, nil
}

