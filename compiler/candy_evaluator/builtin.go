package candy_evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func builtinLen(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument, got %d", len(args))
	}
	switch args[0].Kind {
	case ValString:
		return &Value{Kind: ValInt, I64: int64(len(args[0].Str))}, nil
	case ValArray:
		return &Value{Kind: ValInt, I64: int64(len(args[0].Elems))}, nil
	default:
		return nil, fmt.Errorf("len: unsupported %v", args[0].Kind)
	}
}

func builtinOk(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Ok wants 1 arg, got %d", len(args))
	}
	return &Value{Kind: ValResult, ResOk: true, Res: args[0]}, nil
}

func builtinErr(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Err wants 1 arg, got %d", len(args))
	}
	return &Value{Kind: ValResult, ResOk: false, Err: args[0]}, nil
}

func builtinCwd(args []*Value) (*Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("cwd expects 0 arguments, got %d", len(args))
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValString, Str: wd}, nil
}

func builtinJoinPath(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("joinPath expects at least 2 arguments, got %d", len(args))
	}
	parts := make([]string, 0, len(args))
	for i, a := range args {
		if a == nil || a.Kind != ValString {
			return nil, fmt.Errorf("joinPath arg %d must be string", i+1)
		}
		parts = append(parts, a.Str)
	}
	return &Value{Kind: ValString, Str: filepath.Join(parts...)}, nil
}

func builtinReadFile(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("readFile expects 1 argument, got %d", len(args))
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("readFile path must be string")
	}
	b, err := os.ReadFile(args[0].Str)
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValString, Str: string(b)}, nil
}

func builtinWriteFile(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("writeFile expects 2 arguments, got %d", len(args))
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("writeFile path must be string")
	}
	if args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("writeFile contents must be string")
	}
	if err := os.WriteFile(args[0].Str, []byte(args[1].Str), 0o644); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

func builtinGetEnv(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getEnv expects 1 argument, got %d", len(args))
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("getEnv key must be string")
	}
	v, ok := os.LookupEnv(args[0].Str)
	if !ok {
		return &Value{Kind: ValNull}, nil
	}
	return &Value{Kind: ValString, Str: v}, nil
}

func builtinPrint(args []*Value) (*Value, error) {
	var b strings.Builder
	for i, a := range args {
		if i > 0 {
			b.WriteString(" ")
		}
		if a == nil {
			b.WriteString("null")
			continue
		}
		b.WriteString(a.String())
	}
	fmt.Println(b.String())
	return &Value{Kind: ValNull}, nil
}

func builtinArray(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValInt {
		return nil, fmt.Errorf("array expects 1 int size argument")
	}
	n := args[0].I64
	if n < 0 {
		return nil, fmt.Errorf("array size must be >= 0")
	}
	elems := make([]Value, n)
	for i := int64(0); i < n; i++ {
		elems[i] = Value{Kind: ValNull}
	}
	return &Value{Kind: ValArray, Elems: elems}, nil
}

func builtinBytes(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValInt {
		return nil, fmt.Errorf("bytes expects 1 int size argument")
	}
	n := args[0].I64
	if n < 0 {
		return nil, fmt.Errorf("bytes size must be >= 0")
	}
	elems := make([]Value, n)
	for i := int64(0); i < n; i++ {
		elems[i] = Value{Kind: ValInt, I64: 0}
	}
	return &Value{Kind: ValArray, Elems: elems}, nil
}

// Builtins is the set of pre-defined functions (name -> impl).
var Builtins = map[string]func(args []*Value) (*Value, error){
	"print":     builtinPrint,
	"println":   builtinPrint,
	"len":       builtinLen,
	"length":    builtinLen,
	"ok":        builtinOk,
	"err":       builtinErr,
	"cwd":       builtinCwd,
	"joinpath":  builtinJoinPath,
	"readfile":  builtinReadFile,
	"writefile": builtinWriteFile,
	"getenv":    builtinGetEnv,
	"new":       builtinNew,
	"input":     builtinInput,
	"random":    builtinRandomInt, // alias for random.int
	"array":     builtinArray,
	"bytes":     builtinBytes,
}

func builtinInput(args []*Value) (*Value, error) {
	if len(args) > 0 {
		fmt.Print(args[0].String())
	}
	var line string
	fmt.Scanln(&line)
	return &Value{Kind: ValString, Str: line}, nil
}

func builtinNew(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("new expects 1 argument")
	}
	// If the argument is a class definition (has Env set), instantiate it with no args
	if args[0].Kind == ValStruct && (args[0].St.Def != nil || args[0].St.ClassDef != nil) {
		if args[0].St.Env != nil {
			return nil, fmt.Errorf("new <ClassName> without parens not supported yet, use <ClassName>()")
		}
	}
	return args[0], nil
}

func RegisterBuiltin(name string, fn func(args []*Value) (*Value, error)) {
	Builtins[strings.ToLower(name)] = fn
}
