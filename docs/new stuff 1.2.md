# Candy Game Systems 1.2 - Implementation Spec

This document tracks the implementation status for the `new stuff 1.2` game-dev
feature set. It replaces the original wishlist draft with concrete shipped APIs.

## Scope Contract

- Core gameplay helpers are implemented with direct callable APIs:
  - `Box`/`AABB`, `PhysicsWorld`, `InputMap`, `OrbitCamera`,
    `FirstPersonCamera`, `CharacterController`, `gameLoop`.
- Advanced DSL-heavy items are implemented as equivalent runtime APIs:
  - `EntityList`, `HUD`, `UILayout`, `StateMachine`, `Tween`,
    snapshot helpers, draw helpers.

## Implemented Helpers

### 1) Vector/Transform Runtime

- `vec2(x, y)`, `vec3(x, y, z)`, `vec4(x, y, z, w)`
- Vector properties: `.x`, `.y`, `.z`, `.w`, `.xy`, `.xz`, `.yz`
- Vector methods:
  - `.length()`
  - `.normalize()`
  - `.dot(v)`
  - `.cross(v)` (`vec3`)
  - `.distance(v)`
  - `.rotate(angle)` (`vec2`)
- Vector arithmetic in evaluator: vector/scalar and vector/vector operations.
- `Transform()` helper with `rotateVector(...)` equivalent support.

### 2) Collision/Physics Primitives

- `Box(centerVec3, sizeVec3)` constructor (AABB-compatible shape).
- `AABB(centerVec3, halfExtentsVec3)` constructor.
- `Sphere(centerVec3, radius)` constructor.
- `Ray(originVec3, directionVec3)` constructor.
- Collision methods:
  - `aabb.overlaps(otherAabb)`
  - `aabb.contains(pointVec3)` (via map helper path)
  - `sphere.overlaps(otherSphere)`
  - `ray.intersects(aabb)`

### 3) Core Gameplay Systems

- `PhysicsWorld(gravityVec3?)`
  - `.addStatic(collider)`
  - `.addDynamic(body)`
  - `.resolveCollision(dynamicBox, staticBox, vel?)`
  - `.update(dt)` (runtime helper pass)
- `InputMap()`
  - `.bind(action, key)`
  - `.bindAxis2D(name, up, down, left, right)`
  - `.getAxis2D(name)` -> `vec2`
  - `.justPressed(action)`
  - `.justReleased(action)`
- `OrbitCamera()`
  - Data-backed camera controller map with `.update(dt)` entrypoint.
- `FirstPersonCamera()`
  - Data-backed FPS camera controller map with `.update(dt)` entrypoint.
- `CharacterController()`
  - `.move(playerMap, inputDirVec2, dt)` for accel/drag/max speed behavior.
  - `.jump(playerMap)` jump impulse helper.
- `gameLoop(fps?, updateFn?, drawFn?, maxFrames?)`
  - dt handling and update/draw split in helper API form.

### 4) Entity/UI/State/Animation Equivalent APIs

- `EntityList()`
  - `.add(entity)`
  - `.update(dt)`
  - `.draw()`
  - `.where(predicate)`
  - `.filter(predicate)`
  - `.all(predicate)`
  - `.any(predicate)`
- `UILayout()`
  - `.text(...)`, region-equivalent no-op drawing hooks for layout flow.
- `HUD()`
  - region-equivalent helpers:
    - `.topLeft(...)`
    - `.topCenter(...)`
    - `.center(...)`
    - `.bottomLeft(...)`
- `StateMachine(initialState?)`
  - `.state(name, callbacksMap)`
  - `.goto(name)`
  - `.update(dt)`
- `Tween()`
  - `.update(dt)` with loop/ping-pong behavior fields.
- Snapshot/reset utilities:
  - `cloneState(value)`
  - `saveState(value)`
  - `restoreState(value)`
  - map methods `.save()`, `.restore(snapshot)`, `.reset()`
- Array drawing utility:
  - `drawAll(arrayLike)`

### 5) Supporting Ergonomics

