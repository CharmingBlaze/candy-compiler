package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Map helpers ----

func vec2ToMap(v rl.Vector2) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(v.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(v.Y)},
	})
}

func vec3ToMap(v rl.Vector3) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(v.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(v.Y)},
		"z": {Kind: candy_evaluator.ValFloat, F64: float64(v.Z)},
	})
}

func quatToMap(q rl.Quaternion) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"x": {Kind: candy_evaluator.ValFloat, F64: float64(q.X)},
		"y": {Kind: candy_evaluator.ValFloat, F64: float64(q.Y)},
		"z": {Kind: candy_evaluator.ValFloat, F64: float64(q.Z)},
		"w": {Kind: candy_evaluator.ValFloat, F64: float64(q.W)},
	})
}

func matToMap(m rl.Matrix) *candy_evaluator.Value {
	return vMap(map[string]candy_evaluator.Value{
		"m0":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M0)},
		"m1":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M1)},
		"m2":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M2)},
		"m3":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M3)},
		"m4":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M4)},
		"m5":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M5)},
		"m6":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M6)},
		"m7":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M7)},
		"m8":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M8)},
		"m9":  {Kind: candy_evaluator.ValFloat, F64: float64(m.M9)},
		"m10": {Kind: candy_evaluator.ValFloat, F64: float64(m.M10)},
		"m11": {Kind: candy_evaluator.ValFloat, F64: float64(m.M11)},
		"m12": {Kind: candy_evaluator.ValFloat, F64: float64(m.M12)},
		"m13": {Kind: candy_evaluator.ValFloat, F64: float64(m.M13)},
		"m14": {Kind: candy_evaluator.ValFloat, F64: float64(m.M14)},
		"m15": {Kind: candy_evaluator.ValFloat, F64: float64(m.M15)},
	})
}

func mapFloat(m map[string]candy_evaluator.Value, key string) float32 {
	if v, ok := m[key]; ok {
		if v.Kind == candy_evaluator.ValFloat {
			return float32(v.F64)
		}
		if v.Kind == candy_evaluator.ValInt {
			return float32(v.I64)
		}
	}
	return 0
}

func argQuat(name string, args []*candy_evaluator.Value, i int) (rl.Quaternion, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.Quaternion{}, fmt.Errorf("%s arg %d must be map {x,y,z,w}", name, i+1)
	}
	m := args[i].StrMap
	return rl.Quaternion{
		X: mapFloat(m, "x"),
		Y: mapFloat(m, "y"),
		Z: mapFloat(m, "z"),
		W: mapFloat(m, "w"),
	}, nil
}

func argMatrix(name string, args []*candy_evaluator.Value, i int) (rl.Matrix, error) {
	if i >= len(args) || args[i] == nil || args[i].Kind != candy_evaluator.ValMap {
		return rl.Matrix{}, fmt.Errorf("%s arg %d must be matrix map {m0..m15}", name, i+1)
	}
	m := args[i].StrMap
	return rl.Matrix{
		M0: mapFloat(m, "m0"), M4: mapFloat(m, "m4"), M8: mapFloat(m, "m8"), M12: mapFloat(m, "m12"),
		M1: mapFloat(m, "m1"), M5: mapFloat(m, "m5"), M9: mapFloat(m, "m9"), M13: mapFloat(m, "m13"),
		M2: mapFloat(m, "m2"), M6: mapFloat(m, "m6"), M10: mapFloat(m, "m10"), M14: mapFloat(m, "m14"),
		M3: mapFloat(m, "m3"), M7: mapFloat(m, "m7"), M11: mapFloat(m, "m11"), M15: mapFloat(m, "m15"),
	}, nil
}

// ---- Utils math ----

func builtinMathClamp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("mathClamp", args, 3); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("mathClamp", args, 0)
	mn, _ := getArgFloat("mathClamp", args, 1)
	mx, _ := getArgFloat("mathClamp", args, 2)
	return vFloat(float64(rl.Clamp(float32(v), float32(mn), float32(mx)))), nil
}

func builtinMathLerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("mathLerp", args, 3); err != nil {
		return nil, err
	}
	s, _ := getArgFloat("mathLerp", args, 0)
	e, _ := getArgFloat("mathLerp", args, 1)
	a, _ := getArgFloat("mathLerp", args, 2)
	return vFloat(float64(rl.Lerp(float32(s), float32(e), float32(a)))), nil
}

