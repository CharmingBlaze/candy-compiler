package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Font loading extras ----

func builtinGetFontDefault(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	f := rl.GetFontDefault()
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

func builtinLoadFontEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("loadFontEx expects path, fontSize, [codepoints...]")
	}
	path, err := argString("loadFontEx", args, 0)
	if err != nil {
		return nil, err
	}
	sz, _ := argInt("loadFontEx", args, 1)
	var runes []rune
	if len(args) > 2 && args[2] != nil && args[2].Kind == candy_evaluator.ValArray {
		for _, e := range args[2].Elems {
			if e.Kind == candy_evaluator.ValInt {
				runes = append(runes, rune(e.I64))
			}
		}
	}
	f := rl.LoadFontEx(path, int32(sz), runes)
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

func builtinLoadFontFromImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("loadFontFromImage expects imageId, keyColor, firstChar")
	}
	_, img, err := imageByID("loadFontFromImage", args, 0)
	if err != nil {
		return nil, err
	}
	key, _ := argColorValue("loadFontFromImage", args, 1, rl.Magenta)
	firstChar, _ := argInt("loadFontFromImage", args, 2)
	f := rl.LoadFontFromImage(*img, key, int32(firstChar))
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

func builtinLoadFontFromMemory(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("loadFontFromMemory expects fileType, bytesArr, fontSize, [codepoints]")
	}
	ft, err := argString("loadFontFromMemory", args, 0)
	if err != nil {
		return nil, err
	}
	if args[1] == nil || args[1].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("loadFontFromMemory: arg 2 must be byte array")
	}
	data := make([]byte, len(args[1].Elems))
	for i, e := range args[1].Elems {
		data[i] = byte(e.I64)
	}
	sz, _ := argInt("loadFontFromMemory", args, 2)
	var runes []rune
	if len(args) > 3 && args[3] != nil && args[3].Kind == candy_evaluator.ValArray {
		for _, e := range args[3].Elems {
			if e.Kind == candy_evaluator.ValInt {
				runes = append(runes, rune(e.I64))
			}
		}
	}
	f := rl.LoadFontFromMemory(ft, data, int32(sz), runes)
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

func builtinIsFontValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isFontValid", args, 1); err != nil {
		return nil, err
	}
	_, f, err := fontByID("isFontValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsFontValid(f)), nil
}

// ---- Text drawing extras ----

func builtinDrawTextPro(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("drawTextPro expects fontId, text, position, origin, rotation, fontSize, spacing, [tint]")
	}
	_, f, err := fontByID("drawTextPro", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("drawTextPro", args, 1)
	pos, _ := argVector2("drawTextPro", args, 2)
	origin, _ := argVector2("drawTextPro", args, 3)
	rot, _ := getArgFloat("drawTextPro", args, 4)
	fs, _ := getArgFloat("drawTextPro", args, 5)
	sp, _ := getArgFloat("drawTextPro", args, 6)
	c, _ := argColorValue("drawTextPro", args, 7, rl.RayWhite)
	rl.DrawTextPro(f, text, pos, origin, float32(rot), float32(fs), float32(sp), c)
	return null(), nil
}

func builtinDrawTextCodepoint(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawTextCodepoint expects fontId, codepoint, position, fontSize, [tint]")
	}
	_, f, err := fontByID("drawTextCodepoint", args, 0)
	if err != nil {
		return nil, err
	}
	cp, _ := argInt("drawTextCodepoint", args, 1)
	pos, _ := argVector2("drawTextCodepoint", args, 2)
	fs, _ := getArgFloat("drawTextCodepoint", args, 3)
	c, _ := argColorValue("drawTextCodepoint", args, 4, rl.RayWhite)
	rl.DrawTextCodepoint(f, rune(cp), pos, float32(fs), c)
	return null(), nil
}

func builtinDrawTextCodepoints(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawTextCodepoints expects fontId, codepointsArr, position, fontSize, spacing, [tint]")
	}
	_, f, err := fontByID("drawTextCodepoints", args, 0)
	if err != nil {
		return nil, err
	}
	if args[1] == nil || args[1].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("drawTextCodepoints: arg 2 must be int array")
	}
	var runes []rune
	for _, e := range args[1].Elems {
		if e.Kind == candy_evaluator.ValInt {
			runes = append(runes, rune(e.I64))
		}
	}
	pos, _ := argVector2("drawTextCodepoints", args, 2)
	fs, _ := getArgFloat("drawTextCodepoints", args, 3)
	sp, _ := getArgFloat("drawTextCodepoints", args, 4)
	c, _ := argColorValue("drawTextCodepoints", args, 5, rl.RayWhite)
	rl.DrawTextCodepoints(f, runes, pos, float32(fs), float32(sp), c)
	return null(), nil
}

