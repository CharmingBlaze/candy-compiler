# Candy: Essential C-Interop Additions

This guide documents the essential language/runtime additions for practical C wrapper compatibility while keeping Candy simple.

## Included additions

- `null` checks for optional/failed C results.
- `try/catch/finally` for recoverable interop failures.
- `with name = expr { ... }` resource-scoped usage with automatic cleanup/nulling.
- String interpolation (`"Score: {score}"`).
- Lambda support (`(x) => x * 2`) and named `fun` callbacks.
- Bitwise operators for C-style flags:
  - infix: `|`, `&`, `^`, `<<`, `>>`
  - prefix: `~`
- Buffer helpers:
  - `array(size)` -> fixed-size Candy array (null-initialized)
  - `bytes(size)` -> fixed-size int array initialized with zeros (byte buffer usage)
- `.candylib` import + `extern` synthesis for native builds via `candywrap`.

## Opaque handle model

Pointers from C are intentionally treated as opaque handles in Candy user code. You pass them around and back into externs; no pointer arithmetic is exposed in normal Candy usage.

## What users write

```candy
import "mylib.candylib"

tex = LoadTexture("player.png")
if tex == null {
  print("load failed")
}

flags = FLAG_FULLSCREEN | FLAG_VSYNC
if (flags & FLAG_FULLSCREEN) != 0 {
  print("fullscreen")
}

with file = OpenFile("data.txt") {
  content = ReadAll(file)
  print("bytes: {len(content)}")
}
```

## Notes on current scope

- These additions target practical C interop ergonomics, not raw unsafe pointer manipulation.
- Complex ABI shapes (variadics/function pointers/unions/callback marshalling) still depend on wrapper strategy and are guarded by `candywrap` diagnostics where unsupported.
- For distribution and release bundling of `candy` + `candywrap`, see `docs/DISTRIBUTION.md`.
