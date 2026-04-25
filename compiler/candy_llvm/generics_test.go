package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_typecheck"
	"strings"
	"testing"
)

func TestGenerateIR_Generics(t *testing.T) {
	src := `
struct Box<T> {
    value: T;
}

val b = Box<float> { value: 3.14 };
val v = b.value;
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Errorf("Parser error: %v", err)
		}
		t.FailNow()
	}

	// Typecheck to trigger monomorphization
	candy_typecheck.CheckProgram(prog)
	t.Logf("Program after typecheck:\n%s", prog.String())

	c := New()
	ir, err := c.GenerateIR(prog)
	if err != nil {
		t.Fatalf("IR generation failed: %v", err)
	}

	t.Logf("Generated IR:\n%s", ir)

	// Verify specialized struct exists in IR (monomorph uses lowercased name; LLVM omits extra spaces)
	if !strings.Contains(ir, "%box_float = type {double}") {
		t.Errorf("Expected specialized struct box_float in IR, not found")
	}
}
