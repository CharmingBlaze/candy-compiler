# CANDY LANGUAGE SPECIFICATION
# Ultra-Simple Syntax - Easy as Eating Candy!

================================================================================
## PHILOSOPHY: IF A 12-YEAR-OLD CAN'T UNDERSTAND IT, IT'S TOO COMPLEX
================================================================================


## VARIABLES - NO TYPES NEEDED!

x = 10
name = "Sugar"
speed = 5.5
active = true

// That's it! Candy figures out the type automatically.


## MATH - JUST LIKE SCHOOL

result = 5 + 3
difference = 10 - 4
product = 6 * 7
quotient = 20 / 4

x = x + 1  // Add one
x++        // Even easier!

y = y - 1  // Subtract one
y--        // Even easier!


## PRINT STUFF

print("Hello!")           // Print text
print(42)                // Print number
print("Score: " + score) // Print with variables


## GET INPUT

name = input("What's your name? ")
age = input("How old are you? ")


## IF STATEMENTS - SUPER SIMPLE

if x > 5 {
  print("Big number!")
}

if score == 100 {
  print("Perfect!")
} else {
  print("Keep trying!")
}

if lives > 0 {
  print("Still alive!")
} else if lives == 0 {
  print("Game over!")
} else {
  print("Huh?")
}


## LOOPS - EASY PEASY

// Count to 10
for i = 0 to 10 {
  print(i)
}

// Count backwards
for i = 10 to 0 step -1 {
  print(i)
}

// Repeat something 5 times
repeat 5 {
  print("Yay!")
}

// While loop
while playing {
  updateGame()
}

// Loop through a list
candies = ["Skittles", "M&Ms", "Twix"]
for candy in candies {
  print(candy)
}


## LISTS (ARRAYS) - SIMPLE!

// Make a list
scores = [10, 20, 30, 40, 50]
names = ["Alice", "Bob", "Charlie"]

// Get item from list
print(scores[0])  // First item: 10
print(names[2])   // Third item: Charlie

// Change item
scores[0] = 100

// Add to list
scores.add(60)

// Remove from list
scores.remove(2)  // Remove third item

// How many items?
print(scores.count)


## FUNCTIONS - EASY TO MAKE

fun sayHi() {
  print("Hello!")
}

fun greet(name) {
  print("Hello " + name)
}

fun add(a, b) {
  return a + b
}

// Use them
sayHi()
greet("Candy")
result = add(5, 3)


## RANDOM NUMBERS

dice = random(1, 6)        // Random from 1 to 6
coin = random(0, 1)        // 0 or 1
percent = random(0, 100)   // 0 to 100


## SIMPLE GRAPHICS (NO SETUP NEEDED!)

// Open a window - that's it!
window(800, 600, "My Game")

// Main game loop - runs forever
loop {
  
  // Clear screen
  clear(WHITE)
  
  // Draw stuff
  circle(400, 300, 50, RED)
  box(100, 100, 200, 100, BLUE)
  text("Hello!", 100, 50, 20, BLACK)
  
  // Show everything
  show()
}


## COMPLETE SIMPLE EXAMPLES


### Example 1: Moving Ball

window(800, 600, "Moving Ball")

x = 400
y = 300

loop {
  
  // Move with arrow keys
  if key(LEFT) {
    x = x - 5
  }
  if key(RIGHT) {
    x = x + 5
  }
  if key(UP) {
    y = y - 5
  }
  if key(DOWN) {
    y = y + 5
  }
  
  // Draw
  clear(WHITE)
  circle(x, y, 25, RED)
  show()
}


### Example 2: Clicker Game

window(800, 600, "Cookie Clicker")

score = 0

loop {
  
  // Click to score
  if clicked() {
    score = score + 1
  }
  
  // Draw
  clear(SKYBLUE)
  circle(400, 300, 100, BROWN)
  text("Score: " + score, 300, 100, 30, BLACK)
  text("Click the cookie!", 280, 500, 20, GRAY)
  show()
}


### Example 3: Catch Game

window(800, 600, "Catch the Candy")

playerX = 400
candyX = random(0, 800)
candyY = 0
score = 0

loop {
  
  // Move player
  if key(LEFT) {
    playerX = playerX - 5
  }
  if key(RIGHT) {
    playerX = playerX + 5
  }
  
  // Move candy down
  candyY = candyY + 3
  
  // Check if caught
  if candyY > 550 and candyX > playerX - 30 and candyX < playerX + 30 {
    score = score + 1
    candyY = 0
    candyX = random(0, 800)
  }
  
  // Reset if missed
  if candyY > 600 {
    candyY = 0
    candyX = random(0, 800)
  }
  
  // Draw
  clear(SKYBLUE)
  box(playerX - 30, 550, 60, 20, BLUE)
  circle(candyX, candyY, 10, PINK)
  text("Score: " + score, 10, 10, 20, BLACK)
  show()
}