// ---- Font info / metrics ----

func builtinSetTextLineSpacing(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setTextLineSpacing", args, 1); err != nil {
		return nil, err
	}
	sp, _ := argInt("setTextLineSpacing", args, 0)
	rl.SetTextLineSpacing(int(sp))
	return null(), nil
}

func builtinMeasureTextEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("measureTextEx expects fontId, text, fontSize, spacing")
	}
	_, f, err := fontByID("measureTextEx", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("measureTextEx", args, 1)
	fs, _ := getArgFloat("measureTextEx", args, 2)
	sp, _ := getArgFloat("measureTextEx", args, 3)
	v2 := rl.MeasureTextEx(f, text, float32(fs), float32(sp))
	return vec2ToMap(v2), nil
}

func builtinGetGlyphIndex(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGlyphIndex", args, 2); err != nil {
		return nil, err
	}
	_, f, err := fontByID("getGlyphIndex", args, 0)
	if err != nil {
		return nil, err
	}
	cp, _ := argInt("getGlyphIndex", args, 1)
	return vInt(int64(rl.GetGlyphIndex(f, int32(cp)))), nil
}

func builtinGetGlyphInfo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGlyphInfo", args, 2); err != nil {
		return nil, err
	}
	_, f, err := fontByID("getGlyphInfo", args, 0)
	if err != nil {
		return nil, err
	}
	cp, _ := argInt("getGlyphInfo", args, 1)
	gi := rl.GetGlyphInfo(f, int32(cp))
	return vMap(map[string]candy_evaluator.Value{
		"value":    {Kind: candy_evaluator.ValInt, I64: int64(gi.Value)},
		"offsetX":  {Kind: candy_evaluator.ValInt, I64: int64(gi.OffsetX)},
		"offsetY":  {Kind: candy_evaluator.ValInt, I64: int64(gi.OffsetY)},
		"advanceX": {Kind: candy_evaluator.ValInt, I64: int64(gi.AdvanceX)},
	}), nil
}

func builtinGetGlyphAtlasRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGlyphAtlasRec", args, 2); err != nil {
		return nil, err
	}
	_, f, err := fontByID("getGlyphAtlasRec", args, 0)
	if err != nil {
		return nil, err
	}
	cp, _ := argInt("getGlyphAtlasRec", args, 1)
	rec := rl.GetGlyphAtlasRec(f, int32(cp))
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.X)},
		"y":      {Kind: candy_evaluator.ValFloat, F64: float64(rec.Y)},
		"width":  {Kind: candy_evaluator.ValFloat, F64: float64(rec.Width)},
		"height": {Kind: candy_evaluator.ValFloat, F64: float64(rec.Height)},
	}), nil
}

// ---- Text string utility functions (implemented natively in Go) ----
// The C-side TextXxx buffer functions are not exposed by the cgo binding;
// we implement equivalent behaviour using Go's standard library.

func builtinTextIsEqual(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textIsEqual", args, 2); err != nil {
		return nil, err
	}
	a, _ := argString("textIsEqual", args, 0)
	b, _ := argString("textIsEqual", args, 1)
	return vBool(a == b), nil
}

func builtinTextLength(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textLength", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textLength", args, 0)
	return vInt(int64(len(s))), nil
}

func builtinTextSubtext(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textSubtext", args, 3); err != nil {
		return nil, err
	}
	s, _ := argString("textSubtext", args, 0)
	pos, _ := argInt("textSubtext", args, 1)
	length, _ := argInt("textSubtext", args, 2)
	start := int(pos)
	if start < 0 {
		start = 0
	}
	if start >= len(s) {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: ""}, nil
	}
	end := start + int(length)
	if end > len(s) {
		end = len(s)
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: s[start:end]}, nil
}

func builtinTextRemoveSpaces(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textRemoveSpaces", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textRemoveSpaces", args, 0)
	parts := strings.Fields(s)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.Join(parts, "")}, nil
}

