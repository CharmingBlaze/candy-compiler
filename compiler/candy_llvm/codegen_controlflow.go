package candy_llvm

import (
	"candy/candy_ast"
	"fmt"
)

func (c *Compiler) freshLabel(prefix string) string {
	c.strCount++
	return fmt.Sprintf("%s%d", prefix, c.strCount)
}

func (c *Compiler) compileBlock(b *candy_ast.BlockStatement) {
	if b == nil {
		return
	}
	for _, s := range b.Statements {
		c.compileStatement(s)
	}
}

// coerceForStore coerces a compiled RHS to match an alloca slot type.
func (c *Compiler) coerceForStore(rhs value, ltyp string) value {
	if rhs.typ == ltyp {
		return rhs
	}
	if rhs.typ == "%any" {
		return c.unboxValueFromAny(rhs, ltyp)
	}
	if ltyp == "i64" {
		return c.valueToI64(rhs)
	}
	if ltyp == "double" {
		if rhs.typ == "i64" {
			r := c.nextReg()
			c.emit("%s = sitofp i64 %s to double", r, rhs.reg)
			return value{reg: r, typ: "double"}
		}
		if rhs.typ == "i1" {
			zz := c.nextReg()
			c.emit("%s = zext i1 %s to i64", zz, rhs.reg)
			r2 := c.nextReg()
			c.emit("%s = sitofp i64 %s to double", r2, zz)
			return value{reg: r2, typ: "double"}
		}
		return c.valueToI64(rhs) // i8* etc. — likely wrong, caller will fail
	}
	if ltyp == "i1" {
		return c.valueAsCond(rhs)
	}
	return rhs
}

// valueToI64 widens a scalar value to i64 for control flow and switch subjects.
func (c *Compiler) valueToI64(v value) value {
	switch v.typ {
	case "i64":
		return v
	case "i1":
		r := c.nextReg()
		c.emit("%s = zext i1 %s to i64", r, v.reg)
		return value{reg: r, typ: "i64"}
	case "double":
		r := c.nextReg()
		c.emit("%s = fptosi double %s to i64", r, v.reg)
		return value{reg: r, typ: "i64"}
	default:
		c.addErr(fmt.Errorf("Native codegen: need scalar for i64 conversion, got %s", v.typ))
		return value{reg: "0", typ: "i64"}
	}
}

// valueAsCond converts a value to i1 for branching.
func (c *Compiler) valueAsCond(v value) value {
	if v.typ == "i1" {
		return v
	}
	if v.typ == "i64" {
		r := c.nextReg()
		c.emit("%s = icmp ne i64 %s, 0", r, v.reg)
		return value{reg: r, typ: "i1"}
	}
	if v.typ == "double" {
		r := c.nextReg()
		c.emit("%s = fcmp one double %s, 0.0", r, v.reg)
		return value{reg: r, typ: "i1"}
	}
	if v.typ == "i8*" {
		r := c.nextReg()
		c.emit("%s = icmp ne i8* %s, null", r, v.reg)
		return value{reg: r, typ: "i1"}
	}
	c.addErr(fmt.Errorf("Native codegen: value type %q cannot be used as condition", v.typ))
	zero := c.nextReg()
	c.emit("%s = icmp eq i1 0, 0", zero)
	return value{reg: zero, typ: "i1"}
}

func (c *Compiler) compileWhile(t *candy_ast.WhileStatement) {
	if t == nil || t.Body == nil {
		return
	}
	header := c.freshLabel("wh")
	body := c.freshLabel("wb")
	merge := c.freshLabel("wm")

	c.emit("br label %%%s", header)
	c.out.WriteString(header + ":\n")
	condV := c.valueAsCond(c.compileExpression(t.Condition))
	c.emit("br i1 %s, label %%%s, label %%%s", condV.reg, body, merge)
	c.out.WriteString(body + ":\n")
	c.compileBlock(t.Body)
	c.emit("br label %%%s", header)
	c.out.WriteString(merge + ":\n")
}

func (c *Compiler) compileDoWhile(t *candy_ast.DoWhileStatement) {
	if t == nil || t.Body == nil {
		return
	}
	body := c.freshLabel("dwb")
	condL := c.freshLabel("dwc")
	merge := c.freshLabel("dwm")

	c.emit("br label %%%s", body)
	c.out.WriteString(body + ":\n")
	c.compileBlock(t.Body)
	c.emit("br label %%%s", condL)
	c.out.WriteString(condL + ":\n")
	condV := c.valueAsCond(c.compileExpression(t.Condition))
	c.emit("br i1 %s, label %%%s, label %%%s", condV.reg, body, merge)
	c.out.WriteString(merge + ":\n")
}

