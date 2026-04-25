package candy_bindgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseManifestFileAndResolve(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "demo.candylib")
	content := `{
  "library":"demo",
  "externs":[{"name":"sum","return_type":"int","params":[{"name":"a","type":"int"},{"name":"b","type":"int"}]}],
  "compile":{"glue_sources":["demo_glue.c"],"include_dirs":["include"]},
  "link":{"lib_dirs":["lib"],"libs":["demo"]}
}`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := ParseManifestFile(p)
	if err != nil {
		t.Fatalf("ParseManifestFile: %v", err)
	}
	if len(m.Externs) != 1 || m.Externs[0].Name != "sum" {
		t.Fatalf("unexpected externs: %#v", m.Externs)
	}
	if len(m.Compile.GlueSources) != 1 || !filepath.IsAbs(m.Compile.GlueSources[0]) {
		t.Fatalf("expected absolute glue source, got: %#v", m.Compile.GlueSources)
	}
}

func TestUnsafeABI(t *testing.T) {
	ex := ExternBinding{
		Name:       "qsort",
		ReturnType: "void",
		Params: []ExternParam{
			{Name: "cmp", Type: "int (*)(const void*, const void*)"},
		},
	}
	if why := IsUnsafeABI(ex); why == "" {
		t.Fatal("expected unsafe ABI")
	}
}
