package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// builtinBeginMode3D sets the active 3D camera and starts 3D mode.
// Usage: beginMode3D(camPosX, camPosY, camPosZ, targetX, targetY, targetZ, fovy)
func builtinBeginMode3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("beginMode3D", args, 7); err != nil {
		return nil, err
	}
	cx, _ := getArgFloat("beginMode3D", args, 0)
	cy, _ := getArgFloat("beginMode3D", args, 1)
	cz, _ := getArgFloat("beginMode3D", args, 2)
	tx, _ := getArgFloat("beginMode3D", args, 3)
	ty, _ := getArgFloat("beginMode3D", args, 4)
	tz, _ := getArgFloat("beginMode3D", args, 5)
	fovy, _ := getArgFloat("beginMode3D", args, 6)
	activeCamera3D = rl.Camera3D{
		Position:   rl.NewVector3(float32(cx), float32(cy), float32(cz)),
		Target:     rl.NewVector3(float32(tx), float32(ty), float32(tz)),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       float32(fovy),
		Projection: rl.CameraPerspective,
	}
	rl.BeginMode3D(activeCamera3D)
	return null(), nil
}

// builtinEndMode3D ends 3D mode.
func builtinEndMode3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EndMode3D()
	return null(), nil
}

// builtinDrawCube draws a color-filled cube.
// Usage: drawCube(x, y, z, width, height, length, color)
func builtinDrawCube(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawCube expects x, y, z, width, height, depth, [color]")
	}
	x, _ := getArgFloat("drawCube", args, 0)
	y, _ := getArgFloat("drawCube", args, 1)
	z, _ := getArgFloat("drawCube", args, 2)
	w, _ := getArgFloat("drawCube", args, 3)
	h, _ := getArgFloat("drawCube", args, 4)
	d, _ := getArgFloat("drawCube", args, 5)
	c, _ := argColor("drawCube", args, 6, rl.Maroon)
	rl.DrawCube(rl.NewVector3(float32(x), float32(y), float32(z)), float32(w), float32(h), float32(d), c)
	return null(), nil
}

// builtinDrawCubeWires draws cube wireframe edges.
// Usage: drawCubeWires(x, y, z, width, height, depth, color)
func builtinDrawCubeWires(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawCubeWires expects x, y, z, width, height, depth, [color]")
	}
	x, _ := getArgFloat("drawCubeWires", args, 0)
	y, _ := getArgFloat("drawCubeWires", args, 1)
	z, _ := getArgFloat("drawCubeWires", args, 2)
	w, _ := getArgFloat("drawCubeWires", args, 3)
	h, _ := getArgFloat("drawCubeWires", args, 4)
	d, _ := getArgFloat("drawCubeWires", args, 5)
	c, _ := argColor("drawCubeWires", args, 6, rl.DarkGray)
	rl.DrawCubeWires(rl.NewVector3(float32(x), float32(y), float32(z)), float32(w), float32(h), float32(d), c)
	return null(), nil
}

// builtinDrawPlane draws a color-filled plane.
// Usage: drawPlane(x, y, z, width, length, color)
func builtinDrawPlane(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawPlane expects x, y, z, width, length, [color]")
	}
	x, _ := getArgFloat("drawPlane", args, 0)
	y, _ := getArgFloat("drawPlane", args, 1)
	z, _ := getArgFloat("drawPlane", args, 2)
	w, _ := getArgFloat("drawPlane", args, 3)
	l, _ := getArgFloat("drawPlane", args, 4)
	c, _ := argColor("drawPlane", args, 5, rl.Green)
	rl.DrawPlane(rl.NewVector3(float32(x), float32(y), float32(z)), rl.NewVector2(float32(w), float32(l)), c)
	return null(), nil
}

// builtinDrawSphere draws a color-filled sphere.
// Usage: drawSphere(x, y, z, radius, color)
func builtinDrawSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawSphere expects x, y, z, radius, [color]")
	}
	x, _ := getArgFloat("drawSphere", args, 0)
	y, _ := getArgFloat("drawSphere", args, 1)
	z, _ := getArgFloat("drawSphere", args, 2)
	r, _ := getArgFloat("drawSphere", args, 3)
	c, _ := argColor("drawSphere", args, 4, rl.RayWhite)
	rl.DrawSphere(rl.NewVector3(float32(x), float32(y), float32(z)), float32(r), c)
	return null(), nil
}