### Example 4: Simple Animation

window(800, 600, "Bouncing Ball")

x = 400
y = 300
speedX = 5
speedY = 3

loop {
  
  // Move ball
  x = x + speedX
  y = y + speedY
  
  // Bounce off walls
  if x < 0 or x > 800 {
    speedX = speedX * -1
  }
  if y < 0 or y > 600 {
    speedY = speedY * -1
  }
  
  // Draw
  clear(WHITE)
  circle(x, y, 20, RED)
  show()
}


## SUPER SIMPLE REFERENCE


### Drawing Commands

circle(x, y, radius, color)
box(x, y, width, height, color)
line(x1, y1, x2, y2, color)
text(message, x, y, size, color)
image(filename, x, y)


### Colors (Built-in)

RED, BLUE, GREEN, YELLOW, ORANGE, PURPLE, PINK
WHITE, BLACK, GRAY
SKYBLUE, BROWN, GOLD


### Input Commands

key(LEFT)       // Is left arrow pressed?
key(RIGHT)      // Is right arrow pressed?
key(UP)         // Is up arrow pressed?
key(DOWN)       // Is down arrow pressed?
key(SPACE)      // Is space pressed?
clicked()       // Was mouse clicked?
mouseX()        // Mouse X position
mouseY()        // Mouse Y position


### Sounds (Simple!)

play("sound.wav")       // Play a sound
music("song.mp3")       // Play background music
stopMusic()             // Stop music


### Math Functions

random(min, max)        // Random number
sqrt(x)                 // Square root
abs(x)                  // Absolute value
sin(x), cos(x), tan(x)  // Trig functions


### String Functions

upper("hello")          // "HELLO"
lower("HELLO")          // "hello"
length("candy")         // 5


### Utilities

wait(1.0)              // Wait 1 second
fps(60)                // Set frames per second
exit()                 // Close program


================================================================================
## SIMPLIFIED COMPILER DESIGN - TRANSPILE TO C!
================================================================================

The simplest approach: Convert Candy → C → Executable


### Translation Table

| Candy Code | C Code |
|------------|--------|
| `x = 10`   | `int x = 10;` |
| `print("Hi")` | `printf("Hi\n");` |
| `print(x)` | `printf("%d\n", x);` |
| `if x > 5 {` | `if (x > 5) {` |
| `for i = 0 to 10 {` | `for (int i = 0; i <= 10; i++) {` |
| `repeat 5 {` | `for (int _i = 0; _i < 5; _i++) {` |
| `loop {` | `while (1) {` |
| `window(800,600,"Game")` | `InitWindow(800,600,"Game"); SetTargetFPS(60);` |
| `clear(WHITE)` | `BeginDrawing(); ClearBackground(WHITE);` |
| `show()` | `EndDrawing();` |
| `circle(x,y,r,c)` | `DrawCircle(x, y, r, c);` |
| `box(x,y,w,h,c)` | `DrawRectangle(x, y, w, h, c);` |
| `key(LEFT)` | `IsKeyDown(KEY_LEFT)` |
| `clicked()` | `IsMouseButtonPressed(MOUSE_LEFT_BUTTON)` |
| `random(1,10)` | `GetRandomValue(1, 10)` |


### Simple Transpiler (Go)

```go
package main

import (
    "fmt"
    "regexp"
    "strings"
)

func TranspileToC(candy string) string {
    c := "#include <stdio.h>\n#include \"raylib.h\"\n\nint main() {\n"
    
    // Simple pattern replacements
    candy = strings.ReplaceAll(candy, "print(", "printf(")
    candy = strings.ReplaceAll(candy, "window(", "InitWindow(")
    candy = strings.ReplaceAll(candy, "loop {", "while (!WindowShouldClose()) {")
    candy = strings.ReplaceAll(candy, "clear(", "BeginDrawing(); ClearBackground(")
    candy = strings.ReplaceAll(candy, "show()", "EndDrawing()")
    candy = strings.ReplaceAll(candy, "circle(", "DrawCircle(")
    candy = strings.ReplaceAll(candy, "box(", "DrawRectangle(")
    candy = strings.ReplaceAll(candy, "text(", "DrawText(")
    candy = strings.ReplaceAll(candy, "key(LEFT)", "IsKeyDown(KEY_LEFT)")
    candy = strings.ReplaceAll(candy, "key(RIGHT)", "IsKeyDown(KEY_RIGHT)")
    candy = strings.ReplaceAll(candy, "key(UP)", "IsKeyDown(KEY_UP)")
    candy = strings.ReplaceAll(candy, "key(DOWN)", "IsKeyDown(KEY_DOWN)")
    candy = strings.ReplaceAll(candy, "key(SPACE)", "IsKeyDown(KEY_SPACE)")
    candy = strings.ReplaceAll(candy, "clicked()", "IsMouseButtonPressed(MOUSE_LEFT_BUTTON)")
    candy = strings.ReplaceAll(candy, "random(", "GetRandomValue(")
    
    // Handle for loops: for i = 0 to 10
    forPattern := regexp.MustCompile(`for (\w+) = (\d+) to (\d+)`)
    candy = forPattern.ReplaceAllString(candy, "for (int $1 = $2; $1 <= $3; $1++)")
    
    // Handle repeat loops: repeat 5
    repeatPattern := regexp.MustCompile(`repeat (\d+)`)
    candy = repeatPattern.ReplaceAllString(candy, "for (int _i = 0; _i < $1; _i++)")
    
    c += candy
    c += "\n    CloseWindow();\n    return 0;\n}"
    
    return c
}
```