func (c *Compiler) compileCFor(t *candy_ast.CForStatement) {
	if t == nil || t.Body == nil {
		return
	}
	if t.Init != nil {
		c.compileStatement(t.Init)
	}
	header := c.freshLabel("cfh")
	body := c.freshLabel("cfb")
	merge := c.freshLabel("cfm")

	c.emit("br label %%%s", header)
	c.out.WriteString(header + ":\n")
	if t.Cond == nil {
		c.emit("br label %%%s", body)
	} else {
		cv := c.valueAsCond(c.compileExpression(t.Cond))
		c.emit("br i1 %s, label %%%s, label %%%s", cv.reg, body, merge)
	}
	c.out.WriteString(body + ":\n")
	c.compileBlock(t.Body)
	if t.Post != nil {
		_ = c.compileExpression(t.Post)
	}
	c.emit("br label %%%s", header)
	c.out.WriteString(merge + ":\n")
}

func (c *Compiler) compileForBASIC(t *candy_ast.ForStatement) {
	if t == nil {
		return
	}
	if t.Iterable != nil {
		c.addErr(fmt.Errorf("Native codegen: for-in is not supported yet"))
		return
	}
	if t.Body == nil || t.Var == nil || t.Start == nil || t.End == nil {
		return
	}
	if t.Step != nil {
		if il, ok := t.Step.(*candy_ast.IntegerLiteral); ok && il.Value == 0 {
			c.addErr(fmt.Errorf("Native codegen: FOR step cannot be 0"))
			return
		}
	}
	varName := t.Var.Value
	startReg := c.valueToI64(c.compileExpression(t.Start)).reg
	endReg := c.valueToI64(c.compileExpression(t.End)).reg
	var stepReg string
	if t.Step == nil {
		stepReg = "1"
	} else {
		stepReg = c.valueToI64(c.compileExpression(t.Step)).reg
	}

	ptr := "%" + varName
	c.emit("%s = alloca i64", ptr)
	c.emit("store i64 %s, i64* %s", startReg, ptr)
	c.vars[varName] = ptr + "|i64"

	// end and step in registers (constants may be plain immediates; LLVM accepts "store i64 1")

	header := c.freshLabel("fbh")
	body := c.freshLabel("fbb")
	inc := c.freshLabel("fbi")
	merge := c.freshLabel("fbm")

	c.emit("br label %%%s", header)
	c.out.WriteString(header + ":\n")
	iVal := c.nextReg()
	c.emit("%s = load i64, i64* %s", iVal, ptr)
	// (step>0 && i<=end) || (step<0 && i>=end)
	spos := c.nextReg()
	c.emit("%s = icmp sgt i64 %s, 0", spos, stepReg)
	sneg := c.nextReg()
	c.emit("%s = icmp slt i64 %s, 0", sneg, stepReg)
	le := c.nextReg()
	c.emit("%s = icmp sle i64 %s, %s", le, iVal, endReg)
	ge := c.nextReg()
	c.emit("%s = icmp sge i64 %s, %s", ge, iVal, endReg)
	c1 := c.nextReg()
	c.emit("%s = and i1 %s, %s", c1, spos, le)
	c2 := c.nextReg()
	c.emit("%s = and i1 %s, %s", c2, sneg, ge)
	cont := c.nextReg()
	c.emit("%s = or i1 %s, %s", cont, c1, c2)
	c.emit("br i1 %s, label %%%s, label %%%s", cont, body, merge)
	c.out.WriteString(body + ":\n")
	c.compileBlock(t.Body)
	c.emit("br label %%%s", inc)
	c.out.WriteString(inc + ":\n")
	i2 := c.nextReg()
	c.emit("%s = load i64, i64* %s", i2, ptr)
	i3 := c.nextReg()
	c.emit("%s = add i64 %s, %s", i3, i2, stepReg)
	c.emit("store i64 %s, i64* %s", i3, ptr)
	c.emit("br label %%%s", header)
	c.out.WriteString(merge + ":\n")
}

func (c *Compiler) switchPatternI64(p candy_ast.Expression) (value, bool) {
	if p == nil {
		return value{}, false
	}
	switch e := p.(type) {
	case *candy_ast.IntegerLiteral:
		return value{reg: fmt.Sprintf("%d", e.Value), typ: "i64"}, true
	case *candy_ast.Boolean:
		v := "0"
		if e.Value {
			v = "1"
		}
		return value{reg: v, typ: "i64"}, true
	case *candy_ast.Identifier:
		return c.valueToI64(c.compileExpression(e)), true
	default:
		c.addErr(fmt.Errorf("Native codegen: switch case pattern (type %T) not supported yet", p))
		return value{}, false
	}
}

