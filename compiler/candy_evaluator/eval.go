package candy_evaluator

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_stdlib"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Eval runs a program in a fresh environment.
func Eval(n *candy_ast.Program, pre *Env) (result *Value, rerr error) {
	if pre == nil {
		pre = &Env{Store: make(map[string]*Value)}
	}
	if pre.Cwd == "" {
		if wd, err := os.Getwd(); err == nil {
			pre.Cwd = wd
		}
	}
	if pre.Imported == nil {
		pre.Imported = make(map[string]bool)
	}
	registerPrelude(pre)
	var out *Value
	for _, s := range n.Statements {
		if s == nil {
			continue
		}
		r, err := evalStatement(s, pre)
		if err != nil {
			return out, err
		}
		if rw, ok := r.(ReturnWrap); ok {
			return rw.V, nil
		}
		if v, ok := r.(*Value); ok {
			out = v
		}
	}
	runDefers(pre)
	return out, rerr
}

func runDefers(e *Env) {
	for i := len(e.Defers) - 1; i >= 0; i-- {
		evalExpression(e.Defers[i].Call, e)
	}
}

// RuntimeError is a runtime error with position info.
type RuntimeError struct {
	Msg  string
	Line int
	Col  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("Runtime Error: %s at line %d, col %d", e.Msg, e.Line, e.Col)
	}
	return "Runtime Error: " + e.Msg
}

func newError(msg string, n candy_ast.Node) error {
	tok := candy_ast.GetToken(n)
	return &RuntimeError{Msg: msg, Line: tok.Line, Col: tok.Col}
}

func evalStatement(s candy_ast.Statement, e *Env) (any, error) {
	switch t := s.(type) {
	case *candy_ast.ExpressionStatement:
		if t.Expression == nil {
			return nil, nil
		}
		return evalExpression(t.Expression, e)
	case *candy_ast.ValStatement:
		v, err := evalExpression(t.Value, e)
		if err != nil {
			return nil, err
		}
		e.Set(t.Name.Value, v)
		return v, nil
	case *candy_ast.VarStatement:
		var v *Value
		var err error
		if t.Value != nil {
			v, err = evalExpression(t.Value, e)
			if err != nil {
				return nil, err
			}
		} else {
			v = &Value{Kind: ValNull}
		}
		e.Set(t.Name.Value, v)
		return v, nil
	case *candy_ast.DeferStatement:
		e.Defers = append(e.Defers, t)
		return nil, nil
	case *candy_ast.ReturnStatement:
		v := &Value{Kind: ValNull}
		var err error
		if t.ReturnValue != nil {
			v, err = evalExpression(t.ReturnValue, e)
			if err != nil {
				return nil, err
			}
		}
		return ReturnWrap{V: v}, nil
	case *candy_ast.BlockStatement:
		ne := e.NewEnclosed()
		var last any
		for _, s2 := range t.Statements {
			if s2 == nil {
				continue
			}
			r, err := evalStatement(s2, ne)
			if err != nil {
				return nil, err
			}
			if rw, ok := r.(ReturnWrap); ok {
				runDefers(ne)
				return rw, nil
			}
			if b, ok := r.(BreakWrap); ok {
				return b, nil
			}
			if c, ok := r.(ContinueWrap); ok {
				return c, nil
			}
			last = r
		}
		runDefers(ne)
		return last, nil
	case *candy_ast.FunctionStatement:
		built := &Value{Kind: ValFunction, Fn: &functionVal{Stmt: t, Env: e.NewEnclosed(), Outer: e}}
		built.Fn.Env.Set(t.Name.Value, built)
		built.Fn.Outer = e
		e.Set(t.Name.Value, built)
		return built, nil
	case *candy_ast.IfExpression:
		cond, err := evalExpression(t.Condition, e)
		if err != nil {
			return nil, err
		}
		ne := e.NewEnclosed()
		if cond.Truthy() {
			if t.Consequence != nil {
				return evalStatement(t.Consequence, ne)
			}
			return nil, nil
		}
		if t.Alternative != nil {
			return evalStatement(t.Alternative, ne)
		}
		return nil, nil
	case *candy_ast.ImportStatement:
		if t.From != "" {
			if err := evalImport(t.From, e); err != nil {
				return nil, err
			}
			for _, sym := range t.Symbols {
				if v, ok := e.Get(sym); ok {
					e.Set(sym, v)
					continue
				}
				// Resolve from module value if exported as module object.
				if modv, ok := e.Get(t.From); ok && modv != nil && modv.Kind == ValModule && modv.Mod != nil {
					if c, ok2 := modv.Mod.Consts[sym]; ok2 && c != nil {
						e.Set(sym, c)
						continue
					}
					if fn, ok2 := lookupModFn(modv.Mod, sym); ok2 {
						e.Set(sym, &Value{
							Kind: ValFunction,
							Builtin: func(args []*Value) (*Value, error) {
								return fn(args)
							},
						})
						continue
					}
				}
				return nil, &RuntimeError{Msg: "from-import symbol not found: " + sym}
			}
			return nil, nil
		}
		if err := evalImport(t.Path, e); err != nil {
			return nil, err
		}
		if t.Alias != "" {
			if v, ok := e.Get(t.Path); ok && v != nil {
				e.Set(t.Alias, v)
			} else {
				// For common stdlib `import math as m` style, alias existing module name.
				last := t.Path
				if strings.Contains(last, ".") {
					parts := strings.Split(last, ".")
					last = parts[len(parts)-1]
				}
				if v2, ok2 := e.Get(last); ok2 && v2 != nil {
					e.Set(t.Alias, v2)
				}
			}
		}
		return nil, nil
	case *candy_ast.StructStatement:
		e.Set(t.Name.Value, &Value{Kind: ValStruct, St: &structVal{Def: t, Env: e, Data: make(map[string]Value)}})
		return nil, nil
	case *candy_ast.EnumStatement:
		return nil, evalEnumStatement(t, e)
	case *candy_ast.ForStatement:
		if t.Iterable != nil {
			return evalForIn(t, e)
		}
		return evalForTo(t, e)
	case *candy_ast.CForStatement:
		return evalCFor(t, e)
	case *candy_ast.WhileStatement:
		return evalWhile(t, e)
	case *candy_ast.SwitchStatement:
		return evalSwitch(t, e)
	case *candy_ast.TryStatement:
		return evalTry(t, e)
	case *candy_ast.ClassStatement:
		return evalClassStatement(t, e)
	case *candy_ast.ObjectStatement:
		return evalObjectStatement(t, e)
	case *candy_ast.ForEachStatement:
		return evalForEach(t, e)
	case *candy_ast.RepeatStatement:
		return evalRepeat(t, e)
	case *candy_ast.LoopStatement:
		return evalLoop(t, e)
	case *candy_ast.BreakStatement:
		return BreakWrap{}, nil
	case *candy_ast.ContinueStatement:
		return ContinueWrap{}, nil
	case *candy_ast.DeleteStatement:
		return evalDelete(t, e)
	case *candy_ast.WithStatement:
		v, err := evalExpression(t.Value, e)
		if err != nil {
			return nil, err
		}
		ne := e.NewEnclosed()
		ne.Set(t.Name.Value, v)
		r, err := evalStatement(t.Body, ne)
		if cleanupErr := cleanupWithBinding(t.Name.Value, ne); cleanupErr != nil && err == nil {
			err = cleanupErr
		}
		if err != nil {
			return nil, err
		}
		return r, nil
	case *candy_ast.LibraryStatement:
		if t.Body == nil {
			return nil, nil
		}
		var last any
		for _, st := range t.Body.Statements {
			r, err := evalStatement(st, e)
			if err != nil {
				return nil, err
			}
			last = r
		}
		return last, nil
	default:
		return nil, fmt.Errorf("unhandled statement %T", s)
	}
}

