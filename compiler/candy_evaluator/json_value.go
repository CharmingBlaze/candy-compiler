package candy_evaluator

import "fmt"

// valueToJSONable converts a runtime value to a JSON-encodable Go value.
func valueToJSONable(v *Value) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch v.Kind {
	case ValNull:
		return nil, nil
	case ValBool:
		return v.B, nil
	case ValInt:
		return v.I64, nil
	case ValFloat:
		return v.F64, nil
	case ValString:
		return v.Str, nil
	case ValArray:
		a := make([]interface{}, 0, len(v.Elems))
		for i := range v.Elems {
			x, err := valueToJSONable(&v.Elems[i])
			if err != nil {
				return nil, err
			}
			a = append(a, x)
		}
		return a, nil
	case ValMap:
		if v.StrMap == nil {
			return map[string]interface{}{}, nil
		}
		m := make(map[string]interface{}, len(v.StrMap))
		for k, raw := range v.StrMap {
			inner, err := valueToJSONable(&raw)
			if err != nil {
				return nil, err
			}
			m[k] = inner
		}
		return m, nil
	default:
		return nil, fmt.Errorf("value cannot be JSON-encoded: %s", v.String())
	}
}

// jsonableToValue converts a decoded json.RawMessage or interface{} tree into a Value.
func jsonableToValue(x interface{}) (*Value, error) {
	switch t := x.(type) {
	case nil:
		return &Value{Kind: ValNull}, nil
	case bool:
		return &Value{Kind: ValBool, B: t}, nil
	case float64:
		// json numbers decode to float64
		f := t
		if f == float64(int64(f)) && f >= -9e15 && f <= 9e15 {
			return &Value{Kind: ValInt, I64: int64(f)}, nil
		}
		return &Value{Kind: ValFloat, F64: f}, nil
	case string:
		return &Value{Kind: ValString, Str: t}, nil
	case []interface{}:
		out := make([]Value, 0, len(t))
		for _, e := range t {
			iv, err := jsonableToValue(e)
			if err != nil {
				return nil, err
			}
			out = append(out, *iv)
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case map[string]interface{}:
		sm := make(map[string]Value, len(t))
		for k, e := range t {
			iv, err := jsonableToValue(e)
			if err != nil {
				return nil, err
			}
			sm[k] = *iv
		}
		return &Value{Kind: ValMap, StrMap: sm}, nil
	default:
		return &Value{Kind: ValString, Str: fmt.Sprint(t)}, nil
	}
}
