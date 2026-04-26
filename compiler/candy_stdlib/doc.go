// Package candy_stdlib contains interpreter-side module source strings that are
// importable from Candy scripts (for example, `import candy.3d`).
//
// Module registration model:
//   - Host/runtime-backed modules (math, file, json, random, time, fs/rand aliases)
//     are injected from candy_evaluator/prelude.go.
//   - Script-backed modules are registered through Modules[...] in this package.
//
// Script-backed game modules currently shipped:
//   - candy.2d: entity graph, sprites, tilemap helpers, lightweight particles
//   - candy.3d: 3D entities/models/cameras/lights/animation helpers
//   - candy.physics2d and candy.physics3d: kinematic/rigid helpers + raycasts
//   - candy.ui: common UI widgets and layout primitives
//   - candy.proc: noise + dungeon/procedural helpers
//   - candy.vfx: post-process/effect orchestration surface
//   - candy.editor: runtime inspector and simple gizmos
//   - candy.scene, candy.audio, candy.input, candy.resources, candy.save,
//     candy.debug, candy.state: cross-cutting gameplay/runtime helpers
//   - candy.camera: 2D camera utilities
//   - candy.ai: steering primitives + behavior tree nodes
//   - candy.game3d: higher-level third-person camera rig + platform collisions
//   - candy.game: high-level facade for 2D/3D/app/multiplayer bootstrapping
//   - candy.network: ENet-backed client/server helpers
//
// Extra compatibility modules in Modules map:
//   - "std/*" helpers (import with quoted path, e.g. `import "std/time"`)
//   - "candy.math" and small "std/..." convenience wrappers
//
// Keep this package in sync with Modules map entries and prelude-exposed host
// module APIs when extending the standard library surface.
package candy_stdlib
