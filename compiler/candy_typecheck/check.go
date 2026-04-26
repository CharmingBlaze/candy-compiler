package candy_typecheck

import (
	"candy/candy_ast"
	"candy/candy_report"
	"candy/candy_token"
	"fmt"
	"strings"
)

// Checker is a placeholder for Hindley–Milner / full typing; for now a shallow walk.
type Checker struct {
	Issues             []candy_report.Diagnostic
	scopes             []map[string]string
	structs            map[string]*candy_ast.StructStatement
	enums              map[string]*candy_ast.EnumStatement
	classes            map[string]*candy_ast.ClassStatement
	interfaces         map[string]*candy_ast.InterfaceStatement
	traits             map[string]*candy_ast.TraitStatement
	genericStructs     map[string]*candy_ast.StructStatement
	genericFunctions   map[string]*candy_ast.FunctionStatement
	SpecializedStructs []*candy_ast.StructStatement
	returnTypes        []string
}

func (c *Checker) add(msg string, n candy_ast.Node) {
	tok, _ := tokenOfNode(n)
	c.Issues = append(c.Issues, candy_report.Diagnostic{
		Level:   candy_report.Warning, // Typecheck issues currently treated as warnings/non-fatal
		Message: msg,
		Line:    tok.Line,
		Col:     tok.Col,
		Offset:  tok.Offset,
		Length:  len(tok.Literal),
	})
}

func canonType(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return "any"
	}
	return s
}

// CheckProgram walks the AST.
func CheckProgram(p *candy_ast.Program) []candy_report.Diagnostic {
	if p == nil {
		return nil
	}
	c := &Checker{
		structs:          make(map[string]*candy_ast.StructStatement),
		enums:            make(map[string]*candy_ast.EnumStatement),
		classes:          make(map[string]*candy_ast.ClassStatement),
		interfaces:       make(map[string]*candy_ast.InterfaceStatement),
		traits:           make(map[string]*candy_ast.TraitStatement),
		genericStructs:   make(map[string]*candy_ast.StructStatement),
		genericFunctions: make(map[string]*candy_ast.FunctionStatement),
	}

	c.pushScope()
	c.bind("println", "builtin")
	c.bind("print", "builtin")

	// Pass 1: Collect declarations
	for _, s := range p.Statements {
		c.handleDeclStatement(s)
	}

	c.pushScope()
	c.bind("len", "builtin")
	c.bind("ok", "builtin")
	c.bind("err", "builtin")
	c.bind("clock", "builtin")
	c.bind("exit", "builtin")
	c.bind("vec2", "builtin")
	c.bind("vec3", "builtin")
	c.bind("vec4", "builtin")
	c.bind("format", "builtin")
	c.bind("enumerate", "builtin")
	c.bind("box", "builtin")
	c.bind("aabb", "builtin")
	c.bind("sphere", "builtin")
	c.bind("ray", "builtin")
	c.bind("physicsworld", "builtin")
	c.bind("inputmap", "builtin")
	c.bind("orbitcamera", "builtin")
	c.bind("firstpersoncamera", "builtin")
	c.bind("charactercontroller", "builtin")
	c.bind("entitylist", "builtin")
	c.bind("uilayout", "builtin")
	c.bind("hud", "builtin")
	c.bind("statemachine", "builtin")
	c.bind("tween", "builtin")
	c.bind("transform", "builtin")
	c.bind("drawall", "builtin")
	c.bind("gameloop", "builtin")

	// Pass 2: Typecheck statements
	for _, s := range p.Statements {
		c.stmt(s)
	}

	// Append specialized structs so backend sees them
	for _, s := range c.SpecializedStructs {
		p.Statements = append(p.Statements, s)
	}

	return c.Issues
}

