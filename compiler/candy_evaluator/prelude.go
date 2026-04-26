package candy_evaluator

import (
	"math"
	"strings"
)

// registerPrelude injects stdlib module objects and top-level constants (PI, E).
// Called at the start of Eval() so each run has a known surface.
func registerPrelude(e *Env) {
	if e == nil {
		return
	}
	if e.Store == nil {
		e.Store = make(map[string]*Value)
	}
	pi := &Value{Kind: ValFloat, F64: math.Pi}
	ee := &Value{Kind: ValFloat, F64: math.E}
	inf := &Value{Kind: ValFloat, F64: math.Inf(1)}
	e.Set("PI", pi)
	e.Set("E", ee)
	e.Set("pi", pi)
	e.Set("e", ee)
	e.Set("infinity", inf)
	e.Set("Infinity", inf)

	// Keys
	keys := map[string]int{
		"KEY_W": 87, "KEY_A": 65, "KEY_S": 83, "KEY_D": 68,
		"KEY_SPACE": 32, "KEY_ESCAPE": 256, "KEY_ENTER": 257,
		"KEY_LEFT": 263, "KEY_RIGHT": 262, "KEY_UP": 265, "KEY_DOWN": 264,
		// Common gamepad aliases for scripts that use keyboard+pad fallbacks.
		// InputMap currently consumes integer codes via isKeyDown/isKeyPressed.
		"GAMEPAD_BUTTON_A": 32,
		"GAMEPAD_BUTTON_B": 257,
		"GAMEPAD_BUTTON_X": 65,
		"GAMEPAD_BUTTON_Y": 87,
	}
	for k, v := range keys {
		e.Set(k, &Value{Kind: ValInt, I64: int64(v)})
		e.Set(strings.ToLower(k), &Value{Kind: ValInt, I64: int64(v)})
	}
	e.Set("left", &Value{Kind: ValString, Str: "LEFT"})
	e.Set("right", &Value{Kind: ValString, Str: "RIGHT"})
	e.Set("up", &Value{Kind: ValString, Str: "UP"})
	e.Set("down", &Value{Kind: ValString, Str: "DOWN"})
	e.Set("space", &Value{Kind: ValString, Str: "SPACE"})
	e.Set("key_esc", &Value{Kind: ValString, Str: "ESC"})

	// Colors
	colors := []string{"white", "black", "gray", "red", "green", "blue", "yellow", "gold", "orange", "pink", "purple", "skyblue", "brown", "darkgray", "maroon"}
	colorConsts := make(map[string]*Value)
	for _, c := range colors {
		val := &Value{Kind: ValString, Str: c}
		e.Set(c, val)
		colorName := "COLOR_" + strings.ToUpper(c)
		e.Set(colorName, val)
		e.Set(strings.ToLower(colorName), val)
		colorConsts[strings.ToUpper(c)] = val
	}

	mathFn := map[string]func(args []*Value) (*Value, error){
		"sqrt":  builtinSqrt,
		"pow":   builtinPow,
		"abs":   builtinAbsF,
		"floor": builtinFloor,
		"ceil":  builtinCeil,
		"round": builtinRound,
		"sin":   builtinSin,
		"cos":   builtinCos,
		"tan":   builtinTan,
		"min":   builtinMin,
		"max":   builtinMax,
		"clamp": builtinClamp,
		"isInfinite": builtinIsInfinite,
	}
	mathC := map[string]*Value{
		"PI":  pi,
		"E":   ee,
		"pi":  pi,
		"e":   ee,
		"Inf": inf,
		"infinity": inf,
	}
	e.Set("math", newModule("math", mathFn, mathC))

	fileFns := map[string]func(args []*Value) (*Value, error){
		"read":        builtinReadFile,
		"write":       builtinWriteFile,
		"read_file":   builtinReadFile,
		"readFile":    builtinReadFile,
		"write_file":  builtinWriteFile,
		"writeFile":   builtinWriteFile,
		"read_lines":  builtinReadLines,
		"readLines":   builtinReadLines,
		"exists":      builtinFileExists,
		"file_exists": builtinFileExists,
		"fileExists":  builtinFileExists,
		"delete":      builtinDeleteFile,
		"remove":      builtinDeleteFile,
		"delete_file": builtinDeleteFile,
		"list":        builtinListFiles,
		"list_dir":    builtinListFiles,
		"listDir":     builtinListFiles,
		"list_files":  builtinListFiles,
		"listFiles":   builtinListFiles,
	}
	fileV := newModule("file", fileFns, nil)
	e.Set("file", fileV)
	e.Set("fs", fileV)

	jsonFn := map[string]func(args []*Value) (*Value, error){
		"parse":     builtinJsonParse,
		"decode":    builtinJsonParse,
		"stringify": builtinJsonStringify,
		"encode":    builtinJsonStringify,
		"load":      builtinLoadJSON,
		"load_file": builtinLoadJSON,
		"loadFile":  builtinLoadJSON,
		"save":      builtinSaveJSON,
		"save_file": builtinSaveJSON,
		"saveFile":  builtinSaveJSON,
	}
	e.Set("json", newModule("json", jsonFn, nil))

	rndFn := map[string]func(args []*Value) (*Value, error){
		"int":   builtinRandomInt,
		"float": builtinRandomFloat,
		"choice": builtinChoose,
		"sample": builtinChoose,
		"pick":   builtinChoose,
		"shuffle": builtinRandomShuffle,
		"seed":   builtinRngSeed,
	}
	e.Set("random", newModule("random", rndFn, nil))
	e.Set("rand", e.Store["random"])

	timeFn := map[string]func(args []*Value) (*Value, error){
		"millis":    builtinTimeMillis,
		"ms":        builtinTimeMillis,
		"now":       builtinTimeMillis,
		"now_ms":    builtinTimeMillis,
		"nowMs":     builtinTimeMillis,
		"sleep":     builtinSleepMS,
		"sleep_ms":  builtinSleepMS,
		"sleepMs":   builtinSleepMS,
		"sleep_sec": builtinSleepSec,
		"sleepSec":  builtinSleepSec,
		"wait":      builtinSleepMS,
	}
	e.Set("time", newModule("time", timeFn, nil))

	stringFn := map[string]func(args []*Value) (*Value, error){
		"trim":        builtinStringTrim,
		"split":       builtinStringSplit,
		"join":        builtinStringJoin,
		"replace":     builtinStringReplace,
		"lower":       builtinStringLower,
		"upper":       builtinStringUpper,
		"starts_with": builtinStringStartsWith,
		"ends_with":   builtinStringEndsWith,
		"contains":    builtinStringContains,
	}
	e.Set("string", newModule("string", stringFn, nil))
	e.Set("upper", &Value{Kind: ValFunction, Builtin: builtinStringUpper})
	e.Set("lower", &Value{Kind: ValFunction, Builtin: builtinStringLower})

	osFn := map[string]func(args []*Value) (*Value, error){
		"cwd":   builtinOSCwd,
		"chdir": builtinOSChdir,
		"env":   builtinOSEnv,
		"run":   builtinOSRun,
		"mkdir": builtinOSMkdir,
		"rmdir": builtinOSRmdir,
	}
	e.Set("os", newModule("os", osFn, nil))

	pathFn := map[string]func(args []*Value) (*Value, error){
		"join":      builtinPathJoin,
		"basename":  builtinPathBasename,
		"dirname":   builtinPathDirname,
		"ext":       builtinPathExt,
		"normalize": builtinPathNormalize,
	}
	e.Set("path", newModule("path", pathFn, nil))

	collectionsFn := map[string]func(args []*Value) (*Value, error){
		"set":            builtinCollectionsSet,
		"queue":          builtinCollectionsArrayCtor,
		"stack":          builtinCollectionsArrayCtor,
		"deque":          builtinCollectionsArrayCtor,
		"priority_queue": builtinCollectionsPriorityQueue,
	}
	e.Set("collections", newModule("collections", collectionsFn, nil))

	colorFn := map[string]func(args []*Value) (*Value, error){
		"rgb":    builtinColorRGB,
		"rgba":   builtinColorRGBA,
		"hex":    builtinColorHex,
		"lerp":   builtinColorLerp,
		"random": builtinColorRandom,
	}
	e.Set("color", newModule("color", colorFn, colorConsts))
	e.Set("Color", e.Store["color"]) // Alias for documentation parity

	registerENetModule(e)
}

func newModule(name string, fns map[string]func(args []*Value) (*Value, error), consts map[string]*Value) *Value {
	if fns == nil {
		fns = make(map[string]func(args []*Value) (*Value, error))
	}
	return &Value{
		Kind: ValModule,
		Mod: &moduleVal{
			Name:   name,
			Fns:    fns,
			Consts: consts,
		},
	}
}
