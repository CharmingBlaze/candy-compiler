package candy_physics

import (
	"fmt"
	"math"

	"candy/candy_evaluator"
)

// ---- helpers ------------------------------------------------------------

func null() *candy_evaluator.Value { return &candy_evaluator.Value{Kind: candy_evaluator.ValNull} }

func vBool(b bool) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValBool, B: b}
}
func vInt(n int64) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValInt, I64: n}
}
func vFloat(f float64) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: f}
}
func vMap(m map[string]candy_evaluator.Value) *candy_evaluator.Value {
	return &candy_evaluator.Value{Kind: candy_evaluator.ValMap, StrMap: m}
}

func expectArgs(name string, args []*candy_evaluator.Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s expects %d args, got %d", name, n, len(args))
	}
	return nil
}

func argFloat(name string, args []*candy_evaluator.Value, i int) (float64, error) {
	if i >= len(args) || args[i] == nil {
		return 0, fmt.Errorf("%s arg %d must be a number", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValFloat {
		return v.F64, nil
	}
	if v.Kind == candy_evaluator.ValInt {
		return float64(v.I64), nil
	}
	return 0, fmt.Errorf("%s arg %d must be a number", name, i+1)
}

func argInt(name string, args []*candy_evaluator.Value, i int) (int64, error) {
	f, err := argFloat(name, args, i)
	return int64(f), err
}

func argVec3(name string, args []*candy_evaluator.Value, i int) (Vec3, error) {
	if i >= len(args) || args[i] == nil {
		return Vec3{}, fmt.Errorf("%s arg %d must be a {x,y,z} map", name, i+1)
	}
	v := args[i]
	if v.Kind == candy_evaluator.ValMap {
		x, _ := mapFloat(v, "x")
		y, _ := mapFloat(v, "y")
		z, _ := mapFloat(v, "z")
		return Vec3{x, y, z}, nil
	}
	return Vec3{}, fmt.Errorf("%s arg %d must be a {x,y,z} map", name, i+1)
}

func mapFloat(v *candy_evaluator.Value, key string) (float64, bool) {
	if v.Kind != candy_evaluator.ValMap {
		return 0, false
	}
	f, ok := v.StrMap[key]
	if !ok {
		return 0, false
	}
	if f.Kind == candy_evaluator.ValFloat {
		return f.F64, true
	}
	if f.Kind == candy_evaluator.ValInt {
		return float64(f.I64), true
	}
	return 0, false
}

func vec3ToMap(v Vec3) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: v.X},
		"y": {Kind: candy_evaluator.ValFloat, F64: v.Y},
		"z": {Kind: candy_evaluator.ValFloat, F64: v.Z},
	})
}

func quatToMap(q Quat) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: q.X},
		"y": {Kind: candy_evaluator.ValFloat, F64: q.Y},
		"z": {Kind: candy_evaluator.ValFloat, F64: q.Z},
		"w": {Kind: candy_evaluator.ValFloat, F64: q.W},
	})
}

func requireWorld(name string, args []*candy_evaluator.Value, i int) (*World, int64, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return nil, 0, err
	}
	w, ok := GetWorld(id)
	if !ok {
		return nil, 0, fmt.Errorf("%s: invalid world handle %d", name, id)
	}
	return w, id, nil
}

func requireBody(name string, w *World, args []*candy_evaluator.Value, i int) (*Body, int64, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return nil, 0, err
	}
	b, ok := w.GetBody(id)
	if !ok {
		return nil, 0, fmt.Errorf("%s: invalid body handle %d", name, id)
	}
	return b, id, nil
}

// ---- World management ---------------------------------------------------

// physicsCreateWorld([gravX, gravY, gravZ]) → worldId
func builtinPhysicsCreateWorld(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := NewWorld()
	w, _ := GetWorld(id)
	if len(args) >= 3 {
		gx, _ := argFloat("physicsCreateWorld", args, 0)
		gy, _ := argFloat("physicsCreateWorld", args, 1)
		gz, _ := argFloat("physicsCreateWorld", args, 2)
		w.SetGravity(Vec3{gx, gy, gz})
	}
	return vInt(id), nil
}

// physicsDestroyWorld(worldId)
func builtinPhysicsDestroyWorld(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsDestroyWorld", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("physicsDestroyWorld", args, 0)
	DestroyWorld(id)
	return null(), nil
}

