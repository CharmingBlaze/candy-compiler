package candy_evaluator

import (
	"candy/candy_ast"
	"strings"
)

func evalTry(t *candy_ast.TryStatement, e *Env) (any, error) {
	if t == nil {
		return nil, nil
	}
	ne := e.NewEnclosed()
	var res any
	var err0 error
	if t.TryBody != nil {
		res, err0 = runBlockWithValue(t.TryBody, ne)
	}
	if rw, ok := res.(ReturnWrap); ok {
		if t.FinallyBody != nil {
			_, fe := runBlockInEnv(t.FinallyBody, e)
			if fe != nil {
				return nil, fe
			}
		}
		return rw, nil
	}
	if err0 == nil {
		if t.FinallyBody != nil {
			_, fe := runBlockInEnv(t.FinallyBody, e)
			if fe != nil {
				return nil, fe
			}
		}
		return res, nil
	}
	for _, cc := range t.CatchClauses {
		if cc == nil || !catchClauseApplies(cc, err0) {
			continue
		}
		cne := e.NewEnclosed()
		if cc.Identifier != nil {
			cne.Set(cc.Identifier.Value, &Value{Kind: ValString, Str: err0.Error()})
		}
		var cres any
		var err1 error
		if cc.Body != nil {
			cres, err1 = runBlockWithValue(cc.Body, cne)
		}
		if err1 != nil {
			if t.FinallyBody != nil {
				_, fe := runBlockInEnv(t.FinallyBody, e)
				if fe != nil {
					return nil, fe
				}
			}
			return nil, err1
		}
		if rw, ok2 := cres.(ReturnWrap); ok2 {
			if t.FinallyBody != nil {
				_, fe := runBlockInEnv(t.FinallyBody, e)
				if fe != nil {
					return nil, fe
				}
			}
			return rw, nil
		}
		if t.FinallyBody != nil {
			_, fe := runBlockInEnv(t.FinallyBody, e)
			if fe != nil {
				return nil, fe
			}
		}
		return cres, nil
	}
	if t.FinallyBody != nil {
		_, fe := runBlockInEnv(t.FinallyBody, e)
		if fe != nil {
			return nil, fe
		}
	}
	return nil, err0
}

func catchClauseApplies(cc *candy_ast.CatchClause, err error) bool {
	if err == nil {
		return false
	}
	n := strings.TrimSpace(candy_ast.ExprAsSimpleTypeName(cc.Type))
	if n == "" || strings.EqualFold(n, "any") {
		return true
	}
	if strings.EqualFold(n, "Error") || strings.EqualFold(n, "error") {
		return true
	}
	// No typed throw values yet: treat a named catch as documentation and still run it.
	return true
}

// runBlockWithValue runs a block in env and returns the value of the last statement
// (expression, val, var) like top-level eval; it propagates ReturnWrap and errors.
func runBlockWithValue(b *candy_ast.BlockStatement, env *Env) (any, error) {
	if b == nil {
		return nil, nil
	}
	var last any
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
		if r != nil {
			last = r
		}
	}
	return last, nil
}
