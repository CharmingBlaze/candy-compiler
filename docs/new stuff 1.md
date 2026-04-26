# Candy Game Dev Features - Implementation Status

This file is now the implementation and documentation tracker for the game-dev
convenience feature set originally proposed in this document.

## Implemented

- Native vector runtime values: `vec2`, `vec3`, `vec4` constructors.
- Vector arithmetic in evaluator: vector/vector and vector/scalar `+ - * /`.
- Vector members and helpers:
  - Properties: `.x`, `.y`, `.z`, `.w`, `.xy`, `.xz`, `.yz`
  - Methods: `.length()`, `.normalize()`, `.dot(v)`, `.cross(v)`, `.distance(v)`
- Import ergonomics:
  - `import module as alias`
  - `from module import a, b`
- Named call arguments:
  - `fun add(a, b) { ... }`
  - `add(b: 2, a: 3)`
- Indexed for-in support:
  - `for item, index in array { ... }`
- Operator overload execution in runtime for struct operators parsed in struct defs.
- Struct property getter/setter execution in runtime for parsed property blocks.
- Formatting and convenience builtins:
  - `format("HP: {} Score: {}", hp, score)`
  - `enumerate(array)`
- Physics helper constructors:
  - `AABB(centerVec3, halfExtentsVec3)`
  - `Sphere(centerVec3, radius)`
  - `Ray(originVec3, directionVec3)`
- Physics helper methods on physics maps:
  - `aabb.overlaps(otherAabb)`
  - `aabb.contains(pointVec3)`
  - `sphere.overlaps(otherSphere)`
  - `ray.intersects(aabb)`
- Game/ECS convenience helper builtins:
  - `gameLoop(updateFn, drawFn)` (minimal helper)
  - `entity()`, `addComponent(entity, name, comp)`, `getComponent(entity, name)`

## Partial

- Object destructuring shorthand (`{x, y} = obj`) parser support is in progress;
  object-pattern assignment semantics need one more stabilization pass.
- Physics helpers currently expose map-backed runtime objects (dynamic) rather
  than dedicated static AST/typechecker-native object types.
- `enumerate()` is implemented as a runtime helper; loop ergonomics can be
  expanded further for first-index naming conventions.

## Notes for Usage

- The vector and physics features are available in the evaluator path now and
  are intended for game scripting ergonomics.
- Existing Mario/game examples can incrementally migrate from scalar component
  math to vector math as needed.
- The language remains backwards compatible with previous forms.
plaintext# CANDY LANGUAGE - GAME DEVELOPMENT CONVENIENCE FEATURES
# Commands to Make Game Programming Easier

Based on the Mario 64 example, here are missing features that would simplify game code:

================================================================================
## CRITICAL MISSING FEATURES FOR GAME DEV
================================================================================

### 1. VECTOR/POINT TYPES (Built-in)

**Problem:** Currently doing math component-by-component
```candy
// Current (verbose)
px = px + vx * dt
py = py + vy * dt
pz = pz + vz * dt

// Better with vector types
position = position + velocity * dt
```

**Solution: Native Vector Types**
```candy
// Built-in vector types
vec2 playerPos = vec2(100, 200)
vec3 position = vec3(0, 2, 0)
vec4 color = vec4(1, 0, 0, 1)

// Vector operations
position = position + velocity * dt
distance = position.distance(target)
direction = position.normalize()
dotProduct = vecA.dot(vecB)
crossProduct = vecA.cross(vecB)

// Component access
x = position.x
y = position.y
z = position.z

// Or swizzling
xy = position.xy  // Returns vec2
xz = position.xz
```

**Add to checklist:**
☐ vec2, vec3, vec4 types
☐ Vector arithmetic (+, -, *, /)
☐ Vector methods: .length(), .normalize(), .dot(), .cross(), .distance()
☐ Component access: .x, .y, .z, .w
☐ Swizzling: .xy, .xz, .yz, etc.


### 2. MATH LIBRARY (Module)

**Problem:** Having to write `math.sin()` everywhere
```candy
// Current
sinY = math.sin(yaw)
cosY = math.cos(yaw)
moveLen = math.sqrt(moveX * moveX + moveZ * moveZ)

// Better with imports
from math import sin, cos, sqrt
sinY = sin(yaw)
cosY = cos(yaw)
moveLen = sqrt(moveX * moveX + moveZ * moveZ)
```

