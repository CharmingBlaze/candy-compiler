package candy_raylib

import (
	"fmt"

	"candy/candy_evaluator"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- helpers ----------------------------------------------------------------

// argRect reads a {x,y,width,height} map or 4 flat floats starting at args[i].
// Returns the Rectangle and the number of args consumed (1 for map, 4 for flat).
func argRect(fname string, args []*candy_evaluator.Value, i int) (rl.Rectangle, int, error) {
	if i >= len(args) {
		return rl.Rectangle{}, 0, fmt.Errorf("%s: missing bounds argument", fname)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValMap {
		x := float32(mapFloatDef(v.StrMap, "x", 0))
		y := float32(mapFloatDef(v.StrMap, "y", 0))
		w := float32(mapFloatDef(v.StrMap, "width", mapFloatDef(v.StrMap, "w", 0)))
		h := float32(mapFloatDef(v.StrMap, "height", mapFloatDef(v.StrMap, "h", 0)))
		return rl.Rectangle{X: x, Y: y, Width: w, Height: h}, 1, nil
	}
	// flat: x y width height
	if i+3 >= len(args) {
		return rl.Rectangle{}, 0, fmt.Errorf("%s: bounds needs {x,y,width,height} map or 4 numbers", fname)
	}
	x, _ := getArgFloat(fname, args, i)
	y, _ := getArgFloat(fname, args, i+1)
	w, _ := getArgFloat(fname, args, i+2)
	h, _ := getArgFloat(fname, args, i+3)
	return rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}, 4, nil
}

func mapFloatDef(m map[string]candy_evaluator.Value, key string, def float64) float64 {
	if v, ok := m[key]; ok {
		switch v.Kind {
		case candy_evaluator.ValFloat:
			return v.F64
		case candy_evaluator.ValInt:
			return float64(v.I64)
		}
	}
	return def
}

func rectToMap(r rl.Rectangle) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x":      {Kind: candy_evaluator.ValFloat, F64: float64(r.X)},
		"y":      {Kind: candy_evaluator.ValFloat, F64: float64(r.Y)},
		"width":  {Kind: candy_evaluator.ValFloat, F64: float64(r.Width)},
		"height": {Kind: candy_evaluator.ValFloat, F64: float64(r.Height)},
	})
}

// argBool reads a bool from args[i] (bool, int, or float).
func argBool(_ string, args []*candy_evaluator.Value, i int) bool {
	if i >= len(args) || args[i] == nil {
		return false
	}
	v := args[i]
	switch v.Kind {
	case candy_evaluator.ValBool:
		return v.B
	case candy_evaluator.ValInt:
		return v.I64 != 0
	case candy_evaluator.ValFloat:
		return v.F64 != 0
	}
	return false
}

// argFloat is an alias kept for clarity inside this file.
func argFloat(fname string, args []*candy_evaluator.Value, i int) (float64, error) {
	return getArgFloat(fname, args, i)
}

// ---- Global state -----------------------------------------------------------

// guiState holds mutable per-call state (scroll positions, active items etc.)
// keyed by an int64 handle so Candy scripts can manage multiple widgets.
var (
	guiScrollVec        = map[int64]*rl.Vector2{}
	guiScrollNext int64 = 1
	guiActiveNext int64 = 1
	guiStrBuf           = map[int64]string{} // TextBox / TextInputBox / ValueBoxFloat text
	guiStrNext    int64 = 1
)

// guiNewScroll allocates a scroll vector handle.
func guiNewScroll() int64 {
	id := guiScrollNext
	guiScrollNext++
	guiScrollVec[id] = &rl.Vector2{}
	return id
}

func guiScrollOf(id int64) *rl.Vector2 {
	if v, ok := guiScrollVec[id]; ok {
		return v
	}
	v := &rl.Vector2{}
	guiScrollVec[id] = v
	return v
}

// ---- Style / state ----------------------------------------------------------

// guiEnable()
func builtinGuiEnable(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.Enable()
	return null(), nil
}

