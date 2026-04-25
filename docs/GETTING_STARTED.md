# Getting Started with Candy 🍬

Candy is a scripting language built for making games. It reads like plain English, needs no boilerplate, and has a built-in 2D/3D graphics engine (Raylib) so you can draw things on screen within minutes.

This guide goes from zero to a working game — no prior programming experience required.

For the **ultra-simple “kid” dialect** of Candy (teaching order: `loop`, `clear` / `show`, the built-in color names, and full spec examples), read **[CANDY_KID.md](CANDY_KID.md)**. Runnable copies live under `examples/candy/kid_*.candy` (e.g. `kid_moving_ball.candy`, `kid_clicker.candy`, `kid_catch.candy`, `kid_bounce.candy`); build `candy` with `-tags raylib` as below.

For **teaching names** for collisions, `sprite`/`draw`, time (`seconds`, `deltaTime`), `join` / `toNumber` at the top level, and how they map to the real engine, read **[CANDY_KID_EXTENDED.md](CANDY_KID_EXTENDED.md)**. [examples/candy/kid_extended_sandbox.candy](../examples/candy/kid_extended_sandbox.candy) is a one-loop intro; [examples/candy/kid_every_after.candy](../examples/candy/kid_every_after.candy) runs **`every(...)`** and **`after(...)`** with named `fun` callbacks. See the guide for the full map.

For native C bindings generated from headers, see **[CANDYWRAP.md](CANDYWRAP.md)**.
For the team-facing "Single-Command Integration" contract (module-folder standard, ABI guardrails, and manifest boundary), see **[CANDYWRAP_WORKFLOW.md](CANDYWRAP_WORKFLOW.md)**.
For language-level interop additions (`with`, bitwise flags, `array/bytes`, null/try/catch guidance), see **[C_INTEROP_ESSENTIAL_ADDITIONS.md](C_INTEROP_ESSENTIAL_ADDITIONS.md)**.

---

## 1. Install and run

### Step 1 — Install Go

Candy is written in Go. Download it from **https://go.dev/dl/** and run the installer.

### Step 2 — Get the source

```bash
git clone https://github.com/yourname/KGO.git
cd KGO/compiler
```

### Step 3 — Build the Candy executable

On Windows:
```bash
go build -tags raylib -o candy.exe ./cmd/candy
```

On macOS / Linux:
```bash
go build -tags raylib -o candy ./cmd/candy
```

> **Why `-tags raylib`?** Without it you get a plain scripting runtime with no graphics. The `raylib` tag adds the full game API.

### Step 4 — Run your first script

Create `hello.candy`:

```candy
println("Hello from Candy!")
```

Run it:
```bash
./candy hello.candy
```

You should see `Hello from Candy!` printed. That's it — you're running Candy.

---

## 2. The language in five minutes

You don't need to read all of this before making a game. Skim it and come back when you need a specific piece.

### Variables

Assign a value and Candy figures out the type automatically:

```candy
score = 0
playerName = "Hero"
speed = 3.5
alive = true
```

You can also be explicit if you want:

```candy
health as float = 100.0
lives as int = 3
```

### Strings and interpolation

Wrap a variable name in `{}` inside a string to embed it:

```candy
name = "Alice"
println("Hello, {name}!")          // Hello, Alice!
println("Score: {score} pts")      // Score: 0 pts
```

### Basic Math

| Operator | Action | Example |
|:---:|---|---|
| `+` | Add | `x + 10` |
| `-` | Subtract | `x - 5` |
| `*` | Multiply | `x * 2` |
| `/` | Divide | `x / 2` |
| `%` | Modulus | `x % 2` |
| `++` | Increment | `++x` |
| `--` | Decrement | `--x` |
| `x++` | Easier! | `x++` |
| `x--` | Easier! | `x--` |

### Comparing values

| Operator | Result | Example |
|:---:|---|---|
| `==` | Equal | `score == 100` |
| `!=` | Not equal | `lives != 0` |
| `>` | Greater than | `health > 50` |
| `<` | Less than | `y < 0` |
| `>=` | Greater or equal | `x >= 10` |
| `<=` | Less or equal | `x <= 0` |

### If / else

```candy
if health <= 0 {
    println("Game over!")
} else if health < 25 {
    println("Low health!")
} else {
    println("You're fine.")
}
```

