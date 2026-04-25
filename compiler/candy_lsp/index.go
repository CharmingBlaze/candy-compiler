package candy_lsp

import (
	"path/filepath"
	"sort"
	"strings"

	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_parser"
)

type symbolDef struct {
	Name      string
	Kind      string
	URI       string
	Container string
	Range     lspRange
}

func buildDocumentIndex(uri, src string) (symbols []symbolDef, imports []string) {
	p := candy_parser.New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if prog == nil {
		return nil, nil
	}
	for _, st := range prog.Statements {
		switch t := st.(type) {
		case *candy_ast.ImportStatement:
			imports = append(imports, t.Path)
		case *candy_ast.FunctionStatement:
			if t.Name != nil {
				symbols = append(symbols, makeSymbol(uri, t.Name.Value, "fun", t.Name.Token.Line, t.Name.Token.Col, len(t.Name.Value), ""))
			}
			container := ""
			if t.Name != nil {
				container = t.Name.Value
			}
			for _, p0 := range t.Parameters {
				if p0.Name != nil {
					symbols = append(symbols, makeSymbol(uri, p0.Name.Value, "param", p0.Name.Token.Line, p0.Name.Token.Col, len(p0.Name.Value), container))
				}
			}
		case *candy_ast.ValStatement:
			if t.Name != nil {
				symbols = append(symbols, makeSymbol(uri, t.Name.Value, "val", t.Name.Token.Line, t.Name.Token.Col, len(t.Name.Value), ""))
			}
		case *candy_ast.VarStatement:
			if t.Name != nil {
				symbols = append(symbols, makeSymbol(uri, t.Name.Value, "var", t.Name.Token.Line, t.Name.Token.Col, len(t.Name.Value), ""))
			}
		case *candy_ast.StructStatement:
			if t.Name != nil {
				symbols = append(symbols, makeSymbol(uri, t.Name.Value, "struct", t.Name.Token.Line, t.Name.Token.Col, len(t.Name.Value), ""))
			}
			for _, f := range t.Fields {
				if f.Name != nil {
					symbols = append(symbols, makeSymbol(uri, f.Name.Value, "field", f.Name.Token.Line, f.Name.Token.Col, len(f.Name.Value), ""))
				}
			}
		}
	}
	sort.Slice(symbols, func(i, j int) bool {
		if symbols[i].Range.Start.Line != symbols[j].Range.Start.Line {
			return symbols[i].Range.Start.Line < symbols[j].Range.Start.Line
		}
		if symbols[i].Range.Start.Character != symbols[j].Range.Start.Character {
			return symbols[i].Range.Start.Character < symbols[j].Range.Start.Character
		}
		return symbols[i].Name < symbols[j].Name
	})
	return symbols, imports
}

func makeSymbol(uri, name, kind string, line, col, n int, container string) symbolDef {
	if n <= 0 {
		n = len(name)
	}
	sl := max(line-1, 0)
	sc := max(col-1, 0)
	return symbolDef{
		Name:      strings.ToLower(name),
		Kind:      kind,
		URI:       uri,
		Container: container,
		Range: lspRange{
			Start: lspPosition{Line: sl, Character: sc},
			End:   lspPosition{Line: sl, Character: sc + n},
		},
	}
}

func resolveImportURI(_ string, fromPath, imp string) string {
	if strings.HasPrefix(strings.ToLower(imp), "std/") {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(imp), "file:///") {
		return imp
	}
	base := filepath.Dir(fromPath)
	if base == "." || base == "" {
		return ""
	}
	target := filepath.Clean(filepath.Join(base, imp))
	return "file:///" + filepath.ToSlash(target)
}