// builtinDrawSphereEx draws a sphere with custom rings/slices.
// Usage: drawSphereEx(x, y, z, radius, rings, slices, [color])
func builtinDrawSphereEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawSphereEx expects x, y, z, radius, rings, slices, [color]")
	}
	x, _ := getArgFloat("drawSphereEx", args, 0)
	y, _ := getArgFloat("drawSphereEx", args, 1)
	z, _ := getArgFloat("drawSphereEx", args, 2)
	r, _ := getArgFloat("drawSphereEx", args, 3)
	rings, _ := argInt("drawSphereEx", args, 4)
	slices, _ := argInt("drawSphereEx", args, 5)
	c, _ := argColor("drawSphereEx", args, 6, rl.RayWhite)
	rl.DrawSphereEx(rl.NewVector3(float32(x), float32(y), float32(z)), float32(r), int32(rings), int32(slices), c)
	return null(), nil
}

// builtinDrawSphereWires draws sphere wireframe with optional rings/slices.
// Usage: drawSphereWires(x, y, z, radius, [rings], [slices], [color])
func builtinDrawSphereWires(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("drawSphereWires expects x, y, z, radius, [rings], [slices], [color]")
	}
	x, _ := getArgFloat("drawSphereWires", args, 0)
	y, _ := getArgFloat("drawSphereWires", args, 1)
	z, _ := getArgFloat("drawSphereWires", args, 2)
	r, _ := getArgFloat("drawSphereWires", args, 3)
	rings, _ := argInt("drawSphereWires", args, 4)
	if len(args) <= 4 {
		rings = 8
	}
	slices, _ := argInt("drawSphereWires", args, 5)
	if len(args) <= 5 {
		slices = 8
	}
	c, _ := argColor("drawSphereWires", args, 6, rl.DarkGray)
	rl.DrawSphereWires(rl.NewVector3(float32(x), float32(y), float32(z)), float32(r), int32(rings), int32(slices), c)
	return null(), nil
}

// builtinDrawCylinder draws a color-filled cylinder.
// Usage: drawCylinder(x, y, z, radiusT, radiusB, height, color)
func builtinDrawCylinder(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawCylinder expects x, y, z, radiusT, radiusB, height, [color]")
	}
	x, _ := getArgFloat("drawCylinder", args, 0)
	y, _ := getArgFloat("drawCylinder", args, 1)
	z, _ := getArgFloat("drawCylinder", args, 2)
	rt, _ := getArgFloat("drawCylinder", args, 3)
	rb, _ := getArgFloat("drawCylinder", args, 4)
	h, _ := getArgFloat("drawCylinder", args, 5)
	c, _ := argColor("drawCylinder", args, 6, rl.RayWhite)
	rl.DrawCylinder(rl.NewVector3(float32(x), float32(y), float32(z)), float32(rt), float32(rb), float32(h), 16, c)
	return null(), nil
}

// builtinDrawCylinderEx draws a cylinder between two points.
// Usage: drawCylinderEx(x1, y1, z1, x2, y2, z2, startRadius, endRadius, [sides], [color])
func builtinDrawCylinderEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("drawCylinderEx expects x1, y1, z1, x2, y2, z2, startRadius, endRadius, [sides], [color]")
	}
	x1, _ := getArgFloat("drawCylinderEx", args, 0)
	y1, _ := getArgFloat("drawCylinderEx", args, 1)
	z1, _ := getArgFloat("drawCylinderEx", args, 2)
	x2, _ := getArgFloat("drawCylinderEx", args, 3)
	y2, _ := getArgFloat("drawCylinderEx", args, 4)
	z2, _ := getArgFloat("drawCylinderEx", args, 5)
	startR, _ := getArgFloat("drawCylinderEx", args, 6)
	endR, _ := getArgFloat("drawCylinderEx", args, 7)
	sides := int64(16)
	if len(args) > 8 {
		sides, _ = argInt("drawCylinderEx", args, 8)
	}
	c, _ := argColor("drawCylinderEx", args, 9, rl.RayWhite)
	rl.DrawCylinderEx(
		rl.NewVector3(float32(x1), float32(y1), float32(z1)),
		rl.NewVector3(float32(x2), float32(y2), float32(z2)),
		float32(startR),
		float32(endR),
		int32(sides),
		c,
	)
	return null(), nil
}

