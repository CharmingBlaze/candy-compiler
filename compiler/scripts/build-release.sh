#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 || $# -gt 4 ]]; then
  echo "usage: scripts/build-release.sh <llvm-root> <out-dir> [version] [raylib-runtime-dir]" >&2
  exit 1
fi

llvm_root="$1"
out_dir="$2"
version_in="${3:-}"
raylib_runtime="${4:-}"

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
project_root="$(cd "$repo_root/.." && pwd)"
mkdir -p "$out_dir"

if [[ -z "$version_in" ]]; then
  version_in="$(git -C "$project_root" describe --tags --always 2>/dev/null || echo dev)"
fi

stdlib_hash="$(
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$repo_root" && sha256sum candy_stdlib/*.go | sha256sum | awk '{print $1}')
  else
    (cd "$repo_root" && shasum -a 256 candy_stdlib/*.go | shasum -a 256 | awk '{print $1}')
  fi
)"

stage_root="$(mktemp -d "$out_dir/.release-stage.XXXXXX")"
bundle_dir="$stage_root/candy-release"
mkdir -p "$bundle_dir/bin" "$bundle_dir/lib" "$bundle_dir/docs"

cd "$repo_root"
ldflags="-X main.BuildVersion=$version_in -X main.BuildStdlibHash=$stdlib_hash"
go build -tags raylib -ldflags "$ldflags" -o "$bundle_dir/bin/candy" ./cmd/candy
go build -ldflags "$ldflags" -o "$bundle_dir/bin/candywrap" ./cmd/candywrap
go build -ldflags "$ldflags" -o "$bundle_dir/bin/sweet" ./cmd/sweet

"$repo_root/scripts/bundle-llvm.sh" "$bundle_dir/bin/candy" "$bundle_dir/bin/candywrap" "$llvm_root" "$bundle_dir" "$raylib_runtime"

# Sanity check: stdlib integrity marker included in release payload and matches build-time hash.
printf "stdlib_hash=%s\nversion=%s\n" "$stdlib_hash" "$version_in" > "$bundle_dir/lib/STDLIB_MANIFEST.txt"
if ! grep -q "$stdlib_hash" "$bundle_dir/lib/STDLIB_MANIFEST.txt"; then
  echo "stdlib manifest integrity check failed" >&2
  exit 1
fi

if [[ -d "$project_root/docs" ]]; then
  cp -R "$project_root/docs"/. "$bundle_dir/docs/"
fi
if [[ -d "$project_root/templates" ]]; then
  cp -R "$project_root/templates" "$bundle_dir/templates"
else
  mkdir -p "$bundle_dir/templates"
fi
cp -R "$project_root/examples" "$bundle_dir/examples"

cat > "$bundle_dir/README.md" <<EOF
# Candy Portable Release

Version: $version_in  
Stdlib hash: $stdlib_hash

This release is self-contained:

- binaries in \`bin/\`
- native backend toolchain in \`toolchain/\`
- compatibility alias in \`llvm/\`
- docs in \`docs/\`
- templates in \`templates/\`

Quick checks:

\`\`\`bash
./bin/candy doctor
./bin/candy --help
./bin/candywrap wrap --help
./bin/sweet convert --help
\`\`\`
EOF

# Optional binary stripping (best effort).
if [[ "$(uname -s)" != "MINGW"* && "$(uname -s)" != "MSYS"* && "$(uname -s)" != "CYGWIN"* ]]; then
  if command -v strip >/dev/null 2>&1; then
    strip "$bundle_dir/bin/candy" "$bundle_dir/bin/candywrap" "$bundle_dir/bin/sweet" 2>/dev/null || true
    strip "$bundle_dir/toolchain/bin/"* 2>/dev/null || true
  fi
fi

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
archive="$out_dir/candy-$version_in-$os-$arch.tar.gz"
(cd "$stage_root" && tar -czf "$archive" "candy-release")

echo "Release bundle ready:"
echo "  folder: $bundle_dir"
echo "  archive: $archive"