### Build Process

```bash
# 1. Candy to C
./candyc game.candy > game.c

# 2. C to executable
gcc game.c -lraylib -o game

# 3. Run!
./game
```


### Example: Full Translation

**Candy:**
```candy
window(800, 600, "Catch Game")

playerX = 400
score = 0

loop {
  if key(LEFT) {
    playerX = playerX - 5
  }
  if key(RIGHT) {
    playerX = playerX + 5
  }
  
  clear(WHITE)
  box(playerX - 25, 500, 50, 20, BLUE)
  text("Score: " + score, 10, 10, 20, BLACK)
  show()
}
```

**Generated C:**
```c
#include "raylib.h"

int main() {
    InitWindow(800, 600, "Catch Game");
    SetTargetFPS(60);
    
    int playerX = 400;
    int score = 0;
    
    while (!WindowShouldClose()) {
        if (IsKeyDown(KEY_LEFT)) {
            playerX = playerX - 5;
        }
        if (IsKeyDown(KEY_RIGHT)) {
            playerX = playerX + 5;
        }
        
        BeginDrawing();
        ClearBackground(WHITE);
        DrawRectangle(playerX - 25, 500, 50, 20, BLUE);
        DrawText(TextFormat("Score: %d", score), 10, 10, 20, BLACK);
        EndDrawing();
    }
    
    CloseWindow();
    return 0;
}
```


## KEY SIMPLIFICATIONS

1. **No type declarations** - Candy figures it out
2. **No semicolons** - Not needed
3. **Simple keywords** - `loop` instead of `while (!WindowShouldClose())`
4. **Easy graphics** - `circle()` instead of `DrawCircle()`
5. **Built-in colors** - Just say `RED` or `BLUE`
6. **Simple input** - `key(LEFT)` instead of `IsKeyDown(KEY_LEFT)`
7. **Auto-include** - Don't worry about headers
8. **No main()** - Your code IS the main
9. **`for i = 0 to 10`** - Easier than C-style for loops
10. **`repeat 5`** - Super simple repetition

**Result: A language a kid can learn in 30 minutes!** 🍬



# CANDY LANGUAGE SPECIFICATION - EXTENDED EDITION
# Ultra-Simple Syntax - Easy as Eating Candy!

================================================================================
## ADDITIONS TO MAKE CANDY EVEN MORE AWESOME!
================================================================================


## TIMERS & TIME

// Get elapsed time
time = seconds()           // Seconds since program started
delta = deltaTime()        // Time since last frame (for smooth movement)

// Wait/delay
wait(2.5)                 // Pause for 2.5 seconds

// Every X seconds, do something
every 2 {
  spawnEnemy()            // Runs every 2 seconds
}

// Countdown timer
timer = 60                // 60 seconds
loop {
  timer = timer - deltaTime()
  if timer <= 0 {
    print("Time's up!")
    break
  }
}


## COLLISION DETECTION (SUPER EASY!)

// Circle collision
if touching(x1, y1, radius1, x2, y2, radius2) {
  print("Hit!")
}

// Box collision
if boxHit(x1, y1, w1, h1, x2, y2, w2, h2) {
  print("Collision!")
}

// Point in box
if inside(mouseX(), mouseY(), x, y, width, height) {
  print("Mouse over box!")
}


## SIMPLE OBJECTS/ENTITIES

// Define a simple object type
object Player {
  x = 400
  y = 300
  speed = 5
  health = 100
  
  fun move() {
    if key(LEFT) {
      x = x - speed
    }
    if key(RIGHT) {
      x = x + speed
    }
  }
  
