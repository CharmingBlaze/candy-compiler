package candy_evaluator

import (
	"candy/candy_ast"
	"strings"
)

func matchPattern(pat candy_ast.Expression, val *Value, bindings map[string]*Value, env *Env) bool {
	if pat == nil || val == nil {
		return false
	}

	switch p := pat.(type) {
	case *candy_ast.Identifier:
		if p.Value == "_" {
			return true
		}
		// Special case: check if it's a type name (int, string, etc.)
		switch strings.ToLower(p.Value) {
		case "int":
			return val.Kind == ValInt
		case "float":
			return val.Kind == ValFloat
		case "string":
			return val.Kind == ValString
		case "bool":
			return val.Kind == ValBool
		case "array":
			return val.Kind == ValArray
		case "map":
			return val.Kind == ValMap
		}
		// Bind to variable
		bindings[p.Value] = val
		return true

	case *candy_ast.IntegerLiteral:
		return val.Kind == ValInt && val.I64 == p.Value
	case *candy_ast.FloatLiteral:
		return val.Kind == ValFloat && val.F64 == p.Value
	case *candy_ast.StringLiteral:
		return val.Kind == ValString && val.Str == p.Value
	case *candy_ast.Boolean:
		return val.Kind == ValBool && val.B == p.Value
	case *candy_ast.NullLiteral:
		return val.Kind == ValNull

	case *candy_ast.ArrayLiteral:
		if val.Kind != ValArray {
			return false
		}
		// TODO: support ...rest
		if len(p.Elem) != len(val.Elems) {
			return false
		}
		for i, ep := range p.Elem {
			if !matchPattern(ep, &val.Elems[i], bindings, env) {
				return false
			}
		}
		return true

	case *candy_ast.MapLiteral:
		if val.Kind == ValMap {
			for _, pair := range p.Pairs {
				key, err := evalExpression(pair.Key, env)
				if err != nil {
					return false
				}
				ks, _ := mapKeyString(key)
				v, ok := val.StrMap[ks]
				if !ok {
					return false
				}
				if !matchPattern(pair.Value, &v, bindings, env) {
					return false
				}
			}
			return true
		}
		if val.Kind == ValStruct && val.St != nil {
			for _, pair := range p.Pairs {
				ks := candy_ast.ExprAsSimpleTypeName(pair.Key)
				v, ok := val.St.Data[ks]
				if !ok {
					return false
				}
				if !matchPattern(pair.Value, &v, bindings, env) {
					return false
				}
			}
			return true
		}
		return false
	}

	return false
}
