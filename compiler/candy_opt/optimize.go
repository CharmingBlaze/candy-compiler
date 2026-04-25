package candy_opt

import (
	"candy/candy_ast"
	"candy/candy_token"
	"strconv"
)

// OptimizeProgram performs lightweight AST-level optimization passes.
// Current passes:
// - Constant folding for pure literal expressions.
// - Dead branch/loop elimination for constant-boolean conditions.
func OptimizeProgram(program *candy_ast.Program) *candy_ast.Program {
	if program == nil {
		return nil
	}
	ctx := buildInlineContext(program)
	program.Statements = optimizeStatements(program.Statements, ctx)
	return program
}

type optimizeContext struct {
	inlineable map[string]inlineFunction
}

type inlineFunction struct {
	params []string
	body   candy_ast.Expression
}

func buildInlineContext(program *candy_ast.Program) optimizeContext {
	ctx := optimizeContext{inlineable: map[string]inlineFunction{}}
	if program == nil {
		return ctx
	}
	for _, st := range program.Statements {
		fn, ok := st.(*candy_ast.FunctionStatement)
		if !ok || fn == nil || fn.Name == nil || fn.Body == nil {
			continue
		}
		if fn.IsAsync || fn.Suspend || len(fn.TypeParameters) > 0 || fn.ReturnType != nil {
			continue
		}
		if len(fn.Body.Statements) != 1 {
			continue
		}
		ret, ok := fn.Body.Statements[0].(*candy_ast.ReturnStatement)
		if !ok || ret.ReturnValue == nil {
			continue
		}
		paramNames := make([]string, 0, len(fn.Parameters))
		for _, p := range fn.Parameters {
			if p.Name == nil || p.Name.Value == "" || p.Default != nil || p.TypeName != nil {
				paramNames = nil
				break
			}
			paramNames = append(paramNames, p.Name.Value)
		}
		if paramNames == nil {
			continue
		}
		if !isInlineSafeExpr(ret.ReturnValue, paramNames) {
			continue
		}
		// Skip direct recursion.
		if callsNamedFunction(ret.ReturnValue, fn.Name.Value) {
			continue
		}
		ctx.inlineable[fn.Name.Value] = inlineFunction{
			params: paramNames,
			body:   ret.ReturnValue,
		}
	}
	return ctx
}

func optimizeStatements(in []candy_ast.Statement, ctx optimizeContext) []candy_ast.Statement {
	out := make([]candy_ast.Statement, 0, len(in))
	for _, s := range in {
		out = append(out, optimizeStatement(s, ctx)...)
	}
	return out
}

