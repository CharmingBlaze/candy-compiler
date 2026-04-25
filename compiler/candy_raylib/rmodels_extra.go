package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- helpers ----

func meshByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Mesh, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return 0, rl.Mesh{}, err
	}
	m, ok := meshes[id]
	if !ok {
		return 0, rl.Mesh{}, fmt.Errorf("%s: invalid mesh handle %d", name, id)
	}
	return id, m, nil
}

func materialByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Material, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return 0, rl.Material{}, err
	}
	m, ok := materials[id]
	if !ok {
		return 0, rl.Material{}, fmt.Errorf("%s: invalid material handle %d", name, id)
	}
	return id, m, nil
}

func rayCollisionToMap(rc rl.RayCollision) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"hit":      {Kind: candy_evaluator.ValBool, B: rc.Hit},
		"distance": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Distance)},
		"point": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Point.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Point.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Point.Z)},
		}},
		"normal": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Normal.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Normal.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(rc.Normal.Z)},
		}},
	})
}

func boundingBoxToMap(bb rl.BoundingBox) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"min": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Min.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Min.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Min.Z)},
		}},
		"max": {Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Max.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Max.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(bb.Max.Z)},
		}},
	})
}

func argBoundingBox(name string, args []*candy_evaluator.Value, i int) (rl.BoundingBox, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.BoundingBox{}, fmt.Errorf("%s arg %d must be boundingBox map {min:{x,y,z}, max:{x,y,z}}", name, i+1)
	}
	m := args[i].StrMap
	minMap := m["min"].StrMap
	maxMap := m["max"].StrMap
	return rl.NewBoundingBox(
		rl.NewVector3(mapFloat(minMap, "x"), mapFloat(minMap, "y"), mapFloat(minMap, "z")),
		rl.NewVector3(mapFloat(maxMap, "x"), mapFloat(maxMap, "y"), mapFloat(maxMap, "z")),
	), nil
}

func argRay(name string, args []*candy_evaluator.Value, i int) (rl.Ray, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.Ray{}, fmt.Errorf("%s arg %d must be ray map {position:{x,y,z}, direction:{x,y,z}}", name, i+1)
	}
	m := args[i].StrMap
	posMap := m["position"].StrMap
	dirMap := m["direction"].StrMap
	return rl.Ray{
		Position:  rl.NewVector3(mapFloat(posMap, "x"), mapFloat(posMap, "y"), mapFloat(posMap, "z")),
		Direction: rl.NewVector3(mapFloat(dirMap, "x"), mapFloat(dirMap, "y"), mapFloat(dirMap, "z")),
	}, nil
}

// ---- Basic 3D drawing extras ----

func builtinDrawPoint3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("drawPoint3D expects x, y, z, [color]")
	}
	x, _ := getArgFloat("drawPoint3D", args, 0)
	y, _ := getArgFloat("drawPoint3D", args, 1)
	z, _ := getArgFloat("drawPoint3D", args, 2)
	c, _ := argColorValue("drawPoint3D", args, 3, rl.RayWhite)
	rl.DrawPoint3D(rl.NewVector3(float32(x), float32(y), float32(z)), c)
	return null(), nil
}

func builtinDrawCircle3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("drawCircle3D expects cx, cy, cz, radius, axisX, axisY, axisZ, angle, [color]")
	}
	cx, _ := getArgFloat("drawCircle3D", args, 0)
	cy, _ := getArgFloat("drawCircle3D", args, 1)
	cz, _ := getArgFloat("drawCircle3D", args, 2)
	r, _ := getArgFloat("drawCircle3D", args, 3)
	ax, _ := getArgFloat("drawCircle3D", args, 4)
	ay, _ := getArgFloat("drawCircle3D", args, 5)
	az, _ := getArgFloat("drawCircle3D", args, 6)
	angle, _ := getArgFloat("drawCircle3D", args, 7)
	c, _ := argColorValue("drawCircle3D", args, 8, rl.RayWhite)
	rl.DrawCircle3D(
		rl.NewVector3(float32(cx), float32(cy), float32(cz)),
		float32(r),
		rl.NewVector3(float32(ax), float32(ay), float32(az)),
		float32(angle),
		c,
	)
	return null(), nil
}

