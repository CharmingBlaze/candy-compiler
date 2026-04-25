// Package candy_load expands import statements for native builds, mirroring candy_evaluator.evalImport.
package candy_load

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"candy/candy_ast"
	"candy/candy_bindgen"
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_report"
	"candy/candy_stdlib"
	"candy/candy_token"
)

// expander tracks loaded modules for cycle breaking and stdlib de-duplication.
type expander struct {
	imported map[string]struct{}
	ctx      *BuildContext
}

func (e *expander) keyStdlib(path string) string { return "stdlib::" + path }
func (e *expander) keyFile(full string) string   { return "file::" + full }

// ExpandProgramForBuild reads an entry .candy file, parses it, and returns a new program
// with all imports inlined in interpreter order. Disk imports use the entry file’s directory
// and nested imports use each loaded file’s directory, matching the evaluator.
func ExpandProgramForBuild(entryPath string) (*candy_ast.Program, error) {
	prog, _, err := ExpandProgramForBuildWithContext(entryPath)
	return prog, err
}

func ExpandProgramForBuildWithContext(entryPath string) (*candy_ast.Program, *BuildContext, error) {
	abs, err := filepath.Abs(entryPath)
	if err != nil {
		return nil, nil, fmt.Errorf("candy_load: %w", err)
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return nil, nil, err
	}
	return ExpandProgramForBuildFromSourceWithContext(string(b), filepath.Dir(abs), abs)
}

// ExpandProgramForBuildFromSource expands imports given source text and the directory
// to resolve relative file imports (the directory of the file that contained this text).
// originPath is used only for error messages.
func ExpandProgramForBuildFromSource(source, curDir, originPath string) (*candy_ast.Program, error) {
	prog, _, err := ExpandProgramForBuildFromSourceWithContext(source, curDir, originPath)
	return prog, err
}

func ExpandProgramForBuildFromSourceWithContext(source, curDir, originPath string) (*candy_ast.Program, *BuildContext, error) {
	l := candy_lexer.New(source)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, nil, formatParseErrors(originPath, p.Errors())
	}
	ctx := &BuildContext{}
	e := &expander{imported: make(map[string]struct{}), ctx: ctx}
	stmts, err := e.expandStmts(prog.Statements, curDir, originPath)
	if err != nil {
		return nil, nil, err
	}
	return &candy_ast.Program{Statements: stmts}, ctx, nil
}

func formatParseErrors(origin string, diags []candy_report.Diagnostic) error {
	originInfo := ""
	if origin != "" {
		originInfo = origin + ": "
	}
	var b strings.Builder
	for i, d := range diags {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(d.Message)
	}
	return fmt.Errorf("candy_load: %sparse error: %s", originInfo, b.String())
}

