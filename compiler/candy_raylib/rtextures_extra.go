package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Color helpers ----

func colorToMap(c color.RGBA) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"r": {Kind: candy_evaluator.ValInt, I64: int64(c.R)},
		"g": {Kind: candy_evaluator.ValInt, I64: int64(c.G)},
		"b": {Kind: candy_evaluator.ValInt, I64: int64(c.B)},
		"a": {Kind: candy_evaluator.ValInt, I64: int64(c.A)},
	})
}

// argColorValue reads a color from args[i]: string name, {r,g,b,a} map, or fallback.
func argColorValue(_ string, args []*candy_evaluator.Value, i int, fallback color.RGBA) (color.RGBA, error) {
	if i >= len(args) || args[i] == nil {
		return fallback, nil
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValString {
		return colorFrom(v.Str), nil
	}
	if v.Kind == candy_evaluator.ValMap {
		m := v.StrMap
		r := uint8(mapFloat(m, "r"))
		g := uint8(mapFloat(m, "g"))
		b := uint8(mapFloat(m, "b"))
		a := uint8(255)
		if _, ok := m["a"]; ok {
			a = uint8(mapFloat(m, "a"))
		}
		return color.RGBA{R: r, G: g, B: b, A: a}, nil
	}
	return fallback, nil
}

// ---- Image loading extras ----

func builtinLoadImageRaw(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadImageRaw", args, 5); err != nil {
		return nil, err
	}
	path, err := argString("loadImageRaw", args, 0)
	if err != nil {
		return nil, err
	}
	w, _ := argInt("loadImageRaw", args, 1)
	h, _ := argInt("loadImageRaw", args, 2)
	fmt_, _ := argInt("loadImageRaw", args, 3)
	hdr, _ := argInt("loadImageRaw", args, 4)
	img := rl.LoadImageRaw(path, int32(w), int32(h), rl.PixelFormat(int32(fmt_)), int32(hdr))
	id := nextImageID
	nextImageID++
	images[id] = img
	return vInt(id), nil
}

func builtinLoadImageFromTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadImageFromTexture", args, 1); err != nil {
		return nil, err
	}
	texID, _ := argInt("loadImageFromTexture", args, 0)
	tex, ok := textures[texID]
	if !ok {
		return nil, fmt.Errorf("loadImageFromTexture: invalid texture handle %d", texID)
	}
	img := rl.LoadImageFromTexture(tex)
	id := nextImageID
	nextImageID++
	images[id] = img
	return vInt(id), nil
}

func builtinLoadImageFromScreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	img := rl.LoadImageFromScreen()
	id := nextImageID
	nextImageID++
	images[id] = img
	return vInt(id), nil
}

func builtinIsImageValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, img, err := imageByID("isImageValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsImageValid(img)), nil
}

func builtinExportImageAsCode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("exportImageAsCode", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("exportImageAsCode", args, 0)
	if err != nil {
		return nil, err
	}
	path, err := argString("exportImageAsCode", args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(rl.ExportImage(*img, path+"_code")), nil
}

func builtinExportImageToMemory(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("exportImageToMemory", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("exportImageToMemory", args, 0)
	if err != nil {
		return nil, err
	}
	ft, err := argString("exportImageToMemory", args, 1)
	if err != nil {
		return nil, err
	}
	data := rl.ExportImageToMemory(*img, ft)
	elems := make([]candy_evaluator.Value, len(data))
	for i, b := range data {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValInt, I64: int64(b)}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

// ---- Image manipulation extras ----

func builtinImageFromChannel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("imageFromChannel", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("imageFromChannel", args, 0)
	if err != nil {
		return nil, err
	}
	ch, _ := argInt("imageFromChannel", args, 1)
	result := rl.ImageFromChannel(*img, int32(ch))
	id := nextImageID
	nextImageID++
	images[id] = &result
	return vInt(id), nil
}

func builtinImageResizeNN(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("imageResizeNN expects imageId, newWidth, newHeight")
	}
	_, img, err := imageByID("imageResizeNN", args, 0)
	if err != nil {
		return nil, err
	}
	w, _ := argInt("imageResizeNN", args, 1)
	h, _ := argInt("imageResizeNN", args, 2)
	rl.ImageResizeNN(img, int32(w), int32(h))
	return null(), nil
}

func builtinImageKernelConvolution(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("imageKernelConvolution expects imageId, kernelArr")
	}
	_, img, err := imageByID("imageKernelConvolution", args, 0)
	if err != nil {
		return nil, err
	}
	kv := args[1]
	if kv.Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("imageKernelConvolution: kernel must be array of floats")
	}
	kernel := make([]float32, len(kv.Elems))
	for i, e := range kv.Elems {
		switch e.Kind {
		case candy_evaluator.ValFloat:
			kernel[i] = float32(e.F64)
		case candy_evaluator.ValInt:
			kernel[i] = float32(e.I64)
		}
	}
	rl.ImageKernelConvolution(img, kernel)
	return null(), nil
}