// builtinDrawCylinderWires draws a cylinder wireframe.
// Usage: drawCylinderWires(x, y, z, radiusT, radiusB, height, [slices], [color])
func builtinDrawCylinderWires(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawCylinderWires expects x, y, z, radiusT, radiusB, height, [slices], [color]")
	}
	x, _ := getArgFloat("drawCylinderWires", args, 0)
	y, _ := getArgFloat("drawCylinderWires", args, 1)
	z, _ := getArgFloat("drawCylinderWires", args, 2)
	rt, _ := getArgFloat("drawCylinderWires", args, 3)
	rb, _ := getArgFloat("drawCylinderWires", args, 4)
	h, _ := getArgFloat("drawCylinderWires", args, 5)
	slices := int64(16)
	if len(args) > 6 {
		slices, _ = argInt("drawCylinderWires", args, 6)
	}
	c, _ := argColor("drawCylinderWires", args, 7, rl.DarkGray)
	rl.DrawCylinderWires(rl.NewVector3(float32(x), float32(y), float32(z)), float32(rt), float32(rb), float32(h), int32(slices), c)
	return null(), nil
}

// builtinDrawLine3D draws a colored 3D line segment.
// Usage: drawLine3D(x1, y1, z1, x2, y2, z2, [color])
func builtinDrawLine3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawLine3D expects x1, y1, z1, x2, y2, z2, [color]")
	}
	x1, _ := getArgFloat("drawLine3D", args, 0)
	y1, _ := getArgFloat("drawLine3D", args, 1)
	z1, _ := getArgFloat("drawLine3D", args, 2)
	x2, _ := getArgFloat("drawLine3D", args, 3)
	y2, _ := getArgFloat("drawLine3D", args, 4)
	z2, _ := getArgFloat("drawLine3D", args, 5)
	c, _ := argColor("drawLine3D", args, 6, rl.White)
	rl.DrawLine3D(rl.NewVector3(float32(x1), float32(y1), float32(z1)), rl.NewVector3(float32(x2), float32(y2), float32(z2)), c)
	return null(), nil
}

// builtinDrawRay draws a ray from position in direction.
// Usage: drawRay(px, py, pz, dx, dy, dz, [color])
func builtinDrawRay(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawRay expects px, py, pz, dx, dy, dz, [color]")
	}
	px, _ := getArgFloat("drawRay", args, 0)
	py, _ := getArgFloat("drawRay", args, 1)
	pz, _ := getArgFloat("drawRay", args, 2)
	dx, _ := getArgFloat("drawRay", args, 3)
	dy, _ := getArgFloat("drawRay", args, 4)
	dz, _ := getArgFloat("drawRay", args, 5)
	c, _ := argColor("drawRay", args, 6, rl.Maroon)
	ray := rl.Ray{
		Position:  rl.NewVector3(float32(px), float32(py), float32(pz)),
		Direction: rl.NewVector3(float32(dx), float32(dy), float32(dz)),
	}
	rl.DrawRay(ray, c)
	return null(), nil
}

// builtinDrawBoundingBox draws an axis-aligned bounding box.
// Usage: drawBoundingBox(minX, minY, minZ, maxX, maxY, maxZ, [color])
func builtinDrawBoundingBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("drawBoundingBox expects minX, minY, minZ, maxX, maxY, maxZ, [color]")
	}
	minX, _ := getArgFloat("drawBoundingBox", args, 0)
	minY, _ := getArgFloat("drawBoundingBox", args, 1)
	minZ, _ := getArgFloat("drawBoundingBox", args, 2)
	maxX, _ := getArgFloat("drawBoundingBox", args, 3)
	maxY, _ := getArgFloat("drawBoundingBox", args, 4)
	maxZ, _ := getArgFloat("drawBoundingBox", args, 5)
	c, _ := argColor("drawBoundingBox", args, 6, rl.Green)
	box := rl.NewBoundingBox(
		rl.NewVector3(float32(minX), float32(minY), float32(minZ)),
		rl.NewVector3(float32(maxX), float32(maxY), float32(maxZ)),
	)
	rl.DrawBoundingBox(box, c)
	return null(), nil
}

// builtinDrawCapsule draws a filled capsule between two points.
// Usage: drawCapsule(x1, y1, z1, x2, y2, z2, radius, [slices], [rings], [color])
func builtinDrawCapsule(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("drawCapsule expects x1, y1, z1, x2, y2, z2, radius, [slices], [rings], [color]")
	}
	x1, _ := getArgFloat("drawCapsule", args, 0)
	y1, _ := getArgFloat("drawCapsule", args, 1)
	z1, _ := getArgFloat("drawCapsule", args, 2)
	x2, _ := getArgFloat("drawCapsule", args, 3)
	y2, _ := getArgFloat("drawCapsule", args, 4)
	z2, _ := getArgFloat("drawCapsule", args, 5)
	r, _ := getArgFloat("drawCapsule", args, 6)
	slices := int64(8)
	if len(args) > 7 {
		slices, _ = argInt("drawCapsule", args, 7)
	}
	rings := int64(8)
	if len(args) > 8 {
		rings, _ = argInt("drawCapsule", args, 8)
	}
	c, _ := argColor("drawCapsule", args, 9, rl.RayWhite)
	rl.DrawCapsule(
		rl.NewVector3(float32(x1), float32(y1), float32(z1)),
		rl.NewVector3(float32(x2), float32(y2), float32(z2)),
		float32(r),
		int32(slices),
		int32(rings),
		c,
	)
	return null(), nil
}

