package candy_load

import (
	"path/filepath"
	"runtime"
	"strings"

	"candy/candy_bindgen"
)

type BuildContext struct {
	GlueSources []string
	IncludeDirs []string
	CFlags      []string
	LibDirs     []string
	Libs        []string
	LDFlags     []string
	Static      bool
	StaticLibs  []string
}

func (bc *BuildContext) mergeManifest(m *candy_bindgen.Manifest) {
	if bc == nil || m == nil {
		return
	}
	c, l := candy_bindgen.LinkContextForOS(m, runtime.GOOS)
	bc.GlueSources = appendUnique(bc.GlueSources, c.GlueSources...)
	bc.IncludeDirs = appendUnique(bc.IncludeDirs, c.IncludeDirs...)
	bc.CFlags = appendUnique(bc.CFlags, c.CFlags...)
	bc.LibDirs = appendUnique(bc.LibDirs, l.LibDirs...)
	bc.Libs = appendUnique(bc.Libs, l.Libs...)
	bc.LDFlags = appendUnique(bc.LDFlags, l.LDFlags...)
	bc.StaticLibs = appendUnique(bc.StaticLibs, l.StaticLibs...)
	bc.Static = bc.Static || l.Static
}

func appendUnique(dst []string, items ...string) []string {
	seen := map[string]struct{}{}
	for _, v := range dst {
		seen[normalize(v)] = struct{}{}
	}
	for _, it := range items {
		if strings.TrimSpace(it) == "" {
			continue
		}
		k := normalize(it)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		dst = append(dst, it)
	}
	return dst
}

func normalize(v string) string {
	return strings.ToLower(filepath.Clean(v))
}
