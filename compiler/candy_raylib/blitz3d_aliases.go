package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type entityState struct {
	modelID int64
	x, y, z float32
	rx, ry, rz float32
	sx, sy, sz float32
	visible bool
	parent int64
	// 0 = axis cube (default), 1 = sphere (radius ≈ sx)
	kind int32
	// Optional tint when kind is sphere/cube without model
	useTint bool
	tr, tg, tb uint8
}

var (
	entities = map[int64]*entityState{}
	nextEntityID int64 = 1
)

func builtinGraphics3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	w := int64(800)
	h := int64(600)
	if len(args) >= 2 {
		w, _ = argInt("Graphics3D", args, 0)
		h, _ = argInt("Graphics3D", args, 1)
	}
	rl.InitWindow(int32(w), int32(h), "Candy 3D")
	rl.SetTargetFPS(60)
	rl.BeginDrawing()
	candyFrameActive = true
	return null(), nil
}

func builtinCreateCamera(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	activeCamera3D = rl.Camera3D{
		Position:   rl.NewVector3(0, 10, 10),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}
	return vInt(0), nil // Camera is often 0 in simple scripts, or we can use a handle
}

func builtinCreateLight(_ []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(1), nil
}

func builtinCreateCube(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextEntityID
	nextEntityID++
	entities[id] = &entityState{
		sx: 1, sy: 1, sz: 1,
		visible: true,
	}
	// We don't necessarily load a "model" for a cube, we might just draw it in DrawEntities
	return vInt(id), nil
}

func builtinCreateSphere(_ []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextEntityID
	nextEntityID++
	entities[id] = &entityState{
		sx: 1, sy: 1, sz: 1,
		visible: true,
		kind:    1,
	}
	return vInt(id), nil
}

func builtinLoadMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	mv, err := builtinLoadModel(args)
	if err != nil {
		return nil, err
	}
	id := nextEntityID
	nextEntityID++
	entities[id] = &entityState{
		modelID: mv.I64,
		sx: 1, sy: 1, sz: 1,
		visible: true,
	}
	return vInt(id), nil
}

func builtinPositionEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("PositionEntity(id, x, y, z)")
	}
	id, _ := argInt("PositionEntity", args, 0)
	x, _ := getArgFloat("PositionEntity", args, 1)
	y, _ := getArgFloat("PositionEntity", args, 2)
	z, _ := getArgFloat("PositionEntity", args, 3)

	if ent, ok := entities[id]; ok {
		ent.x, ent.y, ent.z = float32(x), float32(y), float32(z)
	}
	return null(), nil
}

func builtinRotateEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("RotateEntity(id, x, y, z)")
	}
	id, _ := argInt("RotateEntity", args, 0)
	x, _ := getArgFloat("RotateEntity", args, 1)
	y, _ := getArgFloat("RotateEntity", args, 2)
	z, _ := getArgFloat("RotateEntity", args, 3)

	if ent, ok := entities[id]; ok {
		ent.rx, ent.ry, ent.rz = float32(x), float32(y), float32(z)
	}
	return null(), nil
}

func builtinScaleEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("ScaleEntity(id, x, y, z)")
	}
	id, _ := argInt("ScaleEntity", args, 0)
	x, _ := getArgFloat("ScaleEntity", args, 1)
	y, _ := getArgFloat("ScaleEntity", args, 2)
	z, _ := getArgFloat("ScaleEntity", args, 3)

	if ent, ok := entities[id]; ok {
		ent.sx, ent.sy, ent.sz = float32(x), float32(y), float32(z)
	}
	return null(), nil
}

func builtinMoveEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("MoveEntity(id, x, y, z)")
	}
	id, _ := argInt("MoveEntity", args, 0)
	x, _ := getArgFloat("MoveEntity", args, 1)
	y, _ := getArgFloat("MoveEntity", args, 2)
	z, _ := getArgFloat("MoveEntity", args, 3)

	if ent, ok := entities[id]; ok {
		ent.x += float32(x)
		ent.y += float32(y)
		ent.z += float32(z)
	}
	return null(), nil
}

func builtinTurnEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("TurnEntity(id, x, y, z)")
	}
	id, _ := argInt("TurnEntity", args, 0)
	x, _ := getArgFloat("TurnEntity", args, 1)
	y, _ := getArgFloat("TurnEntity", args, 2)
	z, _ := getArgFloat("TurnEntity", args, 3)

	if ent, ok := entities[id]; ok {
		ent.rx += float32(x)
		ent.ry += float32(y)
		ent.rz += float32(z)
	}
	return null(), nil
}

