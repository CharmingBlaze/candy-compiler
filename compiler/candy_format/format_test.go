package candy_format

import (
	"strings"
	"testing"

	"candy/candy_lexer"
	"candy/candy_parser"
)

func TestSourceDeterministic(t *testing.T) {
	src := `val x=1;
fun f(n: Int): Int { return n + 1; };
`
	p := candy_parser.New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	out := Source(prog)
	if !strings.Contains(out, "val x = 1;") {
		t.Fatalf("formatted output missing val statement: %q", out)
	}
	if !strings.Contains(out, "fun f(") {
		t.Fatalf("formatted output missing function: %q", out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Fatalf("formatted output should end with newline: %q", out)
	}
}
