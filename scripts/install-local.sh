#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPILER_DIR="$ROOT/compiler"
OUT_DIR="$ROOT/release/bin"

mkdir -p "$OUT_DIR"
cd "$COMPILER_DIR"

echo "Building Candy tools..."
go build -tags raylib -o "$OUT_DIR/candy" ./cmd/candy
go build -o "$OUT_DIR/candywrap" ./cmd/candywrap
go build -o "$OUT_DIR/sweet" ./cmd/sweet

echo
echo "Install complete."
echo "Binaries:"
echo "  $OUT_DIR/candy"
echo "  $OUT_DIR/candywrap"
echo "  $OUT_DIR/sweet"
echo
echo "Optional: add to PATH"
echo "  export PATH=\"$OUT_DIR:\$PATH\""
