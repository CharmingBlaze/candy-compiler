package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"strings"
	"testing"
)

func TestGenerateIR(t *testing.T) {
	input := `
val x = 10 + 20;
return x;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := New()
	ir, err := comp.GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}

	expectedLines := []string{
		"define i64 @main() nounwind willreturn hot inlinehint {",
		"%x = alloca i64",
		"store i64 30",
		"store i64",
		"load i64",
		"ret i64",
		"}",
	}

	for _, line := range expectedLines {
		if !strings.Contains(ir, line) {
			t.Errorf("expected IR to contain %q, but it didn't.\nGot:\n%s", line, ir)
		}
	}
}

func TestGenerateIR_FunctionAttributes(t *testing.T) {
	input := `
fun update(): int { return 1; };
fun math_dot(): float { return 1.0; };
return update();
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "define i64 @update() nounwind willreturn hot inlinehint {") {
		t.Fatalf("expected hot function attrs for update; got:\n%s", ir)
	}
	if !strings.Contains(ir, "define double @math_dot() nounwind willreturn readnone nosync speculatable {") {
		t.Fatalf("expected pure math attrs for math_dot; got:\n%s", ir)
	}
}

func TestGenerateIR_TypedMathFastPaths(t *testing.T) {
	input := `
val a = abs(-7);
val b = min(10, 3);
val c = max(2.0, 5.5);
return a;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	// abs/min/max should lower to compare+select, not unknown runtime calls.
	for _, want := range []string{
		"icmp slt i64",
		"select i1",
		"icmp slt i64",
		"fcmp ogt double",
	} {
		if !strings.Contains(ir, want) {
			t.Fatalf("expected typed fast-path marker %q; got:\n%s", want, ir)
		}
	}
	if strings.Contains(ir, "call i64 @abs(") || strings.Contains(ir, "call i64 @min(") || strings.Contains(ir, "call i64 @max(") {
		t.Fatalf("expected no generic call for abs/min/max fast paths; got:\n%s", ir)
	}
}

func TestGenerateIR_Float(t *testing.T) {
	input := `
val x = 10 + 2.5;
return x;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()

	comp := New()
	ir, err := comp.GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}

	expectedLines := []string{
		"alloca double",
		"store double 12.500000",
		"fptosi double",
		"ret i64",
	}

	for _, line := range expectedLines {
		if !strings.Contains(ir, line) {
			t.Errorf("expected IR to contain %q, but it didn't.\nGot:\n%s", line, ir)
		}
	}
}

func TestGenerateIR_Println(t *testing.T) {
	input := `println("hello world");`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := New()
	ir, err := comp.GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}

	expectedLines := []string{
		"declare i32 @printf",
		"@.str.s1 = private unnamed_addr constant [12 x i8] c\"hello world\\00\"",
		"getelementptr inbounds [12 x i8]",
		"call i32 (i8*, ...) @printf",
	}

	for _, line := range expectedLines {
		if !strings.Contains(ir, line) {
			t.Errorf("expected IR to contain %q, but it didn't.\nGot:\n%s", line, ir)
		}
	}
}

func TestGenerateIR_TypeNarrowing_IsGuard(t *testing.T) {
	input := `
struct Player { int health; };
fun Player_Update(self: Player): int { return self.health; };
val payload = Player{ health: 42 };
if (payload is Player) {
	return payload.health;
};
return 0;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	comp := New()
	ir, err := comp.GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	// payload is inferred as %player* (struct literal), so `payload is Player` is a static check (br i1 1), not %any tag compare.
	expected := []string{
		"%any = type { i8, i8* }",
		"br i1",
		"getelementptr inbounds %player",
	}
	for _, line := range expected {
		if !strings.Contains(ir, line) {
			t.Fatalf("expected IR to contain %q; got:\n%s", line, ir)
		}
	}
}

func TestGenerateIR_While(t *testing.T) {
	input := `
val s = 0;
val n = 3;
while n > 0 {
  s = s + n;
  n = n - 1;
}
return s;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	for _, want := range []string{"br i1", "br label", "icmp", "sub i64"} {
		if !strings.Contains(ir, want) {
			t.Errorf("expected IR to contain %q; got:\n%s", want, ir)
		}
	}
}

func TestGenerateIR_DoWhile(t *testing.T) {
	input := `
val s = 0;
val n = 0;
do {
  s = s + 1;
  n = n + 1;
} while (n < 3);
return s;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "br i1") || !strings.Contains(ir, "icmp") {
		t.Fatalf("expected control-flow IR, got:\n%s", ir)
	}
}

func TestGenerateIR_ForNumeric(t *testing.T) {
	input := `
val s = 0;
for i = 1 to 3 {
  s = s + i;
}
return s;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	for _, want := range []string{"or i1", "icmp sle", "add i64", "br i1"} {
		if !strings.Contains(ir, want) {
			t.Errorf("expected IR to contain %q; got:\n%s", want, ir)
		}
	}
}

func TestGenerateIR_CFor(t *testing.T) {
	input := `
val s = 0;
for (val i = 0; i < 4; i = i + 1) {
  s = s + i;
}
return s;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "cfh") {
		t.Errorf("expected C-style for header label; got:\n%s", ir)
	}
}

func TestGenerateIR_SwitchInt(t *testing.T) {
	input := `
val x = 1;
val r = 0;
switch (x) {
  case 1: { r = 5; }
  default: { r = 9; }
}
return r;
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	ir, err := New().GenerateIR(program)
	if err != nil {
		t.Fatalf("GenerateIR: %v", err)
	}
	if !strings.Contains(ir, "icmp eq i64") {
		t.Errorf("expected integer switch icmp; got:\n%s", ir)
	}
}

func TestGenerateIR_ImportError(t *testing.T) {
	input := `import "x"; return 0;`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	_, err := New().GenerateIR(program)
	if err == nil {
		t.Fatal("expected GenerateIR to fail for import in native path")
	}
}
