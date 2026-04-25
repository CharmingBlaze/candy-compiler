package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// candyFrameActive is true after BeginDrawing until EndDrawing. clear() and flip() coordinate so a frame can start without a stray double-Begin.
var candyFrameActive bool

func windowFlagFromString(name string) (uint32, bool) {
	switch strings.ToLower(name) {
	case "vsync", "flag_vsync_hint":
		return rl.FlagVsyncHint, true
	case "fullscreen", "flag_fullscreen_mode":
		return rl.FlagFullscreenMode, true
	case "resizable", "window_resizable", "flag_window_resizable":
		return rl.FlagWindowResizable, true
	case "undecorated", "window_undecorated", "flag_window_undecorated":
		return rl.FlagWindowUndecorated, true
	case "hidden", "window_hidden", "flag_window_hidden":
		return rl.FlagWindowHidden, true
	case "minimized", "window_minimized", "flag_window_minimized":
		return rl.FlagWindowMinimized, true
	case "maximized", "window_maximized", "flag_window_maximized":
		return rl.FlagWindowMaximized, true
	case "unfocused", "window_unfocused", "flag_window_unfocused":
		return rl.FlagWindowUnfocused, true
	case "topmost", "window_topmost", "flag_window_topmost":
		return rl.FlagWindowTopmost, true
	case "always_run", "window_always_run", "flag_window_always_run":
		return rl.FlagWindowAlwaysRun, true
	case "transparent", "window_transparent", "flag_window_transparent":
		return rl.FlagWindowTransparent, true
	case "highdpi", "window_highdpi", "flag_window_highdpi":
		return rl.FlagWindowHighdpi, true
	case "mouse_passthrough", "window_mouse_passthrough", "flag_window_mouse_passthrough":
		return rl.FlagWindowMousePassthrough, true
	case "borderless", "borderless_windowed_mode", "flag_borderless_windowed_mode":
		return rl.FlagBorderlessWindowedMode, true
	case "msaa4x", "flag_msaa_4x_hint", "flag_msaa4x_hint":
		return rl.FlagMsaa4xHint, true
	case "interlaced", "flag_interlaced_hint":
		return rl.FlagInterlacedHint, true
	default:
		return 0, false
	}
}

func parseWindowFlagsArg(name string, v *candy_evaluator.Value) (uint32, error) {
	if v == nil {
		return 0, fmt.Errorf("%s: flags argument is required", name)
	}
	switch v.Kind {
	case candy_evaluator.ValInt:
		return uint32(v.I64), nil
	case candy_evaluator.ValString:
		f, ok := windowFlagFromString(v.Str)
		if !ok {
			return 0, fmt.Errorf("%s: unknown window flag %q", name, v.Str)
		}
		return f, nil
	case candy_evaluator.ValArray:
		var out uint32
		for i := range v.Elems {
			elem := v.Elems[i]
			flag, err := parseWindowFlagsArg(name, &elem)
			if err != nil {
				return 0, err
			}
			out |= flag
		}
		return out, nil
	default:
		return 0, fmt.Errorf("%s: flags must be int, string, or [flags]", name)
	}
}

// ---- Window-related functions ----

func builtinInitWindow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("initWindow", args, 3); err != nil {
		return nil, err
	}
	w, _ := argInt("initWindow", args, 0)
	h, _ := argInt("initWindow", args, 1)
	title, err := argString("initWindow", args, 2)
	if err != nil {
		return nil, err
	}
	candyFrameActive = false
	rl.InitWindow(int32(w), int32(h), title)
	rl.SetTargetFPS(60)
	return null(), nil
}

func builtinCloseWindow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	candyFrameActive = false
	rl.CloseWindow()
	return null(), nil
}

func builtinIsWindowReady(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowReady()), nil
}

func builtinIsWindowFullscreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowFullscreen()), nil
}

func builtinToggleFullscreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.ToggleFullscreen()
	return null(), nil
}

func builtinSetWindowTitle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowTitle", args, 1); err != nil {
		return nil, err
	}
	title, err := argString("setWindowTitle", args, 0)
	if err != nil {
		return nil, err
	}
	rl.SetWindowTitle(title)
	return null(), nil
}

func builtinSetWindowSize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowSize", args, 2); err != nil {
		return nil, err
	}
	w, _ := argInt("setWindowSize", args, 0)
	h, _ := argInt("setWindowSize", args, 1)
	rl.SetWindowSize(int(w), int(h))
	return null(), nil
}

func builtinMinimizeWindow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.MinimizeWindow()
	return null(), nil
}

func builtinMaximizeWindow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.MaximizeWindow()
	return null(), nil
}

