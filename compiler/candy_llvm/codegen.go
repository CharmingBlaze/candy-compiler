package candy_llvm

import (
	"candy/candy_ast"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Compiler struct {
	out           *strings.Builder
	globals       *strings.Builder
	regCount      int
	strCount      int
	vars          map[string]string              // name -> alloca pointer register|type
	structs       map[string]map[string]int      // name -> field name -> index
	structTy      map[string]string              // name -> LLVM type string e.g. "{ i64, double }"
	structProps   map[string]map[string]string   // name -> prop name -> type
	structOps     map[string]map[string]string   // name -> op -> return type
	structSetters map[string]map[string]struct{} // name -> property name (has user-defined setter) -> mark
	curRet        string                         // non-empty: LLVM return type for the function being emitted ("" = @main)
	typeTags      map[string]int                 // name -> unique integer tag
	tagCount      int
	narrowed      []map[string]string                   // flow-sensitive narrowing scopes: var -> struct type
	structNodes   map[string]*candy_ast.StructStatement // name -> node
	funcRet       map[string]string                     // function name -> LLVM return type
	errs          []error
}

type value struct {
	reg string
	typ string // "i64", "double", "i1", "i8*"
}

func New() *Compiler {
	c := &Compiler{
		out:           &strings.Builder{},
		globals:       &strings.Builder{},
		vars:          make(map[string]string),
		structs:       make(map[string]map[string]int),
		structTy:      make(map[string]string),
		structProps:   make(map[string]map[string]string),
		structOps:     make(map[string]map[string]string),
		structSetters: make(map[string]map[string]struct{}),
		structNodes:   make(map[string]*candy_ast.StructStatement),
		funcRet:       make(map[string]string),
		typeTags: map[string]int{
			"null":   0,
			"int":    1,
			"float":  2,
			"string": 3,
			"bool":   4,
		},
		tagCount: 10, // Start struct tags from 10
	}
	// Define universal Any type
	c.globals.WriteString("%any = type { i8, i8* }\n")
	c.globals.WriteString("declare i64 @candy_dyn_get_i64(i8*, i8*)\n")
	c.globals.WriteString("declare i64 @candy_dyn_call_i64(i8*, i8*)\n")

	// Declare C functions
	c.globals.WriteString("declare i32 @printf(i8*, ...)\n")
	c.globals.WriteString("declare void @exit(i32)\n")
	c.globals.WriteString("declare i64 @strlen(i8*)\n")
	c.globals.WriteString("declare i64 @clock()\n")
	c.globals.WriteString("declare double @llvm.sqrt.f64(double)\n")
	c.globals.WriteString("declare i8* @malloc(i64)\n")
	c.globals.WriteString("declare i8* @strcpy(i8*, i8*)\n")
	c.globals.WriteString("declare i8* @strcat(i8*, i8*)\n")
	c.globals.WriteString("declare i32 @sprintf(i8*, i8*, ...)\n")

	// Runtime helper: string addition
	c.globals.WriteString(`
define i8* @candy_str_add(i8* %s1, i8* %s2) {
  %l1 = call i64 @strlen(i8* %s1)
  %l2 = call i64 @strlen(i8* %s2)
  %l = add i64 %l1, %l2
  %ltot = add i64 %l, 1
  %ptr = call i8* @malloc(i64 %ltot)
  %tmp1 = call i8* @strcpy(i8* %ptr, i8* %s1)
  %tmp2 = call i8* @strcat(i8* %ptr, i8* %s2)
  ret i8* %ptr
}

define i8* @candy_int_to_str(i64 %val) {
  %ptr = call i8* @malloc(i64 32)
  %fmt = getelementptr inbounds [5 x i8], [5 x i8]* @.str.fmt_int_only, i64 0, i64 0
  %tmp = call i32 (i8*, i8*, ...) @sprintf(i8* %ptr, i8* %fmt, i64 %val)
  ret i8* %ptr
}
`)

	// Common format strings
	c.addGlobalStr("fmt_int", "%lld\\0A\\00")
	c.addGlobalStr("fmt_int_only", "%lld\\00")
	c.addGlobalStr("fmt_float", "%f\\0A\\00")
	c.addGlobalStr("fmt_str", "%s\\0A\\00")
	return c
}

func (c *Compiler) addErr(err error) {
	if err == nil {
		return
	}
	c.errs = append(c.errs, err)
}

func stmtTypeName(s candy_ast.Statement) string {
	if s == nil {
		return "<nil>"
	}
	return reflect.TypeOf(s).String()
}

func (c *Compiler) addGlobalStr(name, content string) string {
	id := fmt.Sprintf("@.str.%s", name)
	// Calculate length for LLVM array type [N x i8]
	// Note: We need to be careful with escaping in the content
	rawContent := strings.ReplaceAll(content, "\\00", "")
	rawContent = strings.ReplaceAll(rawContent, "\\0A", "\n")
	l := len(rawContent) + 1 // +1 for null terminator
	c.globals.WriteString(fmt.Sprintf("%s = private unnamed_addr constant [%d x i8] c\"%s\"\n", id, l, content))
	return id
}

func (c *Compiler) nextReg() string {
	c.regCount++
	return fmt.Sprintf("%%%d", c.regCount)
}

func (c *Compiler) emit(format string, a ...any) {
	fmt.Fprintf(c.out, "  "+format+"\n", a...)
}

func (c *Compiler) GenerateIR(program *candy_ast.Program) (string, error) {
	body := &strings.Builder{}
	body.WriteString("define i64 @main() " + c.functionAttributes("main") + " {\n")

	// Temporary swap to capture body
	originalOut := c.out
	c.out = body
	// Pass 1: Collect struct/class definitions
	for _, s := range program.Statements {
		if st, ok := s.(*candy_ast.StructStatement); ok {
			c.structNodes[st.Name.Value] = st
		}
	}
	// Pass 2: Pre-collect function signatures so calls can use correct return types.
	c.collectFunctionSignatures(program)

	for _, stmt := range program.Statements {
		c.compileStatement(stmt)
	}
	if !strings.Contains(c.out.String(), "ret ") {
		c.emit("ret i64 0")
	}
	c.out.WriteString("}\n")

	finalBody := c.out.String()
	c.out = originalOut

	ir := c.globals.String() + "\n" + finalBody
	if len(c.errs) > 0 {
		return "", errors.Join(c.errs...)
	}
	return ir, nil
}

func (c *Compiler) collectFunctionSignatures(program *candy_ast.Program) {
	for _, stmt := range program.Statements {
		switch t := stmt.(type) {
		case *candy_ast.FunctionStatement:
			c.funcRet[c.functionLLVMName(t)] = c.functionReturnType(t)
		case *candy_ast.ExternFunctionStatement:
			if t != nil && t.Function != nil {
				c.funcRet[t.Function.Name.Value] = c.functionReturnType(t.Function)
			}
		case *candy_ast.StructStatement:
			for _, m := range t.Methods {
				c.funcRet[c.functionLLVMName(m)] = c.functionReturnType(m)
			}
		}
	}
}

func (c *Compiler) functionLLVMName(t *candy_ast.FunctionStatement) string {
	name := t.Name.Value
	if t.Receiver != nil {
		name = candy_ast.ExprAsSimpleTypeName(t.Receiver.TypeName) + "_" + name
	}
	return name
}

func (c *Compiler) functionReturnType(t *candy_ast.FunctionStatement) string {
	retTy := "i64"
	if candy_ast.ExprAsSimpleTypeName(t.ReturnType) != "" {
		retTy = c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(t.ReturnType))
	}
	return retTy
}

