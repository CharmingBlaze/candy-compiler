package candy_evaluator

import (
	"candy/candy_ast"
	"fmt"
	"math"
	"strings"
)

func init() {
	RegisterBuiltin("vec2", builtinVecCtor(2))
	RegisterBuiltin("vec3", builtinVecCtor(3))
	RegisterBuiltin("vec4", builtinVecCtor(4))
	RegisterBuiltin("format", builtinFormat)
	RegisterBuiltin("enumerate", builtinEnumerate)
	RegisterBuiltin("AABB", builtinAABB)
	RegisterBuiltin("Sphere", builtinSphere)
	RegisterBuiltin("Ray", builtinRay)
	RegisterBuiltin("Box", builtinBox)
	RegisterBuiltin("PhysicsWorld", builtinPhysicsWorld)
	RegisterBuiltin("InputMap", builtinInputMap)
	RegisterBuiltin("OrbitCamera", builtinOrbitCamera)
	RegisterBuiltin("FirstPersonCamera", builtinFirstPersonCamera)
	RegisterBuiltin("CharacterController", builtinCharacterController)
	RegisterBuiltin("EntityList", builtinEntityList)
	RegisterBuiltin("UILayout", builtinUILayout)
	RegisterBuiltin("HUD", builtinHUD)
	RegisterBuiltin("StateMachine", builtinStateMachine)
	RegisterBuiltin("Tween", builtinTween)
	RegisterBuiltin("Transform", builtinTransform)
	RegisterBuiltin("drawAll", builtinDrawAll)
	RegisterBuiltin("saveState", builtinSaveState)
	RegisterBuiltin("restoreState", builtinRestoreState)
	RegisterBuiltin("cloneState", builtinCloneState)
	RegisterBuiltin("gameLoop", builtinGameLoop)
	RegisterBuiltin("entity", builtinEntity)
	RegisterBuiltin("addComponent", builtinAddComponent)
	RegisterBuiltin("getComponent", builtinGetComponent)
}

func builtinVecCtor(dim int) func(args []*Value) (*Value, error) {
	return func(args []*Value) (*Value, error) {
		if len(args) != dim {
			return nil, fmt.Errorf("vec%d expects %d args", dim, dim)
		}
		out := make([]float64, dim)
		for i, a := range args {
			x, err := f64Arg(a)
			if err != nil {
				return nil, fmt.Errorf("vec%d arg %d must be numeric", dim, i+1)
			}
			out[i] = x
		}
		return &Value{Kind: ValVec, Vec: out}, nil
	}
}

func vecProps(v *Value, name string) *Value {
	if v == nil || v.Kind != ValVec {
		return &Value{Kind: ValNull}
	}
	ln := strings.ToLower(name)
	switch ln {
	case "x":
		if len(v.Vec) > 0 {
			return &Value{Kind: ValFloat, F64: v.Vec[0]}
		}
	case "y":
		if len(v.Vec) > 1 {
			return &Value{Kind: ValFloat, F64: v.Vec[1]}
		}
	case "z":
		if len(v.Vec) > 2 {
			return &Value{Kind: ValFloat, F64: v.Vec[2]}
		}
	case "w":
		if len(v.Vec) > 3 {
			return &Value{Kind: ValFloat, F64: v.Vec[3]}
		}
	case "xy":
		if len(v.Vec) > 1 {
			return &Value{Kind: ValVec, Vec: []float64{v.Vec[0], v.Vec[1]}}
		}
	case "xz":
		if len(v.Vec) > 2 {
			return &Value{Kind: ValVec, Vec: []float64{v.Vec[0], v.Vec[2]}}
		}
	case "yz":
		if len(v.Vec) > 2 {
			return &Value{Kind: ValVec, Vec: []float64{v.Vec[1], v.Vec[2]}}
		}
	}
	return nil
}

func callVecMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(name) {
	case "length":
		return &Value{Kind: ValFloat, F64: vecLen(recv.Vec)}, nil
	case "normalize":
		l := vecLen(recv.Vec)
		if l <= 1e-12 {
			return &Value{Kind: ValVec, Vec: append([]float64(nil), recv.Vec...)}, nil
		}
		out := make([]float64, len(recv.Vec))
		for i := range recv.Vec {
			out[i] = recv.Vec[i] / l
		}
		return &Value{Kind: ValVec, Vec: out}, nil
	case "dot":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValVec || len(args[0].Vec) != len(recv.Vec) {
			return nil, fmt.Errorf("dot expects one vector with same dimension")
		}
		s := 0.0
		for i := range recv.Vec {
			s += recv.Vec[i] * args[0].Vec[i]
		}
		return &Value{Kind: ValFloat, F64: s}, nil
	case "cross":
		if len(recv.Vec) != 3 || len(args) != 1 || args[0] == nil || args[0].Kind != ValVec || len(args[0].Vec) != 3 {
			return nil, fmt.Errorf("cross expects vec3 argument")
		}
		a := recv.Vec
		b := args[0].Vec
		return &Value{Kind: ValVec, Vec: []float64{
			a[1]*b[2] - a[2]*b[1],
			a[2]*b[0] - a[0]*b[2],
			a[0]*b[1] - a[1]*b[0],
		}}, nil
	case "distance":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValVec || len(args[0].Vec) != len(recv.Vec) {
			return nil, fmt.Errorf("distance expects one vector with same dimension")
		}
		sum := 0.0
		for i := range recv.Vec {
			d := recv.Vec[i] - args[0].Vec[i]
			sum += d * d
		}
		return &Value{Kind: ValFloat, F64: math.Sqrt(sum)}, nil
	case "rotate":
		if len(recv.Vec) != 2 || len(args) != 1 {
			return nil, fmt.Errorf("rotate expects vec2.rotate(angle)")
		}
		a, err := f64Arg(args[0])
		if err != nil {
			return nil, err
		}
		c := math.Cos(a)
		s := math.Sin(a)
		x := recv.Vec[0]
		y := recv.Vec[1]
		return &Value{Kind: ValVec, Vec: []float64{x*c - y*s, x*s + y*c}}, nil
	default:
		return nil, fmt.Errorf("vector has no method %s", name)
	}
}

func vecLen(v []float64) float64 {
	sum := 0.0
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

func evalVecInfix(op string, l, r *Value, node candy_ast.Node) (*Value, error) {
	switch op {
	case "==", "!=":
		eq := valueEqual(l, r)
		if op == "==" {
			return &Value{Kind: ValBool, B: eq}, nil
		}
		return &Value{Kind: ValBool, B: !eq}, nil
	}
	// vec +/-/*// vec|scalar support.
	if l.Kind == ValVec && r.Kind == ValVec {
		if len(l.Vec) != len(r.Vec) {
			return nil, newError("vector dimension mismatch", node)
		}
		out := make([]float64, len(l.Vec))
		for i := range out {
			switch op {
			case "+":
				out[i] = l.Vec[i] + r.Vec[i]
			case "-":
				out[i] = l.Vec[i] - r.Vec[i]
			case "*":
				out[i] = l.Vec[i] * r.Vec[i]
			case "/":
				out[i] = l.Vec[i] / r.Vec[i]
			default:
				return nil, newError("unsupported vector operator "+op, node)
			}
		}
		return &Value{Kind: ValVec, Vec: out}, nil
	}
	if l.Kind == ValVec && isNumeric(r) {
		s := r.F64
		if r.Kind == ValInt {
			s = float64(r.I64)
		}
		out := make([]float64, len(l.Vec))
		for i := range out {
			switch op {
			case "*":
				out[i] = l.Vec[i] * s
			case "/":
				out[i] = l.Vec[i] / s
			case "+":
				out[i] = l.Vec[i] + s
			case "-":
				out[i] = l.Vec[i] - s
			default:
				return nil, newError("unsupported vector operator "+op, node)
			}
		}
		return &Value{Kind: ValVec, Vec: out}, nil
	}
	if r.Kind == ValVec && isNumeric(l) {
		return evalVecInfix(op, r, l, node)
	}
	return nil, newError("invalid vector operation", node)
}

func builtinFormat(args []*Value) (*Value, error) {
	if len(args) == 0 || args[0] == nil || args[0].Kind != ValString {
		return nil, fmt.Errorf("format expects template string first")
	}
	s := args[0].Str
	for i := 1; i < len(args); i++ {
		as := args[i].String()
		s = strings.Replace(s, "{}", as, 1)
	}
	return &Value{Kind: ValString, Str: s}, nil
}

func builtinEnumerate(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil || args[0].Kind != ValArray {
		return nil, fmt.Errorf("enumerate expects one array")
	}
	out := make([]Value, 0, len(args[0].Elems))
	for i, el := range args[0].Elems {
		pair := Value{Kind: ValArray, Elems: []Value{{Kind: ValInt, I64: int64(i)}, el}}
		out = append(out, pair)
	}
	return &Value{Kind: ValArray, Elems: out}, nil
}

func builtinAABB(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValVec || args[1].Kind != ValVec {
		return nil, fmt.Errorf("AABB expects (centerVec3, halfExtentsVec3)")
	}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":        {Kind: ValString, Str: "AABB"},
		"center":      *args[0],
		"halfExtents": *args[1],
	}}, nil
}

