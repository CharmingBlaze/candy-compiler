package candy_llvm

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMajorVersion(t *testing.T) {
	if got := majorVersion("clang version 18.1.8"); got != 18 {
		t.Fatalf("expected 18, got %d", got)
	}
	if got := majorVersion("opt version 17.0.6"); got != 17 {
		t.Fatalf("expected 17, got %d", got)
	}
}

func TestProbeToolchain_WithOverrides(t *testing.T) {
	dir := t.TempDir()
	clang := fakeVersionTool(t, dir, "clang", "clang version 17.0.6")
	opt := fakeVersionTool(t, dir, "opt", "opt version 17.0.6")
	t.Setenv(envClangPath, clang)
	t.Setenv(envOptPath, opt)

	probe := ProbeToolchain()
	if !probe.Clang.Found {
		t.Fatalf("expected clang found, err=%s", probe.Clang.Error)
	}
	if !probe.Opt.Found {
		t.Fatalf("expected opt found, err=%s", probe.Opt.Error)
	}
	if probe.Clang.Version == "" || probe.Opt.Version == "" {
		t.Fatalf("expected versions to be detected: clang=%q opt=%q", probe.Clang.Version, probe.Opt.Version)
	}
}

func TestOSSpecificToolchainSuggestion(t *testing.T) {
	s := osSpecificToolchainSuggestion()
	if strings.TrimSpace(s) == "" {
		t.Fatal("expected non-empty suggestion")
	}
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(strings.ToLower(s), "windows") {
			t.Fatalf("expected windows-specific guidance, got: %s", s)
		}
	case "darwin":
		if !strings.Contains(strings.ToLower(s), "macos") {
			t.Fatalf("expected macos-specific guidance, got: %s", s)
		}
	default:
		if !strings.Contains(strings.ToLower(s), "linux") {
			t.Fatalf("expected linux guidance fallback, got: %s", s)
		}
	}
}

func fakeVersionTool(t *testing.T, dir, name, version string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		p := filepath.Join(dir, name+".cmd")
		body := "@echo off\r\necho " + version + "\r\nexit /b 0\r\n"
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	p := filepath.Join(dir, name+".sh")
	body := "#!/usr/bin/env bash\nset -euo pipefail\necho \"" + version + "\"\n"
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

