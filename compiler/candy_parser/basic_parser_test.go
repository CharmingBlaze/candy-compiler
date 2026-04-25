package candy_parser

import (
	"candy/candy_lexer"
	"testing"
)

func TestSimpleCSyntax(t *testing.T) {
	input := `
name = "Alice"
age = 25
int count = 10
float speed = 5.5
player.x = 50
if player.x > 0 {
    y = 20
}

sub greet(name) {
    msg = "Hello, " + name
}

for i = 1 to 10 {
    total = total + i
}

while playing {
    update()
}

struct Vector2 {
    float x, y
    float length() {
        return 1.0
    }
}

for (int i = 0; i < 10; i = i + 1) {
    print i
}

add = (a, b) => a + b
doubled = nums.map(n => n * 2)
defer cleanup()

struct Entity {
    string name
    Vector2 pos
    
    int init(string n, float x, float y) {
        name = n
        pos = Vector2 { x: x, y: y }
        return 0
    }
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 15 {
		t.Fatalf("expected 15 statements, got %d", len(program.Statements))
	}
}

// `in expr` must not treat `a {` or `] {` as a struct literal; the `{` is the for body.
func TestForInIterableStopsBeforeBlockBrace(t *testing.T) {
	for _, src := range []string{
		`for k in a { n = 1 }`,
		`for k in [1, 2, 3] { n = 1 }`,
	} {
		t.Run(src, func(t *testing.T) {
			l := candy_lexer.New(src)
			p := New(l)
			_ = p.ParseProgram()
			checkParserErrors(t, p)
		})
	}
}

func TestFunctionFunAndFuncAreAliases(t *testing.T) {
	for _, kw := range []string{"function", "fun", "func"} {
		t.Run(kw, func(t *testing.T) {
			src := kw + ` add(a, b) { return a + b; }`
			l := candy_lexer.New(src)
			p := New(l)
			_ = p.ParseProgram()
			checkParserErrors(t, p)
		})
	}
}
