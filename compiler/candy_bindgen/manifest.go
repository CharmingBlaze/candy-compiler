package candy_bindgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Manifest struct {
	Library     string             `json:"library"`
	Namespace   string             `json:"namespace,omitempty"`
	UnsafeABI   bool               `json:"unsafe_abi,omitempty"`
	Headers     []string           `json:"headers,omitempty"`
	GeneratedAt string             `json:"generated_at,omitempty"`
	Externs     []ExternBinding    `json:"externs"`
	Types       []TypeDef          `json:"types,omitempty"`
	Enums       []Enum             `json:"enums,omitempty"`
	Constants   []Constant         `json:"constants,omitempty"`
	Compile     CompileOptions     `json:"compile,omitempty"`
	Link        LinkOptions        `json:"link,omitempty"`
	Platforms   map[string]Overlay `json:"platforms,omitempty"`
}

type ExternBinding struct {
	Name       string        `json:"name"`
	Symbol     string        `json:"symbol,omitempty"`
	ReturnType string        `json:"return_type"`
	Params     []ExternParam `json:"params,omitempty"`
	Variadic   bool          `json:"variadic,omitempty"`
}

type ExternParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CompileOptions struct {
	GlueSources []string `json:"glue_sources,omitempty"`
	IncludeDirs []string `json:"include_dirs,omitempty"`
	CFlags      []string `json:"cflags,omitempty"`
}

type LinkOptions struct {
	LibDirs    []string `json:"lib_dirs,omitempty"`
	Libs       []string `json:"libs,omitempty"`
	LDFlags    []string `json:"ldflags,omitempty"`
	Static     bool     `json:"static,omitempty"`
	StaticLibs []string `json:"static_libs,omitempty"`
}

type Overlay struct {
	Compile *CompileOptions `json:"compile,omitempty"`
	Link    *LinkOptions    `json:"link,omitempty"`
}

func ParseManifestFile(path string) (*Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("candylib parse: %w", err)
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	m.ResolveRelative(path)
	return &m, nil
}

func (m *Manifest) Validate() error {
	if strings.TrimSpace(m.Library) == "" {
		return fmt.Errorf("candylib validate: missing library")
	}
	if len(m.Externs) == 0 {
		return fmt.Errorf("candylib validate: externs is empty")
	}
	seen := map[string]struct{}{}
	for i := range m.Externs {
		ex := m.Externs[i]
		if strings.TrimSpace(ex.Name) == "" {
			return fmt.Errorf("candylib validate: extern[%d] missing name", i)
		}
		if strings.TrimSpace(ex.ReturnType) == "" {
			return fmt.Errorf("candylib validate: extern[%d] missing return_type", i)
		}
		if !isValidIdent(ex.Name) {
			return fmt.Errorf("candylib validate: extern[%d] invalid name %q", i, ex.Name)
		}
		if _, ok := seen[strings.ToLower(ex.Name)]; ok {
			return fmt.Errorf("candylib validate: duplicate extern name %q", ex.Name)
		}
		seen[strings.ToLower(ex.Name)] = struct{}{}
		for j := range ex.Params {
			p := ex.Params[j]
			if strings.TrimSpace(p.Name) == "" {
				return fmt.Errorf("candylib validate: extern[%d] param[%d] missing name", i, j)
			}
			if strings.TrimSpace(p.Type) == "" {
				return fmt.Errorf("candylib validate: extern[%d] param[%d] missing type", i, j)
			}
			if !isValidIdent(p.Name) {
				return fmt.Errorf("candylib validate: extern[%d] param[%d] invalid name %q", i, j, p.Name)
			}
		}
	}
	return nil
}

func (m *Manifest) ResolveRelative(manifestPath string) {
	base := filepath.Dir(manifestPath)
	rewrite := func(items []string) []string {
		out := make([]string, 0, len(items))
		for _, it := range items {
			if it == "" || filepath.IsAbs(it) {
				out = append(out, it)
				continue
			}
			out = append(out, filepath.Clean(filepath.Join(base, it)))
		}
		return out
	}
	m.Compile.GlueSources = rewrite(m.Compile.GlueSources)
	m.Compile.IncludeDirs = rewrite(m.Compile.IncludeDirs)
	m.Link.LibDirs = rewrite(m.Link.LibDirs)
	for k, ov := range m.Platforms {
		if ov.Compile != nil {
			ov.Compile.GlueSources = rewrite(ov.Compile.GlueSources)
			ov.Compile.IncludeDirs = rewrite(ov.Compile.IncludeDirs)
		}
		if ov.Link != nil {
			ov.Link.LibDirs = rewrite(ov.Link.LibDirs)
		}
		m.Platforms[k] = ov
	}
}

func isValidIdent(name string) bool {
	return regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(name)
}

func TypeToCandy(cType string) string {
	t := strings.TrimSpace(strings.ToLower(cType))
	t = strings.ReplaceAll(t, "const ", "")
	t = strings.ReplaceAll(t, "volatile ", "")
	t = strings.Join(strings.Fields(t), " ")
	switch t {
	case "void":
		return "void"
	case "bool", "_bool":
		return "bool"
	case "float", "double":
		return "float"
	case "char*", "char *", "const char*", "const char *":
		return "string"
	case "int", "short", "long", "long long",
		"unsigned", "unsigned int", "unsigned short", "unsigned long", "unsigned long long",
		"size_t", "ssize_t", "int8_t", "int16_t", "int32_t", "int64_t",
		"uint8_t", "uint16_t", "uint32_t", "uint64_t":
		return "int"
	}
	if strings.HasSuffix(t, "*") {
		return "int"
	}
	return "int"
}

func IsUnsafeABI(ex ExternBinding) string {
	if ex.Variadic {
		return "variadic externs are not supported"
	}
	for _, p := range ex.Params {
		pt := strings.TrimSpace(strings.ToLower(p.Type))
		if strings.Contains(pt, "(*)") {
			return "function pointer params are not supported"
		}
	}
	rt := strings.TrimSpace(strings.ToLower(ex.ReturnType))
	if strings.Contains(rt, "(*)") {
		return "function pointer returns are not supported"
	}
	return ""
}

func LinkContextForOS(m *Manifest, goos string) (CompileOptions, LinkOptions) {
	c := m.Compile
	l := m.Link
	if ov, ok := m.Platforms[strings.ToLower(goos)]; ok {
		if ov.Compile != nil {
			c.GlueSources = append(c.GlueSources, ov.Compile.GlueSources...)
			c.IncludeDirs = append(c.IncludeDirs, ov.Compile.IncludeDirs...)
			c.CFlags = append(c.CFlags, ov.Compile.CFlags...)
		}
		if ov.Link != nil {
			l.LibDirs = append(l.LibDirs, ov.Link.LibDirs...)
			l.Libs = append(l.Libs, ov.Link.Libs...)
			l.LDFlags = append(l.LDFlags, ov.Link.LDFlags...)
			l.StaticLibs = append(l.StaticLibs, ov.Link.StaticLibs...)
			l.Static = l.Static || ov.Link.Static
		}
	}
	return c, l
}