func builtinEntityX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("EntityX", args, 0)
	if ent, ok := entities[id]; ok {
		return vFloat(float64(ent.x)), nil
	}
	return vFloat(0), nil
}
func builtinEntityY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("EntityY", args, 0)
	if ent, ok := entities[id]; ok {
		return vFloat(float64(ent.y)), nil
	}
	return vFloat(0), nil
}
func builtinEntityZ(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("EntityZ", args, 0)
	if ent, ok := entities[id]; ok {
		return vFloat(float64(ent.z)), nil
	}
	return vFloat(0), nil
}

func builtinPositionCamera(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("PositionCamera(camera, x, y, z)")
	}
	x, _ := getArgFloat("PositionCamera", args, 1)
	y, _ := getArgFloat("PositionCamera", args, 2)
	z, _ := getArgFloat("PositionCamera", args, 3)
	activeCamera3D.Position = rl.NewVector3(float32(x), float32(y), float32(z))
	return null(), nil
}

func builtinPointCamera(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("PointCamera(camera, targetEntityId)")
	}
	tid, _ := argInt("PointCamera", args, 1)
	if ent, ok := entities[tid]; ok {
		activeCamera3D.Target = rl.NewVector3(ent.x, ent.y, ent.z)
	}
	return null(), nil
}

func builtinCameraClsColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("cameraClsColor(camera, r, g, b)")
	}
	r, _ := argInt("cameraClsColor", args, 1)
	g, _ := argInt("cameraClsColor", args, 2)
	b, _ := argInt("cameraClsColor", args, 3)
	blitzFrameClear = rl.NewColor(uint8(r), uint8(g), uint8(b), 255)
	blitzFrameClearValid = true
	return null(), nil
}

func builtinColorEntity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("ColorEntity(id, r, g, b)")
	}
	id, _ := argInt("ColorEntity", args, 0)
	r, _ := argInt("ColorEntity", args, 1)
	g, _ := argInt("ColorEntity", args, 2)
	b, _ := argInt("ColorEntity", args, 3)
	if ent, ok := entities[id]; ok {
		ent.useTint = true
		ent.tr, ent.tg, ent.tb = uint8(r), uint8(g), uint8(b)
	}
	return null(), nil
}

func builtinEntityTexture(_ []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// Reserved for future material binding; models from LoadMesh already carry materials.
	return null(), nil
}

func builtinKeyHit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("keyHit(keyNameOrCode)")
	}
	a := args[0]
	if a.Kind == candy_evaluator.ValInt {
		return vBool(rl.IsKeyPressed(int32(a.I64))), nil
	}
	if a.Kind == candy_evaluator.ValString {
		return vBool(rl.IsKeyPressed(keyCode(a.Str))), nil
	}
	return nil, fmt.Errorf("keyHit: string key name or integer key code")
}

func builtinRenderWorld(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.BeginMode3D(activeCamera3D)
	for _, ent := range entities {
		if !ent.visible {
			continue
		}
		pos := rl.NewVector3(ent.x, ent.y, ent.z)
		if ent.modelID != 0 {
			if m, ok := models[ent.modelID]; ok {
				// DrawModelEx doesn't handle all rotations easily without a matrix
				// Simple version:
				rl.DrawModelEx(m, pos, rl.NewVector3(0, 1, 0), ent.ry, rl.NewVector3(ent.sx, ent.sy, ent.sz), rl.White)
			}
		} else if ent.kind == 1 {
			col := rl.Maroon
			if ent.useTint {
				col = rl.NewColor(ent.tr, ent.tg, ent.tb, 255)
			}
			radius := ent.sx
			if radius <= 0 {
				radius = 1
			}
			rl.DrawSphere(pos, radius, col)
			rl.DrawSphereWires(pos, radius, 8, 8, rl.Black)
		} else {
			col := rl.Maroon
			if ent.useTint {
				col = rl.NewColor(ent.tr, ent.tg, ent.tb, 255)
			}
			rl.DrawCube(pos, ent.sx, ent.sy, ent.sz, col)
			rl.DrawCubeWires(pos, ent.sx, ent.sy, ent.sz, rl.Black)
		}
	}
	rl.EndMode3D()
	return null(), nil
}
