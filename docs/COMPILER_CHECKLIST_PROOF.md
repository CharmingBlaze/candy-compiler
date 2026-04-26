# Compiler Checklist Proof Map

Source checklist: `docs/compiler checklist.md`

This file maps checklist sections to concrete implementation and test evidence.

## 1-14 Core + Interop Syntax

- `1 Variables & types` -> parser/type inference in `compiler/candy_parser/statements_parser.go`, type checking in `compiler/candy_typecheck`.
- `2 Operators` -> token/parser support in `compiler/candy_token/token.go` + `compiler/candy_parser/parser.go`, runtime in `compiler/candy_evaluator/eval.go`, LLVM lowering in `compiler/candy_llvm/codegen*.go`.
- `3 Control flow` -> parse/eval in `compiler/candy_parser/statements_parser.go` + `compiler/candy_evaluator/eval_for.go` + `compiler/candy_evaluator/eval.go`.
- `4 Functions + lambdas` -> parse in `compiler/candy_parser/statements_parser.go` + `compiler/candy_parser/expressions_parser.go`; runtime calls in `compiler/candy_evaluator/eval.go`.
- `5 Objects/structs` -> parse in `compiler/candy_parser/parse_decls.go` and `compiler/candy_parser/expressions_parser.go`; runtime instantiation/methods in `compiler/candy_evaluator/eval.go` + `compiler/candy_evaluator/eval_struct.go`.
- `6 Arrays/lists + array()/bytes()` -> parser/runtime in `compiler/candy_parser/expressions_parser.go` + `compiler/candy_evaluator/builtin.go` + `compiler/candy_evaluator/eval_container_ops.go`.
- `7 Strings` -> interpolation parse in `compiler/candy_parser/expressions_parser.go`; runtime string ops in `compiler/candy_evaluator/eval.go` + `compiler/candy_evaluator/stdlib_builtins.go`.
- Game-dev ergonomics (`docs/new stuff 1.md`) -> vector runtime and helpers in `compiler/candy_evaluator/vec_runtime.go`; physics helper behavior in `compiler/candy_evaluator/physics_runtime_helpers.go`; import ergonomics and named-call parsing in `compiler/candy_parser/parse_simplec_basic.go` + `compiler/candy_parser/expressions_parser.go`.
- Game systems 1.2 (`docs/new stuff 1.2.md`) -> helper runtime surface in `compiler/candy_evaluator/vec_runtime.go` + `compiler/candy_evaluator/eval_container_ops.go`; parser anchors in `compiler/candy_parser/parser_test.go`; helper type registration in `compiler/candy_typecheck/check.go`.
- `8 Null/default` -> `null`, `??`, and `or` fallback behavior in `compiler/candy_evaluator/eval.go`.
- `9 Try/catch/finally` -> parse/eval in `compiler/candy_parser/parse_simplec_basic.go` + `compiler/candy_evaluator/eval_try.go`.
- `10 with/delete` -> parser in `compiler/candy_parser/statements_parser.go`; runtime in `compiler/candy_evaluator/eval.go` + `compiler/candy_evaluator/eval_struct.go`.
- `11 Comments` -> lexer skip rules in `compiler/candy_lexer/lexer.go`.
- `12 Enums` -> parse/eval in `compiler/candy_parser/parse_simplec_basic.go` + `compiler/candy_evaluator/eval.go`.
- `13 Extern declarations` -> parser support for `extern fun ...`, `extern name(...)`, and variadic syntax in `compiler/candy_parser/parse_decls.go`; ABI guardrails in `compiler/candy_bindgen/manifest.go`.
- `14 Library system` -> parser support in `compiler/candy_parser/statements_parser.go`, loader integration in `compiler/candy_load/load.go`.

### 1-14 Direct Runtime/Parser Tests