// physicsSetGravity(worldId, x, y, z)
func builtinPhysicsSetGravity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("physicsSetGravity expects worldId, x, y, z")
	}
	w, _, err := requireWorld("physicsSetGravity", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsSetGravity", args, 1)
	y, _ := argFloat("physicsSetGravity", args, 2)
	z, _ := argFloat("physicsSetGravity", args, 3)
	w.SetGravity(Vec3{x, y, z})
	return null(), nil
}

// physicsGetGravity(worldId) → {x,y,z}
func builtinPhysicsGetGravity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetGravity", args, 1); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetGravity", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(w.GetGravity()), nil
}

// physicsStep(worldId, deltaTime)
func builtinPhysicsStep(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsStep", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsStep", args, 0)
	if err != nil {
		return nil, err
	}
	dt, _ := argFloat("physicsStep", args, 1)
	w.Step(dt)
	return null(), nil
}

// ---- Body creation ------------------------------------------------------

// physicsCreateBox(worldId, x, y, z, halfW, halfH, halfD, [motionType], [isSensor]) → bodyId
func builtinPhysicsCreateBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("physicsCreateBox expects worldId, x, y, z, halfW, halfH, halfD, [motion], [sensor]")
	}
	w, _, err := requireWorld("physicsCreateBox", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsCreateBox", args, 1)
	y, _ := argFloat("physicsCreateBox", args, 2)
	z, _ := argFloat("physicsCreateBox", args, 3)
	hw, _ := argFloat("physicsCreateBox", args, 4)
	hh, _ := argFloat("physicsCreateBox", args, 5)
	hd, _ := argFloat("physicsCreateBox", args, 6)
	motion := MotionDynamic
	if len(args) > 7 {
		m, _ := argInt("physicsCreateBox", args, 7)
		motion = MotionType(m)
	}
	sensor := false
	if len(args) > 8 {
		if args[8] != nil && args[8].Kind == candy_evaluator.ValBool {
			sensor = args[8].B
		}
	}
	id := w.CreateBody(NewBoxShape(hw, hh, hd), Vec3{x, y, z}, motion, sensor)
	return vInt(id), nil
}

// physicsCreateSphere(worldId, x, y, z, radius, [motionType], [isSensor]) → bodyId
func builtinPhysicsCreateSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsCreateSphere expects worldId, x, y, z, radius, [motion], [sensor]")
	}
	w, _, err := requireWorld("physicsCreateSphere", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsCreateSphere", args, 1)
	y, _ := argFloat("physicsCreateSphere", args, 2)
	z, _ := argFloat("physicsCreateSphere", args, 3)
	r, _ := argFloat("physicsCreateSphere", args, 4)
	motion := MotionDynamic
	if len(args) > 5 {
		m, _ := argInt("physicsCreateSphere", args, 5)
		motion = MotionType(m)
	}
	sensor := false
	if len(args) > 6 {
		if args[6] != nil && args[6].Kind == candy_evaluator.ValBool {
			sensor = args[6].B
		}
	}
	id := w.CreateBody(NewSphereShape(r), Vec3{x, y, z}, motion, sensor)
	return vInt(id), nil
}

// physicsCreateCapsule(worldId, x, y, z, radius, halfHeight, [motionType], [isSensor]) → bodyId
func builtinPhysicsCreateCapsule(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("physicsCreateCapsule expects worldId, x, y, z, radius, halfHeight, [motion], [sensor]")
	}
	w, _, err := requireWorld("physicsCreateCapsule", args, 0)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsCreateCapsule", args, 1)
	y, _ := argFloat("physicsCreateCapsule", args, 2)
	z, _ := argFloat("physicsCreateCapsule", args, 3)
	r, _ := argFloat("physicsCreateCapsule", args, 4)
	hh, _ := argFloat("physicsCreateCapsule", args, 5)
	motion := MotionDynamic
	if len(args) > 6 {
		m, _ := argInt("physicsCreateCapsule", args, 6)
		motion = MotionType(m)
	}
	sensor := false
	if len(args) > 7 {
		if args[7] != nil && args[7].Kind == candy_evaluator.ValBool {
			sensor = args[7].B
		}
	}
	id := w.CreateBody(NewCapsuleShape(r, hh), Vec3{x, y, z}, motion, sensor)
	return vInt(id), nil
}