// guiDisable()
func builtinGuiDisable(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.Disable()
	return null(), nil
}

// guiLock()
func builtinGuiLock(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.Lock()
	return null(), nil
}

// guiUnlock()
func builtinGuiUnlock(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.Unlock()
	return null(), nil
}

// guiIsLocked() → bool
func builtinGuiIsLocked(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rg.IsLocked()), nil
}

// guiSetAlpha(alpha)
func builtinGuiSetAlpha(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiSetAlpha expects alpha")
	}
	a, _ := argFloat("guiSetAlpha", args, 0)
	rg.SetAlpha(float32(a))
	return null(), nil
}

// guiSetState(state)
func builtinGuiSetState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiSetState expects state int")
	}
	s, _ := argInt("guiSetState", args, 0)
	rg.SetState(rg.PropertyValue(s))
	return null(), nil
}

// guiGetState() → int
func builtinGuiGetState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rg.GetState())), nil
}

// guiSetFont(fontId)
func builtinGuiSetFont(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiSetFont expects fontId")
	}
	fid, _ := argInt("guiSetFont", args, 0)
	f, ok := fonts[fid]
	if !ok {
		return nil, fmt.Errorf("guiSetFont: unknown fontId %d", fid)
	}
	rg.SetFont(f)
	return null(), nil
}

// guiGetFont() → fontId  (stores in fonts map and returns the handle)
func builtinGuiGetFont(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	f := rg.GetFont()
	id := nextFontID
	nextFontID++
	fonts[id] = f
	return vInt(id), nil
}

// guiSetStyle(control, property, value)
func builtinGuiSetStyle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("guiSetStyle expects control, property, value")
	}
	ctrl, _ := argInt("guiSetStyle", args, 0)
	prop, _ := argInt("guiSetStyle", args, 1)
	val, _ := argInt("guiSetStyle", args, 2)
	rg.SetStyle(rg.ControlID(ctrl), rg.PropertyID(prop), rg.PropertyValue(val))
	return null(), nil
}

// guiGetStyle(control, property) → int
func builtinGuiGetStyle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiGetStyle expects control, property")
	}
	ctrl, _ := argInt("guiGetStyle", args, 0)
	prop, _ := argInt("guiGetStyle", args, 1)
	return vInt(int64(rg.GetStyle(rg.ControlID(ctrl), rg.PropertyID(prop)))), nil
}

// guiLoadStyle(fileName)
func builtinGuiLoadStyle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiLoadStyle expects fileName")
	}
	p, _ := argString("guiLoadStyle", args, 0)
	rg.LoadStyle(p)
	return null(), nil
}

// guiLoadStyleDefault()
func builtinGuiLoadStyleDefault(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.LoadStyleDefault()
	return null(), nil
}

// guiSetTooltip(text)
func builtinGuiSetTooltip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiSetTooltip expects text")
	}
	t, _ := argString("guiSetTooltip", args, 0)
	rg.SetTooltip(t)
	return null(), nil
}

// guiEnableTooltip()
func builtinGuiEnableTooltip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.EnableTooltip()
	return null(), nil
}

// guiDisableTooltip()
func builtinGuiDisableTooltip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rg.DisableTooltip()
	return null(), nil
}

// guiSetIconScale(scale)
func builtinGuiSetIconScale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiSetIconScale expects scale int")
	}
	s, _ := argInt("guiSetIconScale", args, 0)
	rg.SetIconScale(int32(s))
	return null(), nil
}

// ---- Basic controls ---------------------------------------------------------

// guiLabel(x, y, width, height, text)  or  guiLabel({x,y,width,height}, text)
func builtinGuiLabel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiLabel", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiLabel", args, n)
	rg.Label(bounds, text)
	return null(), nil
}

// guiButton(x, y, w, h, text) → bool
func builtinGuiButton(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiButton", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiButton", args, n)
	return vBool(rg.Button(bounds, text)), nil
}

