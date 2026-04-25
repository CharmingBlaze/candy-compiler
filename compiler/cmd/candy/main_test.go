package main

import (
	"bytes"
	"candy/candy_load"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunBuild_WithOverrideClang(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "app.candy")
	if err := os.WriteFile(srcPath, []byte("return 1+2;"), 0o644); err != nil {
		t.Fatal(err)
	}
	clang := mockClang(t, dir)
	t.Setenv("CANDY_CLANG", clang)
	t.Setenv("PATH", "")
	var out, errOut bytes.Buffer
	code := run([]string{"-build", srcPath}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Generated LLVM IR:") {
		t.Fatalf("expected IR output message, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "Generated Native Binary:") {
		t.Fatalf("expected native output message, got: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "app.ll")); err != nil {
		t.Fatalf("expected app.ll: %v", err)
	}
	if _, err := os.Stat(nativePath(filepath.Join(dir, "app"))); err != nil {
		t.Fatalf("expected native binary: %v", err)
	}
}

func TestRunBuild_StdinWithImport_Errors(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"-build"}, strings.NewReader("import \"a.candy\";\nreturn 0;"), &out, &errOut)
	if code != 1 {
		t.Fatalf("expected code 1, got %d out=%s err=%s", code, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "stdin") {
		t.Fatalf("expected stdin import error, got: %s", errOut.String())
	}
}

func TestRunBuild_NoClang_ShowsGuidance(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "app.candy")
	if err := os.WriteFile(srcPath, []byte("return 1;"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANDY_CLANG", "")
	t.Setenv("CANDY_LLVM_BUNDLE", filepath.Join(dir, "missing"))
	t.Setenv("PATH", "")
	var out, errOut bytes.Buffer
	code := run([]string{"-build", srcPath}, strings.NewReader(""), &out, &errOut)
	if code != 1 {
		t.Fatalf("expected code 1 for fail-fast path, got %d", code)
	}
	if !strings.Contains(out.String(), "Candy toolchain doctor") {
		t.Fatalf("expected no clang guidance, got: %s", out.String())
	}
}

func TestRunDoctor_NoClang_Fails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CANDY_CLANG", "")
	t.Setenv("CANDY_LLVM_BUNDLE", filepath.Join(dir, "missing"))
	t.Setenv("PATH", "")
	var out, errOut bytes.Buffer
	code := run([]string{"doctor"}, strings.NewReader(""), &out, &errOut)
	if code != 1 {
		t.Fatalf("expected code 1 when clang missing, got %d", code)
	}
	if !strings.Contains(out.String(), "Candy toolchain doctor") {
		t.Fatalf("expected doctor output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "Status: FAIL") {
		t.Fatalf("expected fail status in doctor output, got: %s", out.String())
	}
}

func TestRunDoctor_WithOverrideClang_Passes(t *testing.T) {
	dir := t.TempDir()
	clang := mockClang(t, dir)
	t.Setenv("CANDY_CLANG", clang)
	var out, errOut bytes.Buffer
	code := run([]string{"doctor"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0 when clang exists, got %d", code)
	}
	if !strings.Contains(out.String(), "Status: PASS") {
		t.Fatalf("expected pass status in doctor output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "Search order:") {
		t.Fatalf("expected search order block, got: %s", out.String())
	}
}

func mockClang(t *testing.T, dir string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		p := filepath.Join(dir, "clang.cmd")
		body := "@echo off\r\nset in=%1\r\nset out=%3\r\ncopy /Y \"%in%\" \"%out%\" >nul\r\nexit /b 0\r\n"
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	p := filepath.Join(dir, "clang.sh")
	body := "#!/usr/bin/env bash\nset -euo pipefail\ncp \"$1\" \"$3\"\n"
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func nativePath(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func TestRun_Repl(t *testing.T) {
	in := strings.NewReader("1+1\n:exit\n")
	var out, errOut bytes.Buffer
	code := run([]string{"-i"}, in, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d err=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "2") {
		t.Fatalf("expected computed 2 in output, got: %q", out.String())
	}
	if !strings.Contains(out.String(), "Bye.") {
		t.Fatalf("expected Bye., got: %q", out.String())
	}
}

func TestRunBuild_SubcommandAlias(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "app.candy")
	if err := os.WriteFile(srcPath, []byte("return 2+3;"), 0o644); err != nil {
		t.Fatal(err)
	}
	clang := mockClang(t, dir)
	t.Setenv("CANDY_CLANG", clang)
	t.Setenv("PATH", "")
	var out, errOut bytes.Buffer
	code := run([]string{"build", srcPath}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Generated Native Binary:") {
		t.Fatalf("expected native output message, got: %s", out.String())
	}
}

func TestRunBuild_SubcommandAlias_WithOutputFlagAfterFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "app.candy")
	if err := os.WriteFile(srcPath, []byte("return 7;"), 0o644); err != nil {
		t.Fatal(err)
	}
	clang := mockClang(t, dir)
	t.Setenv("CANDY_CLANG", clang)
	t.Setenv("PATH", "")
	outExe := filepath.Join(dir, "mygame")
	var out, errOut bytes.Buffer
	code := run([]string{"build", srcPath, "-o", outExe}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d stderr=%s", code, errOut.String())
	}
	if _, err := os.Stat(nativePath(outExe)); err != nil {
		t.Fatalf("expected output binary %s: %v", nativePath(outExe), err)
	}
}

func TestRunBuild_OptimizeFlagSelectsShippingProfile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "app.candy")
	if err := os.WriteFile(srcPath, []byte("return 7;"), 0o644); err != nil {
		t.Fatal(err)
	}
	clang := mockClang(t, dir)
	t.Setenv("CANDY_CLANG", clang)
	t.Setenv("PATH", "")
	var out, errOut bytes.Buffer
	code := run([]string{"-build", "--optimize", "--verbose", srcPath}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Build profile: shipping") {
		t.Fatalf("expected shipping build profile in verbose output, got: %s", out.String())
	}
}

func TestAppendClangBuildContext_StaticLinking(t *testing.T) {
	args := appendClangBuildContext([]string{"app.ll", "-o", "app"}, &candy_load.BuildContext{
		Libs:       []string{"m"},
		Static:     true,
		StaticLibs: []string{"box2d"},
	})
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-static") {
		t.Fatalf("expected -static in args: %s", joined)
	}
	if !strings.Contains(joined, "-Wl,-Bstatic -lbox2d -Wl,-Bdynamic") {
		t.Fatalf("expected static lib wrapping in args: %s", joined)
	}
}

func TestRun_HelpIncludesDoctorCommand(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"-h"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("expected help exit code 0, got %d", code)
	}
	if !strings.Contains(errOut.String(), "candy doctor") {
		t.Fatalf("expected help text to include doctor command, got: %s", errOut.String())
	}
}
