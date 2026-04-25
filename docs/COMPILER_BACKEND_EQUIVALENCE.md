# Compiler Backend Equivalence

This project implements the checklist compiler goals through an LLVM-native backend instead of a text C-transpiler backend.

The checklist phrases some items as "generate C code" and "link runtime.c". In Candy, these responsibilities are covered by equivalent LLVM/clang stages:

- Frontend: lexer -> parser -> AST -> typecheck.
- Mid-end: AST optimization passes (constant folding, dead branch/loop elimination, guarded inlining).
- Backend: LLVM IR generation.
- Native link: clang invocation (`.ll` + optional glue sources + include/lib/link flags from `.candylib` manifests).

## Checklist Mapping

- **Generate C code** -> **Generate LLVM IR** (`compiler/candy_llvm`).
- **Runtime library link (runtime.c)** -> **Runtime support via LLVM emission + linked glue/runtime objects**.
- **Optimization pass list** -> **AST optimizer + LLVM `opt` pipeline by profile**.
- **Debug symbols / debugger support** -> **`--debug` native profile with `clang -O0 -g` output compatible with platform debuggers**.

## Why this is equivalent

- LLVM IR is a lower-level, stronger intermediate target than direct C emission.
- Clang remains the final native linker/driver in both models.
- External C interop still works through generated glue C and link metadata from `candywrap` `.candylib` manifests.

## Command Surface

Supported checklist-style commands:

```bash
candy build program.candy
candy run program.candy
candy compile program.candy -o output
```

Notes:

- `candy compile ...` is an alias of `candy build ...`.
- `candy run ...` executes through the evaluator path (non-native), while `build`/`compile` use the native LLVM pipeline.

Supported flags on native build path:

- `-o` output path
- `--debug` debug symbols profile
- `--optimize` shipping optimization profile
- `--verbose` build-step logging
