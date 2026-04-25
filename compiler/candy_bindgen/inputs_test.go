package candy_bindgen

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExpandInputFiles_DirectoryAndPattern(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.h"), []byte("int a();"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.cpp"), []byte("int b(){return 1;}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := ExpandInputFiles([]string{dir, filepath.Join(dir, "*.h")})
	if len(got) != 2 {
		t.Fatalf("expected 2 C/C++ files, got %d: %v", len(got), got)
	}
}

func TestDiscoverLibraryFiles(t *testing.T) {
	dir := t.TempDir()
	inc := filepath.Join(dir, "include")
	src := filepath.Join(dir, "src")
	if err := os.MkdirAll(inc, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	h := filepath.Join(inc, "lib.h")
	c := filepath.Join(src, "lib.c")
	if err := os.WriteFile(h, []byte("int lib_add(int a, int b);"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(c, []byte("int lib_add(int a, int b){return a+b;}"), 0o644); err != nil {
		t.Fatal(err)
	}
	headers, sources, includeDirs, err := DiscoverLibraryFiles([]string{dir}, "c")
	if err != nil {
		t.Fatalf("DiscoverLibraryFiles: %v", err)
	}
	if len(headers) != 1 || len(sources) != 1 {
		t.Fatalf("unexpected discovered files: h=%v s=%v", headers, sources)
	}
	if len(includeDirs) == 0 {
		t.Fatalf("expected include dirs")
	}
}

func TestExpandInputFiles_DeterministicAcrossInputOrder(t *testing.T) {
	dir := t.TempDir()
	d1 := filepath.Join(dir, "a")
	d2 := filepath.Join(dir, "b")
	if err := os.MkdirAll(d1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(d2, 0o755); err != nil {
		t.Fatal(err)
	}
	f1 := filepath.Join(d1, "x.h")
	f2 := filepath.Join(d2, "y.c")
	if err := os.WriteFile(f1, []byte("int x();"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("int y(){return 1;}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got1 := ExpandInputFiles([]string{d1, d2})
	got2 := ExpandInputFiles([]string{d2, d1})
	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("expected deterministic ordering; got1=%v got2=%v", got1, got2)
	}
}

func TestDiscoverLibraryFiles_DeterministicAcrossRootOrder(t *testing.T) {
	dir := t.TempDir()
	r1 := filepath.Join(dir, "lib1")
	r2 := filepath.Join(dir, "lib2")
	if err := os.MkdirAll(filepath.Join(r1, "include"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(r2, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	h1 := filepath.Join(r1, "include", "a.h")
	h2 := filepath.Join(r2, "b.h")
	s2 := filepath.Join(r2, "src", "b.c")
	if err := os.WriteFile(h1, []byte("int a();"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(h2, []byte("int b();"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s2, []byte("int b(){return 2;}"), 0o644); err != nil {
		t.Fatal(err)
	}

	hA, sA, iA, err := DiscoverLibraryFiles([]string{r1, r2}, "c")
	if err != nil {
		t.Fatal(err)
	}
	hB, sB, iB, err := DiscoverLibraryFiles([]string{r2, r1}, "c")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hA, hB) || !reflect.DeepEqual(sA, sB) || !reflect.DeepEqual(iA, iB) {
		t.Fatalf("expected deterministic discovery\nA: h=%v s=%v i=%v\nB: h=%v s=%v i=%v", hA, sA, iA, hB, sB, iB)
	}
}
