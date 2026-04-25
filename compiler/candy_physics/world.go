// Package candy_physics provides a pure-Go 3D rigid-body physics simulation.
// Zero CGO — works on Windows, macOS, and Linux with plain `go build`.
package candy_physics

import (
	"math"
	"sort"
	"sync"
)

// ---- Vec3 ---------------------------------------------------------------

type Vec3 struct{ X, Y, Z float64 }

func (v Vec3) Add(o Vec3) Vec3      { return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z} }
func (v Vec3) Sub(o Vec3) Vec3      { return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z} }
func (v Vec3) Scale(s float64) Vec3 { return Vec3{v.X * s, v.Y * s, v.Z * s} }
func (v Vec3) Neg() Vec3            { return Vec3{-v.X, -v.Y, -v.Z} }
func (v Vec3) Dot(o Vec3) float64   { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }
func (v Vec3) LenSq() float64       { return v.Dot(v) }
func (v Vec3) Len() float64         { return math.Sqrt(v.LenSq()) }
func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}
func (v Vec3) Normalize() Vec3 {
	l := v.Len()
	if l < 1e-12 {
		return Vec3{}
	}
	return v.Scale(1.0 / l)
}
func (v Vec3) Min(o Vec3) Vec3 {
	return Vec3{math.Min(v.X, o.X), math.Min(v.Y, o.Y), math.Min(v.Z, o.Z)}
}
func (v Vec3) Max(o Vec3) Vec3 {
	return Vec3{math.Max(v.X, o.X), math.Max(v.Y, o.Y), math.Max(v.Z, o.Z)}
}

// ---- Quaternion ---------------------------------------------------------

type Quat struct{ X, Y, Z, W float64 }

func QuatIdentity() Quat { return Quat{0, 0, 0, 1} }

func (q Quat) Mul(r Quat) Quat {
	return Quat{
		q.W*r.X + q.X*r.W + q.Y*r.Z - q.Z*r.Y,
		q.W*r.Y - q.X*r.Z + q.Y*r.W + q.Z*r.X,
		q.W*r.Z + q.X*r.Y - q.Y*r.X + q.Z*r.W,
		q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z,
	}
}

func (q Quat) Normalize() Quat {
	l := math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
	if l < 1e-12 {
		return QuatIdentity()
	}
	return Quat{q.X / l, q.Y / l, q.Z / l, q.W / l}
}

func (q Quat) RotateVec(v Vec3) Vec3 {
	qv := Vec3{q.X, q.Y, q.Z}
	t := qv.Cross(v).Scale(2)
	return v.Add(t.Scale(q.W)).Add(qv.Cross(t))
}

func quatFromAxisAngle(axis Vec3, angle float64) Quat {
	s := math.Sin(angle / 2)
	c := math.Cos(angle / 2)
	n := axis.Normalize()
	return Quat{n.X * s, n.Y * s, n.Z * s, c}
}

// ---- Shape --------------------------------------------------------------

type ShapeType int

const (
	ShapeBox     ShapeType = iota // half-extents stored in HalfExtent
	ShapeSphere                   // radius in HalfExtent.X
	ShapeCapsule                  // HalfExtent.X=radius, HalfExtent.Y=halfHeight along Y
	ShapePlane                    // infinite static ground plane (Y=0 by default, offset in HalfExtent.Y)
)

type Shape struct {
	Type       ShapeType
	HalfExtent Vec3
}

func NewBoxShape(halfW, halfH, halfD float64) Shape {
	return Shape{Type: ShapeBox, HalfExtent: Vec3{halfW, halfH, halfD}}
}
func NewSphereShape(radius float64) Shape {
	return Shape{Type: ShapeSphere, HalfExtent: Vec3{X: radius}}
}
func NewCapsuleShape(radius, halfHeight float64) Shape {
	return Shape{Type: ShapeCapsule, HalfExtent: Vec3{X: radius, Y: halfHeight}}
}
func NewPlaneShape(yOffset float64) Shape {
	return Shape{Type: ShapePlane, HalfExtent: Vec3{Y: yOffset}}
}

