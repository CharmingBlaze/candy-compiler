package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func builtinDrawText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("text: use text(x, y, text, [size], [color]) or text(text, x, y, [size], [color])")
	}
	// Kid / Candy: text("Hello", 10, 50, 20, BLACK) — first arg is the string.
	// Original API: text(x, y, "Hi", 20, color)
	if args[0] != nil && args[0].Kind == candy_evaluator.ValString {
		s := args[0].Str
		x, _ := argInt("text", args, 1)
		y, _ := argInt("text", args, 2)
		fontSize := int64(20)
		if len(args) > 3 {
			fontSize, _ = argInt("text", args, 3)
		}
		c, _ := argColor("text", args, 4, rl.RayWhite)
		rl.DrawText(s, int32(x), int32(y), int32(fontSize), c)
		return null(), nil
	}
	x, _ := argInt("text", args, 0)
	y, _ := argInt("text", args, 1)
	s, err := argString("text", args, 2)
	if err != nil {
		return nil, err
	}
	fontSize := int64(20)
	if len(args) > 3 {
		fontSize, _ = argInt("text", args, 3)
	}
	c, _ := argColor("text", args, 4, rl.RayWhite)
	rl.DrawText(s, int32(x), int32(y), int32(fontSize), c)
	return null(), nil
}

func measureText(name string, args []*candy_evaluator.Value, i int) (int64, error) {
	s, err := argString(name, args, i)
	if err != nil { return 0, err }
	return int64(rl.MeasureText(s, 20)), nil
}

func fontByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Font, error) {
	id, err := argInt(name, args, i)
	if err != nil { return 0, rl.Font{}, err }
	f, ok := fonts[id]
	if !ok { return 0, rl.Font{}, fmt.Errorf("%s: invalid font handle %d", name, id) }
	return id, f, nil
}

func builtinLoadFont(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadFont", args, 1); err != nil { return nil, err }
	path, err := argString("loadFont", args, 0)
	if err != nil { return nil, err }
	f := rl.LoadFont(path)
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

func builtinUnloadFont(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, f, err := fontByID("unloadFont", args, 0)
	if err != nil { return nil, err }
	rl.UnloadFont(f)
	delete(fonts, id)
	return null(), nil
}

func builtinDrawTextEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 { return nil, fmt.Errorf("drawTextEx expects fontId, text, x, y, fontSize, spacing, [color]") }
	_, f, err := fontByID("drawTextEx", args, 0)
	if err != nil { return nil, err }
	s, _ := argString("drawTextEx", args, 1)
	x, _ := getArgFloat("drawTextEx", args, 2)
	y, _ := getArgFloat("drawTextEx", args, 3)
	fs, _ := getArgFloat("drawTextEx", args, 4)
	sp, _ := getArgFloat("drawTextEx", args, 5)
	c, _ := argColor("drawTextEx", args, 6, rl.RayWhite)
	rl.DrawTextEx(f, s, rl.NewVector2(float32(x), float32(y)), float32(fs), float32(sp), c)
	return null(), nil
}

func builtinMeasureText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("measureText", args, 2); err != nil {
		return nil, err
	}
	s, err := argString("measureText", args, 0)
	if err != nil {
		return nil, err
	}
	size, _ := argInt("measureText", args, 1)
	return vInt(int64(rl.MeasureText(s, int32(size)))), nil
}
