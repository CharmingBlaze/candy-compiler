# SimpleC-BASIC: Clean, Simple, Modern

A blend of C's directness, BASIC's readability, with carefully chosen modern features.

> **Beginner path & honest feature matrix:** for “teach it like BASIC, add C/types later” — what the dynamic interpreter does today vs the LLVM build vs not yet — see [BEGINNER_FRIENDLY.md](BEGINNER_FRIENDLY.md) and the long-form checklist in [VISION_ROADMAP.md](VISION_ROADMAP.md). For a **one-page syntax cheat sheet** with **Spec → Candy** status per feature, see **Syntax summary (Spec → Candy)** below. For **optional Raylib builtins** (`window`, `circle`, …) behind the `raylib` build tag, see [EXTENSIONS.md](EXTENSIONS.md).

## Syntax summary (Spec → Candy)

Quick reference: teaching names vs what the **tree-walking interpreter** (`candy run` / default execution) actually runs. For detail, follow the links in **Notes**.

| Feature | Candy syntax (examples) | Interpreter | Notes |
|---------|-------------------------|------------|--------|
| Variables | `x = 10` | **Yes** | Assignment; see [Variables - Simple and Clear](#variables---simple-and-clear). |
| “Typed” variables | `x: int = 10`, type-first `int x = 10`, or untyped | **Partial** | Teaching `x as int = 10` maps to `:` / type-first forms; see [Spec → Candy (variables)](#spec--candy-variables-arrays-pointers). |
| Constants | `const PI = 3.14` | **Yes** (parsed like `val`) | `const` is parsed to the same statement shape as `val`; the dynamic runtime does not enforce **immutability** the way a static const would. |
| Functions | `fun add(a, b) { return a + b }` — also `func` / `function` | **Yes** | See [Functions - BASIC Style with Types](#functions---basic-style-with-types). |
| Subroutines | `sub greet() { print "Hi!" }` | **Yes** | Treated like a function in the dynamic runtime. |
| Structs | `struct Player { x: int; y: int }` and `Player { x: 1, y: 2 }` | **Yes** | See [Simple Structs with Methods](#simple-structs-with-methods). |
| `if` / `else` | `if score > 100 { } else { }` | **Yes** | `if` is a **statement**, not a general expression value in eval; see [Control Flow](#control-flow). |
| `for` with bounds | `for i = 0 to 10 { }`, `for i = a to b step c { }` | **Yes** | See [Loops](#loops). |
| “For each” / iteration | **`for item in items { }`** — not the keywords `for each` | **Yes** | There is **no** `each` token; “read it as for-each” in books, but type `for` + `in` only. |
| `while` | `while cond { }` | **Yes** | |
| Arrays | `numbers = [1, 2, 3]` | **Yes** | [Arrays and Lists](#arrays-and-lists). |
| Maps (dicts) | `scores = { "A": 100 }` or `map { "a": 1 }` | **Yes** | [Dictionaries (Maps)](#dictionaries-maps). |
| Line comments | `// comment` | **Yes** | Lexer skips `//` to end of line. |
| Block comments | `/* … */` | **Yes** | Skipped in the lexer (see [lexer tests](candy_lexer/lexer_test.go)); can span lines. |
| Pointers, `malloc` / `free`, `->` | N/A in user syntax for the VM | **No** | [Pointers, memory, and C-style features](#pointers-memory-and-c-style-features). |
| `defer` | `defer expr` | **No** (eval) | [Parsed only — not run in the interpreter](#defer-parsed-not-run-in-the-interpreter). |
| `null`, `??`, `?.` | `a ?? b`, `x == null` | **`??` Yes; `?.` not** | [Null, `??`, and safe navigation](#null--and-safe-navigation). |
| Imports / stdlib | `import "file.candy"`, `import math`, `math.sqrt(16)` | **Partial** | [Modules, imports, and standard library](#modules-imports-and-standard-library). |
| Graphics (`window`, `circle`, `flip`, …) | `window(800, 600, "Game")` | **No** (default binary) / **Yes** (`go build -tags raylib`) | Go-wrapped Raylib; not in default `Builtins`. See [EXTENSIONS.md](EXTENSIONS.md) and [scratch/bounce_raylib.candy](scratch/bounce_raylib.candy). |

### Key features (teaching vs interpreter)

| Pitch | Interpreter today | Where to read |
|-------|-------------------|---------------|
| **BASIC-style:** `fun` / `sub`, `for` … `to` … `step`, iteration with `for x in y`, `print` | **Yes** | [Loops](#loops), [Output - BASIC style](#output---basic-style) |
| **C-style (full):** raw pointers, `&`/`*`, `malloc` / `free`, all operators, `->` | **No** in the Go VM | [Pointers, memory, and C-style features](#pointers-memory-and-c-style-features) — use values, `.`, lists; **LLVM** is a different story. |
| **Structs, methods, fixed-array teaching** | **Partial** | [Structs](#simple-structs-with-methods); fixed arrays: see variables / LLVM. |
| **“Modern”:** string interpolation, list/string methods, lambdas `(x) => x` | **Yes** | [Strings](#strings), [Lambdas](#lambdas) |
| **`defer`, safe-nav `?.`** | `defer` **not** in eval; `?.` not full | [`defer` section](#defer-parsed-not-run-in-the-interpreter), [Null / `??`](#null--and-safe-navigation) |
| **`try` / `catch` / `??`** | **try/`catch`/`??` Yes** in eval | [Error handling](#error-handling-try--catch--finally) |
| **Native / fast compilation** | **Partial** | Only where **`candy build` / `candy compile`** supports the AST; see [Implementation map](#implementation-map-in-this-repo) and [VISION_ROADMAP](VISION_ROADMAP.md) — not “every program is native”. |
| **Omitted in v1 teaching (generics, heavy operator overloading, async, …)** | **No** / **Partial** as listed in long docs | [Generics](#genericstemplates), [Async](#asyncawait) — not beginner-complete in eval. |
| **Host graphics (Raylib)** | **`window` / `circle` / …** | **Optional** (`-tags raylib`) | [EXTENSIONS.md](EXTENSIONS.md) — same three-tier idea: `.candy` calls builtins; Go wraps C. |

## Language levels (vision)

**Level 1 — Dynamic (beginner):** no types required; `name = "Alice"`, `print x`, lists, and `for` / `while` with a **tree-walking interpreter** (`candy run` / default execution).

**Level 2 — Structure:** structs with fields and methods, gradually introduced.

**Level 3 — Performance:** add `: type` and return types so **`candy build` / `candy compile`** can emit **LLVM IR** and optional native binaries.

### Implementation map (in this repo)

| Area | Role |
|------|------|
| `candy_evaluator` | **Dynamic runtime** — `Value` (int, float, string, bool, array, map, function, struct def), `Env`, **assign** (`a = b`), `print`/`println` builtins |
| `candy_llvm` | **Native** — subset of the same AST; use when you add types and need speed |
| `candy_typecheck` | Optional static checking; **monomorph** for simple generics on the LLVM path |
| `candy_lexer` | Newlines can insert `;` so **semicolons are optional** in many cases |

**Recently in the interpreter:** assignment to `a[i]` and `obj.field` (with **negative** indices, including **rune**-wise string indexing), **array concat** with `+` (e.g. `a + [x]`), **integer range values** `lo..hi` (inclusive both ends, descending if `lo > hi`) usable as `for v in 0..10 { … }`, `struct T { }` + `T { f: v }` instances, string interpolation `"a {x} b"`, `for v in …` (array / string / map keys, or a range), `for i = a to b`, `while`, `try` / `catch` / `finally` (caught `Error` is bound to the `catch` variable as a **string**), struct **method calls** `o.m()` (C-style methods in a struct body), and no extra `null` line from `candy` when the last value is `null`.

**Source files** use the `.candy` extension. Imports in source refer to other modules by path, e.g. `import "other.candy";`.

Gaps to close over time: closer parity between **interpreted** and **compiled** paths (e.g. `try`/`catch` in native code, `switch` in eval, list slicing, `for k,v in map`). (Assignment inside **`for` / `for-in`** bodies uses a fresh scope per iteration, so outer locals are not updated; use **`while`** for simple accumulation, or refactor.) List instance methods (`.add`, `.map`, `.filter`, `sort`, …) and string methods (`upper`, `split`, …) are implemented in the **interpreter**; see [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md) and [VISION_ROADMAP.md](VISION_ROADMAP.md) for the exact matrix.

## Variables - Simple and Clear

Some teaching docs use `name as int = 0` or `scores[10] as int`. In Candy, use the forms below; the **Spec → Candy** table maps the ideas to what the parser and interpreter do today.

### Spec → Candy (variables, arrays, pointers)

| Teaching / spec idiom | Candy you can type today | Interpreter / parser notes |
|------------------------|--------------------------|----------------------------|
| `count as int = 0` | `int count = 0` (type-first) or `val count: int = 0` | `as` is a reserved word but **not** used for `x as type` declarations; use **C-style type-first** or `val` / `var` with **`:`**. |
| `active as bool = true` | `bool active = true` or `val active: bool = true` | Same as above. |
| `const MAX_HEALTH = 100` | `const MAX_HEALTH = 100` | Parsed like `val`; **UPPERCASE is convention only** — the runtime does not enforce immutability on `const` today. |
| `numbers = [1,2,3,4,5]` | `numbers = [1, 2, 3, 4, 5]` | Dynamic list (`ValArray`). |
| `scores[10] as int` / fixed C array | `int scores[10]` | The parser notes an array type (the internal name becomes `int[]`); the **VM does not model a fixed-size C array** — treat as a type hint for tools / future native code, not a bound-checked fixed buffer in the REPL. |
| `matrix[5][5] as float` | **No fixed 2D array type** | Use nested **dynamic** lists, e.g. a list of rows, each row a list of floats, or a flat list + manual indexing. |
| `ptr as int* = null` / `ptr = &value` | Types like `int*` can appear in **type positions** in the parser | **Address-of `&expr` is not implemented** in the tree-walking interpreter (`PrefixExpression` does not handle `&`). Prefer ordinary names and values; do not rely on C pointer semantics in `.candy` on the dynamic path. |

```c
// Simple variables - type inferred
name = "Alice"
score = 100
x = 10.5
playing = true

// Explicit types: C-style (type first) or val/var with ':'
int count = 0
float speed = 5.5
string playerName = "Player"
bool active = true
val count2: int = 0

// Constants (UPPERCASE is style, not a language rule)
const MAX_HEALTH = 100
const PI = 3.14159

// Arrays
numbers = [1, 2, 3, 4, 5]          // dynamic list
int scores[10]                    // parse/type hint; not a real fixed C array in the VM

// Pointers: avoid on the dynamic interpreter — `&` is not supported at runtime
// int* ptr = null;   // may parse in some contexts; &value will not run in eval
```

## Simple Structs with Methods

```c
// Basic struct
struct Vector2 {
    float x
    float y
    
    // Methods right in the struct!
    float length() {
        return sqrt(x * x + y * y)
    }
    
    void normalize() {
        float len = length()
        x = x / len
        y = y / len
    }
    
    Vector2 add(Vector2 other) {
        return Vector2 { x + other.x, y + other.y }
    }
}

// Use it
vec = Vector2 { x: 10, y: 20 }
print vec.length()
vec.normalize()

// Create with constructor-like syntax
pos = Vector2 { 100, 200 }
```

## Structs with Constructors

```c
struct Player {
    string name
    float x, y
    int health
    int maxHealth
    
    // Constructor - called when created
    Player(string n, float startX, float startY) {
        name = n
        x = startX
        y = startY
        health = 100
        maxHealth = 100
    }
    
    // Methods
    void move(float dx, float dy) {
        x += dx
        y += dy
    }
    
    void damage(int amount) {
        health -= amount
        if health < 0 {
            health = 0
        }
    }
    
    bool isAlive() {
        return health > 0
    }
    
    float distanceTo(Player other) {
        dx = x - other.x
        dy = y - other.y
        return sqrt(dx * dx + dy * dy)
    }
}

// Create player
player = Player("Hero", 100, 200)
player.move(10, 5)
player.damage(20)

if player.isAlive() {
    print "Still alive!"
}
```

## Functions - BASIC Style with Types

Keywords **`fun`**, **`function`**, and **`func`** all use the same token: pick whichever reads best. **`sub name(...) { }`** is also supported for routines that do not return a value (BASIC-style); the dynamic interpreter still treats it like a `function` (no separate `void` value is required at runtime for simple scripts).

**Parameter and return types** in `fun`/`function`/`func` forms use a **colon**: `name: type` for each parameter, and **`: returnType` before the `{`**. A teaching syntax like `a as int` or `) as int` is **not** the same as Candy’s surface today — use **`:`** instead. Alternatively, use **C-style** declarations: `int add(int a, int b) { ... }` (type name first, then the function name).

**Default parameters** (e.g. `width = 800` in the parameter list) are **not** implemented in the parser/AST for user functions, and the interpreter does **not** fill in missing arguments. **Workaround:** always pass every argument, or use two differently named functions.

```c
// Type inference
fun multiply(a, b) {
    return a * b
}

// fun + colons: parameters and return type (not `a as int` — use `a: int`)
fun add(a: int, b: int): int {
    return a + b
}

// C-style: type-first function header (choose one style per program)
int add2(int a, int b) {
    return a + b
}

// No return (void-style or sub)
void greet(string name) {
    print "Hello, {name}!"
}

sub greet2(name: string) {
    print "Hello, {name}!"
}

// Default parameter values: not supported — define explicit args only
void createWindow(int width, int height) {
    // call as createWindow(800, 600) from the caller
}

// Calls
result = add(5, 3)
greet("Alice")
```

**Multiple return values** using a struct (example continues below)

```c
// Multiple return values using struct
struct Result {
    int quotient
    int remainder
}

Result divide(int a, int b) {
    return Result { a / b, a % b }
}

result = divide(17, 5)
print result.quotient
print result.remainder
```

## Output - BASIC style

Some teaching materials use `printf` / `sprintf` or a print that does not add a newline. In the **tree-walking interpreter** (`candy run` / default execution), output goes through a small **`print` / `println` builtin** only.

### Spec → Candy (printing)

| Spec / idiom | Candy today | Notes |
|--------------|-------------|--------|
| `print "Hello {name}!"` | `print "Hello {name}!"` | String interpolation: `{expr}` in double-quoted strings; use `\{` in the source to emit a **literal** `{` (see [Strings](#strings) below). |
| `print "a", b` (several values) | `print "a", b` (comma between expressions) | Arguments are **space-joined** on one line, then a **newline** is always printed (see implementation: `print` and `println` both map to the same host routine). |
| `print` with and without newline | `print` / `println` | In the current interpreter, **both always end the line** like a typical `println` — there is no separate `print` that leaves the cursor mid-line. |
| C-style `printf("Score: %d\n", score)` | **No `printf` builtin** | Use interpolation: `print "Score: {score}"`, or build a string and `print` it. |
| `println` for line | `println` | Same as `print` for the interpreter. |

**Runnable example (safe for `candy` today):**

```c
print "Hello World"
print "Score:", score
print "Hello {name}!"
print x, y, z
```

**Not runnable on the dynamic interpreter (parser / future / other tiers may differ):** `printf`, `sprintf`, `fprintf` — not wired as builtins; avoid copy-pasting C-style `printf` samples into `.candy` and expecting them to run unchanged.

## Control Flow

`if` / `else if` / `else` and the one-line `if cond { … }` form are evaluated in the **interpreter**. A teaching **`a ? b : c` ternary** and **`switch` / `case` / `default` are not** implemented in `candy_evaluator` today: they may **parse** in some builds, but **`candy` will not run them** as you would expect. Use the workarounds below for copy-pastable programs.

| Feature | In interpreter? | Use instead (beginner-friendly) |
|---------|----------------|----------------------------------|
| `if` / `else` / `else if` | **Yes** | (examples below) |
| One-liner `if` | **Yes** | `if playing { update() }` |
| Ternary `a ? b : c` | **No** | Assign inside `if` / `else`, or a small `func` that `return`s one branch, or a temporary variable and two branches. |
| `switch` on a value | **No** in eval | A chain of `if` / `else if` / `else` comparing `state` to string or other values. |

**Runnable (`candy` today):**

```c
// IF
if health <= 0 {
    print "Game Over"
}

if score > 1000 {
    print "High score!"
} else if score > 500 {
    print "Good job!"
} else {
    print "Keep trying"
}

// One-liner
if playing { update() }

// Ternary (workaround) — e.g. pick a status string
status = "dead"
if health > 0 {
    status = "alive"
}
```

**Switch (workaround — same idea as `switch` on a string `state`):**

```c
if state == "menu" {
    showMenu()
} else if state == "playing" {
    updateGame()
} else if state == "paused" {
    showPause()
} else {
    print "Unknown"
}
```

> **Parser note:** a `switch { … }` form may appear in grammar or in examples aimed at the LLVM / future path; for the **REPL and interpreter**, prefer `if` / `else if` as above.

## Loops

**`for i = a to b` / `step`:** supported; **endpoints are inclusive** for both directions (e.g. `i = 0` through `9` in `for i = 0 to 9`).

**`for x in y` (iteration):** use **`for item in items`** — there is no separate `for each` keyword; write `for` and `in` only.

**`while`:** supported.

**`do { } while`:** **not** implemented in the tree-walking interpreter (may parse elsewhere) — do not treat as a guaranteed runnable form in `.candy` for the REPL. Use `while` with a condition that is checked at the start, or a `while true` with a `break`-style control when that exists in your program.

**C-style `for (int i = 0; i < 10; i++)` and `for (i as int = 0; …)`:** **not** wired in `candy_evaluator` — use **`for i = 0 to 9`** (or a `while` and manual `i`).

> **for-in and outer variables:** the body of `for v in …` runs in a **new scope** each time; reassigning a loop local does not update an outer `x` — use `while` if you need to accumulate in one scope (as noted in the [introduction](#language-levels-vision)).

| Loop form | In interpreter? | Candy spelling / workaround |
|-----------|-----------------|------------------------------|
| `for i = a to b` / `step` | **Yes** | Inclusive `a` and `b`; step can be negative. |
| `for each x in y` (teaching) | **N/A** | `for x in y` (no `each` keyword) |
| `for x in array` / string / map keys / range | **Yes** | Map iteration is **keys only** (see [VISION_ROADMAP.md](VISION_ROADMAP.md)) |
| `while` | **Yes** | |
| `do { } while` | **No** in eval | `while` + structure your condition |
| C `for(;;)` | **No** in eval | `for`/`while` as above |

**Runnable (`candy` today):**

```c
// BASIC-style for (inclusive bounds)
for i = 0 to 9 {
    print i
}

for i = 10 to 1 step -1 {
    print i
}

// for … in  (not "for each" — the word "each" is not required)
for item in items {
    print item
}

for enemy in enemies {
    enemy.update()
}

// while
while playing {
    update()
    render()
}
```

**C-style and do-while (workaround pattern — not the same syntax, but the usual replacement):**

```c
// Count 0..9 with a while (if you need C-style control)
i = 0
while i < 10 {
    print i
    i = i + 1
}
// For "at least one iteration", structure the test or use a first body before the while.
```

## Arrays and Lists

Teaching samples sometimes use `scores[100] as int` or a **fixed** C array. In Candy, use **`int scores[100]`** (type-first); the **VM still uses a dynamic list** under the hood for the interpreter — it is not a strict fixed buffer. See also *Variables - Simple and Clear*.

### Spec → Candy (lists)

| Teaching / spec | Candy today | Notes |
|------------------|-------------|--------|
| `numbers = [1,2,3]` | `numbers = [1, 2, 3]` | `ValArray`; `+` concatenates two lists. |
| `numbers[0]`, `numbers[-1]` | **Yes** | Negative index counts from the end. |
| `numbers[i] = x` | **Yes** | In-bounds assignment. |
| `numbers.add(x)` / `push` | **Yes** | In-place append (one or more values). |
| `insert(i, v)` | **Yes** | |
| `remove(3)` (by value) | **Yes** | First match by value; use **`remove_at(i)`** to delete by index. |
| `clear()` | **Yes** | |
| `size` / `length` | `.length`, `.size`, `.count` | On the list value. |
| `sort` / `reverse` / `sum` / `min` / `max` | **Yes** | |
| `evens = numbers.filter((n) => …)` / `map` | **Yes** | Pass a function or lambda. |
| `scores[100] as int` | `int scores[100]` | Parse/type name; not a real fixed-size array in the REPL VM. |

**Runnable (`candy` today):**

```c
numbers = [1, 2, 3, 4, 5]
print numbers[0]
print numbers[-1]

numbers[0] = 100
numbers.add(6)
numbers.insert(0, 0)
numbers.remove(3)
numbers.clear()

print numbers.length
size = numbers.length

numbers = [3, 1, 4, 1]
numbers.sort()
numbers.reverse()
total = numbers.sum()

evens = numbers.filter((n) => n % 2 == 0)
doubled = numbers.map((n) => n * 2)

// Type-first declaration (VM is still dynamic in the interpreter)
int board[100]
board[0] = 1
```

## Dictionaries (Maps)

Map literals use **`map { "key": value }`** or similar forms accepted by the parser; string keys are what the **dynamic** `for v in map` iterates over. There is **no** `for key, value in map` with two binders in the tree-walking interpreter — use **`for name in scores`** and index **`scores[name]`**.

### Spec → Candy (maps)

| Teaching / spec | Candy today | Notes |
|-----------------|-------------|--------|
| `{ "Alice": 100, "Bob": 95 }` (one literal) | **Use** `map { "Alice": 100 }` then `scores["Bob"] = 95` | The `map { }` literal in the current parser is effectively **one** `key: value` pair per literal; add more keys with `m["k"] = v`. |
| `scores["Alice"]` | **Yes** | Read / assign. |
| `for name, score in scores` | **Not in interpreter** | **One** loop variable = **keys only**; body: `score = scores[name]` (or pick your names). |
| `keys` / `values` | `m.keys()`, `m.values()` | Sorted key order for `values` in the host. |
| `hasKey` / `contains` key | `m.has("Alice")` or `m.contains("Alice")` | **On a map**, `contains` means **key** presence. On a **string**, `contains` is a **substring** check — different type. |
| `get` with default | `m.get("k", default)` | |

**Runnable (`candy` today):**

```c
scores = map { "Alice": 100 }
scores["Bob"] = 95
scores["Carol"] = 110

print scores["Alice"]
scores["David"] = 88

keys = scores.keys()
values = scores.values()
hasKey = scores.contains("Alice")

for name in scores {
    sc = scores[name]
    print "{name}: {sc}"
}
```

## Strings

String **interpolation** in double-quoted text uses `{expression}`. A **literal `{`** in the source string is written as **`\\{`**. `text.length` is the **rune** count (consistent with `text[i]` indexing). There is **no** `substring(start, end)` method on strings in the interpreter today — use **`s[i]`** for a single rune as a one-character string, build a small loop, or `split` / `replace` as needed.

**Not** copy-pastable in the current REPL: C-style **`char` buffers**, **`sprintf`**, and **triple-quoted** `"""` multiline blocks (treat as future / other tiers). Use [Output - BASIC style](#output---basic-style) and `print "…{x}…"`.

### Spec → Candy (strings)

| Teaching / spec | Candy today | Notes |
|-----------------|-------------|--------|
| `"Hello {name}!"` | **Yes** | Interpolation. |
| `message = """ ... """` | **Not** in the doc’s multiline form | Rely on normal strings or concatenation for now. |
| `text.upper` / `lower` | `text.upper()`, `text.lower()` | Call as methods. |
| `text.length` | **Yes** (property) | Rune count. |
| `text.substring(0, 5)` | **No** as a method | Use `s[0]`, `s[1]`, … for runes, or other workarounds above. |
| `text.split(" ")` | **Yes** | Returns a list of strings. |
| `parts.join("-")` | **Yes** | On a **list** (array), not on a string. |
| `text.contains("x")` on a string | **Yes** | Substring test. (On a **map**, `contains` is for **keys**.) |

**Runnable (`candy` today):**

```c
name = "World"
count = 42
print "Hello {name}! Count: {count}"

text = "hello world"
print text.upper()
print text.lower()
print text.length
print text.contains("world")

parts = text.split(" ")
joined = parts.join("-")

// No substring() — first rune
print text[0]
print "Score: {count}"
```

## Pointers, memory, and C-style features

Teaching materials often show **address-of** (`&`), **dereference** (`*`), **`malloc` / `free`**, **`sizeof`**, and the **arrow** operator (`->`) for pointers to structs. The **tree-walking interpreter** is built on Go with a single **`Value`** model and **garbage collection** — it does **not** expose raw addresses, manual heap lifetimes, or C layout. Use **names**, **structs by value**, and **lists** for resizable data; rely on **`candy build` / `candy compile` (LLVM)** when you need a lower-level story, not on the REPL alone.

### Spec → Candy (pointers and memory)

| Teaching / C-style | In `candy` interpreter today | Notes |
|--------------------|-----------------------------|--------|
| `ptr as int* = &value` | **No** | `&` is not evaluated in [PrefixExpression](candy_evaluator/eval.go) (only `-` and `!`). |
| `*ptr`, `*ptr = x` | **No** | No pointer `Value` kind for user code. |
| `malloc` / `free` | **No** builtins | Use lists, maps, and struct instances; memory is host/GC-backed. |
| `sizeof(T)` | **No** | — |
| `enemy->health` | **No** | Use **`enemy.health`** (dot) on struct instances. |
| `int scores[100]` / `numbers[100] as int` | **Partial** | Type-first / `int x[100]` parses; the VM still uses **dynamic** lists for interpreted runs (see *Variables* and *Arrays and Lists*). |

**Aspirational (not copy-pastable in the REPL):**

```c
// Full C-style — not supported in the dynamic VM
// value = 42
// ptr as int* = &value
// print *ptr
```

**What to use instead (runnable patterns):** ordinary variables, struct fields with `.`, `[]` on lists, and `map { }` / assignment for maps.

## `defer` (parsed, not run in the interpreter)

`defer expr` is **parsed** ([parseDeferStatement](candy_parser/statements_parser.go)) but **`candy_evaluator` does not run `defer` yet** — there is no `DeferStatement` case in `evalStatement`. Do not rely on cleanup at function exit in the REPL until this is implemented. For file-like resources, use **host** helpers (`read_file`, etc.) and **try/finally** where `finally` is supported (see below).

## Error handling (`try` / `catch` / `finally`)

**Supported** in the interpreter: `try { } catch Name e { } finally { }`. The **caught value** is the **string** message of the error (not a rich object type). See existing examples in this file and [eval_try.go](candy_evaluator/eval_try.go).

**Runnable pattern:**

```c
try { nope; } catch Error e { 42; }
```

For `finally` (outer scope), see tests and [eval_try.go](candy_evaluator/eval_try.go).

## Null, `??`, and safe navigation

- **`null`**: use `null` / `ValNull` in expressions; compare with **`== null`** or **`!= null`** where the grammar allows.
- **Nullish coalescing `a ?? b`**: **supported** in the interpreter: if `a` is **null** (and only null for `??`), the result is `b`; otherwise `a`. The **right** side is **not evaluated** when `a` is non-null (short-circuit).
- **`player?.field` (safe navigation)**: **not** fully supported in eval in the same way as Kotlin/JS; use explicit **`if player != null { player.field }`** or `player` checks before dot.

### Spec → Candy (null)

| Feature | Interpreter | Notes |
|---------|------------|--------|
| `x == null` | **Yes** | |
| `a ?? b` | **Yes** | Left nullish only; see above. |
| `player?.health` | **No** / partial | Prefer `if` + dot access. |
| `score = player?.score ?? 0` | **Partial** | Use `??` for the default; for `?.` use `if` on `player` first. |

**Runnable:**

```c
a = null
b = a ?? 42
// b is 42
c = 1
d = c ?? 99
// d is 1
```

## Lambdas

**Supported:** `(a, b) => expr` and **block** form `(a, b) => { return … }` (see [LambdaExpression](candy_evaluator/eval.go)). They work with **user calls** and with **list** methods **`map`**, **`filter`**, **`sort`**, etc., when the receiver and operands are **plain values** (e.g. compare **ints** in `sort`).

| Teaching snippet | Notes |
|------------------|--------|
| `add = (a, b) => a + b` | **Yes** |
| `multiply = (a, b) => { return a * b }` | **Yes** |
| `numbers.map(n => n * 2)` | **Yes** |
| `players.sort((a, b) => a.score > b.score)` | **Caution:** parameters may be `Value` / struct instances; for reliable behavior start with **numeric** or **simple** comparators; struct field access in lambdas depends on the instance shape. |

**Runnable:**

```c
add = (a, b) => a + b
d = add(2, 3)
nums = [3, 1, 4]
doubled = nums.map((n) => n * 2)
evens = nums.filter((n) => n % 2 == 0)
```

## Simple Inheritance

```c
// Base struct
struct GameObject {
    float x, y
    bool active
    
    GameObject(float startX, float startY) {
        x = startX
        y = startY
        active = true
    }
    
    void update() {
        // Base update
    }
    
    void render() {
        // Base render
    }
}

// Inherited struct
struct Enemy : GameObject {
    int health
    int damage
    
    Enemy(float x, float y, int hp) : GameObject(x, y) {
        health = hp
        damage = 10
    }
    
    // Override method
    override void update() {
        GameObject.update()  // Call parent
        
        // Enemy-specific update
        if health <= 0 {
            active = false
        }
    }
    
    void attack(Player player) {
        player.damage(damage)
    }
}

// Use it
enemy = Enemy(100, 200, 50)
enemy.update()
enemy.attack(player)
```

## Memory Management - Simple

```c
// Stack (automatic) - preferred
Player player = Player("Hero", 100, 200)
int numbers[100]

// Heap when needed
int* bigArray = malloc(1000 * sizeof(int))
// ... use it ...
free(bigArray)

// For objects
Player* enemy = malloc(sizeof(Player))
enemy->x = 100
enemy->health = 50
free(enemy)

// Automatic cleanup with defer (modern feature)
void loadGame() {
    FILE* file = fopen("save.dat", "r")
    defer fclose(file)  // Automatically called at function end
    
    // Read data
    // File automatically closed when function returns
}
```

## Memory Management with Smart Pointers

```c
// Stack allocation (automatic cleanup)
Player player = Player("Hero", 100, 200)
Enemy enemy = Enemy(50, 50, 30)

// Heap with automatic cleanup - use 'new'
player = new Player("Hero", 100, 200)
// Automatically freed when out of scope

// Reference to existing data
ref Player p = player
p.health = 50  // Modifies original

// Nullable/optional
maybe Player target = findPlayer(5)
if target != null {
    target.health = 0
}

// Safe navigation
target?.health = 0
print target?.score ?? 0

// Shared ownership (reference counted)
shared Texture tex = new Texture("sprite.png")
shared copy = tex  // Both reference same texture
// Freed when last reference destroyed
```

## Lambda Functions

```c
// Simple lambda
add = (a, b) => a + b
result = add(5, 3)

// With block
multiply = (a, b) => {
    result = a * b
    return result
}

// Use with arrays
numbers = [1, 2, 3, 4, 5]
doubled = numbers.map(n => n * 2)
evens = numbers.filter(n => n % 2 == 0)

// Sort with lambda
players.sort((a, b) => a.score > b.score)
```

## Modules, imports, and standard library

The **parser** and **interpreter** support a small, explicit import story. **Selective imports** (`import m (a, b)`), **`import m as x`**, and user-defined **`module { }` / `export` blocks** are *not* part of what `candy` evaluates today, even if they appear in teaching materials elsewhere in this file. The subsections below separate **what runs** from **aspirational** C-style module syntax.

### How `import` works in `candy_evaluator`

- **Syntax** ([parseImportStatement](candy_parser/parse_simplec_basic.go)): `import` plus **either** a **string** path (e.g. `import "lib.candy";` or `import "std/strings";`) **or** a **dotted** path of identifiers (e.g. `import math;`). Dots join **without** normalizing to slashes — for nested stdlib keys in [modules.go](candy_stdlib/modules.go) (e.g. `std/strings`), use the **quoted** form `import "std/strings"` so the path matches [Lookup](candy_stdlib/modules.go). Semicolons are optional in many cases.
- **Execution** ([evalImport](candy_evaluator/eval.go)):
  1. If the path matches a **stdlib** name in [candy_stdlib/modules.go](candy_stdlib/modules.go) (`math`, `file`, `fs`, `json`, `random`, `rand`, `time`, and `std/...` entries with real source), the embedded file is **parsed and evaluated** in the current environment. Many stdlib “packages” (e.g. `math`) are only **comments** in the map; the real **`math`**, **`file`**, etc. objects come from the **prelude** (see below), so `import "math";` is often **redundant** for calling `math.sqrt` but is still valid.
  2. Otherwise the path is treated as a **file** relative to the evaluator’s current working directory (`env.Cwd`, updated while loading nested file imports). The same absolute path is only run **once** per `Imported` map.
- **Prelude** ([registerPrelude](candy_evaluator/prelude.go)): at the start of every [Eval](candy_evaluator/eval.go), the host injects `math`, `file` (and alias `fs`), `json`, `random` (and `rand`), `time` module **values**, plus global **`PI`**, **`E`**, **`pi`**, **`e`**. You do **not** need to `import` to use `math.sqrt(16)` in the normal REPL/VM — but imports of stdlib **names** still work for consistency and for loading `std/...` modules that have real code.

#### Spec → Candy (imports and modules)

| Teaching / spec | In `candy` interpreter | Notes |
|-----------------|------------------------|--------|
| `import math` / `import "math"` | **Yes** | Unquoted form uses identifier path; string uses literal path. Both can resolve to stdlib [Lookup](candy_stdlib/modules.go). |
| `import "other.candy"` (relative/absolute file) | **Yes** | Resolved from `Cwd` when path is not absolute. |
| `import math (sqrt, max)` | **No** | Not in the grammar. Use `math.sqrt` or a local `f = math.sqrt` assignment. |
| `import math as m` | **No** | Not in the grammar. Use a variable: `m = math` if the prelude already defined `math`. |
| `module math { export … }` | **No** in [evalStatement](candy_evaluator/eval.go) | Parsed in some code paths, **not** executed as a module system in the tree-walking runtime. |
| `import "std/strings"` and other `import "std/…"` | **Yes** (if the key exists and the module source runs) | Keys in [modules.go](candy_stdlib/modules.go) use **slashes**. `import std.strings` (dotted) builds the path `std.strings` and will **not** find `std/strings` in Lookup — use a **string** import. |

**Runnable (prelude, no import required):**

```c
math.sqrt(16)
math.max(10, 20)
```

**Runnable (explicit stdlib import — optional):**

```c
import "math"
math.sqrt(16)
```

**Runnable (load another file from the same directory as `Cwd`):**

```c
// contents of "player.candy" are executed; path is relative to the running program’s CWD
import "player.candy"
```

### Standard library and builtins (interpreter)

Two mechanisms overlap by design: **prelude module objects** (`math.*`, `file.*`, …) and **top-level built-in functions** merged into the same [Eval](candy_evaluator/eval.go) environment. Many file/math/time helpers exist **both** as `name(...)` and as `file.name(...)` / `math.name(...)` where [prelude.go](candy_evaluator/prelude.go) and [stdlib_builtins.go](candy_evaluator/stdlib_builtins.go) agree.

#### Prelude modules — methods on `math`, `file` / `fs`, `json`, `random` / `rand`, `time`

| Module | Callable members (as `module.name(args…)`) | Constants / fields (where applicable) |
|--------|---------------------------------------------|----------------------------------------|
| `math` | `sqrt`, `pow`, `abs`, `floor`, `ceil`, `round`, `sin`, `cos`, `tan`, `min`, `max`, `clamp` | `PI`, `E`, `pi`, `e`, `Inf` |
| `file` / `fs` | `read`, `write`, `read_file` / `readFile`, `write_file` / `writeFile`, `read_lines` / `readLines`, `exists` / `file_exists` / `fileExists`, `delete` / `remove` / `delete_*`, `list` / `list_dir` / `list_files` / … | — |
| `json` | `parse`, `stringify`, `load` / `load_file` / `loadFile`, `save` / `save_file` / `saveFile` | — |
| `random` / `rand` | `int` (inclusive int range **two** args: min, max), `float` (**half-open** float range: `[min,max)` from host impl), `choice` / `sample` (one **array** arg), `seed` | — |
| `time` | `millis` / `ms` / `now_ms` / `nowMs` (current time ms), `sleep` / `sleep_ms` / `sleepMs` / `wait` (ms), `sleep_sec` / `sleepSec` (seconds) | — |

Exact bindings are defined in [candy_evaluator/prelude.go](candy_evaluator/prelude.go).

#### Top-level builtins (no `import` — from [builtin.go](candy_evaluator/builtin.go) + `init` in [stdlib_builtins.go](candy_evaluator/stdlib_builtins.go))

Representative **global** callables: `print`, `println`, `len`, `cwd`, `joinpath` / `joinPath` patterns as registered, `read_file` and aliases, `write_file`, `read_lines`, `json_stringify` / `json_parse`, `load_json` / `save_json`, math as globals (`sqrt`, `sin`, …), `random` / `randomInt` (int between two **inclusive** bounds), `random_float` (float **half-open** interval), `choose`, `sleep` / `sleepMs`, `sleep_sec`, `time_millis`, `assert`, `type` / `typeof`, `is_int`, `is_string`, `range`, `debug`, `ok` / `err`, `getenv` / `readfile` (see merged map) — the merged table in [stdlib_builtins.go](candy_evaluator/stdlib_builtins.go) is authoritative. **Detail tables** (I/O, list methods, etc.) live in [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md).

#### Spec → Candy (common teaching names vs this runtime)

| Teaching name | In interpreter | Use instead |
|---------------|---------------|-------------|
| `printf("format", …)` | **No** | `print(…)` and/or string interpolation `"x = {x}"` |
| `int(x)`, `float(x)` as C-style conversion builtins | **No** | Rely on dynamic values; `type(x)`; no generic `int(...)` in merged [Builtins](candy_evaluator/builtin.go). |
| `random.pick(list)` | **No** (name) | **`random.choice(list)`** or `random.sample(list)` |
| `random.shuffle(arr)` | **No** | Not in [prelude.go](candy_evaluator/prelude.go). |
| `file.copy(…)` | **No** | — |
| `time.now()` | **No** (name) | **`time.millis()`** (or `time.ms` / `time.now_ms` — same implementation) |
| `import m (a, b)` / `import m as x` | **No** | See import table above. |

### Aspirational module syntax (not copy-pastable in the REPL)

The following matches **teaching** / future native stories; **`module` bodies with `export` are not executed** by the tree-walking interpreter’s `evalStatement` switch, and **selective / alias imports are not in the parser**.

```c
// Aspirational — not evaluated as a module in candy_evaluator today
module math {
    const PI = 3.14159
    export float sqrt(float x) { /* ... */ }
    export int max(int a, int b) {
        if a > b { return a; }
        return b
    }
    float internal_helper() { return 0.0; }
}

// import math (sqrt, max)  // not in grammar
// import math as m         // not in grammar

import math
result = math.sqrt(16)
maximum = math.max(10, 20)
```

## Enums and Pattern Matching

```c
// Simple enum
enum GameState {
    Menu,
    Playing,
    Paused,
    GameOver
}

state = GameState.Playing

// Enum with values
enum Color {
    Red = 0xFF0000,
    Green = 0x00FF00,
    Blue = 0x0000FF
}

// Pattern matching on enums
switch state {
    case GameState.Menu: {
        showMenu()
    }
    case GameState.Playing: {
        updateGame()
    }
    case GameState.Paused: {
        showPause()
    }
    case GameState.GameOver: {
        showGameOver()
    }
}

// Match expression (returns value)
message = switch state {
    case GameState.Menu => "Press Start",
    case GameState.Playing => "Playing",
    case GameState.Paused => "Paused",
    case GameState.GameOver => "Game Over"
}

// Match with patterns
switch response {
    case {status: 200, data}: {
        print "Success: {data}"
    }
    case {status: 404}: {
        print "Not found"
    }
    case {status: s} where s >= 500: {
        print "Server error"
    }
}
```

## Error Handling

```c
// Result type for operations that can fail
struct Result<T, E> {
    bool isOk
    T value
    E error
}

Result<int, string> parseNumber(string text) {
    if isNumeric(text) {
        return Result.ok(parseInt(text))
    } else {
        return Result.error("Not a number")
    }
}

// Pattern match on result
result = parseNumber("123")
switch result {
    case {isOk: true, value: v}: {
        print "Got: {v}"
    }
    case {isOk: false, error: e}: {
        print "Error: {e}"
    }
}

// Or use helper methods
if result.isOk() {
    print result.value
} else {
    print result.error
}

// Try-catch for exceptions
try {
    file = openFile("data.txt")
    data = file.read()
    processData(data)
} catch FileException e {
    print "File error: {e.message}"
} catch Exception e {
    print "Error: {e.message}"
} finally {
    cleanup()
}
```

## Generics/Templates

```c
// Generic function (syntax sketch — dynamic interpreter may not run full generics)
T max<T>(T a, T b) {
    if a > b { return a; }
    return b
}

result = max(10, 20)        // int version
result = max(5.5, 3.2)      // float version
result = max("abc", "xyz")  // string version

// Generic struct
struct Container<T> {
    T value
    
    Container(T v) {
        value = v
    }
    
    T get() {
        return value
    }
    
    void set(T v) {
        value = v
    }
}

intBox = Container<int>(42)
strBox = Container<string>("Hello")

// Multiple type parameters
struct Pair<T, U> {
    T first
    U second
    
    Pair(T f, U s) {
        first = f
        second = s
    }
}

pair = Pair<int, string>(42, "Answer")
```

## Properties

```c
struct Player {
    private int _health = 100
    private int _maxHealth = 100
    
    // Property with getter and setter
    int health {
        get {
            return _health
        }
        set {
            _health = clamp(value, 0, _maxHealth)
        }
    }
    
    // Read-only property
    int maxHealth {
        get {
            return _maxHealth
        }
    }
    
    // Computed property
    float healthPercent {
        get {
            return float(_health) / float(_maxHealth)
        }
    }
    
    // Auto-property (generates private field)
    string name { get; set; } = "Player"
    float x { get; set; } = 0
    float y { get; set; } = 0
}

// Use properties like fields
player.health = 150      // Calls setter, clamped to 100
print player.healthPercent  // Calls getter
player.name = "Hero"     // Auto-property setter
```

## Operator Overloading

```c
struct Vector2 {
    float x, y
    
    // Operator overloading
    Vector2 operator+(Vector2 other) {
        return Vector2 { x + other.x, y + other.y }
    }
    
    Vector2 operator-(Vector2 other) {
        return Vector2 { x - other.x, y - other.y }
    }
    
    Vector2 operator*(float scalar) {
        return Vector2 { x * scalar, y * scalar }
    }
    
    bool operator==(Vector2 other) {
        return x == other.x && y == other.y
    }
    
    // Indexer
    float operator[](int index) {
        if index == 0 { return x }
        if index == 1 { return y }
        throw "Index out of range"
    }
}

// Use operators naturally
a = Vector2 { 1, 2 }
b = Vector2 { 3, 4 }
c = a + b           // Vector2 { 4, 6 }
d = a * 2.5         // Vector2 { 2.5, 5 }
same = (a == b)     // false

// Access by index
print a[0]  // 1
print a[1]  // 2
```

## Interfaces

```c
// Define interface
interface IDamageable {
    void takeDamage(int amount)
    int getHealth()
    bool isAlive()
}

// Implement interface
struct Player : IDamageable {
    int health = 100
    
    void takeDamage(int amount) override {
        health -= amount
        if health < 0 { health = 0 }
    }
    
    int getHealth() override {
        return health
    }
    
    bool isAlive() override {
        return health > 0
    }
}

struct Enemy : IDamageable {
    int health = 50
    
    void takeDamage(int amount) override {
        health -= amount
    }
    
    int getHealth() override {
        return health
    }
    
    bool isAlive() override {
        return health > 0
    }
}

// Use interface polymorphically
void damageEntity(ref IDamageable entity, int amount) {
    entity.takeDamage(amount)
    if !entity.isAlive() {
        print "Entity destroyed!"
    }
}

player = Player()
enemy = Enemy()

damageEntity(player, 20)
damageEntity(enemy, 30)
```

## Async/Await

```c
// Async function
async string loadData(string url) {
    response = await http.get(url)
    data = await response.text()
    return data
}

async void loadGame() {
    try {
        data = await loadData("https://api.game.com/save")
        player = parsePlayer(data)
        print "Loaded: {player.name}"
    } catch Exception e {
        print "Failed: {e.message}"
    }
}

// Run async function
run loadGame()

// Parallel execution
async void loadAll() {
    results = await all(
        loadData("url1"),
        loadData("url2"),
        loadData("url3")
    )
    
    for result in results {
        print result
    }
}

// Race condition (first to complete)
async void loadFastest() {
    data = await race(
        loadFromServer(),
        loadFromCache()
    )
    print "Got data: {data}"
}
```

## Attributes/Decorators

```c
// Define custom attribute
struct Serializable {
    string format = "json"
}

struct EditorVisible {
    string category = "General"
}

// Use attributes
[Serializable(format: "binary")]
[EditorVisible(category: "Player")]
struct Player {
    [SaveField]
    string name
    
    [SaveField]
    int score
    
    [DontSave]
    Texture sprite  // Not saved
    
    [Range(0, 100)]
    int health = 100
}

// Query attributes at compile time or runtime
fields = typeof(Player).fields
saveFields = fields.filter(f => f.hasAttribute("SaveField"))

for field in saveFields {
    print "Save field: {field.name}"
}
```

## Reflection

```c
// Get type information
typeInfo = typeof(Player)
print "Type: {typeInfo.name}"

// List fields
for field in typeInfo.fields {
    print "Field: {field.name} : {field.type.name}"
}

// List methods
for method in typeInfo.methods {
    print "Method: {method.name}"
    for param in method.parameters {
        print "  Param: {param.name} : {param.type.name}"
    }
}

// Create instance dynamically
instance = typeInfo.create()

// Invoke method dynamically
method = typeInfo.getMethod("takeDamage")
method.invoke(instance, [10])  // instance.takeDamage(10)

// Get/set field values
field = typeInfo.getField("health")
field.set(instance, 100)
value = field.get(instance)
```

## Standard Library Essentials

```c
// Math
import std.math

result = math.sqrt(16)
result = math.pow(2, 8)
result = math.abs(-5)
result = math.sin(math.PI / 2)
result = math.floor(5.7)
result = math.ceil(5.2)
result = math.round(5.5)
result = math.clamp(value, 0, 100)

// Collections
import std.collections

list = List<int>()
list.add(1)
list.add(2)
list.remove(1)
print list.length

map = Map<string, int>()
map["key"] = 100
if map.contains("key") {
    print map["key"]
}

set = Set<int>()
set.add(1)
set.add(2)
set.add(1)  // Duplicate ignored
print set.size  // 2

// File I/O
import std.io

file = File.open("data.txt", "r")
content = file.readAll()
file.close()

// Or with defer
file = File.open("data.txt", "w")
defer file.close()
file.writeLine("Hello World")

// Time
import std.time

now = Time.now()
print "Current time: {now}"

start = Time.now()
// ... do work ...
elapsed = Time.now() - start
print "Elapsed: {elapsed}ms"

sleep(1000)  // Sleep 1 second

// Random
import std.random

value = random(0, 100)        // Random int 0-100
value = randomFloat(0.0, 1.0) // Random float 0.0-1.0
choice = choose([1, 2, 3, 4]) // Random element
```

## Complete Game Example

```c
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 600

// ===== Vector2 =====
struct Vector2 {
    float x, y
    
    float distance(Vector2 other) {
        dx = x - other.x
        dy = y - other.y
        return sqrt(dx * dx + dy * dy)
    }
}

// ===== GameObject Base =====
struct GameObject {
    Vector2 pos
    bool active
    
    GameObject(float x, float y) {
        pos = Vector2 { x, y }
        active = true
    }
    
    void update() {
        // Override in children
    }
    
    void render() {
        // Override in children
    }
}

// ===== Player =====
struct Player : GameObject {
    int health
    float speed
    int score
    
    Player(float x, float y) : GameObject(x, y) {
        health = 100
        speed = 5.0
        score = 0
    }
    
    override void update() {
        // Input
        if key("LEFT") { pos.x -= speed }
        if key("RIGHT") { pos.x += speed }
        if key("UP") { pos.y -= speed }
        if key("DOWN") { pos.y += speed }
        
        // Clamp to screen
        if pos.x < 0 { pos.x = 0 }
        if pos.x > SCREEN_WIDTH { pos.x = SCREEN_WIDTH }
        if pos.y < 0 { pos.y = 0 }
        if pos.y > SCREEN_HEIGHT { pos.y = SCREEN_HEIGHT }
    }
    
    override void render() {
        screen.circle(pos.x, pos.y, 10, color.green)
    }
    
    void damage(int amount) {
        health -= amount
        if health < 0 { health = 0 }
    }
    
    bool isAlive() {
        return health > 0
    }
}

// ===== Enemy =====
struct Enemy : GameObject {
    Vector2 velocity
    int damage
    
    Enemy(float x, float y, float vx, float vy) : GameObject(x, y) {
        velocity = Vector2 { vx, vy }
        damage = 10
    }
    
    override void update() {
        pos.x += velocity.x
        pos.y += velocity.y
        
        // Bounce off walls
        if pos.x < 0 || pos.x > SCREEN_WIDTH {
            velocity.x = -velocity.x
        }
        if pos.y < 0 || pos.y > SCREEN_HEIGHT {
            velocity.y = -velocity.y
        }
    }
    
    override void render() {
        screen.circle(pos.x, pos.y, 8, color.red)
    }
    
    void respawn() {
        pos.x = rand() % SCREEN_WIDTH
        pos.y = rand() % SCREEN_HEIGHT
        velocity.x = (rand() / (float)RAND_MAX - 0.5) * 4
        velocity.y = (rand() / (float)RAND_MAX - 0.5) * 4
    }
}

// ===== Game Manager =====
struct Game {
    Player player
    Enemy enemies[5]
    bool running
    
    Game() {
        player = Player(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2)
        running = true
        
        // Spawn enemies
        for i = 0 to 4 {
            x = rand() % SCREEN_WIDTH
            y = rand() % SCREEN_HEIGHT
            vx = (rand() / (float)RAND_MAX - 0.5) * 4
            vy = (rand() / (float)RAND_MAX - 0.5) * 4
            enemies[i] = Enemy(x, y, vx, vy)
        }
    }
    
    void update() {
        if key("ESC") {
            running = false
            return
        }
        
        player.update()
        
        // Update and check enemies
        for i = 0 to 4 {
            enemies[i].update()
            
            // Check collision
            dist = player.pos.distance(enemies[i].pos)
            if dist < 18 {
                player.damage(enemies[i].damage)
                player.score += 10
                enemies[i].respawn()
                sound.beep(880, 100)
            }
        }
        
        // Check game over
        if !player.isAlive() {
            print "Game Over! Score: {player.score}"
            running = false
        }
    }
    
    void render() {
        screen.clear(color.navy)
        
        player.render()
        
        for i = 0 to 4 {
            enemies[i].render()
        }
        
        // UI
        screen.text(10, 10, "Health: {player.health}", color.white)
        screen.text(10, 30, "Score: {player.score}", color.white)
        
        screen.flip()
    }
    
    void run() {
        screen.init(SCREEN_WIDTH, SCREEN_HEIGHT, "Game")
        
        while running {
            update()
            render()
            sleep(16)  // ~60 FPS
        }
        
        screen.close()
    }
}

// ===== Main =====
int main() {
    game = Game()
    game.run()
    return 0
}
```

## Key Features

### From BASIC:
- `for i = 0 to 10` loops
- `print` statements
- String interpolation with `{}`
- Simple variable declarations
- `for each` loops

### From C:
- Pointers when needed
- `malloc`/`free` for manual memory
- Fast, compiled performance
- Direct hardware access
- Structs

### Modern Features:
- Smart pointers (`new`, automatic cleanup)
- References (`ref`)
- Nullable types (`maybe`)
- Shared ownership (`shared`)
- Methods in structs
- Constructors
- Simple inheritance with `:` and `override`
- Interfaces
- Generics/Templates
- Properties with `get`/`set`
- Operator overloading
- Lambdas `(x) => x * 2`
- Array methods `.map()`, `.filter()`, `.sort()`
- String methods `.upper()`, `.split()`
- `defer` for automatic cleanup
- Dynamic arrays
- Pattern matching
- Error handling (`try`/`catch`, `Result<T,E>`)
- Async/await
- Modules and imports
- Attributes/decorators
- Reflection

## Design Principles

1. **Simple by default** - Basic things are easy
2. **Powerful when needed** - Advanced features available
3. **Safe but not restrictive** - Smart pointers by default, raw pointers when needed
4. **Fast** - Compiles to native code
5. **Modern** - Has features developers expect
6. **Clean syntax** - Easy to read and write
7. **Practical** - Built for real applications and games
