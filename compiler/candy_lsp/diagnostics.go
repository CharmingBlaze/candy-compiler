package candy_lsp

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_typecheck"
)

// DiagnosticSeverity maps loosely to LSP severities.
type DiagnosticSeverity int

const (
	SeverityError DiagnosticSeverity = 1
	SeverityWarn  DiagnosticSeverity = 2
)

// Diagnostic is a diagnostics-first payload suitable for an LSP bridge.
type Diagnostic struct {
	Severity DiagnosticSeverity
	Message  string
}

// AnalyzeSource returns parser and static-check diagnostics.
func AnalyzeSource(src string) []Diagnostic {
	out := make([]Diagnostic, 0)
	p := candy_parser.New(candy_lexer.New(src))
	prog := p.ParseProgram()
	for _, d := range p.Errors() {
		out = append(out, Diagnostic{Severity: SeverityError, Message: d.Message})
	}
	if len(p.Errors()) > 0 {
		return out
	}
	for _, d := range candy_typecheck.CheckProgram(prog) {
		out = append(out, Diagnostic{Severity: SeverityWarn, Message: d.Message})
	}
	return out
}