  fun draw() {
    circle(x, y, 20, BLUE)
  }
}

// Create and use it
player = Player()
player.move()
player.draw()


## LISTS WITH MORE POWER

enemies = []

// Add items
enemies.add(newEnemy)

// Remove items
enemies.remove(0)           // Remove first
enemies.removeLast()        // Remove last
enemies.clear()             // Remove all

// Check if empty
if enemies.empty() {
  print("No enemies!")
}

// Loop through and modify
for enemy in enemies {
  enemy.x = enemy.x + 1
}

// Find and remove
for i = 0 to enemies.count - 1 {
  if enemies[i].dead {
    enemies.remove(i)
  }
}


## SPRITES & IMAGES (EASY!)

// Load image
player = sprite("player.png")

// Draw sprite
draw(player, x, y)

// Draw with size
draw(player, x, y, width, height)

// Rotate sprite
drawRotated(player, x, y, angle)

// Flip sprite
drawFlipped(player, x, y, flipX, flipY)

// Get sprite size
w = player.width
h = player.height

// Unload when done
unload(player)


## ANIMATION (SPRITE SHEETS)

// Load sprite sheet
hero = spriteSheet("hero.png", 32, 32)  // 32x32 frames

// Set current frame
hero.frame = 5

// Animate
hero.frame = hero.frame + 1
if hero.frame > 7 {
  hero.frame = 0
}

// Simple animation helper
anim = animation(hero, 0, 7, 0.1)  // Frames 0-7, 0.1 sec per frame
anim.play()                         // Play animation


## PARTICLES & EFFECTS

// Create particles
for i = 0 to 20 {
  particle(x, y, random(-5, 5), random(-5, 5), RED, 2.0)
}

// Screen shake
shake(0.5, 10)  // Duration 0.5 sec, intensity 10

// Flash screen
flash(WHITE, 0.2)  // Flash white for 0.2 seconds

// Fade in/out
fadeIn(1.0)    // Fade in over 1 second
fadeOut(1.0)   // Fade out over 1 second


## CAMERA (2D)

// Set camera position
camera(x, y)

// Follow player smoothly
cameraFollow(player.x, player.y, 0.1)  // 0.1 = smoothness

// Zoom
zoom(2.0)      // 2x zoom
zoom(0.5)      // Zoom out

// Shake camera
cameraShake(0.3, 5)  // Shake for 0.3 sec


## TILES & GRID

// Create tile map (2D array)
map = [
  [1, 1, 1, 1, 1],
  [1, 0, 0, 0, 1],
  [1, 0, 2, 0, 1],
  [1, 0, 0, 0, 1],
  [1, 1, 1, 1, 1]
]

// Draw tile map
tileSize = 32
for row = 0 to map.count - 1 {
  for col = 0 to map[row].count - 1 {
    tile = map[row][col]
    if tile == 1 {
      box(col * tileSize, row * tileSize, tileSize, tileSize, GRAY)
    }
    if tile == 2 {
      circle(col * tileSize + 16, row * tileSize + 16, 8, GOLD)
    }
  }
}

// Grid to pixel
px = gridToPixel(col, tileSize)
py = gridToPixel(row, tileSize)

// Pixel to grid
col = pixelToGrid(px, tileSize)
row = pixelToGrid(py, tileSize)


## SAVE & LOAD DATA

// Save high score
save("highscore", 1000)

// Load high score
highscore = load("highscore", 0)  // 0 is default if not found

// Save multiple values
save("player_x", player.x)
save("player_y", player.y)
save("player_health", player.health)

// Load multiple values
player.x = load("player_x", 400)
player.y = load("player_y", 300)
player.health = load("player_health", 100)


## SCENES/STATES

// Define scenes
scene = "menu"

loop {
  
  if scene == "menu" {
    // Menu code
    clear(BLACK)
    text("PRESS SPACE TO START", 200, 300, 30, WHITE)
    
    if key(SPACE) {
      scene = "game"
    }
  }
  
  if scene == "game" {
    // Game code
    updateGame()
    drawGame()
    
    if gameOver {
      scene = "gameover"
    }
  }
  
  if scene == "gameover" {
    // Game over code
    clear(RED)
    text("GAME OVER", 300, 300, 40, WHITE)
    text("SCORE: " + score, 320, 350, 25, WHITE)
    
    if key(SPACE) {
      scene = "menu"
      resetGame()
    }
  }
  
  show()
}


## BETTER STRINGS

// Join strings
message = join("Hello", " ", "World")  // "Hello World"

// Split strings
words = split("apple,banana,cherry", ",")  // ["apple", "banana", "cherry"]

// Contains
if contains("hello world", "world") {
  print("Found it!")
}

