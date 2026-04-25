package candy_llvm

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type BuildProfile int

const (
	BuildDebug BuildProfile = iota
	BuildDevRelease
	BuildShipping
)

func (p BuildProfile) String() string {
	switch p {
	case BuildDebug:
		return "debug"
	case BuildShipping:
		return "shipping"
	default:
		return "dev-release"
	}
}

func GetPassPipeline(profile BuildProfile) []string {
	switch profile {
	case BuildDebug:
		return []string{
			"verify",
		}
	case BuildShipping:
		return []string{
			"verify",
			"mem2reg",
			"sroa",
			"instcombine",
			"reassociate",
			"gvn",
			"simplifycfg",
			"early-cse",
			"jump-threading",
			"licm",
			"loop-simplify",
			"indvars",
			"loop-idiom",
			"loop-rotate",
			"loop-unroll",
			"slp-vectorizer",
			"loop-vectorize",
			"instcombine",
			"simplifycfg",
			"globaldce",
			"constmerge",
		}
	default:
		return []string{
			"mem2reg",
			"instcombine",
			"reassociate",
			"gvn",
			"simplifycfg",
			"licm",
			"loop-simplify",
			"inline",
		}
	}
}

func passSpec(profile BuildProfile) string {
	return strings.Join(GetPassPipeline(profile), ",")
}

// OptimizeIR applies an LLVM opt pass pipeline to textual IR.
// If profile is debug, it returns the input IR unchanged.
func OptimizeIR(ir string, profile BuildProfile) (string, error) {
	if profile == BuildDebug {
		return ir, nil
	}
	optPath, err := ResolveOptPath()
	if err != nil {
		return ir, nil
	}

	cmd := exec.Command(optPath, "-S", "-passes="+passSpec(profile))
	cmd.Stdin = strings.NewReader(ir)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if runErr := cmd.Run(); runErr != nil {
		return ir, fmt.Errorf("opt failed: %w: %s", runErr, strings.TrimSpace(errOut.String()))
	}
	if out.Len() == 0 {
		return ir, fmt.Errorf("opt produced empty output")
	}
	return out.String(), nil
}
