package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func builtinDrawPixel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawPixel expects x, y, [color]")
	}
	x, _ := argInt("drawPixel", args, 0)
	y, _ := argInt("drawPixel", args, 1)
	c, _ := argColor("drawPixel", args, 2, rl.RayWhite)
	rl.DrawPixel(int32(x), int32(y), c)
	return null(), nil
}

func builtinDrawLine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawLine expects x1, y1, x2, y2, [color]")
	}
	x1, _ := argInt("drawLine", args, 0)
	y1, _ := argInt("drawLine", args, 1)
	x2, _ := argInt("drawLine", args, 2)
	y2, _ := argInt("drawLine", args, 3)
	c, _ := argColor("drawLine", args, 4, rl.RayWhite)
	rl.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2), c)
	return null(), nil
}

func builtinDrawCircle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawCircle expects x, y, radius, [color]")
	}
	x, _ := argInt("drawCircle", args, 0)
	y, _ := argInt("drawCircle", args, 1)
	r, _ := getArgFloat("drawCircle", args, 2)
	c, _ := argColor("drawCircle", args, 3, rl.RayWhite)
	rl.DrawCircle(int32(x), int32(y), float32(r), c)
	return null(), nil
}

func builtinDrawCircleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawCircleLines expects x, y, radius, [color]")
	}
	x, _ := argInt("drawCircleLines", args, 0)
	y, _ := argInt("drawCircleLines", args, 1)
	r, _ := getArgFloat("drawCircleLines", args, 2)
	c, _ := argColor("drawCircleLines", args, 3, rl.RayWhite)
	rl.DrawCircleLines(int32(x), int32(y), float32(r), c)
	return null(), nil
}

func builtinDrawEllipse(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawEllipse expects x, y, radiusH, radiusV, [color]")
	}
	x, _ := argInt("drawEllipse", args, 0)
	y, _ := argInt("drawEllipse", args, 1)
	rh, _ := getArgFloat("drawEllipse", args, 2)
	rv, _ := getArgFloat("drawEllipse", args, 3)
	c, _ := argColor("drawEllipse", args, 4, rl.RayWhite)
	rl.DrawEllipse(int32(x), int32(y), float32(rh), float32(rv), c)
	return null(), nil
}

func builtinDrawRing(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRing expects x, y, innerR, outerR, startAngle, endAngle, [segments], [color]")
	}
	x, _ := getArgFloat("drawRing", args, 0)
	y, _ := getArgFloat("drawRing", args, 1)
	ir, _ := getArgFloat("drawRing", args, 2)
	or, _ := getArgFloat("drawRing", args, 3)
	sa, _ := getArgFloat("drawRing", args, 4)
	ea, _ := getArgFloat("drawRing", args, 5)
	seg := int64(36)
	if len(args) > 6 {
		seg, _ = argInt("drawRing", args, 6)
	}
	c, _ := argColor("drawRing", args, 7, rl.RayWhite)
	rl.DrawRing(rl.NewVector2(float32(x), float32(y)), float32(ir), float32(or), float32(sa), float32(ea), int32(seg), c)
	return null(), nil
}

func builtinDrawRectangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawRectangle expects x, y, w, h, [color]")
	}
	x, _ := argInt("drawRectangle", args, 0)
	y, _ := argInt("drawRectangle", args, 1)
	w, _ := argInt("drawRectangle", args, 2)
	h, _ := argInt("drawRectangle", args, 3)
	c, _ := argColor("drawRectangle", args, 4, rl.RayWhite)
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), c)
	return null(), nil
}

func builtinDrawRectangleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawRectangleLines expects x, y, w, h, [color]")
	}
	x, _ := argInt("drawRectangleLines", args, 0)
	y, _ := argInt("drawRectangleLines", args, 1)
	w, _ := argInt("drawRectangleLines", args, 2)
	h, _ := argInt("drawRectangleLines", args, 3)
	c, _ := argColor("drawRectangleLines", args, 4, rl.RayWhite)
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), c)
	return null(), nil
}

func builtinDrawRectangleRounded(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRectangleRounded expects x, y, w, h, roundness, segments, [color]")
	}
	x, _ := getArgFloat("drawRectangleRounded", args, 0)
	y, _ := getArgFloat("drawRectangleRounded", args, 1)
	w, _ := getArgFloat("drawRectangleRounded", args, 2)
	h, _ := getArgFloat("drawRectangleRounded", args, 3)
	r, _ := getArgFloat("drawRectangleRounded", args, 4)
	s, _ := argInt("drawRectangleRounded", args, 5)
	c, _ := argColor("drawRectangleRounded", args, 6, rl.RayWhite)
	rl.DrawRectangleRounded(rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)), float32(r), int32(s), c)
	return null(), nil
}

func builtinDrawTriangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawTriangle expects x1, y1, x2, y2, x3, y3, [color]")
	}
	x1, _ := getArgFloat("drawTriangle", args, 0)
	y1, _ := getArgFloat("drawTriangle", args, 1)
	x2, _ := getArgFloat("drawTriangle", args, 2)
	y2, _ := getArgFloat("drawTriangle", args, 3)
	x3, _ := getArgFloat("drawTriangle", args, 4)
	y3, _ := getArgFloat("drawTriangle", args, 5)
	c, _ := argColor("drawTriangle", args, 6, rl.RayWhite)
	rl.DrawTriangle(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x2), float32(y2)), rl.NewVector2(float32(x3), float32(y3)), c)
	return null(), nil
}

func builtinDrawFPS(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawFPS expects x, y")
	}
	x, _ := argInt("drawFPS", args, 0)
	y, _ := argInt("drawFPS", args, 1)
	rl.DrawFPS(int32(x), int32(y))
	return null(), nil
}

// ---- Collision Detection ----

func builtinCheckCollisionRecs(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionRecs", args, 8); err != nil {
		return nil, err
	}
	x1, _ := getArgFloat("checkCollisionRecs", args, 0)
	y1, _ := getArgFloat("checkCollisionRecs", args, 1)
	w1, _ := getArgFloat("checkCollisionRecs", args, 2)
	h1, _ := getArgFloat("checkCollisionRecs", args, 3)
	x2, _ := getArgFloat("checkCollisionRecs", args, 4)
	y2, _ := getArgFloat("checkCollisionRecs", args, 5)
	w2, _ := getArgFloat("checkCollisionRecs", args, 6)
	h2, _ := getArgFloat("checkCollisionRecs", args, 7)
	res := rl.CheckCollisionRecs(
		rl.NewRectangle(float32(x1), float32(y1), float32(w1), float32(h1)),
		rl.NewRectangle(float32(x2), float32(y2), float32(w2), float32(h2)),
	)
	return vBool(res), nil
}

func builtinCheckCollisionCircles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionCircles", args, 6); err != nil {
		return nil, err
	}
	x1, _ := getArgFloat("checkCollisionCircles", args, 0)
	y1, _ := getArgFloat("checkCollisionCircles", args, 1)
	r1, _ := getArgFloat("checkCollisionCircles", args, 2)
	x2, _ := getArgFloat("checkCollisionCircles", args, 3)
	y2, _ := getArgFloat("checkCollisionCircles", args, 4)
	r2, _ := getArgFloat("checkCollisionCircles", args, 5)
	res := rl.CheckCollisionCircles(
		rl.NewVector2(float32(x1), float32(y1)), float32(r1),
		rl.NewVector2(float32(x2), float32(y2)), float32(r2),
	)
	return vBool(res), nil
}

func builtinCheckCollisionCircleRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionCircleRec", args, 7); err != nil {
		return nil, err
	}
	cx, _ := getArgFloat("checkCollisionCircleRec", args, 0)
	cy, _ := getArgFloat("checkCollisionCircleRec", args, 1)
	cr, _ := getArgFloat("checkCollisionCircleRec", args, 2)
	rx, _ := getArgFloat("checkCollisionCircleRec", args, 3)
	ry, _ := getArgFloat("checkCollisionCircleRec", args, 4)
	rw, _ := getArgFloat("checkCollisionCircleRec", args, 5)
	rh, _ := getArgFloat("checkCollisionCircleRec", args, 6)
	res := rl.CheckCollisionCircleRec(
		rl.NewVector2(float32(cx), float32(cy)), float32(cr),
		rl.NewRectangle(float32(rx), float32(ry), float32(rw), float32(rh)),
	)
	return vBool(res), nil
}

// ---- Splines (Raylib 5.0+) ----

func builtinDrawSplineLinear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineLinear expects pointsArr, thickness, [color]")
	}
	pts, err := parsePointsArray("drawSplineLinear", args[0])
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineLinear", args, 1)
	c, _ := argColor("drawSplineLinear", args, 2, rl.RayWhite)
	rl.DrawSplineLinear(pts, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineBasis(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineBasis expects pointsArr, thickness, [color]")
	}
	pts, err := parsePointsArray("drawSplineBasis", args[0])
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineBasis", args, 1)
	c, _ := argColor("drawSplineBasis", args, 2, rl.RayWhite)
	rl.DrawSplineBasis(pts, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineCatmullRom(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineCatmullRom expects pointsArr, thickness, [color]")
	}
	pts, err := parsePointsArray("drawSplineCatmullRom", args[0])
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineCatmullRom", args, 1)
	c, _ := argColor("drawSplineCatmullRom", args, 2, rl.RayWhite)
	rl.DrawSplineCatmullRom(pts, float32(thick), c)
	return null(), nil
}

// ---- Shapes texture ----

func builtinSetShapesTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setShapesTexture", args, 2); err != nil {
		return nil, err
	}
	texID, _ := argInt("setShapesTexture", args, 0)
	tex, ok := textures[texID]
	if !ok {
		return nil, fmt.Errorf("setShapesTexture: invalid texture handle %d", texID)
	}
	rec, err := argRectangle("setShapesTexture", args, 1)
	if err != nil {
		return nil, err
	}
	rl.SetShapesTexture(tex, rec)
	return null(), nil
}

func builtinGetShapesTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	tex := rl.GetShapesTexture()
	id := nextTextureID
	nextTextureID++
	textures[id] = tex
	return vInt(id), nil
}

func builtinGetShapesTextureRectangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rec := rl.GetShapesTextureRectangle()
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.X)},
		"y":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.Y)},
		"width":  {Kind: candy_evaluator.ValFloat, F64: float64(rec.Width)},
		"height": {Kind: candy_evaluator.ValFloat, F64: float64(rec.Height)},
	}), nil
}

// ---- Draw V / Ex variants ----

func builtinDrawPixelV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("drawPixelV", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("drawPixelV", args, 0)
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawPixelV", args, 1, rl.RayWhite)
	rl.DrawPixelV(v, c)
	return null(), nil
}

func builtinDrawLineV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("drawLineV", args, 3); err != nil {
		return nil, err
	}
	s, err := argVector2("drawLineV", args, 0)
	if err != nil {
		return nil, err
	}
	e, err := argVector2("drawLineV", args, 1)
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawLineV", args, 2, rl.RayWhite)
	rl.DrawLineV(s, e, c)
	return null(), nil
}

func builtinDrawLineEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawLineEx expects startPos, endPos, thick, [color]")
	}
	s, err := argVector2("drawLineEx", args, 0)
	if err != nil {
		return nil, err
	}
	e, err := argVector2("drawLineEx", args, 1)
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawLineEx", args, 2)
	c, _ := argColor("drawLineEx", args, 3, rl.RayWhite)
	rl.DrawLineEx(s, e, float32(thick), c)
	return null(), nil
}

func builtinDrawLineStrip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawLineStrip expects pointsArr, [color]")
	}
	pts, err := parsePointsArray("drawLineStrip", args[0])
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawLineStrip", args, 1, rl.RayWhite)
	rl.DrawLineStrip(pts, c)
	return null(), nil
}

func builtinDrawLineBezier(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawLineBezier expects startPos, endPos, thick, [color]")
	}
	s, err := argVector2("drawLineBezier", args, 0)
	if err != nil {
		return nil, err
	}
	e, err := argVector2("drawLineBezier", args, 1)
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawLineBezier", args, 2)
	c, _ := argColor("drawLineBezier", args, 3, rl.RayWhite)
	rl.DrawLineBezier(s, e, float32(thick), c)
	return null(), nil
}

func builtinDrawLineDashed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawLineDashed expects startPos, endPos, dashSize, spaceSize, [color]")
	}
	s, err := argVector2("drawLineDashed", args, 0)
	if err != nil {
		return nil, err
	}
	e, err := argVector2("drawLineDashed", args, 1)
	if err != nil {
		return nil, err
	}
	dash, _ := argInt("drawLineDashed", args, 2)
	space, _ := argInt("drawLineDashed", args, 3)
	c, _ := argColor("drawLineDashed", args, 4, rl.RayWhite)
	rl.DrawLineDashed(s, e, int32(dash), int32(space), c)
	return null(), nil
}

func builtinDrawCircleV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawCircleV expects center, radius, [color]")
	}
	v, err := argVector2("drawCircleV", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawCircleV", args, 1)
	c, _ := argColor("drawCircleV", args, 2, rl.RayWhite)
	rl.DrawCircleV(v, float32(r), c)
	return null(), nil
}

func builtinDrawCircleGradient(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawCircleGradient expects center, radius, innerColor, outerColor")
	}
	v, err := argVector2("drawCircleGradient", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawCircleGradient", args, 1)
	inner, _ := argColor("drawCircleGradient", args, 2, rl.Red)
	outer, _ := argColor("drawCircleGradient", args, 3, rl.Blue)
	rl.DrawCircleGradient(int32(v.X), int32(v.Y), float32(r), inner, outer)
	return null(), nil
}

func builtinDrawCircleSector(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawCircleSector expects center, radius, startAngle, endAngle, segments, [color]")
	}
	v, err := argVector2("drawCircleSector", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawCircleSector", args, 1)
	sa, _ := getArgFloat("drawCircleSector", args, 2)
	ea, _ := getArgFloat("drawCircleSector", args, 3)
	seg, _ := argInt("drawCircleSector", args, 4)
	c, _ := argColor("drawCircleSector", args, 5, rl.RayWhite)
	rl.DrawCircleSector(v, float32(r), float32(sa), float32(ea), int32(seg), c)
	return null(), nil
}

func builtinDrawCircleSectorLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawCircleSectorLines expects center, radius, startAngle, endAngle, segments, [color]")
	}
	v, err := argVector2("drawCircleSectorLines", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawCircleSectorLines", args, 1)
	sa, _ := getArgFloat("drawCircleSectorLines", args, 2)
	ea, _ := getArgFloat("drawCircleSectorLines", args, 3)
	seg, _ := argInt("drawCircleSectorLines", args, 4)
	c, _ := argColor("drawCircleSectorLines", args, 5, rl.RayWhite)
	rl.DrawCircleSectorLines(v, float32(r), float32(sa), float32(ea), int32(seg), c)
	return null(), nil
}

func builtinDrawCircleLinesV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawCircleLinesV expects center, radius, [color]")
	}
	v, err := argVector2("drawCircleLinesV", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawCircleLinesV", args, 1)
	c, _ := argColor("drawCircleLinesV", args, 2, rl.RayWhite)
	rl.DrawCircleLinesV(v, float32(r), c)
	return null(), nil
}

func builtinDrawEllipseLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawEllipseLines expects x, y, radiusH, radiusV, [color]")
	}
	x, _ := argInt("drawEllipseLines", args, 0)
	y, _ := argInt("drawEllipseLines", args, 1)
	rh, _ := getArgFloat("drawEllipseLines", args, 2)
	rv, _ := getArgFloat("drawEllipseLines", args, 3)
	c, _ := argColor("drawEllipseLines", args, 4, rl.RayWhite)
	rl.DrawEllipseLines(int32(x), int32(y), float32(rh), float32(rv), c)
	return null(), nil
}

func builtinDrawRingLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRingLines expects x, y, innerR, outerR, startAngle, endAngle, [segments], [color]")
	}
	x, _ := getArgFloat("drawRingLines", args, 0)
	y, _ := getArgFloat("drawRingLines", args, 1)
	ir, _ := getArgFloat("drawRingLines", args, 2)
	or2, _ := getArgFloat("drawRingLines", args, 3)
	sa, _ := getArgFloat("drawRingLines", args, 4)
	ea, _ := getArgFloat("drawRingLines", args, 5)
	seg := int64(36)
	if len(args) > 6 {
		seg, _ = argInt("drawRingLines", args, 6)
	}
	c, _ := argColor("drawRingLines", args, 7, rl.RayWhite)
	rl.DrawRingLines(rl.NewVector2(float32(x), float32(y)), float32(ir), float32(or2), float32(sa), float32(ea), int32(seg), c)
	return null(), nil
}

func builtinDrawRectangleV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawRectangleV expects position, size, [color]")
	}
	pos, err := argVector2("drawRectangleV", args, 0)
	if err != nil {
		return nil, err
	}
	size, err := argVector2("drawRectangleV", args, 1)
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawRectangleV", args, 2, rl.RayWhite)
	rl.DrawRectangleV(pos, size, c)
	return null(), nil
}

func builtinDrawRectangleRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawRectangleRec expects rec, [color]")
	}
	rec, err := argRectangle("drawRectangleRec", args, 0)
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawRectangleRec", args, 1, rl.RayWhite)
	rl.DrawRectangleRec(rec, c)
	return null(), nil
}

func builtinDrawRectanglePro(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawRectanglePro expects rec, origin, rotation, [color]")
	}
	rec, err := argRectangle("drawRectanglePro", args, 0)
	if err != nil {
		return nil, err
	}
	origin, err := argVector2("drawRectanglePro", args, 1)
	if err != nil {
		return nil, err
	}
	rot, _ := getArgFloat("drawRectanglePro", args, 2)
	c, _ := argColor("drawRectanglePro", args, 3, rl.RayWhite)
	rl.DrawRectanglePro(rec, origin, float32(rot), c)
	return null(), nil
}

func builtinDrawRectangleGradientV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRectangleGradientV expects x, y, w, h, topColor, bottomColor")
	}
	x, _ := argInt("drawRectangleGradientV", args, 0)
	y, _ := argInt("drawRectangleGradientV", args, 1)
	w, _ := argInt("drawRectangleGradientV", args, 2)
	h, _ := argInt("drawRectangleGradientV", args, 3)
	top, _ := argColor("drawRectangleGradientV", args, 4, rl.Red)
	bot, _ := argColor("drawRectangleGradientV", args, 5, rl.Blue)
	rl.DrawRectangleGradientV(int32(x), int32(y), int32(w), int32(h), top, bot)
	return null(), nil
}