// guiLabelButton(x, y, w, h, text) → bool
func builtinGuiLabelButton(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiLabelButton", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiLabelButton", args, n)
	return vBool(rg.LabelButton(bounds, text)), nil
}

// guiToggle(x, y, w, h, text, active) → bool
func builtinGuiToggle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiToggle", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiToggle", args, n)
	active := argBool("guiToggle", args, n+1)
	return vBool(rg.Toggle(bounds, text, active)), nil
}

// guiToggleGroup(x, y, w, h, text, active) → int  (text is semicolon-separated items)
func builtinGuiToggleGroup(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiToggleGroup", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiToggleGroup", args, n)
	active, _ := argInt("guiToggleGroup", args, n+1)
	return vInt(int64(rg.ToggleGroup(bounds, text, int32(active)))), nil
}

// guiToggleSlider(x, y, w, h, text, active) → int
func builtinGuiToggleSlider(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiToggleSlider", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiToggleSlider", args, n)
	active, _ := argInt("guiToggleSlider", args, n+1)
	return vInt(int64(rg.ToggleSlider(bounds, text, int32(active)))), nil
}

// guiCheckBox(x, y, w, h, text, checked) → bool
func builtinGuiCheckBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiCheckBox", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiCheckBox", args, n)
	checked := argBool("guiCheckBox", args, n+1)
	return vBool(rg.CheckBox(bounds, text, checked)), nil
}

// guiComboBox(x, y, w, h, text, active) → int  (text is semicolon-separated items)
func builtinGuiComboBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiComboBox", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiComboBox", args, n)
	active, _ := argInt("guiComboBox", args, n+1)
	return vInt(int64(rg.ComboBox(bounds, text, int32(active)))), nil
}

// guiDropdownBox(x, y, w, h, text, activeHandle, editMode) → bool (true when value changed)
// activeHandle is an int64 returned by guiNewActive(); editMode bool
func builtinGuiDropdownBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiDropdownBox", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiDropdownBox", args, n)
	handle, _ := argInt("guiDropdownBox", args, n+1)
	editMode := argBool("guiDropdownBox", args, n+2)
	ptr := guiActiveOf(handle)
	changed := rg.DropdownBox(bounds, text, ptr, editMode)
	return vBool(changed), nil
}

// guiNewActive() → handle  — allocate a shared int32 active-index slot
func builtinGuiNewActive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := guiActiveNext
	guiActiveNext++
	_ = guiActiveOf(id) // pre-allocate slot
	return vInt(id), nil
}

// guiGetActive(handle) → int
func builtinGuiGetActive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiGetActive expects handle")
	}
	h, _ := argInt("guiGetActive", args, 0)
	return vInt(int64(*guiActiveOf(h))), nil
}

// guiSetActive(handle, value)
func builtinGuiSetActive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiSetActive expects handle, value")
	}
	h, _ := argInt("guiSetActive", args, 0)
	v, _ := argInt("guiSetActive", args, 1)
	*guiActiveOf(h) = int32(v)
	return null(), nil
}

// guiActiveSlots holds stable int32 values addressable by pointer.
var guiActiveSlots []*int32

func guiActiveOf(id int64) *int32 {
	idx := int(id) - 1
	for len(guiActiveSlots) <= idx {
		v := int32(0)
		guiActiveSlots = append(guiActiveSlots, &v)
	}
	return guiActiveSlots[idx]
}

// ---- Slider family ----------------------------------------------------------

// guiSlider(x,y,w,h, textLeft, textRight, value, min, max) → float
func builtinGuiSlider(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiSlider", args, 0)
	if err != nil {
		return nil, err
	}
	tl, _ := argString("guiSlider", args, n)
	tr, _ := argString("guiSlider", args, n+1)
	val, _ := argFloat("guiSlider", args, n+2)
	mn, _ := argFloat("guiSlider", args, n+3)
	mx, _ := argFloat("guiSlider", args, n+4)
	out := rg.Slider(bounds, tl, tr, float32(val), float32(mn), float32(mx))
	return vFloat(float64(out)), nil
}

