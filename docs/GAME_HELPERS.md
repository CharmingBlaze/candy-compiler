# Candy Game Helpers Built-in Commands

This document defines the high-level game helper API as flat built-in commands (no user-defined structs required to use them). It is intended as a practical reference and implementation target for host/runtime builtins.

Where a helper returns a compound value, it returns a Candy map (for example `{x, y}`), consistent with the rest of the extension surface.

Implementation source: `compiler/candy_raylib/game_helpers.go` (registered by `compiler/candy_raylib/register.go`).

## Design goals

- Keep gameplay code short and beginner-friendly.
- Prefer one built-in call over multi-step boilerplate.
- Make frame-driven behavior explicit (`update*` calls and global per-frame updates).
- Keep naming consistent with existing Candy/Raylib style.

## Implementation status

- All commands listed in this file are registered built-ins.
- Core gameplay helpers (camera/input/collision/timers/tween/animation/path/spatial/particles/scene/save-load/debug) are fully runnable.
- Advanced systems helpers (AI/networking/waves/remote interpolation) currently provide a baseline runtime surface with lightweight default behavior and are designed to be extended without breaking script APIs.

## 1) Camera commands

### Smooth camera follow

- `cameraFollow(targetX, targetY, smoothness)`
  - Sets camera follow target and smoothing speed.
  - Typical range for `smoothness`: `1.0` to `12.0`.
- `cameraSnapTo(x, y)`
  - Instantly moves camera center.
- `cameraBounds(minX, minY, maxX, maxY)`
  - Clamps camera center to world bounds.
- `cameraShake(intensity)`
  - Adds temporary shake impulse. Larger value means stronger shake.
- `cameraZoom(level)`
  - Sets camera zoom (`1.0` = default).
- `cameraPosition() -> {x, y}`
  - Returns current camera center.
- `screenToWorld(x, y) -> {x, y}`
- `worldToScreen(x, y) -> {x, y}`

### Runtime behavior notes

- Follow and shake should be resolved once per frame in runtime update.
- Shake should decay smoothly over time.
- Final camera transform should apply in this order: follow -> shake -> bounds -> zoom.

### 3D camera helpers

- `camera3DFollow(targetX, targetY, targetZ, distance, height, smoothness)`
  - Third-person follow camera that trails and smooths movement.
- `cameraOrbit(targetX, targetY, targetZ, radius, speed)`
  - Continuously orbits around a target at fixed radius.
- `cameraOrbitInput(targetX, targetY, targetZ, radius, sensitivity)`
  - Mouse/right-stick controlled orbit camera.
- `cameraLookAt(x, y, z)`
  - Redirects current 3D camera target.
- `cameraYaw(angleDegrees)`
- `cameraPitch(angleDegrees)`
- `cameraRoll(angleDegrees)`
- `cameraForward() -> {x, y, z}`
- `cameraRight() -> {x, y, z}`
- `cameraUp() -> {x, y, z}`
- `screenToWorld3D(mouseX, mouseY, planeY) -> {x, y, z}`
  - Projects mouse ray to a horizontal plane (`y = planeY`).

### 3D camera behavior notes

- Pitch should clamp (for example `-80` to `80`) to avoid camera flips.
- Orbit input should support optional invert-Y and sensitivity scaling.
- `camera3DFollow` should preserve target visibility by keeping a minimum distance.

## 2) Input commands

- `keyDown(keyName) -> bool`
- `keyPressed(keyName) -> bool`
- `keyReleased(keyName) -> bool`
- `axis(name) -> float`
  - `horizontal`: left/a = `-1`, right/d = `+1`
  - `vertical`: up/w = `-1`, down/s = `+1`
- `combo(k1, k2, ..., kN) -> bool`
  - Buffered sequence check for fighting-game style inputs.
- `mousePressed(button) -> bool`
  - `0` left, `1` right, `2` middle.
- `mouseX() -> float`
- `mouseY() -> float`
- `mousePosition() -> {x, y}`
- `gamepadButton(index, buttonName) -> bool`
- `gamepadAxis(index, stickName) -> {x, y}`

### Runtime behavior notes