### Loops

**Count from A to B** (both ends inclusive):
```candy
for i = 1 to 5 {
    println(i)          // 1 2 3 4 5
}
```

**Step by a custom amount:**
```candy
for i = 0 to 100 step 10 {
    println(i)          // 0 10 20 30 … 100
}
```

**Loop over a list:**
```candy
enemies = ["Slime", "Goblin", "Dragon"]
foreach e in enemies {
    println("Fighting {e}")
}
```

**Loop many times:**
```candy
repeat 5 {
  println("Candy is sweet!")
}
```

**Loop forever until you exit:**
```candy
loop {
  println("Spinning...")
  if key("SPACE") { break }
}
```

### Functions

```candy
fun damage(amount) {
    health = health - amount
    if health < 0 {
        health = 0
    }
}

damage(10)
println("Health: {health}")
```

Functions can return values:

```candy
fun clamp(val, lo, hi) {
    if val < lo { return lo }
    if val > hi { return hi }
    return val
}

speed = clamp(speed, 0, 10)
```

### Lists (arrays)

```candy
items = ["sword", "shield", "potion"]

println(items[0])           // sword
println(len(items))         // 3

items = items + ["bow"]     // add an element
```

Iterate:
```candy
for each item in items {
    println(item)
}
```

### Maps

```candy
player = {name: "Hero", x: 100, y: 200, hp: 100}

println(player.name)        // Hero
player.hp = player.hp - 5
```

---

## 3. Opening a window

Now for the fun part. Every Candy game follows the same three-step structure:

```candy
// 1. Open a window
window(800, 600, "My Game")
fps(60)

// 2. Game loop — runs ~60 times per second
loop {
    clear(BLACK)          // wipe the screen each frame

    // --- draw things here ---
    text("It works!", 20, 20, 32, WHITE)

    show()                  // show everything you drew
}
```

Save as `mygame.candy` and run:
```bash
./candy mygame.candy
```

A window opens showing white text on black. Close it with the × button or press Escape.

### The game loop explained

| Line | What it does |
|---|---|
| `window(800, 600, "Title")` | Opens an 800×600 pixel window |
| `setTargetFPS(60)` | Locks the loop to ~60 updates per second |
| `while !shouldClose()` | Keeps running until the user closes the window |
| `clear("black")` | Fills the screen with a solid colour before drawing |
| `flip()` | Swaps the back-buffer to the screen (shows your drawing) |
| `closeWindow()` | Frees all graphics resources — always call this at the end |

---

## 4. Drawing shapes and text

```candy
window(800, 600, "Shapes")
setTargetFPS(60)

while !shouldClose() {
    clear("darkgray")

    drawCircle(400, 300, 60, "red")             // filled circle
    drawRectangle(100, 100, 200, 80, "blue")    // filled rectangle
    drawLine(0, 0, 800, 600, "yellow")          // diagonal line
    drawText(300, 20, "Hello!", 28, "white")    // text at x,y

    flip()
}
closeWindow()
```

**Coordinate system:** `(0, 0)` is the **top-left** corner. X grows rightward, Y grows downward.

### Colors

Pass a name string or an `{r,g,b,a}` map:

```candy
drawCircle(200, 200, 40, "red")
drawCircle(300, 200, 40, {r:255, g:165, b:0, a:255})   // orange
drawCircle(400, 200, 40, color(128, 0, 128))            // purple
```

Named colors: `"white"`, `"black"`, `"red"`, `"green"`, `"blue"`, `"yellow"`, `"orange"`, `"purple"`, `"pink"`, `"gray"`, `"darkgray"`, `"skyblue"`, `"gold"`, `"lime"`, `"magenta"`, `"brown"`, `"maroon"`, `"darkblue"`, `"darkgreen"`, `"violet"`, `"beige"`, `"transparent"`.

---

## 5. Reading keyboard input

```candy
window(800, 600, "Move the square")
fps(60)

x = 400
y = 300
speed = 4

loop {
    // Move with arrow keys
    if key(RIGHT) { x = x + speed }
    if key(LEFT)  { x = x - speed }
    if key(DOWN)  { y = y + speed }
    if key(UP)    { y = y - speed }

    clear(BLACK)
    box(x - 20, y - 20, 40, 40, LIME)
    text("Arrow keys to move", 10, 10, 20, WHITE)
    show()
}
```

