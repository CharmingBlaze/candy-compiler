package candy_evaluator

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

// evalStructLiteral creates a new struct *Value from a struct type registered in the env.
func evalStructLiteral(t *candy_ast.StructLiteral, e *Env) (*Value, error) {
	if t == nil {
		return nil, nil
	}
	name := candy_ast.ExprAsSimpleTypeName(t.Name)
	if name == "" {
		return nil, &RuntimeError{Msg: "struct literal: missing type name"}
	}
	tmpl, ok := e.Get(name)
	if !ok || tmpl == nil {
		return nil, &RuntimeError{Msg: "unknown struct type: " + name}
	}
	if tmpl.Kind != ValStruct || tmpl.St == nil {
		return nil, fmt.Errorf("struct literal: %s is not a struct type", name)
	}
	// Class/object literal style: `Platform { ... }`
	// Instantiate class/object runtime and then overlay explicit fields.
	if tmpl.St.ClassDef != nil {
		inst, err := instantiateClass(tmpl, nil, e)
		if err != nil {
			return nil, err
		}
		if inst == nil || inst.St == nil {
			return nil, &RuntimeError{Msg: "failed to instantiate class literal: " + name}
		}
		for fname, ex := range t.Fields {
			v, err := evalExpression(ex, e)
			if err != nil {
				return nil, err
			}
			inst.St.Data[fname] = valueToValue(v)
		}
		return inst, nil
	}
	data := make(map[string]Value)
	// Apply struct field defaults first, then override with explicit literal fields.
	if tmpl.St.Def != nil {
		for _, f := range tmpl.St.Def.Fields {
			if f.Name == nil {
				continue
			}
			if f.Init == nil {
				continue
			}
			dv, err := evalExpression(f.Init, e)
			if err != nil {
				return nil, err
			}
			data[f.Name.Value] = valueToValue(dv)
		}
	}
	for fname, ex := range t.Fields {
		v, err := evalExpression(ex, e)
		if err != nil {
			return nil, err
		}
		data[fname] = valueToValue(v)
	}
	return &Value{
		Kind: ValStruct,
		St:   &structVal{Def: tmpl.St.Def, Env: e, Data: data},
	}, nil
}

