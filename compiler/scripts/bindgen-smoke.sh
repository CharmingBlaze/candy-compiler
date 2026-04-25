#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/compiler"

go run ./cmd/candywrap wrap --name mylib --output ../examples/bindgen/out ../examples/bindgen/mylib.h
go run ./cmd/candy -ast ../examples/bindgen/main.candy >/dev/null

echo "bindgen smoke: OK"