func optimizeStatement(stmt candy_ast.Statement, ctx optimizeContext) []candy_ast.Statement {
	switch s := stmt.(type) {
	case *candy_ast.ExpressionStatement:
		s.Expression = optimizeExpr(s.Expression, ctx)
		return []candy_ast.Statement{s}
	case *candy_ast.ValStatement:
		s.Value = optimizeExpr(s.Value, ctx)
		return []candy_ast.Statement{s}
	case *candy_ast.VarStatement:
		s.Value = optimizeExpr(s.Value, ctx)
		return []candy_ast.Statement{s}
	case *candy_ast.ReturnStatement:
		s.ReturnValue = optimizeExpr(s.ReturnValue, ctx)
		return []candy_ast.Statement{s}
	case *candy_ast.BlockStatement:
		s.Statements = optimizeStatements(s.Statements, ctx)
		return []candy_ast.Statement{s}
	case *candy_ast.FunctionStatement:
		if s.Body != nil {
			s.Body.Statements = optimizeStatements(s.Body.Statements, ctx)
		}
		return []candy_ast.Statement{s}
	case *candy_ast.IfExpression:
		s.Condition = optimizeExpr(s.Condition, ctx)
		if s.Consequence != nil {
			s.Consequence.Statements = optimizeStatements(s.Consequence.Statements, ctx)
		}
		if s.Alternative != nil {
			alt := optimizeStatement(s.Alternative, ctx)
			if len(alt) == 1 {
				s.Alternative = alt[0]
			}
		}
		if b, ok := asBool(s.Condition); ok {
			if b {
				if s.Consequence == nil {
					return nil
				}
				return s.Consequence.Statements
			}
			if s.Alternative == nil {
				return nil
			}
			if blk, ok := s.Alternative.(*candy_ast.BlockStatement); ok {
				return blk.Statements
			}
			return []candy_ast.Statement{s.Alternative}
		}
		return []candy_ast.Statement{s}
	case *candy_ast.WhileStatement:
		s.Condition = optimizeExpr(s.Condition, ctx)
		if s.Body != nil {
			s.Body.Statements = optimizeStatements(s.Body.Statements, ctx)
		}
		if b, ok := asBool(s.Condition); ok && !b {
			return nil
		}
		return []candy_ast.Statement{s}
	case *candy_ast.DoWhileStatement:
		s.Condition = optimizeExpr(s.Condition, ctx)
		if s.Body != nil {
			s.Body.Statements = optimizeStatements(s.Body.Statements, ctx)
		}
		if b, ok := asBool(s.Condition); ok && !b && s.Body != nil {
			return s.Body.Statements
		}
		return []candy_ast.Statement{s}
	case *candy_ast.RepeatStatement:
		s.Count = optimizeExpr(s.Count, ctx)
		if s.Body != nil {
			s.Body.Statements = optimizeStatements(s.Body.Statements, ctx)
		}
		if n, ok := asInt(s.Count); ok && n <= 0 {
			return nil
		}
		return []candy_ast.Statement{s}
	default:
		return []candy_ast.Statement{stmt}
	}
}

func optimizeExpr(expr candy_ast.Expression, ctx optimizeContext) candy_ast.Expression {
	switch e := expr.(type) {
	case *candy_ast.CallExpression:
		e.Function = optimizeExpr(e.Function, ctx)
		for i := range e.Arguments {
			e.Arguments[i] = optimizeExpr(e.Arguments[i], ctx)
		}
		if inlined, ok := inlineCall(e, ctx); ok {
			return optimizeExpr(inlined, ctx)
		}
		return e
	case *candy_ast.InfixExpression:
		e.Left = optimizeExpr(e.Left, ctx)
		e.Right = optimizeExpr(e.Right, ctx)
		if folded, ok := foldInfix(e); ok {
			return folded
		}
		return e
	case *candy_ast.PrefixExpression:
		e.Right = optimizeExpr(e.Right, ctx)
		if folded, ok := foldPrefix(e); ok {
			return folded
		}
		return e
	case *candy_ast.TernaryExpression:
		e.Condition = optimizeExpr(e.Condition, ctx)
		e.Consequence = optimizeExpr(e.Consequence, ctx)
		e.Alternative = optimizeExpr(e.Alternative, ctx)
		if b, ok := asBool(e.Condition); ok {
			if b {
				return e.Consequence
			}
			return e.Alternative
		}
		return e
	case *candy_ast.GroupedExpression:
		e.Expr = optimizeExpr(e.Expr, ctx)
		return e.Expr
	default:
		return expr
	}
}

func inlineCall(call *candy_ast.CallExpression, ctx optimizeContext) (candy_ast.Expression, bool) {
	if call == nil || call.Function == nil {
		return nil, false
	}
	fnIdent, ok := call.Function.(*candy_ast.Identifier)
	if !ok || fnIdent == nil {
		return nil, false
	}
	fn, ok := ctx.inlineable[fnIdent.Value]
	if !ok {
		return nil, false
	}
	if len(call.Arguments) != len(fn.params) {
		return nil, false
	}
	subst := make(map[string]candy_ast.Expression, len(fn.params))
	for i, name := range fn.params {
		arg := call.Arguments[i]
		if !isPureExpr(arg) {
			return nil, false
		}
		subst[name] = cloneExpr(arg)
	}
	return substituteExpr(fn.body, subst), true
}

