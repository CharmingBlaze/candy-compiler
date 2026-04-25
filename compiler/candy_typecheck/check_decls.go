package candy_typecheck

import (
	"candy/candy_ast"
	"fmt"
)

func (c *Checker) handleDeclStatement(s candy_ast.Statement) bool {
	switch t := s.(type) {
	case *candy_ast.ImportStatement:
		return true // no-op
	case *candy_ast.PackageStatement:
		return true // no-op
	case *candy_ast.StructStatement:
		if t.Name != nil {
			if len(t.TypeParameters) > 0 {
				c.genericStructs[canonType(t.Name.Value)] = t
			} else {
				c.structs[canonType(t.Name.Value)] = t
				c.bind(t.Name.Value, "type:"+canonType(t.Name.Value))
			}
		}
		return true
	case *candy_ast.ClassStatement:
		if t.Name != nil {
			name := canonType(t.Name.Value)
			c.classes[name] = t
			c.bind(t.Name.Value, "class:"+name)
			if t.Base != nil {
				base := canonType(t.Base.Value)
				if _, ok := c.classes[base]; !ok {
					c.add(fmt.Sprintf("unknown base class: %s", base), t)
				}
			}
		}
		return true
	case *candy_ast.InterfaceStatement:
		if t.Name != nil {
			c.interfaces[canonType(t.Name.Value)] = t
			c.bind(t.Name.Value, "iface:"+canonType(t.Name.Value))
		}
		return true
	case *candy_ast.TraitStatement:
		if t.Name != nil {
			c.traits[canonType(t.Name.Value)] = t
			c.bind(t.Name.Value, "trait:"+canonType(t.Name.Value))
		}
		return true
	case *candy_ast.ObjectStatement:
		if t.Name != nil {
			c.bind(t.Name.Value, "object:"+canonType(t.Name.Value))
		}
		return true
	case *candy_ast.ExternFunctionStatement:
		if t.Function != nil && t.Function.Name != nil {
			c.bind(t.Function.Name.Value, "extern-fn")
			validateExternSignature(c, t)
		}
		return true
	case *candy_ast.ModuleStatement:
		return true
	case *candy_ast.LibraryStatement:
		if t.Body != nil {
			for _, st := range t.Body.Statements {
				c.stmt(st)
			}
		}
		return true
	case *candy_ast.EnumStatement:
		if t.Name != nil {
			n := canonType(t.Name.Value)
			c.enums[n] = t
			c.bind(t.Name.Value, "enum:"+n)
		}
		return true
	case *candy_ast.TryStatement, *candy_ast.RunStatement:
		return true
	default:
		return false
	}
}

func validateExternSignature(c *Checker, ex *candy_ast.ExternFunctionStatement) {
	if c == nil || ex == nil || ex.Function == nil {
		return
	}
	isSupported := func(typeExpr candy_ast.Expression) bool {
		name := canonType(candy_ast.ExprAsSimpleTypeName(typeExpr))
		switch name {
		case "", "int", "float", "double", "bool", "string", "void":
			return true
		default:
			return false
		}
	}
	for _, p := range ex.Function.Parameters {
		if !isSupported(p.TypeName) {
			c.add(fmt.Sprintf("extern param type not supported yet: %s", candy_ast.StringExpr(p.TypeName)), ex)
		}
	}
	if !isSupported(ex.Function.ReturnType) {
		c.add(fmt.Sprintf("extern return type not supported yet: %s", candy_ast.StringExpr(ex.Function.ReturnType)), ex)
	}
}

func tokenOfDeclNode(n candy_ast.Node) (candy_ast.Node, bool) {
	switch n.(type) {
	case *candy_ast.PackageStatement,
		*candy_ast.ClassStatement,
		*candy_ast.ObjectStatement,
		*candy_ast.InterfaceStatement,
		*candy_ast.TraitStatement,
		*candy_ast.ExternFunctionStatement:
		return n, true
	default:
		return nil, false
	}
}
