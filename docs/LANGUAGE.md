# Candy language specification

Candy is a gradually typed, **BASIC-meets-C** scripting language: braces, `fun`, lists, classes, and an optional **BlitzBasic-style 3D** layer on top of Raylib. Source files use the **`.candy`** extension.

**Teaching / minimal subset:** the **kid (ultra-simple) dialect and examples** are documented in [CANDY_KID.md](CANDY_KID.md) with scripts under [examples/candy/](../examples/candy/README.md) (e.g. `kid_*.candy`).

Identifiers are **ASCII case-insensitive** at lex time (for example `KEY_ESC` and `key_esc` are the same name).

---

## Conditional statements

Braces are required on all branches.

```candy
if (health <= 0) {
    println("Game over")
} else if (health < 20) {
    println("Low health")
} else {
    println("You are safe")
}
```

Nested `if` works as usual. For a compact value choice, you can use the ternary operator: `alive ? "yes" : "no"`.

---

## Arithmetic operators

| Operator | Name | Description | Example |
|:---:|---|---|---|
| `+` | Addition | Adds together two values | `x + y` |
| `-` | Subtraction | Subtracts one value from another | `x - y` |
| `*` | Multiplication | Multiplies two values | `x * y` |
| `/` | Division | Divides one value by another | `x / y` |
| `%` | Modulus | Returns the division remainder | `x % y` |
| `x = x + 1` | Increment pattern | Increases a variable by 1 | `x = x + 1` |
| `x++` | Increment | Increases a variable by 1 | `x++` |
| `x = x - 1` | Decrement pattern | Decreases a variable by 1 | `x = x - 1` |
| `x--` | Decrement | Decreases a variable by 1 | `x--` |

---

## Comparison operators

| Operator | Name | Description | Example |
|:---:|---|---|---|
| `==` | Equal to | Returns true if values are equal | `x == y` |
| `!=` | Not equal | Returns true if values are not equal | `x != y` |
| `>` | Greater than | Returns true if left is greater | `x > y` |
| `<` | Less than | Returns true if left is less | `x < y` |
| `>=` | Greater than or equal | Returns true if left is greater or equal | `x >= y` |
| `<=` | Less than or equal | Returns true if left is less or equal | `x <= y` |

---

## Logical operators

You may use C-style or English spellings (all are lowered consistently):

| Style   | Operators        |
|---------|------------------|
| C-style | `&&` `||` `!`    |
| English | `and` `or` `not` |

---

## Loops

### `loop`
Runs forever until you `break` or close the window.
```candy
loop {
    update()
    show()
}
```

### `repeat N`
Runs exactly N times.
```candy
repeat 5 {
    println("Candy!")
}
```

### `while`
```candy
while (health > 0) {
    update()
}
```

### `do` / `while`
```candy
do {
    step()
} while (running)
```

### `for` - inclusive `to` / `step` (BASIC-style)
Bounds are **inclusive**.
```candy
for i = 1 to 5 {
    println(i)
}
```

### `foreach` and `for each` (iterate the candy jar)
```candy
foreach (candy in candyJar) {
    println(candy)
}
```

---

## `break` and `continue`

```candy
i = 0
while (i < 10) {
    println(i)
    i = i + 1
    if (i == 4) {
        break
    }
}
```

```candy
i = 0
while (i < 10) {
    if (i == 4) {
        i = i + 1
        continue
    }
    println(i)
    i = i + 1
}
```

---

## `switch`

Use `case` with a block or a single statement; `default` is optional.

```candy
switch (level) {
    case 1: { println("Start") }
    case 10: { println("Boss") }
    default: { println("Explore") }
}
```

---

## Functions - `fun`

```candy
fun greet() {
    println("Sweet!")
}

fun add(a: int, b: int): int {
    return a + b
}

// Tuple return type and tuple literal (multiple unwraps)
fun getPosition(): (float, float, float) {
    return (x, y, z)
}
```

Alternate spellings: `function`, `func`, `sub` (some forms differ slightly in the parser).

Parameters may use `name: type` or `name as type`. Return type is introduced with `:` or `as` before the `{`.

---

## Classes and objects

Primary-constructor classes map parameters onto instance fields with the same names.

```candy
class CandyBar(name: string, flavor: string, calories: int) {
    fun getInfo(): string {
        return name + " - " + flavor + " (" + calories + " cal)"
    }
    fun unwrap() {
        println("Unwrapping " + name + "!")
    }
}

snickers = CandyBar("Snickers", "Chocolate Peanut", 250)
println(snickers.name)
snickers.unwrap()
```

Extra fields use `var` / `val` inside the class body (plain `name = expr` member lines are not stored as fields today—use `var`).

```candy
class Entity {
    var x = 0.0
    var y = 0.0
    var z = 0.0
    var health = 100

    fun move(dx: float, dy: float, dz: float) {
        x = x + dx
        y = y + dy
        z = z + dz
    }

    fun takeDamage(amount: int) {
        health = health - amount
        if (health <= 0) {
            destroy()
        }
    }

    fun destroy() {
        println("Entity melted!")
    }
}
```