// guiSliderBar(x,y,w,h, textLeft, textRight, value, min, max) → float
func builtinGuiSliderBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiSliderBar", args, 0)
	if err != nil {
		return nil, err
	}
	tl, _ := argString("guiSliderBar", args, n)
	tr, _ := argString("guiSliderBar", args, n+1)
	val, _ := argFloat("guiSliderBar", args, n+2)
	mn, _ := argFloat("guiSliderBar", args, n+3)
	mx, _ := argFloat("guiSliderBar", args, n+4)
	out := rg.SliderBar(bounds, tl, tr, float32(val), float32(mn), float32(mx))
	return vFloat(float64(out)), nil
}

// guiProgressBar(x,y,w,h, textLeft, textRight, value, min, max) → float
func builtinGuiProgressBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiProgressBar", args, 0)
	if err != nil {
		return nil, err
	}
	tl, _ := argString("guiProgressBar", args, n)
	tr, _ := argString("guiProgressBar", args, n+1)
	val, _ := argFloat("guiProgressBar", args, n+2)
	mn, _ := argFloat("guiProgressBar", args, n+3)
	mx, _ := argFloat("guiProgressBar", args, n+4)
	out := rg.ProgressBar(bounds, tl, tr, float32(val), float32(mn), float32(mx))
	return vFloat(float64(out)), nil
}

// guiScrollBar(x,y,w,h, value, min, max) → int
func builtinGuiScrollBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiScrollBar", args, 0)
	if err != nil {
		return nil, err
	}
	val, _ := argInt("guiScrollBar", args, n)
	mn, _ := argInt("guiScrollBar", args, n+1)
	mx, _ := argInt("guiScrollBar", args, n+2)
	out := rg.ScrollBar(bounds, int32(val), int32(mn), int32(mx))
	return vInt(int64(out)), nil
}

// guiSpinner(x,y,w,h, text, activeHandle, min, max, editMode) → bool
func builtinGuiSpinner(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiSpinner", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiSpinner", args, n)
	handle, _ := argInt("guiSpinner", args, n+1)
	mn, _ := argInt("guiSpinner", args, n+2)
	mx, _ := argInt("guiSpinner", args, n+3)
	editMode := argBool("guiSpinner", args, n+4)
	ptr := guiActiveOf(handle)
	changed := rg.Spinner(bounds, text, ptr, int(mn), int(mx), editMode)
	return vBool(changed), nil
}

// guiValueBox(x,y,w,h, text, activeHandle, min, max, editMode) → bool
func builtinGuiValueBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiValueBox", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiValueBox", args, n)
	handle, _ := argInt("guiValueBox", args, n+1)
	mn, _ := argInt("guiValueBox", args, n+2)
	mx, _ := argInt("guiValueBox", args, n+3)
	editMode := argBool("guiValueBox", args, n+4)
	ptr := guiActiveOf(handle)
	changed := rg.ValueBox(bounds, text, ptr, int(mn), int(mx), editMode)
	return vBool(changed), nil
}

// ---- Text input -------------------------------------------------------------

// guiNewTextBox(initialText) → handle
func builtinGuiNewTextBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := guiStrNext
	guiStrNext++
	init := ""
	if len(args) >= 1 && args[0] != nil && args[0].Kind == candy_evaluator.ValString {
		init = args[0].Str
	}
	guiStrBuf[id] = init
	return vInt(id), nil
}

// guiGetText(handle) → string
func builtinGuiGetText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiGetText expects handle")
	}
	h, _ := argInt("guiGetText", args, 0)
	s, ok := guiStrBuf[h]
	if !ok {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: ""}, nil
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: s}, nil
}

// guiSetText(handle, text)
func builtinGuiSetText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiSetText expects handle, text")
	}
	h, _ := argInt("guiSetText", args, 0)
	t, _ := argString("guiSetText", args, 1)
	guiStrBuf[h] = t
	return null(), nil
}