func (c *Checker) stmt(s candy_ast.Statement) {
	if s == nil {
		return
	}
	switch t := s.(type) {
	case *candy_ast.FunctionStatement:
		if t.Name != nil {
			if len(t.TypeParameters) > 0 {
				c.genericFunctions[canonType(t.Name.Value)] = t
			} else {
				rt := "Any"
				if candy_ast.ExprAsSimpleTypeName(t.ReturnType) != "" {
					rt = canonType(candy_ast.ExprAsSimpleTypeName(t.ReturnType))
				}
				c.bind(t.Name.Value, "fn:"+rt)
			}
		}
		for _, a := range t.Attributes {
			c.walkAttributeArgs(a)
		}
		// IsAsync, Exported, Suspend: surface syntax only; no async/visibility semantics in this checker yet.
		c.walkFunctionLike(t)
	case *candy_ast.ModuleStatement:
		c.pushScope()
		if t.Body != nil {
			for _, x := range t.Body.Statements {
				c.stmt(x)
			}
		}
		c.popScope()
	case *candy_ast.EnumStatement:
		for _, v := range t.Variants {
			if v != nil && v.Value != nil {
				c.expr(v.Value)
			}
		}
	case *candy_ast.TryStatement:
		if t.TryBody != nil {
			c.pushScope()
			for _, x := range t.TryBody.Statements {
				c.stmt(x)
			}
			c.popScope()
		}
		for _, cc := range t.CatchClauses {
			if cc == nil {
				continue
			}
			c.pushScope()
			if cc.Type != nil && cc.Identifier != nil {
				c.bind(cc.Identifier.Value, candy_ast.ExprAsSimpleTypeName(cc.Type))
			}
			if cc.Body != nil {
				for _, x := range cc.Body.Statements {
					c.stmt(x)
				}
			}
			c.popScope()
		}
		if t.FinallyBody != nil {
			c.pushScope()
			for _, x := range t.FinallyBody.Statements {
				c.stmt(x)
			}
			c.popScope()
		}
	case *candy_ast.RunStatement:
		c.expr(t.Value)
	case *candy_ast.WithStatement:
		c.pushScope()
		c.expr(t.Value)
		if t.Name != nil {
			c.bind(t.Name.Value, c.inferExprType(t.Value))
		}
		if t.Body != nil {
			for _, x := range t.Body.Statements {
				c.stmt(x)
			}
		}
		c.popScope()
	case *candy_ast.StructStatement:
		c.checkStructBody(t)
	case *candy_ast.BlockStatement:
		c.pushScope()
		for _, x := range t.Statements {
			c.stmt(x)
		}
		c.popScope()
	case *candy_ast.IfExpression:
		c.walkIfExpression(t)
	case *candy_ast.ExpressionStatement:
		c.expr(t.Expression)
	case *candy_ast.ValStatement:
		c.expr(t.Value)
		got := c.inferExprType(t.Value)
		if t.TypeName != nil {
			s := candy_ast.ExprAsSimpleTypeName(t.TypeName)
			if s == "" {
				s = "any"
			}
			exp := canonType(s)
			if got != "any" && !c.typeAssignable(exp, got) {
				c.add(fmt.Sprintf("type mismatch: cannot assign %s to %s", got, exp), t)
			}
			if t.Name != nil {
				c.bind(t.Name.Value, exp)
			}
		} else if t.Name != nil {
			c.bind(t.Name.Value, got)
		}
	case *candy_ast.VarStatement:
		c.expr(t.Value)
		got := c.inferExprType(t.Value)
		if t.TypeName != nil {
			s := candy_ast.ExprAsSimpleTypeName(t.TypeName)
			if s == "" {
				s = "any"
			}
			exp := canonType(s)
			if got != "any" && !c.typeAssignable(exp, got) {
				c.add(fmt.Sprintf("type mismatch: cannot assign %s to %s", got, exp), t)
			}
			if t.Name != nil {
				c.bind(t.Name.Value, exp)
			}
		} else if t.Name != nil {
			c.bind(t.Name.Value, got)
		}
	case *candy_ast.ReturnStatement:
		c.expr(t.ReturnValue)
		got := c.inferExprType(t.ReturnValue)
		exp := c.currentReturnType()
		if exp != "" && exp != "any" && got != "any" && !c.typeAssignable(exp, got) {
			c.add(fmt.Sprintf("return type mismatch: expected %s, got %s", exp, got), t)
		}
	case *candy_ast.SwitchStatement:
		c.expr(t.Subject)
		for _, arm := range t.Cases {
			for _, pat := range arm.Patterns {
				c.expr(pat)
			}
			if arm.Body == nil {
				continue
			}
			if bs, ok := arm.Body.(*candy_ast.BlockStatement); ok {
				for _, x := range bs.Statements {
					c.stmt(x)
				}
				continue
			}
			c.stmt(arm.Body)
		}
	case *candy_ast.DoWhileStatement:
		if t.Body != nil {
			for _, x := range t.Body.Statements {
				c.stmt(x)
			}
		}
		c.expr(t.Condition)
	default:
		_ = c.handleDeclStatement(s)
	}
}

