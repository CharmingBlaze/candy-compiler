package candy_llvm

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

func (c *Compiler) compileCall(t *candy_ast.CallExpression) value {
	if dot, ok := t.Function.(*candy_ast.DotExpression); ok {
		// Method call: v.method(args) -> v_method(v, args)
		left := c.compileExpression(dot.Left)
		if left.typ == "%any" {
			namePtr := c.emitCString(dot.Right.Value)
			dataPtr := c.nextReg()
			c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 1", dataPtr, left.reg)
			raw := c.nextReg()
			c.emit("%s = load i8*, i8** %s", raw, dataPtr)
			out := c.nextReg()
			c.emit("%s = call i64 @candy_dyn_call_i64(i8* %s, i8* %s)", out, raw, namePtr)
			return value{reg: out, typ: "i64"}
		}
		stName := strings.TrimSuffix(strings.TrimPrefix(left.typ, "%"), "*")
		mangled := stName + "_" + dot.Right.Value

		args := []value{left}
		for _, a := range t.Arguments {
			args = append(args, c.compileExpression(a))
		}
		return c.emitCall(mangled, args)
	}

	fnName := ""
	if id, ok := t.Function.(*candy_ast.Identifier); ok {
		fnName = id.Value
	}

	args := []value{}
	for _, a := range t.Arguments {
		args = append(args, c.compileExpression(a))
	}

	// Built-ins
	switch fnName {
	case "print", "println":
		for _, a := range args {
			fmtStr := "@.str.fmt_int"
			switch a.typ {
			case "double":
				fmtStr = "@.str.fmt_float"
			case "i8*":
				fmtStr = "@.str.fmt_str"
			}
			reg := c.nextReg()
			c.emit("%s = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* %s, i64 0, i64 0), %s %s)", reg, fmtStr, a.typ, a.reg)
		}
		return value{reg: "0", typ: "i32"}
	case "exit":
		if len(args) > 0 {
			c.emit("call void @exit(i32 %s)", args[0].reg)
		} else {
			c.emit("call void @exit(i32 0)")
		}
		return value{reg: "0", typ: "void"}
	case "clock":
		reg := c.nextReg()
		c.emit("%s = call i64 @clock()", reg)
		return value{reg: reg, typ: "i64"}
	case "sqrt":
		// Fast typed path: map sqrt(x) directly to LLVM intrinsic.
		if len(args) == 1 {
			a := args[0]
			if a.typ == "i64" {
				conv := c.nextReg()
				c.emit("%s = sitofp i64 %s to double", conv, a.reg)
				a = value{reg: conv, typ: "double"}
			}
			if a.typ == "double" {
				reg := c.nextReg()
				c.emit("%s = call double @llvm.sqrt.f64(double %s)", reg, a.reg)
				return value{reg: reg, typ: "double"}
			}
		}
	case "abs":
		// Typed fast path: avoid dynamic dispatch/boxing.
		if len(args) == 1 {
			a := args[0]
			switch a.typ {
			case "i64":
				neg := c.nextReg()
				c.emit("%s = sub i64 0, %s", neg, a.reg)
				cmp := c.nextReg()
				c.emit("%s = icmp slt i64 %s, 0", cmp, a.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, i64 %s, i64 %s", out, cmp, neg, a.reg)
				return value{reg: out, typ: "i64"}
			case "double":
				neg := c.nextReg()
				c.emit("%s = fsub double 0.0, %s", neg, a.reg)
				cmp := c.nextReg()
				c.emit("%s = fcmp olt double %s, 0.0", cmp, a.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, double %s, double %s", out, cmp, neg, a.reg)
				return value{reg: out, typ: "double"}
			}
		}
	case "min":
		if len(args) == 2 {
			a, b := args[0], args[1]
			// Promote mixed numeric types to double.
			if (a.typ == "double" || b.typ == "double") && (a.typ == "i64" || a.typ == "double") && (b.typ == "i64" || b.typ == "double") {
				if a.typ == "i64" {
					conv := c.nextReg()
					c.emit("%s = sitofp i64 %s to double", conv, a.reg)
					a = value{reg: conv, typ: "double"}
				}
				if b.typ == "i64" {
					conv := c.nextReg()
					c.emit("%s = sitofp i64 %s to double", conv, b.reg)
					b = value{reg: conv, typ: "double"}
				}
			}
			if a.typ == "i64" && b.typ == "i64" {
				cmp := c.nextReg()
				c.emit("%s = icmp slt i64 %s, %s", cmp, a.reg, b.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, i64 %s, i64 %s", out, cmp, a.reg, b.reg)
				return value{reg: out, typ: "i64"}
			}
			if a.typ == "double" && b.typ == "double" {
				cmp := c.nextReg()
				c.emit("%s = fcmp olt double %s, %s", cmp, a.reg, b.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, double %s, double %s", out, cmp, a.reg, b.reg)
				return value{reg: out, typ: "double"}
			}
		}
	case "max":
		if len(args) == 2 {
			a, b := args[0], args[1]
			// Promote mixed numeric types to double.
			if (a.typ == "double" || b.typ == "double") && (a.typ == "i64" || a.typ == "double") && (b.typ == "i64" || b.typ == "double") {
				if a.typ == "i64" {
					conv := c.nextReg()
					c.emit("%s = sitofp i64 %s to double", conv, a.reg)
					a = value{reg: conv, typ: "double"}
				}
				if b.typ == "i64" {
					conv := c.nextReg()
					c.emit("%s = sitofp i64 %s to double", conv, b.reg)
					b = value{reg: conv, typ: "double"}
				}
			}
			if a.typ == "i64" && b.typ == "i64" {
				cmp := c.nextReg()
				c.emit("%s = icmp sgt i64 %s, %s", cmp, a.reg, b.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, i64 %s, i64 %s", out, cmp, a.reg, b.reg)
				return value{reg: out, typ: "i64"}
			}
			if a.typ == "double" && b.typ == "double" {
				cmp := c.nextReg()
				c.emit("%s = fcmp ogt double %s, %s", cmp, a.reg, b.reg)
				out := c.nextReg()
				c.emit("%s = select i1 %s, double %s, double %s", out, cmp, a.reg, b.reg)
				return value{reg: out, typ: "double"}
			}
		}
	}

	return c.emitCall(fnName, args)
}

func (c *Compiler) emitCall(name string, args []value) value {
	argStr := []string{}
	for _, a := range args {
		argStr = append(argStr, fmt.Sprintf("%s %s", a.typ, a.reg))
	}
	reg := c.nextReg()
	retTy := "i64"
	if rt, ok := c.funcRet[name]; ok && rt != "" {
		retTy = rt
	}
	c.emit("%s = call %s @%s(%s)", reg, retTy, name, strings.Join(argStr, ", "))
	return value{reg: reg, typ: retTy}
}
