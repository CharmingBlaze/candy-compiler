#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 || $# -gt 3 ]]; then
  echo "usage: scripts/make-portable-release.sh <llvm-root> <out-dir> [raylib-runtime-dir]" >&2
  exit 1
fi

llvm_root="$1"
out_dir="$2"
raylib_runtime="${3:-}"

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
project_root="$(cd "$repo_root/.." && pwd)"
mkdir -p "$out_dir/bin"

cd "$repo_root"
go build -tags raylib -o "$out_dir/bin/candy" ./cmd/candy
go build -o "$out_dir/bin/candywrap" ./cmd/candywrap
go build -o "$out_dir/bin/sweet" ./cmd/sweet

bundle_dir="$out_dir/portable-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
"$repo_root/scripts/bundle-llvm.sh" "$out_dir/bin/candy" "$out_dir/bin/candywrap" "$llvm_root" "$bundle_dir" "$raylib_runtime"
cp "$out_dir/bin/sweet" "$bundle_dir/"

cp -R "$project_root/examples" "$bundle_dir/examples"
cp -R "$project_root/docs" "$bundle_dir/docs"

echo "Portable release ready at $bundle_dir"