// builtinDrawCapsuleWires draws a capsule wireframe.
// Usage: drawCapsuleWires(x1, y1, z1, x2, y2, z2, radius, [slices], [rings], [color])
func builtinDrawCapsuleWires(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("drawCapsuleWires expects x1, y1, z1, x2, y2, z2, radius, [slices], [rings], [color]")
	}
	x1, _ := getArgFloat("drawCapsuleWires", args, 0)
	y1, _ := getArgFloat("drawCapsuleWires", args, 1)
	z1, _ := getArgFloat("drawCapsuleWires", args, 2)
	x2, _ := getArgFloat("drawCapsuleWires", args, 3)
	y2, _ := getArgFloat("drawCapsuleWires", args, 4)
	z2, _ := getArgFloat("drawCapsuleWires", args, 5)
	r, _ := getArgFloat("drawCapsuleWires", args, 6)
	slices := int64(8)
	if len(args) > 7 {
		slices, _ = argInt("drawCapsuleWires", args, 7)
	}
	rings := int64(8)
	if len(args) > 8 {
		rings, _ = argInt("drawCapsuleWires", args, 8)
	}
	c, _ := argColor("drawCapsuleWires", args, 9, rl.DarkGray)
	rl.DrawCapsuleWires(
		rl.NewVector3(float32(x1), float32(y1), float32(z1)),
		rl.NewVector3(float32(x2), float32(y2), float32(z2)),
		float32(r),
		int32(slices),
		int32(rings),
		c,
	)
	return null(), nil
}

// builtinDrawGrid draws a grid on the ground plane.
// Usage: drawGrid(slices, spacing)
func builtinDrawGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("drawGrid expects slices, spacing")
	}
	slices, _ := argInt("drawGrid", args, 0)
	spacing, _ := getArgFloat("drawGrid", args, 1)
	rl.DrawGrid(int32(slices), float32(spacing))
	return null(), nil
}

// builtinGetMouseRay gets a ray from the current mouse position.
func builtinGetMouseRay(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	pos := rl.GetMousePosition()
	ray := rl.GetMouseRay(pos, activeCamera3D)

	m := make(map[string]candy_evaluator.Value)
	m["position"] = candy_evaluator.Value{
		Kind: candy_evaluator.ValMap,
		StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Position.Z)},
		},
	}
	m["direction"] = candy_evaluator.Value{
		Kind: candy_evaluator.ValMap,
		StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(ray.Direction.Z)},
		},
	}
	return vMap(m), nil
}

// builtinGetRayCollisionBox calculates collision between a ray and a bounding box.
func builtinGetRayCollisionBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getRayCollisionBox", args, 12); err != nil {
		return nil, err
	}
	rpx, _ := getArgFloat("getRayCollisionBox", args, 0)
	rpy, _ := getArgFloat("getRayCollisionBox", args, 1)
	rpz, _ := getArgFloat("getRayCollisionBox", args, 2)
	rdx, _ := getArgFloat("getRayCollisionBox", args, 3)
	rdy, _ := getArgFloat("getRayCollisionBox", args, 4)
	rdz, _ := getArgFloat("getRayCollisionBox", args, 5)
	minx, _ := getArgFloat("getRayCollisionBox", args, 6)
	miny, _ := getArgFloat("getRayCollisionBox", args, 7)
	minz, _ := getArgFloat("getRayCollisionBox", args, 8)
	maxx, _ := getArgFloat("getRayCollisionBox", args, 9)
	maxy, _ := getArgFloat("getRayCollisionBox", args, 10)
	maxz, _ := getArgFloat("getRayCollisionBox", args, 11)

	ray := rl.Ray{
		Position:  rl.NewVector3(float32(rpx), float32(rpy), float32(rpz)),
		Direction: rl.NewVector3(float32(rdx), float32(rdy), float32(rdz)),
	}
	box := rl.NewBoundingBox(rl.NewVector3(float32(minx), float32(miny), float32(minz)), rl.NewVector3(float32(maxx), float32(maxy), float32(maxz)))

	res := rl.GetRayCollisionBox(ray, box)
	m := make(map[string]candy_evaluator.Value)
	m["hit"] = candy_evaluator.Value{Kind: candy_evaluator.ValBool, B: res.Hit}
	m["distance"] = candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: float64(res.Distance)}
	m["point"] = candy_evaluator.Value{
		Kind: candy_evaluator.ValMap,
		StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.Z)},
		},
	}
	return vMap(m), nil
}

