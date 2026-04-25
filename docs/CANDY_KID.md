# Candy language specification (kid / ultra-simple)

**Philosophy:** if a 12-year-old can’t understand it, it is too complex.

**Candy in this repo (KGO):** programs are run by the **Candy interpreter** in `compiler/`. For windows, drawing, and input, build the `candy` tool with **`-tags raylib`**. The sections below are the **language**; a short **Implementation (KGO)** block at the end states what the interpreter does today. A separate **Candy → C (Raylib)** path is a **design sketch** only—there is no `candyc` transpiler in this tree yet.

**Runnable spec examples (Raylib build):** `examples/candy/kid_moving_ball.candy`, `kid_clicker.candy`, `kid_catch.candy`, `kid_bounce.candy` — also listed in [../examples/candy/README.md](../examples/candy/README.md) (see [GETTING_STARTED.md](GETTING_STARTED.md) for the exact `go run` line).

**Extended (games):** [CANDY_KID_EXTENDED.md](CANDY_KID_EXTENDED.md) — collision aliases, `sprite`/`draw`, time helpers, and mapping to the larger Raylib / game-helper surface.

---

## Variables (no type annotations)

```candy
x = 10
name = "Sugar"
speed = 5.5
active = true
```

Candy infers the runtime kind from the value.

## Math (school-style)

```candy
result = 5 + 3
difference = 10 - 4
product = 6 * 7
quotient = 20 / 4

x = x + 1  // add one
x++        // postfix: bump x (also prefix ++x / --x on a simple variable is supported)
y = y - 1
y--
```

String concatenation uses `+` when at least one side is a string (e.g. `"Score: " + score`).

## Print

```candy
print("Hello!")            // line of output
print(42)
print("Score: " + score)
```

In KGO, `print` writes a line (like a single print with a newline to the user).

## Input

```candy
name = input("What's your name? ")
age = input("How old are you? ")
```

`input` can print an optional prompt, then reads a line of text (stdin).

## If

```candy
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
```

Use `and` / `or` / `not` (or `!`) in conditions as supported by the implementation.

## Loops

**Inclusive numeric `for`:**

```candy
// count 0..10
for i = 0 to 10 {
  print(i)
}

// count backwards
for i = 10 to 0 step -1 {
  print(i)
}
```

**`repeat`:**

```candy
repeat 5 {
  print("Yay!")
}
```

**`while`:**

```candy
while playing {
  updateGame()
}
```

**`for` over a list:**

```candy
candies = ["Skittles", "M&Ms", "Twix"]
for candy in candies {
  print(candy)
}
```

**`loop` (game / window):** runs the body in a way that can stop when the window should close in the Raylib build (see **Implementation (KGO)** at the end of this file).

## Lists (arrays)

```candy
scores = [10, 20, 30, 40, 50]
names = ["Alice", "Bob", "Charlie"]

print(scores[0])
print(names[2])
scores[0] = 100
scores.add(60)
scores.remove(2)  // 0-based index: third item; see Implementation (KGO) for exact semantics
print(scores.count)
```

In KGO, **`.remove(i)`** removes the element at **index** `i` (0-based, negative = from the end). To remove the first matching **value**, use **`removeFirst`**, aliases **`remove_value`** / **`removevalue`**. List length: **`.length`**, **`.count`**, or **`len(x)`**.

## Functions

```candy
fun sayHi() {
  print("Hello!")
}

fun greet(name) {
  print("Hello " + name)
}

fun add(a, b) {
  return a + b
}

sayHi()
greet("Candy")
result = add(5, 3)
```

`fun` is the same keyword class as `function` in the full language.

## Random

```candy
dice = random(1, 6)      // inclusive integer range
coin = random(0, 1)
percent = random(0, 100)
```

In KGO the top-level `random(min, max)` uses the interpreter’s inclusive integer RNG (not Raylib’s `GetRandomValue`).

## Simple graphics (Raylib build)

**Open a window:**

```candy
window(800, 600, "My Game")
```

**Main loop and frame:**

```candy
loop {
  clear(WHITE)
  circle(400, 300, 50, RED)
  box(100, 100, 200, 100, BLUE)
  text("Hello!", 100, 50, 20, BLACK)
  show()
}
```

In KGO, `window` also sets a default 60 FPS target. `clear` + `show` are wired so one frame is: start drawing and clear, draw primitives, end drawing (see Implementation).

`text` accepts the kid form **`text(message, x, y, size, color)`** and also the older form **`text(x, y, message, size, color)`**.

`image` draws a file-backed image at `(x, y)` with a path-based cache in KGO.

## Complete examples

### Example 1: moving ball

```candy
window(800, 600, "Moving Ball")

x = 400
y = 300

loop {
  if key(LEFT)  { x = x - 5 }
  if key(RIGHT) { x = x + 5 }
  if key(UP)    { y = y - 5 }
  if key(DOWN)  { y = y + 5 }
  clear(WHITE)
  circle(x, y, 25, RED)
  show()
}
```

A copy lives at `examples/candy/kid_moving_ball.candy` (run with `-tags raylib`).