// physicsCreatePlane(worldId, [yOffset]) → bodyId
func builtinPhysicsCreatePlane(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("physicsCreatePlane expects worldId, [yOffset]")
	}
	w, _, err := requireWorld("physicsCreatePlane", args, 0)
	if err != nil {
		return nil, err
	}
	yOff := 0.0
	if len(args) > 1 {
		yOff, _ = argFloat("physicsCreatePlane", args, 1)
	}
	id := w.CreateBody(NewPlaneShape(yOff), Vec3{}, MotionStatic, false)
	return vInt(id), nil
}

// physicsDestroyBody(worldId, bodyId)
func builtinPhysicsDestroyBody(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsDestroyBody", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsDestroyBody", args, 0)
	if err != nil {
		return nil, err
	}
	bid, _ := argInt("physicsDestroyBody", args, 1)
	w.DestroyBody(bid)
	return null(), nil
}

// ---- Body property get/set ----------------------------------------------

// physicsGetPosition(worldId, bodyId) → {x,y,z}
func builtinPhysicsGetPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetPosition", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetPosition", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetPosition", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(b.Position), nil
}

// physicsSetPosition(worldId, bodyId, x, y, z)
func builtinPhysicsSetPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsSetPosition expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsSetPosition", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetPosition", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsSetPosition", args, 2)
	y, _ := argFloat("physicsSetPosition", args, 3)
	z, _ := argFloat("physicsSetPosition", args, 4)
	b.Position = Vec3{x, y, z}
	return null(), nil
}

// physicsGetVelocity(worldId, bodyId) → {x,y,z}
func builtinPhysicsGetVelocity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetVelocity", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetVelocity", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetVelocity", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(b.LinearVel), nil
}

// physicsSetVelocity(worldId, bodyId, x, y, z)
func builtinPhysicsSetVelocity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsSetVelocity expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsSetVelocity", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetVelocity", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsSetVelocity", args, 2)
	y, _ := argFloat("physicsSetVelocity", args, 3)
	z, _ := argFloat("physicsSetVelocity", args, 4)
	b.LinearVel = Vec3{x, y, z}
	return null(), nil
}

// physicsGetAngularVelocity(worldId, bodyId) → {x,y,z}
func builtinPhysicsGetAngularVelocity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetAngularVelocity", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetAngularVelocity", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetAngularVelocity", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(b.AngularVel), nil
}

// physicsSetAngularVelocity(worldId, bodyId, x, y, z)
func builtinPhysicsSetAngularVelocity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsSetAngularVelocity expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsSetAngularVelocity", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetAngularVelocity", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsSetAngularVelocity", args, 2)
	y, _ := argFloat("physicsSetAngularVelocity", args, 3)
	z, _ := argFloat("physicsSetAngularVelocity", args, 4)
	b.AngularVel = Vec3{x, y, z}
	return null(), nil
}

// physicsGetRotation(worldId, bodyId) → {x,y,z,w}
func builtinPhysicsGetRotation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetRotation", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetRotation", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetRotation", w, args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(b.Rotation), nil
}

// physicsSetRotation(worldId, bodyId, x, y, z, w)
func builtinPhysicsSetRotation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("physicsSetRotation expects worldId, bodyId, x, y, z, w")
	}
	world, _, err := requireWorld("physicsSetRotation", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetRotation", world, args, 1)
	if err != nil {
		return nil, err
	}
	qx, _ := argFloat("physicsSetRotation", args, 2)
	qy, _ := argFloat("physicsSetRotation", args, 3)
	qz, _ := argFloat("physicsSetRotation", args, 4)
	qw, _ := argFloat("physicsSetRotation", args, 5)
	b.Rotation = Quat{qx, qy, qz, qw}.Normalize()
	return null(), nil
}

