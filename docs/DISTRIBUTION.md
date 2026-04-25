# Distribution (Portable Bundles)

Goal: ship Candy so users do **not** need to install toolchains manually.

Portable bundles include:

- `bin/candy` (or `bin/candy.exe`)
- `bin/candywrap` (or `bin/candywrap.exe`)
- `sweet` (or `sweet.exe`)
- bundled `toolchain/` (clang + libs for native build/compile flow)
- compatibility `llvm/` alias/copy for older bundle consumers
- `licenses/`
- optional `raylib-runtime/` folder (if your platform build needs dynamic runtime files)
- `examples/` and `docs/` copied in for out-of-box use

## Windows

From `compiler/`:

```powershell
.\scripts\make-portable-release.ps1 -LlvmRoot "C:\toolchains\llvm" -OutDir ".\dist"
```

Recommended one-command release orchestration (build + bundle + archive):

```powershell
.\scripts\build-release.ps1 -LlvmRoot "C:\toolchains\llvm" -OutDir ".\dist" -Version "v1.2.3"
```

Optional raylib runtime folder:

```powershell
.\scripts\make-portable-release.ps1 -LlvmRoot "C:\toolchains\llvm" -OutDir ".\dist" -RaylibRuntimeDir "C:\raylib-runtime"
```

## Linux / macOS

From `compiler/`:

```bash
./scripts/make-portable-release.sh /opt/llvm ./dist
```

Recommended one-command release orchestration (build + bundle + archive):

```bash
./scripts/build-release.sh /opt/llvm ./dist v1.2.3
```

Optional raylib runtime folder:

```bash
./scripts/make-portable-release.sh /opt/llvm ./dist /opt/raylib-runtime
```

## Release checklist

- Build on each target OS (`windows`, `linux`, `darwin`) to avoid cross-platform runtime surprises.
- Run smoke checks inside the bundle:
  - `bin/candy --help`
  - `bin/candy doctor`
  - `bin/candywrap wrap --help`
  - `sweet convert --help`
  - one sample script from `examples/`
  - one `bin/candy build <file>.candy` flow (or `compile`) using bundled toolchain
- Zip/tar the final `portable-*` directory as your release artifact.
- Or use `build-release` scripts to emit ready-to-upload archive files directly.

## Notes

- `candywrap` can parse headers in `auto`, `libclang`, or `regex` mode.
- If shipping libclang-based parsing as default on user machines, include the matching libclang runtime in your bundle and test on clean machines.
- Toolchain discovery is self-contained:
  - env override (`CANDY_CLANG` / `CANDY_OPT`)
  - then bundled `toolchain/bin` (or compatibility `llvm/bin`)
  - then system `PATH`.
- `build-release` scripts inject build metadata into `candy doctor` output via linker flags:
  - `BuildVersion`
  - `BuildStdlibHash`