func builtinImageRotateCW(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("imageRotateCW", args, 1); err != nil {
		return nil, err
	}
	_, img, err := imageByID("imageRotateCW", args, 0)
	if err != nil {
		return nil, err
	}
	rl.ImageRotateCW(img)
	return null(), nil
}

func builtinImageRotateCCW(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("imageRotateCCW", args, 1); err != nil {
		return nil, err
	}
	_, img, err := imageByID("imageRotateCCW", args, 0)
	if err != nil {
		return nil, err
	}
	rl.ImageRotateCCW(img)
	return null(), nil
}

func builtinLoadImageColors(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadImageColors", args, 1); err != nil {
		return nil, err
	}
	_, img, err := imageByID("loadImageColors", args, 0)
	if err != nil {
		return nil, err
	}
	cols := rl.LoadImageColors(img)
	elems := make([]candy_evaluator.Value, len(cols))
	for i, c := range cols {
		elems[i] = *colorToMap(c)
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

func builtinLoadImagePalette(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// loadImagePalette is not exposed under the cgo/raylib build tag.
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: nil}, nil
}

func builtinGetImageAlphaBorder(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// getImageAlphaBorder is not exposed under the cgo/raylib build tag; returns zero rect.
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: 0},
		"y":      {Kind: candy_evaluator.ValFloat, F64: 0},
		"width":  {Kind: candy_evaluator.ValFloat, F64: 0},
		"height": {Kind: candy_evaluator.ValFloat, F64: 0},
	}), nil
}

func builtinGetImageColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getImageColor", args, 3); err != nil {
		return nil, err
	}
	_, img, err := imageByID("getImageColor", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argInt("getImageColor", args, 1)
	y, _ := argInt("getImageColor", args, 2)
	return colorToMap(rl.GetImageColor(*img, int32(x), int32(y))), nil
}

// ---- Image drawing functions ----

func builtinImageClearBackground(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("imageClearBackground expects imageId, color")
	}
	_, img, err := imageByID("imageClearBackground", args, 0)
	if err != nil {
		return nil, err
	}
	c, _ := argColorValue("imageClearBackground", args, 1, rl.Black)
	rl.ImageClearBackground(img, c)
	return null(), nil
}

func builtinImageDrawPixel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawPixel expects imageId, x, y, color")
	}
	_, img, err := imageByID("imageDrawPixel", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argInt("imageDrawPixel", args, 1)
	y, _ := argInt("imageDrawPixel", args, 2)
	c, _ := argColorValue("imageDrawPixel", args, 3, rl.RayWhite)
	rl.ImageDrawPixel(img, int32(x), int32(y), c)
	return null(), nil
}

func builtinImageDrawPixelV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("imageDrawPixelV expects imageId, position, color")
	}
	_, img, err := imageByID("imageDrawPixelV", args, 0)
	if err != nil {
		return nil, err
	}
	pos, err2 := argVector2("imageDrawPixelV", args, 1)
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("imageDrawPixelV", args, 2, rl.RayWhite)
	rl.ImageDrawPixelV(img, pos, c)
	return null(), nil
}

