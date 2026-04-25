package candy_bindgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteLibraryDocs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.md")
	api := &API{
		Functions: []Function{
			{Name: "add", ReturnType: "int", Params: []Parameter{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}}},
		},
		Constants: []Constant{{Name: "DEMO_VERSION", Value: "1"}},
	}
	m := &Manifest{Library: "demo"}
	if err := WriteLibraryDocs(path, "demo", api, m); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, "add(a: int, b: int)") {
		t.Fatalf("docs missing function signature: %s", s)
	}
}
