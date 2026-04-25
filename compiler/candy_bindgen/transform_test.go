package candy_bindgen

import "testing"

func TestTransformAPI_NamespaceStripIgnore(t *testing.T) {
	api := &API{
		Functions: []Function{
			{Name: "b2World_Create", Symbol: "b2World_Create", ReturnType: "int"},
			{Name: "b2World_Destroy", Symbol: "b2World_Destroy", ReturnType: "void"},
			{Name: "b2Internal_Debug", Symbol: "b2Internal_Debug", ReturnType: "void"},
		},
	}
	warns, err := TransformAPI(api, "box2d", []string{"b2World_"}, []string{"b2Internal_*"})
	if err != nil {
		t.Fatalf("TransformAPI: %v", err)
	}
	if len(warns) == 0 {
		t.Fatalf("expected ignore warning")
	}
	if len(api.Functions) != 2 {
		t.Fatalf("expected 2 functions after ignore, got %d", len(api.Functions))
	}
	if api.Functions[0].Name != "box2d_Create" {
		t.Fatalf("unexpected transformed name: %s", api.Functions[0].Name)
	}
	if api.Functions[0].Symbol != "b2World_Create" {
		t.Fatalf("symbol should preserve original C symbol, got %s", api.Functions[0].Symbol)
	}
}

func TestWriteCXXShimTemplate(t *testing.T) {
	path := t.TempDir() + "/shim.cpp"
	m := &Manifest{
		Externs: []ExternBinding{
			{
				Name:       "box2d_createWorld",
				Symbol:     "box2d_createWorld",
				ReturnType: "void*",
				Params:     []ExternParam{{Name: "gx", Type: "float"}, {Name: "gy", Type: "float"}},
			},
		},
	}
	if err := WriteCXXShimTemplate(path, "box2d", m); err != nil {
		t.Fatalf("WriteCXXShimTemplate: %v", err)
	}
}
