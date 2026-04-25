package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Ray helper ----

func rayToMap(ray rl.Ray) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"position": {
			Kind: candy_evaluator.ValMap,
			StrMap: map[string]candy_evaluator.Value{
				"x": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.X)},
				"y": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.Y)},
				"z": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.Z)},
			},
		},
		"direction": {
			Kind: candy_evaluator.ValMap,
			StrMap: map[string]candy_evaluator.Value{
				"x": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.X)},
				"y": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.Y)},
				"z": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.Z)},
			},
		},
	})
}

// argCamera3D reads a Camera3D from a map arg {posX,posY,posZ,targetX,targetY,targetZ,upX,upY,upZ,fovY,[projection]}
// Falls back to activeCamera3D if arg is absent or null.
func argCamera3D(name string, args []*candy_evaluator.Value, i int) (rl.Camera3D, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind == candy_evaluator.ValNull {
		return activeCamera3D, nil
	}
	if args[i].Kind != candy_evaluator.ValMap {
		return rl.Camera3D{}, fmt.Errorf("%s arg %d must be camera map", name, i+1)
	}
	m := args[i].StrMap
	return rl.Camera3D{
		Position:   rl.NewVector3(mapFloat(m, "posX"), mapFloat(m, "posY"), mapFloat(m, "posZ")),
		Target:     rl.NewVector3(mapFloat(m, "targetX"), mapFloat(m, "targetY"), mapFloat(m, "targetZ")),
		Up:         rl.NewVector3(mapFloatDefault(m, "upX", 0), mapFloatDefault(m, "upY", 1), mapFloatDefault(m, "upZ", 0)),
		Fovy:       mapFloat(m, "fovY"),
		Projection: rl.CameraProjection(int32(mapFloat(m, "projection"))),
	}, nil
}

func argCamera2D(name string, args []*candy_evaluator.Value, i int) (rl.Camera2D, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.Camera2D{}, fmt.Errorf("%s arg %d must be camera2D map {offsetX,offsetY,targetX,targetY,rotation,zoom}", name, i+1)
	}
	m := args[i].StrMap
	return rl.Camera2D{
		Offset:   rl.NewVector2(mapFloat(m, "offsetX"), mapFloat(m, "offsetY")),
		Target:   rl.NewVector2(mapFloat(m, "targetX"), mapFloat(m, "targetY")),
		Rotation: mapFloat(m, "rotation"),
		Zoom:     mapFloatDefault(m, "zoom", 1),
	}, nil
}

func mapFloatDefault(m map[string]candy_evaluator.Value, key string, def float32) float32 {
	if v, ok := m[key]; ok {
		if v.Kind == candy_evaluator.ValFloat {
			return float32(v.F64)
		}
		if v.Kind == candy_evaluator.ValInt {
			return float32(v.I64)
		}
	}
	return def
}

// ---- Drawing extras: scissor, blend ----

func builtinBeginScissorMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("beginScissorMode", args, 4); err != nil {
		return nil, err
	}
	x, _ := argInt("beginScissorMode", args, 0)
	y, _ := argInt("beginScissorMode", args, 1)
	w, _ := argInt("beginScissorMode", args, 2)
	h, _ := argInt("beginScissorMode", args, 3)
	rl.BeginScissorMode(int32(x), int32(y), int32(w), int32(h))
	return null(), nil
}

func builtinEndScissorMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EndScissorMode()
	return null(), nil
}

func builtinBeginBlendMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("beginBlendMode", args, 1); err != nil {
		return nil, err
	}
	mode, _ := argInt("beginBlendMode", args, 0)
	rl.BeginBlendMode(rl.BlendMode(int32(mode)))
	return null(), nil
}

func builtinEndBlendMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EndBlendMode()
	return null(), nil
}

// ---- Screen-space functions ----

func builtinGetScreenToWorldRay(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getScreenToWorldRay", args, 2); err != nil {
		return nil, err
	}
	pos, err := argVector2("getScreenToWorldRay", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera3D("getScreenToWorldRay", args, 1)
	if err != nil {
		return nil, err
	}
	return rayToMap(rl.GetScreenToWorldRay(pos, cam)), nil
}

func builtinGetScreenToWorldRayEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getScreenToWorldRayEx", args, 4); err != nil {
		return nil, err
	}
	pos, err := argVector2("getScreenToWorldRayEx", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera3D("getScreenToWorldRayEx", args, 1)
	if err != nil {
		return nil, err
	}
	w, _ := argInt("getScreenToWorldRayEx", args, 2)
	h, _ := argInt("getScreenToWorldRayEx", args, 3)
	return rayToMap(rl.GetScreenToWorldRayEx(pos, cam, int32(w), int32(h))), nil
}

func builtinGetWorldToScreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getWorldToScreen", args, 2); err != nil {
		return nil, err
	}
	pos, err := argVector3("getWorldToScreen", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera3D("getWorldToScreen", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.GetWorldToScreen(pos, cam)), nil
}

func builtinGetWorldToScreenEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getWorldToScreenEx", args, 4); err != nil {
		return nil, err
	}
	pos, err := argVector3("getWorldToScreenEx", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera3D("getWorldToScreenEx", args, 1)
	if err != nil {
		return nil, err
	}
	w, _ := argInt("getWorldToScreenEx", args, 2)
	h, _ := argInt("getWorldToScreenEx", args, 3)
	return vec2ToMap(rl.GetWorldToScreenEx(pos, cam, int32(w), int32(h))), nil
}

func builtinGetWorldToScreen2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getWorldToScreen2D", args, 2); err != nil {
		return nil, err
	}
	pos, err := argVector2("getWorldToScreen2D", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera2D("getWorldToScreen2D", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.GetWorldToScreen2D(pos, cam)), nil
}

func builtinGetScreenToWorld2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getScreenToWorld2D", args, 2); err != nil {
		return nil, err
	}
	pos, err := argVector2("getScreenToWorld2D", args, 0)
	if err != nil {
		return nil, err
	}
	cam, err := argCamera2D("getScreenToWorld2D", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.GetScreenToWorld2D(pos, cam)), nil
}

func builtinGetCameraMatrix(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	cam, err := argCamera3D("getCameraMatrix", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.GetCameraMatrix(cam)), nil
}

func builtinGetCameraMatrix2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getCameraMatrix2D", args, 1); err != nil {
		return nil, err
	}
	cam, err := argCamera2D("getCameraMatrix2D", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.GetCameraMatrix2D(cam)), nil
}

// ---- Frame control ----

func builtinSwapScreenBuffer(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.SwapScreenBuffer()
	return null(), nil
}

func builtinPollInputEvents(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.PollInputEvents()
	return null(), nil
}

func builtinWaitTime(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("waitTime", args, 1); err != nil {
		return nil, err
	}
	s, _ := getArgFloat("waitTime", args, 0)
	rl.WaitTime(float64(s))
	return null(), nil
}

// ---- Random ----

func builtinSetRandomSeed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setRandomSeed", args, 1); err != nil {
		return nil, err
	}
	seed, _ := argInt("setRandomSeed", args, 0)
	rand.Seed(int64(seed)) //nolint:staticcheck
	return null(), nil
}

func builtinGetRandomValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getRandomValue", args, 2); err != nil {
		return nil, err
	}
	mn, _ := argInt("getRandomValue", args, 0)
	mx, _ := argInt("getRandomValue", args, 1)
	return vInt(int64(rl.GetRandomValue(int32(mn), int32(mx)))), nil
}

// ---- Misc ----

func builtinSetConfigFlags(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setConfigFlags", args, 1); err != nil {
		return nil, err
	}
	flags, _ := argInt("setConfigFlags", args, 0)
	rl.SetConfigFlags(uint32(flags))
	return null(), nil
}

func builtinOpenURL(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("openURL", args, 1); err != nil {
		return nil, err
	}
	url, err := argString("openURL", args, 0)
	if err != nil {
		return nil, err
	}
	rl.OpenURL(url)
	return null(), nil
}

func builtinSetTraceLogLevel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setTraceLogLevel", args, 1); err != nil {
		return nil, err
	}
	level, _ := argInt("setTraceLogLevel", args, 0)
	rl.SetTraceLogLevel(rl.TraceLogLevel(int32(level)))
	return null(), nil
}

// ---- Shader extras ----

func builtinLoadShaderFromMemory(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadShaderFromMemory", args, 2); err != nil {
		return nil, err
	}
	vs, _ := argString("loadShaderFromMemory", args, 0)
	fs, _ := argString("loadShaderFromMemory", args, 1)
	s := rl.LoadShaderFromMemory(vs, fs)
	id := nextShaderID
	nextShaderID++
	shaders[id] = s
	return vInt(id), nil
}

func builtinIsShaderValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := shaderByID("isShaderValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsShaderValid(s)), nil
}

func builtinGetShaderLocation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getShaderLocation", args, 2); err != nil {
		return nil, err
	}
	_, s, err := shaderByID("getShaderLocation", args, 0)
	if err != nil {
		return nil, err
	}
	name, err := argString("getShaderLocation", args, 1)
	if err != nil {
		return nil, err
	}
	return vInt(int64(rl.GetShaderLocation(s, name))), nil
}