// worldAABB returns the axis-aligned bounding box of the shape at the given position.
func (s Shape) worldAABB(pos Vec3) (min, max Vec3) {
	switch s.Type {
	case ShapeSphere:
		r := s.HalfExtent.X
		half := Vec3{r, r, r}
		return pos.Sub(half), pos.Add(half)
	case ShapeCapsule:
		r := s.HalfExtent.X
		h := s.HalfExtent.Y + r
		half := Vec3{r, h, r}
		return pos.Sub(half), pos.Add(half)
	case ShapePlane:
		big := 1e9
		return Vec3{-big, pos.Y + s.HalfExtent.Y - 0.01, -big},
			Vec3{big, pos.Y + s.HalfExtent.Y + 0.01, big}
	default: // box
		return pos.Sub(s.HalfExtent), pos.Add(s.HalfExtent)
	}
}

// ---- Motion type --------------------------------------------------------

type MotionType int

const (
	MotionStatic MotionType = iota
	MotionDynamic
	MotionKinematic
)

// ---- Body ---------------------------------------------------------------

type Body struct {
	ID       int64
	Shape    Shape
	Motion   MotionType
	IsSensor bool

	Position   Vec3
	Rotation   Quat
	LinearVel  Vec3
	AngularVel Vec3
	Force      Vec3
	Torque     Vec3

	Mass    float64
	InvMass float64
	// InvInertia stores the diagonal of the inverse inertia tensor in body-local space.
	// Using a diagonal is exact for axis-aligned shapes and avoids a full 3x3 matrix.
	InvInertia  Vec3
	Restitution float64
	Friction    float64
	LinearDrag  float64
	AngularDrag float64

	// Sleep system: body stops being integrated when it has been nearly
	// still for SleepDelay seconds.
	Sleeping   bool
	SleepTimer float64

	UserData int64
	Active   bool
}

func newBody(id int64, shape Shape, pos Vec3, motion MotionType, isSensor bool) *Body {
	mass := 1.0
	invMass := 1.0
	if motion == MotionStatic {
		mass = 0
		invMass = 0
	}
	b := &Body{
		ID:          id,
		Shape:       shape,
		Motion:      motion,
		IsSensor:    isSensor,
		Position:    pos,
		Rotation:    QuatIdentity(),
		Mass:        mass,
		InvMass:     invMass,
		Restitution: 0.3,
		Friction:    0.5,
		LinearDrag:  0.01,
		AngularDrag: 0.05,
		Active:      true,
	}
	if motion == MotionDynamic {
		b.computeInertia()
	}
	return b
}

// computeInertia calculates the diagonal inverse inertia tensor for the body's shape.
// Must be called whenever Mass or Shape changes.
func (b *Body) computeInertia() {
	if b.Mass <= 0 {
		b.InvInertia = Vec3{}
		return
	}
	switch b.Shape.Type {
	case ShapeBox:
		w := b.Shape.HalfExtent.X * 2
		h := b.Shape.HalfExtent.Y * 2
		d := b.Shape.HalfExtent.Z * 2
		m := b.Mass
		Ixx := (m / 12.0) * (h*h + d*d)
		Iyy := (m / 12.0) * (w*w + d*d)
		Izz := (m / 12.0) * (w*w + h*h)
		b.InvInertia = Vec3{1 / Ixx, 1 / Iyy, 1 / Izz}
	case ShapeSphere:
		r := b.Shape.HalfExtent.X
		I := (2.0 / 5.0) * b.Mass * r * r
		b.InvInertia = Vec3{1 / I, 1 / I, 1 / I}
	case ShapeCapsule:
		// Approximate: solid cylinder formula
		r := b.Shape.HalfExtent.X
		h := b.Shape.HalfExtent.Y * 2
		m := b.Mass
		Iy := 0.5 * m * r * r
		Ixz := (m / 12.0) * (3*r*r + h*h)
		b.InvInertia = Vec3{1 / Ixz, 1 / Iy, 1 / Ixz}
	default:
		// Fallback: unit inertia
		b.InvInertia = Vec3{1, 1, 1}
	}
}

// mulInvInertia applies the diagonal inverse inertia tensor to a vector.
func (b *Body) mulInvInertia(v Vec3) Vec3 {
	return Vec3{b.InvInertia.X * v.X, b.InvInertia.Y * v.Y, b.InvInertia.Z * v.Z}
}

// ApplyForce accumulates a world-space force (applied during next Step).
func (b *Body) ApplyForce(f Vec3) { b.Force = b.Force.Add(f) }

// ApplyImpulse immediately changes linear velocity.
func (b *Body) ApplyImpulse(imp Vec3) {
	if b.InvMass > 0 {
		b.LinearVel = b.LinearVel.Add(imp.Scale(b.InvMass))
	}
}