func builtinDrawTriangle3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 9 {
		return nil, fmt.Errorf("drawTriangle3D expects x1,y1,z1, x2,y2,z2, x3,y3,z3, [color]")
	}
	v1x, _ := getArgFloat("drawTriangle3D", args, 0)
	v1y, _ := getArgFloat("drawTriangle3D", args, 1)
	v1z, _ := getArgFloat("drawTriangle3D", args, 2)
	v2x, _ := getArgFloat("drawTriangle3D", args, 3)
	v2y, _ := getArgFloat("drawTriangle3D", args, 4)
	v2z, _ := getArgFloat("drawTriangle3D", args, 5)
	v3x, _ := getArgFloat("drawTriangle3D", args, 6)
	v3y, _ := getArgFloat("drawTriangle3D", args, 7)
	v3z, _ := getArgFloat("drawTriangle3D", args, 8)
	c, _ := argColorValue("drawTriangle3D", args, 9, rl.RayWhite)
	rl.DrawTriangle3D(
		rl.NewVector3(float32(v1x), float32(v1y), float32(v1z)),
		rl.NewVector3(float32(v2x), float32(v2y), float32(v2z)),
		rl.NewVector3(float32(v3x), float32(v3y), float32(v3z)),
		c,
	)
	return null(), nil
}

func builtinDrawCubeV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawCubeV expects position{x,y,z}, size{x,y,z}, [color]")
	}
	pos, err := argVector3("drawCubeV", args, 0)
	if err != nil {
		return nil, err
	}
	size, err2 := argVector3("drawCubeV", args, 1)
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("drawCubeV", args, 2, rl.Maroon)
	rl.DrawCubeV(pos, size, c)
	return null(), nil
}

func builtinDrawCubeWiresV(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawCubeWiresV expects position{x,y,z}, size{x,y,z}, [color]")
	}
	pos, err := argVector3("drawCubeWiresV", args, 0)
	if err != nil {
		return nil, err
	}
	size, err2 := argVector3("drawCubeWiresV", args, 1)
	if err2 != nil {
		return nil, err2
	}
	c, _ := argColorValue("drawCubeWiresV", args, 2, rl.DarkGray)
	rl.DrawCubeWiresV(pos, size, c)
	return null(), nil
}

func builtinDrawCylinderWiresEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("drawCylinderWiresEx expects x1,y1,z1, x2,y2,z2, startRadius, endRadius, [sides], [color]")
	}
	x1, _ := getArgFloat("drawCylinderWiresEx", args, 0)
	y1, _ := getArgFloat("drawCylinderWiresEx", args, 1)
	z1, _ := getArgFloat("drawCylinderWiresEx", args, 2)
	x2, _ := getArgFloat("drawCylinderWiresEx", args, 3)
	y2, _ := getArgFloat("drawCylinderWiresEx", args, 4)
	z2, _ := getArgFloat("drawCylinderWiresEx", args, 5)
	sr, _ := getArgFloat("drawCylinderWiresEx", args, 6)
	er, _ := getArgFloat("drawCylinderWiresEx", args, 7)
	sides := int64(16)
	if len(args) > 8 {
		sides, _ = argInt("drawCylinderWiresEx", args, 8)
	}
	c, _ := argColorValue("drawCylinderWiresEx", args, 9, rl.DarkGray)
	rl.DrawCylinderWiresEx(
		rl.NewVector3(float32(x1), float32(y1), float32(z1)),
		rl.NewVector3(float32(x2), float32(y2), float32(z2)),
		float32(sr), float32(er), int32(sides), c,
	)
	return null(), nil
}

// ---- Model extras ----

func builtinIsModelValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isModelValid", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("isModelValid", args, 0)
	m, ok := models[id]
	if !ok {
		return vBool(false), nil
	}
	return vBool(rl.IsModelValid(m)), nil
}

func builtinGetModelBoundingBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getModelBoundingBox", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("getModelBoundingBox", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("getModelBoundingBox: invalid model %d", id)
	}
	return boundingBoxToMap(rl.GetModelBoundingBox(m)), nil
}

func builtinLoadModelFromMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadModelFromMesh", args, 1); err != nil {
		return nil, err
	}
	_, mesh, err := meshByID("loadModelFromMesh", args, 0)
	if err != nil {
		return nil, err
	}
	m := rl.LoadModelFromMesh(mesh)
	id := nextModelID
	nextModelID++
	models[id] = m
	return vInt(id), nil
}