func builtinMathNormalize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("mathNormalize", args, 3); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("mathNormalize", args, 0)
	s, _ := getArgFloat("mathNormalize", args, 1)
	e, _ := getArgFloat("mathNormalize", args, 2)
	return vFloat(float64(rl.Normalize(float32(v), float32(s), float32(e)))), nil
}

func builtinMathRemap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("mathRemap", args, 5); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("mathRemap", args, 0)
	is, _ := getArgFloat("mathRemap", args, 1)
	ie, _ := getArgFloat("mathRemap", args, 2)
	os, _ := getArgFloat("mathRemap", args, 3)
	oe, _ := getArgFloat("mathRemap", args, 4)
	return vFloat(float64(rl.Remap(float32(v), float32(is), float32(ie), float32(os), float32(oe)))), nil
}

func builtinMathWrap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("mathWrap", args, 3); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("mathWrap", args, 0)
	mn, _ := getArgFloat("mathWrap", args, 1)
	mx, _ := getArgFloat("mathWrap", args, 2)
	return vFloat(float64(rl.Wrap(float32(v), float32(mn), float32(mx)))), nil
}

func builtinFloatEquals(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("floatEquals", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("floatEquals", args, 0)
	y, _ := getArgFloat("floatEquals", args, 1)
	return vBool(rl.FloatEquals(float32(x), float32(y))), nil
}

// ---- Vector2 math ----

func builtinVector2Zero(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.Vector2Zero()), nil
}

func builtinVector2One(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec2ToMap(rl.Vector2One()), nil
}

func builtinVector2Add(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Add", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Add", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Add", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Add(v1, v2)), nil
}

func builtinVector2AddValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2AddValue", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2AddValue", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector2AddValue", args, 1)
	return vec2ToMap(rl.Vector2AddValue(v, float32(a))), nil
}

func builtinVector2Subtract(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Subtract", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Subtract", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Subtract", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Subtract(v1, v2)), nil
}

func builtinVector2SubtractValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2SubtractValue", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2SubtractValue", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("vector2SubtractValue", args, 1)
	return vec2ToMap(rl.Vector2SubtractValue(v, float32(s))), nil
}

func builtinVector2Length(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Length", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Length", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2Length(v))), nil
}

func builtinVector2LengthSqr(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2LengthSqr", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2LengthSqr", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2LengthSqr(v))), nil
}

func builtinVector2DotProduct(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2DotProduct", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2DotProduct", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2DotProduct", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2DotProduct(v1, v2))), nil
}

func builtinVector2CrossProduct(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2CrossProduct", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2CrossProduct", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2CrossProduct", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2CrossProduct(v1, v2))), nil
}

func builtinVector2Distance(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Distance", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Distance", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Distance", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2Distance(v1, v2))), nil
}

func builtinVector2DistanceSqr(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2DistanceSqr", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2DistanceSqr", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2DistanceSqr", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2DistanceSqr(v1, v2))), nil
}

func builtinVector2Angle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Angle", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Angle", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Angle", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2Angle(v1, v2))), nil
}

func builtinVector2LineAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2LineAngle", args, 2); err != nil {
		return nil, err
	}
	s, err := argVector2("vector2LineAngle", args, 0)
	if err != nil {
		return nil, err
	}
	e, err := argVector2("vector2LineAngle", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector2LineAngle(s, e))), nil
}

func builtinVector2Scale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Scale", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Scale", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("vector2Scale", args, 1)
	return vec2ToMap(rl.Vector2Scale(v, float32(s))), nil
}

func builtinVector2Multiply(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Multiply", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Multiply", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Multiply", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Multiply(v1, v2)), nil
}

func builtinVector2Negate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Negate", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Negate", args, 0)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Negate(v)), nil
}

func builtinVector2Divide(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Divide", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Divide", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Divide", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Divide(v1, v2)), nil
}

func builtinVector2Normalize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Normalize", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Normalize", args, 0)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Normalize(v)), nil
}