// ApplyImpulseAtPoint applies a linear + angular impulse at a world-space point.
func (b *Body) ApplyImpulseAtPoint(imp, contactPoint Vec3) {
	if b.InvMass <= 0 {
		return
	}
	b.LinearVel = b.LinearVel.Add(imp.Scale(b.InvMass))
	r := contactPoint.Sub(b.Position)
	b.AngularVel = b.AngularVel.Add(b.mulInvInertia(r.Cross(imp)))
}

// ApplyTorque accumulates angular force.
func (b *Body) ApplyTorque(t Vec3) { b.Torque = b.Torque.Add(t) }

// ---- Contact / hit types ------------------------------------------------

type ContactEvent struct {
	BodyA, BodyB int64
	Normal       Vec3
	Depth        float64
	Point        Vec3
}

type RayHit struct {
	BodyID   int64
	Point    Vec3
	Normal   Vec3
	Distance float64
}

// ---- Collision helpers --------------------------------------------------

func aabbOverlap(aminX, aminY, aminZ, amaxX, amaxY, amaxZ,
	bminX, bminY, bminZ, bmaxX, bmaxY, bmaxZ float64) bool {
	return amaxX >= bminX && aminX <= bmaxX &&
		amaxY >= bminY && aminY <= bmaxY &&
		amaxZ >= bminZ && aminZ <= bmaxZ
}

// sphereVsSphere: returns normal (a→b), depth, ok
func sphereVsSphere(a, b *Body) (Vec3, float64, bool) {
	delta := b.Position.Sub(a.Position)
	dist := delta.Len()
	sum := a.Shape.HalfExtent.X + b.Shape.HalfExtent.X
	if dist >= sum {
		return Vec3{}, 0, false
	}
	depth := sum - dist
	var normal Vec3
	if dist < 1e-10 {
		normal = Vec3{0, 1, 0}
	} else {
		normal = delta.Scale(1.0 / dist)
	}
	return normal, depth, true
}

// boxVsBox using SAT on 3 axes
func boxVsBox(a, b *Body) (Vec3, float64, bool) {
	d := b.Position.Sub(a.Position)
	he := a.Shape.HalfExtent
	hb := b.Shape.HalfExtent

	overlapX := he.X + hb.X - math.Abs(d.X)
	if overlapX <= 0 {
		return Vec3{}, 0, false
	}
	overlapY := he.Y + hb.Y - math.Abs(d.Y)
	if overlapY <= 0 {
		return Vec3{}, 0, false
	}
	overlapZ := he.Z + hb.Z - math.Abs(d.Z)
	if overlapZ <= 0 {
		return Vec3{}, 0, false
	}

	// smallest penetration axis
	if overlapX <= overlapY && overlapX <= overlapZ {
		nx := 1.0
		if d.X < 0 {
			nx = -1
		}
		return Vec3{nx, 0, 0}, overlapX, true
	}
	if overlapY <= overlapX && overlapY <= overlapZ {
		ny := 1.0
		if d.Y < 0 {
			ny = -1
		}
		return Vec3{0, ny, 0}, overlapY, true
	}
	nz := 1.0
	if d.Z < 0 {
		nz = -1
	}
	return Vec3{0, 0, nz}, overlapZ, true
}

// sphereVsBox: normal points from box centre to sphere centre
func sphereVsBox(sph, box *Body) (Vec3, float64, bool) {
	// closest point on box to sphere centre
	bMin := box.Position.Sub(box.Shape.HalfExtent)
	bMax := box.Position.Add(box.Shape.HalfExtent)
	closest := Vec3{
		math.Max(bMin.X, math.Min(sph.Position.X, bMax.X)),
		math.Max(bMin.Y, math.Min(sph.Position.Y, bMax.Y)),
		math.Max(bMin.Z, math.Min(sph.Position.Z, bMax.Z)),
	}
	diff := sph.Position.Sub(closest)
	distSq := diff.LenSq()
	r := sph.Shape.HalfExtent.X
	if distSq >= r*r {
		return Vec3{}, 0, false
	}
	dist := math.Sqrt(distSq)
	depth := r - dist
	var normal Vec3
	if dist < 1e-10 {
		normal = Vec3{0, 1, 0}
	} else {
		normal = diff.Scale(1.0 / dist)
	}
	return normal, depth, true
}