func builtinImageDrawLine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("imageDrawLine expects imageId, x1, y1, x2, y2, color")
	}
	_, img, err := imageByID("imageDrawLine", args, 0)
	if err != nil {
		return nil, err
	}
	x1, _ := argInt("imageDrawLine", args, 1)
	y1, _ := argInt("imageDrawLine", args, 2)
	x2, _ := argInt("imageDrawLine", args, 3)
	y2, _ := argInt("imageDrawLine", args, 4)
	c, _ := argColorValue("imageDrawLine", args, 5, rl.RayWhite)
	rl.ImageDrawLine(img, int32(x1), int32(y1), int32(x2), int32(y2), c)
	return null(), nil
}

func builtinImageDrawLineV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawLineV expects imageId, start, end, color")
	}
	_, img, err := imageByID("imageDrawLineV", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := argVector2("imageDrawLineV", args, 1)
	e, _ := argVector2("imageDrawLineV", args, 2)
	c, _ := argColorValue("imageDrawLineV", args, 3, rl.RayWhite)
	rl.ImageDrawLineV(img, s, e, c)
	return null(), nil
}

func builtinImageDrawLineEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("imageDrawLineEx expects imageId, start, end, thick, color")
	}
	_, img, err := imageByID("imageDrawLineEx", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := argVector2("imageDrawLineEx", args, 1)
	e, _ := argVector2("imageDrawLineEx", args, 2)
	thick, _ := argInt("imageDrawLineEx", args, 3)
	c, _ := argColorValue("imageDrawLineEx", args, 4, rl.RayWhite)
	rl.ImageDrawLineEx(img, s, e, int32(thick), c)
	return null(), nil
}

func builtinImageDrawCircle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("imageDrawCircle expects imageId, centerX, centerY, radius, color")
	}
	_, img, err := imageByID("imageDrawCircle", args, 0)
	if err != nil {
		return nil, err
	}
	cx, _ := argInt("imageDrawCircle", args, 1)
	cy, _ := argInt("imageDrawCircle", args, 2)
	r, _ := argInt("imageDrawCircle", args, 3)
	c, _ := argColorValue("imageDrawCircle", args, 4, rl.RayWhite)
	rl.ImageDrawCircle(img, int32(cx), int32(cy), int32(r), c)
	return null(), nil
}

func builtinImageDrawCircleV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawCircleV expects imageId, center, radius, color")
	}
	_, img, err := imageByID("imageDrawCircleV", args, 0)
	if err != nil {
		return nil, err
	}
	center, _ := argVector2("imageDrawCircleV", args, 1)
	r, _ := argInt("imageDrawCircleV", args, 2)
	c, _ := argColorValue("imageDrawCircleV", args, 3, rl.RayWhite)
	rl.ImageDrawCircleV(img, center, int32(r), c)
	return null(), nil
}

func builtinImageDrawCircleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("imageDrawCircleLines expects imageId, centerX, centerY, radius, color")
	}
	_, img, err := imageByID("imageDrawCircleLines", args, 0)
	if err != nil {
		return nil, err
	}
	cx, _ := argInt("imageDrawCircleLines", args, 1)
	cy, _ := argInt("imageDrawCircleLines", args, 2)
	r, _ := argInt("imageDrawCircleLines", args, 3)
	c, _ := argColorValue("imageDrawCircleLines", args, 4, rl.RayWhite)
	rl.ImageDrawCircleLines(img, int32(cx), int32(cy), int32(r), c)
	return null(), nil
}

func builtinImageDrawCircleLinesV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawCircleLinesV expects imageId, center, radius, color")
	}
	_, img, err := imageByID("imageDrawCircleLinesV", args, 0)
	if err != nil {
		return nil, err
	}
	center, _ := argVector2("imageDrawCircleLinesV", args, 1)
	r, _ := argInt("imageDrawCircleLinesV", args, 2)
	c, _ := argColorValue("imageDrawCircleLinesV", args, 3, rl.RayWhite)
	rl.ImageDrawCircleLinesV(img, center, int32(r), c)
	return null(), nil
}

func builtinImageDrawRectangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("imageDrawRectangle expects imageId, x, y, w, h, color")
	}
	_, img, err := imageByID("imageDrawRectangle", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argInt("imageDrawRectangle", args, 1)
	y, _ := argInt("imageDrawRectangle", args, 2)
	w, _ := argInt("imageDrawRectangle", args, 3)
	h, _ := argInt("imageDrawRectangle", args, 4)
	c, _ := argColorValue("imageDrawRectangle", args, 5, rl.RayWhite)
	rl.ImageDrawRectangle(img, int32(x), int32(y), int32(w), int32(h), c)
	return null(), nil
}

func builtinImageDrawRectangleV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawRectangleV expects imageId, position, size, color")
	}
	_, img, err := imageByID("imageDrawRectangleV", args, 0)
	if err != nil {
		return nil, err
	}
	pos, _ := argVector2("imageDrawRectangleV", args, 1)
	size, _ := argVector2("imageDrawRectangleV", args, 2)
	c, _ := argColorValue("imageDrawRectangleV", args, 3, rl.RayWhite)
	rl.ImageDrawRectangleV(img, pos, size, c)
	return null(), nil
}

func builtinImageDrawRectangleRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("imageDrawRectangleRec expects imageId, rec, color")
	}
	_, img, err := imageByID("imageDrawRectangleRec", args, 0)
	if err != nil {
		return nil, err
	}
	rec, err2 := argRectangle("imageDrawRectangleRec", args, 1)
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("imageDrawRectangleRec", args, 2, rl.RayWhite)
	rl.ImageDrawRectangleRec(img, rec, c)
	return null(), nil
}

func builtinImageDrawRectangleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("imageDrawRectangleLines expects imageId, rec, thick, color")
	}
	_, img, err := imageByID("imageDrawRectangleLines", args, 0)
	if err != nil {
		return nil, err
	}
	rec, err2 := argRectangle("imageDrawRectangleLines", args, 1)
	if err2 != nil {
		return nil, err2
	}
	thick, _ := argInt("imageDrawRectangleLines", args, 2)
	c, _ := argColorValue("imageDrawRectangleLines", args, 3, rl.RayWhite)
	rl.ImageDrawRectangleLines(img, rec, int(thick), c)
	return null(), nil
}

func builtinImageDrawTriangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("imageDrawTriangle expects imageId, v1, v2, v3, color")
	}
	_, img, err := imageByID("imageDrawTriangle", args, 0)
	if err != nil {
		return nil, err
	}
	v1, _ := argVector2("imageDrawTriangle", args, 1)
	v2, _ := argVector2("imageDrawTriangle", args, 2)
	v3, _ := argVector2("imageDrawTriangle", args, 3)
	c, _ := argColorValue("imageDrawTriangle", args, 4, rl.RayWhite)
	rl.ImageDrawTriangle(img, v1, v2, v3, c)
	return null(), nil
}

func builtinImageDrawTriangleEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("imageDrawTriangleEx expects imageId, v1, v2, v3, c1, c2, c3")
	}
	_, img, err := imageByID("imageDrawTriangleEx", args, 0)
	if err != nil {
		return nil, err
	}
	v1, _ := argVector2("imageDrawTriangleEx", args, 1)
	v2, _ := argVector2("imageDrawTriangleEx", args, 2)
	v3, _ := argVector2("imageDrawTriangleEx", args, 3)
	c1, _ := argColorValue("imageDrawTriangleEx", args, 4, rl.Red)
	c2, _ := argColorValue("imageDrawTriangleEx", args, 5, rl.Green)
	c3, _ := argColorValue("imageDrawTriangleEx", args, 6, rl.Blue)
	rl.ImageDrawTriangleEx(img, v1, v2, v3, c1, c2, c3)
	return null(), nil
}

func builtinImageDrawTriangleLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("imageDrawTriangleLines expects imageId, v1, v2, v3, color")
	}
	_, img, err := imageByID("imageDrawTriangleLines", args, 0)
	if err != nil {
		return nil, err
	}
	v1, _ := argVector2("imageDrawTriangleLines", args, 1)
	v2, _ := argVector2("imageDrawTriangleLines", args, 2)
	v3, _ := argVector2("imageDrawTriangleLines", args, 3)
	c, _ := argColorValue("imageDrawTriangleLines", args, 4, rl.RayWhite)
	rl.ImageDrawTriangleLines(img, v1, v2, v3, c)
	return null(), nil
}