func returnsNullLiteral(b *candy_ast.BlockStatement) bool {
	if b == nil {
		return false
	}
	for _, st := range b.Statements {
		switch t := st.(type) {
		case *candy_ast.ReturnStatement:
			if _, ok := t.ReturnValue.(*candy_ast.NullLiteral); ok {
				return true
			}
		case *candy_ast.BlockStatement:
			if returnsNullLiteral(t) {
				return true
			}
		case *candy_ast.IfExpression:
			if statementReturnsNullLiteral(t.Consequence) || statementReturnsNullLiteral(t.Alternative) {
				return true
			}
		}
	}
	return false
}

func statementReturnsNullLiteral(s candy_ast.Statement) bool {
	if s == nil {
		return false
	}
	if b, ok := s.(*candy_ast.BlockStatement); ok {
		return returnsNullLiteral(b)
	}
	if ie, ok := s.(*candy_ast.IfExpression); ok {
		if statementReturnsNullLiteral(ie.Consequence) || statementReturnsNullLiteral(ie.Alternative) {
			return true
		}
	}
	return false
}

func tokenOfNode(n candy_ast.Node) (candy_token.Token, bool) {
	switch t := n.(type) {
	case *candy_ast.FunctionStatement:
		return t.Token, true
	case *candy_ast.BlockStatement:
		return t.Token, true
	case *candy_ast.IfExpression:
		return t.Token, true
	case *candy_ast.ExpressionStatement:
		return t.Token, true
	case *candy_ast.ValStatement:
		return t.Token, true
	case *candy_ast.VarStatement:
		return t.Token, true
	case *candy_ast.ReturnStatement:
		return t.Token, true
	case *candy_ast.ImportStatement:
		return t.Token, true
	case *candy_ast.StructStatement:
		return t.Token, true
	case *candy_ast.PackageStatement:
		return t.Token, true
	case *candy_ast.ClassStatement:
		return t.Token, true
	case *candy_ast.ObjectStatement:
		return t.Token, true
	case *candy_ast.InterfaceStatement:
		return t.Token, true
	case *candy_ast.TraitStatement:
		return t.Token, true
	case *candy_ast.ExternFunctionStatement:
		return t.Token, true
	case *candy_ast.ModuleStatement:
		return t.Token, true
	case *candy_ast.EnumStatement:
		return t.Token, true
	case *candy_ast.TryStatement:
		return t.Token, true
	case *candy_ast.RunStatement:
		return t.Token, true
	case *candy_ast.AwaitExpression:
		return t.Token, true
	case *candy_ast.AssignExpression:
		return t.Token, true
	case *candy_ast.IsExpression:
		return t.Token, true
	case *candy_ast.LambdaExpression:
		return t.Token, true
	case *candy_ast.TupleLiteral:
		return t.Token, true
	case *candy_ast.TupleTypeExpression:
		return t.Token, true
	case *candy_ast.PrefixExpression:
		return t.Token, true
	case *candy_ast.PostfixExpression:
		return t.Token, true
	case *candy_ast.InfixExpression:
		return t.Token, true
	case *candy_ast.GroupedExpression:
		return t.Token, true
	case *candy_ast.CallExpression:
		return t.Token, true
	case *candy_ast.ArrayLiteral:
		return t.Token, true
	case *candy_ast.MapLiteral:
		return t.Token, true
	case *candy_ast.IndexExpression:
		return t.Token, true
	case *candy_ast.WhenExpression:
		return t.Token, true
	case *candy_ast.MatchExpression:
		return t.Token, true
	case *candy_ast.IntegerLiteral:
		return t.Token, true
	case *candy_ast.FloatLiteral:
		return t.Token, true
	case *candy_ast.StringLiteral:
		return t.Token, true
	case *candy_ast.Boolean:
		return t.Token, true
	case *candy_ast.NullLiteral:
		return t.Token, true
	case *candy_ast.Identifier:
		return t.Token, true
	default:
		return candy_token.Token{}, false
	}
}