func builtinBox(args []*Value) (*Value, error) {
	// Box(centerVec3, sizeVec3) -> stores as AABB-style center+halfExtents.
	if len(args) != 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValVec || args[1].Kind != ValVec || len(args[0].Vec) != 3 || len(args[1].Vec) != 3 {
		return nil, fmt.Errorf("Box expects (centerVec3, sizeVec3)")
	}
	size := args[1].Vec
	half := []float64{size[0] * 0.5, size[1] * 0.5, size[2] * 0.5}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":        {Kind: ValString, Str: "AABB"},
		"center":      *args[0],
		"size":        *args[1],
		"halfExtents": {Kind: ValVec, Vec: half},
	}}, nil
}

func builtinSphere(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValVec || !isNumeric(args[1]) {
		return nil, fmt.Errorf("Sphere expects (centerVec3, radius)")
	}
	r := args[1].F64
	if args[1].Kind == ValInt {
		r = float64(args[1].I64)
	}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":   {Kind: ValString, Str: "Sphere"},
		"center": *args[0],
		"radius": {Kind: ValFloat, F64: r},
	}}, nil
}

func builtinRay(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValVec || args[1].Kind != ValVec {
		return nil, fmt.Errorf("Ray expects (originVec3, directionVec3)")
	}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":      {Kind: ValString, Str: "Ray"},
		"origin":    *args[0],
		"direction": *args[1],
	}}, nil
}

func builtinGameLoop(args []*Value) (*Value, error) {
	// gameLoop(fps?, updateFn?, drawFn?, maxFrames?)
	fps := int64(60)
	argIx := 0
	if len(args) > 0 && args[0] != nil && args[0].Kind == ValInt {
		fps = args[0].I64
		argIx = 1
	}
	var updateFn *Value
	var drawFn *Value
	if argIx < len(args) {
		updateFn = args[argIx]
	}
	if argIx+1 < len(args) {
		drawFn = args[argIx+1]
	}
	maxFrames := int64(1)
	if argIx+2 < len(args) && args[argIx+2] != nil && args[argIx+2].Kind == ValInt {
		maxFrames = args[argIx+2].I64
	}
	dt := 1.0 / float64(fps)
	if dt <= 0 {
		dt = 0.016
	}
	if dt > 0.05 {
		dt = 0.05
	}
	for i := int64(0); i < maxFrames; i++ {
		if updateFn != nil {
			_, _ = InvokeCallable(updateFn, []*Value{{Kind: ValFloat, F64: dt}})
		}
		if drawFn != nil {
			_, _ = InvokeCallable(drawFn, nil)
		}
	}
	return &Value{Kind: ValNull}, nil
}

