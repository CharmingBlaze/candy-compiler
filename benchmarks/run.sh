#!/usr/bin/env bash
set -euo pipefail

echo "=== Candy Benchmarks ==="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
COMPILER_DIR="${ROOT_DIR}/compiler"
SRC="${SCRIPT_DIR}/vec3_bench.candy"

cd "${COMPILER_DIR}"

echo "Debug build (-debug):"
go run ./cmd/candy -build -debug "${SRC}"
time ./../benchmarks/vec3_bench

echo ""
echo "Release build (-release):"
go run ./cmd/candy -build -release "${SRC}"
time ./../benchmarks/vec3_bench

echo ""
echo "Compare elapsed times above for rough speedup ratio."