// builtinGetRayCollisionSphere calculates collision between a ray and a sphere.
func builtinGetRayCollisionSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("getRayCollisionSphere", args, 10); err != nil {
		return nil, err
	}
	rpx, _ := getArgFloat("getRayCollisionSphere", args, 0)
	rpy, _ := getArgFloat("getRayCollisionSphere", args, 1)
	rpz, _ := getArgFloat("getRayCollisionSphere", args, 2)
	rdx, _ := getArgFloat("getRayCollisionSphere", args, 3)
	rdy, _ := getArgFloat("getRayCollisionSphere", args, 4)
	rdz, _ := getArgFloat("getRayCollisionSphere", args, 5)
	cx, _ := getArgFloat("getRayCollisionSphere", args, 6)
	cy, _ := getArgFloat("getRayCollisionSphere", args, 7)
	cz, _ := getArgFloat("getRayCollisionSphere", args, 8)
	radius, _ := getArgFloat("getRayCollisionSphere", args, 9)

	ray := rl.Ray{
		Position:  rl.NewVector3(float32(rpx), float32(rpy), float32(rpz)),
		Direction: rl.NewVector3(float32(rdx), float32(rdy), float32(rdz)),
	}

	res := rl.GetRayCollisionSphere(ray, rl.NewVector3(float32(cx), float32(cy), float32(cz)), float32(radius))
	m := make(map[string]candy_evaluator.Value)
	m["hit"] = candy_evaluator.Value{Kind: candy_evaluator.ValBool, B: res.Hit}
	m["distance"] = candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: float64(res.Distance)}
	m["point"] = candy_evaluator.Value{
		Kind: candy_evaluator.ValMap,
		StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.X)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.Y)},
			"z": {Kind: candy_evaluator.ValFloat, F64: float64(res.Point.Z)},
		},
	}
	return vMap(m), nil
}

// builtinCheckCollisionBoxes checks collision between two bounding boxes.
func builtinCheckCollisionBoxes(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionBoxes", args, 12); err != nil {
		return nil, err
	}
	min1x, _ := getArgFloat("checkCollisionBoxes", args, 0)
	min1y, _ := getArgFloat("checkCollisionBoxes", args, 1)
	min1z, _ := getArgFloat("checkCollisionBoxes", args, 2)
	max1x, _ := getArgFloat("checkCollisionBoxes", args, 3)
	max1y, _ := getArgFloat("checkCollisionBoxes", args, 4)
	max1z, _ := getArgFloat("checkCollisionBoxes", args, 5)
	min2x, _ := getArgFloat("checkCollisionBoxes", args, 6)
	min2y, _ := getArgFloat("checkCollisionBoxes", args, 7)
	min2z, _ := getArgFloat("checkCollisionBoxes", args, 8)
	max2x, _ := getArgFloat("checkCollisionBoxes", args, 9)
	max2y, _ := getArgFloat("checkCollisionBoxes", args, 10)
	max2z, _ := getArgFloat("checkCollisionBoxes", args, 11)

	res := rl.CheckCollisionBoxes(
		rl.NewBoundingBox(rl.NewVector3(float32(min1x), float32(min1y), float32(min1z)), rl.NewVector3(float32(max1x), float32(max1y), float32(max1z))),
		rl.NewBoundingBox(rl.NewVector3(float32(min2x), float32(min2y), float32(min2z)), rl.NewVector3(float32(max2x), float32(max2y), float32(max2z))),
	)
	return vBool(res), nil
}

