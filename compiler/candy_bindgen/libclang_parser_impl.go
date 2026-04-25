//go:build libclang

package candy_bindgen

import (
	"fmt"
	"strings"

	clang "github.com/go-clang/clang-v14/clang"
)

func parseHeadersLibclangImpl(headerFiles []string, opts ParseOptions) (*API, []string, error) {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()
	args := make([]string, 0, len(opts.IncludeDirs)+len(opts.Defines))
	lang := strings.ToLower(strings.TrimSpace(opts.Language))
	if lang == "c++" || lang == "cpp" || lang == "cxx" {
		args = append(args, "-x", "c++")
	}
	for _, inc := range opts.IncludeDirs {
		args = append(args, "-I"+inc)
	}
	for _, d := range opts.Defines {
		args = append(args, "-D"+d)
	}
	api := &API{}
	var warnings []string
	for _, hf := range headerFiles {
		tu := idx.ParseTranslationUnit(hf, args, nil, clang.TranslationUnit_None)
		if tu.IsNil() {
			return nil, warnings, fmt.Errorf("libclang: failed to parse %s", hf)
		}
		cursor := tu.TranslationUnitCursor()
		cursor.Visit(func(c, parent clang.Cursor) clang.ChildVisitResult {
			if c.IsNull() {
				return clang.ChildVisit_Continue
			}
			loc := c.Location()
			file, _, _, _ := loc.FileLocation()
			if !file.IsNil() && !strings.EqualFold(file.Name(), hf) {
				return clang.ChildVisit_Continue
			}
			switch c.Kind() {
			case clang.Cursor_FunctionDecl:
				fn := Function{
					Name:       c.Spelling(),
					Symbol:     c.Spelling(),
					ReturnType: c.ResultType().Spelling(),
					Comment:    c.BriefComment(),
				}
				num := c.NumArguments()
				for i := uint32(0); i < uint32(num); i++ {
					arg := c.Argument(i)
					name := arg.Spelling()
					if strings.TrimSpace(name) == "" {
						name = fmt.Sprintf("arg%d", i+1)
					}
					fn.Params = append(fn.Params, Parameter{Name: name, Type: arg.Type().Spelling()})
				}
				ex := ExternBinding{Name: fn.Name, Symbol: fn.Symbol, ReturnType: fn.ReturnType}
				for _, p := range fn.Params {
					ex.Params = append(ex.Params, ExternParam{Name: p.Name, Type: p.Type})
				}
				if why := IsUnsafeABI(ex); why == "" || opts.UnsafeABI {
					if why != "" && opts.UnsafeABI {
						warnings = append(warnings, "including "+fn.Name+" in unsafe_abi mode: "+why)
					}
					api.Functions = append(api.Functions, fn)
				} else {
					warnings = append(warnings, "skipping "+fn.Name+": "+why)
				}
			case clang.Cursor_TypedefDecl:
				ct := c.TypedefDeclUnderlyingType().Spelling()
				api.Types = append(api.Types, TypeDef{Name: c.Spelling(), CType: ct, CandyType: TypeToCandy(ct)})
			case clang.Cursor_EnumDecl:
				en := Enum{Name: c.Spelling()}
				c.Visit(func(ch, _ clang.Cursor) clang.ChildVisitResult {
					if ch.Kind() == clang.Cursor_EnumConstantDecl {
						en.Values = append(en.Values, EnumValue{Name: ch.Spelling(), Value: fmt.Sprintf("%d", ch.EnumConstantDeclValue())})
					}
					return clang.ChildVisit_Continue
				})
				api.Enums = append(api.Enums, en)
			}
			return clang.ChildVisit_Continue
		})
		tu.Dispose()
	}
	dedupeAPI(api)
	return api, warnings, nil
}