func builtinEntity(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"components": {Kind: ValMap, StrMap: map[string]Value{}},
	}}, nil
}

func builtinAddComponent(args []*Value) (*Value, error) {
	if len(args) != 3 || args[0] == nil || args[0].Kind != ValMap || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("addComponent(entity, name, component)")
	}
	e := args[0]
	comps, ok := e.StrMap["components"]
	if !ok || comps.Kind != ValMap || comps.StrMap == nil {
		comps = Value{Kind: ValMap, StrMap: map[string]Value{}}
	}
	comps.StrMap[args[1].Str] = *args[2]
	e.StrMap["components"] = comps
	return e, nil
}

func builtinGetComponent(args []*Value) (*Value, error) {
	if len(args) != 2 || args[0] == nil || args[0].Kind != ValMap || args[1] == nil || args[1].Kind != ValString {
		return nil, fmt.Errorf("getComponent(entity, name)")
	}
	comps, ok := args[0].StrMap["components"]
	if !ok || comps.Kind != ValMap || comps.StrMap == nil {
		return &Value{Kind: ValNull}, nil
	}
	if v, ok := comps.StrMap[args[1].Str]; ok {
		return &v, nil
	}
	return &Value{Kind: ValNull}, nil
}

func builtinPhysicsWorld(args []*Value) (*Value, error) {
	g := []float64{0, -9.8, 0}
	if len(args) >= 1 && args[0] != nil && args[0].Kind == ValVec && len(args[0].Vec) == 3 {
		g = append([]float64(nil), args[0].Vec...)
	}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":    {Kind: ValString, Str: "PhysicsWorld"},
		"gravity": {Kind: ValVec, Vec: g},
		"static":  {Kind: ValArray, Elems: []Value{}},
		"dynamic": {Kind: ValArray, Elems: []Value{}},
	}}, nil
}

func builtinInputMap(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":    {Kind: ValString, Str: "InputMap"},
		"actions": {Kind: ValMap, StrMap: map[string]Value{}},
		"axes2d":  {Kind: ValMap, StrMap: map[string]Value{}},
		"axes1d":  {Kind: ValMap, StrMap: map[string]Value{}},
		"contexts": {Kind: ValArray, Elems: []Value{}},
	}}, nil
}

func builtinOrbitCamera(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":        {Kind: ValString, Str: "OrbitCamera"},
		"target":      {Kind: ValVec, Vec: []float64{0, 0, 0}},
		"distance":    {Kind: ValFloat, F64: 11.0},
		"yaw":         {Kind: ValFloat, F64: 0.35},
		"pitch":       {Kind: ValFloat, F64: 0.35},
		"sensitivity": {Kind: ValFloat, F64: 0.007},
		"zoomSpeed":   {Kind: ValFloat, F64: 1.4},
		"minPitch":    {Kind: ValFloat, F64: -1.1},
		"maxPitch":    {Kind: ValFloat, F64: 1.1},
		"minZoom":     {Kind: ValFloat, F64: 5.0},
		"maxZoom":     {Kind: ValFloat, F64: 18.0},
	}}, nil
}

func builtinFirstPersonCamera(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":        {Kind: ValString, Str: "FirstPersonCamera"},
		"position":    {Kind: ValVec, Vec: []float64{0, 2, 0}},
		"yaw":         {Kind: ValFloat, F64: 0},
		"pitch":       {Kind: ValFloat, F64: 0},
		"sensitivity": {Kind: ValFloat, F64: 0.007},
	}}, nil
}

