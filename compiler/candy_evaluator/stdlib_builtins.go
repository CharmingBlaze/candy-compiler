package candy_evaluator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	processStartTime = time.Now()
	lastDeltaCall    = processStartTime
	exitProgram      = os.Exit
	globalSaveData   = map[string]Value{}
)

// --- small helpers for builtins ---

func arg1(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return &Value{Kind: ValNull}, nil
	}
	return args[0], nil
}

func arg0(args []*Value) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 arguments, got %d", len(args))
	}
	return nil
}

func f64Arg(v *Value) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("expected number")
	}
	switch v.Kind {
	case ValInt:
		return float64(v.I64), nil
	case ValFloat:
		return v.F64, nil
	default:
		return 0, fmt.Errorf("expected number")
	}
}

func i64Arg(v *Value) (int64, error) {
	if v == nil {
		return 0, fmt.Errorf("expected integer")
	}
	if v.Kind == ValInt {
		return v.I64, nil
	}
	if v.Kind == ValFloat {
		return int64(v.F64), nil
	}
	return 0, fmt.Errorf("expected integer")
}

// --- I/O (also registered as snake_case aliases in init) ---

func builtinReadLines(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("readLines path must be string")
	}
	b, e := os.ReadFile(a.Str)
	if e != nil {
		return nil, e
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	out := make([]Value, 0, len(lines))
	for _, l := range lines {
		out = append(out, Value{Kind: ValString, Str: l})
	}
	return &Value{Kind: ValArray, Elems: out}, nil
}

func builtinFileExists(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("fileExists path must be string")
	}
	_, e := os.Stat(a.Str)
	if e == nil {
		return &Value{Kind: ValBool, B: true}, nil
	}
	if os.IsNotExist(e) {
		return &Value{Kind: ValBool, B: false}, nil
	}
	return nil, e
}

func builtinDeleteFile(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("deleteFile path must be string")
	}
	if e := os.Remove(a.Str); e != nil {
		return nil, e
	}
	return &Value{Kind: ValNull}, nil
}

func builtinAppendFile(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("appendFile: path, content")
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("appendFile path must be string")
	}
	if args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("appendFile content must be string")
	}
	f, err := os.OpenFile(args[0].Str, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := f.WriteString(args[1].Str); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

func builtinListFiles(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("listFiles path must be string")
	}
	ents, e := os.ReadDir(a.Str)
	if e != nil {
		return nil, e
	}
	names := make([]string, 0, len(ents))
	for _, ent := range ents {
		names = append(names, ent.Name())
	}
	sort.Strings(names)
	out := make([]Value, 0, len(names))
	for _, n := range names {
		out = append(out, Value{Kind: ValString, Str: n})
	}
	return &Value{Kind: ValArray, Elems: out}, nil
}

// --- JSON ---

func builtinJsonStringify(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("jsonStringify needs at least 1 arg")
	}
	js, e := valueToJSONable(args[0])
	if e != nil {
		return nil, e
	}
	b, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValString, Str: string(b)}, nil
}

func builtinJsonParse(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("jsonParse: string")
	}
	var x interface{}
	if e := json.Unmarshal([]byte(a.Str), &x); e != nil {
		return nil, e
	}
	return jsonableToValue(x)
}

func builtinLoadJSON(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("loadJson: path string")
	}
	b, e := os.ReadFile(a.Str)
	if e != nil {
		return nil, e
	}
	var x interface{}
	if e := json.Unmarshal(b, &x); e != nil {
		return nil, e
	}
	return jsonableToValue(x)
}

func builtinSaveJSON(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("saveJson: path, value")
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("saveJson: path string")
	}
	js, e := valueToJSONable(args[1])
	if e != nil {
		return nil, e
	}
	b, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(args[0].Str, b, 0o644); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

// --- math (single-arity helpers) ---

func oneF64(fn func(float64) float64) func([]*Value) (*Value, error) {
	return func(args []*Value) (*Value, error) {
		a, err := arg1(args)
		if err != nil {
			return nil, err
		}
		x, err := f64Arg(a)
		if err != nil {
			return nil, err
		}
		return &Value{Kind: ValFloat, F64: fn(x)}, nil
	}
}