- Variables/types/functions: `TestValStatements`, `TestVarStatements`, `TestParse_EquivalentVariablesAndFunctions`, `TestTypeAnnotatedVariables`.
- Core syntax + declarations: `TestSimpleCSyntax`, `TestKotlinStyleDeclarationsParse`, `TestKotlinStyleDeclarationsNoSemicolons`, `TestFunctionFunAndFuncAreAliases`.
- Nullability/receiver/default/null-coalesce: `TestNullableTypeSuffix`, `TestNullableTypeSuffix_MixedCaseInputCanonicalized`, `TestReceiverSyntaxFunction`, `TestEval_NullishCoalesce`, `TestEval_DefaultValueOperatorOrWithNull`, `TestEval_DefaultValueOperatorOrDoesNotOverrideZero`.
- Control flow/loops: `TestOptionalSemicolons`, `TestForInIterableStopsBeforeBlockBrace`, `TestEval_WhileAccumulate`, `TestEval_ForInDoesNotUpdateOuterByAssignment`.
- Objects/structs/enums/with/bitwise: `TestInlineObjectLiteralParsing`, `TestFullSpecStructsAndProperties`, `TestParseWithStatementAndBitwise`, `TestEval_ObjectDeclarationInstantiationAndMethodCall`, `TestEval_BitwiseOps`, `TestFullSpecEnums`, `TestEval_EnumValuesAndAccess`.
- Error handling: `TestFullSpecTryCatch`, `TestEval_TryCatch`, `TestEval_TryFinally`.
- Extern/library syntax: `TestParseLibrarySyntax`, `TestParseExternWithoutFunKeyword`, `TestParseExternVariadicSignature`.
- Switch/case form parity: `TestSwitchCaseColonStyleParsing`, `TestEval_SwitchCaseColonStyle`.
- Tuple destructure + ignored `_`: `TestEval_MultipleReturnDestructureAndIgnoreUnderscore`.
- Exclusive ranges/slicing: `TestExclusiveRangeToken`, `TestExclusiveRangeAndSliceParsing`, `TestEval_ExclusiveRangeAndArraySlice`.
- Extended method ergonomics: `TestEval_StringAndArrayExtendedMethods` (string `indexOf`/`substring`, array `reduce/find/all/any/unique`).
- Membership ergonomics: `TestInOperatorParsing`, `TestEval_InOperatorForArrayMapAndString` (`in` for array/map/string membership).
- Optional syntax coverage: `TestNotInOperatorParsing`, `TestEval_NotInOperator`, `TestTernaryOperatorParsing`, `TestEval_TernaryOperator`.
- Game-dev additions: `TestImportAliasAndFromImportParsing`, `TestNamedCallArgumentParsing`, `TestEval_VecBuiltinsAndOps`, `TestEval_NamedArgumentsOnUserFunction`.

## 15-19 Builtins + File IO

- `15 I/O` -> `print/println/input/readLine` in `compiler/candy_evaluator/builtin.go` + `compiler/candy_evaluator/stdlib_builtins.go`.
- `16 Math` -> `abs/sqrt/pow/min/max/round/floor/ceil/sin/cos/tan/random/clamp/lerp` in `compiler/candy_evaluator/stdlib_builtins.go`.
- `17 Conversions` -> `int/float/string/bool` and `toInt/toFloat/toString/toBool` in `compiler/candy_evaluator/stdlib_builtins.go`.
- `18 Utility` -> `wait`, `exit`, `seconds`, `deltaTime` in `compiler/candy_evaluator/stdlib_builtins.go`.
- `19 File operations` -> `readFile/writeFile/appendFile/fileExists` and persistence helpers (`save/load`) in `compiler/candy_evaluator/builtin.go`, `compiler/candy_evaluator/stdlib_builtins.go`, and `compiler/candy_evaluator/prelude.go`.

### 15-19 Direct Runtime Tests

- `compiler/candy_evaluator/eval_test.go::TestEval_ChecklistBuiltins_MathAndConversions`
- `compiler/candy_evaluator/eval_test.go::TestEval_ChecklistBuiltins_FileAndPersistence`
- `compiler/candy_evaluator/eval_test.go::TestEval_ExitBuiltin_UsesExitHook`
- `compiler/candy_evaluator/eval_test.go::TestEval_SecondsAndDeltaTimeBuiltins`

## 20-30 Backend, Architecture, Build, Debug

