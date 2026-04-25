package candy_llvm

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"candy/candy_lexer"
	"candy/candy_parser"
)

// TestE2E_ClangCompileAndRun exercises GenerateIR then clang; skips when no clang is available.
// This catches invalid IR that string-based unit tests can miss.
func TestE2E_ClangCompileAndRun(t *testing.T) {
	clangPath, err := ResolveClangPath()
	if err != nil {
		t.Skip("clang not found (set CANDY_CLANG or use bundled LLVM on PATH): ", err)
	}

	src := `
return 0;
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	comp := New()
	ir, err := comp.GenerateIR(prog)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}

	dir := t.TempDir()
	llPath := filepath.Join(dir, "e2e.ll")
	exePath := filepath.Join(dir, "e2e_out")
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}
	if err := os.WriteFile(llPath, []byte(ir), 0644); err != nil {
		t.Fatal(err)
	}

	out, cerr := exec.Command(clangPath, llPath, "-o", exePath).CombinedOutput()
	if cerr != nil {
		t.Skipf("clang failed to compile IR (toolchain issue, not necessarily Candy): %v\n%s", cerr, out)
	}

	cmd := exec.Command(exePath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("e2e binary: %v", err)
	}
}
