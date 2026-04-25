# Candywrap / Sweet (MVP+)

`sweet` is the C/C++ converter CLI name.
`candywrap` remains fully supported as a compatible command.

`sweet`/`candywrap` generates Candy-native binding metadata from C/C++ headers so `candy build` / `candy compile` can import those bindings and forward compile/link flags.

Team contract and architecture framing:
- See **[CANDYWRAP_WORKFLOW.md](CANDYWRAP_WORKFLOW.md)** for the Single-Command Integration model, compiler-vs-candywrap responsibility split, and module-folder portability standard.

## What it generates

For `sweet convert --name mylib --output ./bindings mylib.h`, the tool writes:

- `mylib.candylib` (JSON manifest)
- `mylib_glue.c` (glue source placeholder/shims)
- `mylib.candy` (optional helper stub)
- `mylib_namespace.candy` (optional namespace convenience wrappers when `--namespace` is set)

## CLI surface

`sweet convert` supports:

- `--name` library name
- `--output` output directory
- `--include` include directories (comma-separated)
- input args can be files, directories, and wildcard patterns (`*.h`, `*.c`, `*.cpp`)
- `--all` scan provided roots and wrap the whole library automatically
- `--root` explicit root directories/files for whole-library scan
- `--define` preprocessor defines (comma-separated)
- `--profile` built-in starter profiles (`raylib`, `sqlite`, `curl`)
- `--config` YAML config path
- `--parser` parser engine (`auto`, `libclang`, `regex`)
- `--lang` header language (`c`, `c++`)
- `--cxx-std` C++ language standard for compile flags (`c++17`, `c++20`, ...)
- `--namespace` namespace prefix for generated externs (`ns_function`)
- `--strip-prefix` remove C symbol prefixes before generating Candy names
- `--ignore` ignore functions by regex or wildcard patterns
- `--unsafe-abi` include risky signatures (variadics/function-pointers) in manifest for manual shim workflows
- `--static` request fully static native linking where supported
- `--static-lib` force specific libs to link via static mode
- `--link-lib`, `--link-lib-dir`, `--link-ldflag` add explicit link metadata
- `--simple` stricter safe subset for extern generation
- `--smart` enable convenience wrapper generation path
- `--stub` emit helper `.candy` stub
- `--docs` emit `<library>.md` API documentation
- `--cxx-shim` emit `<lib>_cxx_shim.cpp` starter template for C++ bridge wrappers

## `.candylib` schema (MVP)

- `library`: logical library name
- `namespace` (optional): namespace prefix applied to generated Candy extern names
- `headers`: source headers used during generation
- `generated_at`: RFC3339 generation timestamp
- `externs[]`: exported function metadata
  - `name`: Candy-visible function name
  - `symbol` (optional): C symbol override
  - `return_type`: C return type
  - `params[]`: `{ name, type }`
  - `variadic` (optional): rejected by MVP ABI guardrails
- `compile`:
  - `glue_sources[]`
  - `include_dirs[]`
  - `cflags[]`
- `link`:
  - `lib_dirs[]`
  - `libs[]`
  - `ldflags[]`
  - `static` (optional): request static linking
  - `static_libs[]` (optional): libraries to force in static mode
- `platforms` (optional):
  - `windows` / `linux` / `darwin` overlays for `compile` and `link`
- `types[]`, `enums[]`, `constants[]` metadata extracted from headers (best effort)

## Import/build flow

1. Candy source imports a `.candylib` path.
2. Import expansion loads and validates the manifest.
3. `externs[]` are synthesized into `extern fun ...` AST declarations.
4. Native build (`candy build` / `candy compile`) extends `clang` args with:
   - glue sources
   - include directories / compile flags
   - library directories / `-l` libs / linker flags

## ABI limits (MVP)

The MVP accepts safe primitive-oriented signatures and rejects:

- Variadic functions
- Function pointers in params/returns
- Complex ABI forms that cannot be lowered safely yet

Diagnostics are emitted both in bindgen filtering and compiler/typecheck validation.

This gives broad C-library coverage for common API shapes. For hard ABI edges (C varargs, function pointers, heavy macro APIs, C++ templates/classes), use `--unsafe-abi` plus generated C/C++ shim files to bridge into stable C wrappers.

## Naming and namespace strategy

To make wrappers easier to use across large C APIs:

- Keep original C symbol in `symbol`.
- Generate Candy-safe `name` values.
- Optionally strip long prefixes and add a namespace prefix:
  - C symbol `b2World_CreateBody`
  - with `--strip-prefix b2World_ --namespace box2d`
  - generated Candy name: `box2d_CreateBody`

If namespace is set and stubs are enabled, Sweet also emits `*_namespace.candy` with short helper wrappers:

- generated extern: `box2d_CreateBody(...)`
- convenience wrapper: `CreateBody(...)`

## Example

```bash
cd compiler
go run ./cmd/sweet convert --name mylib --namespace mylib --output ../examples/bindgen/out ../examples/bindgen/mylib.h
go run ./cmd/candy build ../examples/bindgen/main.candy
```

Config-driven run:

```bash
go run ./cmd/sweet convert --config ./candywrap.yaml
```

Prefer libclang when available:

```bash
go run -tags libclang ./cmd/sweet convert --parser libclang --lang c++ --name raylib --output ./bindings ./raylib.h
```

Box2D-style namespacing example:

```bash
cd examples/bindgen/box2d
../../compiler/sweet convert --config ./candywrap.yaml --output ./out
```

Convert all C/C++ files from a folder:

```bash
sweet convert --lang c++ --parser libclang --name mylib --namespace mylib --output ./out ./third_party/mylib
```

Or using wildcard patterns:

```bash
sweet convert --name mylib --output ./out "./third_party/mylib/*.h" "./third_party/mylib/*.cpp"
```

Whole-library one-command mode:

```bash
sweet convert --all --root ./third_party/mylib --name mylib --namespace mylib --output ./out
```

Nuklear-style end-to-end example in this repo:

```bash
cd compiler
go run ./cmd/candywrap wrap --config ../examples/bindgen/nuklear/candywrap.yaml --output ../examples/bindgen/nuklear/out --docs --stub
go run ./cmd/candy build ../examples/bindgen/nuklear/main.candy -o ../examples/bindgen/nuklear/out/nuklear_demo
../examples/bindgen/nuklear/out/nuklear_demo
```

Expected output includes:
- `[nk] init`
- `[nk] begin ...`
- `GUI event: button clicked`
- `press count = 1`

## C++ libraries

`sweet` can parse many C++ headers best with `--parser libclang --lang c++`, but stable cross-platform interop still works best through a C ABI boundary:

1. Generate normal manifest/docs/stub.
2. Add `--cxx-shim` to emit a starter `*_cxx_shim.cpp` template.
3. Implement each `extern "C"` wrapper in that shim to call into C++ classes/APIs.
4. Link resulting object/library through the generated `.candylib`.

Recommended C++ flags for conversion:

```bash
sweet convert --lang c++ --parser libclang --cxx-std c++20 --cxx-shim ...
```

Example files:

- `examples/bindgen/mylib.h`
- `examples/bindgen/main.candy`

