package candy_typecheck

import (
	"strings"
	"testing"

	"candy/candy_lexer"
	"candy/candy_parser"
)

// TestCheckProgram_FullSpecComposite parses fragments aligned with candy_parser/full_spec_test.go
// and runs the shallow typechecker. Expects no panics; issue count is bounded (many names are
// still unknown in this pass).
func TestCheckProgram_FullSpecComposite(t *testing.T) {
	t.Parallel()
	const src = `
package demo;

[Serializable(format: "binary")]
struct Data {
	[SaveField]
	int id = 0
}

struct Player {
	private int _health = 100
	int health {
		get { return _health }
		set { _health = 1 }
	}
	string name { get; set; } = "Player"
	Vector2 operator+(Vector2 other) {
		return other
	}
}

module math {
	const PI = 3.14159
	export float sqrt(float x) { return x; }
}

enum GameState {
	Menu,
	Playing,
	Paused = 10,
	GameOver
}

try {
	noop()
} catch FileException e {
	print(1)
} finally {
	cleanup()
}

async string loadData(string url) {
	return await url
}
run loadData("https://api")
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	// Shallow check: no crash; diagnostics stay finite (unknown ids / builtins differ by snippet).
	if len(issues) > 400 {
		if len(issues) > 0 {
			t.Fatalf("unexpectedly many typecheck issues: %d: %q", len(issues), issues[0].Message)
		}
		t.Fatalf("unexpectedly many typecheck issues: %d", len(issues))
	}
	// At least one expected unknown-name style warning from minimal stubs.
	if len(issues) == 0 {
		t.Fatal("expected at least one typecheck issue for the composite stub program")
	}
}

// TestCheckProgram_FullSpecMemorySlice covers full_spec memory modifiers: parsed exprs must be walked.
func TestCheckProgram_FullSpecMemorySlice(t *testing.T) {
	t.Parallel()
	const src = `
val x: Int = 1
int* p = 1
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	_ = CheckProgram(prog)
}

// TestCheckProgram_AwaitInExpr ensures await/assign and receiver paths do not skip expr kinds.
func TestCheckProgram_AwaitInExpr(t *testing.T) {
	t.Parallel()
	const src = `
struct S { }
struct M {
	unit f(S self) {
		_ = await 1
		_ = notdeclared
	}
}
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	// 'a' unknown; await inner literal ok.
	foundUnknown := false
	for _, d := range issues {
		if strings.Contains(d.Message, "unknown") {
			foundUnknown = true
			break
		}
	}
	if !foundUnknown {
		t.Fatalf("expected an unknown-identifier issue, got %v", issues)
	}
}
