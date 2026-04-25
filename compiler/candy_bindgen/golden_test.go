package candy_bindgen

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGoldens = flag.Bool("goldens:update", false, "update candy_bindgen golden files")

func TestGolden_ManifestAndGlue(t *testing.T) {
	t.Setenv(envFixedGeneratedAt, "2026-01-01T00:00:00Z")

	api := &API{
		Functions: []Function{
			{
				Name:       "demo_add",
				Symbol:     "demo_add",
				ReturnType: "int",
				Params: []Parameter{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "int"},
				},
			},
			{
				Name:       "demo_print",
				Symbol:     "demo_print",
				ReturnType: "void",
				Params: []Parameter{
					{Name: "message", Type: "const char*"},
				},
			},
		},
	}

	manifest := BuildManifestFromAPI("demo", "demo", []string{"demo.h"}, api)
	dir := t.TempDir()
	gotManifest := filepath.Join(dir, "demo.candylib")
	gotGlue := filepath.Join(dir, "demo_glue.c")
	if err := WriteManifest(gotManifest, manifest); err != nil {
		t.Fatal(err)
	}
	if err := WriteGlue(gotGlue, "demo", manifest.Externs); err != nil {
		t.Fatal(err)
	}

	expectManifest := filepath.Join("testdata", "golden", "demo.candylib")
	expectGlue := filepath.Join("testdata", "golden", "demo_glue.c")
	assertGoldenFile(t, expectManifest, gotManifest)
	assertGoldenFile(t, expectGlue, gotGlue)
}

func TestGolden_TransformCollisionAndReservedName(t *testing.T) {
	t.Setenv(envFixedGeneratedAt, "2026-01-01T00:00:00Z")
	api := &API{
		Functions: []Function{
			{Name: "b2World_Create", Symbol: "b2World_Create", ReturnType: "int"},
			{Name: "b2Body_Create", Symbol: "b2Body_Create", ReturnType: "int"},
			{Name: "b2World_end", Symbol: "b2World_end", ReturnType: "void"},
		},
	}
	warns, err := TransformAPI(api, "box2d", []string{"b2World_", "b2Body_"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = warns
	manifest := BuildManifestFromAPI("box2d", "box2d", []string{"box2d.h"}, api)
	dir := t.TempDir()
	gotManifest := filepath.Join(dir, "box2d.candylib")
	gotGlue := filepath.Join(dir, "box2d_glue.c")
	if err := WriteManifest(gotManifest, manifest); err != nil {
		t.Fatal(err)
	}
	if err := WriteGlue(gotGlue, "box2d", manifest.Externs); err != nil {
		t.Fatal(err)
	}
	assertGoldenFile(t, filepath.Join("testdata", "golden", "box2d.candylib"), gotManifest)
	assertGoldenFile(t, filepath.Join("testdata", "golden", "box2d_glue.c"), gotGlue)
}

func TestGolden_UnsafeABIVariadicManifest(t *testing.T) {
	t.Setenv(envFixedGeneratedAt, "2026-01-01T00:00:00Z")
	api := &API{
		Functions: []Function{
			{
				Name:       "printf",
				Symbol:     "printf",
				ReturnType: "int",
				Variadic:   true,
				Params: []Parameter{
					{Name: "fmt", Type: "const char*"},
				},
			},
		},
	}
	manifest := BuildManifestFromAPI("cstd", "", []string{"stdio.h"}, api)
	manifest.UnsafeABI = true
	dir := t.TempDir()
	gotManifest := filepath.Join(dir, "cstd.candylib")
	if err := WriteManifest(gotManifest, manifest); err != nil {
		t.Fatal(err)
	}
	assertGoldenFile(t, filepath.Join("testdata", "golden", "cstd_unsafe_abi.candylib"), gotManifest)
}

func assertGoldenFile(t *testing.T, expectedPath, gotPath string) {
	t.Helper()
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatal(err)
	}
	if *updateGoldens {
		if os.Getenv("CI") != "" {
			t.Fatalf("refusing to update goldens in CI")
		}
		if err := os.MkdirAll(filepath.Dir(expectedPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(expectedPath, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with -goldens:update to create/update)", expectedPath, err)
	}
	if string(want) != string(got) {
		t.Fatalf("golden mismatch for %s (run with -goldens:update)", expectedPath)
	}
}