func builtinVector2Transform(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Transform", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Transform", args, 0)
	if err != nil {
		return nil, err
	}
	m, err := argMatrix("vector2Transform", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Transform(v, m)), nil
}

func builtinVector2Lerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Lerp", args, 3); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Lerp", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Lerp", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector2Lerp", args, 2)
	return vec2ToMap(rl.Vector2Lerp(v1, v2, float32(a))), nil
}

func builtinVector2Reflect(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Reflect", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Reflect", args, 0)
	if err != nil {
		return nil, err
	}
	n, err := argVector2("vector2Reflect", args, 1)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Reflect(v, n)), nil
}

func builtinVector2Rotate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Rotate", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Rotate", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector2Rotate", args, 1)
	return vec2ToMap(rl.Vector2Rotate(v, float32(a))), nil
}

func builtinVector2MoveTowards(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2MoveTowards", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2MoveTowards", args, 0)
	if err != nil {
		return nil, err
	}
	t, err := argVector2("vector2MoveTowards", args, 1)
	if err != nil {
		return nil, err
	}
	d, _ := getArgFloat("vector2MoveTowards", args, 2)
	return vec2ToMap(rl.Vector2MoveTowards(v, t, float32(d))), nil
}

func builtinVector2Invert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Invert", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Invert", args, 0)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Invert(v)), nil
}

func builtinVector2Clamp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Clamp", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2Clamp", args, 0)
	if err != nil {
		return nil, err
	}
	mn, err := argVector2("vector2Clamp", args, 1)
	if err != nil {
		return nil, err
	}
	mx, err := argVector2("vector2Clamp", args, 2)
	if err != nil {
		return nil, err
	}
	return vec2ToMap(rl.Vector2Clamp(v, mn, mx)), nil
}

func builtinVector2ClampValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2ClampValue", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector2("vector2ClampValue", args, 0)
	if err != nil {
		return nil, err
	}
	mn, _ := getArgFloat("vector2ClampValue", args, 1)
	mx, _ := getArgFloat("vector2ClampValue", args, 2)
	return vec2ToMap(rl.Vector2ClampValue(v, float32(mn), float32(mx))), nil
}

func builtinVector2Equals(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector2Equals", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector2("vector2Equals", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector2("vector2Equals", args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(rl.Vector2Equals(v1, v2)), nil
}

// ---- Vector3 math ----

func builtinVector3Zero(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec3ToMap(rl.Vector3Zero()), nil
}

func builtinVector3One(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vec3ToMap(rl.Vector3One()), nil
}

func builtinVector3Add(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Add", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Add", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Add", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Add(v1, v2)), nil
}

func builtinVector3AddValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3AddValue", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3AddValue", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector3AddValue", args, 1)
	return vec3ToMap(rl.Vector3AddValue(v, float32(a))), nil
}

func builtinVector3Subtract(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Subtract", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Subtract", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Subtract", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Subtract(v1, v2)), nil
}

func builtinVector3SubtractValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3SubtractValue", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3SubtractValue", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("vector3SubtractValue", args, 1)
	return vec3ToMap(rl.Vector3SubtractValue(v, float32(s))), nil
}

func builtinVector3Scale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Scale", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Scale", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("vector3Scale", args, 1)
	return vec3ToMap(rl.Vector3Scale(v, float32(s))), nil
}

func builtinVector3Multiply(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Multiply", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Multiply", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Multiply", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Multiply(v1, v2)), nil
}

func builtinVector3CrossProduct(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3CrossProduct", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3CrossProduct", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3CrossProduct", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3CrossProduct(v1, v2)), nil
}

func builtinVector3Perpendicular(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Perpendicular", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Perpendicular", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Perpendicular(v)), nil
}

func builtinVector3Length(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Length", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Length", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3Length(v))), nil
}

func builtinVector3LengthSqr(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3LengthSqr", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3LengthSqr", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3LengthSqr(v))), nil
}

func builtinVector3DotProduct(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3DotProduct", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3DotProduct", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3DotProduct", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3DotProduct(v1, v2))), nil
}

func builtinVector3Distance(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Distance", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Distance", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Distance", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3Distance(v1, v2))), nil
}

func builtinVector3DistanceSqr(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3DistanceSqr", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3DistanceSqr", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3DistanceSqr", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3DistanceSqr(v1, v2))), nil
}