// builtinCheckCollisionSpheres checks collision between two spheres.
func builtinCheckCollisionSpheres(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("checkCollisionSpheres", args, 8); err != nil {
		return nil, err
	}
	cx1, _ := getArgFloat("checkCollisionSpheres", args, 0)
	cy1, _ := getArgFloat("checkCollisionSpheres", args, 1)
	cz1, _ := getArgFloat("checkCollisionSpheres", args, 2)
	radius1, _ := getArgFloat("checkCollisionSpheres", args, 3)
	cx2, _ := getArgFloat("checkCollisionSpheres", args, 4)
	cy2, _ := getArgFloat("checkCollisionSpheres", args, 5)
	cz2, _ := getArgFloat("checkCollisionSpheres", args, 6)
	radius2, _ := getArgFloat("checkCollisionSpheres", args, 7)

	res := rl.CheckCollisionSpheres(
		rl.NewVector3(float32(cx1), float32(cy1), float32(cz1)), float32(radius1),
		rl.NewVector3(float32(cx2), float32(cy2), float32(cz2)), float32(radius2),
	)
	return vBool(res), nil
}

// ---- External Models ----

// builtinLoadModel loads a 3D model from a file.
func builtinLoadModel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadModel", args, 1); err != nil {
		return nil, err
	}
	path, err := argString("loadModel", args, 0)
	if err != nil {
		return nil, err
	}

	m := rl.LoadModel(path)
	id := nextModelID
	models[id] = m
	nextModelID++
	return vInt(id), nil
}

// builtinUnloadModel unloads a 3D model from memory.
func builtinUnloadModel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadModel", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("unloadModel", args, 0)
	m, ok := models[id]
	if ok {
		rl.UnloadModel(m)
		delete(models, id)
	}
	return null(), nil
}

// builtinDrawModel draws a 3D model.
func builtinDrawModel(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawModel expects modelId, x, y, z, scale, [color]")
	}
	id, _ := argInt("drawModel", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModel: invalid model %d", id)
	}

	x, _ := getArgFloat("drawModel", args, 1)
	y, _ := getArgFloat("drawModel", args, 2)
	z, _ := getArgFloat("drawModel", args, 3)
	scale, _ := getArgFloat("drawModel", args, 4)
	c, _ := argColor("drawModel", args, 5, rl.White)

	rl.DrawModel(m, rl.NewVector3(float32(x), float32(y), float32(z)), float32(scale), c)
	return null(), nil
}

// builtinDrawModelEx draws a 3D model with extended parameters (rotation, scale).
func builtinDrawModelEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 11 {
		return nil, fmt.Errorf("drawModelEx expects modelId, x, y, z, rotX, rotY, rotZ, angle, scaleX, scaleY, scaleZ, [color]")
	}
	id, _ := argInt("drawModelEx", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModelEx: invalid model %d", id)
	}

	px, _ := getArgFloat("drawModelEx", args, 1)
	py, _ := getArgFloat("drawModelEx", args, 2)
	pz, _ := getArgFloat("drawModelEx", args, 3)
	rx, _ := getArgFloat("drawModelEx", args, 4)
	ry, _ := getArgFloat("drawModelEx", args, 5)
	rz, _ := getArgFloat("drawModelEx", args, 6)
	angle, _ := getArgFloat("drawModelEx", args, 7)
	sx, _ := getArgFloat("drawModelEx", args, 8)
	sy, _ := getArgFloat("drawModelEx", args, 9)
	sz, _ := getArgFloat("drawModelEx", args, 10)
	c, _ := argColor("drawModelEx", args, 11, rl.White)

	rl.DrawModelEx(m, rl.NewVector3(float32(px), float32(py), float32(pz)), rl.NewVector3(float32(rx), float32(ry), float32(rz)), float32(angle), rl.NewVector3(float32(sx), float32(sy), float32(sz)), c)
	return null(), nil
}

// builtinDrawModelWires draws a model in wireframe.
func builtinDrawModelWires(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawModelWires expects modelId, x, y, z, scale, [color]")
	}
	id, _ := argInt("drawModelWires", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModelWires: invalid model %d", id)
	}
	x, _ := getArgFloat("drawModelWires", args, 1)
	y, _ := getArgFloat("drawModelWires", args, 2)
	z, _ := getArgFloat("drawModelWires", args, 3)
	scale, _ := getArgFloat("drawModelWires", args, 4)
	c, _ := argColor("drawModelWires", args, 5, rl.White)
	rl.DrawModelWires(m, rl.NewVector3(float32(x), float32(y), float32(z)), float32(scale), c)
	return null(), nil
}