// guiTextBox(x,y,w,h, textHandle, maxLen, editMode) → bool (true while editing)
func builtinGuiTextBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiTextBox", args, 0)
	if err != nil {
		return nil, err
	}
	handle, _ := argInt("guiTextBox", args, n)
	maxLen, _ := argInt("guiTextBox", args, n+1)
	editMode := argBool("guiTextBox", args, n+2)
	s, ok := guiStrBuf[handle]
	if !ok {
		s = ""
	}
	editing := rg.TextBox(bounds, &s, int(maxLen), editMode)
	guiStrBuf[handle] = s
	return vBool(editing), nil
}

// ---- Container / layout -----------------------------------------------------

// guiWindowBox(x,y,w,h, title) → bool (true when close button clicked)
func builtinGuiWindowBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiWindowBox", args, 0)
	if err != nil {
		return nil, err
	}
	title, _ := argString("guiWindowBox", args, n)
	return vBool(rg.WindowBox(bounds, title)), nil
}

// guiGroupBox(x,y,w,h, text)
func builtinGuiGroupBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiGroupBox", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiGroupBox", args, n)
	rg.GroupBox(bounds, text)
	return null(), nil
}

// guiLine(x,y,w,h, text)
func builtinGuiLine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiLine", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiLine", args, n)
	rg.Line(bounds, text)
	return null(), nil
}

// guiPanel(x,y,w,h, text)
func builtinGuiPanel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiPanel", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiPanel", args, n)
	rg.Panel(bounds, text)
	return null(), nil
}

// guiStatusBar(x,y,w,h, text)
func builtinGuiStatusBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiStatusBar", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiStatusBar", args, n)
	rg.StatusBar(bounds, text)
	return null(), nil
}

// guiDummyRec(x,y,w,h, text)
func builtinGuiDummyRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiDummyRec", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiDummyRec", args, n)
	rg.DummyRec(bounds, text)
	return null(), nil
}

// guiScrollPanel(x,y,w,h, text, contentW, contentH, scrollHandle) → {x,y,w,h} visible view
func builtinGuiScrollPanel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiScrollPanel", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiScrollPanel", args, n)
	cw, _ := argFloat("guiScrollPanel", args, n+1)
	ch, _ := argFloat("guiScrollPanel", args, n+2)
	handle, _ := argInt("guiScrollPanel", args, n+3)
	content := rl.Rectangle{Width: float32(cw), Height: float32(ch)}
	scrollPtr := guiScrollOf(handle)
	var view rl.Rectangle
	rg.ScrollPanel(bounds, text, content, scrollPtr, &view)
	return rectToMap(view), nil
}

// guiNewScroll() → handle   — allocate a scroll state slot
func builtinGuiNewScroll(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(guiNewScroll()), nil
}

// guiGetScroll(handle) → {x, y}
func builtinGuiGetScroll(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiGetScroll expects handle")
	}
	h, _ := argInt("guiGetScroll", args, 0)
	v := guiScrollOf(h)
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(v.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(v.Y)},
	}), nil
}

// ---- List / grid ------------------------------------------------------------

// guiListView(x,y,w,h, text, scrollHandle, active) → int
// text is semicolon-separated items
func builtinGuiListView(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiListView", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiListView", args, n)
	scrollHandle, _ := argInt("guiListView", args, n+1)
	active, _ := argInt("guiListView", args, n+2)
	scrollPtr := guiActiveOf(scrollHandle)
	result := rg.ListView(bounds, text, scrollPtr, int32(active))
	return vInt(int64(result)), nil
}

// guiGrid(x,y,w,h, text, spacing, subdivs) → {x,y} cell mouse is over, or {-1,-1}
func builtinGuiGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiGrid", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiGrid", args, n)
	spacing, _ := argFloat("guiGrid", args, n+1)
	subdivs, _ := argInt("guiGrid", args, n+2)
	var mouseCell rl.Vector2
	rg.Grid(bounds, text, float32(spacing), int32(subdivs), &mouseCell)
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(mouseCell.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(mouseCell.Y)},
	}), nil
}