func isInlineSafeExpr(expr candy_ast.Expression, params []string) bool {
	if expr == nil {
		return false
	}
	allowed := map[string]struct{}{}
	for _, p := range params {
		allowed[p] = struct{}{}
	}
	var walk func(candy_ast.Expression) bool
	walk = func(e candy_ast.Expression) bool {
		switch t := e.(type) {
		case *candy_ast.IntegerLiteral, *candy_ast.FloatLiteral, *candy_ast.Boolean, *candy_ast.StringLiteral:
			return true
		case *candy_ast.Identifier:
			_, ok := allowed[t.Value]
			return ok
		case *candy_ast.GroupedExpression:
			return t.Expr != nil && walk(t.Expr)
		case *candy_ast.PrefixExpression:
			return t.Right != nil && walk(t.Right)
		case *candy_ast.InfixExpression:
			return t.Left != nil && t.Right != nil && walk(t.Left) && walk(t.Right)
		case *candy_ast.TernaryExpression:
			return t.Condition != nil && t.Consequence != nil && t.Alternative != nil &&
				walk(t.Condition) && walk(t.Consequence) && walk(t.Alternative)
		default:
			return false
		}
	}
	return walk(expr)
}

func callsNamedFunction(expr candy_ast.Expression, name string) bool {
	switch t := expr.(type) {
	case *candy_ast.CallExpression:
		if id, ok := t.Function.(*candy_ast.Identifier); ok && id != nil && id.Value == name {
			return true
		}
		for _, a := range t.Arguments {
			if callsNamedFunction(a, name) {
				return true
			}
		}
		return callsNamedFunction(t.Function, name)
	case *candy_ast.GroupedExpression:
		return callsNamedFunction(t.Expr, name)
	case *candy_ast.PrefixExpression:
		return callsNamedFunction(t.Right, name)
	case *candy_ast.InfixExpression:
		return callsNamedFunction(t.Left, name) || callsNamedFunction(t.Right, name)
	case *candy_ast.TernaryExpression:
		return callsNamedFunction(t.Condition, name) ||
			callsNamedFunction(t.Consequence, name) ||
			callsNamedFunction(t.Alternative, name)
	default:
		return false
	}
}

func isPureExpr(expr candy_ast.Expression) bool {
	switch t := expr.(type) {
	case *candy_ast.IntegerLiteral, *candy_ast.FloatLiteral, *candy_ast.Boolean, *candy_ast.StringLiteral, *candy_ast.Identifier:
		return true
	case *candy_ast.GroupedExpression:
		return t.Expr != nil && isPureExpr(t.Expr)
	case *candy_ast.PrefixExpression:
		return t.Right != nil && isPureExpr(t.Right)
	case *candy_ast.InfixExpression:
		return t.Left != nil && t.Right != nil && isPureExpr(t.Left) && isPureExpr(t.Right)
	case *candy_ast.TernaryExpression:
		return t.Condition != nil && t.Consequence != nil && t.Alternative != nil &&
			isPureExpr(t.Condition) && isPureExpr(t.Consequence) && isPureExpr(t.Alternative)
	default:
		return false
	}
}

func cloneExpr(expr candy_ast.Expression) candy_ast.Expression {
	switch t := expr.(type) {
	case *candy_ast.IntegerLiteral:
		return &candy_ast.IntegerLiteral{Token: t.Token, Value: t.Value}
	case *candy_ast.FloatLiteral:
		return &candy_ast.FloatLiteral{Token: t.Token, Value: t.Value}
	case *candy_ast.Boolean:
		return &candy_ast.Boolean{Token: t.Token, Value: t.Value}
	case *candy_ast.StringLiteral:
		return &candy_ast.StringLiteral{Token: t.Token, Value: t.Value}
	case *candy_ast.Identifier:
		return &candy_ast.Identifier{Token: t.Token, Value: t.Value, IsPointer: t.IsPointer}
	case *candy_ast.GroupedExpression:
		return &candy_ast.GroupedExpression{Token: t.Token, Expr: cloneExpr(t.Expr)}
	case *candy_ast.PrefixExpression:
		return &candy_ast.PrefixExpression{Token: t.Token, Operator: t.Operator, Right: cloneExpr(t.Right)}
	case *candy_ast.InfixExpression:
		return &candy_ast.InfixExpression{Token: t.Token, Operator: t.Operator, Left: cloneExpr(t.Left), Right: cloneExpr(t.Right)}
	case *candy_ast.TernaryExpression:
		return &candy_ast.TernaryExpression{
			Token:       t.Token,
			Condition:   cloneExpr(t.Condition),
			Consequence: cloneExpr(t.Consequence),
			Alternative: cloneExpr(t.Alternative),
		}
	default:
		return expr
	}
}