// evalDot reads obj.field
func evalDot(t *candy_ast.DotExpression, e *Env) (*Value, error) {
	if t == nil || t.Right == nil {
		return &Value{Kind: ValNull}, nil
	}
	left, err := evalExpression(t.Left, e)
	if err != nil {
		return nil, err
	}
	if left == nil || left.Kind == ValNull {
		if t.IsSafe {
			return &Value{Kind: ValNull}, nil
		}
		return nil, &RuntimeError{Msg: "null receiver for ." + t.Right.Value}
	}
	if left.Kind == ValModule && left.Mod != nil {
		key := t.Right.Value
		if left.Mod.Consts != nil {
			for k, c := range left.Mod.Consts {
				if strings.EqualFold(k, key) {
					return c, nil
				}
			}
		}
		return &Value{Kind: ValNull}, nil
	}
	if left.Kind == ValArray {
		return arrayProps(left, t.Right.Value), nil
	}
	if left.Kind == ValVec {
		if p := vecProps(left, t.Right.Value); p != nil {
			return p, nil
		}
		return nil, &RuntimeError{Msg: "vector has no property: " + t.Right.Value}
	}
	if left.Kind == ValString {
		if p := stringProps(left, t.Right.Value); p != nil {
			return p, nil
		}
		return nil, &RuntimeError{Msg: "string has no property: " + t.Right.Value}
	}
	if left.Kind == ValMap {
		key := t.Right.Value
		if left.StrMap != nil {
			if v, ok := left.StrMap[key]; ok {
				return ptrVal(v), nil
			}
			for k, v := range left.StrMap {
				if strings.EqualFold(k, key) {
					return ptrVal(v), nil
				}
			}
		}
		return &Value{Kind: ValNull}, nil
	}
	if left.Kind != ValStruct || left.St == nil {
		return nil, &RuntimeError{Msg: "not a struct, module, list, or string value for `.` (got " + left.String() + ")"}
	}
	key := t.Right.Value
	typeName := "struct"
	if left.St.Def != nil && left.St.Def.Name != nil {
		typeName = left.St.Def.Name.Value
	}
	if v, ok := left.St.Data[key]; ok {
		return ptrVal(v), nil
	}
	// Inheritance: check base class for members if not found in data
	curr := left.St.ClassDef
	for curr != nil {
		for _, m := range curr.Members {
			switch st := m.(type) {
			case *candy_ast.VarStatement:
				if strings.EqualFold(st.Name.Value, key) {
					// If it's in data, it was already found. If not, maybe it's a default?
					// For now, we only look at Data for instances.
				}
			}
		}
		if curr.Base != nil {
			if baseVal, ok := e.Get(curr.Base.Value); ok && baseVal.Kind == ValStruct && baseVal.St != nil {
				curr = baseVal.St.ClassDef
				continue
			}
		}
		break
	}

	if left.St.Def != nil {
		for _, prop := range left.St.Def.Properties {
			if prop == nil || prop.Name == nil || !strings.EqualFold(prop.Name.Value, key) {
				continue
			}
			if prop.Getter != nil {
				return evalStructPropertyGetter(left, prop, e)
			}
			if v, ok := left.St.Data[prop.Name.Value]; ok {
				return ptrVal(v), nil
			}
		}
	}
	// Also check properties in ClassDef hierarchy
	curr = left.St.ClassDef
	for curr != nil {
		for _, m := range curr.Members {
			if prop, ok := m.(*candy_ast.PropertyStatement); ok {
				if strings.EqualFold(prop.Name.Value, key) {
					if prop.Getter != nil {
						return evalStructPropertyGetter(left, prop, e)
					}
				}
			}
		}
		if curr.Base != nil {
			if baseVal, ok := e.Get(curr.Base.Value); ok && baseVal.Kind == ValStruct && baseVal.St != nil {
				curr = baseVal.St.ClassDef
				continue
			}
		}
		break
	}
	var names []string
	for k := range left.St.Data {
		names = append(names, k)
	}
	for k, v := range left.St.Data {
		if strings.EqualFold(k, key) {
			return ptrVal(v), nil
		}
	}
	return nil, withDidYouMean(typeName, key, names)
}

func evalStructPropertyGetter(inst *Value, prop *candy_ast.PropertyStatement, outer *Env) (*Value, error) {
	if inst == nil || inst.St == nil || prop == nil || prop.Getter == nil {
		return &Value{Kind: ValNull}, nil
	}
	ne := outer.NewEnclosed()
	for k, v := range inst.St.Data {
		vv := v
		ne.Set(k, &vv)
	}
	ne.Set("this", inst)
	for i, st := range prop.Getter.Statements {
		r, err := evalStatement(st, ne)
		if err != nil {
			return nil, err
		}
		if rw, ok := r.(ReturnWrap); ok {
			return rw.V, nil
		}
		// Implicit return
		if i == len(prop.Getter.Statements)-1 {
			if v, ok := r.(*Value); ok {
				return v, nil
			}
		}
	}
	return &Value{Kind: ValNull}, nil
}

func evalStructPropertySetter(inst *Value, prop *candy_ast.PropertyStatement, rhs *Value, outer *Env) error {
	if inst == nil || inst.St == nil || prop == nil {
		return nil
	}
	if prop.Setter == nil {
		if prop.Name != nil {
			inst.St.Data[prop.Name.Value] = valueToValue(rhs)
		}
		return nil
	}
	ne := outer.NewEnclosed()
	for k, v := range inst.St.Data {
		vv := v
		ne.Set(k, &vv)
	}
	ne.Set("this", inst)
	ne.Set("value", rhs)
	for _, st := range prop.Setter.Statements {
		if _, err := evalStatement(st, ne); err != nil {
			return err
		}
	}
	for k := range inst.St.Data {
		if v, ok := ne.Get(k); ok {
			inst.St.Data[k] = *v
		}
	}
	return nil
}