**Solution: Import System**
```candy
// Import entire module
import math

// Import specific functions
from math import sin, cos, tan, sqrt, abs, clamp, lerp

// Import with alias
import math as m
result = m.sin(angle)
```

**Add to checklist:**
☐ `from module import function1, function2`
☐ Import aliases: `import module as alias`


### 3. ENHANCED ARRAY INITIALIZATION

**Problem:** Manually typing out array literals
```candy
// Current (tedious)
platformX = [0.0, 7.0, -7.0, 0.0, 12.0]
platformY = [0.5, 1.25, 1.75, 3.0, 1.0]
platformZ = [0.0, 6.0, -5.5, -10.0, -2.0]

// Better with structured data
platforms = [
  {x: 0.0,  y: 0.5,  z: 0.0,   w: 28.0, h: 1.0,  d: 28.0},
  {x: 7.0,  y: 1.25, z: 6.0,   w: 5.0,  h: 0.5,  d: 5.0},
  {x: -7.0, y: 1.75, z: -5.5,  w: 6.0,  h: 0.5,  d: 6.0},
  {x: 0.0,  y: 3.0,  z: -10.0, w: 10.0, h: 0.75, d: 2.0},
  {x: 12.0, y: 1.0,  z: -2.0,  w: 4.0,  h: 0.5,  d: 4.0}
]

// Access
platform = platforms[i]
drawCube(platform.x, platform.y, platform.z, platform.w, platform.h, platform.d, "gray")
```

**Solution: Array of Objects (already have, but emphasize)**
```candy
// Define platform type
object Platform {
  x = 0.0
  y = 0.0
  z = 0.0
  width = 1.0
  height = 1.0
  depth = 1.0
  
  fun draw(color) {
    drawCube(x, y, z, width, height, depth, color)
  }
}

// Create platforms
platforms = [
  Platform(x: 0, y: 0.5, z: 0, width: 28, height: 1, depth: 28),
  Platform(x: 7, y: 1.25, z: 6, width: 5, height: 0.5, depth: 5)
]

// Use
for platform in platforms {
  platform.draw("gray")
}
```


### 4. NAMED PARAMETERS (For Function Calls)

**Problem:** Long function calls are unclear
```candy
// Current (unclear what numbers mean)
drawCube(platformX[i], platformY[i], platformZ[i], platformW[i], platformH[i], platformD[i], "gray")

// Better with named parameters
drawCube(
  x: platformX[i],
  y: platformY[i],
  z: platformZ[i],
  width: platformW[i],
  height: platformH[i],
  depth: platformD[i],
  color: "gray"
)

// Or mixed positional + named
drawCube(x, y, z, width: w, height: h, depth: d, color: "gray")
```

**Add to checklist:**
☐ Named parameters in function calls
☐ Mixed positional + named parameters


### 5. DESTRUCTURING ASSIGNMENT

**Problem:** Extracting values from objects is verbose
```candy
// Current
platform = platforms[i]
x = platform.x
y = platform.y
z = platform.z
w = platform.width

// Better with destructuring
{x, y, z, width, height, depth} = platforms[i]
drawCube(x, y, z, width, height, depth, "gray")

// Array destructuring (already mentioned)
[r, g, b, a] = color.toArray()
```

**Add to checklist:**
☐ Object destructuring: {x, y} = point
☐ Array destructuring: [a, b, c] = array
☐ Nested destructuring


### 6. ENUM WITH METHODS

**Problem:** Magic numbers everywhere
```candy
// Current
if hp <= 0 {
  lose = true
}

// Better with constants
const MAX_HP = 3
const MIN_HP = 0

if hp <= MIN_HP {
  lose = true
}

// Even better with enum
enum GameState {
  PLAYING,
  WIN,
  LOSE,
  PAUSED
}

gameState = GameState.PLAYING

if gameState == GameState.WIN {
  drawText("You win!", ...)
}
```

**Add to checklist:**
☐ Better enum support (already have basic)


### 7. PROPERTY GETTERS/SETTERS

**Problem:** Can't add computed properties
```candy
// Current - manual calculation
moveLen = math.sqrt(moveX * moveX + moveZ * moveZ)

// Better with getters
object Vector2 {
  x = 0
  y = 0
  
  get length {
    return math.sqrt(x * x + y * y)
  }
  
  set length(newLen) {
    scale = newLen / length
    x = x * scale
    y = y * scale
  }
}

vec = Vector2(x: 3, y: 4)
print(vec.length)  // 5.0 (calculated on access)
vec.length = 10    // Scales vector to length 10
```