var (
	builtinSqrt  = oneF64(math.Sqrt)
	builtinAbsF  = oneF64(func(f float64) float64 { return math.Abs(f) })
	builtinFloor = oneF64(math.Floor)
	builtinCeil  = oneF64(math.Ceil)
	builtinRound = oneF64(math.Round)
	builtinSin   = oneF64(math.Sin)
	builtinCos   = oneF64(math.Cos)
	builtinTan   = oneF64(math.Tan)
)

func builtinPow(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow: two numbers")
	}
	x, e1 := f64Arg(args[0])
	if e1 != nil {
		return nil, e1
	}
	y, e2 := f64Arg(args[1])
	if e2 != nil {
		return nil, e2
	}
	return &Value{Kind: ValFloat, F64: math.Pow(x, y)}, nil
}

func builtinMin(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("min: at least 1 value")
	}
	if args[0].Kind == ValInt {
		best := args[0].I64
		for i := 1; i < len(args); i++ {
			ii, err := i64Arg(args[i])
			if err != nil {
				return nil, err
			}
			if ii < best {
				best = ii
			}
		}
		return &Value{Kind: ValInt, I64: best}, nil
	}
	best, err := f64Arg(args[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(args); i++ {
		x, err2 := f64Arg(args[i])
		if err2 != nil {
			return nil, err2
		}
		if x < best {
			best = x
		}
	}
	return &Value{Kind: ValFloat, F64: best}, nil
}

func builtinMax(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("max: at least 1 value")
	}
	if args[0].Kind == ValInt {
		best := args[0].I64
		for i := 1; i < len(args); i++ {
			ii, err := i64Arg(args[i])
			if err != nil {
				return nil, err
			}
			if ii > best {
				best = ii
			}
		}
		return &Value{Kind: ValInt, I64: best}, nil
	}
	best, err := f64Arg(args[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(args); i++ {
		x, err2 := f64Arg(args[i])
		if err2 != nil {
			return nil, err2
		}
		if x > best {
			best = x
		}
	}
	return &Value{Kind: ValFloat, F64: best}, nil
}

func builtinClamp(args []*Value) (*Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("clamp: value, min, max")
	}
	v, e1 := f64Arg(args[0])
	lo, e2 := f64Arg(args[1])
	hi, e3 := f64Arg(args[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, fmt.Errorf("clamp: numbers only")
	}
	if lo > hi {
		lo, hi = hi, lo
	}
	if v < lo {
		v = lo
	}
	if v > hi {
		v = hi
	}
	if args[0].Kind == ValInt && args[1].Kind == ValInt && args[2].Kind == ValInt {
		return &Value{Kind: ValInt, I64: int64(v + 0.5)}, nil
	}
	return &Value{Kind: ValFloat, F64: v}, nil
}

func builtinLerp(args []*Value) (*Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("lerp: a, b, t")
	}
	a, e1 := f64Arg(args[0])
	b, e2 := f64Arg(args[1])
	t, e3 := f64Arg(args[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, fmt.Errorf("lerp: numbers only")
	}
	v := a + (b-a)*t
	if args[0].Kind == ValInt && args[1].Kind == ValInt && args[2].Kind == ValInt {
		return &Value{Kind: ValInt, I64: int64(v + 0.5)}, nil
	}
	return &Value{Kind: ValFloat, F64: v}, nil
}

// --- random ---

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
var stdinReader = bufio.NewReader(os.Stdin)

func builtinRandomInt(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("random: min, max inclusive")
	}
	lo, e1 := i64Arg(args[0])
	hi, e2 := i64Arg(args[1])
	if e1 != nil || e2 != nil {
		return nil, fmt.Errorf("random: integers")
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	n := hi - lo + 1
	if n <= 0 {
		return &Value{Kind: ValInt, I64: lo}, nil
	}
	return &Value{Kind: ValInt, I64: lo + rng.Int63n(n)}, nil
}

func builtinRandomFloat(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("randomFloat: min, max; uses [min,max) float")
	}
	lo, e1 := f64Arg(args[0])
	hi, e2 := f64Arg(args[1])
	if e1 != nil || e2 != nil {
		return nil, fmt.Errorf("randomFloat: numbers")
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	return &Value{Kind: ValFloat, F64: lo + rng.Float64()*(hi-lo)}, nil
}

func builtinChoose(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValArray || len(a.Elems) == 0 {
		return &Value{Kind: ValNull}, nil
	}
	i := rng.Intn(len(a.Elems))
	return ptrVal(a.Elems[i]), nil
}

func builtinRngSeed(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("seed: int")
	}
	n, e1 := i64Arg(args[0])
	if e1 != nil {
		return nil, e1
	}
	rng = rand.New(rand.NewSource(n))
	return &Value{Kind: ValNull}, nil
}

func builtinRandomShuffle(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValArray {
		return nil, fmt.Errorf("shuffle expects an array")
	}
	rng.Shuffle(len(a.Elems), func(i, j int) {
		a.Elems[i], a.Elems[j] = a.Elems[j], a.Elems[i]
	})
	return &Value{Kind: ValNull}, nil
}

// --- time ---

func builtinSleepMS(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	n, e := i64Arg(a)
	if e != nil {
		return nil, e
	}
	if n < 0 {
		n = 0
	}
	d := time.Duration(n) * time.Millisecond
	// cap to avoid lockup in bad scripts
	if d > 60*time.Second {
		d = 60 * time.Second
	}
	time.Sleep(d)
	return &Value{Kind: ValNull}, nil
}

func builtinSleepSec(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sleepSec: 1 number")
	}
	f, e := f64Arg(args[0])
	if e != nil {
		return nil, e
	}
	if f < 0 {
		f = 0
	}
	d := time.Duration(f * float64(time.Second))
	if d > 60*time.Second {
		d = 60 * time.Second
	}
	time.Sleep(d)
	return &Value{Kind: ValNull}, nil
}