// capsuleVsPlane: capsule is body a, plane body b
func capsuleVsPlane(cap, plane *Body) (Vec3, float64, bool) {
	planeY := plane.Position.Y + plane.Shape.HalfExtent.Y
	r := cap.Shape.HalfExtent.X
	h := cap.Shape.HalfExtent.Y
	// lowest point of capsule
	lowest := cap.Position.Y - h - r
	depth := planeY - lowest
	if depth <= 0 {
		return Vec3{}, 0, false
	}
	return Vec3{0, 1, 0}, depth, true
}

// sphereVsPlane
func sphereVsPlane(sph, plane *Body) (Vec3, float64, bool) {
	planeY := plane.Position.Y + plane.Shape.HalfExtent.Y
	r := sph.Shape.HalfExtent.X
	depth := planeY - (sph.Position.Y - r)
	if depth <= 0 {
		return Vec3{}, 0, false
	}
	return Vec3{0, 1, 0}, depth, true
}

// boxVsPlane
func boxVsPlane(box, plane *Body) (Vec3, float64, bool) {
	planeY := plane.Position.Y + plane.Shape.HalfExtent.Y
	lowest := box.Position.Y - box.Shape.HalfExtent.Y
	depth := planeY - lowest
	if depth <= 0 {
		return Vec3{}, 0, false
	}
	return Vec3{0, 1, 0}, depth, true
}

// ---- World --------------------------------------------------------------

type World struct {
	mu       sync.Mutex
	nextID   int64
	bodies   map[int64]*Body
	gravity  Vec3
	contacts []ContactEvent
	tmpIDs   []int64
	tmpMins  []Vec3
	tmpMaxs  []Vec3
	// Broadphase configuration.
	broadphaseMode BroadphaseMode
	cellSize       float64
	tmpPairKeys    []uint64
	// Spatial-hash scratch (reused between steps to reduce allocations).
	cellBuckets map[[3]int][]int
	cellKeys    [][3]int
	seenPairs   map[uint64]struct{}
}

type BroadphaseMode int

const (
	BroadphaseAABBPrune BroadphaseMode = iota
	BroadphaseSpatialHash
)

func newWorld() *World {
	return &World{
		nextID:         1,
		bodies:         make(map[int64]*Body),
		gravity:        Vec3{0, -9.81, 0},
		broadphaseMode: BroadphaseAABBPrune,
		cellSize:       2.0,
		cellBuckets:    make(map[[3]int][]int),
		seenPairs:      make(map[uint64]struct{}),
	}
}

func (w *World) SetBroadphaseMode(mode BroadphaseMode) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.broadphaseMode = mode
}

func (w *World) SetSpatialHashCellSize(size float64) {
	if size <= 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cellSize = size
}

func (w *World) CreateBody(shape Shape, pos Vec3, motion MotionType, isSensor bool) int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	id := w.nextID
	w.nextID++
	w.bodies[id] = newBody(id, shape, pos, motion, isSensor)
	return id
}

func (w *World) DestroyBody(id int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.bodies, id)
}

func (w *World) GetBody(id int64) (*Body, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	b, ok := w.bodies[id]
	return b, ok
}

func (w *World) SetGravity(g Vec3) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.gravity = g
}

func (w *World) GetGravity() Vec3 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.gravity
}

func (w *World) GetContacts() []ContactEvent {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]ContactEvent, len(w.contacts))
	copy(out, w.contacts)
	return out
}