**Add to checklist:**
☐ Property getters: `get propertyName { }`
☐ Property setters: `set propertyName(value) { }`


### 8. OPERATOR OVERLOADING (For Vectors)

**Problem:** Can't do math with custom types naturally
```candy
// Current (verbose)
position.x = position.x + velocity.x * dt
position.y = position.y + velocity.y * dt
position.z = position.z + velocity.z * dt

// Better with operator overloading
position = position + velocity * dt
```

**Solution:**
```candy
object Vector3 {
  x = 0.0
  y = 0.0
  z = 0.0
  
  operator +(other) {
    return Vector3(x: x + other.x, y: y + other.y, z: z + other.z)
  }
  
  operator *(scalar) {
    return Vector3(x: x * scalar, y: y * scalar, z: z * scalar)
  }
}

// Now can write
newPos = oldPos + velocity * deltaTime
```

**Add to checklist:**
☐ Operator overloading for custom types
☐ Support: +, -, *, /, ==, !=, <, >, etc.


### 9. PHYSICS HELPERS (Built-in or Stdlib)

**Problem:** Collision code is verbose and error-prone
```candy
// Current AABB collision (11 lines!)
fun aabbOverlap(ax, ay, az, ahx, ahy, ahz, bx, by, bz, bhx, bhy, bhz) {
  if absf(ax - bx) > (ahx + bhx) { return false }
  if absf(ay - by) > (ahy + bhy) { return false }
  if absf(az - bz) > (ahz + bhz) { return false }
  return true
}

// Better with built-in types
fun aabbOverlap(a: AABB, b: AABB) {
  return a.overlaps(b)
}

// Or even simpler
if boxA.overlaps(boxB) {
  // collision!
}
```

**Solution: Built-in Physics Types**
```candy
// AABB (Axis-Aligned Bounding Box)
object AABB {
  center: vec3
  halfExtents: vec3
  
  fun overlaps(other: AABB) {
    // Built-in collision check
  }
  
  fun contains(point: vec3) {
    // Point-in-box test
  }
}

// Sphere
object Sphere {
  center: vec3
  radius: float
  
  fun overlaps(other: Sphere) {
    return center.distance(other.center) < (radius + other.radius)
  }
}

// Ray
object Ray {
  origin: vec3
  direction: vec3
  
  fun intersects(aabb: AABB) {
    // Ray-box intersection
  }
}
```

**Add to checklist:**
☐ AABB type with collision methods
☐ Sphere type
☐ Ray type
☐ Raycast helpers


### 10. CONSTANTS IN OBJECTS

**Problem:** Constants scattered everywhere
```candy
// Current
accel = 22.0
drag = 12.0
maxSpeed = 7.0
jumpPower = 9.3

// Better with config objects
object PlayerConfig {
  const ACCELERATION = 22.0
  const DRAG = 12.0
  const MAX_SPEED = 7.0
  const JUMP_POWER = 9.3
  const GRAVITY = 24.0
}

vx = vx + worldX * PlayerConfig.ACCELERATION * dt
```

**Add to checklist:**
☐ const inside objects/classes


### 11. INCREMENTAL OPERATORS ON ARRAYS

**Problem:** Can't easily modify array elements
```candy
// Current (awkward)
temp = enemyX[i]
temp = temp + enemyVX[i] * dt
enemyX[i] = temp

// Better
enemyX[i] += enemyVX[i] * dt
```

**Already have this, but ensure it works on array elements:**
☐ arr[i] += value
☐ arr[i] -= value
☐ arr[i] *= value
☐ arr[i] /= value


### 12. FOREACH WITH INDEX

**Problem:** Need index in foreach loops
```candy
// Current (need manual counter)
i = 0
while i < enemyCount {
  if enemyAlive[i] {
    // process enemy
  }
  i = i + 1
}

// Better with indexed foreach
for enemy, i in enemies {
  if enemy.alive {
    // i is the index
  }
}

// Or enumerate
for i, enemy in enumerate(enemies) {
  // i is index, enemy is value
}
```

**Add to checklist:**
☐ foreach with index: `for item, index in array`
☐ Or `enumerate()` function


### 13. MULTIPLE ASSIGNMENT FROM FUNCTION

