# Candy

Candy is a modern language for games, scripting, and native app tooling.

It is designed to be:

- **Easy to read** for beginners
- **Fast to iterate** like a scripting language
- **Powerful enough** to compile native binaries and integrate C/C++ libraries

In short: Candy gives you a friendly syntax with a serious toolchain.

## Why use Candy

- **Beginner-friendly syntax** with a smooth path to advanced features
- **Game-first workflow** with built-in graphics/audio/input extensions
- **Native build path** for distributable binaries
- **C/C++ ecosystem access** via `candywrap`/`sweet` and `.candylib` manifests
- **Portable toolchain model** with self-contained release bundles

## Quick Start (Using Prebuilt Binaries)

Use the binaries in your release bundle (no compiler source setup required).

### 1) Check environment

```bash
candy doctor
```

### 2) Run a script

```bash
candy examples/candy/kid_moving_ball.candy
```

### 3) Build a native binary

```bash
candy build examples/candy/kid_moving_ball.candy
```

### 4) Wrap a C/C++ library (single-command integration)

```bash
candywrap wrap ./mylib.h --name mylib --output ./candy_modules/mylib --docs --stub
```

Then import the generated `.candylib` in Candy source and build as normal.

## Main Commands

- `candy run <file.candy>`: run via evaluator
- `candy build <file.candy>`: build native binary (LLVM/clang path)
- `candy compile <file.candy>`: alias of `build`
- `candy doctor`: verify toolchain health and bundled metadata
- `candywrap wrap ...`: generate Candy bindings from C/C++ headers
- `sweet convert ...`: alternate C/C++ conversion CLI compatible with candywrap flow

## Recommended Docs

- [Getting Started](docs/GETTING_STARTED.md)
- [Language Guide](docs/LANGUAGE.md)
- [Standard Library](docs/STDLIB.md)
- [Extensions (graphics/input/audio/3D)](docs/EXTENSIONS.md)
- [Candywrap / Sweet](docs/CANDYWRAP.md)
- [Candywrap Workflow Contract](docs/CANDYWRAP_WORKFLOW.md)
- [Distribution Guide](docs/DISTRIBUTION.md)
- [Architecture](ARCHITECTURE.md)

## Repository Layout

- `docs/` - user and team documentation
- `examples/` - runnable scripts and demos
- `compiler/` - internal toolchain source (private/internal development)
- `ARCHITECTURE.md` - architectural overview
- `ARCHITECTURE_HARDENING_PLAN.md` - implementation hardening roadmap

## Project Sharing Best Practice

If your project uses wrapped native libraries, commit the full generated module folder:

```text
candy_modules/<library>/
  <library>.candylib
  <library>_glue.c
  <library>.md
  ...
```

That keeps manifests, glue, and docs in sync across machines.