- `keyPressed` / `keyReleased` depend on previous-frame state.
- Input state should be sampled once per frame, then queried by all script calls.
- Combo buffers should expire old entries using time windows.

## 3) Animation commands

- `animation(name, frames, frameTime, loop) -> animId`
- `playAnimation(animId)`
- `updateAnimation(animId, deltaTime)`
- `animationFrame(animId) -> int`
- `animationDone(animId) -> bool`

### Animation controller helpers

- `animController() -> controllerId`
- `addAnimation(controllerId, name, frames, frameTime, loop)`
- `setAnimation(controllerId, name)`
- `updateAnimController(controllerId, deltaTime)`
- `controllerFrame(controllerId) -> int`

### Runtime behavior notes

- `updateAnimation` should be safe for missing/invalid IDs (no crash).
- One-shot non-loop animations should freeze on final frame and report done.

## 4) Tween commands (smooth interpolation)

- `tweenCreate(start, end, duration, easingType) -> tweenId`
- `updateTween(tweenId, deltaTime)`
- `tweenValue(tweenId) -> float`
- `tweenDone(tweenId) -> bool`
- `tweenTo(object, field, targetValue, duration [, onComplete])`

### Easing types

- `linear`
- `easeIn`
- `easeOut`
- `easeInOut`
- `bounce`
- `elastic`

### Runtime behavior notes

- Clamp tween progress to `[0.0, 1.0]`.
- Invalid easing strings should fall back to `linear`.

## 5) Timer commands

- `createTimer(durationSeconds) -> timerId`
- `updateTimer(timerId, deltaTime)`
- `timerDone(timerId) -> bool`
- `timerRemaining(timerId) -> float`
- `resetTimer(timerId, durationSeconds)`
- `after(delaySeconds, fn)`
- `every(intervalSeconds, fn)`
- `cooldownReady(name) -> bool`
- `cooldownStart(name, durationSeconds)`

### Runtime behavior notes

- Timers and cooldowns should support global per-frame updates for convenience.
- `after` and `every` should run callbacks on the main script thread.

## 6) Particle system

- `particles(x, y) -> emitterId`
- `particleSpread(emitterId, degrees)`
- `particleSpeed(emitterId, min, max)`
- `particleLife(emitterId, minSeconds, maxSeconds)`
- `particleColor(emitterId, startColor, endColor)`
- `particleSize(emitterId, startSize, endSize)`
- `particleGravity(emitterId, gx, gy)`
- `emit(emitterId, count)`
- `updateParticles(emitterId, deltaTime)`
- `drawParticles(emitterId)`

### Presets

- `explosion(x, y, size)`
- `smoke(x, y)`
- `sparkles(x, y)`
- `blood(x, y)`

### Runtime behavior notes

- Emitters should cap total live particles to avoid runaway memory use.
- Particle updates should remove dead particles in-place each frame.

## 7) Collision helpers

- `boxCollision(x1, y1, w1, h1, x2, y2, w2, h2) -> bool`
- `circleCollision(x1, y1, r1, x2, y2, r2) -> bool`
- `pointInBox(px, py, x, y, w, h) -> bool`
- `raycast(x1, y1, x2, y2, obstacles) -> bool|map`
- `inRadius(x, y, radius, objects) -> list`

## 8) Pathfinding

- `pathfindGrid(width, height, tileSize) -> gridId`
- `blockTile(gridId, tileX, tileY)`
- `findPath(gridId, startX, startY, goalX, goalY) -> list`
- `moveTowards(object, targetX, targetY, speed)`

### Runtime behavior notes

- `findPath` should return empty list when no path exists.
- Returned path points should be world-space points for direct movement.

## 9) Spatial partitioning (optimization)

- `spatialGrid(worldWidth, worldHeight, cellSize) -> gridId`
- `insert(gridId, object)`
- `queryGrid(gridId, x, y, w, h) -> list`
- `clearGrid(gridId)`

### Runtime behavior notes

- Typical loop: `clearGrid` -> `insert` all collidables -> `queryGrid` per actor.

## 10) Audio helpers