// Step advances the simulation by dt seconds.
func (w *World) Step(dt float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.contacts = w.contacts[:0]

	// 1. Integrate forces → velocity
	for _, b := range w.bodies {
		if !b.Active || b.Motion != MotionDynamic {
			b.Force = Vec3{}
			b.Torque = Vec3{}
			continue
		}
		// Sleep check: skip integration if body is sleeping
		const sleepLinThresh = 0.02
		const sleepAngThresh = 0.05
		const sleepDelay = 0.5
		if b.LinearVel.Len() < sleepLinThresh && b.AngularVel.Len() < sleepAngThresh &&
			b.Force.LenSq() < 1e-6 {
			b.SleepTimer += dt
			if b.SleepTimer >= sleepDelay {
				b.Sleeping = true
				b.LinearVel = Vec3{}
				b.AngularVel = Vec3{}
				b.Force = Vec3{}
				b.Torque = Vec3{}
				continue
			}
		} else {
			b.SleepTimer = 0
			b.Sleeping = false
		}
		// gravity
		b.LinearVel = b.LinearVel.Add(w.gravity.Scale(dt))
		// applied force
		b.LinearVel = b.LinearVel.Add(b.Force.Scale(b.InvMass * dt))
		// angular acceleration: α = I⁻¹ · τ
		b.AngularVel = b.AngularVel.Add(b.mulInvInertia(b.Torque).Scale(dt))
		// drag
		b.LinearVel = b.LinearVel.Scale(math.Max(0, 1-b.LinearDrag*dt))
		b.AngularVel = b.AngularVel.Scale(math.Max(0, 1-b.AngularDrag*dt))
		b.Force = Vec3{}
		b.Torque = Vec3{}
	}

	// 2. Integrate velocity → position & rotation
	for _, b := range w.bodies {
		if !b.Active || b.Motion == MotionStatic {
			continue
		}
		b.Position = b.Position.Add(b.LinearVel.Scale(dt))
		avLen := b.AngularVel.Len()
		if avLen > 1e-10 {
			angle := avLen * dt
			b.Rotation = b.Rotation.Mul(quatFromAxisAngle(b.AngularVel.Normalize(), angle)).Normalize()
		}
	}

	// 3. Collect sorted pairs for deterministic collision order
	ids := w.tmpIDs[:0]
	for id := range w.bodies {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	mins := w.tmpMins[:0]
	maxs := w.tmpMaxs[:0]
	for i, id := range ids {
		b := w.bodies[id]
		min, max := b.Shape.worldAABB(b.Position)
		if i < len(mins) {
			mins[i] = min
			maxs[i] = max
		} else {
			mins = append(mins, min)
			maxs = append(maxs, max)
		}
	}
	w.tmpIDs = ids
	w.tmpMins = mins
	w.tmpMaxs = maxs

	if w.broadphaseMode == BroadphaseSpatialHash {
		pairs := w.collectSpatialHashPairs(ids, mins, maxs)
		for _, key := range pairs {
			i := int(key >> 32)
			j := int(key & 0xffffffff)
			a := w.bodies[ids[i]]
			b := w.bodies[ids[j]]
			if !a.Active || !b.Active {
				continue
			}
			if a.Motion == MotionStatic && b.Motion == MotionStatic {
				continue
			}
			if !broadPhaseOverlapAABB(mins[i], maxs[i], mins[j], maxs[j]) {
				continue
			}
			w.narrowPhase(a, b)
		}
	} else {
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				a := w.bodies[ids[i]]
				b := w.bodies[ids[j]]
				if !a.Active || !b.Active {
					continue
				}
				if a.Motion == MotionStatic && b.Motion == MotionStatic {
					continue
				}
				if !broadPhaseOverlapAABB(mins[i], maxs[i], mins[j], maxs[j]) {
					continue
				}
				w.narrowPhase(a, b)
			}
		}
	}
}

func spatialCellRange(min, max Vec3, cellSize float64) (int, int, int, int, int, int) {
	cx0 := int(math.Floor(min.X / cellSize))
	cy0 := int(math.Floor(min.Y / cellSize))
	cz0 := int(math.Floor(min.Z / cellSize))
	cx1 := int(math.Floor(max.X / cellSize))
	cy1 := int(math.Floor(max.Y / cellSize))
	cz1 := int(math.Floor(max.Z / cellSize))
	return cx0, cy0, cz0, cx1, cy1, cz1
}

func pairKey(i, j int) uint64 {
	return (uint64(i) << 32) | uint64(j)
}