Inside methods, `this` is the instance (receiver names are supported for future use).

---

## Enumerations

```candy
enum Sweetness {
    Mild,
    Medium,
    ExtraSweet = 2
}

enum CandyType {
    Chocolate = 1,
    Gummy = 2
}

sugar = Sweetness.ExtraSweet
```

---

## Lists (arrays)

```candy
candies = ["Skittles", "M&Ms", "Twix"]
println(candies[0])
println(candies.length)

candies.append("Starburst")
candies.remove_at(0)
```

Use `.add`, `.push`, or `.append` to grow a list; `.remove` drops the first matching value; `.remove_at` / `.delete_at` removes by index.

---

## Types and inference

```candy
playerName = "SugarRush"
health as int = 100
speed: float = 5.5
alive = true
```

Core type names used in signatures are typically lowercase: `int`, `float`, `string`, `bool`.

---

## Built-in commands (stdlib + prelude)

### Output

`print`, `println`, `printf` — join arguments with spaces for `println`.

### Console input and parsing

- `readLine()` — read one line from standard input (no arguments).
- `parseInt(s)` — decimal string to integer (whitespace trimmed).

### Math (top-level and `math.*`)

`sqrt`, `sin`, `cos`, `tan`, `pow`, `abs`, `floor`, `ceil`, `round`, `min`, `max`, `clamp`, and the `math` module.

### Random

- `rand(lo, hi)` — inclusive integer range (same as `random(lo, hi)` / `randomInt` aliases).
- Module `random` still exposes `.int`, `.float`, `.shuffle`, etc.

### Strings

- `len(s)` — length of string or list.
- `toUpper(s)` / `toLower(s)` — aliases for `string.upper` / `string.lower`.
- `string.trim`, `string.split`, `string.join`, … — see [STDLIB.md](STDLIB.md).

### String interpolation

```candy
score = 100
println("Score: {score}")
```

### Memory helpers

`new` / `delete(expr)` exist for host-managed cleanup patterns; many scripts use scopes and lists instead.

### Prelude constants

`PI`, `E`, and **`key_esc`** (integer scan code for Escape, for `keyHit`).

---

## 3D helpers (BlitzBasic-style, Raylib-backed)

Build with **`-tags raylib`** so these builtins exist. They sit alongside the full Raylib flat API (see [EXTENSIONS.md](EXTENSIONS.md)).

| Builtin | Role |
|---------|------|
| `window(w, h, t)` | Opens a window (alias for `Graphics3D`). |
| `fps(60)` | Sets target frames per second. |
| `clear(color)` | Clears the screen with a color constant. |
| `show()` | Swaps buffers and optionally clears (alias for `flip`). |
| `circle(x, y, r, c)` | Draws a filled circle. |
| `box(x,y,w,h,c)` | Draws a filled rectangle. |
| `text(msg,x,y,sz,c)` | Draws text. |
| `key(name)` | `true` if key is down. |
| `clicked()` | `true` if mouse left button pressed. |
| `exit()` | Closes the window. |
| `CreateCamera()` | Resets the default 3D camera (handle `0`). |
| `PositionEntity`, `RotateEntity`, `MoveEntity`, `TurnEntity` | Transform entities. |
| `RenderWorld()` | Draws 3D entities. |

Typical frame loop:

```candy
window(800, 600, "Candy Game")
fps(60)
camera = CreateCamera()
mesh = LoadMesh("assets/model.glb")

loop {
    TurnEntity(mesh, 0, 1, 0)
    RenderWorld()
    show()
    if key("ESCAPE") { exit() }
}
```

Also use `windowShouldClose()`, `beginDrawing`, `endDrawing`, and `clearBackground` when you mix 2D UI with 3D.

---

## Comments

```candy
// Single-line

/*
   Multi-line
*/
```

---

## `main` and top-level scripts

You can wrap entry code in `fun main() { ... }` and call `main()` at the end of the file, or simply write statements at the top level (the `candy` runner executes the program’s statements in order).

---

## Example programs (see `examples/candy/`)

| File | What it shows |
|------|----------------|
| `candy_shop.candy` | Classes, `foreach`, list `append` / `remove_at`, methods. |
| `candy_crusher.candy` | `while`, `readLine`, `parseInt`, `rand`, control flow. |
| `candy_rain.candy` | Blitz-style 3D loop with `CreateSphere`, `ColorEntity`, `RenderWorld`, `flip`. |

---

## Practical notes

- Prefer **`var field = init`** for mutable class fields that must exist on every instance.
- Methods that **`return` early** still flush field mutations back onto the instance (fixed in the evaluator).
- For games beyond the toy entity list, prefer the full Raylib bindings and helpers in [GAME_HELPERS.md](GAME_HELPERS.md).
