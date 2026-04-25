package candy_typecheck

import "candy/candy_ast"

// checkStructBody walks struct members: attributes, field inits, properties, methods, operator overloads.
// The struct is already registered on the Checker in pass 1.
func (c *Checker) checkStructBody(st *candy_ast.StructStatement) {
	if st == nil {
		return
	}
	for _, a := range st.Attributes {
		c.walkAttributeArgs(a)
	}
	for i := range st.Fields {
		f := &st.Fields[i]
		for _, fa := range f.Attributes {
			c.walkAttributeArgs(fa)
		}
		if f.Init != nil {
			c.expr(f.Init)
		}
	}
	for _, p := range st.Properties {
		if p == nil {
			continue
		}
		if p.Getter != nil {
			c.pushScope()
			for _, x := range p.Getter.Statements {
				c.stmt(x)
			}
			c.popScope()
		}
		if p.Setter != nil {
			c.pushScope()
			for _, x := range p.Setter.Statements {
				c.stmt(x)
			}
			c.popScope()
		}
		c.expr(p.DefaultValue)
	}
	for _, m := range st.Methods {
		if m == nil {
			continue
		}
		for _, a := range m.Attributes {
			c.walkAttributeArgs(a)
		}
		c.walkFunctionLike(m)
	}
	for _, op := range st.Operators {
		if op == nil {
			continue
		}
		c.walkOperatorOverload(op)
	}
}

func (c *Checker) walkOperatorOverload(o *candy_ast.OperatorOverloadStatement) {
	if o == nil {
		return
	}
	c.pushScope()
	for _, p := range o.Parameters {
		if p.Name != nil && candy_ast.ExprAsSimpleTypeName(p.TypeName) != "" {
			c.bind(p.Name.Value, candy_ast.ExprAsSimpleTypeName(p.TypeName))
		}
	}
	if o.ReturnType != nil {
		c.returnTypes = append(c.returnTypes, canonType(candy_ast.ExprAsSimpleTypeName(o.ReturnType)))
	} else {
		c.returnTypes = append(c.returnTypes, "any")
	}
	if o.Body != nil {
		for _, x := range o.Body.Statements {
			c.stmt(x)
		}
	}
	if len(c.returnTypes) > 0 {
		c.returnTypes = c.returnTypes[:len(c.returnTypes)-1]
	}
	c.popScope()
}