func builtinImageDrawTriangleFan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("imageDrawTriangleFan expects imageId, pointsArr, color")
	}
	_, img, err := imageByID("imageDrawTriangleFan", args, 0)
	if err != nil {
		return nil, err
	}
	pts, err2 := parsePointsArray("imageDrawTriangleFan", args[1])
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("imageDrawTriangleFan", args, 2, rl.RayWhite)
	rl.ImageDrawTriangleFan(img, pts, c)
	return null(), nil
}

func builtinImageDrawTriangleStrip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("imageDrawTriangleStrip expects imageId, pointsArr, color")
	}
	_, img, err := imageByID("imageDrawTriangleStrip", args, 0)
	if err != nil {
		return nil, err
	}
	pts, err2 := parsePointsArray("imageDrawTriangleStrip", args[1])
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("imageDrawTriangleStrip", args, 2, rl.RayWhite)
	rl.ImageDrawTriangleStrip(img, pts, c)
	return null(), nil
}

func builtinImageDrawText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("imageDrawText expects imageId, text, x, y, fontSize, color")
	}
	_, img, err := imageByID("imageDrawText", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("imageDrawText", args, 1)
	x, _ := argInt("imageDrawText", args, 2)
	y, _ := argInt("imageDrawText", args, 3)
	sz, _ := argInt("imageDrawText", args, 4)
	c, _ := argColorValue("imageDrawText", args, 5, rl.RayWhite)
	rl.ImageDrawText(img, int32(x), int32(y), text, int32(sz), c)
	return null(), nil
}

func builtinImageDrawTextEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("imageDrawTextEx expects imageId, fontId, text, position, fontSize, spacing, [tint]")
	}
	_, img, err := imageByID("imageDrawTextEx", args, 0)
	if err != nil {
		return nil, err
	}
	fontID, _ := argInt("imageDrawTextEx", args, 1)
	fnt, ok := fonts[fontID]
	if !ok {
		return nil, fmt.Errorf("imageDrawTextEx: invalid font handle %d", fontID)
	}
	text, _ := argString("imageDrawTextEx", args, 2)
	pos, _ := argVector2("imageDrawTextEx", args, 3)
	sz, _ := getArgFloat("imageDrawTextEx", args, 4)
	spacing, _ := getArgFloat("imageDrawTextEx", args, 5)
	c, _ := argColorValue("imageDrawTextEx", args, 6, rl.RayWhite)
	rl.ImageDrawTextEx(img, pos, fnt, text, float32(sz), float32(spacing), c)
	return null(), nil
}

// ---- Texture extras ----

func builtinLoadTextureCubemap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadTextureCubemap", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("loadTextureCubemap", args, 0)
	if err != nil {
		return nil, err
	}
	layout, _ := argInt("loadTextureCubemap", args, 1)
	tex := rl.LoadTextureCubemap(img, int32(layout))
	id := nextTextureID
	nextTextureID++
	textures[id] = tex
	return vInt(id), nil
}

func builtinIsTextureValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, tex, err := textureByID("isTextureValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsTextureValid(tex)), nil
}

func builtinIsRenderTextureValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, rt, err := renderTextureByID("isRenderTextureValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsRenderTextureValid(rt)), nil
}

func builtinGenTextureMipmaps(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genTextureMipmaps", args, 1); err != nil {
		return nil, err
	}
	id, tex, err := textureByID("genTextureMipmaps", args, 0)
	if err != nil {
		return nil, err
	}
	rl.GenTextureMipmaps(&tex)
	textures[id] = tex
	return null(), nil
}

func builtinSetTextureFilter(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setTextureFilter", args, 2); err != nil {
		return nil, err
	}
	_, tex, err := textureByID("setTextureFilter", args, 0)
	if err != nil {
		return nil, err
	}
	mode, _ := argInt("setTextureFilter", args, 1)
	rl.SetTextureFilter(tex, rl.TextureFilterMode(int32(mode)))
	return null(), nil
}