func (c *Checker) expr(n candy_ast.Expression) {
	if n == nil {
		return
	}
	switch t := n.(type) {
	case *candy_ast.PrefixExpression:
		c.expr(t.Right)
	case *candy_ast.PostfixExpression:
		c.expr(t.Left)
	case *candy_ast.InfixExpression:
		c.checkInfix(t)
	case *candy_ast.GroupedExpression:
		c.expr(t.Expr)
	case *candy_ast.CallExpression:
		c.expr(t.Function)
		for _, a := range t.Arguments {
			c.expr(a)
		}
		if id, ok := t.Function.(*candy_ast.Identifier); ok {
			switch strings.ToLower(id.Value) {
			case "len", "ok", "err":
				if len(t.Arguments) != 1 {
					c.add(fmt.Sprintf("%s expects exactly 1 argument", id.Value), t)
				}
			}
		}
		for _, ta := range t.TypeArguments {
			c.expr(ta)
		}
	case *candy_ast.ArrayLiteral:
		for _, e := range t.Elem {
			c.expr(e)
		}
	case *candy_ast.MapLiteral:
		for _, p := range t.Pairs {
			c.expr(p.Key)
			c.expr(p.Value)
		}
	case *candy_ast.IndexExpression:
		c.expr(t.Base)
		c.expr(t.Index)
	case *candy_ast.MatchExpression:
		c.expr(t.Subject)
		want := "any"
		for _, b := range t.Branches {
			c.expr(b.Pat)
			c.expr(b.Body)
			bt := c.inferExprType(b.Body)
			if want == "any" {
				want = bt
			} else if bt != "any" && !c.typeAssignable(want, bt) && !c.typeAssignable(bt, want) {
				c.add("match branches should produce consistent types", t)
			}
		}
		c.expr(t.Default)
		if t.Default != nil {
			dt := c.inferExprType(t.Default)
			if want != "any" && dt != "any" && !c.typeAssignable(want, dt) && !c.typeAssignable(dt, want) {
				c.add("match default type should match branch result types", t)
			}
		}
	case *candy_ast.WhenExpression:
		want := "any"
		for _, a := range t.Arms {
			c.expr(a.Cond)
			c.expr(a.Body)
			bt := c.inferExprType(a.Body)
			if want == "any" {
				want = bt
			} else if bt != "any" && !c.typeAssignable(want, bt) && !c.typeAssignable(bt, want) {
				c.add("when arms should produce consistent result types", t)
			}
		}
		c.expr(t.ElseV)
		if t.ElseV != nil {
			et := c.inferExprType(t.ElseV)
			if want != "any" && et != "any" && !c.typeAssignable(want, et) && !c.typeAssignable(et, want) {
				c.add("when else type should match arm result types", t)
			}
		}
	case *candy_ast.Identifier:
		if _, ok := c.lookup(t.Value); !ok {
			c.add(fmt.Sprintf("unknown identifier: %s", t.Value), t)
		}
	case *candy_ast.AssignExpression:
		c.checkAssignExpression(t)
	case *candy_ast.AwaitExpression:
		c.expr(t.Value)
	case *candy_ast.IfExpression:
		c.walkIfExpression(t)
	case *candy_ast.IsExpression:
		c.expr(t.Left)
	case *candy_ast.LambdaExpression:
		c.pushScope()
		for _, p := range t.Parameters {
			if p.Name != nil && candy_ast.ExprAsSimpleTypeName(p.TypeName) != "" {
				c.bind(p.Name.Value, candy_ast.ExprAsSimpleTypeName(p.TypeName))
			}
		}
		c.expr(t.Body)
		c.popScope()
	case *candy_ast.TypeExpression:
		c.resolveTypeExpression(t)
	case *candy_ast.TupleLiteral:
		for _, e := range t.Elems {
			c.expr(e)
		}
	case *candy_ast.StructLiteral:
		c.normalizeStructLiteralName(t)
		stName := canonType(candy_ast.ExprAsSimpleTypeName(t.Name))
		if _, ok := c.structs[stName]; !ok {
			c.add(fmt.Sprintf("unknown struct type: %s", stName), t)
			return
		}
		for fname, valExpr := range t.Fields {
			c.expr(valExpr)
			got := c.inferExprType(valExpr)
			exp, found := c.findMember(stName, fname)
			if !found {
				c.add(fmt.Sprintf("struct %s has no field %s", stName, fname), t)
				continue
			}
			if got != "any" && !c.typeAssignable(exp, got) {
				c.add(fmt.Sprintf("type mismatch for field %s: expected %s, got %s", fname, exp, got), valExpr)
			}
		}
	case *candy_ast.DotExpression:
		c.expr(t.Left)
		lType := c.inferExprType(t.Left)
		if strings.HasPrefix(lType, "vec") {
			return
		}
		if lType != "any" && lType != "builtin" {
			if _, found := c.findMember(lType, t.Right.Value); !found {
				c.add(fmt.Sprintf("type %s has no field or property %s", lType, t.Right.Value), t)
			}
		}
	}
}

