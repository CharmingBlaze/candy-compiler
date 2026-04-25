package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strings"
)

// imagePathToTexID caches path -> handle for Candy image() helper.
var imagePathToTexID = map[string]int64{}

func textureByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Texture2D, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return 0, rl.Texture2D{}, err
	}
	t, ok := textures[id]
	if !ok {
		return 0, rl.Texture2D{}, fmt.Errorf("%s: invalid texture handle %d", name, id)
	}
	return id, t, nil
}

func builtinLoadTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadTexture", args, 1); err != nil {
		return nil, err
	}
	path, err := argString("loadTexture", args, 0)
	if err != nil {
		return nil, err
	}
	tex := rl.LoadTexture(path)
	id := nextTextureID
	nextTextureID++
	textures[id] = tex
	return vInt(id), nil
}

// builtinImage draws a file-backed texture: image("pic.png", x, y) — loads once, then reuses the handle.
func builtinImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("image expects path, x, y, [tint color]")
	}
	path, err := argString("image", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argInt("image", args, 1)
	y, _ := argInt("image", args, 2)
	c, _ := argColor("image", args, 3, rl.White)
	id, ok := imagePathToTexID[path]
	if !ok {
		tex := rl.LoadTexture(path)
		id = nextTextureID
		nextTextureID++
		textures[id] = tex
		imagePathToTexID[path] = id
	}
	tex := textures[id]
	rl.DrawTexture(tex, int32(x), int32(y), c)
	return null(), nil
}

func builtinUnloadTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, tex, err := textureByID("unloadTexture", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadTexture(tex)
	delete(textures, id)
	return null(), nil
}

func builtinIsTextureReady(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, _, err := textureByID("isTextureReady", args, 0)
	if err != nil {
		return nil, err
	}
	return vBool(true), nil
}

func builtinDrawTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 { return nil, fmt.Errorf("drawTexture expects textureId, x, y, [color]") }
	_, tex, err := textureByID("drawTexture", args, 0)
	if err != nil { return nil, err }
	x, _ := argInt("drawTexture", args, 1)
	y, _ := argInt("drawTexture", args, 2)
	c, _ := argColor("drawTexture", args, 3, rl.White)
	rl.DrawTexture(tex, int32(x), int32(y), c)
	return null(), nil
}

func builtinDrawTextureEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 { return nil, fmt.Errorf("drawTextureEx expects textureId, x, y, rotation, scale, [color]") }
	_, tex, err := textureByID("drawTextureEx", args, 0)
	if err != nil { return nil, err }
	x, _ := getArgFloat("drawTextureEx", args, 1)
	y, _ := getArgFloat("drawTextureEx", args, 2)
	rot, _ := getArgFloat("drawTextureEx", args, 3)
	sc, _ := getArgFloat("drawTextureEx", args, 4)
	c, _ := argColor("drawTextureEx", args, 5, rl.White)
	rl.DrawTextureEx(tex, rl.NewVector2(float32(x), float32(y)), float32(rot), float32(sc), c)
	return null(), nil
}

func builtinDrawTextureRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 6 { return nil, fmt.Errorf("drawTextureRec expects textureId, sx, sy, sw, sh, x, y, [color]") }
    _, tex, err := textureByID("drawTextureRec", args, 0)
    if err != nil { return nil, err }
    sx, _ := getArgFloat("drawTextureRec", args, 1)
    sy, _ := getArgFloat("drawTextureRec", args, 2)
    sw, _ := getArgFloat("drawTextureRec", args, 3)
    sh, _ := getArgFloat("drawTextureRec", args, 4)
    x, _ := getArgFloat("drawTextureRec", args, 5)
    y, _ := getArgFloat("drawTextureRec", args, 6)
    c, _ := argColor("drawTextureRec", args, 7, rl.White)
    rl.DrawTextureRec(tex, rl.NewRectangle(float32(sx), float32(sy), float32(sw), float32(sh)), rl.NewVector2(float32(x), float32(y)), c)
    return null(), nil
}