func builtinSetTextureWrap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setTextureWrap", args, 2); err != nil {
		return nil, err
	}
	_, tex, err := textureByID("setTextureWrap", args, 0)
	if err != nil {
		return nil, err
	}
	mode, _ := argInt("setTextureWrap", args, 1)
	rl.SetTextureWrap(tex, rl.TextureWrapMode(int32(mode)))
	return null(), nil
}

func builtinDrawTextureV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawTextureV expects textureId, position, [tint]")
	}
	_, tex, err := textureByID("drawTextureV", args, 0)
	if err != nil {
		return nil, err
	}
	pos, err2 := argVector2("drawTextureV", args, 1)
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("drawTextureV", args, 2, rl.RayWhite)
	rl.DrawTextureV(tex, pos, c)
	return null(), nil
}

func builtinDrawTexturePro(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawTexturePro expects textureId, source, dest, origin, rotation, [tint]")
	}
	_, tex, err := textureByID("drawTexturePro", args, 0)
	if err != nil {
		return nil, err
	}
	src, err2 := argRectangle("drawTexturePro", args, 1)
	if err2 != nil {
		return nil, err2
	}
	dst, err3 := argRectangle("drawTexturePro", args, 2)
	if err3 != nil {
		return nil, err3
	}
	origin, _ := argVector2("drawTexturePro", args, 3)
	rot, _ := getArgFloat("drawTexturePro", args, 4)
	c, _ := argColorValue("drawTexturePro", args, 5, rl.RayWhite)
	rl.DrawTexturePro(tex, src, dst, origin, float32(rot), c)
	return null(), nil
}

func builtinDrawTextureNPatch(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawTextureNPatch expects textureId, nPatchInfo, dest, origin, rotation, [tint]")
	}
	_, tex, err := textureByID("drawTextureNPatch", args, 0)
	if err != nil {
		return nil, err
	}
	// nPatchInfo map: {sourceX,sourceY,sourceW,sourceH,left,top,right,bottom,layout}
	if args[1].Kind != candy_evaluator.ValMap {
		return nil, fmt.Errorf("drawTextureNPatch: arg 2 must be nPatchInfo map")
	}
	nm := args[1].StrMap
	npi := rl.NPatchInfo{
		Source: rl.NewRectangle(
			mapFloat(nm, "sourceX"), mapFloat(nm, "sourceY"),
			mapFloat(nm, "sourceW"), mapFloat(nm, "sourceH"),
		),
		Left:   int32(mapFloat(nm, "left")),
		Top:    int32(mapFloat(nm, "top")),
		Right:  int32(mapFloat(nm, "right")),
		Bottom: int32(mapFloat(nm, "bottom")),
		Layout: rl.NPatchLayout(int32(mapFloat(nm, "layout"))),
	}
	dst, err2 := argRectangle("drawTextureNPatch", args, 2)
	if err2 != nil {
		return nil, err2
	}
	origin, _ := argVector2("drawTextureNPatch", args, 3)
	rot, _ := getArgFloat("drawTextureNPatch", args, 4)
	c, _ := argColorValue("drawTextureNPatch", args, 5, rl.RayWhite)
	rl.DrawTextureNPatch(tex, npi, dst, origin, float32(rot), c)
	return null(), nil
}

// ---- Color/pixel functions ----

func builtinColorIsEqual(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorIsEqual", args, 2); err != nil {
		return nil, err
	}
	c1, _ := argColorValue("colorIsEqual", args, 0, rl.Black)
	c2, _ := argColorValue("colorIsEqual", args, 1, rl.Black)
	return vBool(c1 == c2), nil
}

func builtinFade(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("fade", args, 2); err != nil {
		return nil, err
	}
	c, _ := argColorValue("fade", args, 0, rl.RayWhite)
	alpha, _ := getArgFloat("fade", args, 1)
	return colorToMap(rl.Fade(c, float32(alpha))), nil
}

