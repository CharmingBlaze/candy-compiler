package candy_evaluator

import "math"

func mapVec3(m map[string]Value, key string) ([]float64, bool) {
	v, ok := m[key]
	if !ok || v.Kind != ValVec || len(v.Vec) < 3 {
		return nil, false
	}
	return v.Vec, true
}

func mapRadius(m map[string]Value) (float64, bool) {
	v, ok := m["radius"]
	if !ok {
		return 0, false
	}
	if v.Kind == ValFloat {
		return v.F64, true
	}
	if v.Kind == ValInt {
		return float64(v.I64), true
	}
	return 0, false
}

func aabbOverlapsMap(a, b *Value) (*Value, error) {
	ac, ok1 := mapVec3(a.StrMap, "center")
	ah, ok2 := mapVec3(a.StrMap, "halfExtents")
	bc, ok3 := mapVec3(b.StrMap, "center")
	bh, ok4 := mapVec3(b.StrMap, "halfExtents")
	if !(ok1 && ok2 && ok3 && ok4) {
		return &Value{Kind: ValBool, B: false}, nil
	}
	hit := math.Abs(ac[0]-bc[0]) <= (ah[0]+bh[0]) &&
		math.Abs(ac[1]-bc[1]) <= (ah[1]+bh[1]) &&
		math.Abs(ac[2]-bc[2]) <= (ah[2]+bh[2])
	return &Value{Kind: ValBool, B: hit}, nil
}

func aabbContainsPoint(a *Value, p *Value) (*Value, error) {
	ac, ok1 := mapVec3(a.StrMap, "center")
	ah, ok2 := mapVec3(a.StrMap, "halfExtents")
	if !(ok1 && ok2) || p == nil || p.Kind != ValVec || len(p.Vec) < 3 {
		return &Value{Kind: ValBool, B: false}, nil
	}
	hit := math.Abs(ac[0]-p.Vec[0]) <= ah[0] &&
		math.Abs(ac[1]-p.Vec[1]) <= ah[1] &&
		math.Abs(ac[2]-p.Vec[2]) <= ah[2]
	return &Value{Kind: ValBool, B: hit}, nil
}

func sphereOverlapsMap(a, b *Value) (*Value, error) {
	ac, ok1 := mapVec3(a.StrMap, "center")
	bc, ok2 := mapVec3(b.StrMap, "center")
	ar, ok3 := mapRadius(a.StrMap)
	br, ok4 := mapRadius(b.StrMap)
	if !(ok1 && ok2 && ok3 && ok4) {
		return &Value{Kind: ValBool, B: false}, nil
	}
	dx := ac[0] - bc[0]
	dy := ac[1] - bc[1]
	dz := ac[2] - bc[2]
	dist2 := dx*dx + dy*dy + dz*dz
	r := ar + br
	return &Value{Kind: ValBool, B: dist2 <= r*r}, nil
}

func rayIntersectsAABB(ray, aabb *Value) (*Value, error) {
	ro, ok1 := mapVec3(ray.StrMap, "origin")
	rd, ok2 := mapVec3(ray.StrMap, "direction")
	c, ok3 := mapVec3(aabb.StrMap, "center")
	h, ok4 := mapVec3(aabb.StrMap, "halfExtents")
	if !(ok1 && ok2 && ok3 && ok4) {
		return &Value{Kind: ValBool, B: false}, nil
	}
	minB := []float64{c[0] - h[0], c[1] - h[1], c[2] - h[2]}
	maxB := []float64{c[0] + h[0], c[1] + h[1], c[2] + h[2]}
	tmin := -math.MaxFloat64
	tmax := math.MaxFloat64
	for i := 0; i < 3; i++ {
		if math.Abs(rd[i]) < 1e-9 {
			if ro[i] < minB[i] || ro[i] > maxB[i] {
				return &Value{Kind: ValBool, B: false}, nil
			}
			continue
		}
		t1 := (minB[i] - ro[i]) / rd[i]
		t2 := (maxB[i] - ro[i]) / rd[i]
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		if t1 > tmin {
			tmin = t1
		}
		if t2 < tmax {
			tmax = t2
		}
		if tmin > tmax {
			return &Value{Kind: ValBool, B: false}, nil
		}
	}
	return &Value{Kind: ValBool, B: tmax >= 0}, nil
}