func builtinGetShaderLocationAttrib(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getShaderLocationAttrib", args, 2); err != nil {
		return nil, err
	}
	_, s, err := shaderByID("getShaderLocationAttrib", args, 0)
	if err != nil {
		return nil, err
	}
	name, err := argString("getShaderLocationAttrib", args, 1)
	if err != nil {
		return nil, err
	}
	return vInt(int64(rl.GetShaderLocationAttrib(s, name))), nil
}

func builtinSetShaderValueMatrix(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setShaderValueMatrix", args, 3); err != nil {
		return nil, err
	}
	_, s, err := shaderByID("setShaderValueMatrix", args, 0)
	if err != nil {
		return nil, err
	}
	loc, _ := argInt("setShaderValueMatrix", args, 1)
	mat, err := argMatrix("setShaderValueMatrix", args, 2)
	if err != nil {
		return nil, err
	}
	rl.SetShaderValueMatrix(s, int32(loc), mat)
	return null(), nil
}

func builtinSetShaderValueTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setShaderValueTexture", args, 3); err != nil {
		return nil, err
	}
	_, s, err := shaderByID("setShaderValueTexture", args, 0)
	if err != nil {
		return nil, err
	}
	loc, _ := argInt("setShaderValueTexture", args, 1)
	texID, _ := argInt("setShaderValueTexture", args, 2)
	tex, ok := textures[texID]
	if !ok {
		return nil, fmt.Errorf("setShaderValueTexture: invalid texture handle %d", texID)
	}
	rl.SetShaderValueTexture(s, int32(loc), tex)
	return null(), nil
}

// SetShaderValue: takes shaderId, loc, []float values, uniformType
func builtinSetShaderValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setShaderValue", args, 4); err != nil {
		return nil, err
	}
	_, s, err := shaderByID("setShaderValue", args, 0)
	if err != nil {
		return nil, err
	}
	loc, _ := argInt("setShaderValue", args, 1)
	// arg 2: array of floats or single float/int
	var vals []float32
	v := args[2]
	switch v.Kind {
	case candy_evaluator.ValArray:
		for _, elem := range v.Elems {
			switch elem.Kind {
			case candy_evaluator.ValFloat:
				vals = append(vals, float32(elem.F64))
			case candy_evaluator.ValInt:
				vals = append(vals, float32(elem.I64))
			}
		}
	case candy_evaluator.ValFloat:
		vals = []float32{float32(v.F64)}
	case candy_evaluator.ValInt:
		vals = []float32{float32(v.I64)}
	}
	uniformType, _ := argInt("setShaderValue", args, 3)
	rl.SetShaderValue(s, int32(loc), vals, rl.ShaderUniformDataType(int32(uniformType)))
	return null(), nil
}

// ---- Keyboard extras ----

func builtinIsKeyPressedRepeat(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isKeyPressedRepeat", args, 1); err != nil {
		return nil, err
	}
	k, err := keyArg("isKeyPressedRepeat", args[0])
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsKeyPressedRepeat(k)), nil
}

func builtinGetKeyPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetKeyPressed())), nil
}

func builtinGetCharPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetCharPressed())), nil
}

// ---- Mouse extras ----

func builtinIsMouseButtonUp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMouseButtonUp", args, 1); err != nil {
		return nil, err
	}
	b, _ := argInt("isMouseButtonUp", args, 0)
	return vBool(rl.IsMouseButtonUp(mouseButtonCode(b))), nil
}

func builtinGetMouseDelta(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.GetMouseDelta()), nil
}

func builtinSetMousePosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMousePosition", args, 2); err != nil {
		return nil, err
	}
	x, _ := argInt("setMousePosition", args, 0)
	y, _ := argInt("setMousePosition", args, 1)
	rl.SetMousePosition(int(x), int(y))
	return null(), nil
}

func builtinSetMouseOffset(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMouseOffset", args, 2); err != nil {
		return nil, err
	}
	x, _ := argInt("setMouseOffset", args, 0)
	y, _ := argInt("setMouseOffset", args, 1)
	rl.SetMouseOffset(int(x), int(y))
	return null(), nil
}

func builtinSetMouseScale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMouseScale", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("setMouseScale", args, 0)
	y, _ := getArgFloat("setMouseScale", args, 1)
	rl.SetMouseScale(float32(x), float32(y))
	return null(), nil
}

func builtinGetMouseWheelMoveV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.GetMouseWheelMoveV()), nil
}

func builtinSetMouseCursor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMouseCursor", args, 1); err != nil {
		return nil, err
	}
	c, _ := argInt("setMouseCursor", args, 0)
	rl.SetMouseCursor(int32(c))
	return null(), nil
}

