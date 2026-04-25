package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// registerKidExtended adds teaching-friendly names for CANDY_KID_EXTENDED.md (see docs).
func registerKidExtended() {
	candy_evaluator.RegisterBuiltin("seconds", builtinGetTime)
	candy_evaluator.RegisterBuiltin("deltaTime", builtinGetFrameTime)
	candy_evaluator.RegisterBuiltin("touching", builtinCircleCollision)
	candy_evaluator.RegisterBuiltin("boxHit", builtinBoxCollision)
	candy_evaluator.RegisterBuiltin("inside", builtinPointInBox)
	candy_evaluator.RegisterBuiltin("gridToPixel", builtinGridToPixel)
	candy_evaluator.RegisterBuiltin("pixelToGrid", builtinPixelToGrid)
	candy_evaluator.RegisterBuiltin("sprite", builtinLoadTexture)
	candy_evaluator.RegisterBuiltin("draw", builtinDrawSprite)
	candy_evaluator.RegisterBuiltin("drawRotated", builtinDrawSpriteRotated)
	candy_evaluator.RegisterBuiltin("drawFlipped", builtinDrawSpriteFlipped)
	candy_evaluator.RegisterBuiltin("unload", builtinUnloadTexture)
	candy_evaluator.RegisterBuiltin("distance", builtinDistance2D)
	candy_evaluator.RegisterBuiltin("angleBetween", builtinAngleTo)
	// 2D camera helper (see game_helpers gameCamera)
	candy_evaluator.RegisterBuiltin("camera", builtinCameraSnapTo)
	candy_evaluator.RegisterBuiltin("zoom", builtinCameraZoom)
	candy_evaluator.RegisterBuiltin("shake", builtinScreenShakeKid)
	candy_evaluator.RegisterBuiltin("debugLine", builtinDebugLine2D)
}

func builtinGridToPixel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("gridToPixel", args, 2); err != nil {
		return nil, err
	}
	c, _ := getArgFloat("gridToPixel", args, 0)
	t, _ := getArgFloat("gridToPixel", args, 1)
	return vFloat(c * t), nil
}

func builtinPixelToGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("pixelToGrid", args, 2); err != nil {
		return nil, err
	}
	p, _ := getArgFloat("pixelToGrid", args, 0)
	t, _ := getArgFloat("pixelToGrid", args, 1)
	if t == 0 {
		return vFloat(0), nil
	}
	return vFloat(math.Floor(p / t)), nil
}

// draw(textureId, x, y) or draw(id, x, y, color) or draw(id, x, y, w, h, [color]).
func builtinDrawSprite(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("draw expects textureId, x, y, [color] or + w, h, [color]")
	}
	_, tex, err := textureByID("draw", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := getArgFloat("draw", args, 1)
	y, _ := getArgFloat("draw", args, 2)
	switch len(args) {
	case 3:
		rl.DrawTexture(tex, int32(x), int32(y), rl.White)
		return null(), nil
	case 4:
		c, _ := argColor("draw", args, 3, rl.White)
		rl.DrawTexture(tex, int32(x), int32(y), c)
		return null(), nil
	default:
		w, _ := getArgFloat("draw", args, 3)
		h, _ := getArgFloat("draw", args, 4)
		c, _ := argColor("draw", args, 5, rl.White)
		src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
		dst := rl.NewRectangle(float32(x), float32(y), float32(w), float32(h))
		rl.DrawTexturePro(tex, src, dst, rl.NewVector2(0, 0), 0, c)
		return null(), nil
	}
}

// drawRotated(id, x, y, angleDeg, [scale], [color]).
func builtinDrawSpriteRotated(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawRotated expects textureId, x, y, angleDeg, [scale], [color]")
	}
	_, tex, err := textureByID("drawRotated", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := getArgFloat("drawRotated", args, 1)
	y, _ := getArgFloat("drawRotated", args, 2)
	ang, _ := getArgFloat("drawRotated", args, 3)
	sc := 1.0
	if len(args) > 4 {
		sc, _ = getArgFloat("drawRotated", args, 4)
	}
	c, _ := argColor("drawRotated", args, 5, rl.White)
	rl.DrawTextureEx(tex, rl.NewVector2(float32(x), float32(y)), float32(ang), float32(sc), c)
	return null(), nil
}

// drawFlipped(id, x, y, flipH, flipV, [color]) — flips by negating the source width/height in the source rect.
func builtinDrawSpriteFlipped(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawFlipped expects textureId, x, y, flipH, flipV, [color]")
	}
	_, tex, err := textureByID("drawFlipped", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := getArgFloat("drawFlipped", args, 1)
	y, _ := getArgFloat("drawFlipped", args, 2)
	flipH := args[3] != nil && args[3].Truthy()
	flipV := args[4] != nil && args[4].Truthy()
	c, _ := argColor("drawFlipped", args, 5, rl.White)
	sx := float32(tex.Width)
	sy := float32(tex.Height)
	srcX, srcY, srcW, srcH := float32(0), float32(0), sx, sy
	if flipH {
		srcW = -sx
	}
	if flipV {
		srcH = -sy
	}
	src := rl.NewRectangle(srcX, srcY, srcW, srcH)
	dst := rl.NewRectangle(float32(x), float32(y), float32(math.Abs(float64(sx))), float32(math.Abs(float64(sy))))
	rl.DrawTexturePro(tex, src, dst, rl.NewVector2(0, 0), 0, c)
	return null(), nil
}

// shake(duration, intensity) — second argument sets 2D helper camera shake strength; first is reserved.
func builtinScreenShakeKid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("shake", args, 2); err != nil {
		return nil, err
	}
	_, _ = getArgFloat("shake", args, 0)
	intensity, _ := getArgFloat("shake", args, 1)
	gameCamera.shake = intensity
	if gameCamera.shake < 0 {
		gameCamera.shake = 0
	}
	return null(), nil
}

func builtinDebugLine2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("debugLine expects x1, y1, x2, y2, [color]")
	}
	x1, _ := getArgFloat("debugLine", args, 0)
	y1, _ := getArgFloat("debugLine", args, 1)
	x2, _ := getArgFloat("debugLine", args, 2)
	y2, _ := getArgFloat("debugLine", args, 3)
	color, _ := argColor("debugLine", args, 4, rl.Lime)
	rl.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2), color)
	return null(), nil
}
