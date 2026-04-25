package candy_evaluator

import (
	"candy/candy_ast"
)

func evalSwitch(s *candy_ast.SwitchStatement, e *Env) (any, error) {
	subject, err := evalExpression(s.Subject, e)
	if err != nil {
		return nil, err
	}

	var defaultCase *candy_ast.SwitchCase

	for _, c := range s.Cases {
		if c.IsDefault {
			defaultCase = &c
			continue
		}

		for _, pat := range c.Patterns {
			pVal, err := evalExpression(pat, e)
			if err != nil {
				return nil, err
			}

			if valueEqual(subject, pVal) {
				res, err := evalStatement(c.Body, e)
				if _, ok := res.(BreakWrap); ok {
					return nil, err
				}
				return res, err
			}
		}
	}

	if defaultCase != nil {
		res, err := evalStatement(defaultCase.Body, e)
		if _, ok := res.(BreakWrap); ok {
			return nil, err
		}
		return res, err
	}

	return nil, nil
}
