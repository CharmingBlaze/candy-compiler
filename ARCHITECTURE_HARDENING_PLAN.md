# Architecture Hardening Plan

This plan defines the next implementation sprint(s) to harden Candy's architecture where it matters most for developer experience and long-term maintainability.

Scope priority:

1. `candy_bindgen` stability and determinism
2. Toolchain probe hardening (`clang`/`opt`)
3. LLVM pipeline profile validation with evidence
4. LSP/CLI parity through shared import + manifest resolution
5. Typed-AST bridge (minimal, value-first)

## Goals

- Keep FFI contract deterministic and regression-resistant.
- Fail fast with actionable native-build diagnostics.
- Ensure optimization profiles are benchmarked, not assumed.
- Ensure editor diagnostics and CLI behavior agree.
- Introduce typed-AST only where it improves optimization/codegen decisions.

---

## Sprint Milestones

## Sprint 1 (Week 1): FFI Contract Hardening + Toolchain Preflight

### Milestone 1.1 - `candy_bindgen` golden stability suite

Deliverables:

- Deterministic golden tests for:
  - Name transforms (`--namespace`, `--strip-prefix`, `--ignore`)
  - Collision resolution behavior
  - Reserved Candy identifier handling
  - Variadic/function-pointer ABI guardrail filtering
  - Manifest serialization order and stable formatting
- Fixture corpus for C and C++ flavored headers (safe + unsafe ABI shapes)

Planned files:

- `compiler/candy_bindgen/*_golden_test.go` (new)
- `compiler/candy_bindgen/testdata/golden/**` (new fixtures + expected outputs)

### Milestone 1.2 - Toolchain preflight and fail-fast diagnostics

Deliverables:

- Add a preflight API in `candy_llvm` invoked before native codegen:
  - Detect `clang` and `opt` paths
  - Detect versions
  - Validate minimum compatibility matrix
  - Emit clear remediation guidance
- Integrate into `candy build` / `candy compile` entry path.
- Add a standalone `candy doctor` command that runs probe checks without compiling.

Planned files:

- `compiler/candy_llvm/toolchain_probe.go` (new)
- `compiler/candy_llvm/toolchain.go` (update)
- `compiler/cmd/candy/main.go` (update)

---

## Sprint 2 (Week 2): LLVM Profile Evidence + LSP/CLI Parity

### Milestone 2.1 - LLVM profile evidence matrix

Deliverables:

- Benchmark harness runs each profile (`debug`, `dev-release`, `shipping`) against:
  - Tiny script workload
  - Mid-sized app/module workload
  - FFI-heavy workload (`.candylib` import + glue link)
- Capture metrics:
  - Compile time
  - Runtime throughput/latency
  - Binary size
  - IR size (optional secondary metric)
- Publish report table and baseline thresholds.

Planned files:

- `PERFORMANCE.md` (update with profile breakdown)
- `compiler/scripts/perf-gate.sh` (extend thresholds/modes)
- `compiler/scripts/perf-summary.sh` (extend reporting)

### Milestone 2.2 - LSP/CLI parity path

Deliverables:

- LSP diagnostics path reuses import + `.candylib` resolution logic equivalent to `candy_load`.
- Add parity tests proving same source yields equivalent diagnostics under:
  - CLI (`candy -check` / parse+check path)
  - LSP (`didOpen` / `didChange`)

Planned files:

- `compiler/candy_lsp/*` (update for shared load path)
- `compiler/candy_load/*` (minor extraction if needed for reuse)
- `compiler/candy_lsp/server_test.go` (new parity cases)

---

## Sprint 3 (Week 3): Typed-AST Bridge (Minimal Vertical Slice)

### Milestone 3.1 - Typed metadata on selected AST nodes

Deliverables:

- Start with `CallExpression`-first typed metadata:
  - inferred callee signature shape
  - argument type vector (where confidence is high)
  - return type hint
- Add optional type slots/annotations to the minimal supporting expression set:
  - literals
  - identifiers
  - infix expressions
- Populate these annotations from `candy_typecheck` where confidence is high.
- Make optimizer/codegen consume annotations in one concrete optimization path:
  - static-call lowering for known builtins or known function targets to reduce dynamic dispatch overhead.

Planned files:

- `compiler/candy_ast/*` (minimal node metadata additions)
- `compiler/candy_typecheck/*` (annotation pass output)
- `compiler/candy_opt/*` or `compiler/candy_llvm/*` (first consumer path)

Non-goal:

- Full strict type-system redesign or complete monomorphization framework in this sprint.

---

## Test Matrix

## A) `candy_bindgen` Golden Test Suite

Each test case follows:

- Input header(s)
- CLI/config flags
- Expected generated artifacts:
  - `<lib>.candylib`
  - `<lib>_glue.c`
  - optional `<lib>.candy` / `<lib>_namespace.candy`
  - optional docs markdown
