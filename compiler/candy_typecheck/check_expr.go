package candy_typecheck

import (
	"candy/candy_ast"
	"fmt"
)

// litKind classifies supported literal sub-expressions (unknown = not a simple literal).
type litKind int

const (
	litUnknown litKind = iota
	litInt
	litFloat
	litString
	litBool
	litNull
)

func literalKind(e candy_ast.Expression) litKind {
	if e == nil {
		return litUnknown
	}
	switch t := e.(type) {
	case *candy_ast.IntegerLiteral:
		_ = t
		return litInt
	case *candy_ast.FloatLiteral:
		return litFloat
	case *candy_ast.StringLiteral:
		return litString
	case *candy_ast.Boolean:
		return litBool
	case *candy_ast.NullLiteral:
		return litNull
	default:
		return litUnknown
	}
}

func (c *Checker) checkInfix(t *candy_ast.InfixExpression) {
	if t == nil {
		return
	}
	c.expr(t.Left)
	c.expr(t.Right)
	lt := c.inferExprType(t.Left)
	rt := c.inferExprType(t.Right)
	if _, o := c.findOperatorOverload(canonType(lt), t.Operator); o != nil {
		if len(o.Parameters) > 0 && o.Parameters[0].TypeName != nil {
			want := canonType(candy_ast.ExprAsSimpleTypeName(o.Parameters[0].TypeName))
			if rt != "any" && !c.typeAssignable(want, rt) {
				c.add(fmt.Sprintf("operator %q: right operand type mismatch: expected %s, got %s", t.Operator, want, rt), t)
			}
		}
		return
	}
	l, r := literalKind(t.Left), literalKind(t.Right)
	if l == litUnknown || r == litUnknown {
		return
	}
	switch t.Operator {
	case "+":
		if l == litBool || r == litBool {
			c.add("invalid `+` with boolean literal operand", t)
		}
		if l == litNull || r == litNull {
			c.add("invalid `+` with null literal operand", t)
		}
	case "-", "*", "/":
		if l == litBool || r == litBool || l == litString || r == litString || l == litNull || r == litNull {
			c.add("arithmetic is not defined for these literal kinds", t)
		}
	}
}

// checkAssignExpression walks assignment; validates property/field set value types for dot left-hand side.
func (c *Checker) checkAssignExpression(t *candy_ast.AssignExpression) {
	if t == nil {
		return
	}
	c.expr(t.Value)
	if dot, ok := t.Left.(*candy_ast.DotExpression); ok {
		c.expr(dot.Left)
		lType := c.inferExprType(dot.Left)
		if lType == "any" || lType == "builtin" {
			return
		}
		if memTy, found := c.findMember(lType, dot.Right.Value); found {
			valTy := c.inferExprType(t.Value)
			if valTy != "any" && !c.typeAssignable(memTy, valTy) {
				c.add(fmt.Sprintf("cannot assign %s to %s.%s (expected %s)", valTy, lType, dot.Right.Value, memTy), t)
			}
			return
		}
		c.add(fmt.Sprintf("type %s has no field or property %s", lType, dot.Right.Value), t)
		return
	}
	if id, ok := t.Left.(*candy_ast.Identifier); ok {
		if _, found := c.lookup(id.Value); !found {
			// Implicit declaration!
			rhsType := c.inferExprType(t.Value)
			if rhsType == "" {
				rhsType = "any"
			}
			c.bind(id.Value, rhsType)
			return
		}
	}
	c.expr(t.Left)
}
