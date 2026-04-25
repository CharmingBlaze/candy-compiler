package candy_llvm

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

// findPropertyDefining walks the struct hierarchy and returns the owning struct name
// and the property statement, if any.
func (c *Compiler) findPropertyDefining(st *candy_ast.StructStatement, name string) (owning string, prop *candy_ast.PropertyStatement) {
	if st == nil {
		return "", nil
	}
	for _, p := range st.Properties {
		if p != nil && p.Name != nil && strings.EqualFold(p.Name.Value, name) {
			if st.Name != nil {
				return st.Name.Value, p
			}
		}
	}
	for _, b := range st.Bases {
		if n := c.structNodes[b.Value]; n != nil {
			if o, p := c.findPropertyDefining(n, name); p != nil {
				return o, p
			}
		}
	}
	return "", nil
}

// findOperatorInHierarchy returns the struct that defines the operator overload and the statement.
func (c *Compiler) findOperatorInHierarchy(st *candy_ast.StructStatement, op string) (owning string, o *candy_ast.OperatorOverloadStatement) {
	if st == nil {
		return "", nil
	}
	for _, ol := range st.Operators {
		if ol != nil && ol.Operator == op {
			if st.Name != nil {
				return st.Name.Value, ol
			}
		}
	}
	for _, b := range st.Bases {
		if n := c.structNodes[b.Value]; n != nil {
			if on, oL := c.findOperatorInHierarchy(n, op); oL != nil {
				return on, oL
			}
		}
	}
	return "", nil
}

func (c *Compiler) structOpNameSuffix(op string) string {
	switch op {
	case "+":
		return "add"
	case "-":
		return "sub"
	case "*":
		return "mul"
	case "/":
		return "div"
	default:
		return ""
	}
}

// bitcastStructPtr rewrites a struct pointer to another struct's pointer type (for inheritance subtyping).
func (c *Compiler) bitcastStructPtr(v value, targetStructName string) value {
	want := "%" + targetStructName + "*"
	if v.typ == want {
		return v
	}
	if !strings.HasPrefix(v.typ, "%") || !strings.HasSuffix(v.typ, "*") {
		return v
	}
	r := c.nextReg()
	c.emit("%s = bitcast %s %s to %s", r, v.typ, v.reg, want)
	return value{reg: r, typ: want}
}

// tryCompileStructInfix emits a call to @DefiningStruct_op_* when the left-hand type has (or inherits) a matching operator.
func (c *Compiler) tryCompileStructInfix(t *candy_ast.InfixExpression, left, right value) (value, bool) {
	if t == nil {
		return value{}, false
	}
	suff := c.structOpNameSuffix(t.Operator)
	if suff == "" {
		return value{}, false
	}
	lst, ok := llvmStructPtrName(left.typ)
	if !ok {
		return value{}, false
	}
	st := c.structNodes[lst]
	if st == nil {
		return value{}, false
	}
	own, opStmt := c.findOperatorInHierarchy(st, t.Operator)
	if opStmt == nil || own == "" {
		return value{}, false
	}
	retName := "int"
	if opStmt.ReturnType != nil {
		retName = candy_ast.ExprAsSimpleTypeName(opStmt.ReturnType)
	}
	retLL := c.mapCandyTypeToLlvm(retName)
	recv := c.bitcastStructPtr(left, own)
	args := []string{fmt.Sprintf("%%%s* %s", own, recv.reg)}
	if len(opStmt.Parameters) != 1 {
		return value{}, false
	}
	p0 := opStmt.Parameters[0]
	if candy_ast.ExprAsSimpleTypeName(p0.TypeName) == "" {
		return value{}, false
	}
	wantLL := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p0.TypeName))
	var argR value
	switch {
	case wantLL == right.typ:
		argR = right
	case strings.HasPrefix(wantLL, "%") && strings.HasSuffix(wantLL, "*"):
		wn := strings.TrimSuffix(strings.TrimPrefix(wantLL, "%"), "*")
		argR = c.bitcastStructPtr(right, wn)
	case wantLL == "i64" && right.typ == "i64", wantLL == "double" && right.typ == "double":
		argR = right
	default:
		return value{}, false
	}
	args = append(args, fmt.Sprintf("%s %s", wantLL, argR.reg))
	fn := own + "_op_" + suff
	out := c.nextReg()
	c.emit("%s = call %s @%s(%s)", out, retLL, fn, strings.Join(args, ", "))
	return value{reg: out, typ: retLL}, true
}

// llvmStructPtrName parses "%Point*" into "Point".
func llvmStructPtrName(lt string) (string, bool) {
	lt = strings.TrimSpace(lt)
	if !strings.HasPrefix(lt, "%") || !strings.HasSuffix(lt, "*") {
		return "", false
	}
	inner := strings.TrimSuffix(lt, "*")
	inner = strings.TrimPrefix(inner, "%")
	return inner, true
}