// Replace
fixed = replace("I like cats", "cats", "dogs")  // "I like dogs"

// To number
num = toNumber("123")     // 123
num = toNumber("45.5")    // 45.5


## BETTER MATH

// Clamp value
x = clamp(x, 0, 800)       // Keep x between 0 and 800

// Lerp (smooth movement)
x = lerp(x, targetX, 0.1)  // Move x toward targetX smoothly

// Distance
dist = distance(x1, y1, x2, y2)

// Angle between points
angle = angleBetween(x1, y1, x2, y2)

// Move toward point
speed = 5
x = x + cos(angle) * speed
y = y + sin(angle) * speed

// Round/floor/ceil
rounded = round(3.7)    // 4
down = floor(3.7)       // 3
up = ceil(3.2)          // 4

// Min/max
smaller = min(10, 5)    // 5
bigger = max(10, 5)     // 10


## DEBUGGING HELPERS

// Debug print (shows in corner)
debug("Player X: " + player.x)
debug("FPS: " + getFPS())

// Draw debug info
debugBox(x, y, width, height)     // Draw box outline
debugPoint(x, y, RED)              // Draw debug point
debugLine(x1, y1, x2, y2, GREEN)  // Draw debug line

// Pause game
if key(P) {
  pause()
}


## 3D BASICS (SIMPLE!)

// 3D mode
window3D(800, 600, "3D Game")

// Set up camera
cam = camera3D(10, 10, 10)  // Position
cam.lookAt(0, 0, 0)         // Look at center

loop {
  clear(SKYBLUE)
  
  start3D(cam)
  
  // Draw 3D stuff
  cube(0, 0, 0, 2, RED)
  sphere(3, 0, 0, 1, BLUE)
  floor(10, GREEN)
  
  end3D()
  
  // Draw 2D UI on top
  text("3D Game!", 10, 10, 20, BLACK)
  
  show()
}


## COMPLETE ADVANCED EXAMPLE

```candy
// Platformer Game with all the features!

window(800, 600, "Super Candy Platformer")

// Load assets
playerSprite = sprite("player.png")
coinSound = loadSound("coin.wav")

// Player
object Player {
  x = 100
  y = 400
  vx = 0
  vy = 0
  speed = 5
  jumpPower = 12
  onGround = false
  
  fun update() {
    // Movement
    if key(LEFT) {
      vx = -speed
    }
    if key(RIGHT) {
      vx = speed
    }
    
    // Jump
    if key(SPACE) and onGround {
      vy = -jumpPower
      onGround = false
    }
    
    // Gravity
    vy = vy + 0.5
    
    // Apply velocity
    x = x + vx
    y = y + vy
    
    // Friction
    vx = vx * 0.8
    
    // Check ground
    if y > 500 {
      y = 500
      vy = 0
      onGround = true
    }
    
    // Clamp to screen
    x = clamp(x, 0, 800)
  }
  
  fun draw() {
    draw(playerSprite, x - 16, y - 16, 32, 32)
  }
}

// Coin
object Coin {
  x = 0
  y = 0
  collected = false
  
  fun draw() {
    if not collected {
      circle(x, y, 10, GOLD)
    }
  }
}

// Create game objects
player = Player()
coins = []
score = 0

// Spawn coins
for i = 0 to 10 {
  coin = Coin()
  coin.x = random(50, 750)
  coin.y = random(100, 400)
  coins.add(coin)
}

// Main game loop
loop {
  
  // Update
  player.update()
  
  // Check coin collection
  for coin in coins {
    if not coin.collected {
      if distance(player.x, player.y, coin.x, coin.y) < 20 {
        coin.collected = true
        score = score + 1
        play(coinSound)
        shake(0.1, 3)
      }
    }
  }
  
  // Draw
  clear(SKYBLUE)
  
  // Draw ground
  box(0, 500, 800, 100, GREEN)
  
  // Draw coins
  for coin in coins {
    coin.draw()
  }
  
  // Draw player
  player.draw()
  
  // Draw UI
  text("Score: " + score, 10, 10, 30, BLACK)
  text("FPS: " + getFPS(), 10, 40, 20, GRAY)
  
  show()
}
```


## FULL FEATURE LIST

### Core Features
✓ Simple variables (no types needed)
✓ Math operators (+, -, *, /, %)
✓ If/else statements
✓ Loops (for, while, repeat, foreach)
✓ Functions (fun)
✓ Objects (simple classes)
✓ Arrays/Lists

### Graphics
✓ Window creation
✓ Basic shapes (circle, box, line)
✓ Text rendering
✓ Colors (built-in constants)
✓ Sprites & images
✓ Sprite sheets & animation
✓ Camera (follow, zoom, shake)
✓ Particles & effects
✓ 3D support (basic)

