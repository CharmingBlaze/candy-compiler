# OPTIMIZATION_ROADMAP

## Phase 1: Low-Hanging Fruit (Week 1) - Target: 2-5x speedup

### Completed
- [x] Add canonical LLVM pass pipeline entrypoint (`compiler/candy_llvm/pipeline.go`)
- [x] Add build profiles to CLI (`-debug`, default dev-release, `-release`)
- [x] Wire `opt` resolution and optimization step before writing `.ll`

### High Priority (Do First)
- [x] Implement constant folding for literals
- [x] Remove boxing in typed math operations
- [x] Add function attributes (`hot`, `readnone`, etc.)
- [x] Create microbenchmark harness

### Medium Priority
- [ ] Dead code elimination in frontend lowering
- [x] String interning for color literals
- [ ] Compile-time resource folding

### Nice to Have
- [ ] Performance diagnostics ("add type here for speed")
- [ ] Perf CI gates

## Phase 1 Detailed Tasks

### Task 1: LLVM Pass Pipeline (in progress)
- **Files**: `compiler/candy_llvm/pipeline.go`, `compiler/candy_llvm/toolchain.go`, `compiler/cmd/candy/main.go`
- **Status**: implemented baseline
- **Notes**:
  - Uses `opt -S -passes=...` when available.
  - Gracefully falls back to unoptimized IR if `opt` is missing.
  - Build profiles:
    - `-debug` -> minimal optimization, clang `-O0 -g`
    - default -> dev-release, clang `-O2`
    - `-release` -> shipping, clang `-O3`

### Task 2: Function Attributes
- **Files**: `compiler/candy_llvm/attributes.go`, `compiler/candy_llvm/codegen.go`
- **Status**: baseline implemented
- **Notes**:
  - Emits `nounwind` + `willreturn` for generated functions.
  - Adds `hot inlinehint` for likely frame-loop functions (`main`, `update`, `render`, `tick`, `step`).
  - Adds `readnone nosync speculatable` for likely pure math helpers.

### Task 3: Remove Boxing in Typed Math
- **Planned files**: `compiler/candy_llvm/codegen_calls.go`, `compiler/candy_llvm/math_intrinsics.go`
- **Status**: baseline implemented
- **Notes**:
  - `sqrt` lowers to `@llvm.sqrt.f64` on typed numeric args.
  - `abs`, `min`, and `max` lower to typed compare+select IR without dynamic boxing.
  - Dynamic fallback path remains intact for non-typed/non-numeric cases.

### Task 4: Constant Folding
- **Planned files**: `compiler/candy_llvm/constant_fold.go`, `compiler/candy_raylib/utils.go`
- **Plan**:
  - Fold known literals (colors, arithmetic literals) at compile/lowering time.
  - Preserve runtime fallback for unknown/dynamic values.

### Task 5: Microbenchmark Harness
- **Files**: `benchmarks/vec3_bench.candy`, `benchmarks/run.sh`, `compiler/candy_llvm/bench_codegen_test.go`, `compiler/candy_physics/bench_world_test.go`
- **Status**: baseline implemented
- **Run**:
  - `go test ./candy_llvm -run ^$ -bench BenchmarkGenerateIR_TypedGameModule -benchmem`
  - `go test ./candy_physics -run ^$ -bench BenchmarkWorldStep -benchmem`

## Execution Sequence
1. Finalize and test Task 1 (pipeline correctness and CLI behavior).
2. Implement math intrinsic fast path + typed return tracking.
3. Add constant folding for literal hot paths.
4. Add benchmark automation and capture baseline numbers.

## Phase 2 Early Progress

- Implemented deterministic broadphase AABB pruning in `compiler/candy_physics/world.go`.
- Added scratch-buffer reuse for broadphase data to avoid per-step allocation spikes.
- Current microbench (`BenchmarkWorldStep_1000Bodies`) is within target budget.
- Added experimental deterministic spatial-hash mode (`BroadphaseSpatialHash`) with major speedup, pending allocation optimization before defaulting on.
- Reduced spatial-hash allocation overhead via reusable bucket/key/dedup scratch structures.
- Refactored spatial-hash build to mapless sorted cell entries for very low allocs/op.
- Switched to reusable per-cell bucket traversal for better speed/alloc balance in spatial hash mode.