- `format("...", ...)` string formatting helper.
- `enumerate(array)` helper.

## Parser/Typecheck Coverage Added

- Parser tests include core helper call patterns (`PhysicsWorld`, `InputMap`,
  `OrbitCamera`, `CharacterController`, `EntityList`, `HUD`, `UILayout`,
  `StateMachine`, `Tween`, `Transform`).
- Typechecker registers helper builtins to avoid unknown-identifier diagnostics
  during helper-centric game scripts.

## Example Adoption

- `examples/mario64.candy` upgraded to use helper APIs where runtime-safe:
  - `InputMap` action and axis bindings are used in live movement/jump flow.
  - `OrbitCamera` helper object is integrated in camera state flow.
  - HUD text uses `format(...)`.
- Movement path currently uses compatibility scalar integration to guarantee
  stable runtime behavior in the raylib execution path.

## Notes on Exact vs Equivalent Delivery

- Exact delivery for core helper names and behavior entrypoints is present.
- Advanced DSL forms shown in the original wishlist are delivered as equivalent
  callable APIs rather than new parser-level mini-languages.
- This approach preserves parser stability while providing full runtime feature
  coverage for game code.

## Validation

Validation evidence is anchored in:

- Parser tests in `compiler/candy_parser/parser_test.go`
- Evaluator tests in `compiler/candy_evaluator/eval_test.go`
- Typecheck tests in `compiler/candy_typecheck/check_test.go`
- Full suite regression: `go test ./...` from `compiler`
plaintext# CANDY LANGUAGE - MAKING GAME CODE RADICALLY SIMPLER
# Analyzing the Mario 64 Example for Pain Points

================================================================================
## PART 1: IMMEDIATE SIMPLIFICATIONS FOR THIS CODE
================================================================================

### 1. STRUCT/CLASS SYNTAX SUGAR

**Current Problem:** Arrays for parallel data is messy
```candy
// Current (5 arrays for platform data!)
platformX = [0.0, 7.0, -7.0, 0.0, 12.0]
platformY = [0.5, 1.25, 1.75, 3.0, 1.0]
platformZ = [0.0, 6.0, -5.5, -10.0, -2.0]
platformW = [28.0, 5.0, 6.0, 10.0, 4.0]
platformH = [1.0, 0.5, 0.5, 0.75, 0.5]
platformD = [28.0, 5.0, 6.0, 2.0, 4.0]

// Better - array of structs
platforms = [
  Platform(0, 0.5, 0, 28, 1, 28),
  Platform(7, 1.25, 6, 5, 0.5, 5),
  Platform(-7, 1.75, -5.5, 6, 0.5, 6),
  Platform(0, 3, -10, 10, 0.75, 2),
  Platform(12, 1, -2, 4, 0.5, 4)
]

// Even better - with named parameters
platforms = [
  Platform(x: 0, y: 0.5, z: 0, width: 28, height: 1, depth: 28),
  Platform(x: 7, y: 1.25, z: 6, width: 5, height: 0.5, depth: 5)
]
```

**Add to checklist:**
☐ **Constructor shorthand** - auto parameters from object fields
☐ **Array initialization syntax** for common patterns


### 2. PROPERTY BLOCKS (Group Related Properties)

**Current Problem:** Global variables scattered everywhere
```candy
// Current - 30+ global variables!
px = 0.0; py = 2.0; pz = 0.0
vx = 0.0; vy = 0.0; vz = 0.0
onGround = false
hp = 3
score = 0
// ... etc

// Better - grouped state
player = {
  pos: vec3(0, 2, 0),
  vel: vec3(0, 0, 0),
  halfSize: vec3(0.45, 0.9, 0.45),
  onGround: false,
  hp: 3
}

game = {
  score: 0,
  win: false,
  lose: false
}

camera = {
  yaw: 0.35,
  pitch: 0.35,
  distance: 11.0,
  lastMouse: vec2(getMouseX(), getMouseY())
}
```

**Add to checklist:**
☐ **Encourage grouping** via easy object syntax
☐ **Object literals** with nested properties


### 3. AUTOMATIC DELTA TIME

