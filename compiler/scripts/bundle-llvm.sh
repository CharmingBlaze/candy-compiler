#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 4 || $# -gt 5 ]]; then
  echo "usage: scripts/bundle-llvm.sh <candy-binary> <candywrap-binary> <llvm-root> <out-dir> [raylib-runtime-dir]" >&2
  exit 1
fi

candy_bin="$1"
candywrap_bin="$2"
llvm_root="$3"
out_dir="$4"
raylib_runtime="${5:-}"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ ! -f "$candy_bin" ]]; then
  echo "candy binary not found: $candy_bin" >&2
  exit 1
fi
if [[ ! -f "$candywrap_bin" ]]; then
  echo "candywrap binary not found: $candywrap_bin" >&2
  exit 1
fi
if [[ ! -d "$llvm_root" ]]; then
  echo "llvm root not found: $llvm_root" >&2
  exit 1
fi

manifest="${CANDY_LLVM_MANIFEST:-}"
if [[ -z "$manifest" ]]; then
  case "$(uname -s)" in
    Linux*) manifest="$repo_root/scripts/llvm-manifest.linux-x64.txt" ;;
    Darwin*) manifest="$repo_root/scripts/llvm-manifest.macos-universal.txt" ;;
    *) manifest="" ;;
  esac
fi

if [[ -n "$manifest" && -f "$manifest" ]]; then
  while IFS= read -r pat; do
    [[ -z "$pat" ]] && continue
    if ! compgen -G "$llvm_root/$pat" > /dev/null; then
      echo "missing required LLVM artifact for pattern: $pat" >&2
      exit 1
    fi
  done < "$manifest"
fi

mkdir -p "$out_dir"
mkdir -p "$out_dir/bin"
cp "$candy_bin" "$out_dir/bin/"
cp "$candywrap_bin" "$out_dir/bin/"
mkdir -p "$out_dir/toolchain"
cp -R "$llvm_root/bin" "$out_dir/toolchain/bin"
if [[ -d "$llvm_root/lib" ]]; then
  cp -R "$llvm_root/lib" "$out_dir/toolchain/lib"
fi

# Backward-compatible alias for older bundles/scripts expecting ./llvm.
if [[ ! -e "$out_dir/llvm" ]]; then
  ln -s "toolchain" "$out_dir/llvm" 2>/dev/null || cp -R "$out_dir/toolchain" "$out_dir/llvm"
fi

if [[ -f "$repo_root/licenses/LLVM-LICENSE.txt" ]]; then
  mkdir -p "$out_dir/licenses"
  cp "$repo_root/licenses/LLVM-LICENSE.txt" "$out_dir/licenses/LLVM-LICENSE.txt"
else
  echo "missing required license file: licenses/LLVM-LICENSE.txt" >&2
  exit 1
fi

if [[ -n "$raylib_runtime" && -d "$raylib_runtime" ]]; then
  mkdir -p "$out_dir/raylib-runtime"
  cp -R "$raylib_runtime"/. "$out_dir/raylib-runtime/"
fi

cat > "$out_dir/README_PORTABLE.txt" <<'EOF'
Candy Portable Bundle
=====================

This folder is self-contained. No global installs are required.

Included:
- bin/candy (or bin/candy.exe)
- bin/candywrap (or bin/candywrap.exe)
- sweet (or sweet.exe)
- toolchain/ (clang + toolchain used by candy -build)
- llvm/ (compatibility alias for older bundles)
- licenses/
- optional raylib-runtime/ (if provided at packaging time)

Usage:
- ./bin/candy script.candy
- ./bin/candywrap wrap --name mylib --output ./bindings mylib.h
- ./sweet convert --name mylib --output ./bindings mylib.h
EOF

echo "Bundled package created at $out_dir"