### Input
✓ Keyboard (key states)
✓ Mouse (position, clicks)
✓ Gamepad support (optional)

### Audio
✓ Sound effects
✓ Music playback
✓ Volume control

### Utilities
✓ Timers & deltaTime
✓ Random numbers
✓ Collision detection
✓ Save/Load data
✓ String operations
✓ Math helpers (lerp, clamp, distance)
✓ Debug tools

### Game Dev Helpers
✓ Scenes/States
✓ Tile maps
✓ Object pooling (optional)
✓ Screen effects (shake, flash, fade)

**Candy is now a complete game development language!** 🍬🎮



# CANDY LANGUAGE - ESSENTIAL ADDITIONS FOR C WRAPPER COMPATIBILITY
# Keep it Simple, But Cover the Bases

================================================================================
## PHILOSOPHY: ADD ONLY WHAT'S NEEDED FOR C INTEROP WITHOUT COMPLEXITY
================================================================================


## 1. POINTER/REFERENCE SUPPORT (MINIMAL & SAFE)

### The Problem:
C libraries return pointers that need to be passed around

### The Solution: Opaque Handles

```candy
// In Candy code, pointers are just "handles" - you don't manipulate them
texture = loadTexture("player.png")  // Returns a handle
draw(texture, 100, 100)              // Pass handle to function
unload(texture)                      // Done with it

// Behind the scenes:
// - texture is a pointer
// - Candy doesn't let you do pointer math
// - You just pass it around like a value
// - Simple and safe!
```

**No new syntax needed!** Candy treats pointers as opaque values automatically.


## 2. NULL/NONE HANDLING (OPTIONAL VALUES)

### The Problem:
C functions return NULL on failure

### The Solution: Simple null checks

```candy
// Check if something is null/empty
texture = loadTexture("missing.png")

if texture == null {
  print("Failed to load texture!")
  exit()
}

// Or use default value
texture = loadTexture("missing.png") or defaultTexture

// Check if valid
if texture {
  draw(texture, 100, 100)
}
```

**Keywords added:**
- `null` - represents NULL/nil/empty


## 3. STRUCT/RECORD TYPES (FOR C STRUCTS)

### The Problem:
C libraries use structs for data

### The Solution: Simple record syntax

```candy
// Create a struct inline
player = {
  x: 100,
  y: 200,
  health: 100,
  name: "Hero"
}

// Access fields
print(player.x)
player.health = player.health - 10

// Create from C library
pos = Vector2(10, 20)  // Creates C struct
print(pos.x)           // Access field
```

**No new syntax!** Already supported with `object` and `{}`


## 4. ARRAYS/BUFFERS (FOR C ARRAYS)

### The Problem:
C uses arrays and buffers

### The Solution: Keep existing array syntax

```candy
// Dynamic arrays (already have this)
numbers = [1, 2, 3, 4, 5]
numbers.add(6)

// Fixed-size arrays (for C interop)
buffer = array(100)  // Create array of 100 items
buffer[0] = 42
buffer[50] = 99

// Byte buffers (for binary data)
data = bytes(256)    // 256 byte buffer
data[0] = 65  // 'A'
```

**New functions:**
- `array(size)` - Create fixed-size array
- `bytes(size)` - Create byte buffer


## 5. TYPE ANNOTATIONS (OPTIONAL, FOR C WRAPPERS)

### The Problem:
C needs to know types for FFI

### The Solution: Optional type hints (only in library wrappers)

```candy
// Normal Candy code - NO TYPES NEEDED
x = 10
name = "Player"

// In .candylib wrappers - TYPES SPECIFIED
extern LoadTexture(filename: cstring): Texture
extern DrawTexture(texture: Texture, x: int, y: int, tint: Color): void

// When calling from Candy - TYPES AUTO-CONVERTED
texture = LoadTexture("player.png")  // String → cstring automatically
DrawTexture(texture, 100, 100, WHITE)
```

**Users never write types, only library bindings use them!**


## 6. ERROR HANDLING (SIMPLE TRY/CATCH)

### The Problem:
C libraries can fail

### The Solution: Simple error handling

```candy
// Simple try/catch (optional)
try {
  file = openFile("data.txt")
  content = readFile(file)
} catch {
  print("Error reading file!")
}

// Or check return values
result = openFile("data.txt")
if result.error {
  print("Failed: " + result.error)
} else {
  file = result.value
}

// Or just let it crash (default)
file = openFile("data.txt")  // Crashes with error message if fails
```

**New keywords:**
- `try` - Attempt code that might fail
- `catch` - Handle errors


## 7. MEMORY MANAGEMENT HINTS (FOR ADVANCED USERS)