func builtinRestoreWindow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.RestoreWindow()
	return null(), nil
}

func builtinSetTargetFPS(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setTargetFPS", args, 1); err != nil {
		return nil, err
	}
	fps, _ := argInt("setTargetFPS", args, 0)
	rl.SetTargetFPS(int32(fps))
	return null(), nil
}

func builtinGetFPS(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetFPS())), nil
}

func builtinGetFrameTime(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetFrameTime())), nil
}

func builtinGetTime(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(rl.GetTime()), nil
}

func builtinGetScreenWidth(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetScreenWidth())), nil
}

func builtinGetScreenHeight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetScreenHeight())), nil
}

// ---- Drawing-related functions ----

func builtinBeginDrawing(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if !candyFrameActive {
		rl.BeginDrawing()
		candyFrameActive = true
	}
	return null(), nil
}

func builtinEndDrawing(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if candyFrameActive {
		// `every` / `after` from game helpers need one tick per frame (same dt as the completed frame).
		helperTick(float64(rl.GetFrameTime()))
		rl.EndDrawing()
		candyFrameActive = false
	}
	return null(), nil
}

func builtinClearBackground(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// One frame: clear() = Begin (if needed) + clear color; then draw; then show() = End.
	if !candyFrameActive {
		rl.BeginDrawing()
		candyFrameActive = true
	}
	c, err := argColor("clearBackground", args, 0, rl.Black)
	if err != nil {
		return nil, err
	}
	rl.ClearBackground(c)
	return null(), nil
}

func builtinWindowShouldClose(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.WindowShouldClose()), nil
}

// ---- Input-related functions: Keyboard ----

func builtinIsKeyPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isKeyPressed", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("isKeyPressed", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsKeyPressed(k)), nil
}

func builtinIsKeyDown(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isKeyDown", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("isKeyDown", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsKeyDown(k)), nil
}

func builtinIsKeyReleased(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isKeyReleased", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("isKeyReleased", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsKeyReleased(k)), nil
}

func builtinIsKeyUp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isKeyUp", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("isKeyUp", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsKeyUp(k)), nil
}

func builtinSetExitKey(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setExitKey", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("setExitKey", args[0])
	if err != nil {
		return nil, err
	}
	rl.SetExitKey(k)
	return null(), nil
}

func builtinClicked(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsMouseButtonPressed(rl.MouseLeftButton)), nil
}

// ---- Input-related functions: Mouse ----

func builtinIsMouseButtonPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMouseButtonPressed", args, 1); err != nil {
		return nil, err
	}
	b, _ := argInt("isMouseButtonPressed", args, 0)
	return vBool(rl.IsMouseButtonPressed(mouseButtonCode(b))), nil
}

func builtinIsMouseButtonDown(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMouseButtonDown", args, 1); err != nil {
		return nil, err
	}
	b, _ := argInt("isMouseButtonDown", args, 0)
	return vBool(rl.IsMouseButtonDown(mouseButtonCode(b))), nil
}

func builtinIsMouseButtonReleased(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMouseButtonReleased", args, 1); err != nil {
		return nil, err
	}
	b, _ := argInt("isMouseButtonReleased", args, 0)
	return vBool(rl.IsMouseButtonReleased(mouseButtonCode(b))), nil
}

func builtinGetMouseX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetMouseX())), nil
}

func builtinGetMouseY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetMouseY())), nil
}

func builtinGetMousePosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	pos := rl.GetMousePosition()
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(pos.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(pos.Y)},
	}), nil
}

func builtinGetMouseWheelMove(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetMouseWheelMove())), nil
}

func builtinSetClipboardText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setClipboardText", args, 1); err != nil {
		return nil, err
	}
	t, err := argString("setClipboardText", args, 0)
	if err != nil {
		return nil, err
	}
	rl.SetClipboardText(t)
	return null(), nil
}

func builtinGetClipboardText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: rl.GetClipboardText()}, nil
}

func builtinFlip(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if candyFrameActive {
		helperTick(float64(rl.GetFrameTime()))
		rl.EndDrawing()
	}
	rl.BeginDrawing()
	candyFrameActive = true
	if blitzFrameClearValid {
		rl.ClearBackground(blitzFrameClear)
	}
	return null(), nil
}

func builtinKey(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsKeyDown(args)
}

func builtinShouldClose(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinWindowShouldClose(args)
}

// ---- Gamepad Support ----

func builtinIsGamepadAvailable(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGamepadAvailable", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("isGamepadAvailable", args, 0)
	return vBool(rl.IsGamepadAvailable(int32(id))), nil
}

func builtinIsGamepadButtonPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGamepadButtonPressed", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("isGamepadButtonPressed", args, 0)
	btn, _ := argInt("isGamepadButtonPressed", args, 1)
	return vBool(rl.IsGamepadButtonPressed(int32(id), int32(btn))), nil
}

func builtinIsGamepadButtonDown(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGamepadButtonDown", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("isGamepadButtonDown", args, 0)
	btn, _ := argInt("isGamepadButtonDown", args, 1)
	return vBool(rl.IsGamepadButtonDown(int32(id), int32(btn))), nil
}

func builtinGetGamepadName(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGamepadName", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("getGamepadName", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: rl.GetGamepadName(int32(id))}, nil
}

func builtinGetGamepadAxisValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGamepadAxisValue", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("getGamepadAxisValue", args, 0)
	axis, _ := argInt("getGamepadAxisValue", args, 1)
	return vFloat(float64(rl.GetGamepadAxisMovement(int32(id), int32(axis)))), nil
}

// ---- Screen Management ----

func builtinTakeScreenshot(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("takeScreenshot", args, 1); err != nil {
		return nil, err
	}
	path, err := argString("takeScreenshot", args, 0)
	if err != nil {
		return nil, err
	}
	rl.TakeScreenshot(path)
	return null(), nil
}

// ---- Window state checks ----

func builtinIsWindowHidden(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowHidden()), nil
}

func builtinIsWindowMinimized(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowMinimized()), nil
}

func builtinIsWindowMaximized(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowMaximized()), nil
}

func builtinIsWindowFocused(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowFocused()), nil
}

func builtinIsWindowResized(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsWindowResized()), nil
}

func builtinIsWindowState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isWindowState", args, 1); err != nil {
		return nil, err
	}
	flags, err := parseWindowFlagsArg("isWindowState", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsWindowState(flags)), nil
}

// ---- Window state mutators ----

func builtinSetWindowState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowState", args, 1); err != nil {
		return nil, err
	}
	flags, err := parseWindowFlagsArg("setWindowState", args[0])
	if err != nil {
		return nil, err
	}
	rl.SetWindowState(flags)
	return null(), nil
}

func builtinClearWindowState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("clearWindowState", args, 1); err != nil {
		return nil, err
	}
	flags, err := parseWindowFlagsArg("clearWindowState", args[0])
	if err != nil {
		return nil, err
	}
	rl.ClearWindowState(flags)
	return null(), nil
}

func builtinToggleBorderlessWindowed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.ToggleBorderlessWindowed()
	return null(), nil
}

func builtinSetWindowPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowPosition", args, 2); err != nil {
		return nil, err
	}
	x, _ := argInt("setWindowPosition", args, 0)
	y, _ := argInt("setWindowPosition", args, 1)
	rl.SetWindowPosition(int(x), int(y))
	return null(), nil
}

func builtinSetWindowMonitor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowMonitor", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("setWindowMonitor", args, 0)
	rl.SetWindowMonitor(int(m))
	return null(), nil
}

func builtinSetWindowMinSize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowMinSize", args, 2); err != nil {
		return nil, err
	}
	w, _ := argInt("setWindowMinSize", args, 0)
	h, _ := argInt("setWindowMinSize", args, 1)
	rl.SetWindowMinSize(int(w), int(h))
	return null(), nil
}

func builtinSetWindowMaxSize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowMaxSize", args, 2); err != nil {
		return nil, err
	}
	w, _ := argInt("setWindowMaxSize", args, 0)
	h, _ := argInt("setWindowMaxSize", args, 1)
	rl.SetWindowMaxSize(int(w), int(h))
	return null(), nil
}

func builtinSetWindowOpacity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowOpacity", args, 1); err != nil {
		return nil, err
	}
	op, _ := getArgFloat("setWindowOpacity", args, 0)
	rl.SetWindowOpacity(float32(op))
	return null(), nil
}

func builtinSetWindowFocused(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.SetWindowFocused()
	return null(), nil
}

// ---- Window icon ----

func builtinSetWindowIcon(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowIcon", args, 1); err != nil {
		return nil, err
	}
	id, err := argInt("setWindowIcon", args, 0)
	if err != nil {
		return nil, err
	}
	img, ok := images[id]
	if !ok {
		return nil, fmt.Errorf("setWindowIcon: invalid image handle %d", id)
	}
	rl.SetWindowIcon(*img)
	return null(), nil
}

