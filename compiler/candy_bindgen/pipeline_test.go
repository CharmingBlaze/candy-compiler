package candy_bindgen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"candy/candy_bindgen"
	"candy/candy_llvm"
	"candy/candy_load"
)

func TestWrapImportBuildPipeline(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "mylib.h")
	if err := os.WriteFile(header, []byte("int add(int a, int b);\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	externs, warnings, err := candy_bindgen.ParseHeaderFunctions(header)
	if err != nil {
		t.Fatalf("ParseHeaderFunctions: %v", err)
	}
	if len(warnings) > 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	manifest := candy_bindgen.BuildManifest("mylib", []string{header}, externs)
	manifest.Compile.GlueSources = []string{"mylib_glue.c"}
	manifest.Link.Libs = []string{"mylib"}
	manifestPath := filepath.Join(dir, "mylib.candylib")
	if err := candy_bindgen.WriteManifest(manifestPath, manifest); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}
	if err := candy_bindgen.WriteGlue(filepath.Join(dir, "mylib_glue.c"), "mylib", externs); err != nil {
		t.Fatalf("WriteGlue: %v", err)
	}

	mainPath := filepath.Join(dir, "main.candy")
	src := "import \"mylib.candylib\";\nreturn add(10, 32);"
	if err := os.WriteFile(mainPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	prog, ctx, err := candy_load.ExpandProgramForBuildWithContext(mainPath)
	if err != nil {
		t.Fatalf("ExpandProgramForBuildWithContext: %v", err)
	}
	ir, err := candy_llvm.New().GenerateIR(prog)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "declare i64 @add(i64, i64)") {
		t.Fatalf("extern missing in IR: %s", ir)
	}
	if ctx == nil || len(ctx.GlueSources) == 0 || len(ctx.Libs) == 0 {
		t.Fatalf("expected build context from manifest, got: %#v", ctx)
	}
}

func TestWrap_Box2DStyleNamespaceTransform(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "box2d.h")
	src := `
typedef struct b2World* b2WorldId;
typedef struct b2Body* b2BodyId;
b2WorldId b2World_Create(float gx, float gy);
void b2World_Destroy(b2WorldId world);
b2BodyId b2World_CreateBody(b2WorldId world, float x, float y);
void b2Internal_Debug(void);
`
	if err := os.WriteFile(header, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	api, _, err := candy_bindgen.ParseHeadersWithEngine([]string{header}, candy_bindgen.ParseOptions{}, candy_bindgen.ParserRegex)
	if err != nil {
		t.Fatalf("ParseHeadersWithEngine: %v", err)
	}
	_, err = candy_bindgen.TransformAPI(api, "box2d", []string{"b2World_"}, []string{"b2Internal_*"})
	if err != nil {
		t.Fatalf("TransformAPI: %v", err)
	}
	manifest := candy_bindgen.BuildManifestFromAPI("box2d", "box2d", []string{header}, api)
	if len(manifest.Externs) < 2 {
		t.Fatalf("expected extracted externs, got %d", len(manifest.Externs))
	}
	if manifest.Externs[0].Name == manifest.Externs[0].Symbol {
		t.Fatalf("expected transformed candy name distinct from C symbol")
	}
	if manifest.Namespace != "box2d" {
		t.Fatalf("expected namespace recorded")
	}
}