func (c *Checker) pushScope() { c.scopes = append(c.scopes, map[string]string{}) }
func (c *Checker) popScope() {
	if len(c.scopes) > 0 {
		c.scopes = c.scopes[:len(c.scopes)-1]
	}
}
func (c *Checker) bind(name, ty string) {
	if len(c.scopes) == 0 {
		c.pushScope()
	}
	c.scopes[len(c.scopes)-1][strings.ToLower(name)] = canonType(ty)
}
func (c *Checker) lookup(name string) (string, bool) {
	key := strings.ToLower(name)
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if t, ok := c.scopes[i][key]; ok {
			return t, true
		}
	}
	return "", false
}
func (c *Checker) currentReturnType() string {
	if len(c.returnTypes) == 0 {
		return ""
	}
	return c.returnTypes[len(c.returnTypes)-1]
}

// normalizeStructLiteralName rewrites `Box<float> { }` to a concrete struct name so checks and LLVM match monomorphization.
func (c *Checker) normalizeStructLiteralName(sl *candy_ast.StructLiteral) {
	if sl == nil {
		return
	}
	te, ok := sl.Name.(*candy_ast.TypeExpression)
	if !ok {
		return
	}
	stKey := c.resolveTypeExpression(te)
	if st, has := c.structs[stKey]; has && st.Name != nil {
		sl.Name = &candy_ast.Identifier{Token: te.Token, Value: st.Name.Value}
	}
}

