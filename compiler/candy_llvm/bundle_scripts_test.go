package candy_llvm

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBundleScript_Sh(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash script smoke test skipped on windows")
	}
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root = filepath.Dir(root)
	llvm := mkFakeLLVM(t)
	candyBin := filepath.Join(t.TempDir(), "candy")
	if err := os.WriteFile(candyBin, []byte("candy"), 0o755); err != nil {
		t.Fatal(err)
	}
	candywrapBin := filepath.Join(t.TempDir(), "candywrap")
	if err := os.WriteFile(candywrapBin, []byte("candywrap"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(t.TempDir(), "manifest.txt")
	if err := os.WriteFile(manifest, []byte("bin/clang\nlib/*.so*\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(t.TempDir(), "bundle")
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "bundle-llvm.sh"), candyBin, candywrapBin, llvm, outDir)
	cmd.Env = append(os.Environ(), "CANDY_LLVM_MANIFEST="+manifest)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bundle-llvm.sh failed: %v\n%s", err, string(out))
	}
	if _, err := os.Stat(filepath.Join(outDir, "llvm", "bin", "clang")); err != nil {
		t.Fatalf("missing bundled clang: %v", err)
	}
}

func TestBundleScript_PowerShell(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("powershell script smoke test runs on windows")
	}
	powershell := "powershell"
	if _, err := exec.LookPath(powershell); err != nil {
		t.Skip("powershell not available")
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root = filepath.Dir(root)
	llvm := mkFakeLLVM(t)
	candyBin := filepath.Join(t.TempDir(), "candy.exe")
	if err := os.WriteFile(candyBin, []byte("candy"), 0o755); err != nil {
		t.Fatal(err)
	}
	candywrapBin := filepath.Join(t.TempDir(), "candywrap.exe")
	if err := os.WriteFile(candywrapBin, []byte("candywrap"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(t.TempDir(), "manifest.txt")
	if err := os.WriteFile(manifest, []byte("bin/clang.exe\nlib/*.dll\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(t.TempDir(), "bundle")
	cmd := exec.Command(powershell, "-ExecutionPolicy", "Bypass", "-File", filepath.Join(root, "scripts", "bundle-llvm.ps1"), "-CandyBinary", candyBin, "-CandywrapBinary", candywrapBin, "-LlvmRoot", llvm, "-OutDir", outDir)
	cmd.Env = append(os.Environ(), "CANDY_LLVM_MANIFEST="+manifest)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bundle-llvm.ps1 failed: %v\n%s", err, string(out))
	}
	if _, err := os.Stat(filepath.Join(outDir, "llvm", "bin", "clang.exe")); err != nil {
		t.Fatalf("missing bundled clang.exe: %v", err)
	}
}

func mkFakeLLVM(t *testing.T) string {
	t.Helper()
	llvm := t.TempDir()
	bin := filepath.Join(llvm, "bin")
	lib := filepath.Join(llvm, "lib")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(lib, 0o755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		_ = os.WriteFile(filepath.Join(bin, "clang.exe"), []byte("x"), 0o755)
		_ = os.WriteFile(filepath.Join(lib, "x.dll"), []byte("x"), 0o644)
	} else {
		_ = os.WriteFile(filepath.Join(bin, "clang"), []byte("x"), 0o755)
		_ = os.WriteFile(filepath.Join(lib, "libx.so"), []byte("x"), 0o644)
	}
	return llvm
}
