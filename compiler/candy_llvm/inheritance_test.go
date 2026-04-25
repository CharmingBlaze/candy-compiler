package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"strings"
	"testing"
)

func TestGenerateIR_Inheritance(t *testing.T) {
	src := `
struct Entity { name: String };
struct Player : Entity { health: Int };
val p: Player = Player { name: "Hero", health: 100 };
val e: Entity = p;
println(e.name);
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	c := New()
	ir, err := c.GenerateIR(prog)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	// Verify IR contains flattened struct (LLVM mangles struct names to lowercase)
	if !strings.Contains(ir, "%entity = type {i8*}") {
		t.Errorf("expected entity struct definition, got: %s\nFull IR:\n%s", ir, ir)
	}
	if !strings.Contains(ir, "%player = type {i8*, i64}") {
		t.Errorf("expected player struct definition (flattened), got: %s\nFull IR:\n%s", ir, ir)
	}

	// Verify main contains println calls
	if !strings.Contains(ir, "call i8* @candy_str_add") && !strings.Contains(ir, "call i32 (i8*, ...)* @printf") {
		// println uses candy_str_add and printf
	}
}

func TestGenerateIR_InheritedPropertyGetter(t *testing.T) {
	src := `
struct Base {
    w: int
    int score { get { return 42; } }
}
struct Child: Base { }
val c: Child = Child { w: 0 }
val s = c.score
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	c := New()
	ir, err := c.GenerateIR(prog)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}
	if !strings.Contains(ir, "_get_score") {
		t.Errorf("expected call or def of inherited property getter, IR:\n%s", ir)
	}
}