func builtinTimeMillis(args []*Value) (*Value, error) {
	if err := arg0(args); err != nil {
		return nil, err
	}
	return &Value{Kind: ValInt, I64: time.Now().UnixMilli()}, nil
}

func builtinSeconds(args []*Value) (*Value, error) {
	if err := arg0(args); err != nil {
		return nil, err
	}
	secs := time.Since(processStartTime).Seconds()
	return &Value{Kind: ValFloat, F64: secs}, nil
}

func builtinDeltaTime(args []*Value) (*Value, error) {
	if err := arg0(args); err != nil {
		return nil, err
	}
	now := time.Now()
	dt := now.Sub(lastDeltaCall).Seconds()
	if dt < 0 {
		dt = 0
	}
	lastDeltaCall = now
	return &Value{Kind: ValFloat, F64: dt}, nil
}

// --- assert / type / is_* / range() ---

func builtinAssert(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("assert: condition, optional message")
	}
	if args[0] == nil || !args[0].Truthy() {
		if len(args) > 1 && args[1] != nil && args[1].Kind == ValString {
			return nil, &RuntimeError{Msg: "assert: " + args[1].Str}
		}
		return nil, &RuntimeError{Msg: "assert failed"}
	}
	return &Value{Kind: ValNull}, nil
}

func builtinType(args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type: one value")
	}
	if args[0] == nil {
		return &Value{Kind: ValString, Str: "null"}, nil
	}
	s := "unknown"
	switch args[0].Kind {
	case ValInt:
		s = "int"
	case ValFloat:
		s = "float"
	case ValString:
		s = "string"
	case ValBool:
		s = "bool"
	case ValArray:
		s = "list"
	case ValMap:
		s = "map"
	case ValFunction:
		s = "function"
	case ValStruct:
		s = "struct"
	case ValNull:
		s = "null"
	case ValModule:
		s = "module"
	}
	return &Value{Kind: ValString, Str: s}, nil
}

func builtinIsInt(args []*Value) (*Value, error) { return isKind(ValInt, args) }
func isKind(want ValueKind, args []*Value) (*Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("is*: one value")
	}
	if args[0] == nil {
		return &Value{Kind: ValBool, B: false}, nil
	}
	return &Value{Kind: ValBool, B: args[0].Kind == want}, nil
}
func builtinIsString(args []*Value) (*Value, error)  { return isKind(ValString, args) }
func builtinIsList(args []*Value) (*Value, error)   { return isKind(ValArray, args) }
func builtinIsMap(args []*Value) (*Value, error)    { return isKind(ValMap, args) }
func builtinIsFloat(args []*Value) (*Value, error) { return isKind(ValFloat, args) }
func builtinIsBool(args []*Value) (*Value, error)   { return isKind(ValBool, args) }

