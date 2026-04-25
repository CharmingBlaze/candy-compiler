package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"candy/candy_evaluator"
	"candy/candy_lexer"
	"candy/candy_llvm"
	"candy/candy_load"
	"candy/candy_opt"
	"candy/candy_parser"
	"candy/candy_physics"
	"candy/candy_raylib"
	"candy/candy_report"
	"candy/candy_typecheck"
)

var (
	BuildVersion    = "dev"
	BuildStdlibHash = ""
)

func main() {
	candy_raylib.RegisterBuiltins()
	candy_physics.RegisterBuiltins()
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	args = normalizeCLIArgs(args)
	if len(args) > 0 && strings.EqualFold(strings.TrimSpace(args[0]), "doctor") {
		return runDoctor(stdout)
	}
	fs := flag.NewFlagSet("candy", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Candy CLI")
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "Commands:")
		fmt.Fprintln(stderr, "  candy run <file.candy>            Run through evaluator")
		fmt.Fprintln(stderr, "  candy build <file.candy>          Build native binary (LLVM + clang)")
		fmt.Fprintln(stderr, "  candy compile <file.candy>        Alias of build")
		fmt.Fprintln(stderr, "  candy doctor                      Check native toolchain (clang/opt)")
		fmt.Fprintln(stderr, "  candy init <project_name>         Create starter project")
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "Flags:")
		fs.PrintDefaults()
	}
	printAST := fs.Bool("ast", false, "print AST string and exit 0 (no run)")
	staticCheck := fs.Bool("check", false, "run static checker after parse; prints issues, non-fatal")
	buildNative := fs.Bool("build", false, "generate LLVM IR (.ll file)")
	debugBuild := fs.Bool("debug", false, "debug profile for -build (minimal optimization)")
	releaseBuild := fs.Bool("release", false, "shipping profile for -build (aggressive optimization)")
	optimizeBuild := fs.Bool("optimize", false, "enable shipping optimization profile for -build")
	outputPath := fs.String("o", "", "output binary path for -build")
	verbose := fs.Bool("verbose", false, "print compilation steps")
	interactive := fs.Bool("i", false, "start interactive read-eval-print loop (REPL), ignoring file arg")
	interactiveLong := fs.Bool("repl", false, "same as -i")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		fmt.Fprintln(stderr, err)
		return 2
	}
	path := ""
	tail := fs.Args()
	if len(tail) > 0 {
		path = tail[0]
	}

	// Command: init <project_name>
	if path == "init" {
		if len(tail) < 2 {
			fmt.Fprintln(stderr, "usage: candy init <project_name>")
			return 1
		}
		return initProject(tail[1], stderr)
	}

	if *interactive || *interactiveLong {
		if path != "" {
			fmt.Fprintln(stderr, "candy: -i / -repl is interactive; the file argument is ignored.")
		}
		return runREPL(stdin, stdout, stderr)
	}

	profile := candy_llvm.BuildDevRelease
	if *debugBuild {
		profile = candy_llvm.BuildDebug
	}
	if *releaseBuild {
		profile = candy_llvm.BuildShipping
	}
	if *optimizeBuild && !*debugBuild && !*releaseBuild {
		profile = candy_llvm.BuildShipping
	}

	// Native build: expand imports for a file path; skip double-parse/merge for stdin.
	if *buildNative && path != "" {
		probe := candy_llvm.ProbeToolchain()
		if !probe.Clang.Found {
			printToolchainProbe(stdout, probe)
			return 1
		}
		program, buildCtx, lerr := candy_load.ExpandProgramForBuildWithContext(path)
		if lerr != nil {
			fmt.Fprintln(stderr, lerr)
			return 1
		}
		if *staticCheck {
			issues := candy_typecheck.CheckProgram(program)
			if len(issues) > 0 {
				mergedSrc := program.String() + "\n"
				candy_report.Report(mergedSrc, issues)
			}
		}
		if *printAST {
			fmt.Fprint(stdout, program.String())
			if !strings.HasSuffix(program.String(), "\n") {
				fmt.Fprintln(stdout)
			}
			return 0
		}
		candy_opt.OptimizeProgram(program)
		comp := candy_llvm.New()
		ir, genErr := comp.GenerateIR(program)
		if genErr != nil {
			fmt.Fprintln(stderr, "codegen failed:", genErr)
			return 1
		}
		optimizedIR, optErr := candy_llvm.OptimizeIR(ir, profile)
		if optErr != nil {
			fmt.Fprintln(stderr, "warning: optimizer step failed, using unoptimized IR:", optErr)
		} else {
			ir = optimizedIR
		}
		outPath := "output.ll"
		if path != "" {
			outPath = strings.TrimSuffix(path, filepath.Ext(path)) + ".ll"
		}
		err := os.WriteFile(outPath, []byte(ir), 0644)
		if err != nil {
			fmt.Fprintln(stderr, "failed to write IR:", err)
			return 1
		}
		fmt.Fprintln(stdout, "Generated LLVM IR:", outPath)

		// Attempt to compile with clang. Prefer bundled LLVM toolchain.
		exePath := inferExePath(*outputPath, outPath)
		if *verbose {
			fmt.Fprintf(stdout, "Build profile: %v\n", profile)
			if buildCtx != nil {
				fmt.Fprintf(stdout, "Extra glue sources: %d, libs: %d\n", len(buildCtx.GlueSources), len(buildCtx.Libs))
			}
		}

		clangPath, err := candy_llvm.ResolveClangPath()
		if err != nil {
			fmt.Fprintf(stdout, "\nNo clang found. Candy checks: CANDY_CLANG, bundled ./llvm/bin, then PATH.\n")
			fmt.Fprintf(stdout, "To produce a native binary, bundle LLVM or install it, then run:\n  clang %s -o %s\n", outPath, exePath)
			return 0
		}
		clangArgs := []string{outPath, "-o", exePath}
		clangArgs = appendClangBuildContext(clangArgs, buildCtx)
		switch profile {
		case candy_llvm.BuildDebug:
			clangArgs = append(clangArgs, "-O0", "-g")
		case candy_llvm.BuildShipping:
			clangArgs = append(clangArgs, "-O3")
		default:
			clangArgs = append(clangArgs, "-O2")
		}
		cmd := exec.Command(clangPath, clangArgs...)
		if *verbose {
			fmt.Fprintf(stdout, "clang: %s %s\n", clangPath, strings.Join(clangArgs, " "))
		}
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(stdout, "\nclang invocation failed (%s): %v\n", clangPath, err)
			fmt.Fprintf(stdout, "You can rerun manually:\n  %s %s\n", clangPath, strings.Join(clangArgs, " "))
		} else {
			fmt.Fprintln(stdout, "Generated Native Binary:", exePath)
		}
		return 0
	}

	// Read and parse: interpreter, -build from stdin, or -ast/-check on a file.
	input, err := readSource(path, stdin)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	l := candy_lexer.New(string(input))
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		candy_report.Report(string(input), p.Errors())
		return 1
	}

	if *buildNative {
		probe := candy_llvm.ProbeToolchain()
		if !probe.Clang.Found {
			printToolchainProbe(stdout, probe)
			return 1
		}
		if candy_load.HasTopLevelImport(program) {
			fmt.Fprintln(stderr, "candy: -build with stdin does not support import; pass a .candy file so paths resolve")
			return 1
		}
		if *staticCheck {
			issues := candy_typecheck.CheckProgram(program)
			if len(issues) > 0 {
				candy_report.Report(string(input), issues)
			}
		}
		if *printAST {
			fmt.Fprint(stdout, program.String())
			if !strings.HasSuffix(program.String(), "\n") {
				fmt.Fprintln(stdout)
			}
			return 0
		}
		candy_opt.OptimizeProgram(program)
		comp := candy_llvm.New()
		ir, genErr := comp.GenerateIR(program)
		if genErr != nil {
			fmt.Fprintln(stderr, "codegen failed:", genErr)
			return 1
		}
		optimizedIR, optErr := candy_llvm.OptimizeIR(ir, profile)
		if optErr != nil {
			fmt.Fprintln(stderr, "warning: optimizer step failed, using unoptimized IR:", optErr)
		} else {
			ir = optimizedIR
		}
		outPath := "output.ll"
		err := os.WriteFile(outPath, []byte(ir), 0644)
		if err != nil {
			fmt.Fprintln(stderr, "failed to write IR:", err)
			return 1
		}
		fmt.Fprintln(stdout, "Generated LLVM IR:", outPath)
		exePath := inferExePath(*outputPath, outPath)
		clangPath, err := candy_llvm.ResolveClangPath()
		if err != nil {
			fmt.Fprintf(stdout, "\nNo clang found. Candy checks: CANDY_CLANG, bundled ./llvm/bin, then PATH.\n")
			fmt.Fprintf(stdout, "To produce a native binary, bundle LLVM or install it, then run:\n  clang %s -o %s\n", outPath, exePath)
			return 0
		}
		clangArgs := []string{outPath, "-o", exePath}
		switch profile {
		case candy_llvm.BuildDebug:
			clangArgs = append(clangArgs, "-O0", "-g")
		case candy_llvm.BuildShipping:
			clangArgs = append(clangArgs, "-O3")
		default:
			clangArgs = append(clangArgs, "-O2")
		}
		cmd := exec.Command(clangPath, clangArgs...)
		if *verbose {
			fmt.Fprintf(stdout, "clang: %s %s\n", clangPath, strings.Join(clangArgs, " "))
		}
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(stdout, "\nclang invocation failed (%s): %v\n", clangPath, err)
			fmt.Fprintf(stdout, "You can rerun manually:\n  %s %s\n", clangPath, strings.Join(clangArgs, " "))
		} else {
			fmt.Fprintln(stdout, "Generated Native Binary:", exePath)
		}
		return 0
	}

	if *staticCheck {
		issues := candy_typecheck.CheckProgram(program)
		if len(issues) > 0 {
			candy_report.Report(string(input), issues)
		}
	}
	if *printAST {
		fmt.Fprint(stdout, program.String())
		if !strings.HasSuffix(program.String(), "\n") {
			fmt.Fprintln(stdout)
		}
		return 0
	}

	v, eerr := candy_evaluator.Eval(program, nil)
	if eerr != nil {
		fmt.Fprintln(stderr, eerr)
		return 1
	}
	// Avoid a stray "null" line when the last statement is a call/print (returns null).
	if v != nil && v.Kind != candy_evaluator.ValNull {
		fmt.Fprintln(stdout, v.String())
	}
	return 0
}