// ---- Render Textures ----

func renderTextureByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.RenderTexture2D, error) {
    id, err := argInt(name, args, i)
    if err != nil { return 0, rl.RenderTexture2D{}, err }
    rt, ok := renderTextures[id]
    if !ok { return 0, rl.RenderTexture2D{}, fmt.Errorf("%s: invalid render texture handle %d", name, id) }
    return id, rt, nil
}

func builtinLoadRenderTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("loadRenderTexture", args, 2); err != nil { return nil, err }
    w, _ := argInt("loadRenderTexture", args, 0)
    h, _ := argInt("loadRenderTexture", args, 1)
    rt := rl.LoadRenderTexture(int32(w), int32(h))
    id := nextRenderTextureID
    nextRenderTextureID++
    renderTextures[id] = rt
    return vInt(id), nil
}

func builtinUnloadRenderTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    id, rt, err := renderTextureByID("unloadRenderTexture", args, 0)
    if err != nil { return nil, err }
    rl.UnloadRenderTexture(rt)
    delete(renderTextures, id)
    return null(), nil
}

func builtinBeginTextureMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, rt, err := renderTextureByID("beginTextureMode", args, 0)
    if err != nil { return nil, err }
    rl.BeginTextureMode(rt)
    return null(), nil
}

func builtinEndTextureMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    rl.EndTextureMode()
    return null(), nil
}

func builtinGetRenderTextureTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, rt, err := renderTextureByID("getRenderTextureTexture", args, 0)
    if err != nil { return nil, err }
    // We need to store this texture in the textures map to return a handle
    tid := nextTextureID
    nextTextureID++
    textures[tid] = rt.Texture
    return vInt(tid), nil
}

// ---- Images ----

func imageByID(name string, args []*candy_evaluator.Value, i int) (int64, *rl.Image, error) {
    id, err := argInt(name, args, i)
    if err != nil { return 0, nil, err }
    img, ok := images[id]
    if !ok { return 0, nil, fmt.Errorf("%s: invalid image handle %d", name, id) }
    return id, img, nil
}

func builtinLoadImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("loadImage", args, 1); err != nil { return nil, err }
    path, err := argString("loadImage", args, 0)
    if err != nil { return nil, err }
    img := rl.LoadImage(path)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinUnloadImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    id, img, err := imageByID("unloadImage", args, 0)
    if err != nil { return nil, err }
    rl.UnloadImage(img)
    delete(images, id)
    return null(), nil
}

func builtinLoadTextureFromImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, img, err := imageByID("loadTextureFromImage", args, 0)
    if err != nil { return nil, err }
    tex := rl.LoadTextureFromImage(img)
    id := nextTextureID
    nextTextureID++
    textures[id] = tex
    return vInt(id), nil
}

func builtinImageResize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 3 { return nil, fmt.Errorf("imageResize expects imageId, newWidth, newHeight") }
    _, img, err := imageByID("imageResize", args, 0)
    if err != nil { return nil, err }
    w, _ := argInt("imageResize", args, 1)
    h, _ := argInt("imageResize", args, 2)
    rl.ImageResize(img, int32(w), int32(h))
    return null(), nil
}

func builtinImageColorTint(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 2 { return nil, fmt.Errorf("imageColorTint expects imageId, color") }
    _, img, err := imageByID("imageColorTint", args, 0)
    if err != nil { return nil, err }
    c, _ := argColor("imageColorTint", args, 1, rl.White)
    rl.ImageColorTint(img, c)
    return null(), nil
}

func builtinImageCrop(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("imageCrop expects imageId, x, y, width, height") }
    _, img, err := imageByID("imageCrop", args, 0)
    if err != nil { return nil, err }
    x, _ := getArgFloat("imageCrop", args, 1)
    y, _ := getArgFloat("imageCrop", args, 2)
    w, _ := getArgFloat("imageCrop", args, 3)
    h, _ := getArgFloat("imageCrop", args, 4)
    rl.ImageCrop(img, rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)))
    return null(), nil
}