// builtinRange produces [start, end) with step: range(end), range(start,end), range(start,end,step)
func builtinRangeFunc(args []*Value) (*Value, error) {
	var start, end, step int64 = 0, 0, 1
	switch len(args) {
	case 1:
		hi, e := i64Arg(args[0])
		if e != nil {
			return nil, e
		}
		start, end, step = 0, hi, 1
	case 2:
		lo, e1 := i64Arg(args[0])
		hi, e2 := i64Arg(args[1])
		if e1 != nil || e2 != nil {
			return nil, fmt.Errorf("range: need integers")
		}
		start, end, step = lo, hi, 1
	case 3:
		lo, e1 := i64Arg(args[0])
		hi, e2 := i64Arg(args[1])
		st, e3 := i64Arg(args[2])
		if e1 != nil || e2 != nil || e3 != nil {
			return nil, fmt.Errorf("range: need integers")
		}
		if st == 0 {
			return nil, fmt.Errorf("range: step 0")
		}
		start, end, step = lo, hi, st
	default:
		return nil, fmt.Errorf("range: 1, 2, or 3 arguments (Python-style, end exclusive)")
	}
	elems := make([]Value, 0, 8)
	if step > 0 {
		for i := start; i < end; i += step {
			elems = append(elems, Value{Kind: ValInt, I64: i})
		}
	} else {
		for i := start; i > end; i += step {
			elems = append(elems, Value{Kind: ValInt, I64: i})
		}
	}
	return &Value{Kind: ValArray, Elems: elems}, nil
}

func builtinDebug(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("debug: at least 1 value")
	}
	for _, a := range args {
		if a == nil {
			fmt.Println("null")
			continue
		}
		fmt.Println(a.String())
	}
	return &Value{Kind: ValNull}, nil
}

// --- string module ---

func builtinStringTrim(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("string.trim: expected (string)")
	}
	return &Value{Kind: ValString, Str: strings.TrimSpace(args[0].Str)}, nil
}

func builtinStringSplit(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValString || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("string.split: expected (string, delimiter)")
	}
	parts := strings.Split(args[0].Str, args[1].Str)
	out := make([]Value, 0, len(parts))
	for _, p := range parts {
		out = append(out, Value{Kind: ValString, Str: p})
	}
	return &Value{Kind: ValArray, Elems: out}, nil
}

func builtinStringJoin(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValArray || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("string.join: expected (list, delimiter)")
	}
	parts := make([]string, 0, len(args[0].Elems))
	for _, e := range args[0].Elems {
		parts = append(parts, e.String())
	}
	return &Value{Kind: ValString, Str: strings.Join(parts, args[1].Str)}, nil
}

func builtinStringReplace(args []*Value) (*Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string.replace: expected (string, old, new)")
	}
	for i := 0; i < 3; i++ {
		if args[i] == nil || args[i].Kind != ValString {
			return nil, fmt.Errorf("string.replace: argument %d must be string", i+1)
		}
	}
	return &Value{Kind: ValString, Str: strings.ReplaceAll(args[0].Str, args[1].Str, args[2].Str)}, nil
}

func builtinStringLower(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("string.lower: expected (string)")
	}
	return &Value{Kind: ValString, Str: strings.ToLower(args[0].Str)}, nil
}

func builtinStringUpper(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("string.upper: expected (string)")
	}
	return &Value{Kind: ValString, Str: strings.ToUpper(args[0].Str)}, nil
}

func builtinStringStartsWith(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValString || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("string.starts_with: expected (string, prefix)")
	}
	return &Value{Kind: ValBool, B: strings.HasPrefix(args[0].Str, args[1].Str)}, nil
}

func builtinStringEndsWith(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValString || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("string.ends_with: expected (string, suffix)")
	}
	return &Value{Kind: ValBool, B: strings.HasSuffix(args[0].Str, args[1].Str)}, nil
}