func runDoctor(stdout io.Writer) int {
	probe := candy_llvm.ProbeToolchain()
	printToolchainProbe(stdout, probe)
	if !probe.Clang.Found {
		return 1
	}
	return 0
}

func printToolchainProbe(stdout io.Writer, probe candy_llvm.ProbeResult) {
	fmt.Fprintln(stdout, "Candy toolchain doctor")
	if strings.TrimSpace(BuildVersion) != "" {
		fmt.Fprintf(stdout, "Build version: %s\n", strings.TrimSpace(BuildVersion))
	}
	if strings.TrimSpace(BuildStdlibHash) != "" {
		fmt.Fprintf(stdout, "Stdlib hash: %s\n", strings.TrimSpace(BuildStdlibHash))
	}
	if probe.Compatible {
		fmt.Fprintln(stdout, "Status: PASS")
	} else {
		fmt.Fprintln(stdout, "Status: FAIL")
	}
	if len(probe.SearchPolicy) > 0 {
		fmt.Fprintln(stdout, "Search order:")
		for _, p := range probe.SearchPolicy {
			fmt.Fprintf(stdout, "  - %s\n", p)
		}
	}
	fmt.Fprintf(stdout, "- clang: found=%v path=%s\n", probe.Clang.Found, probe.Clang.Path)
	if probe.Clang.Version != "" {
		fmt.Fprintf(stdout, "  version: %s\n", probe.Clang.Version)
	}
	if probe.Clang.Error != "" {
		fmt.Fprintf(stdout, "  error: %s\n", probe.Clang.Error)
	}
	fmt.Fprintf(stdout, "- opt: found=%v path=%s\n", probe.Opt.Found, probe.Opt.Path)
	if probe.Opt.Version != "" {
		fmt.Fprintf(stdout, "  version: %s\n", probe.Opt.Version)
	}
	if probe.Opt.Error != "" {
		fmt.Fprintf(stdout, "  error: %s\n", probe.Opt.Error)
	}
	if len(probe.Problems) > 0 {
		fmt.Fprintln(stdout, "Problems:")
		for _, p := range probe.Problems {
			fmt.Fprintf(stdout, "  - %s\n", p)
		}
	}
	if len(probe.Suggestions) > 0 {
		fmt.Fprintln(stdout, "Suggestions:")
		for _, s := range probe.Suggestions {
			fmt.Fprintf(stdout, "  - %s\n", s)
		}
	}
}

