# “Truly Beginner-Friendly SimpleC” — vision vs Candy

This is the **single checklist** for the teaching pitch: *BASIC-easy, add C-style types when you want*, with **what the Candy interpreter does today** vs **LLVM** vs **not yet**. For API-level detail (I/O, `math`, list methods, etc.), see [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md). For how to run the tool, see [README.md](README.md). **Concrete syntax equivalents** (e.g. `x as int` teaching forms vs Candy’s `x: int` / `int x`, `sub`, fixed arrays, pointers) are summarized in [LANGUAGE.md](LANGUAGE.md) in *Variables - Simple and Clear* and *Functions - BASIC Style with Types* — use that for copy-pastable examples without duplicating the full matrix here. The same file has **Output - BASIC style**, **Control Flow**, and **Loops** for `print`/`printf`, ternary/`switch`, and `for each` / `do`-`while` / C-style `for` vs what the interpreter runs, plus **Arrays and Lists**, **Dictionaries (Maps)**, and **Strings** for list/map methods, `for k,v` vs keys-only, and `substring` / multiline / `sprintf` reality, and **Modules, imports, and standard library** for `import` / prelude, selective imports, and builtins vs C-style `printf`.

**Legend:** **Yes** = good for the stated use in the tree-walking runtime · **Partial** = some forms or limitations · **No** = not in that layer (or not a match to the sample syntax) · **N/A** = concept differs in Candy

---

## Core philosophy

| Idea | Interpreter | LLVM | Notes |
|------|------------|------|--------|
| Level 1: untyped `name = "Alice"`, `print`, math, lists, `for`/`while`, `func` | **Yes** | Subset | `print` / `println` |
| Level 2: structs, methods, `struct T { }` + `T { f: v }` | **Yes** | **Partial** | See [LANGUAGE.md](LANGUAGE.md) for field syntax; “bare” field names may need the forms the parser accepts |
| Level 3: `: type` and return types, `candy build` / `candy compile` | **Partial** | **Partial** | Typecheck + native path grow over time |

---

## Variables and assignment

| Vision | Interpreter | Notes |
|--------|------------|--------|
| `x = 1` (no types) | **Yes** | |
| `age: int = 25` (typed) | **Partial** | Parser / type path; dynamic path is assignment-heavy; see [LANGUAGE.md](LANGUAGE.md) *Variables* for `x as int` → Candy mapping |
| `MAX_HEALTH = 100` as “constant by capitals” | **No** as a rule | Use `const` / `val` where the grammar supports it; UPPERCASE is style only |

---

## Functions

| Vision | Interpreter | Notes |
|--------|------------|--------|
| `func` / `fun` / `function` | **Yes** | All lex to the same `FUNCTION` token; see [LANGUAGE.md](LANGUAGE.md) *Functions* for `:` vs teaching `as` and for `sub` |
| `sub` (no return) | **Yes** (parse + eval) | Treated like a function in the dynamic runtime |
| Default parameters | **No** in user `fun` in interpreter | Parser/AST does not carry defaults; pass every argument; see *Functions* in [LANGUAGE.md](LANGUAGE.md) |
| `(x) => expr` lambda | **Yes** | Produces a real function value; use with calls and `list.map` / `filter` |
| Multiple return / tuple unpack | **No** | Roadmap; use struct or list |

---

## Control flow