func builtinDrawRectangleGradientH(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRectangleGradientH expects x, y, w, h, leftColor, rightColor")
	}
	x, _ := argInt("drawRectangleGradientH", args, 0)
	y, _ := argInt("drawRectangleGradientH", args, 1)
	w, _ := argInt("drawRectangleGradientH", args, 2)
	h, _ := argInt("drawRectangleGradientH", args, 3)
	left, _ := argColor("drawRectangleGradientH", args, 4, rl.Red)
	right, _ := argColor("drawRectangleGradientH", args, 5, rl.Blue)
	rl.DrawRectangleGradientH(int32(x), int32(y), int32(w), int32(h), left, right)
	return null(), nil
}

func builtinDrawRectangleGradientEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawRectangleGradientEx expects rec, topLeft, bottomLeft, bottomRight, topRight")
	}
	rec, err := argRectangle("drawRectangleGradientEx", args, 0)
	if err != nil {
		return nil, err
	}
	tl, _ := argColor("drawRectangleGradientEx", args, 1, rl.Red)
	bl, _ := argColor("drawRectangleGradientEx", args, 2, rl.Green)
	br, _ := argColor("drawRectangleGradientEx", args, 3, rl.Blue)
	tr, _ := argColor("drawRectangleGradientEx", args, 4, rl.Yellow)
	rl.DrawRectangleGradientEx(rec, tl, bl, tr, br)
	return null(), nil
}

func builtinDrawRectangleLinesEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawRectangleLinesEx expects rec, lineThick, [color]")
	}
	rec, err := argRectangle("drawRectangleLinesEx", args, 0)
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawRectangleLinesEx", args, 1)
	c, _ := argColor("drawRectangleLinesEx", args, 2, rl.RayWhite)
	rl.DrawRectangleLinesEx(rec, float32(thick), c)
	return null(), nil
}

func builtinDrawRectangleRoundedLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawRectangleRoundedLines expects rec, roundness, segments, [color]")
	}
	rec, err := argRectangle("drawRectangleRoundedLines", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawRectangleRoundedLines", args, 1)
	seg, _ := argInt("drawRectangleRoundedLines", args, 2)
	c, _ := argColor("drawRectangleRoundedLines", args, 3, rl.RayWhite)
	rl.DrawRectangleRoundedLines(rec, float32(r), int32(seg), c)
	return null(), nil
}

func builtinDrawRectangleRoundedLinesEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawRectangleRoundedLinesEx expects rec, roundness, segments, lineThick, [color]")
	}
	rec, err := argRectangle("drawRectangleRoundedLinesEx", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("drawRectangleRoundedLinesEx", args, 1)
	seg, _ := argInt("drawRectangleRoundedLinesEx", args, 2)
	thick, _ := getArgFloat("drawRectangleRoundedLinesEx", args, 3)
	c, _ := argColor("drawRectangleRoundedLinesEx", args, 4, rl.RayWhite)
	rl.DrawRectangleRoundedLinesEx(rec, float32(r), int32(seg), float32(thick), c)
	return null(), nil
}

func builtinDrawTriangleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawTriangleLines expects x1,y1, x2,y2, x3,y3, [color]")
	}
	x1, _ := getArgFloat("drawTriangleLines", args, 0)
	y1, _ := getArgFloat("drawTriangleLines", args, 1)
	x2, _ := getArgFloat("drawTriangleLines", args, 2)
	y2, _ := getArgFloat("drawTriangleLines", args, 3)
	x3, _ := getArgFloat("drawTriangleLines", args, 4)
	y3, _ := getArgFloat("drawTriangleLines", args, 5)
	c, _ := argColor("drawTriangleLines", args, 6, rl.RayWhite)
	rl.DrawTriangleLines(
		rl.NewVector2(float32(x1), float32(y1)),
		rl.NewVector2(float32(x2), float32(y2)),
		rl.NewVector2(float32(x3), float32(y3)), c)
	return null(), nil
}

func builtinDrawTriangleFan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawTriangleFan expects pointsArr, [color]")
	}
	pts, err := parsePointsArray("drawTriangleFan", args[0])
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawTriangleFan", args, 1, rl.RayWhite)
	rl.DrawTriangleFan(pts, c)
	return null(), nil
}

func builtinDrawTriangleStrip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawTriangleStrip expects pointsArr, [color]")
	}
	pts, err := parsePointsArray("drawTriangleStrip", args[0])
	if err != nil {
		return nil, err
	}
	c, _ := argColor("drawTriangleStrip", args, 1, rl.RayWhite)
	rl.DrawTriangleStrip(pts, c)
	return null(), nil
}

