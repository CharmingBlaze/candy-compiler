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
		v, err := evalExpression(t.Value, e)
		if err != nil {
			return nil, err
		}
		e.Set(t.Name.Value, v)
		return v, nil
	case *candy_ast.DeferStatement:
		e.Defers = append(e.Defers, t)
		return nil, nil
	case *candy_ast.ReturnStatement:
		if t.ReturnValue == nil {
			return nil, nil
		}
		v, err := evalExpression(t.ReturnValue, e)
		if err != nil {
			return nil, err
		}
		return ReturnWrap{V: v}, nil
	case *candy_ast.BlockStatement:
		ne := e.NewEnclosed()
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
		}
		runDefers(ne)
		return nil, nil
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
		return nil, evalImport(t.Path, e)
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
		for _, a := range t.Arms {
			c, err := evalExpression(a.Cond, e)
			if err != nil {
				return nil, err
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
			p, err2 := evalExpression(b.Pat, e)
			if err2 != nil {
				return nil, err2
			}
			if valueEqual(sub, p) {
				return evalExpression(b.Body, e)
			}
		}
		if t.Default != nil {
			return evalExpression(t.Default, e)
		}
		return &Value{Kind: ValNull}, nil
	case *candy_ast.IndexExpression:
		le, err := evalExpression(t.Base, e)
		if err != nil {
			return nil, err
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
		args := make([]*Value, 0, len(t.Arguments))
		for _, a := range t.Arguments {
			av, err := evalExpression(a, e)
			if err != nil {
				return nil, err
			}
			args = append(args, av)
		}
		if id, ok := t.Function.(*candy_ast.Identifier); ok {
			if f, ok2 := Builtins[strings.ToLower(id.Value)]; ok2 {
				return f(args)
			}
		}
		fnv, err := evalExpression(t.Function, e)
		if err != nil {
			return nil, err
		}
		if fnv == nil {
			return nil, newError("nil call", t)
		}
		if fnv.Kind == ValStruct && (fnv.St.Def != nil || fnv.St.ClassDef != nil) {
			return instantiateClass(fnv, args, e)
		}
		return evalUserFunction(fnv, args)
	default:
		return nil, newError(fmt.Sprintf("unhandled expression %T", t), ex)
	}
}

func instantiateClass(def *Value, args []*Value, e *Env) (*Value, error) {
	instance := &Value{Kind: ValStruct, St: &structVal{Data: make(map[string]Value)}}
	if def.St.Def != nil {
		instance.St.Def = def.St.Def
		// existing struct logic...
		for i, p := range def.St.Def.Fields {
			if i < len(args) {
				instance.St.Data[p.Name.Value] = *args[i]
			} else if p.Init != nil {
				v, _ := evalExpression(p.Init, e)
				instance.St.Data[p.Name.Value] = *v
			}
		}
	} else if def.St.ClassDef != nil {
		instance.St.ClassDef = def.St.ClassDef
		// Primary constructor parameters -> Fields
		for i, p := range def.St.ClassDef.Parameters {
			if i < len(args) {
				instance.St.Data[p.Name.Value] = *args[i]
			} else if p.Default != nil {
				v, _ := evalExpression(p.Default, e)
				instance.St.Data[p.Name.Value] = *v
			}
		}
		// Evaluate class body members (fields/methods)
		for _, m := range def.St.ClassDef.Members {
			switch st := m.(type) {
			case *candy_ast.VarStatement:
				v, _ := evalExpression(st.Value, e)
				instance.St.Data[st.Name.Value] = *v
			case *candy_ast.ValStatement:
				v, _ := evalExpression(st.Value, e)
				instance.St.Data[st.Name.Value] = *v
			case *candy_ast.ExpressionStatement:
				// Support object/class field initializers written as `x = 0`.
				if asg, ok := st.Expression.(*candy_ast.AssignExpression); ok {
					if id, ok2 := asg.Left.(*candy_ast.Identifier); ok2 {
						v, _ := evalExpression(asg.Value, e)
						if v != nil {
							instance.St.Data[id.Value] = *v
						}
					}
				}
			}
			// Methods are handled by lookup in ClassDef during evalMethodCall
		}
	}
	return instance, nil
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
	for _, st := range fnv.Fn.Stmt.Body.Statements {
		r, err2 := evalStatement(st, ne)
		if err2 != nil {
			return nil, err2
		}
		if rw, ok2 := r.(ReturnWrap); ok2 {
			runDefers(ne)
			return rw.V, nil
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

func evalImport(path string, env *Env) error {
	if env == nil {
		return &RuntimeError{Msg: "nil env for import"}
	}
	full := path
	if src, ok := candy_stdlib.Lookup(path); ok {
		l := candy_lexer.New(src)
		p := candy_parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) > 0 {
			var msgs []string
			for _, d := range p.Errors() {
				msgs = append(msgs, d.Message)
			}
			return &RuntimeError{Msg: "stdlib import parse error: " + strings.Join(msgs, "; ")}
		}
		_, eerr := Eval(prog, env)
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