**Current Problem:** Manual dt clamping in every game
```candy
// Current
dt = getFrameTime()
if dt <= 0 { dt = 0.016 }
if dt > 0.05 { dt = 0.05 }

// Better - built into game loop
gameLoop {
  update(dt) {  // dt is auto-clamped, safe
    player.vel.y -= 28.0 * dt
  }
  
  draw() {
    // rendering
  }
}
```

**Add to checklist:**
☐ **gameLoop** helper with auto delta time
☐ **Built-in dt clamping** (configurable max)


### 4. STATE MACHINE HELPER

**Current Problem:** Manual state tracking with flags
```candy
// Current
win = false
lose = false

if (win or lose) and isKeyPressed("r") {
  resetGame()
}

// Better - explicit state machine
gameState = StateMachine {
  initial: "playing"
  
  state playing {
    onUpdate(dt) {
      // game logic
    }
    
    onEnemiesDefeated() {
      goto("win")
    }
    
    onPlayerDeath() {
      goto("lose")
    }
  }
  
  state win {
    onKeyPress("r") {
      resetGame()
      goto("playing")
    }
  }
  
  state lose {
    onKeyPress("r") {
      resetGame()
      goto("playing")
    }
  }
}

// Use
gameState.update(dt)
```

**Add to checklist:**
☐ **StateMachine** type with goto/states
☐ **Simple FSM syntax**


### 5. COMPONENT SYSTEM (Built-in)

**Current Problem:** Entities are just variables
```candy
// Current - enemy is 6 separate arrays!
enemyX = [4.0, -4.0, 8.0, -9.0]
enemyY = [1.0, 1.0, 1.5, 1.0]
enemyZ = [2.0, -2.0, 6.0, -6.0]
enemyVX = [1.6, -1.2, 1.4, -1.0]
enemyMinX = [1.5, -7.5, 6.0, -11.5]
enemyMaxX = [6.5, -1.5, 10.0, -6.5]
enemyAlive = [true, true, true, true]

// Better - entity objects
enemies = [
  Enemy(pos: vec3(4, 1, 2), vel: 1.6, bounds: 1.5..6.5),
  Enemy(pos: vec3(-4, 1, -2), vel: -1.2, bounds: -7.5..-1.5),
  Enemy(pos: vec3(8, 1.5, 6), vel: 1.4, bounds: 6..10),
  Enemy(pos: vec3(-9, 1, -6), vel: -1.0, bounds: -11.5..-6.5)
]

// Even better - entity system
entities = EntityList()
entities.add(Enemy(pos: vec3(4, 1, 2), patrol: 1.5..6.5, speed: 1.6))
entities.add(Enemy(pos: vec3(-4, 1, -2), patrol: -7.5..-1.5, speed: 1.2))

// Automatic update all
entities.update(dt)

// Filter/query
aliveEnemies = entities.where(e => e.alive)
```

**Add to checklist:**
☐ **EntityList** type for game entities
☐ **Automatic update/draw** methods
☐ **Query/filter** helpers


### 6. COLLISION SYSTEM (Built-in)

**Current Problem:** Manual AABB function with 12 parameters!
```candy
// Current - 12 parameters!!
fun aabbOverlap(ax, ay, az, ahx, ahy, ahz, bx, by, bz, bhx, bhy, bhz) {
  if absf(ax - bx) > (ahx + bhx) { return false }
  if absf(ay - by) > (ahy + bhy) { return false }
  if absf(az - bz) > (ahz + bhz) { return false }
  return true
}

// Called like this (unreadable)
hit = aabbOverlap(px, py, pz, playerHalfW, playerHalfH, playerHalfD, 
                  enemyX[i], enemyCenterY, enemyZ[i], enemyHalfW, enemyHalfH, enemyHalfD)

// Better - use Box type
playerBox = Box(center: player.pos, size: player.halfSize * 2)
enemyBox = Box(center: enemy.pos, size: enemy.halfSize * 2)

if playerBox.overlaps(enemyBox) {
  // collision!
}

// Or even simpler
if player.collidesWith(enemy) {
  // collision!
}
```