// guiTabBar(x,y,w,h, tabs, active) → int   tabs is a Candy array of strings
func builtinGuiTabBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiTabBar", args, 0)
	if err != nil {
		return nil, err
	}
	if n >= len(args) {
		return nil, fmt.Errorf("guiTabBar: missing tabs array")
	}
	tabsVal := args[n]
	var texts []string
	switch tabsVal.Kind {
	case candy_evaluator.ValArray:
		for _, e := range tabsVal.Elems {
			if e.Kind == candy_evaluator.ValString {
				texts = append(texts, e.Str)
			}
		}
	case candy_evaluator.ValString:
		texts = append(texts, tabsVal.Str)
	}
	active, _ := argInt("guiTabBar", args, n+1)
	act := int32(active)
	result := rg.TabBar(bounds, texts, &act)
	return vInt(int64(result)), nil
}

// ---- Dialogs ----------------------------------------------------------------

// guiMessageBox(x,y,w,h, title, message, buttons) → int  (-1 no action, 0+ button index)
// buttons is semicolon-separated
func builtinGuiMessageBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiMessageBox", args, 0)
	if err != nil {
		return nil, err
	}
	title, _ := argString("guiMessageBox", args, n)
	message, _ := argString("guiMessageBox", args, n+1)
	buttons, _ := argString("guiMessageBox", args, n+2)
	result := rg.MessageBox(bounds, title, message, buttons)
	return vInt(int64(result)), nil
}

// guiTextInputBox(x,y,w,h, title, message, buttons, textHandle) → int
func builtinGuiTextInputBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiTextInputBox", args, 0)
	if err != nil {
		return nil, err
	}
	title, _ := argString("guiTextInputBox", args, n)
	message, _ := argString("guiTextInputBox", args, n+1)
	buttons, _ := argString("guiTextInputBox", args, n+2)
	handle, _ := argInt("guiTextInputBox", args, n+3)
	s, ok := guiStrBuf[handle]
	if !ok {
		s = ""
	}
	result := rg.TextInputBox(bounds, title, message, buttons, &s, 256, nil)
	guiStrBuf[handle] = s
	return vInt(int64(result)), nil
}

// ---- Color picker -----------------------------------------------------------

// guiColorPicker(x,y,w,h, text, color) → {r,g,b,a}
func builtinGuiColorPicker(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiColorPicker", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiColorPicker", args, n)
	col, _ := argColor("guiColorPicker", args, n+1, rl.White)
	result := rg.ColorPicker(bounds, text, col)
	return colorToMap(result), nil
}

// guiColorPanel(x,y,w,h, text, color) → {r,g,b,a}
func builtinGuiColorPanel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiColorPanel", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiColorPanel", args, n)
	col, _ := argColor("guiColorPanel", args, n+1, rl.White)
	result := rg.ColorPanel(bounds, text, col)
	return colorToMap(result), nil
}

// guiColorBarAlpha(x,y,w,h, text, alpha) → float
func builtinGuiColorBarAlpha(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiColorBarAlpha", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiColorBarAlpha", args, n)
	alpha, _ := argFloat("guiColorBarAlpha", args, n+1)
	result := rg.ColorBarAlpha(bounds, text, float32(alpha))
	return vFloat(float64(result)), nil
}

// guiColorBarHue(x,y,w,h, text, value) → float
func builtinGuiColorBarHue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiColorBarHue", args, 0)
	if err != nil {
		return nil, err
	}
	text, _ := argString("guiColorBarHue", args, n)
	val, _ := argFloat("guiColorBarHue", args, n+1)
	result := rg.ColorBarHue(bounds, text, float32(val))
	return vFloat(float64(result)), nil
}

// ---- Drawing helpers --------------------------------------------------------

