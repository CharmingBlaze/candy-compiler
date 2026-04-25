# Beginner-Friendly SimpleC: Vision and Candy (This Repo)

This document ties the **educational goal** (a language as easy as BASIC with C’s power when you need it) to the **Candy** implementation. Candy is a **modular, gradually typed** interpreter plus an optional **LLVM** backend (`candy build` / `candy compile`).

- **What this is:** A *single map* of “teaching / DX ideas” → **status in the dynamic interpreter** vs **native (LLVM) path** vs **not there yet**, so you can plan lessons and set expectations. The **full pitch-level checklist** (one place for every “Level 1 / 2 / 3” and stdlib item) is [VISION_ROADMAP.md](VISION_ROADMAP.md).
- **What this is not:** A promise that every line of aspirational code in long teaching docs runs unchanged; where Candy differs, the tables and [VISION_ROADMAP.md](VISION_ROADMAP.md) say so.

---

## Core philosophy (vision)

- **Level 1 — Complete beginner (dynamic only):** write like BASIC/Python: assign, `print`, arrays, `for`/`while`, functions, and structs without type annotations. Run with `candy` (no build).
- **Level 2 — Structure:** struct fields, methods, instances.
- **Level 3 — Types when you care:** add `: type` and return types for optimization and for **`candy build` / `candy compile`**.

Candy is explicitly designed to **start without types** and add them when you want.

---

## How to read the status table

| Column | Meaning |
|--------|--------|
| **Interpreter** | The tree-walking runtime (`candy_evaluator`), `candy run` (or `candy` default execution). |
| **LLVM** | `candy build` / `candy compile` via `candy_llvm`; only a **subset** of the AST, with errors for unsupported features. |
| **Parser** | Often accepts more than the interpreter or backend run today. |

### Legend

- **Yes** — practical for the intended use.
- **Partial** — some forms work; others are parse-only or have gaps (see “Notes”).
- **No** — not implemented in that layer, or not a language feature in Candy.
- **N/A** — the concept does not apply.

---

## Feature map (vision → Candy)