func (c *Compiler) compileStatement(stmt candy_ast.Statement) {
	if stmt == nil {
		return
	}
	switch t := stmt.(type) {
	case *candy_ast.ExpressionStatement:
		c.compileExpression(t.Expression)
	case *candy_ast.ValStatement:
		c.compileVarDecl(t.Name.Value, t.TypeName, t.Value)
	case *candy_ast.VarStatement:
		c.compileVarDecl(t.Name.Value, t.TypeName, t.Value)
	case *candy_ast.ReturnStatement:
		c.emit("; compiling return statement")
		if c.curRet == "void" {
			if t.ReturnValue != nil {
				_ = c.compileExpression(t.ReturnValue)
			}
			c.emit("ret void")
			return
		}
		if t.ReturnValue == nil {
			c.addErr(fmt.Errorf("codegen: empty return in non-void function"))
			c.emit("ret %s 0", c.curRet)
			return
		}
		val := c.compileExpression(t.ReturnValue)
		if c.curRet != "" {
			rr := val
			if val.typ != c.curRet {
				if strings.HasPrefix(c.curRet, "%") && strings.HasSuffix(c.curRet, "*") {
					if strings.HasPrefix(val.typ, "%") && strings.HasSuffix(val.typ, "*") {
						wn := strings.TrimSuffix(strings.TrimPrefix(c.curRet, "%"), "*")
						rr = c.bitcastStructPtr(val, wn)
					}
				}
			}
			c.emit("ret %s %s", c.curRet, rr.reg)
			return
		}
		// @main: i64 return convention
		switch val.typ {
		case "double":
			reg := c.nextReg()
			c.emit("%s = fptosi double %s to i64", reg, val.reg)
			c.emit("ret i64 %s", reg)
		case "i1":
			reg := c.nextReg()
			c.emit("%s = zext i1 %s to i64", reg, val.reg)
			c.emit("ret i64 %s", reg)
		case "i8*":
			c.emit("ret i64 0") // Cannot return string as int
		default:
			c.emit("ret i64 %s", val.reg)
		}
	case *candy_ast.IfExpression:
		c.compileIf(t)
	case *candy_ast.StructStatement:
		if len(t.TypeParameters) > 0 {
			return
		}
		stName := t.Name.Value
		fields := make(map[string]int)
		var types []string

		var collect func(s *candy_ast.StructStatement)
		collect = func(s *candy_ast.StructStatement) {
			// Prepend bases
			for _, b := range s.Bases {
				if baseNode, ok := c.structNodes[b.Value]; ok {
					collect(baseNode)
				}
			}
			// Add fields
			for _, f := range s.Fields {
				if _, ok := fields[f.Name.Value]; !ok {
					fields[f.Name.Value] = len(types)
					types = append(types, c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(f.TypeName)))
				}
			}
		}
		collect(t)

		llvmTy := "{" + strings.Join(types, ", ") + "}"
		c.structs[stName] = fields
		c.structTy[stName] = llvmTy
		c.typeTags[stName] = c.tagCount
		c.tagCount++

		// Define type globally
		c.globals.WriteString(fmt.Sprintf("%%%s = type %s\n", stName, llvmTy))

		// Compile methods
		for _, m := range t.Methods {
			c.compileFunction(m)
		}
		// Compile properties
		for _, p := range t.Properties {
			c.compileProperty(stName, p)
		}
		// Compile operator overloads
		for _, o := range t.Operators {
			c.compileOperatorOverload(stName, o)
		}
	case *candy_ast.FunctionStatement:
		c.compileFunction(t)
	case *candy_ast.ExternFunctionStatement:
		c.compileExternFunction(t)
	case *candy_ast.BlockStatement:
		c.compileBlock(t)
	case *candy_ast.WhileStatement:
		c.compileWhile(t)
	case *candy_ast.DoWhileStatement:
		c.compileDoWhile(t)
	case *candy_ast.CForStatement:
		c.compileCFor(t)
	case *candy_ast.ForStatement:
		c.compileForBASIC(t)
	case *candy_ast.SwitchStatement:
		c.compileSwitch(t)
	case *candy_ast.ClassStatement, *candy_ast.ObjectStatement, *candy_ast.InterfaceStatement, *candy_ast.TraitStatement, *candy_ast.PackageStatement:
		c.addErr(fmt.Errorf("Native codegen: %s not supported yet", strings.TrimPrefix(stmtTypeName(t), "*candy_ast.")))
	case *candy_ast.ImportStatement:
		// candy_load.ExpandProgramForBuild should strip these before GenerateIR.
		c.addErr(fmt.Errorf("internal: unexpanded import in codegen (use candy -build with a .candy file)"))
	case *candy_ast.ModuleStatement:
		c.addErr(fmt.Errorf("Native codegen: module is not supported yet"))
	case *candy_ast.EnumStatement, *candy_ast.TryStatement, *candy_ast.RunStatement, *candy_ast.DeferStatement, *candy_ast.PropertyStatement, *candy_ast.OperatorOverloadStatement:
		c.addErr(fmt.Errorf("Native codegen: %s not supported yet", strings.TrimPrefix(stmtTypeName(t), "*candy_ast.")))
	case *candy_ast.WithStatement:
		c.addErr(fmt.Errorf("Native codegen: with is not supported yet"))
	case *candy_ast.LibraryStatement:
		if t.Body != nil {
			for _, st := range t.Body.Statements {
				c.compileStatement(st)
			}
		}
	case *candy_ast.LambdaExpression:
		// Can appear as statement in some parse paths; treat as expr stmt
		_ = c.compileExpression(t)
	default:
		c.addErr(fmt.Errorf("Native codegen: unsupported statement %T", t))
	}
}

