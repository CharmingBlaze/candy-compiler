package candy_evaluator

import "candy/candy_ast"

// evalForIn runs `for v in iterable { body }` (BASIC-style and SimpleC "for x in y").
func evalForIn(t *candy_ast.ForStatement, e *Env) (any, error) {
	if t.Iterable == nil || t.Body == nil || t.Var == nil {
		return nil, nil
	}
	it, err := evalExpression(t.Iterable, e)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, nil
	}
	name := t.Var.Value
	switch it.Kind {
	case ValArray:
		for i, el := range it.Elems {
			ne := newIsolatedLoopEnv(e)
			ee := el
			ne.Set(name, ptrVal(ee))
			if t.ValueVar != nil {
				ne.Set(t.ValueVar.Value, &Value{Kind: ValInt, I64: int64(i)})
			}
			if res, err2 := runBlockInEnv(t.Body, ne); err2 != nil {
				return nil, err2
			} else if rw, ok2 := res.(ReturnWrap); ok2 {
				return rw, nil
			} else if _, ok3 := res.(BreakWrap); ok3 {
				break
			} else if _, ok4 := res.(ContinueWrap); ok4 {
				continue
			}
		}
		return nil, nil
	case ValString:
		idx := 0
		for _, ru := range it.Str {
			ne := newIsolatedLoopEnv(e)
			c := &Value{Kind: ValString, Str: string(ru)}
			ne.Set(name, c)
			if t.ValueVar != nil {
				ne.Set(t.ValueVar.Value, &Value{Kind: ValInt, I64: int64(idx)})
			}
			idx++
			if res, err2 := runBlockInEnv(t.Body, ne); err2 != nil {
				return nil, err2
			} else if rw, ok2 := res.(ReturnWrap); ok2 {
				return rw, nil
			} else if _, ok3 := res.(BreakWrap); ok3 {
				break
			} else if _, ok4 := res.(ContinueWrap); ok4 {
				continue
			}
		}
		return nil, nil
	case ValMap:
		if it.StrMap == nil {
			return nil, nil
		}
		for k, v := range it.StrMap {
			ne := newIsolatedLoopEnv(e)
			ne.Set(name, &Value{Kind: ValString, Str: k})
			if t.ValueVar != nil {
				vv := v
				ne.Set(t.ValueVar.Value, &vv)
			}
			if res, err2 := runBlockInEnv(t.Body, ne); err2 != nil {
				return nil, err2
			} else if rw, ok2 := res.(ReturnWrap); ok2 {
				return rw, nil
			} else if _, ok3 := res.(BreakWrap); ok3 {
				break
			} else if _, ok4 := res.(ContinueWrap); ok4 {
				continue
			}
		}
		return nil, nil
	default:
		return nil, &RuntimeError{Msg: "for-in: not iterable (need array, string, or map)"}
	}
}

// evalForTo runs `for i = a to b [step s] { ... }` (integer range, inclusive of endpoints).
func evalForTo(t *candy_ast.ForStatement, e *Env) (any, error) {
	if t.Start == nil || t.End == nil || t.Body == nil || t.Var == nil {
		return nil, nil
	}
	sv, err := evalExpression(t.Start, e)
	if err != nil {
		return nil, err
	}
	ev, err := evalExpression(t.End, e)
	if err != nil {
		return nil, err
	}
	a, aok := asInt64(sv)
	b, bok := asInt64(ev)
	if !aok || !bok {
		return nil, &RuntimeError{Msg: "for to: need integer start/end"}
	}
	var step int64 = 1
	if t.Step != nil {
		st, err2 := evalExpression(t.Step, e)
		if err2 != nil {
			return nil, err2
		}
		s, ok3 := asInt64(st)
		if !ok3 || s == 0 {
			return nil, &RuntimeError{Msg: "for to: need non-zero integer step"}
		}
		step = s
	}
	vn := t.Var.Value
	if step > 0 {
		for i := a; i <= b; i += step {
			ne := e.NewEnclosed()
			ne.Set(vn, &Value{Kind: ValInt, I64: i})
			if res, err2 := runBlockInEnv(t.Body, ne); err2 != nil {
				return nil, err2
			} else if rw, ok2 := res.(ReturnWrap); ok2 {
				return rw, nil
			} else if _, ok3 := res.(BreakWrap); ok3 {
				break
			} else if _, ok4 := res.(ContinueWrap); ok4 {
				continue
			}
		}
	} else {
		for i := a; i >= b; i += step {
			ne := e.NewEnclosed()
			ne.Set(vn, &Value{Kind: ValInt, I64: i})
			if res, err2 := runBlockInEnv(t.Body, ne); err2 != nil {
				return nil, err2
			} else if rw, ok2 := res.(ReturnWrap); ok2 {
				return rw, nil
			} else if _, ok3 := res.(BreakWrap); ok3 {
				break
			} else if _, ok4 := res.(ContinueWrap); ok4 {
				continue
			}
		}
	}
	return nil, nil
}