// ---- Mesh management ----

func builtinGetMeshBoundingBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getMeshBoundingBox", args, 1); err != nil {
		return nil, err
	}
	_, mesh, err := meshByID("getMeshBoundingBox", args, 0)
	if err != nil {
		return nil, err
	}
	return boundingBoxToMap(rl.GetMeshBoundingBox(mesh)), nil
}

func builtinUnloadMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadMesh", args, 1); err != nil {
		return nil, err
	}
	id, mesh, err := meshByID("unloadMesh", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadMesh(&mesh)
	delete(meshes, id)
	return null(), nil
}

func builtinExportMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("exportMesh", args, 2); err != nil {
		return nil, err
	}
	_, mesh, err := meshByID("exportMesh", args, 0)
	if err != nil {
		return nil, err
	}
	path, _ := argString("exportMesh", args, 1)
	rl.ExportMesh(mesh, path)
	return vBool(true), nil
}

func builtinUpdateMeshBuffer(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("updateMeshBuffer expects meshId, index, dataArr, offset")
	}
	_, mesh, err := meshByID("updateMeshBuffer", args, 0)
	if err != nil {
		return nil, err
	}
	index, _ := argInt("updateMeshBuffer", args, 1)
	if args[2] == nil || args[2].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("updateMeshBuffer: arg 3 must be byte array")
	}
	data := make([]byte, len(args[2].Elems))
	for i, e := range args[2].Elems {
		data[i] = byte(e.I64)
	}
	offset, _ := argInt("updateMeshBuffer", args, 3)
	rl.UpdateMeshBuffer(mesh, int(index), data, int(offset))
	return null(), nil
}

func builtinDrawMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawMesh expects meshId, materialId, [matrix]")
	}
	_, mesh, err := meshByID("drawMesh", args, 0)
	if err != nil {
		return nil, err
	}
	_, mat, err2 := materialByID("drawMesh", args, 1)
	if err2 != nil {
		return nil, err2
	}
	var transform rl.Matrix
	if len(args) > 2 && args[2] != nil && args[2].Kind == candy_evaluator.ValMap {
		transform, _ = argMatrix("drawMesh", args, 2)
	} else {
		transform = rl.MatrixIdentity()
	}
	rl.DrawMesh(mesh, mat, transform)
	return null(), nil
}

// ---- Mesh generation ----

func builtinGenMeshPoly(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshPoly", args, 2); err != nil {
		return nil, err
	}
	sides, _ := argInt("genMeshPoly", args, 0)
	r, _ := getArgFloat("genMeshPoly", args, 1)
	m := rl.GenMeshPoly(int(sides), float32(r))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshPlane(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshPlane", args, 4); err != nil {
		return nil, err
	}
	w, _ := getArgFloat("genMeshPlane", args, 0)
	l, _ := getArgFloat("genMeshPlane", args, 1)
	resX, _ := argInt("genMeshPlane", args, 2)
	resZ, _ := argInt("genMeshPlane", args, 3)
	m := rl.GenMeshPlane(float32(w), float32(l), int(resX), int(resZ))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshCube(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshCube", args, 3); err != nil {
		return nil, err
	}
	w, _ := getArgFloat("genMeshCube", args, 0)
	h, _ := getArgFloat("genMeshCube", args, 1)
	l, _ := getArgFloat("genMeshCube", args, 2)
	m := rl.GenMeshCube(float32(w), float32(h), float32(l))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshSphere", args, 3); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshSphere", args, 0)
	rings, _ := argInt("genMeshSphere", args, 1)
	slices, _ := argInt("genMeshSphere", args, 2)
	m := rl.GenMeshSphere(float32(r), int(rings), int(slices))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshHemiSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshHemiSphere", args, 3); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshHemiSphere", args, 0)
	rings, _ := argInt("genMeshHemiSphere", args, 1)
	slices, _ := argInt("genMeshHemiSphere", args, 2)
	m := rl.GenMeshHemiSphere(float32(r), int(rings), int(slices))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshCylinder(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshCylinder", args, 3); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshCylinder", args, 0)
	h, _ := getArgFloat("genMeshCylinder", args, 1)
	slices, _ := argInt("genMeshCylinder", args, 2)
	m := rl.GenMeshCylinder(float32(r), float32(h), int(slices))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshCone(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshCone", args, 3); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshCone", args, 0)
	h, _ := getArgFloat("genMeshCone", args, 1)
	slices, _ := argInt("genMeshCone", args, 2)
	m := rl.GenMeshCone(float32(r), float32(h), int(slices))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshTorus(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshTorus", args, 4); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshTorus", args, 0)
	size, _ := getArgFloat("genMeshTorus", args, 1)
	radSeg, _ := argInt("genMeshTorus", args, 2)
	sides, _ := argInt("genMeshTorus", args, 3)
	m := rl.GenMeshTorus(float32(r), float32(size), int(radSeg), int(sides))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshKnot(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshKnot", args, 4); err != nil {
		return nil, err
	}
	r, _ := getArgFloat("genMeshKnot", args, 0)
	size, _ := getArgFloat("genMeshKnot", args, 1)
	radSeg, _ := argInt("genMeshKnot", args, 2)
	sides, _ := argInt("genMeshKnot", args, 3)
	m := rl.GenMeshKnot(float32(r), float32(size), int(radSeg), int(sides))
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshHeightmap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshHeightmap", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("genMeshHeightmap", args, 0)
	if err != nil {
		return nil, err
	}
	size, err2 := argVector3("genMeshHeightmap", args, 1)
	if err2 != nil {
		return nil, err2
	}
	m := rl.GenMeshHeightmap(*img, size)
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