- Expected diagnostics/warnings

Case matrix:

1. **Name transforms**
   - Input: prefixed API (`b2World_*`, `b2Body_*`)
   - Flags: `--strip-prefix`, `--namespace`
   - Assert: stable transformed names + deterministic collision fallback.

2. **Reserved keywords**
   - Input: symbols colliding with language keywords (`end`, etc.)
   - Assert: generated names are escaped/renamed consistently and compile.

3. **ABI guardrails**
   - Input: variadics/function pointers/mixed complex signatures
   - Assert:
     - rejected in safe mode
     - included with metadata in `--unsafe-abi` mode
     - warning text deterministic.

4. **C++ shim path**
   - Input: C++ headers, `--lang c++ --cxx-shim`
   - Assert:
     - shim template generated
     - manifest includes shim/glue as expected
     - link metadata includes required stdlib where applicable.

5. **Determinism**
   - Same input run twice
   - Assert byte-for-byte identical outputs (except controlled timestamp fields, if present).

6. **Wildcard/root expansion**
   - Input roots + globs
   - Assert stable discovery ordering and identical outputs across runs.

7. **Golden update workflow**
   - Support an explicit update mode for intentional output changes.
   - Suggested test flag/workflow: `-goldens:update` (or equivalent env var gate).
   - Assert update mode is opt-in and never runs silently in CI.

## B) LSP/CLI Parity Golden Suite

Each case follows:

- Source tree with imports (including `.candylib` where relevant)
- Expected diagnostics list (message + severity + line/col normalization)
- Assert same normalized diagnostics from:
  - CLI check path
  - LSP diagnostic path

Case matrix:

1. Parse error in imported file.
2. Typecheck warning in imported symbol usage.
3. `.candylib` ABI rejection error.
4. Missing import path resolution error.
5. Reserved-name or extern mismatch surfaced consistently.

Normalization rules:

- Compare diagnostic tuples:
  - `kind` (parse/type/load/build)
  - normalized message code/text
  - source file
  - line/column
- Ignore presentation-only differences (caret rendering, ANSI formatting).

---

## Acceptance Criteria (Definition of Done)

## 1) `candy_bindgen` stability

- Golden suite exists and passes in CI.
- No nondeterministic output diffs on repeated generation.
- Collision/reserved-name behavior covered by explicit tests.
- ABI guardrail behavior is validated in both safe and unsafe modes.

## 2) Toolchain preflight

- `candy build` fails before codegen when required tools are unusable.
- Error includes:
  - what was searched
  - what was found (version/path)
  - exact fix suggestions per OS
- Unit tests cover "not found", "version mismatch", and "valid toolchain".

## 3) LLVM profile validation

- Performance report includes all three profiles and three workload classes.
- Thresholds are codified in perf gate.
- At least one documented profile recommendation based on measured data.

## 4) LSP/CLI parity

- Parity suite passes.
- Any known divergence is documented with explicit rationale and issue link.
- Import + manifest resolution path is shared or behaviorally equivalent by test.

## 5) Typed-AST bridge

- Minimal node annotations are produced and consumed in one optimization/codegen path.
- Measurable benefit exists (reduced boxing, fewer dynamic fallbacks, or compile/runtime win).
- No broad type-system rewrite required to ship this milestone.

---

## Toolchain Probe Logic (Cross-Platform Technical Note)

Probe precedence:

1. Explicit env overrides (for reproducibility):
   - `CANDY_CLANG`
   - `CANDY_OPT`
2. Bundled toolchain in repo/release bundle:
   - `./llvm/bin` (platform-specific executable suffix)
3. System `PATH`
4. (Optional) common install locations by OS as fallback hints.

Platform details:

- **Windows**
  - Resolve `.exe` automatically.
  - Check common candidates for user guidance only (LLVM default install paths, MSYS/Chocolatey locations) if PATH lookup fails.
- **macOS**
  - Check PATH; include Xcode/CommandLineTools guidance in remediation.
- **Linux**
  - Check PATH; include package-manager guidance (`apt`, `dnf`, `pacman`) in remediation text.

Version policy:

- Parse `clang --version` and `opt --version`.
- Require same major version family for "strict compatible" mode, or allow a documented tolerance window.
- If mismatch:
  - hard fail in strict mode
  - warning + continue in permissive/dev mode (if explicitly enabled).

Preflight API behavior:

- Returns structured result:
  - found paths
  - parsed versions
  - compatibility verdict
  - actionable messages
- CLI converts that into concise human-readable output.

---

## Risks and Mitigations

- **Risk:** Golden fixtures become brittle.
  - **Mitigation:** normalize expected outputs and isolate intentionally variable fields.