func builtinSetWindowIcons(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setWindowIcons", args, 1); err != nil {
		return nil, err
	}
	v := args[0]
	if v == nil || v.Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("setWindowIcons: expected array of imageIds")
	}
	imgs := make([]rl.Image, 0, len(v.Elems))
	for i := range v.Elems {
		elem := v.Elems[i]
		id, err := argInt("setWindowIcons", []*candy_evaluator.Value{&elem}, 0)
		if err != nil {
			return nil, fmt.Errorf("setWindowIcons: element %d: %w", i, err)
		}
		img, ok := images[id]
		if !ok {
			return nil, fmt.Errorf("setWindowIcons: invalid image handle %d at index %d", id, i)
		}
		imgs = append(imgs, *img)
	}
	rl.SetWindowIcons(imgs, int32(len(imgs)))
	return null(), nil
}

// ---- Render size ----

func builtinGetRenderWidth(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetRenderWidth())), nil
}

func builtinGetRenderHeight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetRenderHeight())), nil
}

// ---- Monitor queries ----

func builtinGetMonitorCount(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetMonitorCount())), nil
}

func builtinGetCurrentMonitor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetCurrentMonitor())), nil
}

func builtinGetMonitorPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorPosition", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorPosition", args, 0)
	pos := rl.GetMonitorPosition(int(m))
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(pos.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(pos.Y)},
	}), nil
}

func builtinGetMonitorWidth(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorWidth", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorWidth", args, 0)
	return vInt(int64(rl.GetMonitorWidth(int(m)))), nil
}

func builtinGetMonitorHeight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorHeight", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorHeight", args, 0)
	return vInt(int64(rl.GetMonitorHeight(int(m)))), nil
}

func builtinGetMonitorPhysicalWidth(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorPhysicalWidth", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorPhysicalWidth", args, 0)
	return vInt(int64(rl.GetMonitorPhysicalWidth(int(m)))), nil
}

func builtinGetMonitorPhysicalHeight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorPhysicalHeight", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorPhysicalHeight", args, 0)
	return vInt(int64(rl.GetMonitorPhysicalHeight(int(m)))), nil
}

func builtinGetMonitorRefreshRate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorRefreshRate", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorRefreshRate", args, 0)
	return vInt(int64(rl.GetMonitorRefreshRate(int(m)))), nil
}

func builtinGetMonitorName(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMonitorName", args, 1); err != nil {
		return nil, err
	}
	m, _ := argInt("getMonitorName", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: rl.GetMonitorName(int(m))}, nil
}

// ---- Window position / DPI / handle ----

func builtinGetWindowPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	pos := rl.GetWindowPosition()
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(pos.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(pos.Y)},
	}), nil
}

func builtinGetWindowScaleDPI(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	dpi := rl.GetWindowScaleDPI()
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(dpi.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(dpi.Y)},
	}), nil
}

func builtinGetWindowHandle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: fmt.Sprintf("%p", rl.GetWindowHandle())}, nil
}

// ---- Clipboard image ----

func builtinGetClipboardImage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	img := rl.GetClipboardImage()
	id := nextImageID
	nextImageID++
	images[id] = &img
	return vInt(id), nil
}

// ---- Event waiting ----

func builtinEnableEventWaiting(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EnableEventWaiting()
	return null(), nil
}

func builtinDisableEventWaiting(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.DisableEventWaiting()
	return null(), nil
}

// ---- Cursor ----

func builtinShowCursor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.ShowCursor()
	return null(), nil
}

func builtinHideCursor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.HideCursor()
	return null(), nil
}

func builtinIsCursorHidden(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsCursorHidden()), nil
}

func builtinEnableCursor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EnableCursor()
	return null(), nil
}

func builtinDisableCursor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.DisableCursor()
	return null(), nil
}

func builtinIsCursorOnScreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsCursorOnScreen()), nil
}

// ---- Camera Support ----

func builtinBeginMode2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("beginMode2D", args, 4); err != nil {
		return nil, err
	}
	offset, _ := argVector2("beginMode2D", args, 0)
	target, _ := argVector2("beginMode2D", args, 1)
	rotation, _ := getArgFloat("beginMode2D", args, 2)
	zoom, _ := getArgFloat("beginMode2D", args, 3)
	cam := rl.NewCamera2D(offset, target, float32(rotation), float32(zoom))
	rl.BeginMode2D(cam)
	return null(), nil
}

func builtinEndMode2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EndMode2D()
	return null(), nil
}

// ---- VR Support ----

func builtinBeginVrSimulatorMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// Note: VR support requires a VR config. For now we provide a stub.
	return null(), nil
}

func builtinEndVrSimulatorMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// rl.EndVrSimulatorMode()
	return null(), nil
}