func builtinStringContains(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValString || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("string.contains: expected (string, substring)")
	}
	return &Value{Kind: ValBool, B: strings.Contains(args[0].Str, args[1].Str)}, nil
}

// --- os module ---

func builtinOSCwd(args []*Value) (*Value, error) {
	if err := arg0(args); err != nil {
		return nil, err
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValString, Str: wd}, nil
}

func builtinOSChdir(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("os.chdir: path must be string")
	}
	if err := os.Chdir(a.Str); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

func builtinOSEnv(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("os.env: variable name must be string")
	}
	if v, ok := os.LookupEnv(a.Str); ok {
		return &Value{Kind: ValString, Str: v}, nil
	}
	return &Value{Kind: ValNull}, nil
}

func builtinOSRun(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("os.run: command must be string")
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", a.Str)
	} else {
		cmd = exec.Command("sh", "-c", a.Str)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("os.run: %w: %s", err, string(out))
	}
	return &Value{Kind: ValString, Str: string(out)}, nil
}

func builtinOSMkdir(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("os.mkdir: path must be string")
	}
	if err := os.MkdirAll(a.Str, 0o755); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

func builtinOSRmdir(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("os.rmdir: path must be string")
	}
	if err := os.RemoveAll(a.Str); err != nil {
		return nil, err
	}
	return &Value{Kind: ValNull}, nil
}

// --- path module ---

func builtinPathJoin(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("path.join: expected at least 2 parts")
	}
	parts := make([]string, 0, len(args))
	for i, a := range args {
		if a == nil || a.Kind != ValString {
			return nil, fmt.Errorf("path.join: arg %d must be string", i+1)
		}
		parts = append(parts, a.Str)
	}
	return &Value{Kind: ValString, Str: filepath.Join(parts...)}, nil
}

func builtinPathBasename(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("path.basename: path must be string")
	}
	return &Value{Kind: ValString, Str: filepath.Base(a.Str)}, nil
}

func builtinPathDirname(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("path.dirname: path must be string")
	}
	return &Value{Kind: ValString, Str: filepath.Dir(a.Str)}, nil
}

func builtinPathExt(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("path.ext: path must be string")
	}
	return &Value{Kind: ValString, Str: filepath.Ext(a.Str)}, nil
}

func builtinPathNormalize(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("path.normalize: path must be string")
	}
	return &Value{Kind: ValString, Str: filepath.Clean(a.Str)}, nil
}

// --- collections module ---

func builtinCollectionsSet(args []*Value) (*Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("collections.set: expected 0 or 1 list")
	}
	out := map[string]Value{}
	if len(args) == 1 {
		if args[0] == nil || args[0].Kind != ValArray {
			return nil, fmt.Errorf("collections.set: argument must be list")
		}
		for _, e := range args[0].Elems {
			out[e.String()] = Value{Kind: ValBool, B: true}
		}
	}
	return &Value{Kind: ValMap, StrMap: out}, nil
}

func builtinCollectionsArrayCtor(args []*Value) (*Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("collections ctor: expected 0 or 1 list")
	}
	if len(args) == 0 {
		return &Value{Kind: ValArray, Elems: []Value{}}, nil
	}
	if args[0] == nil || args[0].Kind != ValArray {
		return nil, fmt.Errorf("collections ctor: argument must be list")
	}
	cpy := make([]Value, len(args[0].Elems))
	copy(cpy, args[0].Elems)
	return &Value{Kind: ValArray, Elems: cpy}, nil
}

func builtinCollectionsPriorityQueue(args []*Value) (*Value, error) {
	v, err := builtinCollectionsArrayCtor(args)
	if err != nil {
		return nil, err
	}
	sort.Slice(v.Elems, func(i, j int) bool {
		a, b := v.Elems[i], v.Elems[j]
		af, ea := f64Arg(&a)
		bf, eb := f64Arg(&b)
		if ea == nil && eb == nil {
			return af < bf
		}
		return a.String() < b.String()
	})
	return v, nil
}

// --- color module ---

func clampColorByte(n int64) int64 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return n
}