func cleanupWithBinding(name string, env *Env) error {
	v, ok := env.Get(name)
	if !ok || v == nil || v.Kind == ValNull {
		return nil
	}
	for _, fnName := range []string{"close", "unload", "release", "dispose"} {
		if fn, ok2 := Builtins[fnName]; ok2 {
			if _, err := fn([]*Value{v}); err == nil {
				env.Set(name, &Value{Kind: ValNull})
				return nil
			}
		}
	}
	env.Set(name, &Value{Kind: ValNull})
	return nil
}

func evalExpression(ex candy_ast.Expression, e *Env) (*Value, error) {
	if ex == nil {
		return &Value{Kind: ValNull}, nil
	}
	switch t := ex.(type) {
	case *candy_ast.NullLiteral:
		return &Value{Kind: ValNull}, nil
	case *candy_ast.StringLiteral:
		return &Value{Kind: ValString, Str: t.Value}, nil
	case *candy_ast.IntegerLiteral:
		return &Value{Kind: ValInt, I64: t.Value}, nil
	case *candy_ast.FloatLiteral:
		return &Value{Kind: ValFloat, F64: t.Value}, nil
	case *candy_ast.Boolean:
		return &Value{Kind: ValBool, B: t.Value}, nil
	case *candy_ast.Identifier:
		if v, ok := e.Get(t.Value); ok {
			return v, nil
		}
		return nil, withUndefinedVar(t.Value, e.AllNameBindings())
	case *candy_ast.PrefixExpression:
		switch t.Operator {
		case "++", "--":
			ident, ok := t.Right.(*candy_ast.Identifier)
			if !ok {
				return nil, newError(t.Operator+" requires variable", t)
			}
			val, ok := e.Get(ident.Value)
			if !ok {
				return nil, withUndefinedVar(ident.Value, e.AllNameBindings())
			}
			diff := int64(1)
			if t.Operator == "--" {
				diff = -1
			}
			if val.Kind == ValInt {
				newVal := &Value{Kind: ValInt, I64: val.I64 + diff}
				e.Update(ident.Value, newVal)
				return newVal, nil
			}
			if val.Kind == ValFloat {
				newVal := &Value{Kind: ValFloat, F64: val.F64 + float64(diff)}
				e.Update(ident.Value, newVal)
				return newVal, nil
			}
			return nil, newError(t.Operator+" requires numeric variable", t)
		}

		r, err := evalExpression(t.Right, e)
		if err != nil {
			return nil, err
		}
		switch t.Operator {
		case "-":
			if r.Kind == ValInt {
				return &Value{Kind: ValInt, I64: -r.I64}, nil
			}
			if r.Kind == ValFloat {
				return &Value{Kind: ValFloat, F64: -r.F64}, nil
			}
		case "!", "not":
			return &Value{Kind: ValBool, B: !r.Truthy()}, nil
		case "~":
			if r.Kind != ValInt {
				return nil, newError("`~` expects int operand", t)
			}
			return &Value{Kind: ValInt, I64: ^r.I64}, nil
		}
		return nil, newError("bad prefix "+t.Operator, t)
	case *candy_ast.PostfixExpression:
		if t.Operator != "++" && t.Operator != "--" {
			return nil, newError("bad postfix "+t.Operator, t)
		}
		ident, ok := t.Left.(*candy_ast.Identifier)
		if !ok {
			return nil, newError(t.Operator+" requires variable", t)
		}
		val, ok2 := e.Get(ident.Value)
		if !ok2 {
			return nil, withUndefinedVar(ident.Value, e.AllNameBindings())
		}
		diff := int64(1)
		if t.Operator == "--" {
			diff = -1
		}
		if val.Kind == ValInt {
			newVal := &Value{Kind: ValInt, I64: val.I64 + diff}
			e.Update(ident.Value, newVal)
			return newVal, nil
		}
		if val.Kind == ValFloat {
			newVal := &Value{Kind: ValFloat, F64: val.F64 + float64(diff)}
			e.Update(ident.Value, newVal)
			return newVal, nil
		}
		return nil, newError(t.Operator+" requires numeric variable", t)
	case *candy_ast.AssignExpression:
		return evalAssign(t, e)
	case *candy_ast.TernaryExpression:
		cond, err := evalExpression(t.Condition, e)
		if err != nil {
			return nil, err
		}
		if cond.Truthy() {
			return evalExpression(t.Consequence, e)
		}
		return evalExpression(t.Alternative, e)
	case *candy_ast.StructLiteral:
		return evalStructLiteral(t, e)
	case *candy_ast.DotExpression:
		return evalDot(t, e)
	case *candy_ast.InterpolatedStringLiteral:
		var b strings.Builder
		for _, part := range t.Parts {
			v, err := evalExpression(part, e)
			if err != nil {
				return nil, err
			}
			if v == nil {
				continue
			}
			b.WriteString(v.String())
		}
		return &Value{Kind: ValString, Str: b.String()}, nil
	case *candy_ast.InfixExpression:
		if t.Operator == "??" {
			l, err := evalExpression(t.Left, e)
			if err != nil {
				return nil, err
			}
			if isNullishValue(l) {
				return evalExpression(t.Right, e)
			}
			return l, nil
		}
		if strings.EqualFold(t.Operator, "and") || t.Operator == "&&" {
			l, err := evalExpression(t.Left, e)
			if err != nil {
				return nil, err
			}
			if !l.Truthy() {
				return l, nil
			}
			return evalExpression(t.Right, e)
		}
		if strings.EqualFold(t.Operator, "or") {
			l, err := evalExpression(t.Left, e)
			if err != nil {
				return nil, err
			}
			// `or` behaves as default-value operator for nullish values.
			// For booleans, keep logical-or short-circuit behavior.
			if l != nil && l.Kind == ValBool {
				if l.B {
					return l, nil
				}
				return evalExpression(t.Right, e)
			}
			if isNullishValue(l) {
				return evalExpression(t.Right, e)
			}
			return l, nil
		}
		if t.Operator == "||" {
			l, err := evalExpression(t.Left, e)
			if err != nil {
				return nil, err
			}
			if l.Truthy() {
				return l, nil
			}
			return evalExpression(t.Right, e)
		}
		l, err := evalExpression(t.Left, e)
		if err != nil {
			return nil, err
		}
		r, err := evalExpression(t.Right, e)
		if err != nil {
			return nil, err
		}
		return evalInfix(t.Operator, l, r, t)
	case *candy_ast.GroupedExpression:
		return evalExpression(t.Expr, e)
	case *candy_ast.RangeExpression:
		l, err := evalExpression(t.Left, e)
		if err != nil {
			return nil, err
		}
		r, err := evalExpression(t.Right, e)
		if err != nil {
			return nil, err
		}
		lo, ok1 := asInt64(l)
		hi, ok2 := asInt64(r)
		if !ok1 || !ok2 {
			return nil, newError("range a..b needs integer endpoints", t)
		}
		exclusive := t.TokenLiteral() == "..<"
		// Inclusive by default; exclusive for `..<`.
		elems := make([]Value, 0, 8)
		if lo <= hi {
			end := hi
			if exclusive {
				end = hi - 1
			}
			for i := lo; i <= end; i++ {
				elems = append(elems, Value{Kind: ValInt, I64: i})
			}
		} else {
			end := hi
			if exclusive {
				end = hi + 1
			}
			for i := lo; i >= end; i-- {
				elems = append(elems, Value{Kind: ValInt, I64: i})
			}
		}
		return &Value{Kind: ValArray, Elems: elems}, nil
	case *candy_ast.ArrayLiteral:
		out := make([]Value, 0, len(t.Elem))
		for _, ex := range t.Elem {
			v, err := evalExpression(ex, e)
			if err != nil {
				return nil, err
			}
			if v == nil {
				v = &Value{Kind: ValNull}
			}
			out = append(out, *v)
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case *candy_ast.MapLiteral:
		sm := make(map[string]Value)
		for _, pr := range t.Pairs {
			kk, err := evalExpression(pr.Key, e)
			if err != nil {
				return nil, err
			}
			vv, err := evalExpression(pr.Value, e)
			if err != nil {
				return nil, err
			}
			ks, kerr := mapKeyString(kk)
			if kerr != nil {
				return nil, kerr
			}
			sm[ks] = *vv
		}
		return &Value{Kind: ValMap, StrMap: sm}, nil
	case *candy_ast.WhenExpression:
		var subject *Value
		var err error
		if t.Subject != nil {
			subject, err = evalExpression(t.Subject, e)
			if err != nil {
				return nil, err
			}
		}
		for _, a := range t.Arms {
			c, err := evalExpression(a.Cond, e)
			if err != nil {
				return nil, err
			}
			if t.Subject != nil {
				if valueEqual(subject, c) {
					return evalExpression(a.Body, e)
				}
				continue
			}
			if c.Truthy() {
				return evalExpression(a.Body, e)
			}
		}
		if t.ElseV != nil {
			return evalExpression(t.ElseV, e)
		}
		return &Value{Kind: ValNull}, nil
	case *candy_ast.MatchExpression:
		sub, err := evalExpression(t.Subject, e)
		if err != nil {
			return nil, err
		}
		for _, b := range t.Branches {
			bindings := make(map[string]*Value)
			if matchPattern(b.Pat, sub, bindings, e) {
				ne := e.NewEnclosed()
				for k, v := range bindings {
					ne.Set(k, v)
				}
				if b.Guard != nil {
					gv, err3 := evalExpression(b.Guard, ne)
					if err3 != nil {
						return nil, err3
					}
					if !gv.Truthy() {
						continue
					}
				}
				return evalExpression(b.Body, ne)
			}
		}
		if t.Default != nil {
			return evalExpression(t.Default, e)
		}
		return &Value{Kind: ValNull}, nil
	case *candy_ast.IfExpression:
		res, err := evalStatement(t, e)
		if err != nil {
			return nil, err
		}
		if v, ok := res.(*Value); ok {
			return v, nil
		}
		return &Value{Kind: ValNull}, nil
	case *candy_ast.IndexExpression:
		le, err := evalExpression(t.Base, e)
		if err != nil {
			return nil, err
		}
		if (le == nil || le.Kind == ValNull) && t.IsSafe {
			return &Value{Kind: ValNull}, nil
		}
		if le == nil {
			return nil, newError("bad index", t)
		}
		if re, ok := t.Index.(*candy_ast.RangeExpression); ok {
			if le.Kind != ValArray && le.Kind != ValString {
				return nil, newError("range index/slice requires array or string", t)
			}
			startV, err := evalExpression(re.Left, e)
			if err != nil {
				return nil, err
			}
			endV, err := evalExpression(re.Right, e)
			if err != nil {
				return nil, err
			}
			startI, ok1 := asInt64(startV)
			endI, ok2 := asInt64(endV)
			if !ok1 || !ok2 {
				return nil, newError("slice bounds must be integers", t)
			}
			exclusive := re.TokenLiteral() == "..<"
			if le.Kind == ValArray {
				n := len(le.Elems)
				st := int(startI)
				en := int(endI)
				if st < 0 {
					st = n + st
				}
				if en < 0 {
					en = n + en
				}
				if !exclusive {
					en++
				}
				if st < 0 {
					st = 0
				}
				if st > n {
					st = n
				}
				if en < st {
					en = st
				}
				if en > n {
					en = n
				}
				out := make([]Value, 0, en-st)
				out = append(out, le.Elems[st:en]...)
				return &Value{Kind: ValArray, Elems: out}, nil
			}
			runes := []rune(le.Str)
			n := len(runes)
			st := int(startI)
			en := int(endI)
			if st < 0 {
				st = n + st
			}
			if en < 0 {
				en = n + en
			}
			if !exclusive {
				en++
			}
			if st < 0 {
				st = 0
			}
			if st > n {
				st = n
			}
			if en < st {
				en = st
			}
			if en > n {
				en = n
			}
			return &Value{Kind: ValString, Str: string(runes[st:en])}, nil
		}
		ix, err := evalExpression(t.Index, e)
		if err != nil {
			return nil, err
		}
		if le.Kind == ValArray && ix.Kind == ValInt {
			i := int(ix.I64)
			if i < 0 {
				i = len(le.Elems) + i
			}
			if i < 0 || i >= len(le.Elems) {
				return nil, newError("index out of range", t)
			}
			vv := le.Elems[i]
			return &vv, nil
		}
		if le.Kind == ValString && ix.Kind == ValInt {
			s := le.Str
			sr := []rune(s)
			i := int(ix.I64)
			if i < 0 {
				i = len(sr) + i
			}
			if i < 0 || i >= len(sr) {
				return nil, newError("string index out of range", t)
			}
			return &Value{Kind: ValString, Str: string(sr[i])}, nil
		}
		if le.Kind == ValMap && le.StrMap != nil {
			ks, kerr := mapKeyString(ix)
			if kerr != nil {
				return nil, kerr
			}
			if v, ok := le.StrMap[ks]; ok {
				v2 := v
				return &v2, nil
			}
			return &Value{Kind: ValNull}, nil
		}
		return nil, newError("bad index", t)
	case *candy_ast.LambdaExpression:
		lamName := &candy_ast.Identifier{Token: t.Token, Value: "<lambda>"}
		retSt := &candy_ast.ReturnStatement{Token: t.Token, ReturnValue: t.Body}
		body := &candy_ast.BlockStatement{Token: t.Token, Statements: []candy_ast.Statement{retSt}}
		fs := &candy_ast.FunctionStatement{Token: t.Token, Name: lamName, Parameters: t.Parameters, Body: body}
		return &Value{Kind: ValFunction, Fn: &functionVal{Stmt: fs, Env: e.NewEnclosed(), Outer: e}}, nil
	case *candy_ast.TupleLiteral:
		var elems []Value
		for _, el := range t.Elems {
			ev, err := evalExpression(el, e)
			if err != nil {
				return nil, err
			}
			elems = append(elems, *ev)
		}
		return &Value{Kind: ValArray, Elems: elems}, nil
	case *candy_ast.CallExpression:
		if dot, ok := t.Function.(*candy_ast.DotExpression); ok {
			return evalMethodCall(dot, t.Arguments, e)
		}
		args, namedArgs, err := evalCallArgs(t.Arguments, e)
		if err != nil {
			return nil, err
		}
		if id, ok := t.Function.(*candy_ast.Identifier); ok {
			// Prefer lexical bindings only when they are callable targets.
			// If a non-callable binding shadows a builtin name (e.g. module `random`),
			// keep builtin-call compatibility for `random(...)`.
			callableBinding := false
			if bv, bound := e.Get(id.Value); bound && bv != nil {
				if (bv.Kind == ValFunction && (bv.Fn != nil || bv.Builtin != nil)) ||
					(bv.Kind == ValStruct && bv.St != nil && (bv.St.Def != nil || bv.St.ClassDef != nil)) {
					callableBinding = true
				}
			}
			if !callableBinding {
				if f, ok2 := Builtins[strings.ToLower(id.Value)]; ok2 {
					if len(namedArgs) > 0 {
						return nil, newError("named arguments are only supported for user functions", t)
					}
					return f(args)
				}
				// Inside object/class methods, allow bare method calls to resolve on `this`.
				// Example: `add(entity)` resolves as `this.add(entity)`.
				// Must not run when id names a class/function in the environment (e.g. Scene() vs this.scene()).
				if _, okThis := e.Get("this"); okThis {
					dot := &candy_ast.DotExpression{
						Token: t.Token,
						Left:  &candy_ast.Identifier{Token: t.Token, Value: "this"},
						Right: &candy_ast.Identifier{Token: t.Token, Value: id.Value},
					}
					return evalMethodCall(dot, t.Arguments, e)
				}
			}
		}
		fnv, err := evalExpression(t.Function, e)
		if err != nil {
			return nil, err
		}
		if t.IsSafe && (fnv == nil || fnv.Kind == ValNull) {
			return &Value{Kind: ValNull}, nil
		}
		if fnv == nil {
			return nil, newError("nil call", t)
		}
		if fnv.Kind == ValStruct && (fnv.St.Def != nil || fnv.St.ClassDef != nil) {
			return instantiateClass(fnv, args, e)
		}
		if len(namedArgs) > 0 {
			if fnv.Kind != ValFunction || fnv.Fn == nil || fnv.Fn.Stmt == nil {
				return nil, newError("named arguments require a user function call target", t)
			}
			args = reorderCallArgs(fnv.Fn.Stmt, args, namedArgs)
		}
		return evalUserFunction(fnv, args)
	default:
		return nil, newError(fmt.Sprintf("unhandled expression %T", t), ex)
	}
}

func instantiateClass(def *Value, args []*Value, e *Env) (*Value, error) {
	instance := &Value{Kind: ValStruct, St: &structVal{Data: make(map[string]Value), Env: def.St.Env}}
	if def.St.Def != nil {
		instance.St.Def = def.St.Def
		// existing struct logic...
		for i, p := range def.St.Def.Fields {
			if i < len(args) {
				instance.St.Data[p.Name.Value] = *args[i]
			} else if p.Init != nil {
				v, _ := evalExpression(p.Init, def.St.Env)
				instance.St.Data[p.Name.Value] = *v
			}
		}
	} else if def.St.ClassDef != nil {
		instance.St.ClassDef = def.St.ClassDef
		// Materialize inherited fields/defaults first so base class state exists.
		populateClassInstanceData(instance, def.St.ClassDef, def.St.Env, e)
		// Primary constructor parameters -> Fields
		for i, p := range def.St.ClassDef.Parameters {
			if i < len(args) {
				instance.St.Data[p.Name.Value] = *args[i]
			} else if p.Default != nil {
				v, err := evalExpression(p.Default, def.St.Env)
				if err != nil {
					return nil, err
				}
				if v != nil {
					instance.St.Data[p.Name.Value] = *v
				}
			}
		}
		// Evaluate class body members (fields/methods)
		for _, m := range def.St.ClassDef.Members {
			switch st := m.(type) {
			case *candy_ast.VarStatement:
				if st.Name == nil || st.Value == nil {
					continue
				}
				v, err := evalExpression(st.Value, def.St.Env)
				if err != nil {
					return nil, err
				}
				if v != nil {
					instance.St.Data[st.Name.Value] = *v
				}
			case *candy_ast.ValStatement:
				if st.Name == nil || st.Value == nil {
					continue
				}
				v, err := evalExpression(st.Value, def.St.Env)
				if err != nil {
					return nil, err
				}
				if v != nil {
					instance.St.Data[st.Name.Value] = *v
				}
			case *candy_ast.ExpressionStatement:
				// Support object/class field initializers written as `x = 0`.
				if asg, ok := st.Expression.(*candy_ast.AssignExpression); ok {
					if id, ok2 := asg.Left.(*candy_ast.Identifier); ok2 {
						v, err := evalExpression(asg.Value, def.St.Env)
						if err != nil {
							return nil, err
						}
						if v != nil {
							instance.St.Data[id.Value] = *v
						}
					}
				}
			}
			// Methods are handled by lookup in ClassDef during evalMethodCall
		}
	}
	// Constructor convention: if class defines `init(...)`, call it after
	// instance data initialization. This enables patterns like Model("path").
	if instance.St != nil && instance.St.ClassDef != nil {
		if err := invokeClassInit(instance, args, e); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func populateClassInstanceData(instance *Value, cls *candy_ast.ClassStatement, classEnv *Env, callEnv *Env) {
	if instance == nil || instance.St == nil || cls == nil {
		return
	}
	// Base first, then derived overrides.
	if cls.Base != nil {
		baseName := cls.Base.Value
		var baseVal *Value
		if classEnv != nil {
			if v, ok := classEnv.Get(baseName); ok {
				baseVal = v
			}
		}
		if baseVal == nil && callEnv != nil {
			if v, ok := callEnv.Get(baseName); ok {
				baseVal = v
			}
		}
		if baseVal != nil && baseVal.Kind == ValStruct && baseVal.St != nil && baseVal.St.ClassDef != nil {
			populateClassInstanceData(instance, baseVal.St.ClassDef, baseVal.St.Env, callEnv)
		}
	}
	for _, m := range cls.Members {
		switch st := m.(type) {
		case *candy_ast.VarStatement:
			if st.Name == nil || st.Value == nil {
				continue
			}
			v, _ := evalExpression(st.Value, classEnv)
			if v != nil {
				instance.St.Data[st.Name.Value] = *v
			}
		case *candy_ast.ValStatement:
			if st.Name == nil || st.Value == nil {
				continue
			}
			v, _ := evalExpression(st.Value, classEnv)
			if v != nil {
				instance.St.Data[st.Name.Value] = *v
			}
		case *candy_ast.ExpressionStatement:
			if asg, ok := st.Expression.(*candy_ast.AssignExpression); ok {
				if id, ok2 := asg.Left.(*candy_ast.Identifier); ok2 {
					v, _ := evalExpression(asg.Value, classEnv)
					if v != nil {
						instance.St.Data[id.Value] = *v
					}
				}
			}
		}
	}
}

func invokeClassInit(instance *Value, args []*Value, outer *Env) error {
	if instance == nil || instance.St == nil || instance.St.ClassDef == nil {
		return nil
	}
	var initFn *candy_ast.FunctionStatement
	curr := instance.St.ClassDef
	for curr != nil && initFn == nil {
		for _, m := range curr.Members {
			if fn, ok := m.(*candy_ast.FunctionStatement); ok && fn.Name != nil && strings.EqualFold(fn.Name.Value, "init") {
				initFn = fn
				break
			}
		}
		if initFn != nil || curr.Base == nil {
			break
		}
		baseName := curr.Base.Value
		if baseName == "" {
			break
		}
		if outer != nil {
			if baseVal, ok := outer.Get(baseName); ok && baseVal != nil && baseVal.Kind == ValStruct && baseVal.St != nil {
				curr = baseVal.St.ClassDef
				continue
			}
		}
		if instance.St.Env != nil {
			if baseVal, ok := instance.St.Env.Get(baseName); ok && baseVal != nil && baseVal.Kind == ValStruct && baseVal.St != nil {
				curr = baseVal.St.ClassDef
				continue
			}
		}
		break
	}
	if initFn == nil || initFn.Body == nil {
		return nil
	}
	baseEnv := instance.St.Env
	if baseEnv == nil {
		baseEnv = outer
	}
	if baseEnv == nil {
		baseEnv = &Env{Store: make(map[string]*Value), Imported: map[string]bool{}}
	}
	ne := baseEnv.NewEnclosed()
	for k, v := range instance.St.Data {
		ptr := new(Value)
		*ptr = v
		ne.Set(k, ptr)
	}
	recvName := "this"
	if initFn.Receiver != nil && initFn.Receiver.Name != nil && initFn.Receiver.Name.Value != "" {
		recvName = initFn.Receiver.Name.Value
	}
	ne.Set(recvName, instance)
	for i, p := range initFn.Parameters {
		if i < len(args) {
			ne.Set(p.Name.Value, args[i])
		} else if p.Default != nil {
			dv, err := evalExpression(p.Default, baseEnv)
			if err != nil {
				return err
			}
			ne.Set(p.Name.Value, dv)
		}
	}
	for _, st := range initFn.Body.Statements {
		if _, err := evalStatement(st, ne); err != nil {
			return err
		}
	}
	for k := range instance.St.Data {
		if v, ok := ne.Get(k); ok {
			instance.St.Data[k] = *v
		}
	}
	return nil
}

// evalUserFunction invokes a user-defined function value (closure).
func evalUserFunction(fnv *Value, args []*Value) (*Value, error) {
	if fnv == nil {
		return nil, &RuntimeError{Msg: "nil call"}
	}
	if fnv.Kind == ValFunction && fnv.Builtin != nil {
		return fnv.Builtin(args)
	}
	if fnv.Kind != ValFunction || fnv.Fn == nil {
		return nil, &RuntimeError{Msg: "not a function"}
	}
	outer := fnv.Fn.Outer
	if outer == nil {
		return nil, &RuntimeError{Msg: "function has no environment"}
	}
	ne := outer.NewEnclosed()
	for i, p0 := range fnv.Fn.Stmt.Parameters {
		if p0.Name == nil {
			continue
		}
		if i < len(args) {
			ne.Set(p0.Name.Value, args[i])
		} else if p0.Default != nil {
			dv, err3 := evalExpression(p0.Default, outer) // Use outer env for default value evaluation
			if err3 != nil {
				return nil, err3
			}
			ne.Set(p0.Name.Value, dv)
		}
	}
	if fnv.Fn.Stmt.Name != nil {
		ne.Set(fnv.Fn.Stmt.Name.Value, fnv)
	}
	if fnv.Fn.Stmt.Body == nil {
		return &Value{Kind: ValNull}, nil
	}
	for i, st := range fnv.Fn.Stmt.Body.Statements {
		r, err2 := evalStatement(st, ne)
		if err2 != nil {
			return nil, err2
		}
		if rw, ok2 := r.(ReturnWrap); ok2 {
			runDefers(ne)
			return rw.V, nil
		}
		// Implicit return: if this is the last statement and it's an expression result
		if i == len(fnv.Fn.Stmt.Body.Statements)-1 {
			if v, ok := r.(*Value); ok {
				runDefers(ne)
				return v, nil
			}
		}
	}
	runDefers(ne)
	return &Value{Kind: ValNull}, nil
}

func mapKeyString(v *Value) (string, error) {
	if v == nil {
		return "", nil
	}
	if v.Kind == ValString {
		return v.Str, nil
	}
	if v.Kind == ValInt {
		return fmt.Sprintf("%d", v.I64), nil
	}
	if v.Kind == ValBool {
		return fmt.Sprintf("%t", v.B), nil
	}
	return "", &RuntimeError{Msg: "map key must be string/int/bool"}
}

func evalCallArgs(argExprs []candy_ast.Expression, e *Env) ([]*Value, map[string]*Value, error) {
	args := make([]*Value, 0, len(argExprs))
	named := make(map[string]*Value)
	for _, a := range argExprs {
		if na, ok := a.(*candy_ast.NamedArgumentExpression); ok {
			if na == nil || na.Name == nil {
				return nil, nil, &RuntimeError{Msg: "invalid named argument"}
			}
			v, err := evalExpression(na.Value, e)
			if err != nil {
				return nil, nil, err
			}
			named[strings.ToLower(na.Name.Value)] = v
			continue
		}
		av, err := evalExpression(a, e)
		if err != nil {
			return nil, nil, err
		}
		args = append(args, av)
	}
	return args, named, nil
}

func reorderCallArgs(fn *candy_ast.FunctionStatement, positional []*Value, named map[string]*Value) []*Value {
	if fn == nil {
		return positional
	}
	out := make([]*Value, 0, len(fn.Parameters))
	pos := 0
	for _, p := range fn.Parameters {
		if p.Name == nil {
			continue
		}
		if pos < len(positional) {
			out = append(out, positional[pos])
			pos++
			continue
		}
		if v, ok := named[strings.ToLower(p.Name.Value)]; ok {
			out = append(out, v)
			continue
		}
		out = append(out, nil)
	}
	return out
}

func evalImport(path string, env *Env) error {
	if env == nil {
		return &RuntimeError{Msg: "nil env for import"}
	}
	full := path
	if src, ok := candy_stdlib.Lookup(path); ok {
		stdlibKey := "__stdlib__:" + path
		if env.Imported[stdlibKey] {
			return nil
		}
		l := candy_lexer.New(src)
		p := candy_parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) > 0 {
			var msgs []string
			for _, d := range p.Errors() {
				lineText := ""
				lines := strings.Split(src, "\n")
				if d.Line > 0 && d.Line <= len(lines) {
					lineText = " -> " + strings.TrimSpace(lines[d.Line-1])
				}
				msgs = append(msgs, fmt.Sprintf("%d:%d: %s%s", d.Line, d.Col, d.Message, lineText))
			}
			return &RuntimeError{Msg: "stdlib import parse error: " + strings.Join(msgs, "; ")}
		}
		env.Imported[stdlibKey] = true
		_, eerr := Eval(prog, env)
		if eerr != nil {
			delete(env.Imported, stdlibKey)
		}
		return eerr
	}
	if !filepath.IsAbs(full) {
		base := env.Cwd
		if base == "" {
			if wd, err := os.Getwd(); err == nil {
				base = wd
			}
		}
		full = filepath.Join(base, path)
	}
	full = filepath.Clean(full)
	if env.Imported[full] {
		return nil
	}
	b, err := os.ReadFile(full)
	if err != nil {
		return err
	}
	env.Imported[full] = true
	oldCwd := env.Cwd
	env.Cwd = filepath.Dir(full)
	defer func() { env.Cwd = oldCwd }()

	l := candy_lexer.New(string(b))
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		var msgs []string
		for _, d := range p.Errors() {
			msgs = append(msgs, d.Message)
		}
		return &RuntimeError{Msg: "import parse error: " + strings.Join(msgs, "; ")}
	}
	_, eerr := Eval(prog, env)
	return eerr
}

func evalStructOperatorOverload(left *Value, op string, right *Value) (*Value, bool) {
	if left == nil || left.St == nil || left.St.Def == nil {
		return nil, false
	}
	for _, ov := range left.St.Def.Operators {
		if ov == nil || ov.Operator != op || ov.Body == nil {
			continue
		}
		outer := left.St.Env
		if outer == nil {
			outer = &Env{Store: map[string]*Value{}}
		}
		ne := outer.NewEnclosed()
		ne.Set("this", left)
		if len(ov.Parameters) > 0 && ov.Parameters[0].Name != nil {
			ne.Set(ov.Parameters[0].Name.Value, right)
		}
		for _, st := range ov.Body.Statements {
			r, err := evalStatement(st, ne)
			if err != nil {
				return &Value{Kind: ValNull}, true
			}
			if rw, ok := r.(ReturnWrap); ok {
				return rw.V, true
			}
		}
		return &Value{Kind: ValNull}, true
	}
	return nil, false
}

// isNullishValue is true for null coalescing: only null (and nil pointer as null).
func isNullishValue(v *Value) bool {
	if v == nil {
		return true
	}
	return v.Kind == ValNull
}

func evalInfix(op string, l, r *Value, node candy_ast.Node) (*Value, error) {
	if l == nil {
		l = &Value{Kind: ValNull}
	}
	if r == nil {
		r = &Value{Kind: ValNull}
	}
	if l.Kind == ValVec || r.Kind == ValVec {
		return evalVecInfix(op, l, r, node)
	}
	if lv, lok := vecLikeFromValue(l); lok {
		if rv, rok := vecLikeFromValue(r); rok {
			if len(lv) != len(rv) {
				return nil, newError("vector operation expects equal dimensions", node)
			}
			switch op {
			case "+", "-":
				out := make([]float64, len(lv))
				for i := range lv {
					if op == "+" {
						out[i] = lv[i] + rv[i]
					} else {
						out[i] = lv[i] - rv[i]
					}
				}
				return &Value{Kind: ValVec, Vec: out}, nil
			}
		}
	}
	if l.Kind == ValStruct && l.St != nil && l.St.Def != nil {
		if ov, ok := evalStructOperatorOverload(l, op, r); ok {
			return ov, nil
		}
	}
	switch op {
	case "==", "!=":
		eq := valueEqual(l, r)
		if op == "==" {
			return &Value{Kind: ValBool, B: eq}, nil
		}
		return &Value{Kind: ValBool, B: !eq}, nil
	}
	lf, rf, useF := toFloats(l, r)
	switch op {
	case "<", ">", "<=", ">=":
		if !isNumeric(l) || !isNumeric(r) {
			return nil, newError("comparison expects numeric operands", node)
		}
		if useF {
			t := false
			switch op {
			case "<":
				t = lf < rf
			case ">":
				t = lf > rf
			case "<=":
				t = lf <= rf
			case ">=":
				t = lf >= rf
			}
			return &Value{Kind: ValBool, B: t}, nil
		}
		if l.Kind == ValInt && r.Kind == ValInt {
			var t bool
			switch op {
			case "<":
				t = l.I64 < r.I64
			case ">":
				t = l.I64 > r.I64
			case "<=":
				t = l.I64 <= r.I64
			case ">=":
				t = l.I64 >= r.I64
			}
			return &Value{Kind: ValBool, B: t}, nil
		}
		return nil, newError("comparison expects numeric operands", node)
	case "+":
		if l.Kind == ValVec && r.Kind == ValVec {
			if len(l.Vec) != len(r.Vec) {
				return nil, newError("vector `+` expects equal dimensions", node)
			}
			out := make([]float64, len(l.Vec))
			for i := range l.Vec {
				out[i] = l.Vec[i] + r.Vec[i]
			}
			return &Value{Kind: ValVec, Vec: out}, nil
		}
		if l.Kind == ValString || r.Kind == ValString {
			return &Value{Kind: ValString, Str: l.String() + r.String()}, nil
		}
		if l.Kind == ValArray && r.Kind == ValArray {
			out := make([]Value, 0, len(l.Elems)+len(r.Elems))
			out = append(out, l.Elems...)
			out = append(out, r.Elems...)
			return &Value{Kind: ValArray, Elems: out}, nil
		}
		if !isNumeric(l) || !isNumeric(r) {
			return nil, newError("`+` expects numeric operands, two strings, or two arrays", node)
		}
		if useF {
			return &Value{Kind: ValFloat, F64: lf + rf}, nil
		}
		if l.Kind == ValInt && r.Kind == ValInt {
			return &Value{Kind: ValInt, I64: l.I64 + r.I64}, nil
		}
		return &Value{Kind: ValFloat, F64: lf + rf}, nil
	case "-", "*", "/", "%", "mod":
		if l.Kind == ValVec && r.Kind == ValVec && op == "-" {
			if len(l.Vec) != len(r.Vec) {
				return nil, newError("vector `-` expects equal dimensions", node)
			}
			out := make([]float64, len(l.Vec))
			for i := range l.Vec {
				out[i] = l.Vec[i] - r.Vec[i]
			}
			return &Value{Kind: ValVec, Vec: out}, nil
		}
		if !isNumeric(l) || !isNumeric(r) {
			return nil, newError("arithmetic expects numeric operands", node)
		}
		if useF {
			switch op {
			case "-":
				return &Value{Kind: ValFloat, F64: lf - rf}, nil
			case "*":
				return &Value{Kind: ValFloat, F64: lf * rf}, nil
			case "/":
				if math.Abs(rf) < 1e-20 {
					return nil, newError("div by 0", node)
				}
				return &Value{Kind: ValFloat, F64: lf / rf}, nil
			}
		}
		if l.Kind == ValInt && r.Kind == ValInt {
			switch op {
			case "-":
				return &Value{Kind: ValInt, I64: l.I64 - r.I64}, nil
			case "*":
				return &Value{Kind: ValInt, I64: l.I64 * r.I64}, nil
			case "/":
				if r.I64 == 0 {
					return nil, newError("div by 0", node)
				}
				return &Value{Kind: ValInt, I64: l.I64 / r.I64}, nil
			case "%", "mod":
				if r.I64 == 0 {
					return nil, newError("mod by 0", node)
				}
				return &Value{Kind: ValInt, I64: l.I64 % r.I64}, nil
			}
		}
	case "|", "&", "^", "<<", ">>":
		if l.Kind != ValInt || r.Kind != ValInt {
			return nil, newError("bitwise ops expect int operands", node)
		}
		switch op {
		case "|":
			return &Value{Kind: ValInt, I64: l.I64 | r.I64}, nil
		case "&":
			return &Value{Kind: ValInt, I64: l.I64 & r.I64}, nil
		case "^":
			return &Value{Kind: ValInt, I64: l.I64 ^ r.I64}, nil
		case "<<":
			return &Value{Kind: ValInt, I64: l.I64 << uint64(r.I64)}, nil
		case ">>":
			return &Value{Kind: ValInt, I64: l.I64 >> uint64(r.I64)}, nil
		}
	case "in":
		if r.Kind == ValArray {
			for i := range r.Elems {
				if valueEqual(l, &r.Elems[i]) {
					return &Value{Kind: ValBool, B: true}, nil
				}
			}
			return &Value{Kind: ValBool, B: false}, nil
		}
		if r.Kind == ValMap && r.StrMap != nil {
			ks, err := mapKeyString(l)
			if err != nil {
				return nil, err
			}
			if _, ok := r.StrMap[ks]; ok {
				return &Value{Kind: ValBool, B: true}, nil
			}
			for k := range r.StrMap {
				if strings.EqualFold(k, ks) {
					return &Value{Kind: ValBool, B: true}, nil
				}
			}
			return &Value{Kind: ValBool, B: false}, nil
		}
		if r.Kind == ValString {
			if l == nil || l.Kind != ValString {
				return nil, newError("`in` with string rhs expects string lhs", node)
			}
			return &Value{Kind: ValBool, B: strings.Contains(r.Str, l.Str)}, nil
		}
		return nil, newError("`in` expects array, map, or string on rhs", node)
	}
	if (op == "%" || op == "mod") && useF {
		return &Value{Kind: ValFloat, F64: math.Mod(lf, rf)}, nil
	}
	return nil, newError("bad infix: "+op, node)
}

func vecLikeFromValue(v *Value) ([]float64, bool) {
	if v == nil {
		return nil, false
	}
	if v.Kind == ValVec && len(v.Vec) > 0 {
		return append([]float64(nil), v.Vec...), true
	}
	if v.Kind != ValMap || v.StrMap == nil {
		return nil, false
	}
	get := func(name string) (float64, bool) {
		for k, vv := range v.StrMap {
			if strings.EqualFold(k, name) {
				switch vv.Kind {
				case ValFloat:
					return vv.F64, true
				case ValInt:
					return float64(vv.I64), true
				}
			}
		}
		return 0, false
	}
	x, okX := get("x")
	y, okY := get("y")
	if !okX || !okY {
		return nil, false
	}
	if z, okZ := get("z"); okZ {
		if w, okW := get("w"); okW {
			return []float64{x, y, z, w}, true
		}
		return []float64{x, y, z}, true
	}
	return []float64{x, y}, true
}

func isNumeric(v *Value) bool {
	return v != nil && (v.Kind == ValInt || v.Kind == ValFloat)
}

// toFloats2 promotes int/float mix to float.
func toFloats2(l, r *Value) (float64, float64) {
	if l.Kind == ValFloat {
		lf := l.F64
		var rf float64
		switch r.Kind {
		case ValFloat:
			rf = r.F64
		case ValInt:
			rf = float64(r.I64)
		default:
			return lf, 0
		}
		return lf, rf
	}
	if r.Kind == ValFloat {
		var lf float64
		if l.Kind == ValInt {
			lf = float64(l.I64)
		} else {
			return 0, r.F64
		}
		return lf, r.F64
	}
	if l.Kind == ValInt && r.Kind == ValInt {
		return float64(l.I64), float64(r.I64)
	}
	return 0, 0
}

func toFloats(l, r *Value) (a, b float64, useFloat bool) {
	if l.Kind == ValFloat || r.Kind == ValFloat {
		a, b = toFloats2(l, r)
		return a, b, true
	}
	if l.Kind == ValInt && r.Kind == ValInt {
		return float64(l.I64), float64(r.I64), false
	}
	a, b = toFloats2(l, r)
	return a, b, l.Kind == ValFloat || r.Kind == ValFloat
}

func valueEqual(l, r *Value) bool {
	if l == nil && r == nil {
		return true
	}
	if l == nil || r == nil {
		return false
	}
	if l.Kind != r.Kind {
		lf, rf := toFloats2(l, r)
		if (l.Kind == ValInt || l.Kind == ValFloat) && (r.Kind == ValInt || r.Kind == ValFloat) {
			return math.Abs(lf-rf) < 1e-9
		}
		return false
	}
	switch l.Kind {
	case ValInt:
		return l.I64 == r.I64
	case ValFloat:
		return math.Abs(l.F64-r.F64) < 1e-9
	case ValString:
		return l.Str == r.Str
	case ValBool:
		return l.B == r.B
	case ValNull:
		return true
	case ValArray:
		if len(l.Elems) != len(r.Elems) {
			return false
		}
		for i := range l.Elems {
			if !valueEqual(&l.Elems[i], &r.Elems[i]) {
				return false
			}
		}
		return true
	case ValVec:
		if len(l.Vec) != len(r.Vec) {
			return false
		}
		for i := range l.Vec {
			if math.Abs(l.Vec[i]-r.Vec[i]) > 1e-9 {
				return false
			}
		}
		return true
	case ValResult:
		if l.ResOk != r.ResOk {
			return false
		}
		if l.ResOk {
			return l.Res != nil && r.Res != nil && valueEqual(l.Res, r.Res)
		}
		return l.Err != nil && r.Err != nil && valueEqual(l.Err, r.Err)
	default:
		return false
	}
}

func evalEnumStatement(t *candy_ast.EnumStatement, e *Env) error {
	if t == nil || t.Name == nil {
		return &RuntimeError{Msg: "enum: missing name"}
	}
	consts := map[string]*Value{}
	next := int64(0)
	for _, v := range t.Variants {
		if v == nil || v.Name == nil {
			continue
		}
		if v.Value != nil {
			ev, err := evalExpression(v.Value, e)
			if err != nil {
				return err
			}
			if ev != nil && ev.Kind == ValInt {
				next = ev.I64
			}
			if ev == nil {
				ev = &Value{Kind: ValNull}
			}
			consts[v.Name.Value] = ev
		} else {
			consts[v.Name.Value] = &Value{Kind: ValInt, I64: next}
		}
		next++
	}
	e.Set(t.Name.Value, &Value{
		Kind: ValModule,
		Mod: &moduleVal{
			Name:   t.Name.Value,
			Consts: consts,
		},
	})
	return nil
}