// builtinDrawModelWiresEx draws a wireframe model with full transform.
func builtinDrawModelWiresEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 11 {
		return nil, fmt.Errorf("drawModelWiresEx expects modelId, x, y, z, rotX, rotY, rotZ, angle, scaleX, scaleY, scaleZ, [color]")
	}
	id, _ := argInt("drawModelWiresEx", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModelWiresEx: invalid model %d", id)
	}
	px, _ := getArgFloat("drawModelWiresEx", args, 1)
	py, _ := getArgFloat("drawModelWiresEx", args, 2)
	pz, _ := getArgFloat("drawModelWiresEx", args, 3)
	rx, _ := getArgFloat("drawModelWiresEx", args, 4)
	ry, _ := getArgFloat("drawModelWiresEx", args, 5)
	rz, _ := getArgFloat("drawModelWiresEx", args, 6)
	angle, _ := getArgFloat("drawModelWiresEx", args, 7)
	sx, _ := getArgFloat("drawModelWiresEx", args, 8)
	sy, _ := getArgFloat("drawModelWiresEx", args, 9)
	sz, _ := getArgFloat("drawModelWiresEx", args, 10)
	c, _ := argColor("drawModelWiresEx", args, 11, rl.White)
	rl.DrawModelWiresEx(m, rl.NewVector3(float32(px), float32(py), float32(pz)), rl.NewVector3(float32(rx), float32(ry), float32(rz)), float32(angle), rl.NewVector3(float32(sx), float32(sy), float32(sz)), c)
	return null(), nil
}

// builtinDrawModelPoints draws model vertices as points.
func builtinDrawModelPoints(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawModelPoints expects modelId, x, y, z, scale, [color]")
	}
	id, _ := argInt("drawModelPoints", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModelPoints: invalid model %d", id)
	}
	x, _ := getArgFloat("drawModelPoints", args, 1)
	y, _ := getArgFloat("drawModelPoints", args, 2)
	z, _ := getArgFloat("drawModelPoints", args, 3)
	scale, _ := getArgFloat("drawModelPoints", args, 4)
	c, _ := argColor("drawModelPoints", args, 5, rl.White)
	rl.DrawModelPoints(m, rl.NewVector3(float32(x), float32(y), float32(z)), float32(scale), c)
	return null(), nil
}

// builtinDrawModelPointsEx draws model vertices with full transform.
func builtinDrawModelPointsEx(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 11 {
		return nil, fmt.Errorf("drawModelPointsEx expects modelId, x, y, z, rotX, rotY, rotZ, angle, scaleX, scaleY, scaleZ, [color]")
	}
	id, _ := argInt("drawModelPointsEx", args, 0)
	m, ok := models[id]
	if !ok {
		return nil, fmt.Errorf("drawModelPointsEx: invalid model %d", id)
	}
	px, _ := getArgFloat("drawModelPointsEx", args, 1)
	py, _ := getArgFloat("drawModelPointsEx", args, 2)
	pz, _ := getArgFloat("drawModelPointsEx", args, 3)
	rx, _ := getArgFloat("drawModelPointsEx", args, 4)
	ry, _ := getArgFloat("drawModelPointsEx", args, 5)
	rz, _ := getArgFloat("drawModelPointsEx", args, 6)
	angle, _ := getArgFloat("drawModelPointsEx", args, 7)
	sx, _ := getArgFloat("drawModelPointsEx", args, 8)
	sy, _ := getArgFloat("drawModelPointsEx", args, 9)
	sz, _ := getArgFloat("drawModelPointsEx", args, 10)
	c, _ := argColor("drawModelPointsEx", args, 11, rl.White)
	rl.DrawModelPointsEx(m, rl.NewVector3(float32(px), float32(py), float32(pz)), rl.NewVector3(float32(rx), float32(ry), float32(rz)), float32(angle), rl.NewVector3(float32(sx), float32(sy), float32(sz)), c)
	return null(), nil
}

// builtinDrawBillboard draws a billboard texture facing the active camera.
func builtinDrawBillboard(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("drawBillboard expects textureId, x, y, z, scale, [color]")
	}
	_, tex, err := textureByID("drawBillboard", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := getArgFloat("drawBillboard", args, 1)
	y, _ := getArgFloat("drawBillboard", args, 2)
	z, _ := getArgFloat("drawBillboard", args, 3)
	scale, _ := getArgFloat("drawBillboard", args, 4)
	c, _ := argColor("drawBillboard", args, 5, rl.White)
	rl.DrawBillboard(activeCamera3D, tex, rl.NewVector3(float32(x), float32(y), float32(z)), float32(scale), c)
	return null(), nil
}