func builtinGenMeshCubicmap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("genMeshCubicmap", args, 2); err != nil {
		return nil, err
	}
	_, img, err := imageByID("genMeshCubicmap", args, 0)
	if err != nil {
		return nil, err
	}
	cubeSize, err2 := argVector3("genMeshCubicmap", args, 1)
	if err2 != nil {
		return nil, err2
	}
	m := rl.GenMeshCubicmap(*img, cubeSize)
	id := nextMeshID
	nextMeshID++
	meshes[id] = m
	return vInt(id), nil
}

// ---- Material management ----

func builtinLoadMaterials(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadMaterials", args, 1); err != nil {
		return nil, err
	}
	path, _ := argString("loadMaterials", args, 0)
	mats := rl.LoadMaterials(path)
	elems := make([]candy_evaluator.Value, len(mats))
	for i, mat := range mats {
		id := nextMaterialID
		nextMaterialID++
		materials[id] = mat
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValInt, I64: id}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

func builtinLoadMaterialDefault(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	mat := rl.LoadMaterialDefault()
	id := nextMaterialID
	nextMaterialID++
	materials[id] = mat
	return vInt(id), nil
}

func builtinIsMaterialValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMaterialValid", args, 1); err != nil {
		return nil, err
	}
	_, mat, err := materialByID("isMaterialValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsMaterialValid(mat)), nil
}

func builtinUnloadMaterial(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadMaterial", args, 1); err != nil {
		return nil, err
	}
	id, mat, err := materialByID("unloadMaterial", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadMaterial(mat)
	delete(materials, id)
	return null(), nil
}

func builtinSetMaterialTexture(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMaterialTexture", args, 3); err != nil {
		return nil, err
	}
	matID, mat, err := materialByID("setMaterialTexture", args, 0)
	if err != nil {
		return nil, err
	}
	mapType, _ := argInt("setMaterialTexture", args, 1)
	_, tex, err2 := textureByID("setMaterialTexture", args, 2)
	if err2 != nil {
		return nil, err2
	}
	rl.SetMaterialTexture(&mat, int32(mapType), tex)
	materials[matID] = mat
	return null(), nil
}

func builtinSetModelMeshMaterial(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setModelMeshMaterial", args, 3); err != nil {
		return nil, err
	}
	modelID, _ := argInt("setModelMeshMaterial", args, 0)
	m, ok := models[modelID]
	if !ok {
		return nil, fmt.Errorf("setModelMeshMaterial: invalid model %d", modelID)
	}
	meshIdx, _ := argInt("setModelMeshMaterial", args, 1)
	matIdx, _ := argInt("setModelMeshMaterial", args, 2)
	rl.SetModelMeshMaterial(&m, int32(meshIdx), int32(matIdx))
	models[modelID] = m
	return null(), nil
}

// ---- Animation extras ----

