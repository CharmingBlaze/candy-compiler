package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"testing"
)

const benchTypedModule = `
fun update(pos: float, vel: float, dt: float): float {
  val p = pos + vel * dt;
  val clamped = max(p, 0.0);
  return sqrt(clamped + 1.0);
}

fun tick(count: int): int {
  val acc = 0;
  for i = 0 to count {
    val a = abs(i - 50);
    val b = min(a, 25);
    val c = max(b, 3);
    // Keep loop body non-trivial for codegen.
    val _x = c + 1;
  }
  return count;
}

return tick(1000);
`

func parseBenchProgram(b *testing.B, src string) any {
	b.Helper()
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		b.Fatalf("parser errors: %v", errs)
	}
	return program
}

func BenchmarkGenerateIR_TypedGameModule(b *testing.B) {
	program, ok := parseBenchProgram(b, benchTypedModule).(interface {
		String() string
	})
	if !ok {
		b.Fatal("unexpected program type")
	}
	_ = program.String()

	l := candy_lexer.New(benchTypedModule)
	p := candy_parser.New(l)
	astProgram := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		b.Fatalf("parser errors: %v", errs)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		if _, err := c.GenerateIR(astProgram); err != nil {
			b.Fatalf("GenerateIR failed: %v", err)
		}
	}
}

