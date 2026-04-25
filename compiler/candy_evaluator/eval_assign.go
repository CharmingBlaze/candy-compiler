package candy_evaluator

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

func valueToValue(v *Value) Value {
	if v == nil {
		return Value{Kind: ValNull}
	}
	return *v
}

// evalAssign handles = for identifiers, a[i], and obj.field (instance = struct { ... }).
func evalAssign(t *candy_ast.AssignExpression, e *Env) (*Value, error) {
	rhs, err := evalExpression(t.Value, e)
	if err != nil {
		return nil, err
	}
	if rhs == nil {
		rhs = &Value{Kind: ValNull}
	}
	op := t.Operator
	if op == "" {
		op = "="
	}
	if op != "=" {
		baseOp, ok := mapAssignToInfix(op)
		if !ok {
			return nil, fmt.Errorf("unsupported assignment operator %q", op)
		}
		lv, err := evalExpression(t.Left, e)
		if err != nil {
			return nil, err
		}
		rhs, err = evalInfix(baseOp, lv, rhs, t)
		if err != nil {
			return nil, err
		}
	}

	switch l := t.Left.(type) {
	case *candy_ast.Identifier:
		if !e.Update(l.Value, rhs) {
			e.Set(l.Value, rhs)
		}
		return rhs, nil
	case *candy_ast.TupleLiteral:
		if rhs.Kind != ValArray {
			return nil, fmt.Errorf("multiple assignment: right hand side must be a tuple/array")
		}
		if len(l.Elems) != len(rhs.Elems) {
			return nil, fmt.Errorf("multiple assignment: size mismatch (%d vs %d)", len(l.Elems), len(rhs.Elems))
		}
		for i, el := range l.Elems {
			id, ok := el.(*candy_ast.Identifier)
			if !ok {
				return nil, fmt.Errorf("multiple assignment: left side must be identifiers")
			}
			if id.Value == "_" {
				continue
			}
			e.Set(id.Value, &rhs.Elems[i])
		}
		return rhs, nil
	case *candy_ast.IndexExpression:
		return evalAssignIndex(l, rhs, e)
	case *candy_ast.DotExpression:
		return evalAssignDot(l, rhs, e)
	default:
		return nil, fmt.Errorf("eval: assignment not supported for %T", t.Left)
	}
}

func mapAssignToInfix(op string) (string, bool) {
	switch op {
	case "+=":
		return "+", true
	case "-=":
		return "-", true
	case "*=":
		return "*", true
	case "/=":
		return "/", true
	default:
		return "", false
	}
}

func evalAssignIndex(t *candy_ast.IndexExpression, rhs *Value, e *Env) (*Value, error) {
	id, ok := t.Base.(*candy_ast.Identifier)
	if !ok {
		return nil, &RuntimeError{Msg: "index assign: only name[index] is supported in the dynamic runtime"}
	}
	place, ok2 := e.Get(id.Value)
	if !ok2 || place == nil {
		return nil, &RuntimeError{Msg: "undefined: " + id.Value}
	}
	ix, err := evalExpression(t.Index, e)
	if err != nil {
		return nil, err
	}
	if ix == nil {
		return nil, &RuntimeError{Msg: "index is null"}
	}
	vv := valueToValue(rhs)
	switch {
	case place.Kind == ValArray && ix.Kind == ValInt:
		i := int(ix.I64)
		if i < 0 {
			i = len(place.Elems) + i
		}
		if i < 0 || i >= len(place.Elems) {
			return nil, &RuntimeError{Msg: "index out of range in assignment"}
		}
		place.Elems[i] = vv
		return rhs, nil
	case place.Kind == ValMap && place.StrMap != nil:
		ks, kerr := mapKeyString(ix)
		if kerr != nil {
			return nil, kerr
		}
		place.StrMap[ks] = vv
		return rhs, nil
	}
	return nil, &RuntimeError{Msg: "index assign: need array or map"}
}

func evalAssignDot(t *candy_ast.DotExpression, rhs *Value, e *Env) (*Value, error) {
	var place *Value
	if id, ok := t.Left.(*candy_ast.Identifier); ok {
		p, ok2 := e.Get(id.Value)
		if !ok2 || p == nil {
			return nil, &RuntimeError{Msg: "undefined: " + id.Value}
		}
		place = p
	} else {
		p, err := evalExpression(t.Left, e)
		if err != nil {
			return nil, err
		}
		place = p
	}
	if (place == nil || place.Kind == ValNull) && t.IsSafe {
		return &Value{Kind: ValNull}, nil
	}
	if place.Kind != ValStruct || place.St == nil {
		return nil, &RuntimeError{Msg: "dot assign: not a struct instance"}
	}
	key := t.Right.Value
	if place.St.Data == nil {
		place.St.Data = make(map[string]Value)
	}
	// Match existing key case
	if _, has := place.St.Data[key]; !has {
		for k := range place.St.Data {
			if strings.EqualFold(k, key) {
				key = k
				break
			}
		}
	}
	place.St.Data[key] = valueToValue(rhs)
	return rhs, nil
}
