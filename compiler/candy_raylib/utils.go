package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	nextTextureID int64 = 1
	textures            = map[int64]rl.Texture2D{}

	nextSoundID int64 = 1
	sounds            = map[int64]rl.Sound{}

	nextMusicID int64 = 1
	musics            = map[int64]rl.Music{}

	nextFontID int64 = 1
	fonts            = map[int64]rl.Font{}

	nextShaderID int64 = 1
	shaders            = map[int64]rl.Shader{}

	nextRenderTextureID int64 = 1
	renderTextures            = map[int64]rl.RenderTexture2D{}

	nextModelID int64 = 1
	models            = map[int64]rl.Model{}

	nextModelAnimID int64 = 1
	modelAnims            = map[int64][]rl.ModelAnimation{}

	nextImageID int64 = 1
	images            = map[int64]*rl.Image{}

	nextMeshID int64 = 1
	meshes           = map[int64]rl.Mesh{}

	nextMaterialID int64 = 1
	materials            = map[int64]rl.Material{}

	nextWaveID int64 = 1
	waves            = map[int64]rl.Wave{}

	nextAudioStreamID int64 = 1
	audioStreams            = map[int64]rl.AudioStream{}

	activeCamera3D rl.Camera3D

	// Blitz-style frame clear: set by cameraClsColor; applied inside flip() after BeginDrawing.
	blitzFrameClearValid bool
	blitzFrameClear     rl.Color

	colorCache = map[string]rl.Color{}
)

func null() *candy_evaluator.Value { return &candy_evaluator.Value{Kind: candy_evaluator.ValNull} }
func vBool(b bool) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValBool, B: b}
}
func vInt(n int64) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValInt, I64: n}
}
func vFloat(f float64) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: f}
}
func vMap(m map[string]candy_evaluator.Value) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValMap, StrMap: m}
}

func expectArgs(name string, args []*candy_evaluator.Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s expects %d args, got %d", name, n, len(args))
	}
	return nil
}