### The Problem:
Some C libraries require manual cleanup

### The Solution: Simple resource management

```candy
// Automatic cleanup (preferred)
texture = loadTexture("player.png")
// ... use texture ...
// Automatically cleaned up when out of scope

// Manual cleanup (when needed)
texture = loadTexture("player.png")
// ... use texture ...
delete texture  // Explicit cleanup

// Or use "with" for automatic cleanup
with file = openFile("data.txt") {
  content = read(file)
  // file automatically closed here
}
```

**New keyword:**
- `with` - Automatic resource cleanup


## 8. CALLBACKS (FOR C CALLBACKS)

### The Problem:
C libraries use function pointers

### The Solution: Simple function passing

```candy
// Define a callback function
fun onButtonClick() {
  print("Button clicked!")
}

// Pass it to C library
button = createButton("Click Me", onButtonClick)

// Or use inline function
button = createButton("Click Me", fun() {
  print("Clicked!")
})

// Or lambda syntax (even simpler)
button = createButton("Click Me", () => {
  print("Clicked!")
})
```

**New syntax:**
- `() => { }` - Lambda/arrow function (shorthand)


## 9. CONSTANTS/ENUMS (FOR C DEFINES AND ENUMS)

### The Problem:
C uses #define and enums

### The Solution: Simple const values

```candy
// In .candylib wrapper
const RED = Color(255, 0, 0, 255)
const BLUE = Color(0, 0, 255, 255)
const KEY_LEFT = 263
const KEY_RIGHT = 262

// Use in Candy
circle(100, 100, 25, RED)
if key(KEY_LEFT) {
  // ...
}

// Or group them
Keys = {
  LEFT: 263,
  RIGHT: 262,
  UP: 265,
  DOWN: 264
}

if key(Keys.LEFT) {
  // ...
}
```

**Already supported!** No changes needed.


## 10. VARIADIC FUNCTIONS (FOR PRINTF-STYLE)

### The Problem:
C has functions like printf(format, ...)

### The Solution: String interpolation instead

```candy
// Instead of printf("%d %s", score, name)
// Use string interpolation (easier!)
print("Score: {score} Name: {name}")

// Or concatenation
print("Score: " + score + " Name: " + name)

// For actual C variadic functions in wrappers
extern printf(format: cstring, ...args): int

// Candy auto-converts to C varargs
printf("Score: %d\n", score)
```

**New feature:**
- `{variable}` - String interpolation


## 11. BITWISE OPERATIONS (FOR FLAGS)

### The Problem:
C uses bitwise operations for flags

### The Solution: Add bitwise operators

```candy
// Bitwise operators
flags = FLAG_FULLSCREEN | FLAG_VSYNC  // OR
flags = flags & ~FLAG_VSYNC           // AND NOT
flags = flags ^ FLAG_DEBUG            // XOR

// Check flag
if flags & FLAG_FULLSCREEN {
  print("Fullscreen mode")
}

// Shifts
value = 1 << 5   // Left shift: 32
value = 64 >> 2  // Right shift: 16
```

**New operators:**
- `|` - Bitwise OR
- `&` - Bitwise AND  
- `^` - Bitwise XOR
- `~` - Bitwise NOT
- `<<` - Left shift
- `>>` - Right shift


## 12. CASTING (FOR TYPE CONVERSION)

### The Problem:
Sometimes need to convert between types

### The Solution: Simple conversion functions

```candy
// Type conversion
number = toInt("123")        // String to int
text = toString(456)         // Int to string
decimal = toFloat("3.14")    // String to float

// Or use built-in functions
x = int("100")
y = float("3.14")
s = string(42)

// For pointers (advanced)
handle = toPointer(12345)    // Rarely needed
```

**Already mostly covered!** Just add consistency.


================================================================================
## COMPLETE UPDATED SYNTAX SUMMARY
================================================================================

### Core Language (No Changes!)
```candy
x = 10                    // Variables
if x > 5 { }             // Conditionals
for i = 0 to 10 { }      // Loops
fun myFunc() { }         // Functions
object Player { }        // Objects
[1, 2, 3]               // Arrays
```

### New Additions (For C Interop)

```candy
// 1. Null handling
if value == null { }
value = getValue() or defaultValue

// 2. Try/catch
try { 
  risky() 
} catch { 
  print("Error!") 
}

// 3. With (auto cleanup)
with file = open("data.txt") {
  content = read(file)
}

// 4. Lambda functions
onClick(() => { print("Click!") })

// 5. String interpolation
print("Score: {score}")

// 6. Bitwise operators
flags = FLAG_A | FLAG_B
if flags & FLAG_A { }

// 7. Fixed arrays
buffer = array(100)
data = bytes(256)

// 8. Type hints (ONLY in .candylib wrappers)
extern myFunc(x: int, name: cstring): pointer
```