func newColorValue(r, g, b, a int64) *Value {
	return &Value{
		Kind: ValMap,
		StrMap: map[string]Value{
			"r": {Kind: ValInt, I64: clampColorByte(r)},
			"g": {Kind: ValInt, I64: clampColorByte(g)},
			"b": {Kind: ValInt, I64: clampColorByte(b)},
			"a": {Kind: ValInt, I64: clampColorByte(a)},
		},
	}
}

func colorFromValue(v *Value) (r, g, b, a int64, err error) {
	if v == nil || v.Kind != ValMap || v.StrMap == nil {
		return 0, 0, 0, 0, fmt.Errorf("expected color map {r,g,b,a}")
	}
	get := func(k string, d int64) (int64, error) {
		if x, ok := v.StrMap[k]; ok {
			return i64Arg(&x)
		}
		return d, nil
	}
	r, err = get("r", 0)
	if err != nil {
		return
	}
	g, err = get("g", 0)
	if err != nil {
		return
	}
	b, err = get("b", 0)
	if err != nil {
		return
	}
	a, err = get("a", 255)
	return
}

func builtinColorRGB(args []*Value) (*Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("color.rgb: expected (r,g,b)")
	}
	r, e1 := i64Arg(args[0])
	g, e2 := i64Arg(args[1])
	b, e3 := i64Arg(args[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, fmt.Errorf("color.rgb: channels must be numbers")
	}
	return newColorValue(r, g, b, 255), nil
}

func builtinColorRGBA(args []*Value) (*Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("color.rgba: expected (r,g,b,a)")
	}
	r, e1 := i64Arg(args[0])
	g, e2 := i64Arg(args[1])
	b, e3 := i64Arg(args[2])
	a, e4 := i64Arg(args[3])
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		return nil, fmt.Errorf("color.rgba: channels must be numbers")
	}
	return newColorValue(r, g, b, a), nil
}

func builtinColorHex(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("color.hex: expected string")
	}
	s := strings.TrimSpace(strings.TrimPrefix(a.Str, "#"))
	if len(s) != 6 && len(s) != 8 {
		return nil, fmt.Errorf("color.hex: expected #RRGGBB or #RRGGBBAA")
	}
	parse := func(x string) (int64, error) {
		n, e := strconv.ParseInt(x, 16, 64)
		return n, e
	}
	r, e1 := parse(s[0:2])
	g, e2 := parse(s[2:4])
	b, e3 := parse(s[4:6])
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, fmt.Errorf("color.hex: invalid hex digits")
	}
	alpha := int64(255)
	if len(s) == 8 {
		aa, e4 := parse(s[6:8])
		if e4 != nil {
			return nil, fmt.Errorf("color.hex: invalid alpha")
		}
		alpha = aa
	}
	return newColorValue(r, g, b, alpha), nil
}

func builtinReadLine(args []*Value) (*Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("readLine expects 0 arguments")
	}
	line, err := stdinReader.ReadString('\n')
	if err != nil && line == "" {
		return nil, err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return &Value{Kind: ValString, Str: line}, nil
}

func builtinParseInt(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	if a.Kind != ValString {
		return nil, fmt.Errorf("parseInt: string")
	}
	s := strings.TrimSpace(a.Str)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValInt, I64: n}, nil
}

func builtinToUpper(args []*Value) (*Value, error) {
	return builtinStringUpper(args)
}

func builtinToLower(args []*Value) (*Value, error) {
	return builtinStringLower(args)
}

func builtinColorLerp(args []*Value) (*Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("color.lerp: expected (a, b, t)")
	}
	r1, g1, b1, a1, e1 := colorFromValue(args[0])
	r2, g2, b2, a2, e2 := colorFromValue(args[1])
	t, e3 := f64Arg(args[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, fmt.Errorf("color.lerp: expected (colorMap, colorMap, number)")
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	lerp := func(x, y int64) int64 {
		return int64(math.Round(float64(x) + (float64(y)-float64(x))*t))
	}
	return newColorValue(lerp(r1, r2), lerp(g1, g2), lerp(b1, b2), lerp(a1, a2)), nil
}

// builtinJoinParts: join("Hello", " ", "World") — concatenate any number of value strings (extended-spec style).
func builtinJoinParts(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: ValString, Str: ""}, nil
	}
	var b strings.Builder
	for _, a := range args {
		if a == nil {
			continue
		}
		b.WriteString(a.String())
	}
	return &Value{Kind: ValString, Str: b.String()}, nil
}

