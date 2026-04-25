package candy_llvm

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBundledClangPath_FromBundleRootEnv(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	clang := filepath.Join(bin, clangBinaryName())
	if err := os.WriteFile(clang, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(envBundleRoot, root)
	p, ok := bundledClangPath()
	if !ok {
		t.Fatal("expected bundled clang to be discovered")
	}
	if filepath.Clean(p) != filepath.Clean(clang) {
		t.Fatalf("expected %q, got %q", clang, p)
	}
}

func TestResolveClangPath_PrefersOverride(t *testing.T) {
	root := t.TempDir()
	override := filepath.Join(root, "custom-clang")
	if runtime.GOOS == "windows" {
		override += ".exe"
	}
	if err := os.WriteFile(override, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(envClangPath, override)
	got, err := ResolveClangPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(override) {
		t.Fatalf("expected override %q, got %q", override, got)
	}
}

func TestResolveClangPath_InvalidOverride(t *testing.T) {
	t.Setenv(envClangPath, filepath.Join(t.TempDir(), "missing-clang"))
	_, err := ResolveClangPath()
	if err == nil {
		t.Fatal("expected error for invalid CANDY_CLANG override")
	}
}

func TestBundledBinaryCandidates_IncludesToolchainAndLLVMLayouts(t *testing.T) {
	exeDir := filepath.Join("release", "bin")
	root := filepath.Join("bundle-root")
	got := bundledBinaryCandidates(exeDir, root, "clang")
	if len(got) == 0 {
		t.Fatal("expected candidate paths")
	}
	expectContains := []string{
		filepath.Clean(filepath.Join(root, "bin", "clang")),
		filepath.Clean(filepath.Join(root, "toolchain", "bin", "clang")),
		filepath.Clean(filepath.Join(root, "llvm", "bin", "clang")),
		filepath.Clean(filepath.Join(exeDir, "toolchain", "bin", "clang")),
		filepath.Clean(filepath.Join(exeDir, "llvm", "bin", "clang")),
		filepath.Clean(filepath.Join("release", "toolchain", "bin", "clang")),
		filepath.Clean(filepath.Join("release", "llvm", "bin", "clang")),
	}
	for _, want := range expectContains {
		found := false
		for _, p := range got {
			if p == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected candidate %q in %v", want, got)
		}
	}
}

func TestResolveClangPath_PrefersBundledOverPATH(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	bundled := filepath.Join(bin, clangBinaryName())
	if err := os.WriteFile(bundled, []byte("bundled"), 0o755); err != nil {
		t.Fatal(err)
	}

	pathDir := t.TempDir()
	pathClang := filepath.Join(pathDir, clangBinaryName())
	if err := os.WriteFile(pathClang, []byte("path"), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv(envClangPath, "")
	t.Setenv(envBundleRoot, root)
	t.Setenv("PATH", pathDir)
	got, err := ResolveClangPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(bundled) {
		t.Fatalf("expected bundled clang %q, got %q", bundled, got)
	}
}

func TestResolveOptPath_PrefersOverrideThenBundledThenPATH(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	bundled := filepath.Join(bin, optBinaryName())
	if err := os.WriteFile(bundled, []byte("bundled"), 0o755); err != nil {
		t.Fatal(err)
	}

	pathDir := t.TempDir()
	pathOpt := filepath.Join(pathDir, optBinaryName())
	if err := os.WriteFile(pathOpt, []byte("path"), 0o755); err != nil {
		t.Fatal(err)
	}

	overrideDir := t.TempDir()
	override := filepath.Join(overrideDir, "custom-opt")
	if runtime.GOOS == "windows" {
		override += ".exe"
	}
	if err := os.WriteFile(override, []byte("override"), 0o755); err != nil {
		t.Fatal(err)
	}

	// override wins
	t.Setenv(envBundleRoot, root)
	t.Setenv("PATH", pathDir)
	t.Setenv(envOptPath, override)
	got, err := ResolveOptPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(override) {
		t.Fatalf("expected override opt %q, got %q", override, got)
	}

	// bundled wins when override missing
	t.Setenv(envOptPath, "")
	got, err = ResolveOptPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(bundled) {
		t.Fatalf("expected bundled opt %q, got %q", bundled, got)
	}

	// PATH wins when bundled missing
	t.Setenv(envBundleRoot, filepath.Join(t.TempDir(), "missing"))
	got, err = ResolveOptPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(pathOpt) {
		t.Fatalf("expected PATH opt %q, got %q", pathOpt, got)
	}
}
