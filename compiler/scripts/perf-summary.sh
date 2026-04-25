#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${OUT_DIR:-.perf-results}"
SUMMARY_FILE="${SUMMARY_FILE:-${OUT_DIR}/perf-summary.md}"

latest_summary="$(ls -1t "${OUT_DIR}"/summary-*.txt 2>/dev/null | head -n 1 || true)"
if [[ -z "${latest_summary}" ]]; then
  echo "No summary files found in ${OUT_DIR}" >&2
  exit 1
fi

read_kv() {
  local key="$1"
  awk -F= -v k="$key" '$1 == k { print $2; exit }' "${latest_summary}"
}

timestamp="$(read_kv timestamp)"
codegen_ns="$(read_kv codegen_ns)"
codegen_allocs="$(read_kv codegen_allocs)"
physics_ns="$(read_kv physics1000_ns)"
physics_allocs="$(read_kv physics1000_allocs)"
max_codegen_ns="$(read_kv max_codegen_ns)"
max_codegen_allocs="$(read_kv max_codegen_allocs)"
max_physics_ns="$(read_kv max_physics1000_ns)"
max_physics_allocs="$(read_kv max_physics1000_allocs)"

mkdir -p "${OUT_DIR}"

cat > "${SUMMARY_FILE}" <<EOF
## Perf Gate Summary

- **Timestamp (UTC):** ${timestamp}

| Benchmark | Measured | Threshold | Status |
|---|---:|---:|---|
| Codegen ns/op | ${codegen_ns} | ${max_codegen_ns} | $( (( codegen_ns <= max_codegen_ns )) && echo "PASS" || echo "FAIL" ) |
| Codegen allocs/op | ${codegen_allocs} | ${max_codegen_allocs} | $( (( codegen_allocs <= max_codegen_allocs )) && echo "PASS" || echo "FAIL" ) |
| Physics 1000 ns/op | ${physics_ns} | ${max_physics_ns} | $( (( physics_ns <= max_physics_ns )) && echo "PASS" || echo "FAIL" ) |
| Physics 1000 allocs/op | ${physics_allocs} | ${max_physics_allocs} | $( (( physics_allocs <= max_physics_allocs )) && echo "PASS" || echo "FAIL" ) |
EOF

echo "Wrote ${SUMMARY_FILE}"