func asInt64(v *Value) (int64, bool) {
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

func evalWhile(t *candy_ast.WhileStatement, e *Env) (any, error) {
	if t == nil {
		return nil, nil
	}
	for {
		cond, err := evalExpression(t.Condition, e)
		if err != nil {
			return nil, err
		}
		if cond == nil || !cond.Truthy() {
			return nil, nil
		}
		if t.Body == nil {
			return nil, nil
		}
		// `WhileStatement.Body` is a block; run it in the same env (not an extra `BlockStatement` scope).
		r, err2 := runBlockInEnv(t.Body, e)
		if err2 != nil {
			return nil, err2
		}
		if rw, ok := r.(ReturnWrap); ok {
			return rw, nil
		} else if _, ok2 := r.(BreakWrap); ok2 {
			break
		} else if _, ok3 := r.(ContinueWrap); ok3 {
			continue
		}
	}
	return nil, nil
}

func runBlockInEnv(b *candy_ast.BlockStatement, env *Env) (any, error) {
	if b == nil {
		return nil, nil
	}
	for _, s := range b.Statements {
		if s == nil {
			continue
		}
		r, err := evalStatement(s, env)
		if err != nil {
			return nil, err
		}
		if rw, ok := r.(ReturnWrap); ok {
			return rw, nil
		}
		if bw, ok := r.(BreakWrap); ok {
			return bw, nil
		}
		if cw, ok := r.(ContinueWrap); ok {
			return cw, nil
		}
	}
	return nil, nil
}

// evalCFor runs `for (init; cond; post) { body }`.
func evalCFor(t *candy_ast.CForStatement, e *Env) (any, error) {
	if t == nil {
		return nil, nil
	}
	ne := e.NewEnclosed()
	if t.Init != nil {
		if _, err := evalStatement(t.Init, ne); err != nil {
			return nil, err
		}
	}
	for {
		if t.Cond != nil {
			cv, err := evalExpression(t.Cond, ne)
			if err != nil {
				return nil, err
			}
			if !cv.Truthy() {
				return nil, nil
			}
		}
		if t.Body != nil {
			r, err := runBlockInEnv(t.Body, ne)
			if err != nil {
				return nil, err
			}
			if rw, ok := r.(ReturnWrap); ok {
				return rw, nil
			}
			if _, ok := r.(BreakWrap); ok {
				return nil, nil
			}
			if _, ok := r.(ContinueWrap); ok {
				goto doPost
			}
		}
	doPost:
		if t.Post != nil {
			if _, err := evalExpression(t.Post, ne); err != nil {
				return nil, err
			}
		}
	}
}

func evalForEach(t *candy_ast.ForEachStatement, e *Env) (any, error) {
	if t.Iterable == nil || t.Body == nil || t.Var == nil {
		return nil, nil
	}
	it, err := evalExpression(t.Iterable, e)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, nil
	}
	name := t.Var.Value
	if it.Kind != ValArray && it.Kind != ValString && it.Kind != ValMap {
		return nil, &RuntimeError{Msg: "foreach: not iterable"}
	}

	// Helper for loop body execution
	runBody := func(v *Value) (any, error) {
		ne := newIsolatedLoopEnv(e)
		ne.Set(name, v)
		return runBlockInEnv(t.Body, ne)
	}

	switch it.Kind {
	case ValArray:
		for _, el := range it.Elems {
			ee := el
			res, err2 := runBody(&ee)
			if err2 != nil {
				return nil, err2
			}
			if rw, ok := res.(ReturnWrap); ok {
				return rw, nil
			}
			if _, ok := res.(BreakWrap); ok {
				break
			}
			if _, ok := res.(ContinueWrap); ok {
				continue
			}
		}
	case ValString:
		for _, ru := range it.Str {
			res, err2 := runBody(&Value{Kind: ValString, Str: string(ru)})
			if err2 != nil {
				return nil, err2
			}
			if rw, ok := res.(ReturnWrap); ok {
				return rw, nil
			}
			if _, ok := res.(BreakWrap); ok {
				break
			}
			if _, ok := res.(ContinueWrap); ok {
				continue
			}
		}
	case ValMap:
		for _, v := range it.StrMap {
			vv := v
			res, err2 := runBody(&vv)
			if err2 != nil {
				return nil, err2
			}
			if rw, ok := res.(ReturnWrap); ok {
				return rw, nil
			}
			if _, ok := res.(BreakWrap); ok {
				break
			}
			if _, ok := res.(ContinueWrap); ok {
				continue
			}
		}
	}
	return nil, nil
}