func builtinImageFlipVertical(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageFlipVertical", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageFlipVertical", args, 0)
    if err != nil { return nil, err }
    rl.ImageFlipVertical(img)
    return null(), nil
}

func builtinImageFlipHorizontal(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageFlipHorizontal", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageFlipHorizontal", args, 0)
    if err != nil { return nil, err }
    rl.ImageFlipHorizontal(img)
    return null(), nil
}

func builtinImageRotate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageRotate", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageRotate", args, 0)
    if err != nil { return nil, err }
    deg, _ := argInt("imageRotate", args, 1)
    rl.ImageRotate(img, int32(deg))
    return null(), nil
}

func builtinImageResizeCanvas(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("imageResizeCanvas expects imageId, width, height, offsetX, offsetY, [fillColor]") }
    _, img, err := imageByID("imageResizeCanvas", args, 0)
    if err != nil { return nil, err }
    w, _ := argInt("imageResizeCanvas", args, 1)
    h, _ := argInt("imageResizeCanvas", args, 2)
    ox, _ := argInt("imageResizeCanvas", args, 3)
    oy, _ := argInt("imageResizeCanvas", args, 4)
    fill, _ := argColor("imageResizeCanvas", args, 5, rl.Black)
    rl.ImageResizeCanvas(img, int32(w), int32(h), int32(ox), int32(oy), fill)
    return null(), nil
}

func builtinImageAlphaCrop(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 2 { return nil, fmt.Errorf("imageAlphaCrop expects imageId, threshold") }
    _, img, err := imageByID("imageAlphaCrop", args, 0)
    if err != nil { return nil, err }
    threshold, _ := getArgFloat("imageAlphaCrop", args, 1)
    rl.ImageAlphaCrop(img, float32(threshold))
    return null(), nil
}

func builtinImageAlphaClear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 3 { return nil, fmt.Errorf("imageAlphaClear expects imageId, color, threshold") }
    _, img, err := imageByID("imageAlphaClear", args, 0)
    if err != nil { return nil, err }
    c, _ := argColor("imageAlphaClear", args, 1, rl.Black)
    threshold, _ := getArgFloat("imageAlphaClear", args, 2)
    rl.ImageAlphaClear(img, c, float32(threshold))
    return null(), nil
}

func builtinImageAlphaMask(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageAlphaMask", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageAlphaMask", args, 0)
    if err != nil { return nil, err }
    _, mask, err := imageByID("imageAlphaMask", args, 1)
    if err != nil { return nil, err }
    rl.ImageAlphaMask(img, mask)
    return null(), nil
}

func builtinImageAlphaPremultiply(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageAlphaPremultiply", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageAlphaPremultiply", args, 0)
    if err != nil { return nil, err }
    rl.ImageAlphaPremultiply(img)
    return null(), nil
}

func builtinImageBlurGaussian(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageBlurGaussian", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageBlurGaussian", args, 0)
    if err != nil { return nil, err }
    blurSize, _ := argInt("imageBlurGaussian", args, 1)
    rl.ImageBlurGaussian(img, int32(blurSize))
    return null(), nil
}

func builtinImageColorInvert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageColorInvert", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageColorInvert", args, 0)
    if err != nil { return nil, err }
    rl.ImageColorInvert(img)
    return null(), nil
}

func builtinImageColorGrayscale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageColorGrayscale", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageColorGrayscale", args, 0)
    if err != nil { return nil, err }
    rl.ImageColorGrayscale(img)
    return null(), nil
}