| Vision | Interpreter | Notes |
|--------|------------|--------|
| `if` / `else` / `else if` | **Yes** | |
| `a ? b : c` ternary | **No** | Use `if` statement or a small `func` |
| `if (expr) stmt` as one-liner | **Partial** | `if` is a **statement** in eval, not a general expression value |
| `switch` / `case` / `default` | **No** in eval | May parse; not evaluated in the interpreter path used for beginners today |
| `for v in x` (array, string, map **keys**, range) | **Yes** | Bodies use a **new scope per iteration**; outer `=` does not “accumulate” the loop var — use `while` for that |
| `for i, v in list` or `for k, v in map` | **No** | Single `for` var only; use `for k in m` and `m[k]` — [LANGUAGE.md](LANGUAGE.md) *Dictionaries (Maps)* |
| `for i in 0..10` / `for i in a..b by c` | **Partial** | `0..9` is an **inclusive** int **array**; `a..b by c` is not a dedicated form (use `for i = a to b step c` for integer steps) |
| `range(0, 10)` (end exclusive) | **Yes** | See [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md#5-ranges) |
| `for (i = 0; i < 10; i++)` C-style | **No** in eval | Parsed in places; not wired in the main interpreter `switch` |
| `while` | **Yes** | |
| `do { } while` | **No** in eval | |
| `or` / `and` as word operators | **No** (typically) | Use `||` / `&&` |

---

## Lists and values

| Vision | Interpreter | Notes |
|--------|------------|--------|
| `[1, 2, 3]`, `a + [x]` | **Yes** | Concatenation |
| `items.add(x)`, `push` / `append` | **Yes** | In-place on the list value |
| `insert`, `remove` (value), `remove_at` (index) | **Yes** | |
| `sort`, `reverse`, `sum`, `min`, `max` | **Yes** | `sort` uses a stable numeric/string order |
| `map((x) => …)` / `filter((x) => …)` | **Yes** | |
| Slicing `a[1..3]`, `a[..2]` | **No** in eval | Use index + `range` |
| List comprehensions `[x for x in …]` | **No** | Use `map` / `filter` or loops |
| Negative index `a[-1]` | **Yes** | String index is by **rune** |

---

## Maps / dicts and strings

| Vision | Interpreter | Notes |
|--------|------------|--------|
| `scores["a"] = 1` | **Yes** | String-keyed map in the VM |
| `m.keys`, `m.get`, `m.merge`, … | **Yes** | See [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md#9-map-dictionary--methods) |
| `for name, score in scores` | **No** | `for` over a map is **keys only** today |
| `"Hello {x}"` and `\\{` literal brace | **Yes** | |
| `text.upper()`, `split`, `length`, `contains` | **Yes** | `length` also as `text.length` (rune count) |
| `??` (nullish coalescing) | **Yes** | [LANGUAGE.md](LANGUAGE.md#null--and-safe-navigation) |
| `?.` (safe navigation) | **No** / **Partial** | Prefer `if` + field access; see [LANGUAGE.md](LANGUAGE.md#null--and-safe-navigation) |
| C pointers `&` / `*`, `->`, `malloc` / `free` / `sizeof` | **No** | GC `Value` model; [LANGUAGE.md](LANGUAGE.md#pointers-memory-and-c-style-features) |
| `or` as a word operator for defaults | **No** | Use `??` or `if`; see *Null* in [LANGUAGE.md](LANGUAGE.md) |

---

## Structs, errors, I/O, stdlib, tooling

| Topic | Interpreter | Notes |
|--------|------------|--------|
| `import "path.candy"`, `import math` (stdlib or identifier path) | **Yes** | [evalImport](candy_evaluator/eval.go) + [modules.go](candy_stdlib/modules.go); Cwd-relative file imports. |
| Prelude `math` / `file` / `json` / `random` / `time` without `import` | **Yes** | [registerPrelude](candy_evaluator/prelude.go) runs every [Eval](candy_evaluator/eval.go). |
| `import m (a, b)` selective / `import m as x` | **No** | Grammar has only string or ident path; [LANGUAGE.md](LANGUAGE.md#modules-imports-and-standard-library). |
| `module` / `export` blocks in user code | **No** in eval | Aspirational; not in `evalStatement` — [LANGUAGE.md](LANGUAGE.md#aspirational-module-syntax-not-copy-pastable-in-the-repl) |
| `printf` (C-style format) | **No** | `print` / `println` + string interpolation. |
| Methods `o.m()` | **Yes** | C-style struct methods in body |
| Embedding / `using` / “simple inheritance” | **Partial** | See typecheck and docs; not a full beginner “always works” story in eval |
| `try` / `catch` / `finally` | **Yes** | Catch name gets **string** error message; **not** in LLVM in general |
| `defer` | **No** in eval | Parsed only; [LANGUAGE.md](LANGUAGE.md#defer-parsed-not-run-in-the-interpreter) |
| File + JSON + math + random + sleep + time ms | **Yes** | [ESSENTIAL_ADDITIONS.md](ESSENTIAL_ADDITIONS.md) |
| `random.seed(n)` (module) | **Yes** | `random.seed` / re-seeds global `rng` |
| Game APIs: `window`, `key`, `circle`, `flip` | **No** (default build) / **Yes** with `-tags raylib` | Optional Raylib host extension — [EXTENSIONS.md](EXTENSIONS.md); default binary has no graphics builtins |
| `import "math"`, `import "file"`, … | **Yes** | `from m import a` is **not** implemented |
| Helpful “Did you mean …” for vars / struct fields | **Yes** | |
| `assert`, `debug` | **Yes** | `breakpoint` not present |
| REPL `:help`, `:vars`, `:exit` | **Yes** | |
| `enum`, package manager, operator overloading | **No** | Roadmap / out of scope for current VM |

---

## Memory model (as in the pitch)

Candy’s **interpreter** is implemented in **Go** with a GC. The long-term SimpleC story (RC, `malloc`/`free` in user code) is **not** the current runtime; teach that as a future *native* story, not a guarantee of the `.candy` REPL. Pointers, manual heap, **`defer` in user code**, and **`?.`** are summarized in [LANGUAGE.md](LANGUAGE.md#pointers-memory-and-c-style-features) and following sections; **`??`** and **`try` / `catch`** are implemented in the interpreter.

---

## Implementation strategy (phases) — this repo

| Phase | Status |
|--------|--------|
| **1** Dynamic `Value` runtime, hash maps, growable lists | **In progress** — `candy_evaluator` |
| **2** Type inference / specialization | **Ongoing** — `candy_typecheck` |
| **3** Optional static typing + LLVM | **`candy build` / `candy compile`** where supported |

---

*When documentation and code disagree, the **test suite and evaluator behavior** are authoritative; this file is the narrative bridge.*