func builtinGetTextBetween(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getTextBetween", args, 3); err != nil {
		return nil, err
	}
	s, _ := argString("getTextBetween", args, 0)
	begin, _ := argString("getTextBetween", args, 1)
	end, _ := argString("getTextBetween", args, 2)
	si := strings.Index(s, begin)
	if si == -1 {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: ""}, nil
	}
	si += len(begin)
	ei := strings.Index(s[si:], end)
	if ei == -1 {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: s[si:]}, nil
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: s[si : si+ei]}, nil
}

func builtinTextReplace(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textReplace", args, 3); err != nil {
		return nil, err
	}
	s, _ := argString("textReplace", args, 0)
	search, _ := argString("textReplace", args, 1)
	repl, _ := argString("textReplace", args, 2)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.ReplaceAll(s, search, repl)}, nil
}

func builtinTextInsert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textInsert", args, 3); err != nil {
		return nil, err
	}
	s, _ := argString("textInsert", args, 0)
	ins, _ := argString("textInsert", args, 1)
	pos, _ := argInt("textInsert", args, 2)
	p := int(pos)
	if p < 0 {
		p = 0
	}
	if p > len(s) {
		p = len(s)
	}
	result := s[:p] + ins + s[p:]
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: result}, nil
}

func builtinTextJoin(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("textJoin expects stringsArr, delimiter")
	}
	if args[0] == nil || args[0].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("textJoin: arg 1 must be string array")
	}
	delim, _ := argString("textJoin", args, 1)
	parts := make([]string, len(args[0].Elems))
	for i, e := range args[0].Elems {
		if e.Kind == candy_evaluator.ValString {
			parts[i] = e.Str
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.Join(parts, delim)}, nil
}

func builtinTextSplit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textSplit", args, 2); err != nil {
		return nil, err
	}
	s, _ := argString("textSplit", args, 0)
	delim, _ := argString("textSplit", args, 1)
	var parts []string
	if delim == "" {
		parts = strings.Split(s, "")
	} else {
		parts = strings.Split(s, delim)
	}
	elems := make([]candy_evaluator.Value, len(parts))
	for i, p := range parts {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: p}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

func builtinTextFindIndex(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textFindIndex", args, 2); err != nil {
		return nil, err
	}
	s, _ := argString("textFindIndex", args, 0)
	search, _ := argString("textFindIndex", args, 1)
	return vInt(int64(strings.Index(s, search))), nil
}

func builtinTextToUpper(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToUpper", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToUpper", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.ToUpper(s)}, nil
}

func builtinTextToLower(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToLower", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToLower", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.ToLower(s)}, nil
}

func builtinTextToPascal(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToPascal", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToPascal", args, 0)
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.Join(words, "")}, nil
}

func builtinTextToSnake(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToSnake", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToSnake", args, 0)
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	result := strings.ReplaceAll(b.String(), " ", "_")
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: result}, nil
}

func builtinTextToCamel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToCamel", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToCamel", args, 0)
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) == 0 {
			continue
		}
		if i == 0 {
			words[i] = strings.ToLower(w)
		} else {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: strings.Join(words, "")}, nil
}

func builtinTextToInteger(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToInteger", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToInteger", args, 0)
	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return vInt(n), nil
}

func builtinTextToFloat(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("textToFloat", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("textToFloat", args, 0)
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return vFloat(f), nil
}

// ---- Codepoint helpers (Go-native) ----

func builtinGetCodepointCount(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getCodepointCount", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("getCodepointCount", args, 0)
	return vInt(int64(len([]rune(s)))), nil
}

func builtinLoadCodepoints(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadCodepoints", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("loadCodepoints", args, 0)
	runes := []rune(s)
	elems := make([]candy_evaluator.Value, len(runes))
	for i, r := range runes {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValInt, I64: int64(r)}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

func builtinCodepointToUTF8(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("codepointToUTF8", args, 1); err != nil {
		return nil, err
	}
	cp, _ := argInt("codepointToUTF8", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: string(rune(cp))}, nil
}

func builtinLoadUTF8(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadUTF8", args, 1); err != nil {
		return nil, err
	}
	if args[0] == nil || args[0].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("loadUTF8: arg 1 must be int (codepoint) array")
	}
	runes := make([]rune, len(args[0].Elems))
	for i, e := range args[0].Elems {
		runes[i] = rune(e.I64)
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: string(runes)}, nil
}

func builtinLoadTextLines(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadTextLines", args, 1); err != nil {
		return nil, err
	}
	s, _ := argString("loadTextLines", args, 0)
	lines := strings.Split(s, "\n")
	elems := make([]candy_evaluator.Value, len(lines))
	for i, l := range lines {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: l}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}
