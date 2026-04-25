package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Struct constructors ----

func builtinVec2(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vec2", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("vec2", args, 0)
	y, _ := getArgFloat("vec2", args, 1)
	return vec2ToMap(rl.NewVector2(float32(x), float32(y))), nil
}

func builtinVec3(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vec3", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("vec3", args, 0)
	y, _ := getArgFloat("vec3", args, 1)
	z, _ := getArgFloat("vec3", args, 2)
	return vec3ToMap(rl.NewVector3(float32(x), float32(y), float32(z))), nil
}

func builtinVec4(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vec4", args, 4); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("vec4", args, 0)
	y, _ := getArgFloat("vec4", args, 1)
	z, _ := getArgFloat("vec4", args, 2)
	w, _ := getArgFloat("vec4", args, 3)
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: x},
		"y": {Kind: candy_evaluator.ValFloat, F64: y},
		"z": {Kind: candy_evaluator.ValFloat, F64: z},
		"w": {Kind: candy_evaluator.ValFloat, F64: w},
	}), nil
}

func builtinRect(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("rect", args, 4); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("rect", args, 0)
	y, _ := getArgFloat("rect", args, 1)
	w, _ := getArgFloat("rect", args, 2)
	h, _ := getArgFloat("rect", args, 3)
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: x},
		"y":      {Kind: candy_evaluator.ValFloat, F64: y},
		"width":  {Kind: candy_evaluator.ValFloat, F64: w},
		"height": {Kind: candy_evaluator.ValFloat, F64: h},
	}), nil
}

func builtinColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) == 1 && args[0] != nil && args[0].Kind == candy_evaluator.ValString {
		c := colorFrom(args[0].Str)
		return colorToMap(c), nil
	}
	if len(args) < 3 {
		return nil, fmt.Errorf("color expects r, g, b, [a] or a color name string")
	}
	r, _ := argInt("color", args, 0)
	g, _ := argInt("color", args, 1)
	b, _ := argInt("color", args, 2)
	a := int64(255)
	if len(args) > 3 {
		a, _ = argInt("color", args, 3)
	}
	return colorToMap(rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a))), nil
}

func builtinBoundingBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) == 2 {
		// vec3, vec3 form
		min, err := argVector3("boundingBox", args, 0)
		if err != nil {
			return nil, err
		}
		max, err2 := argVector3("boundingBox", args, 1)
		if err2 != nil {
			return nil, err2
		}
		return boundingBoxToMap(rl.NewBoundingBox(min, max)), nil
	}
	if err := expectArgs("boundingBox", args, 6); err != nil {
		return nil, err
	}
	minX, _ := getArgFloat("boundingBox", args, 0)
	minY, _ := getArgFloat("boundingBox", args, 1)
	minZ, _ := getArgFloat("boundingBox", args, 2)
	maxX, _ := getArgFloat("boundingBox", args, 3)
	maxY, _ := getArgFloat("boundingBox", args, 4)
	maxZ, _ := getArgFloat("boundingBox", args, 5)
	return boundingBoxToMap(rl.NewBoundingBox(
		rl.NewVector3(float32(minX), float32(minY), float32(minZ)),
		rl.NewVector3(float32(maxX), float32(maxY), float32(maxZ)),
	)), nil
}

func builtinRay(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("ray", args, 6); err != nil {
		return nil, err
	}
	px, _ := getArgFloat("ray", args, 0)
	py, _ := getArgFloat("ray", args, 1)
	pz, _ := getArgFloat("ray", args, 2)
	dx, _ := getArgFloat("ray", args, 3)
	dy, _ := getArgFloat("ray", args, 4)
	dz, _ := getArgFloat("ray", args, 5)
	return vMap(map[string]candy_evaluator.Value{
		"position": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: px},
			"y": {Kind: candy_evaluator.ValFloat, F64: py},
			"z": {Kind: candy_evaluator.ValFloat, F64: pz},
		}},
		"direction": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: dx},
			"y": {Kind: candy_evaluator.ValFloat, F64: dy},
			"z": {Kind: candy_evaluator.ValFloat, F64: dz},
		}},
	}), nil
}

