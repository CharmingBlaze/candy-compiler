# Candy — Essential Additions (Status)

This file tracks the **“Essential Additions”** list (I/O, math, stdlib, errors, lists, time, etc.) and what the **Candy interpreter** provides today. Native LLVM build (`candy build` / `candy compile`) only covers a **subset**; most of this is **dynamic** only unless noted. For **language** behavior that is not a stdlib name — **C-style pointers and `malloc` / `free` (not in eval)**, **`defer` (parsed, not run)**, **`try` / `catch` (string errors)**, **`null` and `??`**, **lambdas** — use the dedicated sections in [LANGUAGE.md](LANGUAGE.md) from **Pointers, memory, and C-style features** through **Lambdas**. For **`import` syntax, prelude vs `import` redundancy, and teaching-only `import m as x` / selective imports**, use **Modules, imports, and standard library** in [LANGUAGE.md](LANGUAGE.md); the tables below name **I/O and helpers** — the exact **`math.*` / `file.*` / …** method keys match [prelude.go](candy_evaluator/prelude.go).

---

## 1. File I/O

| API | Status | Notes |
|-----|--------|--------|
| `read_file` / `readFile` / `readfile` | **Yes** | Also `file.read`, `file.read_file` |
| `write_file` / `writeFile` | **Yes** | `file.write` |
| `read_lines` / `readLines` | **Yes** | Returns a **list of strings** (line breaks normalized) |
| `file_exists` / `fileExists` / `file.exists` | **Yes** | |
| `delete_file` / `file.delete` / `file.remove` | **Yes** | |
| `list_files` / `listDir` / `file.list` | **Yes** | Returns **file names** (not full paths), sorted |
| `load_json` / `json.load` / `json.load_file` | **Yes** | File path → value |
| `save_json` / `json.save` / `json.save_file` | **Yes** | Value → JSON file |
| `json_stringify` / `json_stringify` | **Yes** | `json.stringify` |
| `json_parse` / `json.parse` | **Yes** | String → value |

**Import:** `import "file";` (empty package body; **host prelude** still injects `file` / `fs`).

---

## 2. Math and constants

| API | Status | Notes |
|-----|--------|--------|
| `sqrt`, `abs`, `pow`, `floor`, `ceil`, `round` | **Yes** | Top-level and `math.*` |
| `sin`, `cos`, `tan` | **Yes** | Radians |
| `min` / `max` / `clamp` | **Yes** | `min`/`max` variadic in-host |
| `PI`, `E` | **Yes** | **Globals** and `math.PI`, `math.E` (also `pi` / `e`); `math.Inf` available as a float (positive infinity) |
| `import "math";` | **Yes** | Module object `math` with methods + consts

---

## 3. Random

| API | Status | Notes |
|-----|--------|--------|
| `random(lo, hi)` (int, **inclusive**) | **Yes** | Also `randomInt`, `rand_int` |
| `random_float(lo, hi)` **[lo, hi)** | **Yes** | |
| `choose(list)` / `sample` | **Yes** | `random.choice` on module |
| `import "random";` / `import "rand";` | **Yes** | `random` / `rand` module |
| `random.seed(n)` | **Yes** | Re-seeds the global PRNG (see `candy_evaluator` `builtinRngSeed`) |

---

## 4. String interpolation and escaping

| Feature | Status | Notes |
|---------|--------|--------|
| `"Hello {name}!"` | **Yes** | |
| `"{a + 1}"` | **Yes** | Expression in `{`…`}` |
| Literal brace | **Yes** | In source, **`\{`** in a string keeps a literal `{` in output (and does **not** start interpolation) |

`from str import` / selective imports: **not** implemented; use modules or globals.

---

## 5. Ranges

| Form | Status | Semantics (interpreter) |
|------|--------|-------------------------|
| `0..3` in expressions | **Yes** | Inclusive int array `0,1,2,3` |
| `range(n)` | **Yes** | `0` .. `n-1` (Python-style, end **exclusive**) |
| `range(a, b)` | **Yes** | `a` .. `b-1` step 1 |
| `range(a, b, step)` | **Yes** | Step may be **negative** |
| `0..=10` (extra token) | **Not** a separate form | Use `0..10` (inclusive) or `range(0, 11)` |
| Slicing `a[1..3]`, `a[..3]` | **Not** yet | Index with `a[1]`, or use `range` for indices |

---

## 6. Imports (standard library)

`import "math"`, `import "file"`, `import "json"`, `import "random"`, `import "time"`, and existing `std/…` keys in [candy_stdlib/modules.go](candy_stdlib/modules.go) resolve. **Name bindings** for `math`, `file`, `json`, `random`, `time` are **also** set by the interpreter **prelude**; imports can be used to satisfy style even when they add no new code.

`from m import a, b`: **not** implemented.

---

## 7. Better error messages (beginner-friendly)