func (c *Compiler) compileExternFunction(t *candy_ast.ExternFunctionStatement) {
	if t == nil || t.Function == nil || t.Function.Name == nil {
		return
	}
	fn := t.Function
	retTy := "i64"
	if fn.ReturnType != nil {
		retName := strings.ToLower(candy_ast.ExprAsSimpleTypeName(fn.ReturnType))
		if !isSupportedExternCandyType(retName) {
			c.addErr(fmt.Errorf("extern return type not supported in native codegen: %s", retName))
			return
		}
		retTy = c.mapCandyTypeToLlvm(retName)
	}
	c.funcRet[fn.Name.Value] = retTy
	params := make([]string, 0, len(fn.Parameters))
	for _, p := range fn.Parameters {
		pt := "i64"
		if candy_ast.ExprAsSimpleTypeName(p.TypeName) != "" {
			pName := strings.ToLower(candy_ast.ExprAsSimpleTypeName(p.TypeName))
			if !isSupportedExternCandyType(pName) {
				c.addErr(fmt.Errorf("extern param type not supported in native codegen: %s", pName))
				return
			}
			pt = c.mapCandyTypeToLlvm(pName)
		}
		params = append(params, pt)
	}
	c.globals.WriteString(fmt.Sprintf("declare %s @%s(%s)\n", retTy, fn.Name.Value, strings.Join(params, ", ")))
}

func isSupportedExternCandyType(name string) bool {
	switch strings.ToLower(name) {
	case "", "int", "float", "double", "bool", "string", "void":
		return true
	default:
		return false
	}
}

func (c *Compiler) compileFunction(t *candy_ast.FunctionStatement) {
	name := c.functionLLVMName(t)
	retTy := c.functionReturnType(t)
	c.funcRet[name] = retTy

	params := []string{}
	if t.Receiver != nil {
		params = append(params, fmt.Sprintf("%%%s* %%%s_val", candy_ast.ExprAsSimpleTypeName(t.Receiver.TypeName), t.Receiver.Name.Value))
	}
	for _, p := range t.Parameters {
		params = append(params, fmt.Sprintf("%s %%%s_val", c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p.TypeName)), p.Name.Value))
	}

	// Define function globally
	fmt.Fprintf(c.globals, "define %s @%s(%s) %s {\n", retTy, name, strings.Join(params, ", "), c.functionAttributes(name))

	// Temporarily switch output to globals (actually we need a nested builder or just write to globals)
	// But we need to use c.emit() which writes to c.out.
	// So we swap c.out to a temporary builder, then append to globals.

	oldOut := c.out
	funcBody := &strings.Builder{}
	c.out = funcBody

	// Reset local vars for function scope
	oldVars := c.vars
	oldRegCount := c.regCount
	oldCur := c.curRet
	c.curRet = retTy
	c.vars = make(map[string]string)
	c.regCount = 0

	// Alloca parameters
	if t.Receiver != nil {
		pName := t.Receiver.Name.Value
		pTy := "%" + candy_ast.ExprAsSimpleTypeName(t.Receiver.TypeName) + "*"
		ptr := "%" + pName + "_ptr"
		c.emit("%s = alloca %s", ptr, pTy)
		c.emit("store %s %%%s_val, %s* %s", pTy, pName, pTy, ptr)
		c.vars[pName] = ptr + "|" + pTy
	}
	for _, p := range t.Parameters {
		pName := p.Name.Value
		pTy := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p.TypeName))
		ptr := "%" + pName + "_ptr"
		c.emit("%s = alloca %s", ptr, pTy)
		c.emit("store %s %%%s_val, %s* %s", pTy, pName, pTy, ptr)
		c.vars[pName] = ptr + "|" + pTy
	}

	if t.Body != nil {
		c.compileBlock(t.Body)
	}

	/*
		if !strings.Contains(funcBody.String(), "ret ") {
			if retTy == "void" {
				c.emit("ret void")
			} else {
				c.emit("ret %s 0", retTy)
			}
		}
	*/

	c.globals.WriteString(funcBody.String())
	c.globals.WriteString("}\n")

	c.out = oldOut
	c.vars = oldVars
	c.regCount = oldRegCount
	c.curRet = oldCur
}

func (c *Compiler) compileProperty(stName string, p *candy_ast.PropertyStatement) {
	if c.structProps[stName] == nil {
		c.structProps[stName] = make(map[string]string)
	}
	ty := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p.Type))
	c.structProps[stName][p.Name.Value] = ty

	if p.Getter != nil {
		name := stName + "_get_" + p.Name.Value
		fmt.Fprintf(c.globals, "define %s @%s(%%%s* %%this) %s {\n", ty, name, stName, c.functionAttributes(name))
		oldOut := c.out
		funcBody := &strings.Builder{}
		c.out = funcBody

		oldVars := c.vars
		oldRegCount := c.regCount
		oldCur := c.curRet
		c.curRet = ty
		c.vars = make(map[string]string)
		c.regCount = 0

		c.vars["this"] = "%this_ptr|" + "%" + stName + "*"
		c.emit("%%this_ptr = alloca %%%s*", stName)
		c.emit("store %%%s* %%this, %%%s** %%this_ptr", stName, stName)

		c.compileBlock(p.Getter)

		if !strings.Contains(funcBody.String(), "ret ") {
			c.emit("ret %s 0", ty)
		}
		c.globals.WriteString(funcBody.String())
		c.globals.WriteString("}\n")
		c.out = oldOut
		c.vars = oldVars
		c.regCount = oldRegCount
		c.curRet = oldCur
	}
	if p.Setter != nil {
		if c.structSetters[stName] == nil {
			c.structSetters[stName] = make(map[string]struct{})
		}
		c.structSetters[stName][p.Name.Value] = struct{}{}
		name := stName + "_set_" + p.Name.Value
		fmt.Fprintf(c.globals, "define void @%s(%%%s* %%this, %s %%value) %s {\n", name, stName, ty, c.functionAttributes(name))
		oldOut := c.out
		funcBody := &strings.Builder{}
		c.out = funcBody

		oldVars := c.vars
		oldRegCount := c.regCount
		oldCur := c.curRet
		c.curRet = "void"
		c.vars = make(map[string]string)
		c.regCount = 0

		c.vars["this"] = "%this_ptr|" + "%" + stName + "*"
		c.emit("%%this_ptr = alloca %%%s*", stName)
		c.emit("store %%%s* %%this, %%%s** %%this_ptr", stName, stName)

		c.emit("%%value_ptr = alloca %s", ty)
		c.emit("store %s %%value, %s* %%value_ptr", ty, ty)
		c.vars["value"] = "%value_ptr|" + ty

		c.compileBlock(p.Setter)

		if !strings.Contains(funcBody.String(), "ret ") {
			c.emit("ret void")
		}
		c.globals.WriteString(funcBody.String())
		c.globals.WriteString("}\n")
		c.out = oldOut
		c.vars = oldVars
		c.regCount = oldRegCount
		c.curRet = oldCur
	}
}

