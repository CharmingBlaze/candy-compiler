package candy_evaluator

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

// ValueKind is the kind of a runtime value.
type ValueKind int

const (
	ValNull ValueKind = iota
	ValInt
	ValFloat
	ValString
	ValBool
	ValArray
	ValMap
	ValFunction
	ValStruct
	ValVec
	ValResult // Result<T, E> tagged
	ValModule // built-in stdlib group (e.g. math) — method calls go to host functions
)

// Value is a K-Go runtime value.
type Value struct {
	Kind   ValueKind
	I64    int64
	F64    float64
	Str    string
	B      bool
	Elems  []Value // for array
	StrMap map[string]Value
	Fn     *functionVal
	St     *structVal
	Mod    *moduleVal
	ResOk  bool
	// Result<T,E> stored as: Res field + union with Err Value
	Res, Err *Value
	Builtin  func(args []*Value) (*Value, error)
	Vec      []float64
}

// String for debugging
func (v Value) String() string {
	if v.Kind == ValNull {
		return "null"
	}
	if v.Kind == ValInt {
		return fmt.Sprintf("%d", v.I64)
	}
	if v.Kind == ValFloat {
		return fmt.Sprintf("%g", v.F64)
	}
	if v.Kind == ValString {
		return v.Str
	}
	if v.Kind == ValBool {
		if v.B {
			return "true"
		}
		return "false"
	}
	if v.Kind == ValFunction {
		if v.Fn == nil {
			return "fn<nil>"
		}
		return fmt.Sprintf("fn %s", v.Fn.Stmt.Name.Value)
	}
	if v.Kind == ValArray {
		var s []string
		for _, e := range v.Elems {
			s = append(s, e.String())
		}
		return "[" + strings.Join(s, ", ") + "]"
	}
	if v.Kind == ValMap {
		if v.StrMap == nil {
			return "map(…)"
		}
		return "map(…len=" + fmt.Sprint(len(v.StrMap)) + ")"
	}
	if v.Kind == ValModule {
		if v.Mod == nil {
			return "module<nil>"
		}
		return fmt.Sprintf("module(%s)", v.Mod.Name)
	}
	if v.Kind == ValStruct {
		if v.St == nil {
			return "struct<nil>"
		}
		n := "<anon>"
		if v.St.Def != nil && v.St.Def.Name != nil {
			n = v.St.Def.Name.Value
		}
		return fmt.Sprintf("struct(%s){%d fields}", n, len(v.St.Data))
	}
	if v.Kind == ValVec {
		parts := make([]string, 0, len(v.Vec))
		for _, x := range v.Vec {
			parts = append(parts, fmt.Sprintf("%g", x))
		}
		return "vec(" + strings.Join(parts, ", ") + ")"
	}
	if v.Kind == ValResult {
		if v.ResOk {
			return "Ok(" + v.Res.String() + ")"
		}
		return "Err(" + v.Err.String() + ")"
	}
	return "?"
}

type functionVal struct {
	Stmt  *candy_ast.FunctionStatement
	Env   *Env
	Outer *Env
}

type structVal struct {
	Def      *candy_ast.StructStatement
	ClassDef *candy_ast.ClassStatement
	Env      *Env
	Data     map[string]Value
}

// moduleVal groups host-backed functions under a name (e.g. math, file).
type moduleVal struct {
	Name   string
	Fns    map[string]func(args []*Value) (*Value, error)
	Consts map[string]*Value
}

// Truthy: false for null, false, 0, 0.0, empty string.
func (v *Value) Truthy() bool {
	if v == nil {
		return false
	}
	switch v.Kind {
	case ValNull:
		return false
	case ValBool:
		return v.B
	case ValInt:
		return v.I64 != 0
	case ValFloat:
		return v.F64 != 0
	case ValString:
		return v.Str != ""
	case ValArray:
		return true
	case ValMap:
		return len(v.StrMap) > 0
	default:
		return true
	}
}

// ReturnWrap marks a return value in nested blocks.
type ReturnWrap struct{ V *Value }

// BreakWrap signals a loop break.
type BreakWrap struct{}

// ContinueWrap signals a loop continue.
type ContinueWrap struct{}