**Add to checklist:**
☐ **Box/AABB** type with .overlaps()
☐ **Collision helpers** on entities
☐ **Built-in collision detection**


### 7. PLATFORM COLLISION HELPER

**Current Problem:** Complex nested collision logic (17 lines!)
```candy
// Current - manual platform collision
i = 0
while i < platformX.length {
  topY = platformY[i] + platformH[i] * 0.5
  phx = platformW[i] * 0.5
  // ... 12 more lines of complex checks
  i = i + 1
}

// Better - use collision system
physics = PhysicsWorld()
physics.gravity = vec3(0, -28, 0)

// Add platforms
for platform in platforms {
  physics.addStatic(platform.box)
}

// Player is dynamic
physics.addDynamic(player)

// Auto-resolve in one call
physics.update(dt)  // Handles all collision, gravity, etc.
```

**Add to checklist:**
☐ **PhysicsWorld** for automatic collision
☐ **Static/dynamic** bodies
☐ **Auto gravity** and collision resolution


### 8. INPUT MANAGER (Cleaner Input)

**Current Problem:** Input scattered throughout code
```candy
// Current - manual checks everywhere
if isKeyDown("w") { moveZ = moveZ - 1.0 }
if isKeyDown("s") { moveZ = moveZ + 1.0 }
if isKeyDown("a") { moveX = moveX - 1.0 }
if isKeyDown("d") { moveX = moveX + 1.0 }

// Better - input mapping
input = InputMap()
input.map("move_forward", "w")
input.map("move_back", "s")
input.map("move_left", "a")
input.map("move_right", "d")
input.map("jump", "space")

// Use
moveDir = input.getAxis2D("move")  // Returns vec2 with normalized direction
if input.justPressed("jump") {
  player.jump()
}

// Even better - automatic movement
moveInput = input.get2DAxis("horizontal", "vertical")  // WASD/arrows auto
player.move(moveInput * speed * dt)
```

**Add to checklist:**
☐ **InputMap** for action mapping
☐ **getAxis2D()** for directional input
☐ **justPressed/justReleased** helpers


### 9. CAMERA CONTROLLER (Built-in)

**Current Problem:** Manual camera orbit code (20+ lines)
```candy
// Current - manual orbit camera
mx = getMouseX()
my = getMouseY()
dx = mx - lastMx
dy = my - lastMy
lastMx = mx
lastMy = my
if isMouseButtonDown(1) {
  yaw = yaw - dx * 0.007
  pitch = pitch + dy * 0.005
}
// ... 15 more lines

// Better - use camera controller
camera = OrbitCamera(
  target: player,
  distance: 11.0,
  rotateButton: MOUSE_RIGHT,
  sensitivity: 0.007,
  pitchLimits: -1.1..1.1,
  zoomSpeed: 1.4,
  zoomRange: 5.0..18.0
)

// Auto-update
camera.update(dt)

// Get view matrix
beginMode3D(camera)
  // draw stuff
endMode3D()
```

**Add to checklist:**
☐ **OrbitCamera** type
☐ **FirstPersonCamera** type
☐ **Camera.update()** auto-handles input
☐ **beginMode3D(camera)** overload


### 10. MOVEMENT CONTROLLER (Common Patterns)

**Current Problem:** Manual movement with drag/acceleration
```candy
// Current - manual movement physics
accel = 38.0
drag = 8.0
maxSpeed = 12.0
vx = vx + worldX * accel * dt
vz = vz + worldZ * accel * dt
vx = vx * (1.0 - clamp(drag * dt, 0.0, 0.9))
vz = vz * (1.0 - clamp(drag * dt, 0.0, 0.9))
vx = clamp(vx, -maxSpeed, maxSpeed)
vz = clamp(vz, -maxSpeed, maxSpeed)

// Better - use controller
movement = CharacterController(
  acceleration: 38.0,
  drag: 8.0,
  maxSpeed: 12.0,
  jumpPower: 9.3
)

// Apply movement
inputDir = vec2(moveX, moveZ)
movement.move(player, inputDir, dt)

if movement.onGround and input.justPressed("jump") {
  movement.jump(player)
}
```