// builtinDrawBillboardRec draws a billboard with source rectangle and size.
func builtinDrawBillboardRec(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 10 {
		return nil, fmt.Errorf("drawBillboardRec expects textureId, sx, sy, sw, sh, x, y, z, sizeX, sizeY, [color]")
	}
	_, tex, err := textureByID("drawBillboardRec", args, 0)
	if err != nil {
		return nil, err
	}
	sx, _ := getArgFloat("drawBillboardRec", args, 1)
	sy, _ := getArgFloat("drawBillboardRec", args, 2)
	sw, _ := getArgFloat("drawBillboardRec", args, 3)
	sh, _ := getArgFloat("drawBillboardRec", args, 4)
	x, _ := getArgFloat("drawBillboardRec", args, 5)
	y, _ := getArgFloat("drawBillboardRec", args, 6)
	z, _ := getArgFloat("drawBillboardRec", args, 7)
	sizeX, _ := getArgFloat("drawBillboardRec", args, 8)
	sizeY, _ := getArgFloat("drawBillboardRec", args, 9)
	c, _ := argColor("drawBillboardRec", args, 10, rl.White)
	rl.DrawBillboardRec(
		activeCamera3D,
		tex,
		rl.NewRectangle(float32(sx), float32(sy), float32(sw), float32(sh)),
		rl.NewVector3(float32(x), float32(y), float32(z)),
		rl.NewVector2(float32(sizeX), float32(sizeY)),
		c,
	)
	return null(), nil
}

// builtinDrawBillboardPro draws a billboard with full control over up, size, origin and rotation.
func builtinDrawBillboardPro(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 17 {
		return nil, fmt.Errorf("drawBillboardPro expects textureId, sx, sy, sw, sh, px, py, pz, upX, upY, upZ, sizeX, sizeY, originX, originY, rotation, [color]")
	}
	_, tex, err := textureByID("drawBillboardPro", args, 0)
	if err != nil {
		return nil, err
	}
	sx, _ := getArgFloat("drawBillboardPro", args, 1)
	sy, _ := getArgFloat("drawBillboardPro", args, 2)
	sw, _ := getArgFloat("drawBillboardPro", args, 3)
	sh, _ := getArgFloat("drawBillboardPro", args, 4)
	px, _ := getArgFloat("drawBillboardPro", args, 5)
	py, _ := getArgFloat("drawBillboardPro", args, 6)
	pz, _ := getArgFloat("drawBillboardPro", args, 7)
	upX, _ := getArgFloat("drawBillboardPro", args, 8)
	upY, _ := getArgFloat("drawBillboardPro", args, 9)
	upZ, _ := getArgFloat("drawBillboardPro", args, 10)
	sizeX, _ := getArgFloat("drawBillboardPro", args, 11)
	sizeY, _ := getArgFloat("drawBillboardPro", args, 12)
	originX, _ := getArgFloat("drawBillboardPro", args, 13)
	originY, _ := getArgFloat("drawBillboardPro", args, 14)
	rotation, _ := getArgFloat("drawBillboardPro", args, 15)
	c, _ := argColor("drawBillboardPro", args, 16, rl.White)
	rl.DrawBillboardPro(
		activeCamera3D,
		tex,
		rl.NewRectangle(float32(sx), float32(sy), float32(sw), float32(sh)),
		rl.NewVector3(float32(px), float32(py), float32(pz)),
		rl.NewVector3(float32(upX), float32(upY), float32(upZ)),
		rl.NewVector2(float32(sizeX), float32(sizeY)),
		rl.NewVector2(float32(originX), float32(originY)),
		float32(rotation),
		c,
	)
	return null(), nil
}

// builtinLoadModelAnimations loads model animations from a file.
func builtinLoadModelAnimations(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadModelAnimations", args, 1); err != nil {
		return nil, err
	}
	path, err := argString("loadModelAnimations", args, 0)
	if err != nil {
		return nil, err
	}

	anims := rl.LoadModelAnimations(path)
	id := nextModelAnimID
	modelAnims[id] = anims
	nextModelAnimID++
	return vInt(id), nil
}

// builtinUpdateModelAnimation updates model animation pose.
func builtinUpdateModelAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("updateModelAnimation expects modelId, animsId, frame")
	}
	mid, _ := argInt("updateModelAnimation", args, 0)
	aid, _ := argInt("updateModelAnimation", args, 1)
	frame, _ := argInt("updateModelAnimation", args, 2)

	m, okm := models[mid]
	a, oka := modelAnims[aid]
	if !okm || !oka {
		return nil, fmt.Errorf("updateModelAnimation: invalid model %d or anim %d", mid, aid)
	}

	rl.UpdateModelAnimation(m, a[0], int32(frame))
	return null(), nil
}