func argInt(name string, args []*candy_evaluator.Value, i int) (int64, error) {
	if i >= len(args) || args[i] == nil {
		return 0, fmt.Errorf("%s arg %d must be int/float", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValInt {
		return v.I64, nil
	}
	if v.Kind == candy_evaluator.ValFloat {
		return int64(v.F64), nil
	}
	return 0, fmt.Errorf("%s arg %d must be int/float", name, i+1)
}

func getArgFloat(name string, args []*candy_evaluator.Value, i int) (float64, error) {
	if i >= len(args) || args[i] == nil {
		return 0, fmt.Errorf("%s arg %d must be number", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValInt {
		return float64(v.I64), nil
	}
	if v.Kind == candy_evaluator.ValFloat {
		return v.F64, nil
	}
	return 0, fmt.Errorf("%s arg %d must be int/float", name, i+1)
}

func argString(name string, args []*candy_evaluator.Value, i int) (string, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValString {
		return "", fmt.Errorf("%s arg %d must be string", name, i+1)
	}
	return args[i].Str, nil
}

// keyArg is for keyboard functions: `key(LEFT)` (int) or `key("left")` (string name).
func keyArg(name string, v *candy_evaluator.Value) (int32, error) {
	if v == nil {
		return 0, fmt.Errorf("%s: key argument is required", name)
	}
	switch v.Kind {
	case candy_evaluator.ValInt:
		return int32(v.I64), nil
	case candy_evaluator.ValFloat:
		return int32(v.F64), nil
	case candy_evaluator.ValString:
		return keyCode(v.Str), nil
	default:
		return 0, fmt.Errorf("%s: key must be int or string (key name)", name)
	}
}

func colorFrom(name string) rl.Color {
	key := strings.ToLower(strings.TrimSpace(name))
	if c, ok := colorCache[key]; ok {
		return c
	}

	var c rl.Color
	switch key {
	case "white":
		c = rl.RayWhite
	case "black":
		c = rl.Black
	case "red":
		c = rl.Red
	case "green":
		c = rl.Green
	case "blue":
		c = rl.Blue
	case "yellow":
		c = rl.Yellow
	case "gray":
		c = rl.Gray
	case "magenta":
		c = rl.Magenta
	case "gold":
		c = rl.Gold
	case "lime":
		c = rl.Lime
	case "darkgreen":
		c = rl.DarkGreen
	case "darkblue":
		c = rl.DarkBlue
	case "darkgray":
		c = rl.DarkGray
	case "maroon":
		c = rl.Maroon
	case "navy":
		c = rl.NewColor(0, 0, 128, 255)
	case "orange":
		c = rl.Orange
	case "purple":
		c = rl.Purple
	case "pink":
		c = rl.Pink
	case "brown":
		c = rl.Brown
	case "sky", "skyblue":
		c = rl.SkyBlue
	case "violet":
		c = rl.Violet
	case "beige", "beige_candy", "beige_old":
		c = rl.Beige
	case "lightgray", "light_gray":
		c = rl.LightGray
	case "darkpurple", "dark_purple":
		c = rl.DarkPurple
	case "darkbrown", "dark_brown":
		c = rl.DarkBrown
	case "blank", "transparent":
		c = rl.Blank
	case "raywhite", "ray_white":
		c = rl.RayWhite
	default:
		c = rl.RayWhite
	}
	colorCache[key] = c
	return c
}

func argColor(name string, args []*candy_evaluator.Value, i int, fallback rl.Color) (rl.Color, error) {
	if i >= len(args) || args[i] == nil {
		return fallback, nil
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValString {
		return colorFrom(v.Str), nil
	}
	return fallback, fmt.Errorf("%s arg %d must be color string", name, i+1)
}

func keyCode(name string) int32 {
	switch strings.ToLower(name) {
	case "left":
		return rl.KeyLeft
	case "right":
		return rl.KeyRight
	case "up":
		return rl.KeyUp
	case "down":
		return rl.KeyDown
	case "escape", "esc", "key_esc":
		return rl.KeyEscape
	case "space":
		return rl.KeySpace
	case "enter":
		return rl.KeyEnter
	case "tab":
		return rl.KeyTab
	case "backspace":
		return rl.KeyBackspace
	case "shift":
		return rl.KeyLeftShift
	case "ctrl":
		return rl.KeyLeftControl
	case "alt":
		return rl.KeyLeftAlt
	case "f1":
		return rl.KeyF1
	case "f2":
		return rl.KeyF2
	case "f3":
		return rl.KeyF3
	case "f4":
		return rl.KeyF4
	case "f5":
		return rl.KeyF5
	case "f6":
		return rl.KeyF6
	case "f7":
		return rl.KeyF7
	case "f8":
		return rl.KeyF8
	case "f9":
		return rl.KeyF9
	case "f10":
		return rl.KeyF10
	case "f11":
		return rl.KeyF11
	case "f12":
		return rl.KeyF12
	case "a":
		return rl.KeyA
	case "w":
		return rl.KeyW
	case "s":
		return rl.KeyS
	case "d":
		return rl.KeyD
	default:
		return 0
	}
}

func mouseButtonCode(n int64) rl.MouseButton {
	switch n {
	case 0:
		return rl.MouseButtonLeft
	case 1:
		return rl.MouseButtonRight
	case 2:
		return rl.MouseButtonMiddle
	default:
		return rl.MouseButton(n)
	}
}

func argVector2(name string, args []*candy_evaluator.Value, i int) (rl.Vector2, error) {
	if i >= len(args) || args[i] == nil {
		return rl.Vector2{}, fmt.Errorf("%s arg %d is nil", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValVec {
		if len(v.Vec) < 2 {
			return rl.Vector2{}, fmt.Errorf("%s arg %d vector must have at least 2 elements, got %d", name, i+1, len(v.Vec))
		}
		return rl.NewVector2(float32(v.Vec[0]), float32(v.Vec[1])), nil
	}
	if v.Kind != candy_evaluator.ValMap {
		return rl.Vector2{}, fmt.Errorf("%s arg %d must be map {x, y} or vec2", name, i+1)
	}
	m := v.StrMap
	x := 0.0
	if v, ok := m["x"]; ok && (v.Kind == candy_evaluator.ValFloat || v.Kind == candy_evaluator.ValInt) {
		if v.Kind == candy_evaluator.ValInt {
			x = float64(v.I64)
		} else {
			x = v.F64
		}
	}
	y := 0.0
	if v, ok := m["y"]; ok && (v.Kind == candy_evaluator.ValFloat || v.Kind == candy_evaluator.ValInt) {
		if v.Kind == candy_evaluator.ValInt {
			y = float64(v.I64)
		} else {
			y = v.F64
		}
	}
	return rl.NewVector2(float32(x), float32(y)), nil
}

func argVector3(name string, args []*candy_evaluator.Value, i int) (rl.Vector3, error) {
	if i >= len(args) || args[i] == nil {
		return rl.Vector3{}, fmt.Errorf("%s arg %d is nil", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValVec {
		if len(v.Vec) < 3 {
			return rl.Vector3{}, fmt.Errorf("%s arg %d vector must have at least 3 elements, got %d", name, i+1, len(v.Vec))
		}
		return rl.NewVector3(float32(v.Vec[0]), float32(v.Vec[1]), float32(v.Vec[2])), nil
	}
	if v.Kind != candy_evaluator.ValMap {
		return rl.Vector3{}, fmt.Errorf("%s arg %d must be map {x, y, z} or vec3", name, i+1)
	}
	m := v.StrMap
	var x, y, z float64
	if vv, ok := m["x"]; ok {
		if vv.Kind == candy_evaluator.ValInt {
			x = float64(vv.I64)
		} else {
			x = vv.F64
		}
	}
	if vv, ok := m["y"]; ok {
		if vv.Kind == candy_evaluator.ValInt {
			y = float64(vv.I64)
		} else {
			y = vv.F64
		}
	}
	if vv, ok := m["z"]; ok {
		if vv.Kind == candy_evaluator.ValInt {
			z = float64(vv.I64)
		} else {
			z = vv.F64
		}
	}
	return rl.NewVector3(float32(x), float32(y), float32(z)), nil
}

func argRectangle(name string, args []*candy_evaluator.Value, i int) (rl.Rectangle, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.Rectangle{}, fmt.Errorf("%s arg %d must be map {x, y, width, height}", name, i+1)
	}
	m := args[i].StrMap
	var x, y, w, h float64
	if v, ok := m["x"]; ok {
		if v.Kind == candy_evaluator.ValInt {
			x = float64(v.I64)
		} else {
			x = v.F64
		}
	}
	if v, ok := m["y"]; ok {
		if v.Kind == candy_evaluator.ValInt {
			y = float64(v.I64)
		} else {
			y = v.F64
		}
	}
	if v, ok := m["width"]; ok {
		if v.Kind == candy_evaluator.ValInt {
			w = float64(v.I64)
		} else {
			w = v.F64
		}
	}
	if v, ok := m["height"]; ok {
		if v.Kind == candy_evaluator.ValInt {
			h = float64(v.I64)
		} else {
			h = v.F64
		}
	}
	return rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)), nil
}