func (w *World) collectSpatialHashPairs(ids []int64, mins, maxs []Vec3) []uint64 {
	// Reset reusable buckets.
	for _, k := range w.cellKeys {
		if bucket, ok := w.cellBuckets[k]; ok {
			w.cellBuckets[k] = bucket[:0]
		}
	}
	w.cellKeys = w.cellKeys[:0]

	cellSize := w.cellSize
	if cellSize <= 0 {
		cellSize = 2.0
	}
	for i := 0; i < len(ids); i++ {
		cx0, cy0, cz0, cx1, cy1, cz1 := spatialCellRange(mins[i], maxs[i], cellSize)
		for cx := cx0; cx <= cx1; cx++ {
			for cy := cy0; cy <= cy1; cy++ {
				for cz := cz0; cz <= cz1; cz++ {
					k := [3]int{cx, cy, cz}
					bucket, ok := w.cellBuckets[k]
					if !ok {
						// Reused buckets keep their capacity; fresh buckets start small.
						bucket = make([]int, 0, 8)
						w.cellBuckets[k] = bucket
						w.cellKeys = append(w.cellKeys, k)
					}
					w.cellBuckets[k] = append(w.cellBuckets[k], i)
				}
			}
		}
	}

	// Deterministic cell traversal order.
	sort.Slice(w.cellKeys, func(i, j int) bool {
		a, b := w.cellKeys[i], w.cellKeys[j]
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		if a[1] != b[1] {
			return a[1] < b[1]
		}
		return a[2] < b[2]
	})

	for k := range w.seenPairs {
		delete(w.seenPairs, k)
	}
	keys := w.tmpPairKeys[:0]
	for _, k := range w.cellKeys {
		bucket := w.cellBuckets[k]
		for a := 0; a < len(bucket); a++ {
			i := bucket[a]
			for b := a + 1; b < len(bucket); b++ {
				j := bucket[b]
				ii, jj := i, j
				if ii > jj {
					ii, jj = jj, ii
				}
				k := pairKey(ii, jj)
				if _, ok := w.seenPairs[k]; ok {
					continue
				}
				w.seenPairs[k] = struct{}{}
				keys = append(keys, k)
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	w.tmpPairKeys = keys
	return keys
}

// broadPhaseOverlapAABB is a deterministic, conservative AABB broadphase filter.
// Non-overlapping AABBs cannot collide in narrowphase this step.
func broadPhaseOverlapAABB(amin, amax, bmin, bmax Vec3) bool {
	return aabbOverlap(
		amin.X, amin.Y, amin.Z, amax.X, amax.Y, amax.Z,
		bmin.X, bmin.Y, bmin.Z, bmax.X, bmax.Y, bmax.Z,
	)
}

func (w *World) narrowPhase(a, b *Body) {
	var (
		normal Vec3
		depth  float64
		ok     bool
	)

	// Route to the correct collision test, ensuring plane is always second
	at, bt := a.Shape.Type, b.Shape.Type

	switch {
	case at == ShapePlane && bt != ShapePlane:
		n, d, o := w.testPair(b, a)
		normal, depth, ok = n.Neg(), d, o
	case bt == ShapePlane:
		normal, depth, ok = w.testPair(a, b)
	default:
		normal, depth, ok = w.testPair(a, b)
	}

	if !ok || depth <= 0 {
		return
	}

	// Wake sleeping bodies on contact
	a.Sleeping = false
	a.SleepTimer = 0
	b.Sleeping = false
	b.SleepTimer = 0

	// Compute actual contact point: midpoint on the penetrating surface
	// For body A: move from its centre along -normal by half the depth
	contactPoint := a.Position.Sub(normal.Scale(depth * 0.5))

	contact := ContactEvent{
		BodyA:  a.ID,
		BodyB:  b.ID,
		Normal: normal,
		Depth:  depth,
		Point:  contactPoint,
	}
	w.contacts = append(w.contacts, contact)

	if a.IsSensor || b.IsSensor {
		return
	}

	w.resolve(a, b, normal, depth, contactPoint)
}

func (w *World) testPair(a, b *Body) (Vec3, float64, bool) {
	at, bt := a.Shape.Type, b.Shape.Type
	switch {
	case at == ShapeSphere && bt == ShapeSphere:
		return sphereVsSphere(a, b)
	case at == ShapeSphere && bt == ShapeBox:
		return sphereVsBox(a, b)
	case at == ShapeBox && bt == ShapeSphere:
		n, d, ok := sphereVsBox(b, a)
		return n.Neg(), d, ok
	case at == ShapeBox && bt == ShapeBox:
		return boxVsBox(a, b)
	case at == ShapeSphere && bt == ShapePlane:
		return sphereVsPlane(a, b)
	case at == ShapeBox && bt == ShapePlane:
		return boxVsPlane(a, b)
	case at == ShapeCapsule && bt == ShapePlane:
		return capsuleVsPlane(a, b)
	case at == ShapeCapsule && bt == ShapeSphere:
		// treat capsule as sphere at centre for simplicity
		return sphereVsSphere(a, b)
	case at == ShapeSphere && bt == ShapeCapsule:
		n, d, ok := sphereVsSphere(b, a)
		return n.Neg(), d, ok
	case at == ShapeCapsule && bt == ShapeBox:
		return sphereVsBox(a, b)
	case at == ShapeBox && bt == ShapeCapsule:
		n, d, ok := sphereVsBox(b, a)
		return n.Neg(), d, ok
	case at == ShapeCapsule && bt == ShapeCapsule:
		return sphereVsSphere(a, b)
	default:
		// AABB broad check for unknown combos
		amin, amax := a.Shape.worldAABB(a.Position)
		bmin, bmax := b.Shape.worldAABB(b.Position)
		if !aabbOverlap(amin.X, amin.Y, amin.Z, amax.X, amax.Y, amax.Z,
			bmin.X, bmin.Y, bmin.Z, bmax.X, bmax.Y, bmax.Z) {
			return Vec3{}, 0, false
		}
		return Vec3{0, 1, 0}, 0.01, true
	}
}

func (w *World) resolve(a, b *Body, normal Vec3, depth float64, contact Vec3) {
	// Positional correction (Baumgarte — prevents sinking)
	const slop = 0.005
	const percent = 0.4
	totalInvMass := a.InvMass + b.InvMass
	if totalInvMass < 1e-12 {
		return
	}
	correction := math.Max(depth-slop, 0) * percent / totalInvMass
	if a.Motion == MotionDynamic {
		a.Position = a.Position.Sub(normal.Scale(correction * a.InvMass))
	}
	if b.Motion == MotionDynamic {
		b.Position = b.Position.Add(normal.Scale(correction * b.InvMass))
	}

	// Lever arms from each body centre to the contact point
	rA := contact.Sub(a.Position)
	rB := contact.Sub(b.Position)

	// Velocity at contact point (linear + angular contribution)
	vA := a.LinearVel.Add(a.AngularVel.Cross(rA))
	vB := b.LinearVel.Add(b.AngularVel.Cross(rB))
	relVel := vA.Sub(vB)

	vAlongNormal := relVel.Dot(normal)
	if vAlongNormal > 0 {
		return
	}

	// Effective inverse mass including rotational contribution:
	// 1/mEff = 1/mA + 1/mB + (I⁻¹·(rA×n))×rA·n + (I⁻¹·(rB×n))×rB·n
	rACrossN := rA.Cross(normal)
	rBCrossN := rB.Cross(normal)
	rotTermA := a.mulInvInertia(rACrossN).Cross(rA).Dot(normal)
	rotTermB := b.mulInvInertia(rBCrossN).Cross(rB).Dot(normal)
	effInvMass := totalInvMass + rotTermA + rotTermB
	if effInvMass < 1e-12 {
		return
	}

	e := math.Min(a.Restitution, b.Restitution)
	j := -(1 + e) * vAlongNormal / effInvMass
	impulse := normal.Scale(j)

	// Apply linear + angular impulse
	if a.Motion == MotionDynamic {
		a.LinearVel = a.LinearVel.Sub(impulse.Scale(a.InvMass))
		a.AngularVel = a.AngularVel.Sub(a.mulInvInertia(rA.Cross(impulse)))
	}
	if b.Motion == MotionDynamic {
		b.LinearVel = b.LinearVel.Add(impulse.Scale(b.InvMass))
		b.AngularVel = b.AngularVel.Add(b.mulInvInertia(rB.Cross(impulse)))
	}

	// Friction impulse
	tangent := relVel.Sub(normal.Scale(relVel.Dot(normal))).Normalize()
	if tangent.LenSq() < 1e-12 {
		return
	}
	rACrossT := rA.Cross(tangent)
	rBCrossT := rB.Cross(tangent)
	rotTermAT := a.mulInvInertia(rACrossT).Cross(rA).Dot(tangent)
	rotTermBT := b.mulInvInertia(rBCrossT).Cross(rB).Dot(tangent)
	effInvMassT := totalInvMass + rotTermAT + rotTermBT
	if effInvMassT < 1e-12 {
		return
	}
	jt := -relVel.Dot(tangent) / effInvMassT
	mu := math.Sqrt(a.Friction * b.Friction)
	var frictionImpulse Vec3
	if math.Abs(jt) < j*mu {
		frictionImpulse = tangent.Scale(jt)
	} else {
		frictionImpulse = tangent.Scale(-j * mu)
	}
	if a.Motion == MotionDynamic {
		a.LinearVel = a.LinearVel.Sub(frictionImpulse.Scale(a.InvMass))
		a.AngularVel = a.AngularVel.Sub(a.mulInvInertia(rA.Cross(frictionImpulse)))
	}
	if b.Motion == MotionDynamic {
		b.LinearVel = b.LinearVel.Add(frictionImpulse.Scale(b.InvMass))
		b.AngularVel = b.AngularVel.Add(b.mulInvInertia(rB.Cross(frictionImpulse)))
	}
}

// CastRay returns all hits sorted by distance. origin+direction*maxDist is the ray.
func (w *World) CastRay(origin, dir Vec3, maxDist float64) []RayHit {
	w.mu.Lock()
	defer w.mu.Unlock()

	dir = dir.Normalize()
	var hits []RayHit

	for _, b := range w.bodies {
		if !b.Active {
			continue
		}
		hit, dist, ok := rayVsBody(origin, dir, maxDist, b)
		if ok {
			normal := origin.Add(dir.Scale(dist)).Sub(b.Position).Normalize()
			hits = append(hits, RayHit{
				BodyID:   b.ID,
				Point:    hit,
				Normal:   normal,
				Distance: dist,
			})
		}
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Distance < hits[j].Distance })
	return hits
}

func rayVsBody(origin, dir Vec3, maxDist float64, b *Body) (Vec3, float64, bool) {
	switch b.Shape.Type {
	case ShapeSphere:
		return rayVsSphere(origin, dir, maxDist, b.Position, b.Shape.HalfExtent.X)
	case ShapePlane:
		return rayVsPlane(origin, dir, maxDist, b.Position.Y+b.Shape.HalfExtent.Y)
	default:
		bmin, bmax := b.Shape.worldAABB(b.Position)
		return rayVsAABB(origin, dir, maxDist, bmin, bmax)
	}
}

func rayVsAABB(origin, dir Vec3, maxDist float64, bmin, bmax Vec3) (Vec3, float64, bool) {
	tmin := 0.0
	tmax := maxDist
	for i := 0; i < 3; i++ {
		var o, d, mn, mx float64
		switch i {
		case 0:
			o, d, mn, mx = origin.X, dir.X, bmin.X, bmax.X
		case 1:
			o, d, mn, mx = origin.Y, dir.Y, bmin.Y, bmax.Y
		default:
			o, d, mn, mx = origin.Z, dir.Z, bmin.Z, bmax.Z
		}
		if math.Abs(d) < 1e-12 {
			if o < mn || o > mx {
				return Vec3{}, 0, false
			}
			continue
		}
		t1 := (mn - o) / d
		t2 := (mx - o) / d
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		tmin = math.Max(tmin, t1)
		tmax = math.Min(tmax, t2)
		if tmin > tmax {
			return Vec3{}, 0, false
		}
	}
	pt := origin.Add(dir.Scale(tmin))
	return pt, tmin, true
}

func rayVsSphere(origin, dir Vec3, maxDist float64, centre Vec3, radius float64) (Vec3, float64, bool) {
	oc := origin.Sub(centre)
	b := 2 * oc.Dot(dir)
	c := oc.LenSq() - radius*radius
	disc := b*b - 4*c
	if disc < 0 {
		return Vec3{}, 0, false
	}
	t := (-b - math.Sqrt(disc)) / 2
	if t < 0 {
		t = (-b + math.Sqrt(disc)) / 2
	}
	if t < 0 || t > maxDist {
		return Vec3{}, 0, false
	}
	return origin.Add(dir.Scale(t)), t, true
}

func rayVsPlane(origin, dir Vec3, maxDist, planeY float64) (Vec3, float64, bool) {
	if math.Abs(dir.Y) < 1e-12 {
		return Vec3{}, 0, false
	}
	t := (planeY - origin.Y) / dir.Y
	if t < 0 || t > maxDist {
		return Vec3{}, 0, false
	}
	return origin.Add(dir.Scale(t)), t, true
}

// ---- World registry (multiple worlds keyed by int64 ID) -----------------

var (
	worldMu   sync.Mutex
	worldMap        = map[int64]*World{}
	nextWorld int64 = 1
)

func NewWorld() int64 {
	worldMu.Lock()
	defer worldMu.Unlock()
	id := nextWorld
	nextWorld++
	worldMap[id] = newWorld()
	return id
}

func GetWorld(id int64) (*World, bool) {
	worldMu.Lock()
	defer worldMu.Unlock()
	w, ok := worldMap[id]
	return w, ok
}

func DestroyWorld(id int64) {
	worldMu.Lock()
	defer worldMu.Unlock()
	delete(worldMap, id)
}