func (c *Compiler) compileOperatorOverload(stName string, o *candy_ast.OperatorOverloadStatement) {
	if c.structOps[stName] == nil {
		c.structOps[stName] = make(map[string]string)
	}
	retTy := "i64"
	if o.ReturnType != nil {
		retTy = c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(o.ReturnType))
	}
	c.structOps[stName][o.Operator] = retTy

	opName := "add"
	switch o.Operator {
	case "+":
		opName = "add"
	case "-":
		opName = "sub"
	case "*":
		opName = "mul"
	case "/":
		opName = "div"
	}

	name := stName + "_op_" + opName
	params := []string{fmt.Sprintf("%%%s* %%this", stName)}
	for _, p := range o.Parameters {
		params = append(params, fmt.Sprintf("%s %%%s_val", c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p.TypeName)), p.Name.Value))
	}

	fmt.Fprintf(c.globals, "define %s @%s(%s) %s {\n", retTy, name, strings.Join(params, ", "), c.functionAttributes(name))
	oldOut := c.out
	funcBody := &strings.Builder{}
	c.out = funcBody

	oldVars := c.vars
	oldRegCount := c.regCount
	oldCur := c.curRet
	c.curRet = retTy
	c.vars = make(map[string]string)
	c.regCount = 0

	c.vars["this"] = "%this_ptr|" + "%" + stName + "*"
	c.emit("%%this_ptr = alloca %%%s*", stName)
	c.emit("store %%%s* %%this, %%%s** %%this_ptr", stName, stName)

	for _, p := range o.Parameters {
		pName := p.Name.Value
		pTy := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(p.TypeName))
		ptr := "%" + pName + "_ptr"
		c.emit("%s = alloca %s", ptr, pTy)
		c.emit("store %s %%%s_val, %s* %s", pTy, pName, pTy, ptr)
		c.vars[pName] = ptr + "|" + pTy
	}

	c.compileBlock(o.Body)

	if !strings.Contains(funcBody.String(), "ret ") {
		c.emit("ret %s 0", retTy)
	}
	c.globals.WriteString(funcBody.String())
	c.globals.WriteString("}\n")
	c.out = oldOut
	c.vars = oldVars
	c.regCount = oldRegCount
	c.curRet = oldCur
}

func (c *Compiler) mapCandyTypeToLlvm(tyName string) string {
	switch strings.ToLower(tyName) {
	case "int":
		return "i64"
	case "float", "double":
		return "double"
	case "bool":
		return "i1"
	case "string":
		return "i8*"
	case "any", "dyn":
		return "%any"
	default:
		if _, ok := c.structs[tyName]; ok {
			return "%" + tyName + "*"
		}
		return "i64"
	}
}

func (c *Compiler) compileIf(t *candy_ast.IfExpression) {
	cond := c.compileExpression(t.Condition)

	c.strCount++ // Reuse for labels
	id := c.strCount
	thenLabel := fmt.Sprintf("then%d", id)
	elseLabel := fmt.Sprintf("else%d", id)
	mergeLabel := fmt.Sprintf("merge%d", id)

	if t.Alternative == nil {
		c.emit("br i1 %s, label %%%s, label %%%s", cond.reg, thenLabel, mergeLabel)
	} else {
		c.emit("br i1 %s, label %%%s, label %%%s", cond.reg, thenLabel, elseLabel)
	}

	// Then branch
	c.out.WriteString(fmt.Sprintf("%s:\n", thenLabel))
	if isExpr, ident, narrowedType, ok := c.guardNarrowing(t.Condition); ok {
		_ = isExpr
		c.pushNarrowing(map[string]string{ident: narrowedType})
		defer c.popNarrowing()
	}
	c.compileBlock(t.Consequence)
	c.emit("br label %%%s", mergeLabel)

	// Else branch
	if t.Alternative != nil {
		c.out.WriteString(fmt.Sprintf("%s:\n", elseLabel))
		if altBlock, ok := t.Alternative.(*candy_ast.BlockStatement); ok {
			c.compileBlock(altBlock)
		} else {
			c.compileStatement(t.Alternative)
		}
		c.emit("br label %%%s", mergeLabel)
	}

	// Merge
	c.out.WriteString(fmt.Sprintf("%s:\n", mergeLabel))
}

// structNameForStructLiteral matches generic monomorph names (e.g. Box<float> -> Box_float) when a template is present.
func (c *Compiler) structNameForStructLiteral(e candy_ast.Expression) string {
	if te, ok := e.(*candy_ast.TypeExpression); ok {
		if tpl, ok2 := c.structNodes[te.Name.Value]; ok2 && len(tpl.TypeParameters) > 0 {
			return candy_ast.MonomorphStructName(tpl.Name.Value, te.Arguments)
		}
	}
	return candy_ast.ExprAsSimpleTypeName(e)
}

func (c *Compiler) compileVarDecl(name string, typeName candy_ast.Expression, value candy_ast.Expression) {
	ptr := "%" + name
	val := c.compileExpression(value)
	declared := candy_ast.ExprAsSimpleTypeName(typeName)
	// Infer plain struct type from a struct literal: `val v = S { }` is static S*, not %any.
	if declared == "" {
		if sl, ok := value.(*candy_ast.StructLiteral); ok && sl.Name != nil {
			n := c.structNameForStructLiteral(sl.Name)
			if n != "" && c.structs[n] != nil {
				declared = n
			}
		}
	}
	ty := val.typ
	if declared != "" {
		ty = c.mapCandyTypeToLlvm(declared)
	}

	if ty == "%any" {
		c.emit("%s = alloca %%any", ptr)
		boxed := c.boxValueIntoAny(val)
		c.emit("store %%any %s, %%any* %s", boxed.reg, ptr)
		c.vars[name] = ptr + "|%any"
		return
	}

	c.emit("%s = alloca %s", ptr, ty)
	rhs := val.reg
	if val.typ != ty && val.typ != "" && strings.HasSuffix(val.typ, "*") && strings.HasSuffix(ty, "*") {
		reg := c.nextReg()
		c.emit("%s = bitcast %s %s to %s", reg, val.typ, val.reg, ty)
		rhs = reg
	}
	c.emit("store %s %s, %s* %s", ty, rhs, ty, ptr)
	c.vars[name] = ptr + "|" + ty
}