func substituteExpr(expr candy_ast.Expression, subst map[string]candy_ast.Expression) candy_ast.Expression {
	switch t := expr.(type) {
	case *candy_ast.Identifier:
		if rep, ok := subst[t.Value]; ok {
			return cloneExpr(rep)
		}
		return &candy_ast.Identifier{Token: t.Token, Value: t.Value, IsPointer: t.IsPointer}
	case *candy_ast.GroupedExpression:
		return &candy_ast.GroupedExpression{Token: t.Token, Expr: substituteExpr(t.Expr, subst)}
	case *candy_ast.PrefixExpression:
		return &candy_ast.PrefixExpression{Token: t.Token, Operator: t.Operator, Right: substituteExpr(t.Right, subst)}
	case *candy_ast.InfixExpression:
		return &candy_ast.InfixExpression{
			Token:    t.Token,
			Operator: t.Operator,
			Left:     substituteExpr(t.Left, subst),
			Right:    substituteExpr(t.Right, subst),
		}
	case *candy_ast.TernaryExpression:
		return &candy_ast.TernaryExpression{
			Token:       t.Token,
			Condition:   substituteExpr(t.Condition, subst),
			Consequence: substituteExpr(t.Consequence, subst),
			Alternative: substituteExpr(t.Alternative, subst),
		}
	default:
		return cloneExpr(expr)
	}
}

func foldPrefix(p *candy_ast.PrefixExpression) (candy_ast.Expression, bool) {
	switch p.Operator {
	case "!":
		if b, ok := asBool(p.Right); ok {
			return &candy_ast.Boolean{Token: boolToken(!b), Value: !b}, true
		}
	case "-":
		if i, ok := asInt(p.Right); ok {
			return &candy_ast.IntegerLiteral{Token: intToken(-i), Value: -i}, true
		}
		if f, ok := asFloat(p.Right); ok {
			return &candy_ast.FloatLiteral{Token: floatToken(-f), Value: -f}, true
		}
	case "~":
		if i, ok := asInt(p.Right); ok {
			return &candy_ast.IntegerLiteral{Token: intToken(^i), Value: ^i}, true
		}
	}
	return nil, false
}