**Key constants** — call them as functions (they return the key code number):

| Constant | Key |
|---|---|
| `KEY_RIGHT()` | → arrow |
| `KEY_LEFT()` | ← arrow |
| `KEY_UP()` | ↑ arrow |
| `KEY_DOWN()` | ↓ arrow |
| `KEY_SPACE()` | Space bar |
| `KEY_ENTER()` | Enter |
| `KEY_ESCAPE()` | Escape |
| `KEY_A()` … `KEY_Z()` | Letter keys |

**Input functions:**

| Function | Returns |
|---|---|
| `isKeyDown(key)` | `true` while the key is held |
| `isKeyPressed(key)` | `true` only on the frame the key is first pressed |
| `isKeyReleased(key)` | `true` on the frame the key is released |
| `getMouseX()` / `getMouseY()` | Mouse cursor position |
| `isMouseButtonDown(btn)` | Mouse button held (`MOUSE_LEFT_BUTTON()` etc.) |

---

## 6. Delta time — smooth movement at any frame rate

Using a fixed speed value like `speed = 4` means the game moves faster on a 120 FPS machine than a 30 FPS one. Fix this with **delta time**: multiply speeds by the time since the last frame.

```candy
window(800, 600, "Smooth movement")
setTargetFPS(60)

x = 400.0
y = 300.0
speed = 200.0   // pixels per second

while !shouldClose() {
    dt = getFrameTime()   // seconds since last frame, e.g. 0.016

    if isKeyDown(KEY_RIGHT()) { x = x + speed * dt }
    if isKeyDown(KEY_LEFT())  { x = x - speed * dt }
    if isKeyDown(KEY_DOWN())  { y = y + speed * dt }
    if isKeyDown(KEY_UP())    { y = y - speed * dt }

    clear("black")
    drawRectangle(x - 20, y - 20, 40, 40, "skyblue")
    flip()
}
closeWindow()
```

Now `speed = 200` means "200 pixels per second" regardless of frame rate.

---

## 7. Your first complete game — Catch the Dot

Here is a small but complete game. A dot moves around the screen; click it to score a point. After 10 seconds the game ends and shows your score.

```candy
window(800, 600, "Catch the Dot!")
setTargetFPS(60)

// --- game state ---
score = 0
timeLeft = 10.0
dotX = 400.0
dotY = 300.0
dotR = 25.0
dotSpeedX = 180.0
dotSpeedY = 140.0

fun moveDot(dt) {
    dotX = dotX + dotSpeedX * dt
    dotY = dotY + dotSpeedY * dt

    // bounce off edges
    if dotX - dotR < 0   { dotX = dotR;       dotSpeedX = math.abs(dotSpeedX) }
    if dotX + dotR > 800 { dotX = 800 - dotR; dotSpeedX = -math.abs(dotSpeedX) }
    if dotY - dotR < 0   { dotY = dotR;       dotSpeedY = math.abs(dotSpeedY) }
    if dotY + dotR > 600 { dotY = 600 - dotR; dotSpeedY = -math.abs(dotSpeedY) }
}

fun checkClick() {
    if isMouseButtonPressed(MOUSE_LEFT_BUTTON()) {
        mx = getMouseX()
        my = getMouseY()
        dx = mx - dotX
        dy = my - dotY
        dist = math.sqrt(dx * dx + dy * dy)
        if dist <= dotR {
            score = score + 1
            // speed up slightly each catch
            dotSpeedX = dotSpeedX * 1.05
            dotSpeedY = dotSpeedY * 1.05
        }
    }
}

// --- game loop ---
while !shouldClose() {
    dt = getFrameTime()
    timeLeft = timeLeft - dt

    if timeLeft <= 0 {
        // Game over screen
        clear("black")
        drawText(260, 220, "Time's up!", 48, "yellow")
        drawText(280, 300, "Score: {score}", 36, "white")
        drawText(220, 360, "Close the window to exit", 22, "gray")
        flip()
        continue
    }

    moveDot(dt)
    checkClick()

    clear({r: 20, g: 20, b: 40, a: 255})   // dark navy background
    drawCircle(dotX, dotY, dotR, "red")
    drawText(10, 10, "Score: {score}", 24, "white")
    drawText(10, 40, "Time: {math.floor(timeLeft + 1)}", 24, "yellow")
    drawText(10, 570, "Click the dot!", 18, "gray")
    flip()
}
closeWindow()
```