**Already mentioned, but critical for game dev:**
```candy
// Current
camInfo = calculateCamera()
cx = camInfo.x
cy = camInfo.y
cz = camInfo.z

// Better
cx, cy, cz = calculateCamera()
```


### 14. SHORT-CIRCUIT LOGIC (Already have, but document)

```candy
// Make sure && and || short-circuit properly
if (win or lose) and isKeyPressed("r") {
  resetGame()
}

// 'or' should not evaluate second part if first is true
// 'and' should not evaluate second part if first is false
```


### 15. TERNARY WITH CHAINING

```candy
// Current
if !onGround {
  playerColor = "maroon"
} else {
  playerColor = "red"
}

// Better
playerColor = !onGround ? "maroon" : "red"

// Even better with elvis
playerColor = onGround ? "red" : "maroon"
```


### 16. STRING BUILDER / FORMAT

**Problem:** String concatenation in loops is slow
```candy
// Current (inefficient)
text = "HP: " + toString(hp) + "   Score: " + toString(score)

// Better with format
text = format("HP: {}   Score: {}", hp, score)

// Or string interpolation (already have)
text = "HP: {hp}   Score: {score}"
```


### 17. GAME LOOP HELPERS

**Problem:** Common patterns repeated
```candy
// Built-in game loop helper
gameLoop {
  onUpdate(dt) {
    // Update game state
  }
  
  onDraw() {
    beginDrawing()
    // Draw stuff
    endDrawing()
  }
  
  onExit() {
    // Cleanup
  }
}

// Or simpler
while !shouldClose() {
  update(getFrameTime())
  draw()
}
```


### 18. ENTITY COMPONENT SYSTEM (Optional, Advanced)

```candy
// Simple ECS for game entities
object Entity {
  position: vec3
  velocity: vec3
  components: []
  
  fun addComponent(comp) {
    components.add(comp)
  }
  
  fun getComponent(type) {
    return components.find(c => c.type == type)
  }
  
  fun update(dt) {
    for comp in components {
      comp.update(dt)
    }
  }
}

// Components
object PhysicsComponent {
  gravity = 24.0
  
  fun update(entity, dt) {
    entity.velocity.y -= gravity * dt
    entity.position += entity.velocity * dt
  }
}

object RenderComponent {
  color = "red"
  size: vec3
  
  fun draw(entity) {
    drawCube(entity.position, size, color)
  }
}
```


================================================================================
## REVISED PRIORITY LIST (GAME DEV FOCUSED)
================================================================================

### HIGHEST PRIORITY (Game Dev Must-Have):
1. **Vector types (vec2, vec3, vec4)** - Critical for 3D games
2. **Named parameters** - Readability
3. **Import with 'from'** - Cleaner code
4. **Destructuring** - Less boilerplate
5. **Operator overloading** - Natural vector math
6. **foreach with index** - Common pattern
7. **Property getters/setters** - Computed properties

### HIGH PRIORITY:
8. Physics helpers (AABB, Sphere, Ray)
9. Format strings
10. Multiple return values (already mentioned)
11. const in objects
12. Ternary operator (already mentioned)

### MEDIUM PRIORITY:
13. ECS helpers (optional)
14. Game loop patterns
15. Better math module


================================================================================
## UPDATED FINAL CHECKLIST
================================================================================

### NEW CRITICAL ADDITIONS FOR GAME DEV:

☐ **Vector Types**
  - vec2, vec3, vec4
  - Component access (.x, .y, .z, .w)
  - Swizzling (.xy, .xz, etc.)
  - Vector operations (+, -, *, /, dot, cross, length, normalize)

☐ **Named Parameters**
  - Call functions with name: value
  - Mixed positional + named

☐ **From Import**
  - `from module import func1, func2`
  - Import aliases

☐ **Destructuring**
  - Object: `{x, y} = point`
  - Array: `[a, b] = arr`

☐ **Operator Overloading**
  - Define custom operators for types
  - Support: +, -, *, /, ==, !=, etc.

☐ **Foreach with Index**
  - `for item, index in array`
  - Or `enumerate()` function

☐ **Property Getters/Setters**
  - `get property { }`
  - `set property(value) { }`

☐ **Physics Types** (Optional Stdlib)
  - AABB with .overlaps()
  - Sphere with collision
  - Ray with intersect


**Game dev features added: ~35 new items**
**Revised total: ~175 features**

**Implementation time with game dev focus: 15-18 weeks**

🍬 **Candy: Now even sweeter for game development!** 🎮🍬