- `20 C code generation pipeline` -> LLVM-based native pipeline in `compiler/candy_llvm` + `compiler/cmd/candy/main.go` + backend equivalence doc `docs/COMPILER_BACKEND_EQUIVALENCE.md`.
- `21 Generated-code error handling` -> diagnostics/reporting in `compiler/candy_report/report.go`, parser/typecheck diagnostics, and native build command error surfacing in `compiler/cmd/candy/main.go`.
- `22 Optimizations` -> `compiler/candy_opt/optimize.go` with tests in `compiler/candy_opt/optimize_test.go`.
- `23 Lexer` -> `compiler/candy_lexer/lexer.go` + `compiler/candy_lexer/lexer_test.go`.
- `24 Parser + AST` -> `compiler/candy_parser/*` + `compiler/candy_ast/*` with parser suites.
- `25 Semantic analyzer` -> `compiler/candy_typecheck/*` + typecheck tests.
- `26 Codegen + runtime integration` -> `compiler/candy_llvm/*` + `compiler/candy_load/*` + command integration.
- `27 CLI build/run/compile` -> `compiler/cmd/candy/main.go` and tests in `compiler/cmd/candy/main_test.go`.
- `28 Linking` -> manifest/link context path in `compiler/candy_load/context.go`, load in `compiler/candy_load/load.go`, and native clang invocation in `compiler/cmd/candy/main.go`.
- `29 Error messages` -> parser/typecheck/runtime diagnostics via `compiler/candy_report` + parser/typecheck tests.
- `30 Debugging support` -> `--debug` profile wiring in `compiler/cmd/candy/main.go` (native `-O0 -g` flow).

### Applicability Notes For 519-767

- Checklist items that say "generate C code" or "runtime.c" are treated as **backend-equivalent** in this repository because the production backend is LLVM IR + clang link, not a C transpiler.
- Checklist items naming setjmp/longjmp are interpreted as **error-handling strategy examples**, not mandatory mechanisms; equivalent diagnostics/runtime safety behavior is implemented.
- Summary/phase/feature-count lines (`720-767`) are informational and not executable requirements; no code changes are required for those lines.

### 20-30 Direct Toolchain/Architecture Tests

- Optimization proofs: `TestOptimizeProgram_ConstantFolding`, `TestOptimizeProgram_DeadBranchElimination`, `TestOptimizeProgram_RemoveNeverRunLoop`, `TestOptimizeProgram_InlinesSimplePureFunction`.
- Loader/bindgen integration: `TestWrapImportBuildPipeline`, `TestExpand_ImportedFileParseError`.
- Native backend + clang/toolchain: `TestE2E_ClangCompileAndRun`, `TestBundledClangPath_FromBundleRootEnv`, `TestResolveClangPath_PrefersOverride`, `TestResolveClangPath_InvalidOverride`.
- CLI build/run/flags: `TestRunBuild_WithOverrideClang`, `TestRunBuild_StdinWithImport_Errors`, `TestRunBuild_NoClang_ShowsGuidance`, `TestRunBuild_SubcommandAlias`, `TestRunBuild_SubcommandAlias_WithOutputFlagAfterFile`, `TestRunBuild_OptimizeFlagSelectsShippingProfile`, `TestAppendClangBuildContext_StaticLinking`, `TestRun_Repl`.
- candywrap CLI hardening: `TestRun_AllowsFlagsAfterHeaderPath`.
- Diagnostics/LSP/typecheck anchors: `TestAnalyzeSourceParseError`, `TestAnalyzeSourceTypecheckWarning`, `TestParseDiagnosticMessage`, `TestToLSPDiagnosticDefaults`, `TestServerLifecycleAndPublishDiagnostics`, `TestServerDidClosePublishesEmptyDiagnostics`, `TestServerDedupesUnchangedDiagnostics`, `TestStructInheritanceTypecheck`, `TestKotlinStyleDeclsTypecheck`.

## Test Suite Anchors

- Parser coverage: `compiler/candy_parser/parser_test.go`, `compiler/candy_parser/basic_parser_test.go`, `compiler/candy_parser/full_spec_test.go`, `compiler/candy_parser/interop_additions_test.go`, `compiler/candy_parser/library_syntax_test.go`.
- Runtime coverage: `compiler/candy_evaluator/eval_test.go`.
- Bindgen/interop coverage: `compiler/candy_bindgen/pipeline_test.go`, `compiler/candy_bindgen/header_parser_test.go`, `compiler/candy_load/load_test.go`.
- Optimization/codegen coverage: `compiler/candy_opt/optimize_test.go`, `compiler/candy_llvm/codegen_test.go`, `compiler/candy_llvm/e2e_clang_test.go`.