func builtinUpdateModelAnimationBones(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("updateModelAnimationBones expects modelId, animsId, frame")
	}
	mid, _ := argInt("updateModelAnimationBones", args, 0)
	aid, _ := argInt("updateModelAnimationBones", args, 1)
	frame, _ := argInt("updateModelAnimationBones", args, 2)
	m, okm := models[mid]
	a, oka := modelAnims[aid]
	if !okm || !oka || len(a) == 0 {
		return nil, fmt.Errorf("updateModelAnimationBones: invalid model %d or anim %d", mid, aid)
	}
	rl.UpdateModelAnimationBones(m, a[0], int32(frame))
	return null(), nil
}

func builtinUnloadModelAnimations(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadModelAnimations", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("unloadModelAnimations", args, 0)
	a, ok := modelAnims[id]
	if ok {
		rl.UnloadModelAnimations(a)
		delete(modelAnims, id)
	}
	return null(), nil
}

func builtinIsModelAnimationValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("isModelAnimationValid expects modelId, animsId")
	}
	mid, _ := argInt("isModelAnimationValid", args, 0)
	aid, _ := argInt("isModelAnimationValid", args, 1)
	m, okm := models[mid]
	a, oka := modelAnims[aid]
	if !okm || !oka || len(a) == 0 {
		return vBool(false), nil
	}
	return vBool(rl.IsModelAnimationValid(m, a[0])), nil
}

// ---- Collision extras ----

func builtinCheckCollisionBoxSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("checkCollisionBoxSphere expects boundingBox, cx, cy, cz, radius")
	}
	// Support map form: box{min,max}, cx, cy, cz, radius
	box, err := argBoundingBox("checkCollisionBoxSphere", args, 0)
	if err != nil {
		return nil, err
	}
	cx, _ := getArgFloat("checkCollisionBoxSphere", args, 1)
	cy, _ := getArgFloat("checkCollisionBoxSphere", args, 2)
	cz, _ := getArgFloat("checkCollisionBoxSphere", args, 3)
	r, _ := getArgFloat("checkCollisionBoxSphere", args, 4)
	return vBool(rl.CheckCollisionBoxSphere(box, rl.NewVector3(float32(cx), float32(cy), float32(cz)), float32(r))), nil
}

func builtinGetRayCollisionMesh(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("getRayCollisionMesh expects ray{position,direction}, meshId, [matrix]")
	}
	ray, err := argRay("getRayCollisionMesh", args, 0)
	if err != nil {
		return nil, err
	}
	_, mesh, err2 := meshByID("getRayCollisionMesh", args, 1)
	if err2 != nil {
		return nil, err2
	}
	var transform rl.Matrix
	if len(args) > 2 && args[2] != nil && args[2].Kind == candy_evaluator.ValMap {
		transform, _ = argMatrix("getRayCollisionMesh", args, 2)
	} else {
		transform = rl.MatrixIdentity()
	}
	return rayCollisionToMap(rl.GetRayCollisionMesh(ray, mesh, transform)), nil
}

func builtinGetRayCollisionTriangle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("getRayCollisionTriangle expects ray, p1{x,y,z}, p2{x,y,z}, p3{x,y,z}")
	}
	ray, err := argRay("getRayCollisionTriangle", args, 0)
	if err != nil {
		return nil, err
	}
	p1, _ := argVector3("getRayCollisionTriangle", args, 1)
	p2, _ := argVector3("getRayCollisionTriangle", args, 2)
	p3, _ := argVector3("getRayCollisionTriangle", args, 3)
	return rayCollisionToMap(rl.GetRayCollisionTriangle(ray, p1, p2, p3)), nil
}

func builtinGetRayCollisionQuad(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("getRayCollisionQuad expects ray, p1{x,y,z}, p2{x,y,z}, p3{x,y,z}, p4{x,y,z}")
	}
	ray, err := argRay("getRayCollisionQuad", args, 0)
	if err != nil {
		return nil, err
	}
	p1, _ := argVector3("getRayCollisionQuad", args, 1)
	p2, _ := argVector3("getRayCollisionQuad", args, 2)
	p3, _ := argVector3("getRayCollisionQuad", args, 3)
	p4, _ := argVector3("getRayCollisionQuad", args, 4)
	return rayCollisionToMap(rl.GetRayCollisionQuad(ray, p1, p2, p3, p4)), nil
}