func (c *Compiler) boxValueIntoAny(val value) value {
	if val.typ == "%any" {
		return val
	}

	res := c.nextReg()
	c.emit("%s = alloca %%any", res)

	tagPtr := c.nextReg()
	c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 0", tagPtr, res)

	tag := 0
	switch val.typ {
	case "i64":
		tag = 1
	case "double":
		tag = 2
	case "i8*":
		tag = 3
	case "i1":
		tag = 4
	default:
		if strings.HasSuffix(val.typ, "*") {
			typeName := strings.TrimSuffix(strings.TrimPrefix(val.typ, "%"), "*")
			tag = c.ensureTypeTag(typeName)
		}
	}
	c.emit("store i8 %d, i8* %s", tag, tagPtr)

	dataPtr := c.nextReg()
	c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 1", dataPtr, res)

	raw := c.nextReg()
	switch val.typ {
	case "i64":
		c.emit("%s = inttoptr i64 %s to i8*", raw, val.reg)
	case "double":
		// Bitcast double to i64 first
		bc := c.nextReg()
		c.emit("%s = bitcast double %s to i64", bc, val.reg)
		c.emit("%s = inttoptr i64 %s to i8*", raw, bc)
	case "i1":
		ze := c.nextReg()
		c.emit("%s = zext i1 %s to i64", ze, val.reg)
		c.emit("%s = inttoptr i64 %s to i8*", raw, ze)
	default:
		c.emit("%s = bitcast %s %s to i8*", raw, val.typ, val.reg)
	}
	c.emit("store i8* %s, i8** %s", raw, dataPtr)

	loadRes := c.nextReg()
	c.emit("%s = load %%any, %%any* %s", loadRes, res)
	return value{reg: loadRes, typ: "%any"}
}

func (c *Compiler) coerceToString(v value) value {
	if v.typ == "i8*" {
		return v
	}
	if v.typ == "%any" {
		// Call runtime function to convert any to string
		reg := c.nextReg()
		c.emit("%s = call i8* @candy_any_to_str(%%any %s)", reg, v.reg)
		return value{reg: reg, typ: "i8*"}
	}
	if v.typ == "i64" {
		reg := c.nextReg()
		c.emit("%s = call i8* @candy_int_to_str(i64 %s)", reg, v.reg)
		return value{reg: reg, typ: "i8*"}
	}
	// Fallback to "???" for unknown types
	return c.compileExpression(&candy_ast.StringLiteral{Value: "???"})
}

func (c *Compiler) unboxValueFromAny(v value, targetLLVMTy string) value {
	if v.typ != "%any" {
		return v
	}
	// Note: In a real implementation, we would check the tag.
	// For Level 1, we assume the user knows what they are doing or let it crash.
	dataPtr := c.nextReg()
	c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 1", dataPtr, v.reg)
	raw := c.nextReg()
	c.emit("%s = load i8*, i8** %s", raw, dataPtr)

	res := c.nextReg()
	switch targetLLVMTy {
	case "i64":
		c.emit("%s = ptrtoint i8* %s to i64", res, raw)
	case "double":
		bc := c.nextReg()
		c.emit("%s = ptrtoint i8* %s to i64", bc, raw)
		c.emit("%s = bitcast i64 %s to double", res, bc)
	case "i1":
		bc := c.nextReg()
		c.emit("%s = ptrtoint i8* %s to i64", bc, raw)
		c.emit("%s = trunc i64 %s to i1", res, bc)
	default:
		c.emit("%s = bitcast i8* %s to %s", res, raw, targetLLVMTy)
	}
	return value{reg: res, typ: targetLLVMTy}
}

