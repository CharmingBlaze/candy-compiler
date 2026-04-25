package candy_llvm

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type ToolStatus struct {
	Name    string
	Path    string
	Version string
	Found   bool
	Error   string
}

type ProbeResult struct {
	Clang        ToolStatus
	Opt          ToolStatus
	Compatible   bool
	Problems     []string
	Suggestions  []string
	SearchPolicy []string
}

func ProbeToolchain() ProbeResult {
	res := ProbeResult{
		Compatible: true,
		SearchPolicy: []string{
			"CANDY_* override",
			"bundled llvm/bin",
			"system PATH",
		},
	}

	res.Clang = probeSingleTool("clang", ResolveClangPath)
	res.Opt = probeSingleTool("opt", ResolveOptPath)

	if !res.Clang.Found {
		res.Compatible = false
		res.Problems = append(res.Problems, "clang not found")
	}
	if !res.Opt.Found {
		// opt can be optional for development flows; keep compatibility true
		// unless clang is already missing. Build will run without opt.
		res.Problems = append(res.Problems, "opt not found (IR optimization disabled)")
	}

	if res.Clang.Found && res.Opt.Found {
		cm := majorVersion(res.Clang.Version)
		om := majorVersion(res.Opt.Version)
		if cm > 0 && om > 0 && cm != om {
			res.Problems = append(res.Problems, fmt.Sprintf("clang/opt major mismatch: clang=%d opt=%d", cm, om))
		}
	}

	if !res.Clang.Found {
		res.Suggestions = append(res.Suggestions,
			"Set CANDY_CLANG to an absolute clang path, or install/bundle LLVM.",
		)
		res.Suggestions = append(res.Suggestions, osSpecificToolchainSuggestion())
	}
	if !res.Opt.Found {
		res.Suggestions = append(res.Suggestions,
			"Set CANDY_OPT to an absolute opt path for optimized IR builds.",
		)
	}

	return res
}

func osSpecificToolchainSuggestion() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows: install LLVM and ensure clang.exe/opt.exe are reachable from PATH."
	case "darwin":
		return "macOS: install Xcode Command Line Tools or LLVM (brew) and verify clang/opt visibility."
	default:
		return "Linux: install clang/llvm via your package manager and verify PATH visibility."
	}
}

func probeSingleTool(name string, resolve func() (string, error)) ToolStatus {
	st := ToolStatus{Name: name}
	p, err := resolve()
	if err != nil {
		st.Error = err.Error()
		return st
	}
	st.Path = p
	st.Found = true
	st.Version = queryToolVersion(p)
	return st
}

func queryToolVersion(path string) string {
	cmd := exec.Command(path, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	return line
}

var reMajor = regexp.MustCompile(`\b([0-9]{1,3})\b`)

func majorVersion(versionLine string) int {
	m := reMajor.FindStringSubmatch(versionLine)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