// physicsSetRotationFromAxisAngle(worldId, bodyId, axisX, axisY, axisZ, angleDegrees)
func builtinPhysicsSetRotationFromAxisAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("physicsSetRotationFromAxisAngle expects worldId, bodyId, axisX, axisY, axisZ, angleDegrees")
	}
	world, _, err := requireWorld("physicsSetRotationFromAxisAngle", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetRotationFromAxisAngle", world, args, 1)
	if err != nil {
		return nil, err
	}
	ax, _ := argFloat("physicsSetRotationFromAxisAngle", args, 2)
	ay, _ := argFloat("physicsSetRotationFromAxisAngle", args, 3)
	az, _ := argFloat("physicsSetRotationFromAxisAngle", args, 4)
	deg, _ := argFloat("physicsSetRotationFromAxisAngle", args, 5)
	b.Rotation = quatFromAxisAngle(Vec3{ax, ay, az}, deg*math.Pi/180.0)
	return null(), nil
}

// physicsApplyForce(worldId, bodyId, x, y, z)
func builtinPhysicsApplyForce(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsApplyForce expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsApplyForce", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsApplyForce", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsApplyForce", args, 2)
	y, _ := argFloat("physicsApplyForce", args, 3)
	z, _ := argFloat("physicsApplyForce", args, 4)
	b.ApplyForce(Vec3{x, y, z})
	return null(), nil
}

// physicsApplyImpulse(worldId, bodyId, x, y, z)
func builtinPhysicsApplyImpulse(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsApplyImpulse expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsApplyImpulse", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsApplyImpulse", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsApplyImpulse", args, 2)
	y, _ := argFloat("physicsApplyImpulse", args, 3)
	z, _ := argFloat("physicsApplyImpulse", args, 4)
	b.ApplyImpulse(Vec3{x, y, z})
	return null(), nil
}

// physicsApplyForceAtPoint(worldId, bodyId, fx, fy, fz, px, py, pz)
// Applies a force at a world-space point, generating both linear force and torque.
func builtinPhysicsApplyForceAtPoint(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("physicsApplyForceAtPoint expects worldId, bodyId, fx, fy, fz, px, py, pz")
	}
	w, _, err := requireWorld("physicsApplyForceAtPoint", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsApplyForceAtPoint", w, args, 1)
	if err != nil {
		return nil, err
	}
	fx, _ := argFloat("physicsApplyForceAtPoint", args, 2)
	fy, _ := argFloat("physicsApplyForceAtPoint", args, 3)
	fz, _ := argFloat("physicsApplyForceAtPoint", args, 4)
	px, _ := argFloat("physicsApplyForceAtPoint", args, 5)
	py, _ := argFloat("physicsApplyForceAtPoint", args, 6)
	pz, _ := argFloat("physicsApplyForceAtPoint", args, 7)
	force := Vec3{fx, fy, fz}
	point := Vec3{px, py, pz}
	b.ApplyForce(force)
	b.ApplyTorque(point.Sub(b.Position).Cross(force))
	return null(), nil
}

// physicsApplyImpulseAtPoint(worldId, bodyId, ix, iy, iz, px, py, pz)
// Applies an impulse at a specific world-space point, imparting both linear and angular change.
func builtinPhysicsApplyImpulseAtPoint(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("physicsApplyImpulseAtPoint expects worldId, bodyId, ix, iy, iz, px, py, pz")
	}
	w, _, err := requireWorld("physicsApplyImpulseAtPoint", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsApplyImpulseAtPoint", w, args, 1)
	if err != nil {
		return nil, err
	}
	ix, _ := argFloat("physicsApplyImpulseAtPoint", args, 2)
	iy, _ := argFloat("physicsApplyImpulseAtPoint", args, 3)
	iz, _ := argFloat("physicsApplyImpulseAtPoint", args, 4)
	px, _ := argFloat("physicsApplyImpulseAtPoint", args, 5)
	py, _ := argFloat("physicsApplyImpulseAtPoint", args, 6)
	pz, _ := argFloat("physicsApplyImpulseAtPoint", args, 7)
	b.ApplyImpulseAtPoint(Vec3{ix, iy, iz}, Vec3{px, py, pz})
	return null(), nil
}

// physicsApplyTorque(worldId, bodyId, x, y, z)
func builtinPhysicsApplyTorque(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("physicsApplyTorque expects worldId, bodyId, x, y, z")
	}
	w, _, err := requireWorld("physicsApplyTorque", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsApplyTorque", w, args, 1)
	if err != nil {
		return nil, err
	}
	x, _ := argFloat("physicsApplyTorque", args, 2)
	y, _ := argFloat("physicsApplyTorque", args, 3)
	z, _ := argFloat("physicsApplyTorque", args, 4)
	b.ApplyTorque(Vec3{x, y, z})
	return null(), nil
}

