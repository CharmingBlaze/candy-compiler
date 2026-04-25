# Candy Performance Tracking

## Current Status

Updated: 2026-04-24

| Benchmark | Current | Target | Status |
|---|---:|---:|---|
| Vec3 Math | Not measured | <2 ns/op (typed hot loop) | 🟡 Harness pending |
| Physics 1K | 6.302 ms/step (`BenchmarkWorldStep_1000Bodies`) | <8 ms/frame | 🟢 On target (microbench) |
| Physics 1K (SpatialHash exp) | 0.301 ms/step (`BenchmarkWorldStep_1000Bodies_SpatialHash`) | <8 ms/frame | 🟢 Very fast, allocs improving |
| Compile 5K LOC | Not measured | <1s | 🔴 Not measured |
| Binary Size | Not measured | <10MB | 🔴 Not measured |

## Phase Progress

- [ ] Phase 1: Low-hanging fruit (2-5x)
- [ ] Phase 2: Core optimization (10x)
- [ ] Phase 3: C++ parity (15-35x)

## Recent Wins

- Added build profiles and LLVM pass-pipeline integration entrypoint.
- Added benchmark harnesses for LLVM codegen and physics stepping.
- Captured first benchmark smoke numbers:
  - `BenchmarkGenerateIR_TypedGameModule`: `24.7 us/op`, `29.2 KB/op`, `424 allocs/op`
  - `BenchmarkWorldStep_1000Bodies`: `6.302 ms/op`, `10.1 KB/op`, `4 allocs/op`
  - `BenchmarkWorldStep_1000Bodies_SpatialHash`: `0.301 ms/op`, `85.7 KB/op`, `297 allocs/op`

## Blockers

- Need benchmark baselines on representative scenes.
- CI gates currently use absolute thresholds; still need historical trend reporting.

## Benchmark Commands

- `go test ./compiler/candy_llvm -run ^$ -bench BenchmarkGenerateIR_TypedGameModule -benchmem`
- `go test ./compiler/candy_physics -run ^$ -bench BenchmarkWorldStep -benchmem`

## CI Perf Gate

- Workflow: `.github/workflows/perf-gate.yml`
- Script: `compiler/scripts/perf-gate.sh`
- Summary script: `compiler/scripts/perf-summary.sh`
- Enforces baseline thresholds for:
  - `BenchmarkGenerateIR_TypedGameModule`
  - `BenchmarkWorldStep_1000Bodies`
- Uploads raw benchmark logs as CI artifacts in `compiler/.perf-results/`
- Publishes a markdown benchmark table in the GitHub Actions job summary
- Upserts a sticky PR comment with the latest perf summary on pull requests
