package candy_bindgen

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ExpandInputFiles expands input paths/patterns/directories into concrete C/C++ files.
// Supported:
// - literal files
// - wildcard patterns (e.g. "*.h", "src/**/*.c" where shell expanded patterns are passed)
// - directories (recursively scanned)
func ExpandInputFiles(inputs []string) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(p string) {
		p = filepath.Clean(strings.TrimSpace(p))
		if p == "" {
			return
		}
		if abs, err := filepath.Abs(p); err == nil {
			p = abs
		}
		if _, ok := seen[strings.ToLower(p)]; ok {
			return
		}
		seen[strings.ToLower(p)] = struct{}{}
		out = append(out, p)
	}
	for _, raw := range inputs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if hasWildcard(raw) {
			matches, _ := filepath.Glob(raw)
			for _, m := range matches {
				fi, err := os.Stat(m)
				if err != nil {
					continue
				}
				if fi.IsDir() {
					walkDirCFiles(m, add)
					continue
				}
				if isCInputFile(m) {
					add(m)
				}
			}
			continue
		}
		fi, err := os.Stat(raw)
		if err != nil {
			continue
		}
		if fi.IsDir() {
			walkDirCFiles(raw, add)
			continue
		}
		if isCInputFile(raw) {
			add(raw)
		}
	}
	slices.Sort(out)
	return out
}

func hasWildcard(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

func walkDirCFiles(root string, add func(string)) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if isCInputFile(path) {
			add(path)
		}
		return nil
	})
}

func isCInputFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".h", ".hh", ".hpp", ".hxx", ".c", ".cc", ".cpp", ".cxx":
		return true
	default:
		return false
	}
}

func IsSourceFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".c", ".cc", ".cpp", ".cxx":
		return true
	default:
		return false
	}
}

func IsHeaderFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".h", ".hh", ".hpp", ".hxx":
		return true
	default:
		return false
	}
}

// DiscoverLibraryFiles scans library roots and returns header/source/include file lists.
func DiscoverLibraryFiles(roots []string, lang string) (headers []string, sources []string, includeDirs []string, err error) {
	roots = ExpandInputFiles(roots)
	if len(roots) == 0 {
		return nil, nil, nil, fmt.Errorf("no library roots/files found")
	}
	seenH := map[string]struct{}{}
	seenS := map[string]struct{}{}
	seenI := map[string]struct{}{}
	addHeader := func(p string) {
		k := strings.ToLower(filepath.Clean(p))
		if _, ok := seenH[k]; ok {
			return
		}
		seenH[k] = struct{}{}
		headers = append(headers, p)
		dir := filepath.Dir(p)
		dk := strings.ToLower(filepath.Clean(dir))
		if _, ok := seenI[dk]; !ok {
			seenI[dk] = struct{}{}
			includeDirs = append(includeDirs, dir)
		}
	}
	addSource := func(p string) {
		k := strings.ToLower(filepath.Clean(p))
		if _, ok := seenS[k]; ok {
			return
		}
		seenS[k] = struct{}{}
		sources = append(sources, p)
	}
	isCXX := strings.EqualFold(lang, "c++") || strings.EqualFold(lang, "cpp") || strings.EqualFold(lang, "cxx")
	for _, r := range roots {
		fi, statErr := os.Stat(r)
		if statErr != nil {
			continue
		}
		if fi.IsDir() {
			_ = filepath.WalkDir(r, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil || d.IsDir() {
					return nil
				}
				if IsHeaderFile(path) {
					addHeader(path)
					return nil
				}
				if IsSourceFile(path) {
					ext := strings.ToLower(filepath.Ext(path))
					if isCXX {
						if ext == ".cpp" || ext == ".cc" || ext == ".cxx" {
							addSource(path)
						}
					} else if ext == ".c" {
						addSource(path)
					}
				}
				return nil
			})
			continue
		}
		if IsHeaderFile(r) {
			addHeader(r)
			continue
		}
		if IsSourceFile(r) {
			addSource(r)
		}
	}
	slices.Sort(headers)
	slices.Sort(sources)
	slices.Sort(includeDirs)
	return headers, sources, includeDirs, nil
}