func (e *expander) expandStmts(stmts []candy_ast.Statement, curDir, errOrigin string) ([]candy_ast.Statement, error) {
	if len(stmts) == 0 {
		return nil, nil
	}
	var out []candy_ast.Statement
	for _, s := range stmts {
		if s == nil {
			continue
		}
		if imp, ok := s.(*candy_ast.ImportStatement); ok {
			nested, err := e.loadImport(imp.Path, curDir, errOrigin)
			if err != nil {
				return nil, err
			}
			out = append(out, nested...)
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

// loadImport returns a flattened, fully expanded list of statements for the import.
func (e *expander) loadImport(path, curDir, errOrigin string) ([]candy_ast.Statement, error) {
	// 1) Standard library
	if src, ok := candy_stdlib.Lookup(path); ok {
		return e.loadFromStdlib(path, src, errOrigin)
	}

	// 2) File on disk
	var full string
	if filepath.IsAbs(path) {
		full = filepath.Clean(path)
	} else {
		if curDir == "" {
			if wd, gerr := os.Getwd(); gerr == nil {
				full = filepath.Clean(filepath.Join(wd, path))
			} else {
				full = filepath.Clean(path)
			}
		} else {
			full = filepath.Clean(filepath.Join(curDir, path))
		}
	}

	if _, dup := e.imported[e.keyFile(full)]; dup {
		return nil, nil
	}
	if strings.EqualFold(filepath.Ext(full), ".candylib") {
		e.imported[e.keyFile(full)] = struct{}{}
		return e.loadFromCandyLib(full)
	}

	b, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("candy_load: read import %q (resolving from %q): %w", full, errOrigin, err)
	}
	l := candy_lexer.New(string(b))
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, formatParseErrors(full, p.Errors())
	}

	e.imported[e.keyFile(full)] = struct{}{}
	subDir := filepath.Dir(full)
	return e.expandStmts(prog.Statements, subDir, full)
}

func (e *expander) loadFromCandyLib(full string) ([]candy_ast.Statement, error) {
	m, err := candy_bindgen.ParseManifestFile(full)
	if err != nil {
		return nil, fmt.Errorf("candy_load: read candylib %q: %w", full, err)
	}
	if e.ctx != nil {
		e.ctx.mergeManifest(m)
	}
	stmts := make([]candy_ast.Statement, 0, len(m.Externs))
	for _, ex := range m.Externs {
		if why := candy_bindgen.IsUnsafeABI(ex); why != "" && !m.UnsafeABI {
			return nil, fmt.Errorf("candy_load: candylib extern %q rejected: %s", ex.Name, why)
		}
		stmts = append(stmts, externFromManifest(ex))
	}
	return stmts, nil
}

func externFromManifest(ex candy_bindgen.ExternBinding) *candy_ast.ExternFunctionStatement {
	params := make([]candy_ast.Parameter, 0, len(ex.Params))
	for _, p := range ex.Params {
		params = append(params, candy_ast.Parameter{
			Token: candy_token.Token{Type: candy_token.IDENT, Literal: p.Name},
			Name:  &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: p.Name}, Value: p.Name},
			TypeName: &candy_ast.TypeExpression{
				Token: candy_token.Token{Type: candy_token.IDENT, Literal: candy_bindgen.TypeToCandy(p.Type)},
				Name:  &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: candy_bindgen.TypeToCandy(p.Type)}, Value: candy_bindgen.TypeToCandy(p.Type)},
			},
		})
	}
	returnTy := candy_bindgen.TypeToCandy(ex.ReturnType)
	fn := &candy_ast.FunctionStatement{
		Token: candy_token.Token{Type: candy_token.FUNCTION, Literal: "fun"},
		Name:  &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: ex.Name}, Value: ex.Name},
		Parameters: params,
		ReturnType: &candy_ast.TypeExpression{
			Token: candy_token.Token{Type: candy_token.IDENT, Literal: returnTy},
			Name:  &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: returnTy}, Value: returnTy},
		},
	}
	return &candy_ast.ExternFunctionStatement{
		Token:    candy_token.Token{Type: candy_token.EXTERN, Literal: "extern"},
		Function: fn,
	}
}

func (e *expander) loadFromStdlib(path, source, _ string) ([]candy_ast.Statement, error) {
	key := e.keyStdlib(path)
	if _, dup := e.imported[key]; dup {
		return nil, nil
	}

	l := candy_lexer.New(source)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, formatParseErrors("stdlib:"+path, p.Errors())
	}
	e.imported[key] = struct{}{}

	// Like eval: nested file imports in stdlib source resolve from process CWD.
	wd, _ := os.Getwd()
	if wd == "" {
		wd = "."
	}
	return e.expandStmts(prog.Statements, wd, "stdlib:"+path)
}

// HasTopLevelImport reports whether the program has any *ImportStatement at the top level.
func HasTopLevelImport(prog *candy_ast.Program) bool {
	if prog == nil {
		return false
	}
	for _, s := range prog.Statements {
		if _, ok := s.(*candy_ast.ImportStatement); ok {
			return true
		}
	}
	return false
}