- **Risk:** LSP parity introduces latency.
  - **Mitigation:** cache loaded modules and diagnostics hashes; keep incremental path.
- **Risk:** Typed-AST scope creep.
  - **Mitigation:** lock to one vertical slice and one consumer path until value is proven.

---

## Execution Checklist

- [x] Create bindgen golden fixture corpus and expected artifacts.
- [x] Implement deterministic output guards and tests.
- [x] Add toolchain preflight API + CLI integration.
- [ ] Add profile benchmark matrix and perf report updates.
- [ ] Implement LSP/CLI parity harness with shared load behavior.
- [ ] Implement minimal typed-AST metadata + single consumer optimization.
- [~] Publish sprint summary and update roadmap docs.

---

## Live Status (Sprint 1)

Updated: 2026-04-25

Completed:

- `ARCH-101` Build golden fixture corpus for bindgen.
- `ARCH-102` Add golden runner and deterministic output assertions.
- `ARCH-103` Validate transform collisions and reserved keyword handling.
- `ARCH-104` ABI guardrail behavior tests (safe vs unsafe).
- `ARCH-105` Add bindgen golden update mode (+ CI guard against update mode in CI).
- `ARCH-201` Implement structured toolchain probe API.
- `ARCH-202` Integrate preflight into CLI native build path (fail-fast when clang missing).
- `ARCH-204` Add `candy doctor` command for environment health.
- `ARCH-203` Cross-platform resolution hardening (OS-specific toolchain guidance + resolution precedence tests).
- CI enforcement: cross-platform workflow step verifies bindgen golden outputs.

Not started (by design for later sprints):

- `ARCH-301`, `ARCH-302`, `ARCH-303`
- `ARCH-401`, `ARCH-402`, `ARCH-403`
- `ARCH-501`, `ARCH-502`, `ARCH-503`

---

## Ticket-Ready Breakdown

Use this section as direct issue/board input.

Conventions:

- IDs are suggested issue keys (`ARCH-*`).
- Estimates are ideal engineering days.
- "Depends on" defines hard ordering.

## Epic A - Bindgen Determinism and Contract Stability

### ARCH-101: Build golden fixture corpus for bindgen

- **Scope**
  - Create fixture headers for C/C++ edge cases.
  - Add expected manifest/glue/stub/docs outputs in `testdata/golden`.
- **Estimate**: 1.5 days
- **Depends on**: none
- **Done when**
  - Fixture set includes transforms, reserved names, ABI edge cases, wildcard/root discovery.

### ARCH-102: Add golden runner and deterministic output assertions

- **Scope**
  - Implement golden test harness (`*_golden_test.go`).
  - Normalize any volatile fields (timestamps, path separators) before comparison.
- **Estimate**: 1.5 days
- **Depends on**: `ARCH-101`
- **Done when**
  - Repeated generation yields no diff in CI.

### ARCH-103: Validate transform collisions and reserved keyword handling

- **Scope**
  - Explicit tests for namespace + strip-prefix collisions.
  - Reserved identifier rewrite behavior verified and documented.
- **Estimate**: 1 day
- **Depends on**: `ARCH-101`, `ARCH-102`
- **Done when**
  - Collision and reserved-name scenarios have deterministic expected outputs.

### ARCH-104: ABI guardrail behavior tests (safe vs unsafe)

- **Scope**
  - Verify safe mode rejection + `--unsafe-abi` inclusion behavior.
  - Verify warning text stability.
- **Estimate**: 1 day
- **Depends on**: `ARCH-102`
- **Done when**
  - Guardrail scenarios pass in both modes with expected diagnostics.

## Epic B - Toolchain Preflight and Fail-Fast UX

### ARCH-201: Implement structured toolchain probe API

- **Scope**
  - Add probe for `clang` and `opt` paths, versions, and compatibility verdict.
  - Return structured result object.
- **Estimate**: 1.5 days
- **Depends on**: none
- **Done when**
  - Unit tests cover not-found, mismatch, and valid cases.

### ARCH-202: Integrate preflight into CLI native build path

- **Scope**
  - Run probe before IR generation in `candy build` / `compile`.
  - Render concise remediation output by platform.
- **Estimate**: 1 day
- **Depends on**: `ARCH-201`
- **Done when**
  - Native build exits early with actionable message when probe fails.

### ARCH-204: Add `candy doctor` command for environment health

- **Scope**
  - Add CLI command to run preflight checks independently.
  - Print structured pass/fail report for `clang`, `opt`, version compatibility, and search locations.
- **Estimate**: 0.75 day
- **Depends on**: `ARCH-201`
- **Done when**
  - `candy doctor` runs without source input and returns non-zero on hard failures.

### ARCH-203: Cross-platform resolution hardening

- **Scope**
  - Validate path precedence: env override -> bundled -> PATH.
  - Add OS-specific guidance for Windows/macOS/Linux.