func builtinVector3Angle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Angle", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Angle", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Angle", args, 1)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.Vector3Angle(v1, v2))), nil
}

func builtinVector3Negate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Negate", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Negate", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Negate(v)), nil
}

func builtinVector3Divide(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Divide", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Divide", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Divide", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Divide(v1, v2)), nil
}

func builtinVector3Normalize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Normalize", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Normalize", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Normalize(v)), nil
}

func builtinVector3Project(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Project", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Project", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Project", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Project(v1, v2)), nil
}

func builtinVector3Reject(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Reject", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Reject", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Reject", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Reject(v1, v2)), nil
}

func builtinVector3Transform(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Transform", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Transform", args, 0)
	if err != nil {
		return nil, err
	}
	m, err := argMatrix("vector3Transform", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Transform(v, m)), nil
}

func builtinVector3RotateByQuaternion(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3RotateByQuaternion", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3RotateByQuaternion", args, 0)
	if err != nil {
		return nil, err
	}
	q, err := argQuat("vector3RotateByQuaternion", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3RotateByQuaternion(v, q)), nil
}

func builtinVector3RotateByAxisAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3RotateByAxisAngle", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3RotateByAxisAngle", args, 0)
	if err != nil {
		return nil, err
	}
	ax, err := argVector3("vector3RotateByAxisAngle", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector3RotateByAxisAngle", args, 2)
	return vec3ToMap(rl.Vector3RotateByAxisAngle(v, ax, float32(a))), nil
}

func builtinVector3Lerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Lerp", args, 3); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Lerp", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Lerp", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("vector3Lerp", args, 2)
	return vec3ToMap(rl.Vector3Lerp(v1, v2, float32(a))), nil
}

func builtinVector3Reflect(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Reflect", args, 2); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Reflect", args, 0)
	if err != nil {
		return nil, err
	}
	n, err := argVector3("vector3Reflect", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Reflect(v, n)), nil
}

func builtinVector3Min(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Min", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Min", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Min", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Min(v1, v2)), nil
}

func builtinVector3Max(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Max", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Max", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Max", args, 1)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Max(v1, v2)), nil
}

func builtinVector3Barycenter(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Barycenter", args, 4); err != nil {
		return nil, err
	}
	p, err := argVector3("vector3Barycenter", args, 0)
	if err != nil {
		return nil, err
	}
	a, err := argVector3("vector3Barycenter", args, 1)
	if err != nil {
		return nil, err
	}
	b, err := argVector3("vector3Barycenter", args, 2)
	if err != nil {
		return nil, err
	}
	c, err := argVector3("vector3Barycenter", args, 3)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Barycenter(p, a, b, c)), nil
}

func builtinVector3Unproject(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Unproject", args, 3); err != nil {
		return nil, err
	}
	src, err := argVector3("vector3Unproject", args, 0)
	if err != nil {
		return nil, err
	}
	proj, err := argMatrix("vector3Unproject", args, 1)
	if err != nil {
		return nil, err
	}
	view, err := argMatrix("vector3Unproject", args, 2)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Unproject(src, proj, view)), nil
}

func builtinVector3Invert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Invert", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Invert", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Invert(v)), nil
}

func builtinVector3Clamp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Clamp", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Clamp", args, 0)
	if err != nil {
		return nil, err
	}
	mn, err := argVector3("vector3Clamp", args, 1)
	if err != nil {
		return nil, err
	}
	mx, err := argVector3("vector3Clamp", args, 2)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.Vector3Clamp(v, mn, mx)), nil
}

func builtinVector3ClampValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3ClampValue", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3ClampValue", args, 0)
	if err != nil {
		return nil, err
	}
	mn, _ := getArgFloat("vector3ClampValue", args, 1)
	mx, _ := getArgFloat("vector3ClampValue", args, 2)
	return vec3ToMap(rl.Vector3ClampValue(v, float32(mn), float32(mx))), nil
}

func builtinVector3Equals(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Equals", args, 2); err != nil {
		return nil, err
	}
	v1, err := argVector3("vector3Equals", args, 0)
	if err != nil {
		return nil, err
	}
	v2, err := argVector3("vector3Equals", args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(rl.Vector3Equals(v1, v2)), nil
}

