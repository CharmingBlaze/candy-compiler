package candy_physics

import "testing"

func setupDynamicBodies(world *World, n int) {
	cols := 32
	for i := 0; i < n; i++ {
		x := float64(i % cols)
		z := float64(i / cols)
		y := 10.0 + float64(i%7)
		world.CreateBody(NewSphereShape(0.5), Vec3{X: x, Y: y, Z: z}, MotionDynamic, false)
	}
}

func BenchmarkWorldStep_100Bodies(b *testing.B) {
	w := newWorld()
	setupDynamicBodies(w, 100)
	dt := 1.0 / 60.0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Step(dt)
	}
}

func BenchmarkWorldStep_500Bodies(b *testing.B) {
	w := newWorld()
	setupDynamicBodies(w, 500)
	dt := 1.0 / 60.0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Step(dt)
	}
}

func BenchmarkWorldStep_1000Bodies(b *testing.B) {
	w := newWorld()
	setupDynamicBodies(w, 1000)
	dt := 1.0 / 60.0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Step(dt)
	}
}

func BenchmarkWorldStep_1000Bodies_SpatialHash(b *testing.B) {
	w := newWorld()
	w.SetBroadphaseMode(BroadphaseSpatialHash)
	w.SetSpatialHashCellSize(2.0)
	setupDynamicBodies(w, 1000)
	dt := 1.0 / 60.0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Step(dt)
	}
}

