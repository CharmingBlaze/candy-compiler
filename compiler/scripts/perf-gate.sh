#!/usr/bin/env bash
set -euo pipefail

# Simple absolute performance gates for CI.
# Tune thresholds as benchmarks stabilize.
MAX_CODEGEN_NS="${MAX_CODEGEN_NS:-40000}"
MAX_CODEGEN_ALLOCS="${MAX_CODEGEN_ALLOCS:-500}"
MAX_PHYSICS_1000_NS="${MAX_PHYSICS_1000_NS:-8000000}"
MAX_PHYSICS_1000_ALLOCS="${MAX_PHYSICS_1000_ALLOCS:-10}"
OUT_DIR="${OUT_DIR:-.perf-results}"

extract_metric() {
  local line="$1"
  local unit="$2"
  awk -v u="$unit" '{
    for (i = 1; i <= NF; i++) {
      if ($i == u && i > 1) {
        print $(i-1)
        exit
      }
    }
  }' <<<"$line"
}

mkdir -p "${OUT_DIR}"
run_ts="$(date -u +%Y%m%dT%H%M%SZ)"
codegen_log="${OUT_DIR}/codegen-${run_ts}.txt"
physics_log="${OUT_DIR}/physics-${run_ts}.txt"
summary_log="${OUT_DIR}/summary-${run_ts}.txt"

echo "Running codegen benchmark..."
codegen_out="$(go test ./candy_llvm -run '^$' -bench BenchmarkGenerateIR_TypedGameModule -benchmem -benchtime=100ms)"
echo "$codegen_out"
printf '%s\n' "$codegen_out" > "$codegen_log"
codegen_line="$(awk '/^BenchmarkGenerateIR_TypedGameModule/ { print; exit }' <<<"$codegen_out")"
if [[ -z "$codegen_line" ]]; then
  echo "ERROR: missing BenchmarkGenerateIR_TypedGameModule output"
  exit 1
fi
codegen_ns="$(extract_metric "$codegen_line" "ns/op")"
codegen_allocs="$(extract_metric "$codegen_line" "allocs/op")"

echo "Running physics benchmark..."
physics_out="$(go test ./candy_physics -run '^$' -bench BenchmarkWorldStep_1000Bodies -benchmem -benchtime=100ms)"
echo "$physics_out"
printf '%s\n' "$physics_out" > "$physics_log"
physics_line="$(awk '/^BenchmarkWorldStep_1000Bodies/ { print; exit }' <<<"$physics_out")"
if [[ -z "$physics_line" ]]; then
  echo "ERROR: missing BenchmarkWorldStep_1000Bodies output"
  exit 1
fi
physics_ns="$(extract_metric "$physics_line" "ns/op")"
physics_allocs="$(extract_metric "$physics_line" "allocs/op")"

echo "=== Perf Gate Summary ==="
echo "Codegen: ${codegen_ns} ns/op, ${codegen_allocs} allocs/op (max ${MAX_CODEGEN_NS}, ${MAX_CODEGEN_ALLOCS})"
echo "Physics1000: ${physics_ns} ns/op, ${physics_allocs} allocs/op (max ${MAX_PHYSICS_1000_NS}, ${MAX_PHYSICS_1000_ALLOCS})"
{
  echo "timestamp=${run_ts}"
  echo "codegen_ns=${codegen_ns}"
  echo "codegen_allocs=${codegen_allocs}"
  echo "physics1000_ns=${physics_ns}"
  echo "physics1000_allocs=${physics_allocs}"
  echo "max_codegen_ns=${MAX_CODEGEN_NS}"
  echo "max_codegen_allocs=${MAX_CODEGEN_ALLOCS}"
  echo "max_physics1000_ns=${MAX_PHYSICS_1000_NS}"
  echo "max_physics1000_allocs=${MAX_PHYSICS_1000_ALLOCS}"
} > "$summary_log"
echo "Saved raw benchmark logs to ${OUT_DIR}"

fail=0
if (( codegen_ns > MAX_CODEGEN_NS )); then
  echo "FAIL: codegen ns/op regression threshold exceeded"
  fail=1
fi
if (( codegen_allocs > MAX_CODEGEN_ALLOCS )); then
  echo "FAIL: codegen allocs/op threshold exceeded"
  fail=1
fi
if (( physics_ns > MAX_PHYSICS_1000_NS )); then
  echo "FAIL: physics 1000 ns/op threshold exceeded"
  fail=1
fi
if (( physics_allocs > MAX_PHYSICS_1000_ALLOCS )); then
  echo "FAIL: physics 1000 allocs/op threshold exceeded"
  fail=1
fi

if (( fail != 0 )); then
  exit 1
fi

echo "Perf gates passed."
