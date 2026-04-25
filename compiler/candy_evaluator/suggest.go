package candy_evaluator

import (
	"fmt"
	"strings"
)

// levenshteinDistance returns an edit distance between two strings (small inputs only).
func levenshteinDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// trim run cost for very long names
	if la > 64 || lb > 64 {
		if strings.EqualFold(a, b) {
			return 0
		}
		return 999
	}
	rows := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		rows[j] = j
	}
	for i := 1; i <= la; i++ {
		prev := rows[0]
		rows[0] = i
		for j := 1; j <= lb; j++ {
			cur := rows[j]
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			ins := rows[j-1] + 1
			del := rows[j] + 1
			x := prev + cost
			rows[j] = min3(ins, del, x)
			prev = cur
		}
	}
	return rows[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// bestFuzzyName picks the closest name from options if within maxDist.
func bestFuzzyName(want string, options []string) (string, int) {
	want = strings.TrimSpace(want)
	if want == "" || len(options) == 0 {
		return "", 999
	}
	best := ""
	bd := 999
	for _, o := range options {
		if o == want {
			return o, 0
		}
		if strings.EqualFold(o, want) {
			return o, 0
		}
		d := levenshteinDistance(strings.ToLower(want), strings.ToLower(o))
		if d < bd {
			bd, best = d, o
		}
	}
	max := 2
	if len(want) > 6 {
		max = 3
	}
	if len(want) > 10 {
		max = 4
	}
	if bd > max {
		return "", bd
	}
	return best, bd
}

func withDidYouMean(what, missing string, cands []string) error {
	if s, _ := bestFuzzyName(missing, cands); s != "" && s != missing {
		return &RuntimeError{Msg: fmt.Sprintf("%s has no field %q. Did you mean %q?", what, missing, s)}
	}
	return &RuntimeError{Msg: fmt.Sprintf("%s has no field %q", what, missing)}
}

func withUndefinedVar(name string, cands []string) error {
	if s, _ := bestFuzzyName(name, cands); s != "" {
		return &RuntimeError{Msg: fmt.Sprintf("undefined: %q. Did you mean %q?", name, s)}
	}
	return &RuntimeError{Msg: "undefined: " + name}
}