func foldInfix(in *candy_ast.InfixExpression) (candy_ast.Expression, bool) {
	if li, lok := asInt(in.Left); lok {
		if ri, rok := asInt(in.Right); rok {
			switch in.Operator {
			case "+":
				return &candy_ast.IntegerLiteral{Token: intToken(li + ri), Value: li + ri}, true
			case "-":
				return &candy_ast.IntegerLiteral{Token: intToken(li - ri), Value: li - ri}, true
			case "*":
				return &candy_ast.IntegerLiteral{Token: intToken(li * ri), Value: li * ri}, true
			case "/":
				if ri != 0 {
					return &candy_ast.IntegerLiteral{Token: intToken(li / ri), Value: li / ri}, true
				}
			case "%":
				if ri != 0 {
					return &candy_ast.IntegerLiteral{Token: intToken(li % ri), Value: li % ri}, true
				}
			case "|":
				return &candy_ast.IntegerLiteral{Token: intToken(li | ri), Value: li | ri}, true
			case "&":
				return &candy_ast.IntegerLiteral{Token: intToken(li & ri), Value: li & ri}, true
			case "^":
				return &candy_ast.IntegerLiteral{Token: intToken(li ^ ri), Value: li ^ ri}, true
			case "<<":
				return &candy_ast.IntegerLiteral{Token: intToken(li << ri), Value: li << ri}, true
			case ">>":
				return &candy_ast.IntegerLiteral{Token: intToken(li >> ri), Value: li >> ri}, true
			case "==":
				return &candy_ast.Boolean{Token: boolToken(li == ri), Value: li == ri}, true
			case "!=":
				return &candy_ast.Boolean{Token: boolToken(li != ri), Value: li != ri}, true
			case "<":
				return &candy_ast.Boolean{Token: boolToken(li < ri), Value: li < ri}, true
			case "<=":
				return &candy_ast.Boolean{Token: boolToken(li <= ri), Value: li <= ri}, true
			case ">":
				return &candy_ast.Boolean{Token: boolToken(li > ri), Value: li > ri}, true
			case ">=":
				return &candy_ast.Boolean{Token: boolToken(li >= ri), Value: li >= ri}, true
			}
		}
	}
	if lf, lok := asFloat(in.Left); lok {
		if rf, rok := asFloat(in.Right); rok {
			switch in.Operator {
			case "+":
				return &candy_ast.FloatLiteral{Token: floatToken(lf + rf), Value: lf + rf}, true
			case "-":
				return &candy_ast.FloatLiteral{Token: floatToken(lf - rf), Value: lf - rf}, true
			case "*":
				return &candy_ast.FloatLiteral{Token: floatToken(lf * rf), Value: lf * rf}, true
			case "/":
				if rf != 0 {
					return &candy_ast.FloatLiteral{Token: floatToken(lf / rf), Value: lf / rf}, true
				}
			case "==":
				return &candy_ast.Boolean{Token: boolToken(lf == rf), Value: lf == rf}, true
			case "!=":
				return &candy_ast.Boolean{Token: boolToken(lf != rf), Value: lf != rf}, true
			case "<":
				return &candy_ast.Boolean{Token: boolToken(lf < rf), Value: lf < rf}, true
			case "<=":
				return &candy_ast.Boolean{Token: boolToken(lf <= rf), Value: lf <= rf}, true
			case ">":
				return &candy_ast.Boolean{Token: boolToken(lf > rf), Value: lf > rf}, true
			case ">=":
				return &candy_ast.Boolean{Token: boolToken(lf >= rf), Value: lf >= rf}, true
			}
		}
	}
	if lb, lok := asBool(in.Left); lok {
		if rb, rok := asBool(in.Right); rok {
			switch in.Operator {
			case "&&":
				return &candy_ast.Boolean{Token: boolToken(lb && rb), Value: lb && rb}, true
			case "||":
				return &candy_ast.Boolean{Token: boolToken(lb || rb), Value: lb || rb}, true
			case "==":
				return &candy_ast.Boolean{Token: boolToken(lb == rb), Value: lb == rb}, true
			case "!=":
				return &candy_ast.Boolean{Token: boolToken(lb != rb), Value: lb != rb}, true
			}
		}
	}
	return nil, false
}

func asBool(expr candy_ast.Expression) (bool, bool) {
	v, ok := expr.(*candy_ast.Boolean)
	if !ok || v == nil {
		return false, false
	}
	return v.Value, true
}

func asInt(expr candy_ast.Expression) (int64, bool) {
	v, ok := expr.(*candy_ast.IntegerLiteral)
	if !ok || v == nil {
		return 0, false
	}
	return v.Value, true
}

func asFloat(expr candy_ast.Expression) (float64, bool) {
	v, ok := expr.(*candy_ast.FloatLiteral)
	if !ok || v == nil {
		return 0, false
	}
	return v.Value, true
}

func boolToken(v bool) candy_token.Token {
	if v {
		return candy_token.Token{Type: candy_token.TRUE, Literal: "true"}
	}
	return candy_token.Token{Type: candy_token.FALSE, Literal: "false"}
}

func intToken(v int64) candy_token.Token {
	return candy_token.Token{Type: candy_token.INT, Literal: strconv.FormatInt(v, 10)}
}

func floatToken(v float64) candy_token.Token {
	return candy_token.Token{Type: candy_token.FLOAT, Literal: strconv.FormatFloat(v, 'f', -1, 64)}
}
