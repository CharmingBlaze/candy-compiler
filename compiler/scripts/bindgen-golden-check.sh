#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/compiler"

echo "Running candy_bindgen golden drift check..."
go test ./candy_bindgen -run TestGolden_
echo "bindgen golden check: OK"

