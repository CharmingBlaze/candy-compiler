package candy_parser

import (
	"candy/candy_lexer"
	"testing"
)

func TestParseLibrarySyntax(t *testing.T) {
	src := `
library "raylib" {
  type Color {
    r: int
    g: int
    b: int
    a: int
  }
  extern fun DrawTexture(id: int, x: int, y: int): void {}
}
`
	p := New(candy_lexer.New(src))
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	if prog == nil || len(prog.Statements) != 1 {
		t.Fatalf("expected one top-level library statement")
	}
}