- `playSound(path [, volume] [, pitch])`
- `playMusic(path)`
- `pauseMusic()`
- `resumeMusic()`
- `stopMusic()`
- `fadeMusic(seconds)`
- `fadeMusicIn(seconds)`
- `volume(level)`
- `playSound3D(path, x, y, z, maxDistance)`

## 11) Scene / state management

- `scene(name, callbacksMap)`
- `startScene(name)`
- `switchScene(name)`
- `pauseScene()`
- `resumeScene()`

### Callback map keys

- `init: fun() { ... }`
- `update: fun(dt) { ... }`
- `draw: fun() { ... }`

## 12) Save / load system

- `save(key, value)`
- `load(key [, default]) -> value`
- `deleteSave(key)`
- `saveExists(key) -> bool`

### Runtime behavior notes

- Values should round-trip through JSON-compatible maps/lists/scalars.
- Missing key should return provided default, otherwise `null`.

## 13) Debug helpers

- `debugText(x, y, text)`
- `debugBox(x, y, w, h, color)`
- `debugCircle(x, y, radius, color)`
- `debugPath(path, color)`
- `startProfile(name)`
- `endProfile(name)`
- `printProfiles()`

## 14) 3D transform and movement helpers

- `rotateX(object, degrees)`
- `rotateY(object, degrees)`
- `rotateZ(object, degrees)`
- `rotateTo(object, pitch, yaw, roll, speed)`
  - Smoothly rotates object toward target Euler angles.
- `lookAt(object, targetX, targetY, targetZ [, upX, upY, upZ])`
  - Sets orientation so object faces target.
- `orbit(object, centerX, centerY, centerZ, radius, speed [, axis])`
  - Moves object around center on selected axis (`"x"`, `"y"`, `"z"`; default `"y"`).
- `moveForward(object, speed, deltaTime)`
- `moveRight(object, speed, deltaTime)`
- `moveUp(object, speed, deltaTime)`
- `faceVelocity(object, lerpSpeed)`
  - Auto-rotate toward current velocity direction.
- `clampPosition(object, minX, minY, minZ, maxX, maxY, maxZ)`

## 15) Combat and projectile helpers

- `projectile(x, y, z, dirX, dirY, dirZ, speed, lifeSeconds) -> projectileId`
- `updateProjectiles(deltaTime)`
- `drawProjectiles()`
- `projectileHit(projectileId, targetList [, radius]) -> bool|map`
- `hitscan(originX, originY, originZ, dirX, dirY, dirZ, maxDistance, targets) -> map|null`
- `lockOnNearest(x, y, z, radius, targets) -> target|null`
- `damage(target, amount [, damageType])`
- `heal(target, amount)`
- `alive(target) -> bool`
- `team(target) -> string|int`
- `setTeam(target, value)`

## 16) Spawning and wave helpers

- `spawn(pathOrPrefab, x, y [, z]) -> entity`
- `despawn(entity)`
- `spawnAtMarker(markerName, prefab) -> entity`
- `poolCreate(prefab, capacity) -> poolId`
- `poolGet(poolId, x, y [, z]) -> entity`
- `poolRelease(poolId, entity)`
- `waveCreate() -> waveId`
- `waveAdd(waveId, prefab, count, intervalSeconds)`
- `waveStart(waveId)`
- `waveDone(waveId) -> bool`

## 17) AI behavior helpers

- `stateMachine() -> smId`
- `stateAdd(smId, name, callbacksMap)`
- `stateSet(smId, name)`
- `stateUpdate(smId, dt)`
- `patrol(points, speed) -> behaviorId`
- `chase(target, speed, stopDistance) -> behaviorId`
- `flee(target, speed, safeDistance) -> behaviorId`
- `wander(speed, radius) -> behaviorId`
- `lineOfSight(x1, y1, z1, x2, y2, z2, obstacles) -> bool`
- `canSee(observer, target, fovDegrees, viewDistance, obstacles) -> bool`

## 18) UI and quest helpers