// builtinToNumber: parse a string as int or float (Candy extended spec).
func builtinToNumber(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("toNumber: one string")
	}
	s := strings.TrimSpace(args[0].Str)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	if f == float64(int64(f)) {
		return &Value{Kind: ValInt, I64: int64(f)}, nil
	}
	return &Value{Kind: ValFloat, F64: f}, nil
}

func builtinToInt(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	switch a.Kind {
	case ValInt:
		return &Value{Kind: ValInt, I64: a.I64}, nil
	case ValFloat:
		return &Value{Kind: ValInt, I64: int64(a.F64)}, nil
	case ValBool:
		if a.B {
			return &Value{Kind: ValInt, I64: 1}, nil
		}
		return &Value{Kind: ValInt, I64: 0}, nil
	case ValString:
		i, e := strconv.ParseInt(strings.TrimSpace(a.Str), 10, 64)
		if e != nil {
			f, e2 := strconv.ParseFloat(strings.TrimSpace(a.Str), 64)
			if e2 != nil {
				return nil, e
			}
			return &Value{Kind: ValInt, I64: int64(f)}, nil
		}
		return &Value{Kind: ValInt, I64: i}, nil
	default:
		return nil, fmt.Errorf("toInt unsupported type")
	}
}

func builtinToFloat(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	switch a.Kind {
	case ValFloat:
		return &Value{Kind: ValFloat, F64: a.F64}, nil
	case ValInt:
		return &Value{Kind: ValFloat, F64: float64(a.I64)}, nil
	case ValBool:
		if a.B {
			return &Value{Kind: ValFloat, F64: 1}, nil
		}
		return &Value{Kind: ValFloat, F64: 0}, nil
	case ValString:
		f, e := strconv.ParseFloat(strings.TrimSpace(a.Str), 64)
		if e != nil {
			return nil, e
		}
		return &Value{Kind: ValFloat, F64: f}, nil
	default:
		return nil, fmt.Errorf("toFloat unsupported type")
	}
}

func builtinToString(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	return &Value{Kind: ValString, Str: a.String()}, nil
}

func builtinToBool(args []*Value) (*Value, error) {
	a, err := arg1(args)
	if err != nil {
		return nil, err
	}
	switch a.Kind {
	case ValBool:
		return &Value{Kind: ValBool, B: a.B}, nil
	case ValInt:
		return &Value{Kind: ValBool, B: a.I64 != 0}, nil
	case ValFloat:
		return &Value{Kind: ValBool, B: a.F64 != 0}, nil
	case ValString:
		s := strings.ToLower(strings.TrimSpace(a.Str))
		if s == "" || s == "0" || s == "false" || s == "no" || s == "off" || s == "null" {
			return &Value{Kind: ValBool, B: false}, nil
		}
		return &Value{Kind: ValBool, B: true}, nil
	default:
		return &Value{Kind: ValBool, B: a.Truthy()}, nil
	}
}

func builtinExitProg(args []*Value) (*Value, error) {
	if len(args) == 0 {
		exitProgram(0)
	}
	code, err := i64Arg(args[0])
	if err != nil {
		return nil, err
	}
	exitProgram(int(code))
	return &Value{Kind: ValNull}, nil
}

func builtinSave(args []*Value) (*Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("save expects key, value")
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("save key must be string")
	}
	globalSaveData[args[0].Str] = valueToValue(args[1])
	return &Value{Kind: ValNull}, nil
}

func builtinLoad(args []*Value) (*Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("load expects key, [default]")
	}
	if args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("load key must be string")
	}
	if v, ok := globalSaveData[args[0].Str]; ok {
		return ptrVal(v), nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return &Value{Kind: ValNull}, nil
}

