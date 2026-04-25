package candy_bindgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHeaderFunctions(t *testing.T) {
	dir := t.TempDir()
	h := filepath.Join(dir, "demo.h")
	src := `
int add(int a, int b);
const char* name_of(int id);
int printf(const char* fmt, ...);
`
	if err := os.WriteFile(h, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	fns, warns, err := ParseHeaderFunctions(h)
	if err != nil {
		t.Fatalf("ParseHeaderFunctions: %v", err)
	}
	if len(fns) != 2 {
		t.Fatalf("expected 2 safe funcs, got %d (%v)", len(fns), warns)
	}
}

func TestParseHeaders_UnsafeABIIncludesVariadic(t *testing.T) {
	dir := t.TempDir()
	h := filepath.Join(dir, "demo.h")
	src := `
int add(int a, int b);
int printf(const char* fmt, ...);
`
	if err := os.WriteFile(h, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	api, warns, err := ParseHeaders([]string{h}, ParseOptions{UnsafeABI: true})
	if err != nil {
		t.Fatalf("ParseHeaders: %v", err)
	}
	if len(api.Functions) != 2 {
		t.Fatalf("expected 2 funcs in unsafe mode, got %d (%v)", len(api.Functions), warns)
	}
	var foundVariadic bool
	for _, fn := range api.Functions {
		if fn.Name == "printf" {
			foundVariadic = true
			if !fn.Variadic {
				t.Fatalf("expected printf variadic=true")
			}
			if len(fn.Params) != 1 {
				t.Fatalf("expected variadic params to exclude ellipsis, got %d", len(fn.Params))
			}
		}
	}
	if !foundVariadic {
		t.Fatalf("expected printf in parsed functions")
	}
}
