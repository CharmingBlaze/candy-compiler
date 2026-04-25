# Candy Architecture

This document explains the project architecture for engineers onboarding to Candy.

## Overview

Candy is a Go-based language toolchain with:

- A frontend: lexer -> parser -> AST -> typecheck/diagnostics
- An interpreter runtime for fast execution and scripting workflows
- A native backend: AST optimizations -> LLVM IR -> clang link
- C/C++ interop tooling (`candywrap` / `sweet`) based on `.candylib` manifests
- Editor tooling through a minimal LSP implementation

## Repository Layers

### Frontend

- `compiler/candy_token`: token types and keyword/operator constants.
- `compiler/candy_lexer`: lexical analysis and semicolon insertion behavior.
- `compiler/candy_ast`: AST nodes, generally one file per node family.
- `compiler/candy_parser`: Pratt/TDOP parser with precedence-driven expression parsing.
- `compiler/candy_typecheck`: semantic/type analysis and diagnostics.
- `compiler/candy_report`: structured error reporting.

### Runtime

- `compiler/candy_evaluator`: AST interpreter, value model, environments, builtins, control flow.
- `compiler/candy_stdlib`: built-in module source packages.
- `compiler/candy_raylib`, `compiler/candy_physics`, `compiler/candy_enet`: host-integrated modules and game/network helpers.

### Native Backend

- `compiler/candy_opt`: AST-level optimization passes (constant folding, dead branch elimination, guarded inlining).
- `compiler/candy_llvm`: LLVM IR generation and optimization pipeline integration.
- `compiler/cmd/candy`: CLI orchestration for parsing, checking, running, and building.
- `compiler/candy_load`: import expansion and build context synthesis for native build/link.

### Interop Tooling

- `compiler/candy_bindgen`: C/C++ API extraction and generation of manifests/glue/docs.
- `compiler/cmd/candywrap`: wrapper generator CLI (`wrap` command).
- `compiler/cmd/sweet`: converter/wrapper CLI (`convert` / `wrap` style).

### Tooling

- `compiler/candy_lsp`: diagnostics-first language server.
- `compiler/scripts`: bundling, smoke, and performance helper scripts.

## Data Flow

### Interpreter Path (`candy run`)

1. Source is read.
2. Lexer produces tokens.
3. Parser builds AST.
4. Optional static check runs.
5. Evaluator executes AST directly.

### Native Build Path (`candy build` / `candy compile`)

1. Entry file is parsed and imports are expanded.
2. `.candylib` imports are loaded and extern declarations are synthesized.
3. Build context is gathered (include dirs, libs, flags, glue sources).
4. Toolchain preflight validates `clang` availability (and reports `opt` status).
5. AST optimization pass runs.
6. LLVM IR is generated.
7. Optional `opt` pass pipeline is applied by build profile.
8. IR is compiled/linked with `clang` into a native binary.

## Parser Design

Candy uses a Pratt parser:

- Prefix and infix parse functions are registered by token type.
- A precedence table drives associativity and parse order.
- Expression extension is localized and predictable (add token, precedence, parse fn registration).

This design keeps language evolution fast for operators and expression-level features.

## Runtime Model

The evaluator uses:

- A tagged `Value` model (`null`, numeric, string, bool, array, map, function, struct/object, module, result).
- Lexical `Env` chains for scope and closures.
- Control-flow wrappers for `return`, `break`, `continue`.

This keeps execution semantics explicit and testable at AST level.

## Typechecking Model

Typechecking runs as AST walks with scoped symbol tables:

- Declaration collection and statement/expression checking
- Type inference where possible
- Assignability and structural checks across supported constructs
- Diagnostic emission with source locations

Current behavior emphasizes practical diagnostics and compatibility with active language evolution.

## Optimization Strategy

Candy uses a two-stage optimization strategy:

1. AST-level passes in `candy_opt`:
   - Constant folding
   - Dead branch/loop elimination
   - Guarded simple-function inlining
2. LLVM `opt` pass pipeline in `candy_llvm`:
   - Profile-driven (`debug`, `dev-release`, `shipping`)

This balances fast semantic cleanup with backend-level code quality.

## C/C++ Interop Strategy

`candywrap` / `sweet` generate:

- `<lib>.candylib` manifest
- `<lib>_glue.c`
- Optional `.candy` stubs and namespace helpers
- Optional docs markdown and C++ shim template

At build time, `.candylib` imports are resolved into:

- `extern` AST declarations
- Native compile/link metadata used to extend clang arguments

This keeps the core compiler generic while enabling broad external library support.

## Build Profiles

Build behavior is profile-driven:

- `debug`: low optimization, debug symbols
- `dev-release`: balanced optimization
- `shipping`: aggressive optimization

Profiles affect both LLVM pass selection and clang flags.

## Toolchain Diagnostics

Use `candy doctor` to run native toolchain checks without compiling source files.

- Verifies `clang` resolution using the same search order as native build.
- Reports `opt` availability/version for IR optimization support.
- Prints actionable remediation guidance when required tools are missing.

## LSP Capabilities

The LSP server currently focuses on:

- Incremental diagnostics
- Hover info
- Go-to-definition
- Workspace symbol indexing from open documents

It is intentionally lightweight and robust for daily editing workflows.

## Testing and Performance

The project includes:

- Package-level unit tests across parser/evaluator/typecheck/backend/bindgen
- End-to-end tests for interop and native build flows
- Benchmark harnesses and CI perf gate scripts (`compiler/scripts/perf-gate.sh`)

## Design Principles

- Keep frontend modular and easy to extend.
- Preserve clear separation between interpreter and native build paths.
- Treat interop as metadata-driven, not hardcoded in compiler internals.
- Prefer deterministic diagnostics and reproducible build behavior.