func builtinColorToInt(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorToInt", args, 1); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorToInt", args, 0, rl.RayWhite)
	return vInt(int64(rl.ColorToInt(c))), nil
}

func builtinColorNormalize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorNormalize", args, 1); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorNormalize", args, 0, rl.RayWhite)
	v4 := rl.ColorNormalize(c)
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(v4.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(v4.Y)},
		"z": {Kind: candy_evaluator.ValFloat, F64: float64(v4.Z)},
		"w": {Kind: candy_evaluator.ValFloat, F64: float64(v4.W)},
	}), nil
}

func builtinColorFromNormalized(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorFromNormalized", args, 1); err != nil {
		return nil, err
	}
	v4, err := argQuat("colorFromNormalized", args, 0)
	if err != nil {
		return nil, err
	}
	return colorToMap(rl.ColorFromNormalized(v4)), nil
}

func builtinColorToHSV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorToHSV", args, 1); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorToHSV", args, 0, rl.RayWhite)
	v3 := rl.ColorToHSV(c)
	return vec3ToMap(v3), nil
}

func builtinColorFromHSV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorFromHSV", args, 3); err != nil {
		return nil, err
	}
	h, _ := getArgFloat("colorFromHSV", args, 0)
	s, _ := getArgFloat("colorFromHSV", args, 1)
	v, _ := getArgFloat("colorFromHSV", args, 2)
	return colorToMap(rl.ColorFromHSV(float32(h), float32(s), float32(v))), nil
}

func builtinColorTint(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorTint", args, 2); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorTint", args, 0, rl.RayWhite)
	tint, _ := argColorValue("colorTint", args, 1, rl.RayWhite)
	return colorToMap(rl.ColorTint(c, tint)), nil
}

func builtinColorBrightness(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorBrightness", args, 2); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorBrightness", args, 0, rl.RayWhite)
	f, _ := getArgFloat("colorBrightness", args, 1)
	return colorToMap(rl.ColorBrightness(c, float32(f))), nil
}

func builtinColorContrast(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorContrast", args, 2); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorContrast", args, 0, rl.RayWhite)
	f, _ := getArgFloat("colorContrast", args, 1)
	return colorToMap(rl.ColorContrast(c, float32(f))), nil
}

func builtinColorAlpha(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorAlpha", args, 2); err != nil {
		return nil, err
	}
	c, _ := argColorValue("colorAlpha", args, 0, rl.RayWhite)
	a, _ := getArgFloat("colorAlpha", args, 1)
	return colorToMap(rl.ColorAlpha(c, float32(a))), nil
}

func builtinColorAlphaBlend(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorAlphaBlend", args, 3); err != nil {
		return nil, err
	}
	dst, _ := argColorValue("colorAlphaBlend", args, 0, rl.Black)
	src, _ := argColorValue("colorAlphaBlend", args, 1, rl.RayWhite)
	tint, _ := argColorValue("colorAlphaBlend", args, 2, rl.RayWhite)
	return colorToMap(rl.ColorAlphaBlend(src, dst, tint)), nil
}

func builtinColorLerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("colorLerp", args, 3); err != nil {
		return nil, err
	}
	c1, _ := argColorValue("colorLerp", args, 0, rl.Black)
	c2, _ := argColorValue("colorLerp", args, 1, rl.RayWhite)
	f, _ := getArgFloat("colorLerp", args, 2)
	return colorToMap(rl.ColorLerp(c1, c2, float32(f))), nil
}

func builtinGetColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getColor", args, 1); err != nil {
		return nil, err
	}
	hex, _ := argInt("getColor", args, 0)
	return colorToMap(rl.GetColor(uint(hex))), nil
}

func builtinGetPixelDataSize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getPixelDataSize", args, 3); err != nil {
		return nil, err
	}
	w, _ := argInt("getPixelDataSize", args, 0)
	h, _ := argInt("getPixelDataSize", args, 1)
	fmt_, _ := argInt("getPixelDataSize", args, 2)
	return vInt(int64(rl.GetPixelDataSize(int32(w), int32(h), int32(fmt_)))), nil
}