### Example 2: clicker

```candy
window(800, 600, "Cookie Clicker")

score = 0

loop {
  if clicked() {
    score = score + 1
  }
  clear(SKYBLUE)
  circle(400, 300, 100, BROWN)
  text("Score: " + score, 300, 100, 30, BLACK)
  text("Click the cookie!", 280, 500, 20, GRAY)
  show()
}
```

### Example 3: catch

```candy
window(800, 600, "Catch the Candy")

playerX = 400
candyX = random(0, 800)
candyY = 0
score = 0

loop {
  if key(LEFT)  { playerX = playerX - 5 }
  if key(RIGHT) { playerX = playerX + 5 }
  candyY = candyY + 3
  if candyY > 550 and candyX > playerX - 30 and candyX < playerX + 30 {
    score = score + 1
    candyY = 0
    candyX = random(0, 800)
  }
  if candyY > 600 {
    candyY = 0
    candyX = random(0, 800)
  }
  clear(SKYBLUE)
  box(playerX - 30, 550, 60, 20, BLUE)
  circle(candyX, candyY, 10, PINK)
  text("Score: " + score, 10, 10, 20, BLACK)
  show()
}
```

### Example 4: bouncing ball

```candy
window(800, 600, "Bouncing Ball")

x = 400
y = 300
speedX = 5
speedY = 3

loop {
  x = x + speedX
  y = y + speedY
  if x < 0 or x > 800 {
    speedX = speedX * -1
  }
  if y < 0 or y > 600 {
    speedY = speedY * -1
  }
  clear(WHITE)
  circle(x, y, 20, RED)
  show()
}
```

## Super simple reference

| Area | What you write |
|------|------------------|
| **Drawing** | `circle(x, y, r, c)`, `box(x, y, w, h, c)`, `line(x1, y1, x2, y2, c)`, `text(msg, x, y, size, c)`, `image(file, x, y)` |
| **Colors (prelude names)** | `RED`, `BLUE`, `GREEN`, `YELLOW`, `ORANGE`, `PURPLE`, `PINK`, `WHITE`, `BLACK`, `GRAY`, `SKYBLUE`, `BROWN`, `GOLD` (see Implementation) |
| **Input** | `key(LEFT)` … (also string names), `clicked()`, `mouseX()` / `mouseY()` |
| **Sound (helpers)** | `play("x.wav")`, `music("song.mp3")`, `stopMusic()` |
| **Math (typical)** | `random(a, b)`, `sqrt`, `abs`, `sin`, `cos`, `tan` (exact set: prelude + `math` module in KGO) |
| **Strings** | `upper("hello")`, `lower("HELLO")`, `len` / `length` for size |
| **Utils** | `wait(seconds)`, `fps(60)`, `exit()` (closes window in the Raylib build) |

## Key simplifications (design goals)

1. No type declarations in source (runtime values carry the kind).  
2. No semicolons.  
3. `loop` instead of spelling out “while the window is open” every time.  
4. Short names for drawing (`circle` vs engine names).  
5. Named colors instead of only RGB.  
6. `key(LEFT)` instead of low-level scancodes in teaching material.  
7. “No `main`” in spirit: a script is the program.  
8. `for i = 0 to 10` and `repeat 5` are easy to read.

---

## Simplified compiler design: Candy → C (illustration only)

The natural teaching implementation is: **Candy source → C + Raylib → native executable**. That is **not** built in KGO; the project runs the **interpreter** and optional **LLVM** paths. The table and sketch below show how a future transpiler *could* map the kid subset.

### Translation table (illustrative)

| Candy | C (conceptual) |
|-------|----------------|
| `x = 10` | `int x = 10;` (or a tagged value if you add dynamic typing) |
| `print("Hi")` | `printf` / your runtime |
| `if x > 5 {` | `if (x > 5) {` |
| `for i = 0 to 10 {` | `for (int i = 0; i <= 10; i++)` |
| `repeat 5 {` | `for (int _i = 0; _i < 5; _i++)` |
| `loop {` (window) | `while (!WindowShouldClose())` |
| `window(800, 600, "G")` | `InitWindow(...); SetTargetFPS(60);` |
| `clear(WHITE)` | `BeginDrawing(); ClearBackground(WHITE);` (must pair with `show`) |
| `show()` | `EndDrawing();` |
| `circle(…)` | `DrawCircle(…)` |
| `box(…)` | `DrawRectangle(…)` |
| `key(LEFT)` | `IsKeyDown(KEY_LEFT)` |
| `clicked()` | `IsMouseButtonPressed(MOUSE_BUTTON_LEFT)` (naming per Raylib) |
| `random(1, 10)` | `GetRandomValue(1, 10)` in a C+Raylib build |

### Sketch: toy string-based transpiler in Go (illustration only; not the KGO compiler)

A real `candyc` would use a real parser, symbol table, and correct lowering for strings, `text`, and `+`. The following is the common teaching sketch (regex replace) from the spec:

```go
package main

import (
	"regexp"
	"strings"
)

func TranspileToC(candy string) string {
	c := "#include <stdio.h>\n#include \"raylib.h\"\n\nint main() {\n"

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

	forPattern := regexp.MustCompile(`for (\w+) = (\d+) to (\d+)`)
	candy = forPattern.ReplaceAllString(candy, "for (int $1 = $2; $1 <= $3; $1++)")

	repeatPattern := regexp.MustCompile(`repeat (\d+)`)
	candy = repeatPattern.ReplaceAllString(candy, "for (int _i = 0; _i < $1; _i++)")

	c += candy
	c += "\n    CloseWindow();\n    return 0;\n}\n"
	return c
}
```

**Hypothetical build (if you had emitted C):**

```bash
# 1. Candy to C
./candyc game.candy > game.c

# 2. C to executable
gcc game.c -lraylib -o game

# 3. Run
./game
```

**Candy (kid catch fragment):**

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

**Possible generated C (Raylib; string formatting in real code would use `TextFormat` or similar):**

```c
#include <stdio.h>
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

---

## Implementation (KGO) — what actually runs

| Topic | Behavior |
|-------|----------|
| **Run** | `go run -tags raylib ./cmd/candy <file.candy>` from `compiler/` (or your installed `candy` with `raylib` build tags). |
| **Frame** | `clear` ensures a begin + clear; `show` ends the frame. `flip()` is for the alternate `while !shouldClose()` + flip style. |
| **`key`** | `key(LEFT)` uses prelude constants or `key("left")` string; integer scancodes are allowed. |
| **`.remove(n)`** | **Index** `n` in KGO; use `removeFirst` to drop first matching value. |
| **C transpiler** | **Not present**; the section above is a roadmap / teaching aid. |
| **More API** | [GAME_HELPERS.md](GAME_HELPERS.md), `compiler/candy_raylib/register.go`, and `prelude` in `compiler/candy_evaluator/prelude.go`. |

### Spec coverage (this document vs the KGO interpreter)

**Documented in this file:** the full kid syntax, examples, reference tables, and the C/Raylib *illustration* (no `candyc` tool). **Runnable copies** of examples 1–4 are in `examples/candy/kid_*.candy` (Raylib build).

| Area | Status in KGO |
|------|----------------|
| Variables, `//` comments, math, string `+` when a side is a string | Implemented |
| `x++` / `x--` (postfix) and `++x` / `--x` (prefix) on a variable | Implemented (interpreter) |
| `print` | Implemented; each call prints a **line** (with newline) |
| `input("…?")` | Implemented (stdin) |
| `if` / `else` / `else if`, `and` / `or` / `not` | Implemented |
| `for i = a to b` / `step`, `repeat n`, `while`, `for x in list` | Implemented |
| `loop { }` (game loop; exits with window in Raylib build) | Implemented |
| Lists: `[...]`, index, `.add`, **`.remove(i)` by 0-based index**, `.count` / `.length`, `len` | Implemented; for “remove first value `x`” use **`removeFirst`**, not `remove` |
| `fun` … / `return` | Implemented (same as `function` in the full language) |
| `random(a, b)` inclusive integers | Implemented (interpreter RNG, not Raylib `GetRandomValue` unless you wrap it yourself) |
| `window`, `clear`, `show`, `circle`, `box`, `line`, `text` (kid order **or** `x, y, msg, …`), `image` | Implemented with **`-tags raylib`**; `window` also sets 60 FPS by default; `text` has two call shapes (see [GETTING_STARTED](GETTING_STARTED.md)) |
| Color names `RED` … | Set in the prelude (string tokens resolved by the graphics builtins) |
| `key(LEFT)` …, `clicked`, `mouseX` / `mouseY` | Implemented (Raylib) |
| `play` / `music` / `stopMusic` | Implemented via game/Raylib helpers (paths, streaming; see [GAME_HELPERS](GAME_HELPERS.md)) |
| `sqrt`, `abs`, `sin`, `cos`, `tan` | Available from prelude / builtins and Raylib `math` registrations; see `compiler/candy_evaluator/prelude.go` and `compiler/candy_raylib/register.go` for the exact set |
| `upper` / `lower` / `length` (or `len`) | `upper` / `lower` and **`len` = `length`** for strings and arrays |
| `wait(1.0)` seconds, `fps(60)`, `exit()` | **`wait` is registered in the Raylib build** and uses the engine’s wait-in-seconds; **`fps`** sets target FPS; **`exit`** closes the window |
| **Candy → C** (`candyc`, `gcc` pipeline, toy `TranspileToC`) | **Not implemented in this repo**; only described here as a possible future compiler |
| **LLVM / native** codegen for postfix or full game API | **Partial**; prefer the **interpreter** for the kid spec |

**Bottom line:** the **teaching spec** and **KGO** match for everything above except: **no C transpiler** yet; **list `remove` is by index** (not “find value”); **random** is the interpreter’s integer RNG; **`text`** supports a legacy parameter order; **`wait`** in the “utilities” list assumes the **Raylib** runtime so it matches “seconds of delay”.

This document is the kid-facing **Candy** surface; the full product may include more keywords and features described elsewhere in `docs/`.