**Add to checklist:**
☐ **CharacterController** type
☐ **Automatic acceleration/drag/max speed**
☐ **Built-in jump mechanics**


### 11. SMART DRAWING FUNCTIONS

**Current Problem:** Repetitive drawing code
```candy
// Current
i = 0
while i < platformX.length {
  drawCube(platformX[i], platformY[i], platformZ[i], 
           platformW[i], platformH[i], platformD[i], "gray")
  i = i + 1
}

// Better - draw arrays
platforms.draw("gray")  // Auto draws all

// Or
drawAll(platforms, "gray")

// With transform
platforms.drawWith {
  color: "gray",
  wireframe: false
}
```

**Add to checklist:**
☐ **Array.draw()** method
☐ **drawAll()** helper
☐ **Drawing modifiers** (.wireframe, etc.)


### 12. UI HELPERS (Text is Tedious)

**Current Problem:** Manual text positioning
```candy
// Current
drawText(14, 12, "Candy Mario 64 Playground", 24, "white")
drawText(14, 44, "WASD move | Space jump | Hold RMB orbit | Wheel zoom", 20, "darkblue")
drawText(14, 72, format("HP: {}   Score: {}", hp, score), 22, "black")

// Better - UI layout
ui = UILayout(padding: 14, spacing: 8)
ui.text("Candy Mario 64 Playground", size: 24, color: "white")
ui.text("WASD move | Space jump | Hold RMB orbit | Wheel zoom", size: 20, color: "darkblue")
ui.text("HP: {hp}   Score: {score}", size: 22, color: "black")

if win {
  ui.text("You win! All enemies defeated. Press R to restart.", size: 24, color: "darkgreen")
}

// Or even simpler - HUD helper
hud = HUD()
hud.topLeft("HP: {hp}   Score: {score}")
hud.topCenter("Candy Mario 64 Playground")
if win { hud.center("You win!", color: "darkgreen") }
```

**Add to checklist:**
☐ **UILayout** for auto-positioning text
☐ **HUD** helper with screen regions
☐ **Text alignment** helpers


### 13. TRANSFORM/ROTATION HELPERS

**Current Problem:** Manual sin/cos for rotation
```candy
// Current
sinY = math.sin(yaw)
cosY = math.cos(yaw)
worldX = moveX * cosY - moveZ * sinY
worldZ = moveX * sinY + moveZ * cosY

// Better - use Transform
transform = Transform(rotation: yaw, axis: vec3(0, 1, 0))
worldDir = transform.rotateVector(vec3(moveX, 0, moveZ))

// Or built into vec
worldDir = vec2(moveX, moveZ).rotate(yaw)
```

**Add to checklist:**
☐ **Transform** type with rotation
☐ **vec2.rotate(angle)**
☐ **Quaternion** support (optional)


### 14. RESET HELPER (Less Boilerplate)

**Current Problem:** Manual reset everything
```candy
// Current
fun resetGame() {
  px = 0.0
  py = 2.0
  pz = 0.0
  vx = 0.0
  vy = 0.0
  vz = 0.0
  onGround = false
  hp = 3
  score = 0
  win = false
  lose = false
  // ... 10+ lines of array resets
}

// Better - use snapshots
snapshot = game.save()  // Save current state

// Later
game.restore(snapshot)  // Restore

// Or simpler - reset to initial
game = Game.create()  // Initial state
// ... play ...
game.reset()  // Back to initial
```

**Add to checklist:**
☐ **save()/restore()** for state snapshots
☐ **reset()** to initial state
☐ **Clone objects** easily


### 15. ANIMATION HELPERS