func builtinCharacterController(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":         {Kind: ValString, Str: "CharacterController"},
		"acceleration": {Kind: ValFloat, F64: 38.0},
		"drag":         {Kind: ValFloat, F64: 8.0},
		"maxSpeed":     {Kind: ValFloat, F64: 12.0},
		"jumpPower":    {Kind: ValFloat, F64: 9.3},
		"gravity":      {Kind: ValFloat, F64: 28.0},
		"onGround":     {Kind: ValBool, B: false},
	}}, nil
}

func builtinEntityList(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":     {Kind: ValString, Str: "EntityList"},
		"entities": {Kind: ValArray, Elems: []Value{}},
	}}, nil
}

func builtinUILayout(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":    {Kind: ValString, Str: "UILayout"},
		"x":       {Kind: ValInt, I64: 14},
		"y":       {Kind: ValInt, I64: 14},
		"spacing": {Kind: ValInt, I64: 8},
	}}, nil
}

func builtinHUD(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":    {Kind: ValString, Str: "HUD"},
		"padding": {Kind: ValInt, I64: 14},
		"spacing": {Kind: ValInt, I64: 8},
	}}, nil
}

func builtinStateMachine(args []*Value) (*Value, error) {
	state := "playing"
	if len(args) > 0 && args[0] != nil && args[0].Kind == ValString {
		state = args[0].Str
	}
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":    {Kind: ValString, Str: "StateMachine"},
		"current": {Kind: ValString, Str: state},
		"states":  {Kind: ValMap, StrMap: map[string]Value{}},
	}}, nil
}

func builtinTween(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":      {Kind: ValString, Str: "Tween"},
		"from":      {Kind: ValFloat, F64: 0},
		"to":        {Kind: ValFloat, F64: 1},
		"duration":  {Kind: ValFloat, F64: 1},
		"time":      {Kind: ValFloat, F64: 0},
		"pingpong":  {Kind: ValBool, B: false},
		"loop":      {Kind: ValBool, B: false},
		"direction": {Kind: ValFloat, F64: 1},
	}}, nil
}

func builtinTransform(args []*Value) (*Value, error) {
	return &Value{Kind: ValMap, StrMap: map[string]Value{
		"type":     {Kind: ValString, Str: "Transform"},
		"rotation": {Kind: ValFloat, F64: 0},
	}}, nil
}

func builtinDrawAll(args []*Value) (*Value, error) {
	if len(args) == 0 || args[0] == nil || args[0].Kind != ValArray {
		return nil, fmt.Errorf("drawAll expects an array")
	}
	for i := range args[0].Elems {
		el := args[0].Elems[i]
		if el.Kind == ValMap {
			v := el
			_, _ = callMapMethod(&v, "draw", nil, &Env{Store: map[string]*Value{}})
		}
	}
	return &Value{Kind: ValNull}, nil
}

func builtinCloneState(args []*Value) (*Value, error) {
	if len(args) != 1 || args[0] == nil {
		return nil, fmt.Errorf("cloneState expects one value")
	}
	c := deepCloneValue(args[0])
	return c, nil
}

func builtinSaveState(args []*Value) (*Value, error) {
	return builtinCloneState(args)
}

func builtinRestoreState(args []*Value) (*Value, error) {
	return builtinCloneState(args)
}

func deepCloneValue(v *Value) *Value {
	if v == nil {
		return &Value{Kind: ValNull}
	}
	out := *v
	if v.Kind == ValArray {
		out.Elems = make([]Value, len(v.Elems))
		for i := range v.Elems {
			c := deepCloneValue(&v.Elems[i])
			out.Elems[i] = *c
		}
	}
	if v.Kind == ValMap {
		out.StrMap = map[string]Value{}
		for k, vv := range v.StrMap {
			c := deepCloneValue(&vv)
			out.StrMap[k] = *c
		}
	}
	if v.Kind == ValVec {
		out.Vec = append([]float64(nil), v.Vec...)
	}
	return &out
}
