package candy_llvm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	envClangPath   = "CANDY_CLANG"
	envOptPath     = "CANDY_OPT"
	envBundleRoot  = "CANDY_LLVM_BUNDLE"
	bundledLLVMSeg = "llvm"
	bundledToolSeg = "toolchain"
)

// ResolveClangPath returns the preferred clang executable path.
//
// Resolution order:
// 1) CANDY_CLANG override
// 2) bundled LLVM next to the candy executable
// 3) clang from PATH
func ResolveClangPath() (string, error) {
	if p := os.Getenv(envClangPath); p != "" {
		if !isFile(p) {
			return "", fmt.Errorf("%s is set but not a valid file: %s", envClangPath, p)
		}
		return p, nil
	}
	if p, ok := bundledClangPath(); ok {
		return p, nil
	}
	return exec.LookPath(clangBinaryName())
}

// ResolveOptPath returns the preferred llvm-opt executable path.
//
// Resolution order:
// 1) CANDY_OPT override
// 2) bundled LLVM next to the candy executable
// 3) opt from PATH
func ResolveOptPath() (string, error) {
	if p := os.Getenv(envOptPath); p != "" {
		if !isFile(p) {
			return "", fmt.Errorf("%s is set but not a valid file: %s", envOptPath, p)
		}
		return p, nil
	}
	if p, ok := bundledOptPath(); ok {
		return p, nil
	}
	return exec.LookPath(optBinaryName())
}

func bundledClangPath() (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	base := filepath.Dir(exe)
	for _, p := range bundledBinaryCandidates(base, os.Getenv(envBundleRoot), clangBinaryName()) {
		if isFile(p) {
			return p, true
		}
	}
	return "", false
}

func bundledOptPath() (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	base := filepath.Dir(exe)
	for _, p := range bundledBinaryCandidates(base, os.Getenv(envBundleRoot), optBinaryName()) {
		if isFile(p) {
			return p, true
		}
	}
	return "", false
}

func bundledBinaryCandidates(exeDir, bundleRoot, binary string) []string {
	candidates := make([]string, 0, 10)
	add := func(p string) {
		if p == "" {
			return
		}
		candidates = append(candidates, filepath.Clean(p))
	}
	if bundleRoot != "" {
		// Support both "root/bin" and release-root layouts ("root/toolchain/bin", "root/llvm/bin").
		add(filepath.Join(bundleRoot, "bin", binary))
		add(filepath.Join(bundleRoot, bundledToolSeg, "bin", binary))
		add(filepath.Join(bundleRoot, bundledLLVMSeg, "bin", binary))
	}
	// Executable-local portable layout (old and new):
	// <root>/candy + <root>/llvm/bin OR <root>/toolchain/bin
	add(filepath.Join(exeDir, bundledToolSeg, "bin", binary))
	add(filepath.Join(exeDir, bundledLLVMSeg, "bin", binary))
	// Bin-folder layout:
	// <root>/bin/candy + <root>/toolchain/bin OR <root>/llvm/bin
	parent := filepath.Dir(exeDir)
	add(filepath.Join(parent, bundledToolSeg, "bin", binary))
	add(filepath.Join(parent, bundledLLVMSeg, "bin", binary))
	return candidates
}

func clangBinaryName() string {
	if runtime.GOOS == "windows" {
		return "clang.exe"
	}
	return "clang"
}

func optBinaryName() string {
	if runtime.GOOS == "windows" {
		return "opt.exe"
	}
	return "opt"
}

func isFile(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !st.IsDir()
}
