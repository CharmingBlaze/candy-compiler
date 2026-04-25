package candy_bindgen

import (
	"fmt"
	"regexp"
	"strings"
)

// TransformAPI applies ignore/rename namespace transformations for generated extern names.
// It preserves each original symbol name for linking.
func TransformAPI(api *API, namespace string, stripPrefixes []string, ignorePatterns []string) ([]string, error) {
	if api == nil {
		return nil, nil
	}
	ns := sanitizeIdent(namespace)
	var warnings []string
	compiledIgnore := make([]*regexp.Regexp, 0, len(ignorePatterns))
	for _, raw := range ignorePatterns {
		p := strings.TrimSpace(raw)
		if p == "" {
			continue
		}
		// Allow simple wildcard patterns like "foo*" by translating to regex.
		if strings.Contains(p, "*") && !strings.ContainsAny(p, "[](){}+?|\\") {
			p = "^" + strings.ReplaceAll(regexp.QuoteMeta(p), "\\*", ".*") + "$"
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore pattern %q: %w", raw, err)
		}
		compiledIgnore = append(compiledIgnore, re)
	}

	seen := map[string]struct{}{}
	outFns := make([]Function, 0, len(api.Functions))
	for _, fn := range api.Functions {
		originalName := strings.TrimSpace(fn.Name)
		if originalName == "" {
			continue
		}
		origSymbol := strings.TrimSpace(fn.Symbol)
		if origSymbol == "" {
			origSymbol = originalName
		}
		if shouldIgnore(originalName, origSymbol, compiledIgnore) {
			warnings = append(warnings, "ignored function "+originalName+" by pattern")
			continue
		}
		name := stripFirstPrefix(originalName, stripPrefixes)
		name = sanitizeIdent(name)
		if name == "" {
			warnings = append(warnings, "skipping function "+originalName+": became empty after transforms")
			continue
		}
		if ns != "" {
			name = ns + "_" + name
		}
		lower := strings.ToLower(name)
		if _, ok := seen[lower]; ok {
			fallback := sanitizeIdent(originalName)
			if ns != "" {
				fallback = ns + "_" + fallback
			}
			fallbackLower := strings.ToLower(fallback)
			if _, exists := seen[fallbackLower]; exists {
				warnings = append(warnings, "skipping function "+originalName+": duplicate generated name "+name)
				continue
			}
			warnings = append(warnings, "name collision for "+originalName+", using fallback "+fallback)
			name = fallback
			lower = fallbackLower
		}
		seen[lower] = struct{}{}
		fn.Symbol = origSymbol
		fn.Name = name
		outFns = append(outFns, fn)
	}
	api.Functions = outFns
	return warnings, nil
}

func shouldIgnore(name, symbol string, patterns []*regexp.Regexp) bool {
	if len(patterns) == 0 {
		return false
	}
	for _, re := range patterns {
		if re.MatchString(name) || re.MatchString(symbol) {
			return true
		}
	}
	return false
}

func stripFirstPrefix(s string, prefixes []string) string {
	out := s
	for _, p := range prefixes {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(out, p) {
			out = strings.TrimPrefix(out, p)
			break
		}
	}
	return out
}

func sanitizeIdent(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range s {
		isAlpha := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		isNum := r >= '0' && r <= '9'
		if isAlpha || r == '_' || (i > 0 && isNum) {
			b.WriteRune(r)
			continue
		}
		if i == 0 && isNum {
			b.WriteRune('_')
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	out := b.String()
	for strings.Contains(out, "__") {
		out = strings.ReplaceAll(out, "__", "_")
	}
	return strings.Trim(out, "_")
}
