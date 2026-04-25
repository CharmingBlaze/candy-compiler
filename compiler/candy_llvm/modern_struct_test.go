package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"strings"
	"testing"
)

func TestGenerateIR_ModernStructs(t *testing.T) {
	src := `
struct Vector2 {
    x: float
    y: float

    float length {
        get { return 5.0; }
    }

    Vector2 operator+(Vector2 other) {
        return Vector2 { x: this.x + other.x, y: this.y + other.y }
    }
}

val v1: Vector2 = Vector2 { x: 1.0, y: 2.0 }
val v2: Vector2 = Vector2 { x: 3.0, y: 4.0 }
val v3 = v1 + v2
val l = v1.length
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

	if !strings.Contains(ir, "define double @") || !strings.Contains(ir, "_get_length(") {
		t.Errorf("expected double getter, got: %s", ir)
	}
	if !strings.Contains(ir, "_op_add(") {
		t.Errorf("expected operator _op_add, got: %s", ir)
	}
	// v1, v2 typed as struct ptr => infix uses overload; dot uses static getter
	if !strings.Contains(ir, "call") || !strings.Contains(ir, "_op_add") {
		t.Errorf("expected call to _op_add in @main, got: %s", ir)
	}
	if !strings.Contains(ir, "call") || !strings.Contains(ir, "_get_length") {
		t.Errorf("expected call to _get_length in @main, got: %s", ir)
	}
}

func TestGenerateIR_StructLiteralInferenceNoAny(t *testing.T) {
	src := `
struct Vector2 {
    x: float
    y: float
    Vector2 operator+(Vector2 other) {
        return Vector2 { x: this.x + other.x, y: this.y + other.y }
    }
}
val v1 = Vector2 { x: 1.0, y: 2.0 }
val v2 = Vector2 { x: 3.0, y: 4.0 }
val v3 = v1 + v2
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
	if strings.Contains(ir, "alloca %any") {
		t.Errorf("did not expect any-typed alloca for inferred struct locals, got: %s", ir)
	}
	if !strings.Contains(ir, "_op_add(") {
		t.Errorf("expected operator call, got: %s", ir)
	}
}