| Case | Status |
|------|--------|
| **Undefined variable** | Suggests a **similar name** from the current environment chain (edit distance) |
| **Struct field typo** | **Error** (not `null`); “Did you mean `field`?” when a close match exists |
| | |

---

## 8. List / array — methods (instance calls)

On a **list** value, `a.method(…)`:

- `add(…)` / `push` / `append` — one or more values, **in place** (beginner teaching uses `items.add(6)`).  
- `insert(i, v)` — insert at index (negative index from end).  
- `remove(x)` — remove **first** value equal to `x` (by `valueEqual`).  
- `remove_at` / `delete_at` / `splice1` — remove by index.  
- `sort` / `reverse` / `sum` / `min` / `max`  
- `join(sep)` — string from elements’ `String()` with separator  
- `map(fn)` / `filter(fn)` — `fn` is a user function or lambda `(x) => …`  
- `contains(x)` / `include`  
- `index_of` / `index`  
- `is_empty` (also properties `.length` / `.size` / `.count` / `.is_empty` on the list **without** `()`)  
- `clear` (in-place)  
- `first` / `last`  
- `shuffle` (in-place)  

Slicing and list comprehensions: **not** in this list.

**Append idiom (still valid):** `a = a + [x]`.

---

## 8b. String — instance methods and properties

On a **string** value, `s.length` (rune count) and `s.is_empty`; methods:

- `upper` / `lower`  
- `split(sep?)` — default separator **space** if omitted  
- `contains` / `includes`  
- `starts_with` / `ends_with`  
- `trim` / `strip`  
- `replace(old, new)` — `strings.ReplaceAll`  

---

## 9. Map (dictionary) — methods

`map { … }` value; then `m.keys()`, `m.values()`, `m.has(…)`, `m.get(…[, default])`, `m.remove` / `delete`, `m.clear`, `m.merge(other)`.

---

## 10. Time / sleep

| API | Status |
|-----|--------|
| `sleep` / `sleepMs` (milliseconds) | **Yes**; capped to **60s** to avoid runaway programs |
| `sleep_sec` | **Yes**; same cap |
| `time_millis` / `timeMillis` / `time_ms` | **Yes**; wall-clock ms since epoch |
| `time` module: `time.sleep`, `time.millis`, … | **Yes** |
| `time()` as only **millis** | Use `time_millis` — a dedicated **`time` module** shadows a bare `time()` if we ever reintroduce it; **do not** use a global `time` **call**; use `time_millis` or `import "time"`. |
| `datetime` object / `now.year` | **Not** implemented |

---

## 11. Assert / debug

| API | Status |
|-----|--------|
| `assert(cond, "msg?")` | **Yes** |
| `debug(a, b, …)` | **Yes**; prints `Value.String()` |
| `breakpoint()` | **Not** (use `debug`) |

---

## 12. Enum, multiple return, REPL, package manager

- **REPL** — **Yes** (`candy -i` / `candy -repl`: `:help`, `:vars`, `:exit` / `:q`, persistent `Env` per session). One line per “program” for now; no multiline input buffer.
- **Enum**, **multi-return** / **tuple binding**, **candy.toml** / `candy add` — **not** implemented; still roadmap.
- **Lambdas** `(a, b) => expr` in expression position — **Yes** in the interpreter; needed for `list.map` / `list.filter` and higher-order call style.

---

## 13. Type helpers

| API | Status |
|-----|--------|
| `type` / `typeof` | **Yes** → `"int"`, `"float"`, `"string"`, `"bool"`, `"list"`, `"map"`, … |
| `is_int`, `is_string`, `is_list`, `is_map`, `is_float`, `is_bool` | **Yes** |

---

## 14. Comments

`//` and `/* */` are supported. `///` doc: treated like `//` (no special AST).

---

## 15. Standard organization

- **Host (Go) builtins and prelude** in `candy_evaluator` (`prelude.go`, `stdlib_builtins.go`, `builtin.go`, `eval_container_ops.go`, …).  
- **Stubs** for `import "math"`, `import "file"`, etc. in `candy_stdlib/modules.go`.  
- The document’s **separate** `.candy` file per area (`string.candy`, `list.candy`, …) is *not* how the current VM works; the **interpreter** wires the surface directly.

---

## 16. Quick reference — globals of note

**Always available (no import):** `print` / `println`, `len`, I/O/JSON/math/random/time/sleep helpers, `type`, `assert`, `debug`, `range(…)`, and module names injected by the prelude: **`math`**, **`file`**, **`fs`**, **`json`**, **`random`**, **`rand`**, **`time`**, plus **`PI`**, **`E`**, **`pi`**, **`e`**.

*Last updated: string and extended list instance methods, lambda evaluation, `random.seed`, `evalUserFunction` (shared with user calls and `map`/`filter`), plus the rest of the prelude/stdio surface documented above.*