func builtinImageColorContrast(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageColorContrast", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageColorContrast", args, 0)
    if err != nil { return nil, err }
    contrast, _ := getArgFloat("imageColorContrast", args, 1)
    rl.ImageColorContrast(img, float32(contrast))
    return null(), nil
}

func builtinImageColorBrightness(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageColorBrightness", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageColorBrightness", args, 0)
    if err != nil { return nil, err }
    brightness, _ := argInt("imageColorBrightness", args, 1)
    rl.ImageColorBrightness(img, int32(brightness))
    return null(), nil
}

func builtinImageColorReplace(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 3 { return nil, fmt.Errorf("imageColorReplace expects imageId, color, replaceColor") }
    _, img, err := imageByID("imageColorReplace", args, 0)
    if err != nil { return nil, err }
    from, _ := argColor("imageColorReplace", args, 1, rl.White)
    to, _ := argColor("imageColorReplace", args, 2, rl.Black)
    rl.ImageColorReplace(img, from, to)
    return null(), nil
}

func builtinImageMipmaps(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageMipmaps", args, 1); err != nil { return nil, err }
    _, img, err := imageByID("imageMipmaps", args, 0)
    if err != nil { return nil, err }
    rl.ImageMipmaps(img)
    return null(), nil
}

func builtinImageDither(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageDither", args, 5); err != nil { return nil, err }
    _, img, err := imageByID("imageDither", args, 0)
    if err != nil { return nil, err }
    rbpp, _ := argInt("imageDither", args, 1)
    gbpp, _ := argInt("imageDither", args, 2)
    bbpp, _ := argInt("imageDither", args, 3)
    abpp, _ := argInt("imageDither", args, 4)
    rl.ImageDither(img, int32(rbpp), int32(gbpp), int32(bbpp), int32(abpp))
    return null(), nil
}

func builtinImageCopy(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageCopy", args, 1); err != nil { return nil, err }
    _, src, err := imageByID("imageCopy", args, 0)
    if err != nil { return nil, err }
    img := rl.ImageCopy(src)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinImageFromImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("imageFromImage expects imageId, x, y, width, height") }
    _, src, err := imageByID("imageFromImage", args, 0)
    if err != nil { return nil, err }
    x, _ := getArgFloat("imageFromImage", args, 1)
    y, _ := getArgFloat("imageFromImage", args, 2)
    w, _ := getArgFloat("imageFromImage", args, 3)
    h, _ := getArgFloat("imageFromImage", args, 4)
    img := rl.ImageFromImage(*src, rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)))
    id := nextImageID
    nextImageID++
    images[id] = &img
    return vInt(id), nil
}

func builtinImageDraw(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 10 { return nil, fmt.Errorf("imageDraw expects dstImageId, srcImageId, sx, sy, sw, sh, dx, dy, dw, dh, [tint]") }
    _, dst, err := imageByID("imageDraw", args, 0)
    if err != nil { return nil, err }
    _, src, err := imageByID("imageDraw", args, 1)
    if err != nil { return nil, err }
    sx, _ := getArgFloat("imageDraw", args, 2)
    sy, _ := getArgFloat("imageDraw", args, 3)
    sw, _ := getArgFloat("imageDraw", args, 4)
    sh, _ := getArgFloat("imageDraw", args, 5)
    dx, _ := getArgFloat("imageDraw", args, 6)
    dy, _ := getArgFloat("imageDraw", args, 7)
    dw, _ := getArgFloat("imageDraw", args, 8)
    dh, _ := getArgFloat("imageDraw", args, 9)
    tint, _ := argColor("imageDraw", args, 10, rl.White)
    rl.ImageDraw(dst, src, rl.NewRectangle(float32(sx), float32(sy), float32(sw), float32(sh)), rl.NewRectangle(float32(dx), float32(dy), float32(dw), float32(dh)), tint)
    return null(), nil
}