func (c *Compiler) compileExpression(expr candy_ast.Expression) value {
	switch t := expr.(type) {
	case *candy_ast.IntegerLiteral:
		return value{reg: fmt.Sprintf("%d", t.Value), typ: "i64"}
	case *candy_ast.FloatLiteral:
		return value{reg: fmt.Sprintf("%f", t.Value), typ: "double"}
	case *candy_ast.Boolean:
		bVal := "0"
		if t.Value {
			bVal = "1"
		}
		return value{reg: bVal, typ: "i1"}
	case *candy_ast.StringLiteral:
		c.strCount++
		id := fmt.Sprintf("s%d", c.strCount)
		content := t.Value + "\\00"
		ptr := c.addGlobalStr(id, content)
		// Get pointer to first element: bitcast or getelementptr
		reg := c.nextReg()
		l := len(t.Value) + 1
		c.emit("%s = getelementptr inbounds [%d x i8], [%d x i8]* %s, i64 0, i64 0", reg, l, l, ptr)
		return value{reg: reg, typ: "i8*"}
	case *candy_ast.InterpolatedStringLiteral:
		if len(t.Parts) == 0 {
			return c.compileExpression(&candy_ast.StringLiteral{Value: ""})
		}
		res := c.compileExpression(t.Parts[0])
		for i := 1; i < len(t.Parts); i++ {
			right := c.compileExpression(t.Parts[i])
			res = c.coerceToString(res)
			right = c.coerceToString(right)
			reg := c.nextReg()
			c.emit("%s = call i8* @candy_str_add(i8* %s, i8* %s)", reg, res.reg, right.reg)
			res = value{reg: reg, typ: "i8*"}
		}
		return res
	case *candy_ast.Identifier:
		if data, ok := c.vars[t.Value]; ok {
			parts := strings.Split(data, "|")
			ptr, typ := parts[0], parts[1]
			if typ == "%any" {
				if narrowed, ok := c.lookupNarrowed(t.Value); ok {
					dataPtr := c.nextReg()
					c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 1", dataPtr, ptr)
					raw := c.nextReg()
					c.emit("%s = load i8*, i8** %s", raw, dataPtr)
					casted := c.nextReg()
					target := "%" + narrowed + "*"
					c.emit("%s = bitcast i8* %s to %s", casted, raw, target)
					return value{reg: casted, typ: target}
				}
				// Keep dynamic values boxed outside narrowed scopes.
				return value{reg: ptr, typ: "%any"}
			}
			reg := c.nextReg()
			c.emit("%s = load %s, %s* %s", reg, typ, typ, ptr)
			return value{reg: reg, typ: typ}
		}
		c.addErr(fmt.Errorf("codegen: undefined variable %q", t.Value))
		return value{reg: "0", typ: "i64"}
	case *candy_ast.CallExpression:
		return c.compileCall(t)
	case *candy_ast.InfixExpression:
		return c.compileInfix(t)
	case *candy_ast.StructLiteral:
		stName := c.structNameForStructLiteral(t.Name)
		llvmTy := "%" + stName
		ptr := c.nextReg()
		c.emit("%s = alloca %s", ptr, llvmTy)

		fieldIndices := c.structs[stName]
		for fname, valExpr := range t.Fields {
			idx := fieldIndices[fname]
			val := c.compileExpression(valExpr)

			fptr := c.nextReg()
			c.emit("%s = getelementptr inbounds %s, %s* %s, i32 0, i32 %d", fptr, llvmTy, llvmTy, ptr, idx)
			c.emit("store %s %s, %s* %s", val.typ, val.reg, val.typ, fptr)
		}
		// Load the full struct or just return the pointer?
		// For now, we return the pointer but mark the type as the struct name
		return value{reg: ptr, typ: "%" + stName + "*"}
	case *candy_ast.IsExpression:
		// Type guard: supports dynamic `%any` checks and static pointer checks.
		left := c.compileExpression(t.Left)
		if t.TypeName == nil {
			return value{reg: "0", typ: "i1"}
		}
		if left.typ == "%any" {
			tagPtr := c.nextReg()
			c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 0", tagPtr, left.reg)
			tag := c.nextReg()
			c.emit("%s = load i8, i8* %s", tag, tagPtr)
			want := c.ensureTypeTag(candy_ast.ExprAsSimpleTypeName(t.TypeName))
			ok := c.nextReg()
			c.emit("%s = icmp eq i8 %s, %d", ok, tag, want)
			return value{reg: ok, typ: "i1"}
		}
		wantPtr := "%" + candy_ast.ExprAsSimpleTypeName(t.TypeName) + "*"
		if left.typ == wantPtr {
			return value{reg: "1", typ: "i1"}
		}
		return value{reg: "0", typ: "i1"}
	case *candy_ast.GroupedExpression:
		if t.Expr == nil {
			return value{reg: "0", typ: "i64"}
		}
		return c.compileExpression(t.Expr)
	case *candy_ast.PrefixExpression:
		right := c.compileExpression(t.Right)
		switch t.Operator {
		case "-":
			switch right.typ {
			case "i64":
				r := c.nextReg()
				c.emit("%s = sub i64 0, %s", r, right.reg)
				return value{reg: r, typ: "i64"}
			case "double":
				r := c.nextReg()
				c.emit("%s = fsub double 0.0, %s", r, right.reg)
				return value{reg: r, typ: "double"}
			}
			c.addErr(fmt.Errorf("Native codegen: unary - expects i64 or double, got %s", right.typ))
			return value{reg: "0", typ: "i64"}
		case "!":
			if right.typ == "i1" {
				r := c.nextReg()
				c.emit("%s = xor i1 %s, 1", r, right.reg)
				return value{reg: r, typ: "i1"}
			}
			if right.typ == "i64" {
				r := c.nextReg()
				c.emit("%s = icmp eq i64 %s, 0", r, right.reg)
				return value{reg: r, typ: "i1"}
			}
			c.addErr(fmt.Errorf("Native codegen: unary ! expects i1 or i64, got %s", right.typ))
			return value{reg: "0", typ: "i1"}
		case "~":
			if right.typ == "i64" {
				r := c.nextReg()
				c.emit("%s = xor i64 %s, -1", r, right.reg)
				return value{reg: r, typ: "i64"}
			}
			c.addErr(fmt.Errorf("Native codegen: unary ~ expects i64, got %s", right.typ))
			return value{reg: "0", typ: "i64"}
		case "+":
			if right.typ == "i64" || right.typ == "double" || right.typ == "i1" || right.typ == "i8*" {
				return right
			}
			c.addErr(fmt.Errorf("Native codegen: unary + on %s", right.typ))
			return value{reg: "0", typ: "i64"}
		default:
			c.addErr(fmt.Errorf("Native codegen: unknown unary operator %q", t.Operator))
			return value{reg: "0", typ: "i64"}
		}
	case *candy_ast.PostfixExpression:
		if t.Operator == "++" || t.Operator == "--" {
			c.addErr(fmt.Errorf("Native codegen: postfix ++/-- (use interpreter or extend codegen)"))
			return value{reg: "0", typ: "i64"}
		}
		c.addErr(fmt.Errorf("Native codegen: unknown postfix %q", t.Operator))
		return value{reg: "0", typ: "i64"}
	case *candy_ast.IndexExpression:
		// String as i8* + index i64: single-byte load, zero-extended to i64.
		if t.Base == nil || t.Index == nil {
			return value{reg: "0", typ: "i64"}
		}
		base := c.compileExpression(t.Base)
		if base.typ != "i8*" {
			c.addErr(fmt.Errorf("Native codegen: index is only supported for string (i8*) for now, got %s", base.typ))
			return value{reg: "0", typ: "i64"}
		}
		idxv := c.valueToI64(c.compileExpression(t.Index))
		gep := c.nextReg()
		c.emit("%s = getelementptr inbounds i8, i8* %s, i64 %s", gep, base.reg, idxv.reg)
		b := c.nextReg()
		c.emit("%s = load i8, i8* %s", b, gep)
		ze := c.nextReg()
		c.emit("%s = zext i8 %s to i64", ze, b)
		return value{reg: ze, typ: "i64"}
	case *candy_ast.AssignExpression:
		if t == nil || t.Left == nil || t.Value == nil {
			return value{reg: "0", typ: "i64"}
		}
		if id, ok := t.Left.(*candy_ast.Identifier); ok {
			data, found := c.vars[id.Value]
			if !found {
				// Implicit declaration!
				c.compileVarDecl(id.Value, nil, t.Value)
				return c.compileExpression(id)
			}
			parts := strings.Split(data, "|")
			ptr, ltyp := parts[0], parts[1]
			rhs := c.compileExpression(t.Value)
			if ltyp == "%any" {
				boxed := c.boxValueIntoAny(rhs)
				c.emit("store %%any %s, %%any* %s", boxed.reg, ptr)
				return boxed
			}
			rr := c.coerceForStore(rhs, ltyp)
			c.emit("store %s %s, %s* %s", ltyp, rr.reg, ltyp, ptr)
			return rr
		}
		if dot, ok := t.Left.(*candy_ast.DotExpression); ok {
			left := c.compileExpression(dot.Left)
			if strings.HasSuffix(left.typ, "*") {
				lst, ok0 := llvmStructPtrName(left.typ)
				if !ok0 {
					return value{reg: "0", typ: "i64"}
				}
				st0 := c.structNodes[lst]
				if st0 == nil {
					return value{reg: "0", typ: "i64"}
				}
				own, prop := c.findPropertyDefining(st0, dot.Right.Value)
				if prop != nil {
					if prop.Setter == nil {
						c.addErr(fmt.Errorf("codegen: property %s has no setter for assignment", dot.Right.Value))
						return value{reg: "0", typ: "i64"}
					}
					rhs := c.compileExpression(t.Value)
					recv := c.bitcastStructPtr(left, own)
					pTy := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(prop.Type))
					rr := c.coerceForStore(rhs, pTy)
					c.emit("call void @%s_set_%s(%%%s* %s, %s %s)", own, dot.Right.Value, own, recv.reg, pTy, rr.reg)
					return rhs
				}
				// struct field
				rhs := c.compileExpression(t.Value)
				stName := lst
				if idx, ok1 := c.structs[stName][dot.Right.Value]; ok1 {
					fTy := "i64"
					tyStr := c.structTy[stName]
					tyStr = strings.Trim(tyStr, "{}")
					parts := strings.Split(tyStr, ",")
					if idx < len(parts) {
						fTy = strings.TrimSpace(parts[idx])
					}
					rhs = c.coerceForStore(rhs, fTy)
					ptr := c.nextReg()
					c.emit("%s = getelementptr inbounds %%%s, %%%s* %s, i32 0, i32 %d", ptr, stName, stName, left.reg, idx)
					c.emit("store %s %s, %s* %s", fTy, rhs.reg, fTy, ptr)
					return rhs
				}
			}
		}
		c.addErr(fmt.Errorf("Native codegen: assign requires identifier or dot-expr on the left, got %T", t.Left))
		return value{reg: "0", typ: "i64"}
	case *candy_ast.DotExpression:
		left := c.compileExpression(t.Left)
		if left.typ == "%any" {
			namePtr := c.emitCString(t.Right.Value)
			dataPtr := c.nextReg()
			c.emit("%s = getelementptr inbounds %%any, %%any* %s, i32 0, i32 1", dataPtr, left.reg)
			raw := c.nextReg()
			c.emit("%s = load i8*, i8** %s", raw, dataPtr)
			out := c.nextReg()
			c.emit("%s = call i64 @candy_dyn_get_i64(i8* %s, i8* %s)", out, raw, namePtr)
			return value{reg: out, typ: "i64"}
		}
		if !strings.HasSuffix(left.typ, "*") {
			return value{reg: "0", typ: "i64"}
		}
		// Extract struct name from "%S*"
		stName := strings.TrimSuffix(strings.TrimPrefix(left.typ, "%"), "*")
		st0 := c.structNodes[stName]
		if st0 != nil {
			own, prop := c.findPropertyDefining(st0, t.Right.Value)
			if prop != nil {
				if prop.Getter == nil {
					c.addErr(fmt.Errorf("codegen: property %s is missing a getter in LLVM path", t.Right.Value))
					return value{reg: "0", typ: "i64"}
				}
				propLL := c.mapCandyTypeToLlvm(candy_ast.ExprAsSimpleTypeName(prop.Type))
				recv := c.bitcastStructPtr(left, own)
				reg := c.nextReg()
				c.emit("%s = call %s @%s_get_%s(%%%s* %s)", reg, propLL, own, t.Right.Value, own, recv.reg)
				return value{reg: reg, typ: propLL}
			}
		}
		idx, ok := c.structs[stName][t.Right.Value]
		if !ok {
			return value{reg: "0", typ: "i64"}
		}

		// Find the field type
		// Quick way: extract from the { ... } string
		// Better way: use a map. I already have mapCandyTypeToLlvm for the definition
		// Let's just use the index to find the type from the definition

		// Need to know the field's LLVM type
		// For simplicity, I'll re-extract or store it.
		// Let's just assume i64 for now or extract from structTy

		fTy := "i64" // Default fallback
		// Extract type from structTy[stName]
		tyStr := c.structTy[stName] // "{ i64, double }"
		tyStr = strings.Trim(tyStr, "{}")
		parts := strings.Split(tyStr, ",")
		if idx < len(parts) {
			fTy = strings.TrimSpace(parts[idx])
		}

		ptr := c.nextReg()
		c.emit("%s = getelementptr inbounds %%%s, %%%s* %s, i32 0, i32 %d", ptr, stName, stName, left.reg, idx)
		reg := c.nextReg()
		c.emit("%s = load %s, %s* %s", reg, fTy, fTy, ptr)
		return value{reg: reg, typ: fTy}
	case *candy_ast.TernaryExpression:
		return c.compileTernary(t)
	default:
		c.addErr(fmt.Errorf("Native codegen: unsupported expression %T", expr))
		return value{reg: "0", typ: "i64"}
	}
}