// ---- Gamepad extras ----

func builtinIsGamepadButtonReleased(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGamepadButtonReleased", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("isGamepadButtonReleased", args, 0)
	btn, _ := argInt("isGamepadButtonReleased", args, 1)
	return vBool(rl.IsGamepadButtonReleased(int32(id), int32(btn))), nil
}

func builtinIsGamepadButtonUp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGamepadButtonUp", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("isGamepadButtonUp", args, 0)
	btn, _ := argInt("isGamepadButtonUp", args, 1)
	return vBool(rl.IsGamepadButtonUp(int32(id), int32(btn))), nil
}

func builtinGetGamepadButtonPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetGamepadButtonPressed())), nil
}

func builtinGetGamepadAxisCount(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getGamepadAxisCount", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("getGamepadAxisCount", args, 0)
	return vInt(int64(rl.GetGamepadAxisCount(int32(id)))), nil
}

func builtinSetGamepadMappings(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setGamepadMappings", args, 1); err != nil {
		return nil, err
	}
	m, err := argString("setGamepadMappings", args, 0)
	if err != nil {
		return nil, err
	}
	return vInt(int64(rl.SetGamepadMappings(m))), nil
}

func builtinSetGamepadVibration(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setGamepadVibration", args, 4); err != nil {
		return nil, err
	}
	id, _ := argInt("setGamepadVibration", args, 0)
	l, _ := getArgFloat("setGamepadVibration", args, 1)
	r, _ := getArgFloat("setGamepadVibration", args, 2)
	d, _ := getArgFloat("setGamepadVibration", args, 3)
	rl.SetGamepadVibration(int32(id), float32(l), float32(r), float32(d))
	return null(), nil
}

// ---- Touch ----

func builtinGetTouchX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetTouchX())), nil
}

func builtinGetTouchY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetTouchY())), nil
}

func builtinGetTouchPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getTouchPosition", args, 1); err != nil {
		return nil, err
	}
	idx, _ := argInt("getTouchPosition", args, 0)
	return vec2ToMap(rl.GetTouchPosition(int32(idx))), nil
}

func builtinGetTouchPointId(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getTouchPointId", args, 1); err != nil {
		return nil, err
	}
	idx, _ := argInt("getTouchPointId", args, 0)
	return vInt(int64(rl.GetTouchPointId(int32(idx)))), nil
}

func builtinGetTouchPointCount(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetTouchPointCount())), nil
}

// ---- Gestures ----

func builtinSetGesturesEnabled(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setGesturesEnabled", args, 1); err != nil {
		return nil, err
	}
	flags, _ := argInt("setGesturesEnabled", args, 0)
	rl.SetGesturesEnabled(uint32(flags))
	return null(), nil
}

func builtinIsGestureDetected(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isGestureDetected", args, 1); err != nil {
		return nil, err
	}
	g, _ := argInt("isGestureDetected", args, 0)
	return vBool(rl.IsGestureDetected(rl.Gestures(g))), nil
}

func builtinGetGestureDetected(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(rl.GetGestureDetected())), nil
}

func builtinGetGestureHoldDuration(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetGestureHoldDuration())), nil
}

func builtinGetGestureDragVector(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.GetGestureDragVector()), nil
}

func builtinGetGestureDragAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetGestureDragAngle())), nil
}

func builtinGetGesturePinchVector(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.GetGesturePinchVector()), nil
}

func builtinGetGesturePinchAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetGesturePinchAngle())), nil
}

// ---- Camera update (rcamera) ----

func builtinUpdateCamera(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateCamera", args, 1); err != nil {
		return nil, err
	}
	mode, _ := argInt("updateCamera", args, 0)
	rl.UpdateCamera(&activeCamera3D, rl.CameraMode(int32(mode)))
	return null(), nil
}

func builtinUpdateCameraPro(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateCameraPro", args, 3); err != nil {
		return nil, err
	}
	movement, err := argVector3("updateCameraPro", args, 0)
	if err != nil {
		return nil, err
	}
	rotation, err := argVector3("updateCameraPro", args, 1)
	if err != nil {
		return nil, err
	}
	zoom, _ := getArgFloat("updateCameraPro", args, 2)
	rl.UpdateCameraPro(&activeCamera3D, movement, rotation, float32(zoom))
	return null(), nil
}

// ---- File dropped ----

func builtinIsFileDropped(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsFileDropped()), nil
}

func builtinLoadDroppedFiles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	files := rl.LoadDroppedFiles()
	elems := make([]candy_evaluator.Value, len(files))
	for i, f := range files {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: f}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}