func builtinVector3Refract(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("vector3Refract", args, 3); err != nil {
		return nil, err
	}
	v, err := argVector3("vector3Refract", args, 0)
	if err != nil {
		return nil, err
	}
	n, err := argVector3("vector3Refract", args, 1)
	if err != nil {
		return nil, err
	}
	r, _ := getArgFloat("vector3Refract", args, 2)
	return vec3ToMap(rl.Vector3Refract(v, n, float32(r))), nil
}

// ---- Matrix math ----

func builtinMatrixDeterminant(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixDeterminant", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("matrixDeterminant", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.MatrixDeterminant(m))), nil
}

func builtinMatrixTrace(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixTrace", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("matrixTrace", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.MatrixTrace(m))), nil
}

func builtinMatrixTranspose(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixTranspose", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("matrixTranspose", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixTranspose(m)), nil
}

func builtinMatrixInvert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixInvert", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("matrixInvert", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixInvert(m)), nil
}

func builtinMatrixIdentity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return matToMap(rl.MatrixIdentity()), nil
}

func builtinMatrixAdd(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixAdd", args, 2); err != nil {
		return nil, err
	}
	l, err := argMatrix("matrixAdd", args, 0)
	if err != nil {
		return nil, err
	}
	r, err := argMatrix("matrixAdd", args, 1)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixAdd(l, r)), nil
}

func builtinMatrixSubtract(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixSubtract", args, 2); err != nil {
		return nil, err
	}
	l, err := argMatrix("matrixSubtract", args, 0)
	if err != nil {
		return nil, err
	}
	r, err := argMatrix("matrixSubtract", args, 1)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixSubtract(l, r)), nil
}

func builtinMatrixMultiply(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixMultiply", args, 2); err != nil {
		return nil, err
	}
	l, err := argMatrix("matrixMultiply", args, 0)
	if err != nil {
		return nil, err
	}
	r, err := argMatrix("matrixMultiply", args, 1)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixMultiply(l, r)), nil
}

func builtinMatrixTranslate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixTranslate", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("matrixTranslate", args, 0)
	y, _ := getArgFloat("matrixTranslate", args, 1)
	z, _ := getArgFloat("matrixTranslate", args, 2)
	return matToMap(rl.MatrixTranslate(float32(x), float32(y), float32(z))), nil
}

func builtinMatrixRotate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotate", args, 2); err != nil {
		return nil, err
	}
	ax, err := argVector3("matrixRotate", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("matrixRotate", args, 1)
	return matToMap(rl.MatrixRotate(ax, float32(a))), nil
}

func builtinMatrixRotateX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotateX", args, 1); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("matrixRotateX", args, 0)
	return matToMap(rl.MatrixRotateX(float32(a))), nil
}

func builtinMatrixRotateY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotateY", args, 1); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("matrixRotateY", args, 0)
	return matToMap(rl.MatrixRotateY(float32(a))), nil
}

func builtinMatrixRotateZ(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotateZ", args, 1); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("matrixRotateZ", args, 0)
	return matToMap(rl.MatrixRotateZ(float32(a))), nil
}

func builtinMatrixRotateXYZ(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotateXYZ", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("matrixRotateXYZ", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixRotateXYZ(v)), nil
}

func builtinMatrixRotateZYX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixRotateZYX", args, 1); err != nil {
		return nil, err
	}
	v, err := argVector3("matrixRotateZYX", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixRotateZYX(v)), nil
}

func builtinMatrixScale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixScale", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("matrixScale", args, 0)
	y, _ := getArgFloat("matrixScale", args, 1)
	z, _ := getArgFloat("matrixScale", args, 2)
	return matToMap(rl.MatrixScale(float32(x), float32(y), float32(z))), nil
}

func builtinMatrixFrustum(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixFrustum", args, 6); err != nil {
		return nil, err
	}
	l, _ := getArgFloat("matrixFrustum", args, 0)
	r, _ := getArgFloat("matrixFrustum", args, 1)
	b, _ := getArgFloat("matrixFrustum", args, 2)
	t, _ := getArgFloat("matrixFrustum", args, 3)
	n, _ := getArgFloat("matrixFrustum", args, 4)
	f, _ := getArgFloat("matrixFrustum", args, 5)
	return matToMap(rl.MatrixFrustum(float32(l), float32(r), float32(b), float32(t), float32(n), float32(f))), nil
}