func parsePixelFormat(v *candy_evaluator.Value) (rl.PixelFormat, error) {
    if v == nil {
        return 0, fmt.Errorf("pixel format is required")
    }
    if v.Kind == candy_evaluator.ValInt {
        return rl.PixelFormat(v.I64), nil
    }
    if v.Kind != candy_evaluator.ValString {
        return 0, fmt.Errorf("pixel format must be int or string")
    }
    switch strings.ToLower(v.Str) {
    case "grayscale", "uncompressed_grayscale":
        return rl.UncompressedGrayscale, nil
    case "gray_alpha", "uncompressed_gray_alpha":
        return rl.UncompressedGrayAlpha, nil
    case "r5g6b5", "uncompressed_r5g6b5":
        return rl.UncompressedR5g6b5, nil
    case "r8g8b8", "uncompressed_r8g8b8":
        return rl.UncompressedR8g8b8, nil
    case "r5g5b5a1", "uncompressed_r5g5b5a1":
        return rl.UncompressedR5g5b5a1, nil
    case "r4g4b4a4", "uncompressed_r4g4b4a4":
        return rl.UncompressedR4g4b4a4, nil
    case "rgba", "r8g8b8a8", "uncompressed_r8g8b8a8":
        return rl.UncompressedR8g8b8a8, nil
    default:
        return 0, fmt.Errorf("unsupported pixel format %q", v.Str)
    }
}

func builtinImageFormat(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("imageFormat", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("imageFormat", args, 0)
    if err != nil { return nil, err }
    pf, err := parsePixelFormat(args[1])
    if err != nil { return nil, err }
    rl.ImageFormat(img, pf)
    return null(), nil
}

func builtinImageToPOT(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 1 { return nil, fmt.Errorf("imageToPOT expects imageId, [fillColor]") }
    _, img, err := imageByID("imageToPOT", args, 0)
    if err != nil { return nil, err }
    fill, _ := argColor("imageToPOT", args, 1, rl.Black)
    rl.ImageToPOT(img, fill)
    return null(), nil
}

func builtinImageText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 1 { return nil, fmt.Errorf("imageText expects text, [fontSize], [color]") }
    text, err := argString("imageText", args, 0)
    if err != nil { return nil, err }
    fontSize := int64(20)
    if len(args) > 1 {
        fontSize, _ = argInt("imageText", args, 1)
    }
    tint, _ := argColor("imageText", args, 2, rl.White)
    img := rl.ImageText(text, int32(fontSize), tint)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinImageTextEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 4 { return nil, fmt.Errorf("imageTextEx expects fontId, text, fontSize, spacing, [color]") }
    _, f, err := fontByID("imageTextEx", args, 0)
    if err != nil { return nil, err }
    text, err := argString("imageTextEx", args, 1)
    if err != nil { return nil, err }
    fontSize, _ := getArgFloat("imageTextEx", args, 2)
    spacing, _ := getArgFloat("imageTextEx", args, 3)
    tint, _ := argColor("imageTextEx", args, 4, rl.White)
    img := rl.ImageTextEx(f, text, float32(fontSize), float32(spacing), tint)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinExportImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("exportImage", args, 2); err != nil { return nil, err }
    _, img, err := imageByID("exportImage", args, 0)
    if err != nil { return nil, err }
    path, err := argString("exportImage", args, 1)
    if err != nil { return nil, err }
    return vBool(rl.ExportImage(*img, path)), nil
}

func builtinGenImageColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 2 { return nil, fmt.Errorf("genImageColor expects width, height, [color]") }
    w, _ := argInt("genImageColor", args, 0)
    h, _ := argInt("genImageColor", args, 1)
    c, _ := argColor("genImageColor", args, 2, rl.Black)
    img := rl.GenImageColor(int(w), int(h), c)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageGradientLinear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("genImageGradientLinear expects width, height, direction, startColor, endColor") }
    w, _ := argInt("genImageGradientLinear", args, 0)
    h, _ := argInt("genImageGradientLinear", args, 1)
    dir, _ := argInt("genImageGradientLinear", args, 2)
    start, _ := argColor("genImageGradientLinear", args, 3, rl.White)
    end, _ := argColor("genImageGradientLinear", args, 4, rl.Black)
    img := rl.GenImageGradientLinear(int(w), int(h), int(dir), start, end)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageGradientRadial(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("genImageGradientRadial expects width, height, density, innerColor, outerColor") }
    w, _ := argInt("genImageGradientRadial", args, 0)
    h, _ := argInt("genImageGradientRadial", args, 1)
    density, _ := getArgFloat("genImageGradientRadial", args, 2)
    inner, _ := argColor("genImageGradientRadial", args, 3, rl.White)
    outer, _ := argColor("genImageGradientRadial", args, 4, rl.Black)
    img := rl.GenImageGradientRadial(int(w), int(h), float32(density), inner, outer)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageGradientSquare(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("genImageGradientSquare expects width, height, density, innerColor, outerColor") }
    w, _ := argInt("genImageGradientSquare", args, 0)
    h, _ := argInt("genImageGradientSquare", args, 1)
    density, _ := getArgFloat("genImageGradientSquare", args, 2)
    inner, _ := argColor("genImageGradientSquare", args, 3, rl.White)
    outer, _ := argColor("genImageGradientSquare", args, 4, rl.Black)
    img := rl.GenImageGradientSquare(int(w), int(h), float32(density), inner, outer)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageChecked(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 6 { return nil, fmt.Errorf("genImageChecked expects width, height, checksX, checksY, color1, color2") }
    w, _ := argInt("genImageChecked", args, 0)
    h, _ := argInt("genImageChecked", args, 1)
    cx, _ := argInt("genImageChecked", args, 2)
    cy, _ := argInt("genImageChecked", args, 3)
    c1, _ := argColor("genImageChecked", args, 4, rl.White)
    c2, _ := argColor("genImageChecked", args, 5, rl.Black)
    img := rl.GenImageChecked(int(w), int(h), int(cx), int(cy), c1, c2)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageWhiteNoise(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 3 { return nil, fmt.Errorf("genImageWhiteNoise expects width, height, factor") }
    w, _ := argInt("genImageWhiteNoise", args, 0)
    h, _ := argInt("genImageWhiteNoise", args, 1)
    factor, _ := getArgFloat("genImageWhiteNoise", args, 2)
    img := rl.GenImageWhiteNoise(int(w), int(h), float32(factor))
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImagePerlinNoise(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 5 { return nil, fmt.Errorf("genImagePerlinNoise expects width, height, offsetX, offsetY, scale") }
    w, _ := argInt("genImagePerlinNoise", args, 0)
    h, _ := argInt("genImagePerlinNoise", args, 1)
    ox, _ := argInt("genImagePerlinNoise", args, 2)
    oy, _ := argInt("genImagePerlinNoise", args, 3)
    scale, _ := getArgFloat("genImagePerlinNoise", args, 4)
    img := rl.GenImagePerlinNoise(int(w), int(h), int(ox), int(oy), float32(scale))
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageCellular(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 3 { return nil, fmt.Errorf("genImageCellular expects width, height, tileSize") }
    w, _ := argInt("genImageCellular", args, 0)
    h, _ := argInt("genImageCellular", args, 1)
    tileSize, _ := argInt("genImageCellular", args, 2)
    img := rl.GenImageCellular(int(w), int(h), int(tileSize))
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}

func builtinGenImageText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("genImageText", args, 3); err != nil { return nil, err }
    w, _ := argInt("genImageText", args, 0)
    h, _ := argInt("genImageText", args, 1)
    text, err := argString("genImageText", args, 2)
    if err != nil { return nil, err }
    img := rl.GenImageText(int(w), int(h), text)
    id := nextImageID
    nextImageID++
    images[id] = img
    return vInt(id), nil
}