- **Estimate**: 1 day
- **Depends on**: `ARCH-201`
- **Done when**
  - Tests validate expected resolution precedence and message content.

## Epic C - LLVM Profile Evidence and Perf Gates

### ARCH-301: Add profile benchmark matrix runner

- **Scope**
  - Script/harness to run tiny/mid/FFI workloads across all profiles.
  - Emit structured artifacts for CI.
- **Estimate**: 1.5 days
- **Depends on**: none
- **Done when**
  - Matrix output available for CI and local run.

### ARCH-105: Add bindgen golden update mode

- **Scope**
  - Add explicit golden-update mechanism for intentional output changes.
  - Ensure CI forbids update mode and validates committed golden outputs.
- **Estimate**: 0.75 day
- **Depends on**: `ARCH-102`
- **Done when**
  - Developers can refresh goldens intentionally with one command/flag.
  - CI fails if golden outputs drift.

### ARCH-302: Extend perf summary and gate thresholds

- **Scope**
  - Update perf scripts and report tables with profile dimensions.
  - Set baseline thresholds.
- **Estimate**: 1 day
- **Depends on**: `ARCH-301`
- **Done when**
  - CI fails on regression against defined profile baselines.

### ARCH-303: Publish profile recommendation doc update

- **Scope**
  - Update `PERFORMANCE.md` with measured trade-offs and recommended default profile guidance.
- **Estimate**: 0.5 day
- **Depends on**: `ARCH-301`, `ARCH-302`
- **Done when**
  - Recommendations are backed by current benchmark data.

## Epic D - LSP/CLI Diagnostic Parity

### ARCH-401: Extract/reuse shared import + manifest resolution path

- **Scope**
  - Reuse `candy_load` behavior from LSP diagnostics path or implement equivalent shared helper.
- **Estimate**: 1.5 days
- **Depends on**: none
- **Done when**
  - LSP diagnostics resolve imports/manifests with same rules as CLI.

### ARCH-402: Add parity golden diagnostics suite

- **Scope**
  - Build source-tree fixtures and expected normalized diagnostic tuples.
  - Compare LSP and CLI outputs.
- **Estimate**: 1.5 days
- **Depends on**: `ARCH-401`
- **Done when**
  - Parity tests pass, or any differences are explicitly documented.

### ARCH-403: Cache and performance safeguards for LSP parity path

- **Scope**
  - Ensure incremental LSP performance remains acceptable.
  - Add caching/hash checks where needed.
- **Estimate**: 1 day
- **Depends on**: `ARCH-401`
- **Done when**
  - No meaningful editor latency regression from parity improvements.

## Epic E - Typed-AST Bridge (Minimal Vertical Slice)

### ARCH-501: Define typed-annotation schema for hot expression nodes

- **Scope**
  - Add optional type metadata for selected node classes only, with `CallExpression` as the first-class target.
  - Document annotation confidence levels (`known`, `inferred`, `dynamic`).
- **Estimate**: 1 day
- **Depends on**: none
- **Done when**
  - Schema merged with zero parser behavior impact.

### ARCH-502: Populate annotations in typecheck pass

- **Scope**
  - Emit annotations for literals/idents/infix/call where confidence is high.
- **Estimate**: 1.5 days
- **Depends on**: `ARCH-501`
- **Done when**
  - Annotation coverage tests exist for selected node set.

### ARCH-503: Consume typed annotations in one optimization/codegen path

- **Scope**
  - Use annotations in one measurable path:
    - static-call lowering for known function targets from `CallExpression` annotations.
  - Track measurable impact (dispatch overhead reduction and/or runtime win on targeted microbench).
- **Estimate**: 1.5 days
- **Depends on**: `ARCH-502`
- **Done when**
  - Measurable improvement is documented and regression-tested.

---

## Suggested Sprint Packing

### Week 1 (Execution-Focused)

- `ARCH-101`, `ARCH-102`, `ARCH-103`, `ARCH-104`, `ARCH-105`
- `ARCH-201`, `ARCH-202`, `ARCH-204`

Stretch:

- `ARCH-203`

### Week 2 (Validation + Parity)

- `ARCH-301`, `ARCH-302`, `ARCH-303`
- `ARCH-401`, `ARCH-402`

Stretch:

- `ARCH-403`

### Week 3 (Typed-AST Vertical Slice)

- `ARCH-501`, `ARCH-502`, `ARCH-503`

---

## Status Template (for Daily/PR Updates)

Use this format in standups and PR descriptions:

- **Ticket**: `ARCH-XXX`
- **State**: `todo | in_progress | blocked | done`
- **What changed**
- **Evidence** (tests, benchmark output, fixture diff)
- **Risk/Blocker**
- **Next step**