func builtinMatrixPerspective(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixPerspective", args, 4); err != nil {
		return nil, err
	}
	fov, _ := getArgFloat("matrixPerspective", args, 0)
	asp, _ := getArgFloat("matrixPerspective", args, 1)
	near, _ := getArgFloat("matrixPerspective", args, 2)
	far, _ := getArgFloat("matrixPerspective", args, 3)
	return matToMap(rl.MatrixPerspective(float32(fov), float32(asp), float32(near), float32(far))), nil
}

func builtinMatrixOrtho(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixOrtho", args, 6); err != nil {
		return nil, err
	}
	l, _ := getArgFloat("matrixOrtho", args, 0)
	r, _ := getArgFloat("matrixOrtho", args, 1)
	b, _ := getArgFloat("matrixOrtho", args, 2)
	t, _ := getArgFloat("matrixOrtho", args, 3)
	n, _ := getArgFloat("matrixOrtho", args, 4)
	f, _ := getArgFloat("matrixOrtho", args, 5)
	return matToMap(rl.MatrixOrtho(float32(l), float32(r), float32(b), float32(t), float32(n), float32(f))), nil
}

func builtinMatrixLookAt(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixLookAt", args, 3); err != nil {
		return nil, err
	}
	eye, err := argVector3("matrixLookAt", args, 0)
	if err != nil {
		return nil, err
	}
	target, err := argVector3("matrixLookAt", args, 1)
	if err != nil {
		return nil, err
	}
	up, err := argVector3("matrixLookAt", args, 2)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.MatrixLookAt(eye, target, up)), nil
}

// ---- Quaternion math ----

func builtinQuaternionAdd(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionAdd", args, 2); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionAdd", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionAdd", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionAdd(q1, q2)), nil
}

func builtinQuaternionAddValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionAddValue", args, 2); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionAddValue", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("quaternionAddValue", args, 1)
	return quatToMap(rl.QuaternionAddValue(q, float32(a))), nil
}

func builtinQuaternionSubtract(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionSubtract", args, 2); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionSubtract", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionSubtract", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionSubtract(q1, q2)), nil
}

func builtinQuaternionSubtractValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionSubtractValue", args, 2); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionSubtractValue", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("quaternionSubtractValue", args, 1)
	return quatToMap(rl.QuaternionSubtractValue(q, float32(s))), nil
}

func builtinQuaternionIdentity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return quatToMap(rl.QuaternionIdentity()), nil
}

func builtinQuaternionLength(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionLength", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionLength", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.QuaternionLength(q))), nil
}

func builtinQuaternionNormalize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionNormalize", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionNormalize", args, 0)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionNormalize(q)), nil
}

func builtinQuaternionInvert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionInvert", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionInvert", args, 0)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionInvert(q)), nil
}

func builtinQuaternionMultiply(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionMultiply", args, 2); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionMultiply", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionMultiply", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionMultiply(q1, q2)), nil
}

func builtinQuaternionScale(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionScale", args, 2); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionScale", args, 0)
	if err != nil {
		return nil, err
	}
	s, _ := getArgFloat("quaternionScale", args, 1)
	return quatToMap(rl.QuaternionScale(q, float32(s))), nil
}

func builtinQuaternionDivide(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionDivide", args, 2); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionDivide", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionDivide", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionDivide(q1, q2)), nil
}

func builtinQuaternionLerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionLerp", args, 3); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionLerp", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionLerp", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("quaternionLerp", args, 2)
	return quatToMap(rl.QuaternionLerp(q1, q2, float32(a))), nil
}

func builtinQuaternionNlerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionNlerp", args, 3); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionNlerp", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionNlerp", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("quaternionNlerp", args, 2)
	return quatToMap(rl.QuaternionNlerp(q1, q2, float32(a))), nil
}

func builtinQuaternionSlerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionSlerp", args, 3); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionSlerp", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionSlerp", args, 1)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("quaternionSlerp", args, 2)
	return quatToMap(rl.QuaternionSlerp(q1, q2, float32(a))), nil
}

func builtinQuaternionFromVector3ToVector3(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionFromVector3ToVector3", args, 2); err != nil {
		return nil, err
	}
	from, err := argVector3("quaternionFromVector3ToVector3", args, 0)
	if err != nil {
		return nil, err
	}
	to, err := argVector3("quaternionFromVector3ToVector3", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionFromVector3ToVector3(from, to)), nil
}

func builtinQuaternionFromMatrix(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionFromMatrix", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("quaternionFromMatrix", args, 0)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionFromMatrix(m)), nil
}

func builtinQuaternionToMatrix(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionToMatrix", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionToMatrix", args, 0)
	if err != nil {
		return nil, err
	}
	return matToMap(rl.QuaternionToMatrix(q)), nil
}

func builtinQuaternionFromAxisAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionFromAxisAngle", args, 2); err != nil {
		return nil, err
	}
	ax, err := argVector3("quaternionFromAxisAngle", args, 0)
	if err != nil {
		return nil, err
	}
	a, _ := getArgFloat("quaternionFromAxisAngle", args, 1)
	return quatToMap(rl.QuaternionFromAxisAngle(ax, float32(a))), nil
}

func builtinQuaternionToAxisAngle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionToAxisAngle", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionToAxisAngle", args, 0)
	if err != nil {
		return nil, err
	}
	var outAxis rl.Vector3
	var outAngle float32
	rl.QuaternionToAxisAngle(q, &outAxis, &outAngle)
	return vMap(map[string]candy_evaluator.Value{
		"axis":  {Kind: candy_evaluator.ValMap, StrMap: vec3ToMap(outAxis).StrMap},
		"angle": {Kind: candy_evaluator.ValFloat, F64: float64(outAngle)},
	}), nil
}

func builtinQuaternionFromEuler(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionFromEuler", args, 3); err != nil {
		return nil, err
	}
	pitch, _ := getArgFloat("quaternionFromEuler", args, 0)
	yaw, _ := getArgFloat("quaternionFromEuler", args, 1)
	roll, _ := getArgFloat("quaternionFromEuler", args, 2)
	return quatToMap(rl.QuaternionFromEuler(float32(pitch), float32(yaw), float32(roll))), nil
}

func builtinQuaternionToEuler(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionToEuler", args, 1); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionToEuler", args, 0)
	if err != nil {
		return nil, err
	}
	return vec3ToMap(rl.QuaternionToEuler(q)), nil
}

func builtinQuaternionTransform(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionTransform", args, 2); err != nil {
		return nil, err
	}
	q, err := argQuat("quaternionTransform", args, 0)
	if err != nil {
		return nil, err
	}
	m, err := argMatrix("quaternionTransform", args, 1)
	if err != nil {
		return nil, err
	}
	return quatToMap(rl.QuaternionTransform(q, m)), nil
}

func builtinQuaternionEquals(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("quaternionEquals", args, 2); err != nil {
		return nil, err
	}
	q1, err := argQuat("quaternionEquals", args, 0)
	if err != nil {
		return nil, err
	}
	q2, err := argQuat("quaternionEquals", args, 1)
	if err != nil {
		return nil, err
	}
	return vBool(rl.QuaternionEquals(q1, q2)), nil
}

func builtinMatrixDecompose(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("matrixDecompose", args, 1); err != nil {
		return nil, err
	}
	m, err := argMatrix("matrixDecompose", args, 0)
	if err != nil {
		return nil, err
	}
	var t rl.Vector3
	var r rl.Quaternion
	var s rl.Vector3
	rl.MatrixDecompose(m, &t, &r, &s)
	return vMap(map[string]candy_evaluator.Value{
		"translation": {Kind: candy_evaluator.ValMap, StrMap: vec3ToMap(t).StrMap},
		"rotation":    {Kind: candy_evaluator.ValMap, StrMap: quatToMap(r).StrMap},
		"scale":       {Kind: candy_evaluator.ValMap, StrMap: vec3ToMap(s).StrMap},
	}), nil
}

// suppress unused import warning
var _ = fmt.Sprintf