Save as `catch.candy` and run it. Try catching the dot before time runs out!

---

## 8. Loading images and sounds

### Images

```candy
// Load a PNG and draw it
tex = loadTexture("player.png")

while !shouldClose() {
    clear("black")
    drawTexture(tex, 100, 200)      // draw at x=100, y=200
    flip()
}
unloadTexture(tex)
closeWindow()
```

### Sounds

```candy
initAudioDevice()
sfx = loadSound("jump.wav")
music = loadMusicStream("theme.ogg")
playMusicStream(music)

while !shouldClose() {
    updateMusicStream(music)        // call every frame to keep music playing

    if isKeyPressed(KEY_SPACE()) {
        playSound(sfx)
    }

    clear("black")
    flip()
}
unloadSound(sfx)
unloadMusicStream(music)
closeAudioDevice()
closeWindow()
```

---

## 9. Where to go next

| Goal | Read |
|---|---|
| All drawing functions (shapes, text, textures, 3D) | [EXTENSIONS.md](EXTENSIONS.md) |
| Full language reference (structs, closures, imports, …) | [LANGUAGE.md](LANGUAGE.md) |
| Math, file, JSON, random, and other stdlib modules | [STDLIB.md](STDLIB.md) |
| Built-in physics engine (gravity, collisions, raycasting) | [EXTENSIONS.md → Physics](EXTENSIONS.md#physics) |
| Networking (multiplayer with ENet) | [EXTENSIONS.md → ENet](EXTENSIONS.md#enet-networking-module) |
| Building a native binary with LLVM | Run `candy -help` or read the root README |

### Useful patterns to study

- `examples/candy/kid_moving_ball.candy` — [kid spec](CANDY_KID.md): `loop` + `clear` / `show` + `key`
- `examples/candy/kid_clicker.candy` — `clicked()` and score
- `examples/candy/kid_catch.candy` — `random`, collision idea
- `examples/candy/kid_bounce.candy` — `or`, bouncing
- `compiler/scratch/physics_demo.candy` — falling boxes, raycasts
- `compiler/scratch/enet_server_demo.candy` — multiplayer server
- `compiler/scratch/enet_client_demo.candy` — multiplayer client

---

## 10. Native build command matrix

Use these checklist-style commands for the native compiler flow:

- `candy build program.candy`
  - Compiles to LLVM IR and tries to emit a native binary.
- `candy run program.candy`
  - Runs the script directly in the evaluator path.
- `candy compile program.candy -o output`
  - Alias of `build`, with explicit output path.

Build-path flags:

- `-o <path>`
  - Output native binary path.
- `--debug`
  - Debug profile (`clang -O0 -g`).
- `--optimize`
  - Shipping optimization profile.
- `--verbose`
  - Prints compilation/link commands and profile details.

Toolchain health:

- `candy doctor`
  - Runs native toolchain preflight without compiling.
  - Shows `clang`/`opt` discovery status and search order.
  - Returns non-zero when required native build tooling is missing.

Examples:

```bash
# Build with defaults
candy build game.candy

# Verify toolchain before building
candy doctor

# Build optimized shipping binary
candy build --optimize game.candy

# Build debug binary to custom path
candy compile game.candy --debug -o build/game-debug

# Run through evaluator
candy run game.candy
```

---

## Quick reference card

```candy
// Window
window(w, h, "Title")      fps(60)
loop { ... }                exit()
getFrameTime()              getFPS()

// Draw (inside loop, between clear() and show())
clear(BLACK)                show()
text("txt", x, y, size, WHITE)
box(x, y, w, h, BLUE)
circle(x, y, r, RED)
line(x1, y1, x2, y2, GREEN)  // or drawLine(...)

// Input
key(RIGHT)                  clicked()
mouseX()     mouseY()       input("Name? ")

// Math
math.sqrt(x)   math.abs(x)   random(1, 10)
length(list)                 upper("candy")

// Print
print("hello")              println("Score: {score}")
```