func builtinDrawPoly(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawPoly expects center, sides, radius, rotation, [color]")
	}
	v, err := argVector2("drawPoly", args, 0)
	if err != nil {
		return nil, err
	}
	sides, _ := argInt("drawPoly", args, 1)
	r, _ := getArgFloat("drawPoly", args, 2)
	rot, _ := getArgFloat("drawPoly", args, 3)
	c, _ := argColor("drawPoly", args, 4, rl.RayWhite)
	rl.DrawPoly(v, int32(sides), float32(r), float32(rot), c)
	return null(), nil
}

func builtinDrawPolyLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawPolyLines expects center, sides, radius, rotation, [color]")
	}
	v, err := argVector2("drawPolyLines", args, 0)
	if err != nil {
		return nil, err
	}
	sides, _ := argInt("drawPolyLines", args, 1)
	r, _ := getArgFloat("drawPolyLines", args, 2)
	rot, _ := getArgFloat("drawPolyLines", args, 3)
	c, _ := argColor("drawPolyLines", args, 4, rl.RayWhite)
	rl.DrawPolyLines(v, int32(sides), float32(r), float32(rot), c)
	return null(), nil
}

func builtinDrawPolyLinesEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawPolyLinesEx expects center, sides, radius, rotation, lineThick, [color]")
	}
	v, err := argVector2("drawPolyLinesEx", args, 0)
	if err != nil {
		return nil, err
	}
	sides, _ := argInt("drawPolyLinesEx", args, 1)
	r, _ := getArgFloat("drawPolyLinesEx", args, 2)
	rot, _ := getArgFloat("drawPolyLinesEx", args, 3)
	thick, _ := getArgFloat("drawPolyLinesEx", args, 4)
	c, _ := argColor("drawPolyLinesEx", args, 5, rl.RayWhite)
	rl.DrawPolyLinesEx(v, int32(sides), float32(r), float32(rot), float32(thick), c)
	return null(), nil
}

// ---- Spline extras ----

func builtinDrawSplineBezierQuadratic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineBezierQuadratic expects pointsArr, thickness, [color]")
	}
	pts, err := parsePointsArray("drawSplineBezierQuadratic", args[0])
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineBezierQuadratic", args, 1)
	c, _ := argColor("drawSplineBezierQuadratic", args, 2, rl.RayWhite)
	rl.DrawSplineBezierQuadratic(pts, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineBezierCubic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineBezierCubic expects pointsArr, thickness, [color]")
	}
	pts, err := parsePointsArray("drawSplineBezierCubic", args[0])
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineBezierCubic", args, 1)
	c, _ := argColor("drawSplineBezierCubic", args, 2, rl.RayWhite)
	rl.DrawSplineBezierCubic(pts, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineSegmentLinear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawSplineSegmentLinear expects p1, p2, thick, [color]")
	}
	p1, err := argVector2("drawSplineSegmentLinear", args, 0)
	if err != nil {
		return nil, err
	}
	p2, err := argVector2("drawSplineSegmentLinear", args, 1)
	if err != nil {
		return nil, err
	}
	thick, _ := getArgFloat("drawSplineSegmentLinear", args, 2)
	c, _ := argColor("drawSplineSegmentLinear", args, 3, rl.RayWhite)
	rl.DrawSplineSegmentLinear(p1, p2, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineSegmentBasis(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawSplineSegmentBasis expects p1, p2, p3, p4, thick, [color]")
	}
	p1, _ := argVector2("drawSplineSegmentBasis", args, 0)
	p2, _ := argVector2("drawSplineSegmentBasis", args, 1)
	p3, _ := argVector2("drawSplineSegmentBasis", args, 2)
	p4, _ := argVector2("drawSplineSegmentBasis", args, 3)
	thick, _ := getArgFloat("drawSplineSegmentBasis", args, 4)
	c, _ := argColor("drawSplineSegmentBasis", args, 5, rl.RayWhite)
	rl.DrawSplineSegmentBasis(p1, p2, p3, p4, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineSegmentCatmullRom(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawSplineSegmentCatmullRom expects p1, p2, p3, p4, thick, [color]")
	}
	p1, _ := argVector2("drawSplineSegmentCatmullRom", args, 0)
	p2, _ := argVector2("drawSplineSegmentCatmullRom", args, 1)
	p3, _ := argVector2("drawSplineSegmentCatmullRom", args, 2)
	p4, _ := argVector2("drawSplineSegmentCatmullRom", args, 3)
	thick, _ := getArgFloat("drawSplineSegmentCatmullRom", args, 4)
	c, _ := argColor("drawSplineSegmentCatmullRom", args, 5, rl.RayWhite)
	rl.DrawSplineSegmentCatmullRom(p1, p2, p3, p4, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineSegmentBezierQuadratic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawSplineSegmentBezierQuadratic expects p1, c2, p3, thick, [color]")
	}
	p1, _ := argVector2("drawSplineSegmentBezierQuadratic", args, 0)
	p2, _ := argVector2("drawSplineSegmentBezierQuadratic", args, 1)
	p3, _ := argVector2("drawSplineSegmentBezierQuadratic", args, 2)
	thick, _ := getArgFloat("drawSplineSegmentBezierQuadratic", args, 3)
	c, _ := argColor("drawSplineSegmentBezierQuadratic", args, 4, rl.RayWhite)
	rl.DrawSplineSegmentBezierQuadratic(p1, p2, p3, float32(thick), c)
	return null(), nil
}

func builtinDrawSplineSegmentBezierCubic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawSplineSegmentBezierCubic expects p1, c2, c3, p4, thick, [color]")
	}
	p1, _ := argVector2("drawSplineSegmentBezierCubic", args, 0)
	p2, _ := argVector2("drawSplineSegmentBezierCubic", args, 1)
	p3, _ := argVector2("drawSplineSegmentBezierCubic", args, 2)
	p4, _ := argVector2("drawSplineSegmentBezierCubic", args, 3)
	thick, _ := getArgFloat("drawSplineSegmentBezierCubic", args, 4)
	c, _ := argColor("drawSplineSegmentBezierCubic", args, 5, rl.RayWhite)
	rl.DrawSplineSegmentBezierCubic(p1, p2, p3, p4, float32(thick), c)
	return null(), nil
}