// physicsSetMass(worldId, bodyId, mass)
func builtinPhysicsSetMass(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetMass", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetMass", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetMass", w, args, 1)
	if err != nil {
		return nil, err
	}
	mass, _ := argFloat("physicsSetMass", args, 2)
	b.Mass = mass
	if mass > 0 {
		b.InvMass = 1.0 / mass
	} else {
		b.InvMass = 0
	}
	b.computeInertia()
	return null(), nil
}

// physicsGetMass(worldId, bodyId) → float
func builtinPhysicsGetMass(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetMass", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetMass", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetMass", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(b.Mass), nil
}

// physicsSetRestitution(worldId, bodyId, value)  0=no bounce, 1=perfect
func builtinPhysicsSetRestitution(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetRestitution", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetRestitution", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetRestitution", w, args, 1)
	if err != nil {
		return nil, err
	}
	v, _ := argFloat("physicsSetRestitution", args, 2)
	b.Restitution = v
	return null(), nil
}

// physicsSetFriction(worldId, bodyId, value)  0=frictionless, 1=high
func builtinPhysicsSetFriction(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetFriction", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetFriction", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetFriction", w, args, 1)
	if err != nil {
		return nil, err
	}
	v, _ := argFloat("physicsSetFriction", args, 2)
	b.Friction = v
	return null(), nil
}

// physicsSetLinearDrag(worldId, bodyId, value)
func builtinPhysicsSetLinearDrag(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetLinearDrag", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetLinearDrag", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetLinearDrag", w, args, 1)
	if err != nil {
		return nil, err
	}
	v, _ := argFloat("physicsSetLinearDrag", args, 2)
	b.LinearDrag = v
	return null(), nil
}

// physicsSetAngularDrag(worldId, bodyId, value)
func builtinPhysicsSetAngularDrag(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetAngularDrag", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetAngularDrag", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetAngularDrag", w, args, 1)
	if err != nil {
		return nil, err
	}
	v, _ := argFloat("physicsSetAngularDrag", args, 2)
	b.AngularDrag = v
	return null(), nil
}

// physicsSetActive(worldId, bodyId, bool)
func builtinPhysicsSetActive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetActive", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetActive", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetActive", w, args, 1)
	if err != nil {
		return nil, err
	}
	if args[2] != nil && args[2].Kind == candy_evaluator.ValBool {
		b.Active = args[2].B
	}
	return null(), nil
}

// physicsIsActive(worldId, bodyId) → bool
func builtinPhysicsIsActive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsIsActive", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsIsActive", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsIsActive", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(b.Active), nil
}

// physicsIsSleeping(worldId, bodyId) → bool
func builtinPhysicsIsSleeping(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsIsSleeping", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsIsSleeping", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsIsSleeping", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(b.Sleeping), nil
}

// physicsWakeBody(worldId, bodyId) — forces a sleeping body awake
func builtinPhysicsWakeBody(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsWakeBody", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsWakeBody", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsWakeBody", w, args, 1)
	if err != nil {
		return nil, err
	}
	b.Sleeping = false
	b.SleepTimer = 0
	return null(), nil
}

// physicsSetUserData(worldId, bodyId, value)
func builtinPhysicsSetUserData(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsSetUserData", args, 3); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsSetUserData", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsSetUserData", w, args, 1)
	if err != nil {
		return nil, err
	}
	ud, _ := argInt("physicsSetUserData", args, 2)
	b.UserData = ud
	return null(), nil
}

// physicsGetUserData(worldId, bodyId) → int
func builtinPhysicsGetUserData(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetUserData", args, 2); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetUserData", args, 0)
	if err != nil {
		return nil, err
	}
	b, _, err := requireBody("physicsGetUserData", w, args, 1)
	if err != nil {
		return nil, err
	}
	return vInt(b.UserData), nil
}

// ---- Contacts -----------------------------------------------------------