func arrayProps(v *Value, name string) *Value {
	ln := strings.ToLower(name)
	switch ln {
	case "length", "size", "count", "len":
		return &Value{Kind: ValInt, I64: int64(len(v.Elems))}
	case "is_empty", "isempty", "empty":
		return &Value{Kind: ValBool, B: len(v.Elems) == 0}
	default:
		return &Value{Kind: ValNull}
	}
}

// stringProps returns length / is_empty for bare property access: text.length
func stringProps(v *Value, name string) *Value {
	ln := strings.ToLower(name)
	if v == nil || v.Kind != ValString {
		return &Value{Kind: ValNull}
	}
	s := v.Str
	runes := len([]rune(s))
	switch ln {
	case "length", "size", "count", "len":
		return &Value{Kind: ValInt, I64: int64(runes)}
	case "is_empty", "isempty", "empty":
		return &Value{Kind: ValBool, B: s == ""}
	default:
		return nil
	}
}

// ptrVal returns a pointer to a copy of v for safe use in caller (avoids taking address of map slot).
func ptrVal(v Value) *Value {
	vv := v
	return &vv
}

// evalMethodCall runs `value.method(args)` for struct body methods. Binds `this` (or receiver name) to the instance.
func evalMethodCall(dot *candy_ast.DotExpression, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	if dot == nil || dot.Right == nil {
		return &Value{Kind: ValNull}, nil
	}
	recv, err := evalExpression(dot.Left, e)
	if err != nil {
		return nil, err
	}
	if recv == nil || recv.Kind == ValNull {
		if dot.IsSafe {
			return &Value{Kind: ValNull}, nil
		}
		return nil, &RuntimeError{Msg: "method call: null on the left of `.`"}
	}
	// stdlib host modules, arrays, maps
	if recv.Kind == ValModule && recv.Mod != nil {
		name := dot.Right.Value
		if fn, ok := lookupModFn(recv.Mod, name); ok {
			args, err2 := evalArgs(argExprs, e)
			if err2 != nil {
				return nil, err2
			}
			return fn(args)
		}
		return nil, &RuntimeError{Msg: recv.Mod.Name + " has no method " + name}
	}
	if recv.Kind == ValArray {
		return callArrayMethod(recv, dot.Right.Value, argExprs, e)
	}
	if recv.Kind == ValVec {
		return callVecMethod(recv, dot.Right.Value, argExprs, e)
	}
	if recv.Kind == ValMap {
		return callMapMethod(recv, dot.Right.Value, argExprs, e)
	}
	if recv.Kind == ValString {
		return callStringMethod(recv, dot.Right.Value, argExprs, e)
	}
	if recv.Kind != ValStruct || recv.St == nil {
		return nil, &RuntimeError{Msg: "method call: need struct, module, list, map, or string on the left of `.`"}
	}
	name := dot.Right.Value
	var method *candy_ast.FunctionStatement
	if recv.St.Def != nil {
		for _, m := range recv.St.Def.Methods {
			if m != nil && m.Name != nil && strings.EqualFold(m.Name.Value, name) {
				method = m
				break
			}
		}
	}
	// Search ClassDef and its inheritance chain
	currCls := recv.St.ClassDef
	for currCls != nil && method == nil {
		for _, m := range currCls.Members {
			if ms, ok := m.(*candy_ast.FunctionStatement); ok {
				if ms.Name != nil && strings.EqualFold(ms.Name.Value, name) {
					method = ms
					break
				}
			}
		}
		if method == nil && currCls.Base != nil {
			if baseVal, ok := e.Get(currCls.Base.Value); ok && baseVal.Kind == ValStruct && baseVal.St != nil {
				currCls = baseVal.St.ClassDef
			} else {
				break
			}
		} else {
			break
		}
	}

	if method == nil {
		avail := []string{}
		if recv.St.ClassDef != nil {
			for _, m := range recv.St.ClassDef.Members {
				if ms, ok := m.(*candy_ast.FunctionStatement); ok {
					avail = append(avail, fmt.Sprintf("%s(%T)", ms.Name.Value, m))
				} else {
					avail = append(avail, fmt.Sprintf("%T", m))
				}
			}
		}
		return nil, &RuntimeError{Msg: fmt.Sprintf("no such method: %s (available: %v)", name, avail)}
	}
	args, namedArgs, err := evalCallArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	if len(namedArgs) > 0 {
		args = reorderCallArgs(method, args, namedArgs)
	}
	outer := recv.St.Env
	if outer == nil {
		outer = e
	}
	ne := outer.NewEnclosed()
	// Bind all instance data fields into the environment (heap pointers so mutations persist).
	for k, v := range recv.St.Data {
		ptr := new(Value)
		*ptr = v
		ne.Set(k, ptr)
	}

	recvName := "this"
	if method.Receiver != nil && method.Receiver.Name != nil {
		recvName = method.Receiver.Name.Value
	}
	ne.Set(recvName, recv)

	// Bind 'super' if there is a base class
	if recv.St.ClassDef != nil && recv.St.ClassDef.Base != nil {
		if baseVal, ok := e.Get(recv.St.ClassDef.Base.Value); ok && baseVal.Kind == ValStruct && baseVal.St != nil {
			superWrapper := &Value{
				Kind: ValStruct,
				St: &structVal{
					ClassDef: baseVal.St.ClassDef,
					Data:     recv.St.Data,
					Env:      recv.St.Env,
				},
			}
			ne.Set("super", superWrapper)
		}
	}
	for i, p0 := range method.Parameters {
		if i < len(args) {
			ne.Set(p0.Name.Value, args[i])
		} else if p0.Default != nil {
			dv, err := evalExpression(p0.Default, outer)
			if err != nil {
				return nil, err
			}
			ne.Set(p0.Name.Value, dv)
		}
	}
	if method.Body == nil {
		return &Value{Kind: ValNull}, nil
	}
	// Instance fields may be mutated (e.g. list.remove_at); write back even on early return.
	defer func() {
		for k := range recv.St.Data {
			if v, ok := ne.Get(k); ok {
				recv.St.Data[k] = *v
			}
		}
	}()
	for i, st := range method.Body.Statements {
		r, err3 := evalStatement(st, ne)
		if err3 != nil {
			return nil, err3
		}
		if rw, ok2 := r.(ReturnWrap); ok2 {
			return rw.V, nil
		}
		// Implicit return
		if i == len(method.Body.Statements)-1 {
			if v, ok := r.(*Value); ok {
				return v, nil
			}
		}
	}
	return &Value{Kind: ValNull}, nil
}

