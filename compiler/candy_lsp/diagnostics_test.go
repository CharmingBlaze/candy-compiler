package candy_lsp

import "testing"

func TestAnalyzeSourceParseError(t *testing.T) {
	diags := AnalyzeSource("val x = ;")
	if len(diags) == 0 {
		t.Fatal("expected diagnostics for parse error")
	}
}

func TestAnalyzeSourceTypecheckWarning(t *testing.T) {
	diags := AnalyzeSource("return true + 1;")
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}
}
