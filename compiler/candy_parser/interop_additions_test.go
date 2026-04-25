package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"testing"
)

func TestParseWithStatementAndBitwise(t *testing.T) {
	src := `
with file = open("a.txt") {
  flags = 1 | 2
  flags = flags & ~2
  flags += 1
}
`
	p := New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	if prog == nil || len(prog.Statements) == 0 {
		t.Fatalf("empty program")
	}
}

func TestParseExternWithoutFunKeyword(t *testing.T) {
	src := `extern native_add(a: int, b: int): int`
	p := New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	if _, ok := prog.Statements[0].(*candy_ast.ExternFunctionStatement); !ok {
		t.Fatalf("expected extern function statement, got %T", prog.Statements[0])
	}
}

func TestParseExternVariadicSignature(t *testing.T) {
	src := `extern printf(format: cstring, ...args): int`
	p := New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	ex, ok := prog.Statements[0].(*candy_ast.ExternFunctionStatement)
	if !ok || ex.Function == nil {
		t.Fatalf("expected extern function statement, got %T", prog.Statements[0])
	}
	if !ex.Function.Variadic {
		t.Fatalf("expected variadic extern function")
	}
}
