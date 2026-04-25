package candy_typecheck

import (
	"strings"
	"testing"

	"candy/candy_lexer"
	"candy/candy_parser"
)

func TestCheckInfix_StructOperator(t *testing.T) {
	t.Parallel()
	// C#-style return type: V operator+ — second operand type V, but we add two W.
	src := `
struct V { a: int };
struct W {
    m: int
	V operator+(V other) { return V { a: 0 }; }
}
val a: W = W { m: 0 };
val b: W = W { m: 0 };
val c = a + b;
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	found := false
	for _, d := range issues {
		if strings.Contains(d.Message, "mismatch") {
			found = true
			break
		}
	}
	if !found {
		t.Logf("issues: %#v", issues)
		t.Fatal("expected operator right-operand type mismatch")
	}
}