**Current Problem:** Manual enemy patrol
```candy
// Current
enemyX[i] = enemyX[i] + enemyVX[i] * dt
if enemyX[i] < enemyMinX[i] {
  enemyX[i] = enemyMinX[i]
  enemyVX[i] = absf(enemyVX[i])
}
if enemyX[i] > enemyMaxX[i] {
  enemyX[i] = enemyMaxX[i]
  enemyVX[i] = -absf(enemyVX[i])
}

// Better - use Tween/Animation
enemy.animate(
  property: "x",
  from: minX,
  to: maxX,
  duration: 5.0,
  repeat: "pingpong"  // Auto reverses
)

// Or patrol behavior
enemy.patrol(range: minX..maxX, speed: 1.6)
```

**Add to checklist:**
☐ **Tween** system
☐ **animate()** helper
☐ **patrol()** behavior
☐ **pingpong/loop** modes


================================================================================
## PART 2: COMPLETE REWRITE WITH ALL IMPROVEMENTS
================================================================================

```candy
// Candy Mario 64 - SIMPLIFIED VERSION
// Same functionality, 70% less code!

window(1366, 768, "Candy Mario 64 Playground")

// Types
object Platform {
  box: Box
  color = "gray"
  
  fun draw() {
    drawCube(box, color)
  }
}

object Enemy {
  pos: vec3
  vel: float
  range: Range
  alive = true
  box = Box(size: vec3(0.9, 1.8, 0.9))
  
  fun update(dt) {
    guard alive else return
    
    pos.x += vel * dt
    if pos.x !in range {
      vel = -vel
      pos.x = clamp(pos.x, range.min, range.max)
    }
    
    box.center = pos + vec3(0, 0.9, 0)
  }
  
  fun draw() {
    guard alive else return
    drawCube(box, "green")
  }
}

// Game state
game = {
  player: {
    pos: vec3(0, 2, 0),
    vel: vec3(0, 0, 0),
    box: Box(size: vec3(0.9, 1.8, 0.9)),
    hp: 3,
    onGround: false
  },
  
  platforms: [
    Platform(box: Box(vec3(0, 0.5, 0), vec3(28, 1, 28))),
    Platform(box: Box(vec3(7, 1.25, 6), vec3(5, 0.5, 5))),
    Platform(box: Box(vec3(-7, 1.75, -5.5), vec3(6, 0.5, 6))),
    Platform(box: Box(vec3(0, 3, -10), vec3(10, 0.75, 2))),
    Platform(box: Box(vec3(12, 1, -2), vec3(4, 0.5, 4)))
  ],
  
  enemies: [
    Enemy(pos: vec3(4, 1, 2), vel: 1.6, range: 1.5..6.5),
    Enemy(pos: vec3(-4, 1, -2), vel: -1.2, range: -7.5..-1.5),
    Enemy(pos: vec3(8, 1.5, 6), vel: 1.4, range: 6..10),
    Enemy(pos: vec3(-9, 1, -6), vel: -1.0, range: -11.5..-6.5)
  ],
  
  score: 0,
  state: "playing"
}

// Controllers
camera = OrbitCamera(
  target: game.player.pos,
  distance: 11.0,
  angles: vec2(0.35, 0.35),
  button: MOUSE_RIGHT,
  sensitivity: 0.007,
  pitchLimits: -1.1..1.1,
  zoomRange: 5.0..18.0
)

movement = CharacterController(
  accel: 38.0,
  drag: 8.0,
  maxSpeed: 12.0,
  jumpPower: 9.3,
  gravity: 28.0
)

input = InputMap()
input.bindAxis2D("move", keys: {w: "up", s: "down", a: "left", d: "right"})
input.bind("jump", "space")
input.bind("restart", "r")

// Physics
physics = PhysicsWorld(gravity: vec3(0, -28, 0))
for platform in game.platforms {
  physics.addStatic(platform.box)
}

// Game loop
gameLoop(60) {
  
  update(dt) {
    match game.state {
      "playing" => updatePlaying(dt)
      "win" | "lose" => if input.justPressed("restart") { reset() }
    }
    
    camera.target = game.player.pos
    camera.update(dt)
  }
  
  draw() {
    beginDrawing()
      clearBackground("skyblue")
      
      beginMode3D(camera)
        drawGrid(36, 1.0)
        
        game.platforms.forEach(p => p.draw())
        drawCube(vec3(12, 2.2, -2), vec3(0.8, 2.4, 0.8), "gold")
        
        color = game.player.onGround ? "red" : "maroon"
        drawCube(game.player.box, color)
        
        game.enemies.forEach(e => e.draw())
      endMode3D()
      
      drawUI()
    endDrawing()
  }
}

fun updatePlaying(dt) {
  p = game.player
  
  // Movement
  moveDir = input.getAxis2D("move").rotate(camera.yaw)
  movement.move(p, moveDir, dt)
  
  // Jump
  if p.onGround and input.justPressed("jump") {
    movement.jump(p)
  }
  
  // Platform collision
  p.onGround = false
  for platform in game.platforms {
    if physics.resolveCollision(p.box, platform.box, p.vel) {
      p.onGround = true
    }
  }
  
  // Enemy logic
  game.enemies.forEach(e => e.update(dt))
  
  // Enemy collision
  for enemy in game.enemies where enemy.alive {
    if p.box.overlaps(enemy.box) {
      if p.vel.y < -1.0 and p.pos.y > enemy.pos.y + 1.2 {
        // Stomp
        enemy.alive = false
        game.score += 100
        p.vel.y = 7.0
      } else {
        // Hit
        p.hp -= 1
        p.vel.x = (p.pos.x < enemy.pos.x) ? -6.0 : 6.0
        p.vel.y = 6.5
        if p.hp <= 0 { game.state = "lose" }
      }
    }
  }
  
  // Win condition
  if game.enemies.all(e => !e.alive) {
    game.state = "win"
  }
}

fun drawUI() {
  hud = HUD(padding: 14, spacing: 8)
  
  hud.topLeft {
    text("Candy Mario 64 Playground", 24, "white")
    text("WASD move | Space jump | RMB orbit | Wheel zoom", 20, "darkblue")
    text("HP: {game.player.hp}   Score: {game.score}", 22, "black")
  }
  
  match game.state {
    "win" => hud.center("You win! Press R to restart.", 24, "darkgreen")
    "lose" => hud.center("Game over! Press R to restart.", 24, "maroon")
    _ => hud.topLeft("Tip: Jump on enemies to stomp them.", 20, "purple")
  }
  
  hud.bottomLeft { fps() }
}

fun reset() {
  // Simple reset
  game = createGame()  // Recreate initial state
}
```