func (c *Compiler) compileSwitch(sw *candy_ast.SwitchStatement) {
	if sw == nil {
		return
	}
	sv := c.compileExpression(sw.Subject)
	if sv.typ == "i8*" {
		c.addErr(fmt.Errorf("Native codegen: string switch is not supported yet"))
	}
	subj := c.valueToI64(sv)

	var defaults []candy_ast.SwitchCase
	var cases []candy_ast.SwitchCase
	for i := range sw.Cases {
		ce := &sw.Cases[i]
		if ce.IsDefault {
			defaults = append(defaults, *ce)
		} else {
			cases = append(cases, *ce)
		}
	}
	merge := c.freshLabel("swm")
	n := len(cases)
	defLabel := c.freshLabel("swd")

	// no non-default cases: optional default, then merge
	if n == 0 {
		if len(defaults) == 0 {
			c.out.WriteString(merge + ":\n")
			return
		}
		c.emit("br label %%%s", defLabel)
		c.out.WriteString(defLabel + ":\n")
		bd := &defaults[len(defaults)-1]
		if bd.Body != nil {
			if b, ok := bd.Body.(*candy_ast.BlockStatement); ok {
				c.compileBlock(b)
			} else {
				c.compileStatement(bd.Body)
			}
		}
		c.emit("br label %%%s", merge)
		c.out.WriteString(merge + ":\n")
		return
	}

	labels := make([]string, n)
	for i := 0; i < n; i++ {
		labels[i] = c.freshLabel("swc")
	}
	testLabels := make([]string, n)
	for i := 0; i < n; i++ {
		testLabels[i] = c.freshLabel("swt")
	}
	c.emit("br label %%%s", testLabels[0])

	for i := 0; i < n; i++ {
		c.out.WriteString(testLabels[i] + ":\n")
		var lastOr string
		for j, p := range cases[i].Patterns {
			pv, ok := c.switchPatternI64(p)
			eq := c.nextReg()
			if ok {
				c.emit("%s = icmp eq i64 %s, %s", eq, subj.reg, pv.reg)
			} else {
				c.emit("%s = icmp eq i1 0, 1", eq) // false
			}
			if j == 0 {
				lastOr = eq
			} else {
				orr := c.nextReg()
				c.emit("%s = or i1 %s, %s", orr, lastOr, eq)
				lastOr = orr
			}
		}
		var nextT string
		if i+1 < n {
			nextT = testLabels[i+1]
		} else if len(defaults) > 0 {
			nextT = defLabel
		} else {
			nextT = merge
		}
		if len(cases[i].Patterns) == 0 {
			c.emit("br label %%%s", nextT)
		} else {
			c.emit("br i1 %s, label %%%s, label %%%s", lastOr, labels[i], nextT)
		}
	}

	for i := 0; i < n; i++ {
		c.out.WriteString(labels[i] + ":\n")
		if cases[i].Body != nil {
			if b, ok := cases[i].Body.(*candy_ast.BlockStatement); ok {
				c.compileBlock(b)
			} else {
				c.compileStatement(cases[i].Body)
			}
		}
		c.emit("br label %%%s", merge)
	}

	if len(defaults) > 0 {
		c.out.WriteString(defLabel + ":\n")
		bd := &defaults[len(defaults)-1]
		if bd.Body != nil {
			if b, ok := bd.Body.(*candy_ast.BlockStatement); ok {
				c.compileBlock(b)
			} else {
				c.compileStatement(bd.Body)
			}
		}
		c.emit("br label %%%s", merge)
	}

	c.out.WriteString(merge + ":\n")
}

func (c *Compiler) compileTernary(t *candy_ast.TernaryExpression) value {
	cond := c.valueAsCond(c.compileExpression(t.Condition))
	id := c.freshLabel("tern")
	thenL := id + "t"
	elseL := id + "e"
	mergeL := id + "m"

	c.emit("br i1 %s, label %%%s, label %%%s", cond.reg, thenL, elseL)

	c.out.WriteString(thenL + ":\n")
	v1 := c.compileExpression(t.Consequence)
	c.emit("br label %%%s", mergeL)

	c.out.WriteString(elseL + ":\n")
	v2 := c.compileExpression(t.Alternative)
	c.emit("br label %%%s", mergeL)

	c.out.WriteString(mergeL + ":\n")
	// Use PHI if types match
	if v1.typ == v2.typ {
		resReg := c.nextReg()
		c.emit("%s = phi %s [ %s, %%%s ], [ %s, %%%s ]", resReg, v1.typ, v1.reg, thenL, v2.reg, elseL)
		return value{reg: resReg, typ: v1.typ}
	}
	// Fallback or bitcast to i64?
	return v1
}