func (c *Compiler) compileInfix(t *candy_ast.InfixExpression) value {
	if folded, ok := c.tryFoldLiteralInfix(t); ok {
		return folded
	}

	left := c.compileExpression(t.Left)
	right := c.compileExpression(t.Right)

	if v, ok := c.tryCompileStructInfix(t, left, right); ok {
		return v
	}

	if left.typ == "%any" || right.typ == "%any" {
		// Dynamic dispatch for math
		lAny := c.boxValueIntoAny(left)
		rAny := c.boxValueIntoAny(right)

		fnName := ""
		switch t.Operator {
		case "+":
			fnName = "candy_dyn_add"
		case "-":
			fnName = "candy_dyn_sub"
		case "*":
			fnName = "candy_dyn_mul"
		case "/":
			fnName = "candy_dyn_div"
		}

		if fnName != "" {
			reg := c.nextReg()
			c.emit("%s = call %%any @%s(%%any %s, %%any %s)", reg, fnName, lAny.reg, rAny.reg)
			return value{reg: reg, typ: "%any"}
		}
	}

	// Promotion
	if left.typ == "double" || right.typ == "double" {
		if left.typ == "i64" {
			reg := c.nextReg()
			c.emit("%s = sitofp i64 %s to double", reg, left.reg)
			left = value{reg: reg, typ: "double"}
		}
		if right.typ == "i64" {
			reg := c.nextReg()
			c.emit("%s = sitofp i64 %s to double", reg, right.reg)
			right = value{reg: reg, typ: "double"}
		}
	}

	isFloat := left.typ == "double"

	switch t.Operator {
	case "+":
		if left.typ == "i8*" || right.typ == "i8*" {
			// String concatenation
			lStr := left
			if lStr.typ == "i64" {
				reg := c.nextReg()
				c.emit("%s = call i8* @candy_int_to_str(i64 %s)", reg, lStr.reg)
				lStr = value{reg: reg, typ: "i8*"}
			}
			rStr := right
			if rStr.typ == "i64" {
				reg := c.nextReg()
				c.emit("%s = call i8* @candy_int_to_str(i64 %s)", reg, rStr.reg)
				rStr = value{reg: reg, typ: "i8*"}
			}

			reg := c.nextReg()
			c.emit("%s = call i8* @candy_str_add(i8* %s, i8* %s)", reg, lStr.reg, rStr.reg)
			return value{reg: reg, typ: "i8*"}
		}
		op := "add"
		if isFloat {
			op = "fadd"
		}
		reg := c.nextReg()
		c.emit("%s = %s %s %s, %s", reg, op, left.typ, left.reg, right.reg)
		return value{reg: reg, typ: left.typ}
	case "-":
		op := "sub"
		if isFloat {
			op = "fsub"
		}
		reg := c.nextReg()
		c.emit("%s = %s %s %s, %s", reg, op, left.typ, left.reg, right.reg)
		return value{reg: reg, typ: left.typ}
	case "*":
		op := "mul"
		if isFloat {
			op = "fmul"
		}
		reg := c.nextReg()
		c.emit("%s = %s %s %s, %s", reg, op, left.typ, left.reg, right.reg)
		return value{reg: reg, typ: left.typ}
	case "/":
		op := "sdiv"
		if isFloat {
			op = "fdiv"
		}
		reg := c.nextReg()
		c.emit("%s = %s %s %s, %s", reg, op, left.typ, left.reg, right.reg)
		return value{reg: reg, typ: left.typ}
	case "==", "!=", "<", ">":
		op := "icmp"
		cond := "eq"
		switch t.Operator {
		case "!=":
			cond = "ne"
		case "<":
			cond = "slt"
		case ">":
			cond = "sgt"
		}
		if isFloat {
			op = "fcmp"
			cond = "oeq"
			switch t.Operator {
			case "!=":
				cond = "one"
			case "<":
				cond = "olt"
			case ">":
				cond = "ogt"
			}
		}
		reg := c.nextReg()
		c.emit("%s = %s %s %s %s, %s", reg, op, cond, left.typ, left.reg, right.reg)
		return value{reg: reg, typ: "i1"}
	case "|":
		if left.typ != "i64" || right.typ != "i64" {
			c.addErr(fmt.Errorf("Native codegen: bitwise | expects i64 operands"))
			return value{reg: "0", typ: "i64"}
		}
		reg := c.nextReg()
		c.emit("%s = or i64 %s, %s", reg, left.reg, right.reg)
		return value{reg: reg, typ: "i64"}
	case "&":
		if left.typ != "i64" || right.typ != "i64" {
			c.addErr(fmt.Errorf("Native codegen: bitwise & expects i64 operands"))
			return value{reg: "0", typ: "i64"}
		}
		reg := c.nextReg()
		c.emit("%s = and i64 %s, %s", reg, left.reg, right.reg)
		return value{reg: reg, typ: "i64"}
	case "^":
		if left.typ != "i64" || right.typ != "i64" {
			c.addErr(fmt.Errorf("Native codegen: bitwise ^ expects i64 operands"))
			return value{reg: "0", typ: "i64"}
		}
		reg := c.nextReg()
		c.emit("%s = xor i64 %s, %s", reg, left.reg, right.reg)
		return value{reg: reg, typ: "i64"}
	case "<<":
		if left.typ != "i64" || right.typ != "i64" {
			c.addErr(fmt.Errorf("Native codegen: bitwise << expects i64 operands"))
			return value{reg: "0", typ: "i64"}
		}
		reg := c.nextReg()
		c.emit("%s = shl i64 %s, %s", reg, left.reg, right.reg)
		return value{reg: reg, typ: "i64"}
	case ">>":
		if left.typ != "i64" || right.typ != "i64" {
			c.addErr(fmt.Errorf("Native codegen: bitwise >> expects i64 operands"))
			return value{reg: "0", typ: "i64"}
		}
		reg := c.nextReg()
		c.emit("%s = ashr i64 %s, %s", reg, left.reg, right.reg)
		return value{reg: reg, typ: "i64"}
	default:
		c.addErr(fmt.Errorf("Native codegen: unsupported infix operator %q", t.Operator))
		return value{reg: "0", typ: "i64"}
	}
}