// physicsGetContacts(worldId) → [{bodyA, bodyB, normalX, normalY, normalZ, depth, pointX, pointY, pointZ}, ...]
func builtinPhysicsGetContacts(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetContacts", args, 1); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetContacts", args, 0)
	if err != nil {
		return nil, err
	}
	contacts := w.GetContacts()
	items := make([]candy_evaluator.Value, len(contacts))
	for i, c := range contacts {
		items[i] = candy_evaluator.Value{
			Kind: candy_evaluator.ValMap,
			StrMap: map[string]candy_evaluator.Value{
				"bodyA":   {Kind: candy_evaluator.ValInt, I64: c.BodyA},
				"bodyB":   {Kind: candy_evaluator.ValInt, I64: c.BodyB},
				"normalX": {Kind: candy_evaluator.ValFloat, F64: c.Normal.X},
				"normalY": {Kind: candy_evaluator.ValFloat, F64: c.Normal.Y},
				"normalZ": {Kind: candy_evaluator.ValFloat, F64: c.Normal.Z},
				"depth":   {Kind: candy_evaluator.ValFloat, F64: c.Depth},
				"pointX":  {Kind: candy_evaluator.ValFloat, F64: c.Point.X},
				"pointY":  {Kind: candy_evaluator.ValFloat, F64: c.Point.Y},
				"pointZ":  {Kind: candy_evaluator.ValFloat, F64: c.Point.Z},
			},
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: items}, nil
}

// ---- Raycasting ---------------------------------------------------------

// physicsCastRay(worldId, ox, oy, oz, dx, dy, dz, maxDist) → [{bodyId, distance, pointX, pointY, pointZ, normalX, normalY, normalZ}, ...]
func builtinPhysicsCastRay(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("physicsCastRay expects worldId, ox, oy, oz, dx, dy, dz, maxDist")
	}
	w, _, err := requireWorld("physicsCastRay", args, 0)
	if err != nil {
		return nil, err
	}
	ox, _ := argFloat("physicsCastRay", args, 1)
	oy, _ := argFloat("physicsCastRay", args, 2)
	oz, _ := argFloat("physicsCastRay", args, 3)
	dx, _ := argFloat("physicsCastRay", args, 4)
	dy, _ := argFloat("physicsCastRay", args, 5)
	dz, _ := argFloat("physicsCastRay", args, 6)
	maxD, _ := argFloat("physicsCastRay", args, 7)
	hits := w.CastRay(Vec3{ox, oy, oz}, Vec3{dx, dy, dz}, maxD)
	items := make([]candy_evaluator.Value, len(hits))
	for i, h := range hits {
		items[i] = candy_evaluator.Value{
			Kind: candy_evaluator.ValMap,
			StrMap: map[string]candy_evaluator.Value{
				"bodyId":   {Kind: candy_evaluator.ValInt, I64: h.BodyID},
				"distance": {Kind: candy_evaluator.ValFloat, F64: h.Distance},
				"pointX":   {Kind: candy_evaluator.ValFloat, F64: h.Point.X},
				"pointY":   {Kind: candy_evaluator.ValFloat, F64: h.Point.Y},
				"pointZ":   {Kind: candy_evaluator.ValFloat, F64: h.Point.Z},
				"normalX":  {Kind: candy_evaluator.ValFloat, F64: h.Normal.X},
				"normalY":  {Kind: candy_evaluator.ValFloat, F64: h.Normal.Y},
				"normalZ":  {Kind: candy_evaluator.ValFloat, F64: h.Normal.Z},
			},
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: items}, nil
}

// physicsCastRayFirst(worldId, ox, oy, oz, dx, dy, dz, maxDist) → map or null
func builtinPhysicsCastRayFirst(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	result, err := builtinPhysicsCastRay(args)
	if err != nil {
		return nil, err
	}
	if len(result.Elems) == 0 {
		return null(), nil
	}
	first := result.Elems[0]
	return &first, nil
}

// physicsGetBodyCount(worldId) → int
func builtinPhysicsGetBodyCount(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("physicsGetBodyCount", args, 1); err != nil {
		return nil, err
	}
	w, _, err := requireWorld("physicsGetBodyCount", args, 0)
	if err != nil {
		return nil, err
	}
	w.mu.Lock()
	count := int64(len(w.bodies))
	w.mu.Unlock()
	return vInt(count), nil
}

// ---- Motion type constants as builtins ----------------------------------

func builtinPhysicsMotionStatic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(MotionStatic)), nil
}
func builtinPhysicsMotionDynamic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(MotionDynamic)), nil
}
func builtinPhysicsMotionKinematic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vInt(int64(MotionKinematic)), nil
}