// ---- Spline point evaluation ----

func builtinGetSplinePointLinear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getSplinePointLinear", args, 3); err != nil {
		return nil, err
	}
	p1, _ := argVector2("getSplinePointLinear", args, 0)
	p2, _ := argVector2("getSplinePointLinear", args, 1)
	t, _ := getArgFloat("getSplinePointLinear", args, 2)
	return vec2ToMap(rl.GetSplinePointLinear(p1, p2, float32(t))), nil
}

func builtinGetSplinePointBasis(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getSplinePointBasis", args, 5); err != nil {
		return nil, err
	}
	p1, _ := argVector2("getSplinePointBasis", args, 0)
	p2, _ := argVector2("getSplinePointBasis", args, 1)
	p3, _ := argVector2("getSplinePointBasis", args, 2)
	p4, _ := argVector2("getSplinePointBasis", args, 3)
	t, _ := getArgFloat("getSplinePointBasis", args, 4)
	return vec2ToMap(rl.GetSplinePointBasis(p1, p2, p3, p4, float32(t))), nil
}

func builtinGetSplinePointCatmullRom(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getSplinePointCatmullRom", args, 5); err != nil {
		return nil, err
	}
	p1, _ := argVector2("getSplinePointCatmullRom", args, 0)
	p2, _ := argVector2("getSplinePointCatmullRom", args, 1)
	p3, _ := argVector2("getSplinePointCatmullRom", args, 2)
	p4, _ := argVector2("getSplinePointCatmullRom", args, 3)
	t, _ := getArgFloat("getSplinePointCatmullRom", args, 4)
	return vec2ToMap(rl.GetSplinePointCatmullRom(p1, p2, p3, p4, float32(t))), nil
}

func builtinGetSplinePointBezierQuad(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getSplinePointBezierQuad", args, 4); err != nil {
		return nil, err
	}
	p1, _ := argVector2("getSplinePointBezierQuad", args, 0)
	p2, _ := argVector2("getSplinePointBezierQuad", args, 1)
	p3, _ := argVector2("getSplinePointBezierQuad", args, 2)
	t, _ := getArgFloat("getSplinePointBezierQuad", args, 3)
	return vec2ToMap(rl.GetSplinePointBezierQuad(p1, p2, p3, float32(t))), nil
}

func builtinGetSplinePointBezierCubic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getSplinePointBezierCubic", args, 5); err != nil {
		return nil, err
	}
	p1, _ := argVector2("getSplinePointBezierCubic", args, 0)
	p2, _ := argVector2("getSplinePointBezierCubic", args, 1)
	p3, _ := argVector2("getSplinePointBezierCubic", args, 2)
	p4, _ := argVector2("getSplinePointBezierCubic", args, 3)
	t, _ := getArgFloat("getSplinePointBezierCubic", args, 4)
	return vec2ToMap(rl.GetSplinePointBezierCubic(p1, p2, p3, p4, float32(t))), nil
}

// ---- Collision extras ----

func builtinCheckCollisionCircleLine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("checkCollisionCircleLine expects center, radius, p1, p2")
	}
	center, err := argVector2("checkCollisionCircleLine", args, 0)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("checkCollisionCircleLine", args, 1)
	p1, err := argVector2("checkCollisionCircleLine", args, 2)
	if err != nil {
		return nil, err
	}
	p2, err := argVector2("checkCollisionCircleLine", args, 3)
	if err != nil {
		return nil, err
	}
	return vBool(rl.CheckCollisionCircleLine(center, float32(r), p1, p2)), nil
}

func builtinCheckCollisionPointRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionPointRec", args, 2); err != nil {
		return nil, err
	}
	pt, err := argVector2("checkCollisionPointRec", args, 0)
	if err != nil {
		return nil, err
	}
	rec, err := argRectangle("checkCollisionPointRec", args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(rl.CheckCollisionPointRec(pt, rec)), nil
}

func builtinCheckCollisionPointCircle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("checkCollisionPointCircle expects point, center, radius")
	}
	pt, err := argVector2("checkCollisionPointCircle", args, 0)
	if err != nil {
		return nil, err
	}
	center, err := argVector2("checkCollisionPointCircle", args, 1)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("checkCollisionPointCircle", args, 2)
	return vBool(rl.CheckCollisionPointCircle(pt, center, float32(r))), nil
}

func builtinCheckCollisionPointTriangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionPointTriangle", args, 4); err != nil {
		return nil, err
	}
	pt, _ := argVector2("checkCollisionPointTriangle", args, 0)
	p1, _ := argVector2("checkCollisionPointTriangle", args, 1)
	p2, _ := argVector2("checkCollisionPointTriangle", args, 2)
	p3, _ := argVector2("checkCollisionPointTriangle", args, 3)
	return vBool(rl.CheckCollisionPointTriangle(pt, p1, p2, p3)), nil
}

func builtinCheckCollisionPointLine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("checkCollisionPointLine expects point, p1, p2, threshold")
	}
	pt, _ := argVector2("checkCollisionPointLine", args, 0)
	p1, _ := argVector2("checkCollisionPointLine", args, 1)
	p2, _ := argVector2("checkCollisionPointLine", args, 2)
	thresh, _ := argInt("checkCollisionPointLine", args, 3)
	return vBool(rl.CheckCollisionPointLine(pt, p1, p2, int32(thresh))), nil
}

func builtinCheckCollisionPointPoly(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("checkCollisionPointPoly expects point, pointsArr")
	}
	pt, err := argVector2("checkCollisionPointPoly", args, 0)
	if err != nil {
		return nil, err
	}
	pts, err := parsePointsArray("checkCollisionPointPoly", args[1])
	if err != nil {
		return nil, err
	}
	return vBool(rl.CheckCollisionPointPoly(pt, pts)), nil
}

func builtinCheckCollisionLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionLines", args, 4); err != nil {
		return nil, err
	}
	s1, _ := argVector2("checkCollisionLines", args, 0)
	e1, _ := argVector2("checkCollisionLines", args, 1)
	s2, _ := argVector2("checkCollisionLines", args, 2)
	e2, _ := argVector2("checkCollisionLines", args, 3)
	var pt rl.Vector2
	hit := rl.CheckCollisionLines(s1, e1, s2, e2, &pt)
	m := map[string]candy_evaluator.Value{
		"hit": {Kind: candy_evaluator.ValBool, B: hit},
		"point": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(pt.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(pt.Y)},
		}},
	}
	return vMap(m), nil
}

func builtinGetCollisionRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getCollisionRec", args, 2); err != nil {
		return nil, err
	}
	r1, err := argRectangle("getCollisionRec", args, 0)
	if err != nil {
		return nil, err
	}
	r2, err := argRectangle("getCollisionRec", args, 1)
	if err != nil {
		return nil, err
	}
	rec := rl.GetCollisionRec(r1, r2)
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.X)},
		"y":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.Y)},
		"width":  {Kind: candy_evaluator.ValFloat, F64: float64(rec.Width)},
		"height": {Kind: candy_evaluator.ValFloat, F64: float64(rec.Height)},
	}), nil
}

func parsePointsArray(name string, v *candy_evaluator.Value) ([]rl.Vector2, error) {
	if v.Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("%s expects array of maps {x, y}", name)
	}
	pts := make([]rl.Vector2, len(v.Elems))
	for i, e := range v.Elems {
		if e.Kind != candy_evaluator.ValMap {
			return nil, fmt.Errorf("%s: element %d is not a map {x, y}", name, i)
		}
		x := 0.0
		if vx, ok := e.StrMap["x"]; ok {
			if vx.Kind == candy_evaluator.ValInt {
				x = float64(vx.I64)
			} else {
				x = vx.F64
			}
		}
		y := 0.0
		if vy, ok := e.StrMap["y"]; ok {
			if vy.Kind == candy_evaluator.ValInt {
				y = float64(vy.I64)
			} else {
				y = vy.F64
			}
		}
		pts[i] = rl.NewVector2(float32(x), float32(y))
	}
	return pts, nil
}