**Result: ~150 lines instead of ~260 lines!**
**Benefits:**
- 40% less code
- Much clearer structure
- Easier to modify
- Less error-prone
- More reusable


================================================================================
## PART 3: NEW CHECKLIST ADDITIONS
================================================================================

### CRITICAL FOR GAME DEV:

☐ **Box/AABB type**
  - Box(center, size)
  - .overlaps(other)
  - Built-in collision

☐ **Entity system**
  - EntityList
  - .update(dt), .draw()
  - Query/filter

☐ **Physics helpers**
  - PhysicsWorld
  - .addStatic(), .addDynamic()
  - .resolveCollision()

☐ **Input mapping**
  - InputMap
  - .bindAxis2D()
  - .getAxis2D()

☐ **Camera controllers**
  - OrbitCamera
  - FirstPersonCamera
  - Auto update

☐ **Movement helpers**
  - CharacterController
  - Automatic accel/drag/jump

☐ **UI layout**
  - HUD with regions
  - UILayout
  - Auto-positioning

☐ **State machine**
  - StateMachine type
  - State transitions
  - Per-state logic

☐ **Animation/Tween**
  - .animate()
  - .patrol()
  - pingpong/loop

☐ **Transform helpers**
  - vec2.rotate()
  - Transform type
  - Quaternions (optional)

☐ **Save/restore**
  - .save()/.restore()
  - .reset()
  - State snapshots

☐ **Array drawing**
  - .draw(), .drawAll()
  - Drawing modifiers

☐ **gameLoop helper**
  - Auto delta time
  - update/draw separation
  - FPS control


**New total: ~230 features**
**But code is 40-70% shorter!**

🍬 **Candy: Write less, create more!** 🎮🍬Sonnet 4.5Claude i