func builtinCamera3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("camera3D expects posX,posY,posZ, targetX,targetY,targetZ, fovy, [projection]")
	}
	cx, _ := getArgFloat("camera3D", args, 0)
	cy, _ := getArgFloat("camera3D", args, 1)
	cz, _ := getArgFloat("camera3D", args, 2)
	tx, _ := getArgFloat("camera3D", args, 3)
	ty, _ := getArgFloat("camera3D", args, 4)
	tz, _ := getArgFloat("camera3D", args, 5)
	fovy, _ := getArgFloat("camera3D", args, 6)
	proj := int64(rl.CameraPerspective)
	if len(args) > 7 {
		proj, _ = argInt("camera3D", args, 7)
	}
	return vMap(map[string]candy_evaluator.Value{
		"posX": {Kind: candy_evaluator.ValFloat, F64: cx},
		"posY": {Kind: candy_evaluator.ValFloat, F64: cy},
		"posZ": {Kind: candy_evaluator.ValFloat, F64: cz},
		"targetX":    {Kind: candy_evaluator.ValFloat, F64: tx},
		"targetY":    {Kind: candy_evaluator.ValFloat, F64: ty},
		"targetZ":    {Kind: candy_evaluator.ValFloat, F64: tz},
		"fovy":       {Kind: candy_evaluator.ValFloat, F64: fovy},
		"projection": {Kind: candy_evaluator.ValInt, I64: proj},
	}), nil
}

func builtinCamera2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("camera2D expects offsetX,offsetY, targetX,targetY, [rotation], [zoom]")
	}
	ox, _ := getArgFloat("camera2D", args, 0)
	oy, _ := getArgFloat("camera2D", args, 1)
	tx, _ := getArgFloat("camera2D", args, 2)
	ty, _ := getArgFloat("camera2D", args, 3)
	rot := 0.0
	if len(args) > 4 {
		rot, _ = getArgFloat("camera2D", args, 4)
	}
	zoom := 1.0
	if len(args) > 5 {
		zoom, _ = getArgFloat("camera2D", args, 5)
	}
	return vMap(map[string]candy_evaluator.Value{
		"offsetX": {Kind: candy_evaluator.ValFloat, F64: ox},
		"offsetY": {Kind: candy_evaluator.ValFloat, F64: oy},
		"targetX":  {Kind: candy_evaluator.ValFloat, F64: tx},
		"targetY":  {Kind: candy_evaluator.ValFloat, F64: ty},
		"rotation": {Kind: candy_evaluator.ValFloat, F64: rot},
		"zoom":     {Kind: candy_evaluator.ValFloat, F64: zoom},
	}), nil
}

// ---- Named color builtins ----
// Each returns a {r,g,b,a} map identical to colorToMap(rl.XxxColor).

func namedColorBuiltin(c rl.Color) func([]*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return func(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
		return colorToMap(c), nil
	}
}

var (
	builtinColorLightGray  = namedColorBuiltin(rl.LightGray)
	builtinColorGray       = namedColorBuiltin(rl.Gray)
	builtinColorDarkGray   = namedColorBuiltin(rl.DarkGray)
	builtinColorYellow     = namedColorBuiltin(rl.Yellow)
	builtinColorGold       = namedColorBuiltin(rl.Gold)
	builtinColorOrange     = namedColorBuiltin(rl.Orange)
	builtinColorPink       = namedColorBuiltin(rl.Pink)
	builtinColorRed        = namedColorBuiltin(rl.Red)
	builtinColorMaroon     = namedColorBuiltin(rl.Maroon)
	builtinColorGreen      = namedColorBuiltin(rl.Green)
	builtinColorLime       = namedColorBuiltin(rl.Lime)
	builtinColorDarkGreen  = namedColorBuiltin(rl.DarkGreen)
	builtinColorSkyBlue    = namedColorBuiltin(rl.SkyBlue)
	builtinColorBlue       = namedColorBuiltin(rl.Blue)
	builtinColorDarkBlue   = namedColorBuiltin(rl.DarkBlue)
	builtinColorPurple     = namedColorBuiltin(rl.Purple)
	builtinColorViolet     = namedColorBuiltin(rl.Violet)
	builtinColorDarkPurple = namedColorBuiltin(rl.DarkPurple)
	builtinColorBeige      = namedColorBuiltin(rl.Beige)
	builtinColorBrown      = namedColorBuiltin(rl.Brown)
	builtinColorDarkBrown  = namedColorBuiltin(rl.DarkBrown)
	builtinColorWhite      = namedColorBuiltin(rl.White)
	builtinColorBlack      = namedColorBuiltin(rl.Black)
	builtinColorBlank      = namedColorBuiltin(rl.Blank)
	builtinColorMagenta    = namedColorBuiltin(rl.Magenta)
	builtinColorRayWhite   = namedColorBuiltin(rl.RayWhite)
)

// builtinGetNamedColor returns a {r,g,b,a} map for any named color string.
func builtinGetNamedColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getNamedColor", args, 1); err != nil {
		return nil, err
	}
	name, err := argString("getNamedColor", args, 0)
	if err != nil {
		return nil, err
	}
	return colorToMap(colorFrom(name)), nil
}