| Topic | Vision / teaching idea | Interpreter | LLVM | Notes |
|-------|-------------------------|------------|------|--------|
| **No mandatory types** | `name = "Alice"`, `score = 100` | **Yes** | For typed subset | `val` / `int x =` etc. in grammar; dynamic path uses `=` heavily. |
| **Print** | `print` / `println` | **Yes** | Varies by printf subset | String interpolation: `"Hello, {name}!"` works. |
| **String interpolation** | `"{x}"` in strings | **Yes** | **Partial** | As documented in `LANGUAGE.md`. |
| **Arrays** | `[1,2,3]` | **Yes** | **Partial** | `a + [x]` concatenates; `a.add(…)` / `push` / `append` (in place); `map`/`filter` with `(x) =>` lambdas; `sort`/`reverse`/`sum`/`min`/`max` — see [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md#8-list--array--methods-instance-calls). |
| **Int ranges** | `0..10`, `3..0` (counting / backwards) | **Yes** | **Partial** | `a..b` evaluates to a **contiguous int array, inclusive** on both ends. Use with `for v in 0..10 { … }`. The `for i in 0..7 by 2` style is **not** a separate form (use `for i = 0 to 6 step 2` for integer step). |
| **Negative index** | `a[-1]` (last) | **Yes** | **Partial** | For strings, index is by **rune** (works well with non-ASCII). |
| **Maps / dicts** | `{"a": 1}` | **Yes** (interpreted) | **Partial** | `for k in map` in interpreter iterates **string keys** (BASIC-style loop over a map, not `for name, score`). |
| **Functions** | `func` / `fun` / `function` / `sub` | **Yes** | **Partial** | All map to the same `FUNCTION` token. |
| **Default parameters** | `f(w = 800, h = 600)` | **Parser** (often) | TBD | Confirm per form in real samples; not guaranteed in eval+LLVM. |
| **Arrow / lambdas** | `(x) => x * 2` | **Yes** | TBD | Evaluated as a real function value (closure); use in calls, `a.map(…)`, `a.filter(…)`. |
| **Control: `if` / `else` / `else if`** | | **Yes** (statement) | **Partial** | Ternary `a ? b : c` is **not** Candy as of this writing; use an `if` *statement* or a small helper. |
| **`if` as an expression (value)** | | **No** in eval | TBD | `if` in expression position is not a general value-yielding form in the interpreter today. |
| **Switch** | `switch` / `case` / `default` | **Not in interpreter** | **Partial** / gaps | May parse; eval needs explicit support. |
| **Loops: `for v in x`** | | **Yes** | Gaps in LLVM | `for v in` iterable: array, string, map keys, or an **int range** (`0..2`). Bodies get a **fresh env per iteration**; outer variables are **not** reassigned; use `while` to accumulate in one outer scope. |
| **Loops: `for i = a to b [step s]`** | | **Yes** | **Partial** | Integer inclusive range. |
| **C-style `for(;;)`** | | **No** in interpreter | TBD | Parsed elsewhere; not wired in `candy_evaluator` switch. |
| **While** | | **Yes** | **Partial** | |
| **Do-while** | | **No** in interpreter | TBD | |
| **Structs, fields, methods** | C-style in body | **Yes** (read `LANGUAGE.md`) | **Partial** | Struct fields in the parser are often `name: type` (or declared forms); the “untyped `name` only” field from pure teaching examples may need adaptation to Candy’s current grammar. |
| **Method calls** | `o.move(1,2)` | **Yes** | **Partial** | C-style `int move(...)` in struct, etc. |
| **Inheritance / `using` / embedding** | | **Parser** / typecheck in places | TBD | Not a complete beginner “works everywhere” story in eval. |
| **Try / catch / finally** | | **Yes** in interpreter (dynamic) | **Not** in LLVM (yet) | Catch type is a hint; the bound name receives the **string** of the error. No typed exception hierarchy. |
| **Defer** | | **Not in interpreter** (parsed only) | TBD | [LANGUAGE.md](LANGUAGE.md#defer-parsed-not-run-in-the-interpreter) — use **try**/**finally** for now |
| **`open`, game APIs, `key()`, `window()`** | | **No** in default `candy` | **No** | Optional **Raylib** builtins when built with **`-tags raylib`** — [EXTENSIONS.md](EXTENSIONS.md); otherwise treat game samples as a **teaching target**, not shipped I/O. |
| **Null / optional** | | **`??` Yes**; **`?.` not** full in eval | TBD | [LANGUAGE.md](LANGUAGE.md#null--and-safe-navigation) — `null` checks, **`a ?? b`**; use **`if` + dot** instead of `?.` |
| **Logical operators** | | `&&`, `||` | **Partial** | English words `or` / `and` in place of `||` / `&&` are not necessarily keywords. |
| **“Constant by uppercase”** | | **Not** a separate rule in Candy | N/A | Use `const` / `val` for immutability where supported. |
| **Memory model** (RC, `malloc`/`free`, raw pointers) | | **No** the teaching story is not the runtime; Go GC backs the Go VM | N/A | [LANGUAGE.md](LANGUAGE.md#pointers-memory-and-c-style-features) — no `&` / `*` / `->` in the interpreter `Value` model |
| **Semicolons** | | **Optional in many cases** (lexer can insert) | N/A | |
| **Import file / stdlib** | `import "a.candy"`, `import math` | **Yes** | **Partial** (LLVM) | Prelude also injects `math`, `file`, `json`, …; see [LANGUAGE.md](LANGUAGE.md#modules-imports-and-standard-library). |
| **Selective / `import m as x`** | | **No** | TBD | Use `math.sqrt` or assign `m = math`. |
| **`module` / `export` (user package syntax)** | | **No** in eval | TBD | Treated as aspirational in [LANGUAGE.md](LANGUAGE.md#aspirational-module-syntax-not-copy-pastable-in-the-repl). |
| **`printf` vs `print`** | C-style `printf` | **No** | TBD | Use `print` / `println` and `"{x}"` interpolation. |

---

## Aspirational samples vs Candy today (short)

- **List sugar:** `items.add(6)`, `.map(…)`, huge method suites → **not** the dynamic runtime today. Prefer **`a = a + [x]`** and plain loops, or pre-build literals.
- **BASIC-like `for i, v in list` or `for k, v in map`** → **not** a single form yet; use a counted `for` / index or `for k in m` and `m[k]`.
- **Slicing `a[1..3]`** → not evaluated as a slice op; use `range` and indexes.
- **Ternary** → use **`if` statement** or a small `func` / lambda.
- **The full game at the end of the vision** → use as a design compass; you still need a host (graphics/input) that does not ship inside `candy` itself.

---

## Reference docs in this repo

| File | Role |
|------|------|
| [LANGUAGE.md](LANGUAGE.md) | Deeper language + implementation notes, gaps, and examples. |
| [README.md](README.md) | How to run `candy`, native build/compile flow, and the LSP. |
| [VISION_ROADMAP.md](VISION_ROADMAP.md) | One checklist for the “SimpleC for beginners” pitch vs Candy (levels, loops, stdlib, gaps). |
| [`.cursorrules`](.cursorrules) | Project/module layout and conventions for contributors. |

---

## Implementation strategy (from the vision) — where the repo is

- **Phase 1 — Dynamic runtime (Value bag, hash maps, growable arrays):** The interpreter follows this: `candy_evaluator` + `candy_token` / `candy_lexer` / `candy_parser` / `candy_ast`.
- **Phase 2 — Type inference (specialized paths):** Ongoing; see `candy_typecheck` and the LLVM path.
- **Phase 3 — Optional static typing + LLVM:** `candy_llvm`, `candy build` / `candy compile`, env vars in `README.md`.

The **“Key advantages”** checklist in the original pitch (BASIC-style loops, interpolation, structs, etc.) is largely aligned with the **interpreter** column above, with the important caveat that **library sugar and host I/O** are the biggest gaps, not the existence of a parser.

---

*Last updated with: **string** and extended **list** instance methods (`.add`, `.map`/`filter` with lambdas, `sort`/`sum`/…), **`random.seed`**, [VISION_ROADMAP.md](VISION_ROADMAP.md), and the items in [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md) (I/O, math, maps, `range`, `try`/`catch`, error hints, REPL, …).*