// guiDrawIcon(iconId, x, y, pixelSize, color)
func builtinGuiDrawIcon(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("guiDrawIcon expects iconId, x, y, pixelSize, color")
	}
	icon, _ := argInt("guiDrawIcon", args, 0)
	x, _ := argInt("guiDrawIcon", args, 1)
	y, _ := argInt("guiDrawIcon", args, 2)
	sz, _ := argInt("guiDrawIcon", args, 3)
	col, _ := argColor("guiDrawIcon", args, 4, rl.White)
	rg.DrawIcon(rg.IconID(icon), int32(x), int32(y), int32(sz), col)
	return null(), nil
}

// guiDrawRectangle(x,y,w,h, borderWidth, borderColor, fillColor)
func builtinGuiDrawRectangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	bounds, n, err := argRect("guiDrawRectangle", args, 0)
	if err != nil {
		return nil, err
	}
	bw, _ := argInt("guiDrawRectangle", args, n)
	bc, _ := argColor("guiDrawRectangle", args, n+1, rl.Black)
	fc, _ := argColor("guiDrawRectangle", args, n+2, rl.White)
	rg.DrawRectangle(bounds, int32(bw), bc, fc)
	return null(), nil
}

// guiDrawText(text, x,y,w,h, alignment, color)
// alignment: 0=left, 1=center, 2=right
func builtinGuiDrawText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiDrawText expects text, bounds..., alignment, color")
	}
	text, _ := argString("guiDrawText", args, 0)
	bounds, n, err := argRect("guiDrawText", args, 1)
	if err != nil {
		return nil, err
	}
	align, _ := argInt("guiDrawText", args, 1+n)
	col, _ := argColor("guiDrawText", args, 1+n+1, rl.Black)
	rg.DrawText(text, bounds, int32(align), col)
	return null(), nil
}

// guiGetTextBounds(control, x,y,w,h) → {x,y,width,height}
func builtinGuiGetTextBounds(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiGetTextBounds expects control, bounds...")
	}
	ctrl, _ := argInt("guiGetTextBounds", args, 0)
	bounds, _, err := argRect("guiGetTextBounds", args, 1)
	if err != nil {
		return nil, err
	}
	result := rg.GetTextBounds(rg.ControlID(ctrl), bounds)
	return rectToMap(result), nil
}

// guiGetTextWidth(text) → int
func builtinGuiGetTextWidth(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("guiGetTextWidth expects text")
	}
	t, _ := argString("guiGetTextWidth", args, 0)
	return vInt(int64(rg.GetTextWidth(t))), nil
}

// guiIconText(iconId, text) → string  (prepends icon code to text)
func builtinGuiIconText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiIconText expects iconId, text")
	}
	icon, _ := argInt("guiIconText", args, 0)
	text, _ := argString("guiIconText", args, 1)
	result := rg.IconText(rg.IconID(icon), text)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: result}, nil
}

// guiLoadIcons(fileName, loadNames)
func builtinGuiLoadIcons(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiLoadIcons expects fileName, loadNames bool")
	}
	p, _ := argString("guiLoadIcons", args, 0)
	load := argBool("guiLoadIcons", args, 1)
	rg.LoadIcons(p, load)
	return null(), nil
}

// guiFade(color, alpha) → {r,g,b,a}
func builtinGuiFade(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiFade expects color, alpha")
	}
	col, _ := argColor("guiFade", args, 0, rl.White)
	alpha, _ := argFloat("guiFade", args, 1)
	result := rg.Fade(col, float32(alpha))
	return colorToMap(result), nil
}

// guiGetColor(control, property) → {r,g,b,a}
func builtinGuiGetColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("guiGetColor expects control, property")
	}
	ctrl, _ := argInt("guiGetColor", args, 0)
	prop, _ := argInt("guiGetColor", args, 1)
	result := rg.GetColor(rg.ControlID(ctrl), rg.PropertyID(prop))
	return colorToMap(result), nil
}