func inferExePath(outputPath, irPath string) string {
	if strings.TrimSpace(outputPath) != "" {
		if filepath.Ext(outputPath) == "" && os.PathSeparator == '\\' {
			return outputPath + ".exe"
		}
		return outputPath
	}
	exePath := strings.TrimSuffix(irPath, ".ll")
	if filepath.Ext(exePath) == "" && os.PathSeparator == '\\' {
		exePath += ".exe"
	}
	return exePath
}

func normalizeCLIArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}
	cmd := strings.ToLower(strings.TrimSpace(args[0]))
	switch cmd {
	case "build", "compile":
		rest := args[1:]
		file := ""
		if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
			file = rest[0]
			rest = rest[1:]
		}
		out := append([]string{"-build"}, rest...)
		if file != "" {
			out = append(out, file)
		}
		return out
	case "run":
		return args[1:]
	default:
		return args
	}
}

func appendClangBuildContext(clangArgs []string, buildCtx *candy_load.BuildContext) []string {
	if buildCtx == nil {
		return clangArgs
	}
	clangArgs = append(clangArgs, buildCtx.GlueSources...)
	for _, dir := range buildCtx.IncludeDirs {
		clangArgs = append(clangArgs, "-I"+dir)
	}
	for _, flag := range buildCtx.CFlags {
		clangArgs = append(clangArgs, flag)
	}
	for _, dir := range buildCtx.LibDirs {
		clangArgs = append(clangArgs, "-L"+dir)
	}
	for _, lib := range buildCtx.Libs {
		clangArgs = append(clangArgs, "-l"+lib)
	}
	if buildCtx.Static {
		clangArgs = append(clangArgs, "-static")
	}
	if len(buildCtx.StaticLibs) > 0 {
		clangArgs = append(clangArgs, "-Wl,-Bstatic")
		for _, lib := range buildCtx.StaticLibs {
			clangArgs = append(clangArgs, "-l"+lib)
		}
		clangArgs = append(clangArgs, "-Wl,-Bdynamic")
	}
	for _, flag := range buildCtx.LDFlags {
		clangArgs = append(clangArgs, flag)
	}
	return clangArgs
}

func readSource(path string, stdin io.Reader) ([]byte, error) {
	if path == "" {
		b, err := io.ReadAll(stdin)
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			return nil, fmt.Errorf("no input (provide a file or pipe stdin)")
		}
		return b, nil
	}
	return os.ReadFile(path)
}

func initProject(name string, stderr io.Writer) int {
	err := os.Mkdir(name, 0755)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create directory %s: %v\n", name, err)
		return 1
	}
	os.Mkdir(filepath.Join(name, "assets"), 0755)

	mainContent := `// ` + name + ` - A Candy Game
window(800, 600, "` + name + `")

while !shouldClose() {
    clear("black")
    text(20, 20, "Welcome to ` + name + `!", "white")
    flip()
}

closeWindow()
`
	err = os.WriteFile(filepath.Join(name, "main.candy"), []byte(mainContent), 0644)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create main.candy: %v\n", err)
		return 1
	}

	fmt.Printf("Created project %s\nRun it with: candy ./%s/main.candy\n", name, name)
	return 0
}
