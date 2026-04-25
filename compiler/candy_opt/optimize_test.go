package candy_opt

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_parser"
	"testing"
)

func parseProgram(t *testing.T, src string) *candy_ast.Program {
	t.Helper()
	p := candy_parser.New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}
	return prog
}

func TestOptimizeProgram_ConstantFolding(t *testing.T) {
	prog := parseProgram(t, `val n = (2 + 3) * 4`)
	OptimizeProgram(prog)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	v, ok := prog.Statements[0].(*candy_ast.ValStatement)
	if !ok {
		t.Fatalf("expected val statement, got %T", prog.Statements[0])
	}
	lit, ok := v.Value.(*candy_ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected folded integer literal, got %T", v.Value)
	}
	if lit.Value != 20 {
		t.Fatalf("expected folded value 20, got %d", lit.Value)
	}
}

func TestOptimizeProgram_DeadBranchElimination(t *testing.T) {
	prog := parseProgram(t, `
if false {
  val a = 1
} else {
  val b = 2
}
`)
	OptimizeProgram(prog)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement after DCE, got %d", len(prog.Statements))
	}
	v, ok := prog.Statements[0].(*candy_ast.ValStatement)
	if !ok || v.Name == nil || v.Name.Value != "b" {
		t.Fatalf("expected else-branch val b statement, got %T", prog.Statements[0])
	}
}

func TestOptimizeProgram_RemoveNeverRunLoop(t *testing.T) {
	prog := parseProgram(t, `
while false {
  val x = 1
}
val ok = 1
`)
	OptimizeProgram(prog)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected only trailing statement, got %d", len(prog.Statements))
	}
	v, ok := prog.Statements[0].(*candy_ast.ValStatement)
	if !ok || v.Name == nil || v.Name.Value != "ok" {
		t.Fatalf("expected val ok, got %T", prog.Statements[0])
	}
}

func TestOptimizeProgram_InlinesSimplePureFunction(t *testing.T) {
	prog := parseProgram(t, `
fun add1(x) { return x + 1 }
val out = add1(41)
`)
	OptimizeProgram(prog)
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Statements))
	}
	v, ok := prog.Statements[1].(*candy_ast.ValStatement)
	if !ok {
		t.Fatalf("expected val statement, got %T", prog.Statements[1])
	}
	lit, ok := v.Value.(*candy_ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected fully folded inlined integer literal, got %T", v.Value)
	}
	if lit.Value != 42 {
		t.Fatalf("expected inlined/folded value 42, got %d", lit.Value)
	}
}
