package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"strings"
	"testing"
)

func TestGenerateIR_PrefixMinusAndBang(t *testing.T) {
	src := `
	val a = 3;
	return -a;
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	ir, err := New().GenerateIR(prog)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, " sub i64 ") {
		t.Errorf("expected sub i64 for unary minus, got: %s", ir)
	}
}

func TestGenerateIR_PrefixBangInt(t *testing.T) {
	src := `val x = 0; return !x;`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	ir, err := New().GenerateIR(prog)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "icmp eq i64") {
		t.Errorf("expected icmp for ! on int, got: %s", ir)
	}
}

func TestGenerateIR_StringIndex(t *testing.T) {
	src := `val s = "ab"; return s[0];`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	ir, err := New().GenerateIR(prog)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "getelementptr inbounds i8") {
		t.Errorf("expected gep for string index, got: %s", ir)
	}
	if !strings.Contains(ir, "zext i8") {
		t.Errorf("expected zext i8 for char load, got: %s", ir)
	}
}