- `button(id, x, y, w, h, label) -> bool`
- `slider(id, x, y, w, min, max, value) -> float`
- `healthBar(x, y, w, h, current, max, color)`
- `floatingText(text, x, y [, z], seconds [, color])`
- `minimap(worldX, worldY, worldW, worldH, screenX, screenY, screenW, screenH)`
- `questAdd(id, title, description)`
- `questComplete(id)`
- `questState(id) -> string`
- `questStep(id, stepText)`

## 19) Runtime utility helpers

- `tag(entity, name)`
- `untag(entity, name)`
- `withTag(name) -> list`
- `distance2D(x1, y1, x2, y2) -> float`
- `distance3D(x1, y1, z1, x2, y2, z2) -> float`
- `angleTo(x1, y1, x2, y2) -> float`
- `lerp(a, b, t) -> float`
- `remap(value, inMin, inMax, outMin, outMax) -> float`
- `chance(percent) -> bool`
- `randomPointInCircle(cx, cy, radius) -> {x, y}`
- `randomPointInSphere(cx, cy, cz, radius) -> {x, y, z}`

## 20) Networking-friendly helpers (optional)

- `netId(entity) -> int`
- `setNetOwner(entity, peerId)`
- `snapshot(entity, fields) -> map`
- `interpolateRemote(entity, snapshotA, snapshotB, alpha)`
- `predict(entity, inputState, dt)`
- `reconcile(entity, authoritativeSnapshot)`

These are intended as optional helpers that layer on top of ENET or future netcode modules.

## Complete usage example

```candy
window(800, 600, "My Game")
setTargetFPS(60)

player = {x: 400, y: 300, speed: 200}
enemies = []
score = 0

cameraBounds(0, 0, 2000, 2000)

playerAnim = animController()
addAnimation(playerAnim, "idle", [0, 1], 0.3, true)
addAnimation(playerAnim, "run", [2, 3, 4], 0.1, true)

explosionEmitter = particles(0, 0)
particleSpeed(explosionEmitter, 100, 300)
particleLife(explosionEmitter, 0.5, 1.0)

while !shouldClose() {
  dt = deltaTime()

  moveX = axis("horizontal")
  moveY = axis("vertical")
  player.x = player.x + moveX * player.speed * dt
  player.y = player.y + moveY * player.speed * dt

  if moveX != 0 || moveY != 0 {
    setAnimation(playerAnim, "run")
  } else {
    setAnimation(playerAnim, "idle")
  }
  updateAnimController(playerAnim, dt)

  cameraFollow(player.x, player.y, 5.0)

  if mousePressed(0) {
    m = screenToWorld(mouseX(), mouseY())
    // user-defined shoot() helper
    shoot(player.x, player.y, m.x, m.y)
    cameraShake(3.0)
  }

  for each enemy in enemies {
    if circleCollision(player.x, player.y, 20, enemy.x, enemy.y, 20) {
      explosion(enemy.x, enemy.y, 30)
      enemies.remove(enemy)
      score = score + 10
    }
  }

  updateParticles(explosionEmitter, dt)

  clear("skyblue")
  sprite(player.x, player.y, "player.png", controllerFrame(playerAnim))

  for each enemy in enemies {
    circle(enemy.x, enemy.y, 20, "red")
  }

  drawParticles(explosionEmitter)
  uiText(10, 10, "Score: {score}", 24, "white")
  flip()
}
```

## Implementation priority

### Phase 1 (Essential)

- Camera follow, shake, bounds
- Input: `keyPressed`, `axis`, mouse helpers
- Basic collision: box and circle
- Timer and cooldown system

### Phase 2 (High priority)

- Animation controller
- Tween system
- Particle emitter
- Audio playback helpers

### Phase 3 (Nice to have)

- Pathfinding
- Spatial grid
- Scene management
- Save/load system

### Phase 4 (Advanced 3D and systems)

- 3D camera orbit/follow/look controls
- 3D rotate/look-at/orbit movement helpers
- Combat/projectile helpers
- AI/state/wave/spawn helper layers

## Notes on naming and compatibility

- This API intentionally uses simple camelCase helper names.
- Existing low-level Raylib wrappers remain valid and can be mixed with these helpers.
- If both low-level and helper APIs are present, helper calls should layer on top of low-level runtime primitives rather than replacing them.
