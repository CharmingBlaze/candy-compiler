# Candywrap Workflow: Single-Command Integration

This document explains the team contract for turning external C/C++ libraries into first-class Candy modules.

## Core Idea

Candywrap follows a **headers-in, ecosystem-out** model:

- Input: C/C++ headers (plus optional config/flags)
- Output: a portable module folder containing manifest, glue, and docs

The critical boundary is the generated `.candylib` manifest.

## Responsibility Split

### Candywrap responsibility

- Parse headers (`auto`, `libclang`, `regex`)
- Transform API names safely (namespace/strip-prefix/ignore/collision handling)
- Generate package artifacts:
  - `<lib>.candylib`
  - `<lib>_glue.c`
  - optional `<lib>.candy`
  - optional `<lib>_namespace.candy`
  - optional `<lib>.md` docs
- Enforce/annotate ABI guardrails during generation

### Compiler responsibility

- Load and validate `.candylib` at import/build time
- Synthesize `extern` declarations from manifest metadata
- Extend native build/link context from manifest compile/link fields
- Fail fast on unsupported ABI metadata with clean diagnostics

The compiler should stay C/C++-agnostic: it consumes manifest contracts, not raw headers.

## Single-Command Cycle

Example:

```bash
candywrap wrap ./foo.h --name foo --output ./candy_modules/foo --docs --stub
```

Optional advanced flags:

```bash
candywrap wrap ./foo.h --name foo --output ./candy_modules/foo \
  --parser libclang --lang c++ --namespace foo --strip-prefix foo_ --cxx-shim
```

Result:

- The output folder is a drop-in Candy module package.
- User imports the generated `.candylib` in source.
- `candy build`/`compile` handles linking through manifest metadata.

## Module Folder Standard (Portability Rule)

Treat `--output` as a package root that should be checked into source control when shipping a project.

Recommended convention:

```text
project/
  candy_modules/
    foo/
      foo.candylib
      foo_glue.c
      foo.candy              (optional)
      foo_namespace.candy    (optional)
      foo.md                 (optional)
```

Why:

- Keeps glue/docs/manifest in sync
- Avoids hidden machine-local binding state
- Makes project sharing reproducible

## ABI Guardrails Are a Feature

Early ABI rejection is intentional safety, not a limitation.

Instead of failing late with obscure linker/runtime errors, the toolchain rejects risky unsupported signatures early (or requires explicit `--unsafe-abi` flow + shim), with actionable diagnostics.

This is a DX and reliability win.

## Docs-as-Discovery (LSP Hook Opportunity)

Because Candywrap emits docs consistently, generated markdown can be used as a discovery index.

Potential future enhancement:

- LSP hover for imported externs can pull function signatures/examples from generated `<lib>.md`.
- This provides first-class docs for third-party libraries with no custom per-library editor plugin.

## Team Rules

- Keep `.candylib` schema strict and versioned.
- Prefer deterministic output (golden tests in CI).
- Do not move C/C++ parsing logic into compiler build/link path.
- Keep compiler-side manifest validation comprehensive and readable.

## User-facing one-line summary

"Point Candywrap at your headers; it generates manifest, glue, and docs, and your library becomes importable as a native Candy module without manual binding code."