func init() {
	// merge into main Builtins map
	x := map[string]func(args []*Value) (*Value, error){
		"read_file": builtinReadFile, "readFile": builtinReadFile, "readfile": builtinReadFile,
		"write_file": builtinWriteFile, "writeFile": builtinWriteFile, "writefile": builtinWriteFile,
		"read_lines": builtinReadLines, "readLines": builtinReadLines,
		"file_exists": builtinFileExists, "fileExists": builtinFileExists, "path_exists": builtinFileExists,
		"delete_file":  builtinDeleteFile, "deleteFile": builtinDeleteFile,
		"append_file":  builtinAppendFile, "appendFile": builtinAppendFile,
		"list_files":   builtinListFiles, "listFiles": builtinListFiles, "list_dir": builtinListFiles, "listDir": builtinListFiles,
		"json_stringify":  builtinJsonStringify, "jsonStringify": builtinJsonStringify,
		"json_parse":   builtinJsonParse, "jsonParse": builtinJsonParse,
		"load_json":    builtinLoadJSON, "loadJson": builtinLoadJSON,
		"save_json":    builtinSaveJSON, "saveJson": builtinSaveJSON,
		"sqrt":         builtinSqrt, "abs": builtinAbsF,
		"pow":  builtinPow, "floor": builtinFloor, "ceil":  builtinCeil, "round":  builtinRound,
		"sin": builtinSin, "cos": builtinCos, "tan": builtinTan, "lerp": builtinLerp,
		"min": builtinMin, "max": builtinMax, "clamp": builtinClamp,
		"random":       builtinRandomInt, "randomInt":  builtinRandomInt, "rand_int": builtinRandomInt,
		"random_float": builtinRandomFloat, "randomFloat": builtinRandomFloat, "rand_float": builtinRandomFloat,
		"choose":  builtinChoose,
		"sleep":     builtinSleepMS, "sleepMs":  builtinSleepMS, "wait_ms": builtinSleepMS,
		"sleep_sec": builtinSleepSec, "sleepSec": builtinSleepSec, "wait_sec": builtinSleepSec,
		"time_millis": builtinTimeMillis, "timeMillis": builtinTimeMillis, "time_ms": builtinTimeMillis,
		"seconds": builtinSeconds, "deltaTime": builtinDeltaTime, "deltatime": builtinDeltaTime,
		"assert": builtinAssert, "type": builtinType, "typeof": builtinType,
		"is_int":  builtinIsInt, "isInt":  builtinIsInt, "is_string": builtinIsString, "isString": builtinIsString,
		"is_list": builtinIsList, "isList": builtinIsList, "is_map":  builtinIsMap, "isMap": builtinIsMap, "is_float": builtinIsFloat, "is_float64": builtinIsFloat,
		"is_bool":  builtinIsBool, "isBool":  builtinIsBool,
		"range":  builtinRangeFunc, "array_range": builtinRangeFunc, "iota_range": builtinRangeFunc,
		"debug":  builtinDebug,
		"printf": builtinPrint,
		"shuffle": builtinRandomShuffle,
		"pick": builtinChoose,
		"read_line": builtinReadLine, "readline": builtinReadLine, "readLine": builtinReadLine,
		"parse_int": builtinParseInt, "parseInt": builtinParseInt, "parseint": builtinParseInt,
		"toupper": builtinToUpper, "toUpper": builtinToUpper,
		"tolower": builtinToLower, "toLower": builtinToLower,
		"rand":   builtinRandomInt,
		"join":   builtinJoinParts,
		"split":  builtinStringSplit,
		"replace":   builtinStringReplace,
		"contains":  builtinStringContains,
		"toNumber":  builtinToNumber, "tonumber": builtinToNumber, "to_number": builtinToNumber,
		"toInt": builtinToInt, "to_int": builtinToInt, "int": builtinToInt,
		"toFloat": builtinToFloat, "to_float": builtinToFloat, "float": builtinToFloat,
		"toString": builtinToString, "to_string": builtinToString, "string": builtinToString,
		"toBool": builtinToBool, "to_bool": builtinToBool, "bool": builtinToBool,
		"wait": builtinSleepSec, "exit": builtinExitProg,
		"save": builtinSave, "load": builtinLoad,
	}
	// use existing read/write names from builtin — fix references: use builtinReadFile
	for k, v := range x {
		Builtins[strings.ToLower(k)] = v
	}
}