================================================================================
## TRANSLATION TABLE UPDATES (For Transpiler)
================================================================================

| Candy Code | C Code |
|------------|--------|
| `value == null` | `value == NULL` |
| `try { code } catch { handle }` | `if (setjmp(jmp_buf) == 0) { code } else { handle }` |
| `with x = open() { use }` | `{ Type x = open(); use; cleanup(x); }` |
| `() => { print("Hi") }` | `void lambda_1() { printf("Hi\n"); }` |
| `print("x={x}")` | `printf("x=%d\n", x);` |
| `flags \| FLAG_A` | `flags \| FLAG_A` |
| `array(100)` | `int arr[100];` or `malloc(100*sizeof(int))` |
| `bytes(256)` | `uint8_t buf[256];` or `malloc(256)` |


================================================================================
## EXAMPLE: FULL C WRAPPER USAGE
================================================================================

### raylib.candylib (Generated by candy-bindgen)

```candy
library "raylib" {
  
  // Types
  type Texture {
    id: uint
    width: int
    height: int
  }
  
  type Color {
    r: byte
    g: byte
    b: byte
    a: byte
  }
  
  // Constants
  const RED = Color(255, 0, 0, 255)
  const BLUE = Color(0, 0, 255, 255)
  const WHITE = Color(255, 255, 255, 255)
  const KEY_LEFT = 263
  const KEY_RIGHT = 262
  
  // External functions
  extern InitWindow(width: int, height: int, title: cstring): void
  extern CloseWindow(): void
  extern WindowShouldClose(): bool
  extern BeginDrawing(): void
  extern EndDrawing(): void
  extern ClearBackground(color: Color): void
  extern LoadTexture(filename: cstring): Texture
  extern UnloadTexture(texture: Texture): void
  extern DrawTexture(texture: Texture, x: int, y: int, tint: Color): void
  extern IsKeyDown(key: int): bool
  
  // Candy-friendly wrappers
  fun window(w, h, title) {
    InitWindow(w, h, title)
    SetTargetFPS(60)
  }
  
  fun circle(x, y, r, color) {
    DrawCircle(x, y, r, color)
  }
  
  fun key(keyCode) {
    return IsKeyDown(keyCode)
  }
}
```


### game.candy (User code - SIMPLE!)

```candy
import "raylib"

window(800, 600, "My Game")

// Load resources with error handling
playerTexture = LoadTexture("player.png")
if playerTexture == null {
  print("Error: Could not load player.png!")
  exit()
}

x = 400
score = 0

// Game loop
loop {
  
  // Input
  if key(KEY_LEFT) {
    x = x - 5
  }
  if key(KEY_RIGHT) {
    x = x + 5
  }
  
  // Draw
  BeginDrawing()
  ClearBackground(WHITE)
  
  DrawTexture(playerTexture, x, 300, WHITE)
  text("Score: {score}", 10, 10, 20, BLACK)
  
  EndDrawing()
}

// Cleanup
UnloadTexture(playerTexture)
CloseWindow()
```


================================================================================
## FINAL ADDITIONS SUMMARY
================================================================================

### MUST ADD (For C Interop):
1. ✓ `null` keyword - for NULL checks
2. ✓ `try/catch` - simple error handling
3. ✓ `with` - automatic resource cleanup
4. ✓ `() => {}` - lambda syntax (optional, use `fun()` works too)
5. ✓ `{variable}` - string interpolation
6. ✓ Bitwise operators (`|`, `&`, `^`, `~`, `<<`, `>>`)
7. ✓ `array(size)` and `bytes(size)` - fixed buffers

### DON'T ADD:
- ✗ Manual memory management (use automatic)
- ✗ Pointer arithmetic (unsafe)
- ✗ Complex type system (keep simple)
- ✗ Manual struct definitions (use object)
- ✗ Preprocessor (not needed)
- ✗ Templates/generics (too complex)
- ✗ Operator overloading (confusing)
- ✗ Multiple inheritance (complicated)


================================================================================
## THE RESULT
================================================================================

**Candy remains:**
- ✓ Simple enough for kids
- ✓ No semicolons
- ✓ No explicit types (in user code)
- ✓ Automatic memory management
- ✓ Easy graphics and games

**But now can:**
- ✓ Use ANY C library
- ✓ Handle errors properly
- ✓ Work with pointers safely
- ✓ Manage resources
- ✓ Interface with low-level code

**Total new syntax: ~7 simple additions**
**Learning curve: Still 30 minutes!**

🍬 **Candy stays sweet, but can now talk to C!** 🍬
