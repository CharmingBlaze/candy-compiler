package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_AllowsFlagsAfterHeaderPath(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "demo.h")
	outDir := filepath.Join(dir, "out")
	if err := os.WriteFile(header, []byte("int add(int a, int b);\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code := run([]string{
		"wrap",
		header,
		"--name", "mylib",
		"--output", outDir,
		"--docs=false",
		"--stub=false",
	}, os.Stdout, os.Stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, want 0", code)
	}
	if _, err := os.Stat(filepath.Join(outDir, "mylib.candylib")); err != nil {
		t.Fatalf("expected manifest generated with --name override: %v", err)
	}
}