func (c *Checker) inferExprType(e candy_ast.Expression) string {
	switch t := e.(type) {
	case *candy_ast.IntegerLiteral:
		return "int"
	case *candy_ast.FloatLiteral:
		return "float"
	case *candy_ast.StringLiteral:
		return "string"
	case *candy_ast.Boolean:
		return "bool"
	case *candy_ast.NullLiteral:
		return "null"
	case *candy_ast.Identifier:
		if ty, ok := c.lookup(t.Value); ok {
			return ty
		}
		return "any"
	case *candy_ast.InfixExpression:
		l, r := c.inferExprType(t.Left), c.inferExprType(t.Right)
		// Check for operator overloads
		if ty, found := c.findOperator(l, t.Operator); found {
			return ty
		}
		if strings.HasPrefix(l, "vec") || strings.HasPrefix(r, "vec") {
			if strings.HasPrefix(l, "vec") {
				return l
			}
			return r
		}
		if t.Operator == "+" && l == "string" && r == "string" {
			return "string"
		}
		if l == "float" || r == "float" {
			return "float"
		}
		if l == "int" && r == "int" {
			return "int"
		}
		return "any"
	case *candy_ast.StructLiteral:
		c.normalizeStructLiteralName(t)
		return canonType(candy_ast.ExprAsSimpleTypeName(t.Name))
	case *candy_ast.DotExpression:
		leftType := c.inferExprType(t.Left)
		if strings.HasPrefix(leftType, "vec") {
			switch strings.ToLower(t.Right.Value) {
			case "x", "y", "z", "w", "length", "dot", "cross", "normalize", "distance":
				if t.Right.Value == "normalize" {
					return leftType
				}
				return "float"
			case "xy":
				return "vec2"
			case "xz", "yz":
				return "vec2"
			}
		}
		if ty, found := c.findMember(leftType, t.Right.Value); found {
			return ty
		}
		return "any"
	case *candy_ast.IsExpression:
		return "bool"
	case *candy_ast.AwaitExpression:
		inner := c.inferExprType(t.Value)
		if inner == "" {
			return "any"
		}
		return inner
	case *candy_ast.AssignExpression:
		return c.inferExprType(t.Value)
	case *candy_ast.IfExpression:
		return "any"
	case *candy_ast.LambdaExpression:
		return "any"
	case *candy_ast.TupleLiteral:
		return "tuple"
	case *candy_ast.TypeExpression:
		return c.resolveTypeExpression(t)
	default:
		if ce, ok := e.(*candy_ast.CallExpression); ok {
			if id, ok2 := ce.Function.(*candy_ast.Identifier); ok2 {
				switch strings.ToLower(id.Value) {
				case "vec2":
					return "vec2"
				case "vec3":
					return "vec3"
				case "vec4":
					return "vec4"
				}
			}
		}
		return "any"
	}
}
func (c *Checker) findMember(stName, fName string) (string, bool) {
	st, ok := c.structs[canonType(stName)]
	if !ok {
		return "", false
	}
	// Check fields
	for _, f := range st.Fields {
		if strings.EqualFold(f.Name.Value, fName) {
			return canonType(candy_ast.ExprAsSimpleTypeName(f.TypeName)), true
		}
	}
	// Check properties
	for _, p := range st.Properties {
		if strings.EqualFold(p.Name.Value, fName) {
			return canonType(candy_ast.ExprAsSimpleTypeName(p.Type)), true
		}
	}
	// Check methods
	for _, m := range st.Methods {
		if strings.EqualFold(m.Name.Value, fName) {
			return "builtin", true // Treat method as a callable builtin for now
		}
	}
	// Check bases recursively
	for _, b := range st.Bases {
		if ty, found := c.findMember(b.Value, fName); found {
			return ty, found
		}
	}
	return "", false
}

func (c *Checker) isSubtype(child, parent string) bool {
	child = canonType(child)
	parent = canonType(parent)
	if child == parent {
		return true
	}
	st, ok := c.structs[child]
	if ok {
		for _, b := range st.Bases {
			if c.isSubtype(b.Value, parent) {
				return true
			}
		}
	}
	cl, ok := c.classes[child]
	if ok && cl.Base != nil {
		if c.isSubtype(cl.Base.Value, parent) {
			return true
		}
	}
	return false
}

func (c *Checker) typeAssignable(expect, got string) bool {
	expect = canonType(expect)
	got = canonType(got)
	if expect == got || expect == "any" || got == "any" {
		return true
	}
	if strings.HasSuffix(expect, "?") {
		inner := strings.TrimSuffix(expect, "?")
		if got == "null" || c.isSubtype(got, inner) {
			return true
		}
	}
	return c.isSubtype(got, expect)
}