func evalClassStatement(t *candy_ast.ClassStatement, e *Env) (*Value, error) {
	if t == nil || t.Name == nil {
		return nil, nil
	}
	v := &Value{Kind: ValStruct, St: &structVal{ClassDef: t, Env: e, Data: make(map[string]Value)}}
	e.Set(t.Name.Value, v)
	return v, nil
}

func evalObjectStatement(t *candy_ast.ObjectStatement, e *Env) (*Value, error) {
	if t == nil || t.Name == nil {
		return nil, nil
	}
	cls := &candy_ast.ClassStatement{
		Token:   t.Token,
		Name:    t.Name,
		Base:    t.Base,
		Members: t.Members,
	}
	// Objects are singletons: we create the definition AND an instance immediately.
	def := &Value{Kind: ValStruct, St: &structVal{ClassDef: cls, Env: e, Data: make(map[string]Value)}}
	instance, err := instantiateClass(def, nil, e)
	if err != nil {
		return nil, err
	}
	e.Set(t.Name.Value, instance)
	return instance, nil
}

func evalDelete(t *candy_ast.DeleteStatement, e *Env) (*Value, error) {
	if t.Value == nil {
		return &Value{Kind: ValNull}, nil
	}
	// delete(v) -> v = null
	if id, ok := t.Value.(*candy_ast.Identifier); ok {
		e.Set(id.Value, &Value{Kind: ValNull})
	}
	return &Value{Kind: ValNull}, nil
}