func (c *Compiler) tryFoldLiteralInfix(t *candy_ast.InfixExpression) (value, bool) {
	if t == nil || t.Left == nil || t.Right == nil {
		return value{}, false
	}

	li, liok := t.Left.(*candy_ast.IntegerLiteral)
	ri, riok := t.Right.(*candy_ast.IntegerLiteral)
	if liok && riok {
		lv, rv := li.Value, ri.Value
		switch t.Operator {
		case "+":
			return value{reg: fmt.Sprintf("%d", lv+rv), typ: "i64"}, true
		case "-":
			return value{reg: fmt.Sprintf("%d", lv-rv), typ: "i64"}, true
		case "*":
			return value{reg: fmt.Sprintf("%d", lv*rv), typ: "i64"}, true
		case "/":
			if rv == 0 {
				return value{}, false
			}
			return value{reg: fmt.Sprintf("%d", lv/rv), typ: "i64"}, true
		case "==":
			if lv == rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "!=":
			if lv != rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "<":
			if lv < rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case ">":
			if lv > rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		}
	}

	lf, lfok := t.Left.(*candy_ast.FloatLiteral)
	rf, rfok := t.Right.(*candy_ast.FloatLiteral)
	if lfok && rfok {
		lv, rv := lf.Value, rf.Value
		switch t.Operator {
		case "+":
			return value{reg: fmt.Sprintf("%f", lv+rv), typ: "double"}, true
		case "-":
			return value{reg: fmt.Sprintf("%f", lv-rv), typ: "double"}, true
		case "*":
			return value{reg: fmt.Sprintf("%f", lv*rv), typ: "double"}, true
		case "/":
			if rv == 0 {
				return value{}, false
			}
			return value{reg: fmt.Sprintf("%f", lv/rv), typ: "double"}, true
		case "==":
			if lv == rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "!=":
			if lv != rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "<":
			if lv < rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case ">":
			if lv > rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		}
	}

	// Mixed numeric literals (int/float): fold to double.
	if (liok || lfok) && (riok || rfok) {
		lv := 0.0
		rv := 0.0
		if liok {
			lv = float64(li.Value)
		} else if lfok {
			lv = lf.Value
		}
		if riok {
			rv = float64(ri.Value)
		} else if rfok {
			rv = rf.Value
		}
		switch t.Operator {
		case "+":
			return value{reg: fmt.Sprintf("%f", lv+rv), typ: "double"}, true
		case "-":
			return value{reg: fmt.Sprintf("%f", lv-rv), typ: "double"}, true
		case "*":
			return value{reg: fmt.Sprintf("%f", lv*rv), typ: "double"}, true
		case "/":
			if rv == 0 {
				return value{}, false
			}
			return value{reg: fmt.Sprintf("%f", lv/rv), typ: "double"}, true
		case "==":
			if lv == rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "!=":
			if lv != rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "<":
			if lv < rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case ">":
			if lv > rv {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		}
	}

	ls, lsok := t.Left.(*candy_ast.StringLiteral)
	rs, rsok := t.Right.(*candy_ast.StringLiteral)
	if lsok && rsok {
		switch t.Operator {
		case "+":
			return c.compileExpression(&candy_ast.StringLiteral{Value: ls.Value + rs.Value}), true
		case "==":
			if ls.Value == rs.Value {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "!=":
			if ls.Value != rs.Value {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		}
	}

	lb, lbok := t.Left.(*candy_ast.Boolean)
	rb, rbok := t.Right.(*candy_ast.Boolean)
	if lbok && rbok {
		switch t.Operator {
		case "==":
			if lb.Value == rb.Value {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		case "!=":
			if lb.Value != rb.Value {
				return value{reg: "1", typ: "i1"}, true
			}
			return value{reg: "0", typ: "i1"}, true
		}
	}

	return value{}, false
}
