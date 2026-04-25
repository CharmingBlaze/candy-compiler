package candy_evaluator

import (
	"candy/candy_ast"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func lookupModFn(m *moduleVal, name string) (func([]*Value) (*Value, error), bool) {
	ln := strings.ToLower(name)
	if m.Fns != nil {
		for k, f := range m.Fns {
			if strings.ToLower(k) == ln {
				return f, true
			}
		}
	}
	return nil, false
}

func evalArgs(argExprs []candy_ast.Expression, e *Env) ([]*Value, error) {
	var args []*Value
	for _, a := range argExprs {
		av, err2 := evalExpression(a, e)
		if err2 != nil {
			return nil, err2
		}
		args = append(args, av)
	}
	return args, nil
}

func callArrayMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	elems := recv.Elems
	ln := strings.ToLower(name)
	switch ln {
	case "contains", "include":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "contains: one value"}
		}
		for _, e0 := range elems {
			if valueEqual(&e0, args[0]) {
				return &Value{Kind: ValBool, B: true}, nil
			}
		}
		return &Value{Kind: ValBool, B: false}, nil
	case "index_of", "indexof", "index":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "index_of: one value"}
		}
		for i := range elems {
			if valueEqual(&elems[i], args[0]) {
				return &Value{Kind: ValInt, I64: int64(i)}, nil
			}
		}
		return &Value{Kind: ValInt, I64: -1}, nil
	case "is_empty", "isempty", "empty":
		return &Value{Kind: ValBool, B: len(elems) == 0}, nil
	case "length", "count", "size":
		return &Value{Kind: ValInt, I64: int64(len(elems))}, nil
	case "first", "head":
		if len(elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		return ptrVal(elems[0]), nil
	case "last", "tail":
		if len(elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		return ptrVal(elems[len(elems)-1]), nil
	case "clear", "empty_all":
		recv.Elems = recv.Elems[:0] // in-place: receiver must be shared *Value
		return recv, nil
	case "shuffle":
		rg := rand.New(rand.NewSource(time.Now().UnixNano()))
		ix := make([]int, len(recv.Elems))
		for i := range ix {
			ix[i] = i
		}
		rg.Shuffle(len(ix), func(i, j int) { ix[i], ix[j] = ix[j], ix[i] })
		ne := make([]Value, len(recv.Elems))
		for to, from := range ix {
			ne[to] = recv.Elems[from]
		}
		recv.Elems = ne
		return recv, nil
	case "add", "push", "append":
		for _, a := range args {
			recv.Elems = append(recv.Elems, valueToValue(a))
		}
		return recv, nil
	case "insert":
		if len(args) != 2 {
			return nil, &RuntimeError{Msg: "insert: (index, value)"}
		}
		ix, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "insert: int index"}
		}
		i := int(ix)
		if i < 0 {
			i = len(recv.Elems) + i
		}
		if i < 0 || i > len(recv.Elems) {
			return nil, &RuntimeError{Msg: "insert: index out of range"}
		}
		v := valueToValue(args[1])
		recv.Elems = append(recv.Elems, Value{})
		copy(recv.Elems[i+1:], recv.Elems[i:])
		recv.Elems[i] = v
		return recv, nil
	case "remove":
		// Candy: scores.remove(2) removes the element at index 2 (third item, 0-based).
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove: one index"}
		}
		ix, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "remove: int index"}
		}
		i := int(ix)
		if i < 0 {
			i = len(recv.Elems) + i
		}
		if i < 0 || i >= len(recv.Elems) {
			return nil, &RuntimeError{Msg: "remove: index out of range"}
		}
		recv.Elems = append(recv.Elems[:i], recv.Elems[i+1:]...)
		return recv, nil
	case "remove_first", "remove_value", "removevalue":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "removeFirst: one value to delete (first match)"}
		}
		for j := range recv.Elems {
			if valueEqual(&recv.Elems[j], args[0]) {
				recv.Elems = append(recv.Elems[:j], recv.Elems[j+1:]...)
				return recv, nil
			}
		}
		return recv, nil
	case "remove_last", "removelast", "pop":
		if len(recv.Elems) == 0 {
			return recv, nil
		}
		recv.Elems = recv.Elems[:len(recv.Elems)-1]
		return recv, nil
	case "remove_at", "delete_at", "splice1":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove_at: one index"}
		}
		ix, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "remove_at: int index"}
		}
		i := int(ix)
		if i < 0 {
			i = len(recv.Elems) + i
		}
		if i < 0 || i >= len(recv.Elems) {
			return nil, &RuntimeError{Msg: "remove_at: index out of range"}
		}
		recv.Elems = append(recv.Elems[:i], recv.Elems[i+1:]...)
		return recv, nil
	case "sort":
		sort.Slice(recv.Elems, func(i, j int) bool {
			return valueOrderLessForSort(&recv.Elems[i], &recv.Elems[j])
		})
		return recv, nil
	case "reverse":
		for l, r := 0, len(recv.Elems)-1; l < r; l, r = l+1, r-1 {
			recv.Elems[l], recv.Elems[r] = recv.Elems[r], recv.Elems[l]
		}
		return recv, nil
	case "sum":
		if len(recv.Elems) == 0 {
			return &Value{Kind: ValInt, I64: 0}, nil
		}
		allInt := true
		var isum int64
		var fsum float64
		for i := range recv.Elems {
			v := recv.Elems[i]
			switch v.Kind {
			case ValInt:
				if allInt {
					isum += v.I64
				} else {
					fsum += float64(v.I64)
				}
			case ValFloat:
				if allInt {
					fsum = float64(isum) + v.F64
					isum = 0
					allInt = false
				} else {
					fsum += v.F64
				}
			default:
				return nil, &RuntimeError{Msg: "sum: need numeric elements"}
			}
		}
		if allInt {
			return &Value{Kind: ValInt, I64: isum}, nil
		}
		return &Value{Kind: ValFloat, F64: fsum}, nil
	case "max", "min":
		if len(recv.Elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		best := recv.Elems[0]
		isMax := strings.EqualFold(ln, "max")
		for i := 1; i < len(recv.Elems); i++ {
			cur := recv.Elems[i]
			if isMax {
				if valueOrderLessForSort(&best, &cur) {
					best = cur
				}
			} else {
				if valueOrderLessForSort(&cur, &best) {
					best = cur
				}
			}
		}
		return ptrVal(best), nil
	case "join":
		if len(args) < 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "join: separator string"}
		}
		sep := args[0].Str
		var b strings.Builder
		for j, v := range recv.Elems {
			if j > 0 {
				b.WriteString(sep)
			}
			p := ptrVal(v)
			if p == nil {
				b.WriteString("null")
			} else {
				b.WriteString(p.String())
			}
		}
		return &Value{Kind: ValString, Str: b.String()}, nil
	case "map":
		if len(args) != 1 || args[0] == nil {
			return nil, &RuntimeError{Msg: "map: one function"}
		}
		if args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "map: need function (e.g. (x) => x + 1)"}
		}
		out := make([]Value, 0, len(recv.Elems))
		for i := range recv.Elems {
			ev := &recv.Elems[i]
			mv, err := evalUserFunction(args[0], []*Value{ptrVal(*ev)})
			if err != nil {
				return nil, err
			}
			out = append(out, valueToValue(mv))
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "filter":
		if len(args) != 1 || args[0] == nil {
			return nil, &RuntimeError{Msg: "filter: one function"}
		}
		if args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "filter: need function (e.g. (n) => n > 0)"}
		}
		var out []Value
		for i := range recv.Elems {
			ev := &recv.Elems[i]
			fv, err := evalUserFunction(args[0], []*Value{ptrVal(*ev)})
			if err != nil {
				return nil, err
			}
			if fv.Truthy() {
				out = append(out, recv.Elems[i])
			}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "reduce":
		if len(args) != 2 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "reduce: (function, initial)"}
		}
		acc := args[1]
		for i := range recv.Elems {
			nv, err := evalUserFunction(args[0], []*Value{acc, ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			acc = nv
		}
		return acc, nil
	case "find":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "find: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if okv.Truthy() {
				return ptrVal(recv.Elems[i]), nil
			}
		}
		return &Value{Kind: ValNull}, nil
	case "all":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "all: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if !okv.Truthy() {
				return &Value{Kind: ValBool, B: false}, nil
			}
		}
		return &Value{Kind: ValBool, B: true}, nil
	case "any":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "any: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if okv.Truthy() {
				return &Value{Kind: ValBool, B: true}, nil
			}
		}
		return &Value{Kind: ValBool, B: false}, nil
	case "unique":
		out := make([]Value, 0, len(recv.Elems))
		for i := range recv.Elems {
			dup := false
			for j := range out {
				if valueEqual(&out[j], &recv.Elems[i]) {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, recv.Elems[i])
			}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	default:
		return nil, &RuntimeError{Msg: "array has no method: " + name}
	}
}

func callMapMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	if recv.StrMap == nil {
		recv.StrMap = make(map[string]Value)
	}
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	m := recv.StrMap
	ln := strings.ToLower(name)
	switch ln {
	case "keys", "key_list":
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		out := make([]Value, len(ks))
		for i, k := range ks {
			out[i] = Value{Kind: ValString, Str: k}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "values", "vals":
		// stable order: sort by key
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		out := make([]Value, len(ks))
		for i, k := range ks {
			out[i] = m[k]
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "has", "contains", "include":
		if len(args) < 1 {
			return nil, &RuntimeError{Msg: "has: one key"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		_, ok := m[ks]
		if !ok {
			// case-insensitive
			for km := range m {
				if strings.EqualFold(km, ks) {
					ok = true
					break
				}
			}
		}
		return &Value{Kind: ValBool, B: ok}, nil
	case "get", "at":
		if len(args) < 1 {
			return nil, &RuntimeError{Msg: "get: key, optional default"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		if v, ok := m[ks]; ok {
			return ptrVal(v), nil
		}
		for km, v := range m {
			if strings.EqualFold(km, ks) {
				return ptrVal(v), nil
			}
		}
		if len(args) >= 2 {
			return args[1], nil
		}
		return &Value{Kind: ValNull}, nil
	case "remove", "delete", "rm":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove: one key"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		if _, ok := m[ks]; ok {
			delete(m, ks)
			return recv, nil
		}
		for k := range m {
			if strings.EqualFold(k, ks) {
				delete(m, k)
				return recv, nil
			}
		}
		return recv, nil
	case "clear", "empty":
		for k := range m {
			delete(m, k)
		}
		return recv, nil
	case "merge", "put_all", "putall":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "merge: one map"}
		}
		if args[0] == nil || args[0].Kind != ValMap || args[0].StrMap == nil {
			return nil, &RuntimeError{Msg: "merge: need map"}
		}
		for k, v := range args[0].StrMap {
			m[k] = v
		}
		return recv, nil
	default:
		return nil, &RuntimeError{Msg: "map has no method: " + name}
	}
}

// asInt64Value coerces a value to int index (for list insert/remove).
func asInt64Value(v *Value) (int64, bool) {
	if v == nil {
		return 0, true
	}
	if v.Kind == ValInt {
		return v.I64, true
	}
	if v.Kind == ValFloat {
		return int64(v.F64), true
	}
	return 0, false
}

// valueOrderLessForSort is a total order for sort / max / min: numeric when possible, else String().
func valueOrderLessForSort(a, b *Value) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	lf, rf, useF := toFloats(a, b)
	if (a.Kind == ValInt || a.Kind == ValFloat) && (b.Kind == ValInt || b.Kind == ValFloat) {
		if useF {
			return lf < rf
		}
		if a.Kind == ValInt && b.Kind == ValInt {
			return a.I64 < b.I64
		}
		return lf < rf
	}
	if a.Kind == ValString && b.Kind == ValString {
		return a.Str < b.Str
	}
	return a.String() < b.String()
}

func callStringMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	if recv == nil || recv.Kind != ValString {
		return nil, &RuntimeError{Msg: "string method: need string"}
	}
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	s := recv.Str
	ln := strings.ToLower(name)
	switch ln {
	case "upper", "toupper", "to_upper":
		return &Value{Kind: ValString, Str: strings.ToUpper(s)}, nil
	case "lower", "tolower", "to_lower":
		return &Value{Kind: ValString, Str: strings.ToLower(s)}, nil
	case "split":
		sep := " "
		if len(args) >= 1 && args[0] != nil && args[0].Kind == ValString {
			sep = args[0].Str
		}
		parts := strings.Split(s, sep)
		out := make([]Value, len(parts))
		for i, p := range parts {
			out[i] = Value{Kind: ValString, Str: p}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "contains", "includes":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "contains: substring string"}
		}
		return &Value{Kind: ValBool, B: strings.Contains(s, args[0].Str)}, nil
	case "starts_with", "startswith", "prefix":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "starts_with: string prefix"}
		}
		return &Value{Kind: ValBool, B: strings.HasPrefix(s, args[0].Str)}, nil
	case "ends_with", "endswith", "suffix":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "ends_with: string suffix"}
		}
		return &Value{Kind: ValBool, B: strings.HasSuffix(s, args[0].Str)}, nil
	case "trim", "strip", "trim_space":
		return &Value{Kind: ValString, Str: strings.TrimSpace(s)}, nil
	case "replace":
		if len(args) < 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValString || args[1].Kind != ValString {
			return nil, &RuntimeError{Msg: "replace: (old, new) strings"}
		}
		return &Value{Kind: ValString, Str: strings.ReplaceAll(s, args[0].Str, args[1].Str)}, nil
	case "index_of", "indexof":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "indexOf: substring string"}
		}
		return &Value{Kind: ValInt, I64: int64(strings.Index(s, args[0].Str))}, nil
	case "substring":
		if len(args) < 1 || len(args) > 2 {
			return nil, &RuntimeError{Msg: "substring: start, [end]"}
		}
		start, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "substring: start int"}
		}
		runes := []rune(s)
		st := int(start)
		if st < 0 {
			st = len(runes) + st
		}
		if st < 0 {
			st = 0
		}
		if st > len(runes) {
			st = len(runes)
		}
		en := len(runes)
		if len(args) == 2 {
			end, ok := asInt64Value(args[1])
			if !ok {
				return nil, &RuntimeError{Msg: "substring: end int"}
			}
			en = int(end)
			if en < 0 {
				en = len(runes) + en
			}
			if en < 0 {
				en = 0
			}
			if en > len(runes) {
				en = len(runes)
			}
		}
		if en < st {
			en = st
		}
		return &Value{Kind: ValString, Str: string(runes[st:en])}, nil
	default:
		return nil, &RuntimeError{Msg: "string has no method: " + name}
	}
}
