package candy_load

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"candy/candy_llvm"
)

func TestExpand_TwoFileImport_GeneratesF(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, "lib.candy")
	if err := os.WriteFile(lib, []byte("fun f(): int { return 1; }"), 0o644); err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(dir, "main.candy")
	if err := os.WriteFile(mainPath, []byte("import \"lib.candy\";\nreturn f();"), 0o644); err != nil {
		t.Fatal(err)
	}
	prog, err := ExpandProgramForBuild(mainPath)
	if err != nil {
		t.Fatalf("ExpandProgramForBuild: %v", err)
	}
	if len(prog.Statements) < 1 {
		t.Fatalf("empty program")
	}
	comp := candy_llvm.New()
	ir, gerr := comp.GenerateIR(prog)
	if gerr != nil {
		t.Fatalf("GenerateIR: %v", gerr)
	}
	if !strings.Contains(ir, "define i64 @f()") {
		t.Fatalf("expected f in IR, got:\n%s", ir)
	}
}

func TestExpand_ImportCandyLib_SynthExternAndContext(t *testing.T) {
	dir := t.TempDir()
	libPath := filepath.Join(dir, "mylib.candylib")
	libJSON := `{
  "library":"mylib",
  "externs":[{"name":"native_add","return_type":"int","params":[{"name":"a","type":"int"},{"name":"b","type":"int"}]}],
  "compile":{"glue_sources":["mylib_glue.c"],"include_dirs":["include"],"cflags":["-DMYLIB=1"]},
  "link":{"lib_dirs":["lib"],"libs":["mylib"],"ldflags":["-Wl,--as-needed"]}
}`
	if err := os.WriteFile(libPath, []byte(libJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(dir, "main.candy")
	mainSrc := "import \"mylib.candylib\";\nreturn native_add(1, 2);"
	if err := os.WriteFile(mainPath, []byte(mainSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	prog, ctx, err := ExpandProgramForBuildWithContext(mainPath)
	if err != nil {
		t.Fatalf("ExpandProgramForBuildWithContext: %v", err)
	}
	if prog == nil || len(prog.Statements) == 0 {
		t.Fatalf("expected synthesized statements")
	}
	comp := candy_llvm.New()
	ir, gerr := comp.GenerateIR(prog)
	if gerr != nil {
		t.Fatalf("GenerateIR: %v", gerr)
	}
	if !strings.Contains(ir, "declare i64 @native_add(i64, i64)") {
		t.Fatalf("expected extern in IR, got:\n%s", ir)
	}
	if ctx == nil || len(ctx.GlueSources) == 0 || len(ctx.Libs) == 0 {
		t.Fatalf("expected non-empty build context, got: %#v", ctx)
	}
}

func TestExpand_ImportCycle_NoInfiniteLoop(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.candy")
	b := filepath.Join(dir, "b.candy")
	// b imports a; a imports b
	if err := os.WriteFile(b, []byte("import \"a.candy\";\nreturn 0;"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(a, []byte("import \"b.candy\";\nreturn 0;"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ExpandProgramForBuild(a)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestExpand_ImportedFileParseError(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.candy")
	if err := os.WriteFile(bad, []byte("this is not valid candy {"), 0o644); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(dir, "entry.candy")
	if err := os.WriteFile(entry, []byte("import \"bad.candy\";\nreturn 0;"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ExpandProgramForBuild(entry)
	if err == nil {
		t.Fatal("expected error from bad import")
	}
	if !strings.Contains(err.Error(), "bad.candy") && !strings.Contains(err.Error(), "candy_load:") {
		t.Fatalf("error should mention load/parse, got: %v", err)
	}
}
