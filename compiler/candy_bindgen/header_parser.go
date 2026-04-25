package candy_bindgen

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var cFnDecl = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_\s\*]*?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*\((.*)\)\s*;\s*$`)
var cFnDef = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_\s\*]*?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*\((.*)\)\s*\{\s*$`)
var cTypedefDecl = regexp.MustCompile(`^\s*typedef\s+(.+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*;\s*$`)
var cDefineDecl = regexp.MustCompile(`^\s*#define\s+([A-Za-z_][A-Za-z0-9_]*)\s+(.+)\s*$`)
var cEnumStart = regexp.MustCompile(`^\s*enum\s+([A-Za-z_][A-Za-z0-9_]*)?\s*\{\s*$`)
var cEnumValue = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)(\s*=\s*([^,]+))?\s*,?\s*$`)
var cEnumEnd = regexp.MustCompile(`^\s*\}\s*([A-Za-z_][A-Za-z0-9_]*)?\s*;\s*$`)

func ParseHeaderFunctions(path string) ([]ExternBinding, []string, error) {
	api, warns, err := ParseHeaders([]string{path}, ParseOptions{})
	if err != nil {
		return nil, nil, err
	}
	out := make([]ExternBinding, 0, len(api.Functions))
	for _, fn := range api.Functions {
		ex := ExternBinding{
			Name:       fn.Name,
			Symbol:     fn.Symbol,
			ReturnType: fn.ReturnType,
			Params:     make([]ExternParam, 0, len(fn.Params)),
		}
		for _, p := range fn.Params {
			ex.Params = append(ex.Params, ExternParam{Name: p.Name, Type: p.Type})
		}
		out = append(out, ex)
	}
	return out, warns, nil
}

func ParseHeaders(headerFiles []string, opts ParseOptions) (*API, []string, error) {
	api := &API{}
	var warnings []string
	for _, path := range headerFiles {
		if err := parseSingleHeader(path, api, &warnings, opts); err != nil {
			return nil, warnings, err
		}
	}
	dedupeAPI(api)
	return api, warnings, nil
}

func parseSingleHeader(path string, api *API, warnings *[]string, opts ParseOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var inEnum bool
	var enumName string
	var enumVals []EnumValue
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}
		if inEnum {
			if m := cEnumEnd.FindStringSubmatch(line); len(m) == 2 {
				n := enumName
				if strings.TrimSpace(m[1]) != "" {
					n = strings.TrimSpace(m[1])
				}
				api.Enums = append(api.Enums, Enum{Name: n, Values: enumVals})
				inEnum = false
				enumName = ""
				enumVals = nil
				continue
			}
			if m := cEnumValue.FindStringSubmatch(line); len(m) >= 2 {
				val := EnumValue{Name: m[1]}
				if len(m) >= 4 && strings.TrimSpace(m[3]) != "" {
					val.Value = strings.TrimSpace(m[3])
				}
				enumVals = append(enumVals, val)
			}
			continue
		}
		if m := cEnumStart.FindStringSubmatch(line); len(m) == 2 {
			inEnum = true
			enumName = strings.TrimSpace(m[1])
			enumVals = nil
			continue
		}
		if m := cDefineDecl.FindStringSubmatch(line); len(m) == 3 {
			api.Constants = append(api.Constants, Constant{Name: strings.TrimSpace(m[1]), Value: strings.TrimSpace(m[2]), Type: "macro"})
			continue
		}
		if m := cTypedefDecl.FindStringSubmatch(line); len(m) == 3 {
			ct := strings.TrimSpace(m[1])
			name := strings.TrimSpace(m[2])
			api.Types = append(api.Types, TypeDef{Name: name, CType: ct, CandyType: TypeToCandy(ct)})
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		m := cFnDecl.FindStringSubmatch(line)
		if len(m) != 4 {
			m = cFnDef.FindStringSubmatch(line)
		}
		if len(m) != 4 {
			continue
		}
		retTy := strings.TrimSpace(m[1])
		name := strings.TrimSpace(m[2])
		args := strings.TrimSpace(m[3])
		fn := Function{
			Name:       name,
			Symbol:     name,
			ReturnType: retTy,
		}
		isVariadic := strings.Contains(args, "...")
		if isVariadic {
			if opts.UnsafeABI {
				*warnings = append(*warnings, "including variadic function "+name+" due to unsafe_abi mode")
			} else {
				*warnings = append(*warnings, "skipping variadic function "+name)
				continue
			}
		}
		cleanArgs := stripVariadic(args)
		if args != "" && args != "void" {
			parts := splitArgs(cleanArgs)
			for i, p := range parts {
				typ, pname := splitParam(p, i+1)
				fn.Params = append(fn.Params, Parameter{Name: pname, Type: typ})
			}
		}
		fn.Variadic = isVariadic
		ex := ExternBinding{Name: fn.Name, Symbol: fn.Symbol, ReturnType: fn.ReturnType, Variadic: fn.Variadic}
		for _, p := range fn.Params {
			ex.Params = append(ex.Params, ExternParam{Name: p.Name, Type: p.Type})
		}
		if why := IsUnsafeABI(ex); why != "" && !opts.UnsafeABI {
			*warnings = append(*warnings, "skipping "+name+": "+why)
			continue
		} else if why != "" && opts.UnsafeABI {
			*warnings = append(*warnings, "including "+name+" in unsafe_abi mode: "+why)
		}
		api.Functions = append(api.Functions, fn)
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func dedupeAPI(api *API) {
	seenFn := map[string]struct{}{}
	outFn := make([]Function, 0, len(api.Functions))
	for _, fn := range api.Functions {
		k := strings.ToLower(fn.Name)
		if _, ok := seenFn[k]; ok {
			continue
		}
		seenFn[k] = struct{}{}
		outFn = append(outFn, fn)
	}
	api.Functions = outFn
}

func splitArgs(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		it = strings.TrimSpace(it)
		if it != "" {
			out = append(out, it)
		}
	}
	return out
}

func stripVariadic(s string) string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "...") {
			continue
		}
		out = append(out, p)
	}
	return strings.Join(out, ", ")
}

func splitParam(param string, idx int) (string, string) {
	param = strings.Join(strings.Fields(param), " ")
	if strings.Contains(param, "(*)") {
		return param, "arg" + itoa(idx)
	}
	parts := strings.Fields(param)
	if len(parts) == 1 {
		return parts[0], "arg" + itoa(idx)
	}
	name := parts[len(parts)-1]
	typ := strings.TrimSpace(strings.TrimSuffix(param, name))
	if strings.HasPrefix(name, "*") {
		name = strings.TrimPrefix(name, "*")
		typ = strings.TrimSpace(typ + " *")
	}
	return strings.TrimSpace(typ), sanitizeName(name, idx)
}

func sanitizeName(s string, idx int) string {
	s = strings.TrimSpace(strings.Trim(s, "*"))
	if s == "" {
		return "arg" + itoa(idx)
	}
	ok := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(s)
	if !ok {
		return "arg" + itoa(idx)
	}
	return s
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + (v % 10))
		v /= 10
	}
	return string(b[i:])
}
