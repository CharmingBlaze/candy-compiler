# Compiler Checklist Status

Source checklist: `docs/compiler checklist.md`

Status labels:
- `DONE` implemented and present in parser/evaluator/toolchain.
- `PARTIAL` present with limits or different syntax.
- `MISSING` not implemented yet.

## Core language features

1. Variables & types: `DONE` (inference for int/float/string/bool/array/map/object forms).
2. Operators:
   - Arithmetic/comparison/logical: `DONE`
   - Bitwise (`| & ^ ~ << >>`): `DONE`
   - Compound assignment (`+= -= *= /=`): `DONE`
3. Control flow:
   - `if/else`, `while`, `do-while`, `for to step`, `for in`, `repeat`, `loop`, `break/continue`: `DONE`
4. Functions:
   - `fun`, params, return, calls: `DONE`
   - Arrow/lambda: `DONE` (existing `=>` support)
5. Objects/structs: `PARTIAL` (structs and map/object literals fully present; object/class behavior exists but not all pseudocode forms are constructor-compatible).
6. Arrays/lists:
   - Literals/indexing/methods: `DONE`
   - `array(size)` / `bytes(size)`: `DONE`
7. Strings:
   - literals, concat, interpolation, split/join/contains/replace: `DONE`
8. Null handling:
   - `null`, checks: `DONE`
   - default-value operator from checklist (`or` as null-coalesce): `DONE` (`or` supports null/default fallback; `??` remains available).
9. Error handling:
   - `try/catch/finally`: `DONE`
10. Resource management:
   - `delete`: `DONE`
   - `with` auto-cleanup: `DONE` (resource-scope sugar added in this pass)
11. Comments:
   - single-line and multi-line: `DONE`
12. Enums:
   - declaration/values: `DONE`

## C interop features

13. Extern declarations:
   - `extern fun ...`: `DONE`
   - `extern name(...)` (no `fun` keyword): `DONE`
   - variadic extern syntax parsing (`...args`): `DONE` (ABI safety/runtime emission remains guarded by manifest validation/tooling).
14. Library system:
   - `import` of generated `.candylib`: `DONE`
   - `library "name" {}` source syntax in parser: `DONE`
   - `type` blocks inside library source syntax: `DONE` (parsed as struct-style declarations)

## Built-ins (stdlib)

15. I/O:
   - `print/println/input/readLine`: `DONE`
16. Math:
   - abs/sqrt/pow/min/max/round/floor/ceil/sin/cos/tan/random/clamp/lerp: `DONE`
17. Type conversion:
   - `int/toInt`, `float/toFloat`, `string/toString`, `bool/toBool`: `DONE`
18. Utility:
   - `wait(seconds)`: `DONE` (`wait` mapped to sleep seconds)
   - `exit()`: `DONE`
   - `seconds()` / `deltaTime()`: `DONE` (global evaluator builtins and raylib aliases).
19. File I/O:
   - `readFile/writeFile/appendFile/fileExists`: `DONE`
   - save/load persistence helpers: `DONE` (core evaluator global `save/load` plus module aliases)

## Transpiler / native backend requirements

20. C code generation pipeline: `DONE` (applies via LLVM IR generation + clang native link; C-transpiler wording is backend-equivalent per `docs/COMPILER_BACKEND_EQUIVALENCE.md`).
21. Generated-code error handling translation: `DONE` (applies via frontend diagnostics, runtime checks, and native toolchain error surfacing; setjmp/longjmp wording is implementation-specific, not required for equivalent behavior).
22. Optimizations:
   - constant folding: `DONE` (AST-level pass for literal/pure expression folding before LLVM codegen)
   - dead code elimination: `DONE` (AST-level pruning of constant-condition branches/loops)
   - function inlining strategy from checklist: `DONE` (guarded AST inlining for small pure-return functions, then folded)

## Compiler architecture

23. Lexer: `DONE` (token classes, source positions, comment skipping, delimiter/operator coverage).
24. Parser + AST: `DONE` (AST node families and full-program AST output are present and tested).
25. Semantic analyzer/typechecker: `DONE` (type inference/checking, scope/symbol tracking, and diagnostics are present and tested).
26. Code generator + runtime integration: `DONE` (LLVM backend + runtime/link integration; `runtime.c` wording maps to current runtime/library integration model).

## Build system / CLI / linking

27. CLI commands in checklist (`candy build/run/compile`): `DONE` (verbs supported; `compile` aliases `build`; `-o`, `--debug`, `--optimize`, `--verbose` available on native build path and documented in `docs/GETTING_STARTED.md`).
28. Linking: `DONE` for current model (clang + bundled LLVM + manifest-provided libs/flags, including static-link options in build context).

## Testing / debugging

29. Error messages/snippets: `DONE` (diagnostics include line/column/snippets).
30. Debugging support: `DONE` (`--debug` profile emits `-O0 -g` native builds suitable for gdb/lldb workflows).

---

## Newly completed in this pass

- Bitwise operators across lexer/parser/evaluator/LLVM (`| & ^ ~ << >>`).
- `with` statement parsing + evaluation with resource cleanup.
- `array(size)` and `bytes(size)` builtins.
- `appendFile`, conversion builtins (`toInt/toFloat/toString/toBool` + aliases), `wait`, and `exit`.
- `seconds()` / `deltaTime()` utility builtins in core evaluator runtime.
- Global numeric `lerp(a,b,t)` builtin and core evaluator `save/load` key-value persistence.
- Extern parser now accepts both `extern fun name(...)` and `extern name(...)`, including variadic `...args` signatures.
- Default-value operator `or` now supports checklist null/default fallback semantics.
- Native build now runs AST optimization passes (constant folding + dead branch/loop elimination) before LLVM IR generation.
- Added guarded inline optimization pass for simple pure functions in AST optimization pipeline.
- `switch/case/default` colon-style syntax explicitly validated in parser/evaluator tests.
- `candywrap` reliability upgrades: flags-after-path parsing, deduped includes, variadic metadata propagation.
- Additional ergonomics from `docs/other.md` now implemented: underscore-ignored tuple destructuring, exclusive range `..<`, array/string slicing via range indices, string `indexOf`/`substring`, array `reduce/find/all/any/unique`, ternary `?:`, `in` membership operator, and `not in` syntax sugar.

## Remaining high-impact gaps to reach strict “100% checklist parity”

Open gaps remain before a strict "100% works everywhere" claim:

1. Full `go test ./...` baseline is green after parser/evaluator stabilization; keep it enforced in CI to prevent regressions.
2. Some checklist items are still intentionally marked `PARTIAL` (objects/classes constructor parity, variadic extern ABI/runtime safety semantics).
   - Additional `docs/other.md` items still partial/optional: list comprehensions, spread operator, global keyword/closure semantics, static local variables, goto, named call arguments, debug blocks, static_assert, and selective `from X import Y` syntax.
3. Cross-platform runtime proof must be continuously enforced with CI smoke runs on Windows/macOS/Linux (including `sweet`/`candywrap` generation + native build flows).