func snapshotVisibleBindings(e *Env) map[string]*Value {
	out := make(map[string]*Value)
	var chain []*Env
	for cur := e; cur != nil; cur = cur.Parent {
		chain = append(chain, cur)
	}
	for i := len(chain) - 1; i >= 0; i-- {
		for k, v := range chain[i].Store {
			if v == nil {
				out[k] = &Value{Kind: ValNull}
				continue
			}
			vv := *v
			out[k] = &vv
		}
	}
	return out
}

func newIsolatedLoopEnv(parent *Env) *Env {
	ne := parent.NewEnclosed()
	// Snapshot visible bindings into the loop-body scope so assignment
	// stays iteration-local instead of mutating outer bindings.
	for k, vv := range snapshotVisibleBindings(parent) {
		ne.Set(k, vv)
	}
	return ne
}

func evalRepeat(t *candy_ast.RepeatStatement, e *Env) (any, error) {
	if t.Count == nil || t.Body == nil {
		return nil, nil
	}
	cv, err := evalExpression(t.Count, e)
	if err != nil {
		return nil, err
	}
	count, ok := asInt64(cv)
	if !ok {
		return nil, &RuntimeError{Msg: "repeat: need integer count"}
	}
	for i := int64(0); i < count; i++ {
		res, err2 := runBlockInEnv(t.Body, e)
		if err2 != nil {
			return nil, err2
		}
		if rw, ok := res.(ReturnWrap); ok {
			return rw, nil
		}
		if _, ok := res.(BreakWrap); ok {
			break
		}
		if _, ok := res.(ContinueWrap); ok {
			continue
		}
	}
	return nil, nil
}

// evalLoop runs `loop { ... }`. If Raylib is linked, exits when the window should close.
func evalLoop(t *candy_ast.LoopStatement, e *Env) (any, error) {
	if t.Body == nil {
		return nil, nil
	}
	for {
		if fn, ok := Builtins["shouldclose"]; ok {
			v, err := fn(nil)
			if err == nil && v != nil && v.Truthy() {
				break
			}
		}
		res, err := runBlockInEnv(t.Body, e)
		if err != nil {
			return nil, err
		}
		if rw, ok := res.(ReturnWrap); ok {
			return rw, nil
		}
		if _, ok := res.(BreakWrap); ok {
			break
		}
		if _, ok := res.(ContinueWrap); ok {
			continue
		}
	}
	return nil, nil
}