func (c *Checker) findOperator(stName, op string) (string, bool) {
	if _, o := c.findOperatorOverload(canonType(stName), op); o != nil {
		if o.ReturnType == nil {
			return "any", true
		}
		return canonType(candy_ast.ExprAsSimpleTypeName(o.ReturnType)), true
	}
	return "", false
}

// findOperatorOverload returns the defining struct name and the operator statement, or nil.
func (c *Checker) findOperatorOverload(stName, op string) (string, *candy_ast.OperatorOverloadStatement) {
	n := c.structs[stName]
	if n == nil {
		return "", nil
	}
	for _, o := range n.Operators {
		if o.Operator == op {
			return stName, o
		}
	}
	for _, b := range n.Bases {
		if sn, o := c.findOperatorOverload(canonType(b.Value), op); o != nil {
			return sn, o
		}
	}
	return "", nil
}

// walkIfExpression is shared for if used as a statement or as a subexpression.
func (c *Checker) walkIfExpression(t *candy_ast.IfExpression) {
	if t == nil {
		return
	}
	c.expr(t.Condition)
	var narrowedVar string
	var narrowedTo string
	if isExpr, ok := t.Condition.(*candy_ast.IsExpression); ok {
		if ident, ok := isExpr.Left.(*candy_ast.Identifier); ok {
			narrowedVar = ident.Value
			narrowedTo = canonType(candy_ast.ExprAsSimpleTypeName(isExpr.TypeName))
		}
	}
	if t.Consequence != nil {
		if narrowedVar != "" {
			c.pushScope()
			c.bind(narrowedVar, narrowedTo)
		}
		c.stmt(t.Consequence)
		if narrowedVar != "" {
			c.popScope()
		}
	}
	if t.Alternative != nil {
		if altBlock, ok := t.Alternative.(*candy_ast.BlockStatement); ok {
			for _, x := range altBlock.Statements {
				c.stmt(x)
			}
		} else {
			c.stmt(t.Alternative)
		}
	}
}

// walkFunctionLike typechecks a function/method body: parameters, optional receiver, return stack.
// It does not bind the function name in the outer scope (the stmt handler does for top-level decls).
func (c *Checker) walkFunctionLike(fn *candy_ast.FunctionStatement) {
	if fn == nil {
		return
	}
	if candy_ast.ExprAsSimpleTypeName(fn.ReturnType) != "" && !strings.HasSuffix(canonType(candy_ast.ExprAsSimpleTypeName(fn.ReturnType)), "?") {
		if returnsNullLiteral(fn.Body) {
			c.add("function has non-nullable return type but returns null", fn)
		}
	}
	c.pushScope()
	if fn.Receiver != nil && fn.Receiver.Name != nil && candy_ast.ExprAsSimpleTypeName(fn.Receiver.TypeName) != "" {
		c.bind(fn.Receiver.Name.Value, candy_ast.ExprAsSimpleTypeName(fn.Receiver.TypeName))
	}
	for _, param := range fn.Parameters {
		if param.Name != nil && candy_ast.ExprAsSimpleTypeName(param.TypeName) != "" {
			c.bind(param.Name.Value, candy_ast.ExprAsSimpleTypeName(param.TypeName))
		}
	}
	if candy_ast.ExprAsSimpleTypeName(fn.ReturnType) != "" {
		c.returnTypes = append(c.returnTypes, canonType(candy_ast.ExprAsSimpleTypeName(fn.ReturnType)))
	} else {
		c.returnTypes = append(c.returnTypes, "any")
	}
	if fn.Body != nil {
		for _, x := range fn.Body.Statements {
			c.stmt(x)
		}
	}
	if len(c.returnTypes) > 0 {
		c.returnTypes = c.returnTypes[:len(c.returnTypes)-1]
	}
	c.popScope()
}

func (c *Checker) walkAttributeArgs(a *candy_ast.Attribute) {
	if a == nil {
		return
	}
	for _, e := range a.Arguments {
		c.expr(e)
	}
}
