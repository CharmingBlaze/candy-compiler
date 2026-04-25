package candy_parser

import (
	"candy/candy_lexer"
	"testing"
)

func TestFullSpecModules(t *testing.T) {
	input := `
module math {
    const PI = 3.14159
    export float sqrt(float x) { }
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestFullSpecEnums(t *testing.T) {
	input := `
enum GameState {
    Menu,
    Playing,
    Paused = 10,
    GameOver
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestFullSpecStructsAndProperties(t *testing.T) {
	input := `
struct Player : IDamageable, IRenderable {
    private int _health = 100
    
    int health {
        get { return _health }
        set { _health = value }
    }
    
    string name { get; set; } = "Player"

    Vector2 operator+(Vector2 other) {
        return Vector2 { x: x + other.x, y: y + other.y }
    }
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestFullSpecTryCatch(t *testing.T) {
	input := `
try {
    file = openFile("data.txt")
} catch FileException e {
    print e.message
} finally {
    cleanup()
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestFullSpecAsyncAwait(t *testing.T) {
	input := `
async string loadData(string url) {
    response = await http.get(url)
    return await response.text()
}

run loadData("https://api")
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestFullSpecMemoryModifiers(t *testing.T) {
	input := `
ref Player p = player
shared Texture tex = new Texture("sprite.png")
maybe Player target = null
int* ptr = &value
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestFullSpecAttributes(t *testing.T) {
	input := `
[Serializable(format: "binary")]
struct Data {
    [SaveField]
    int id
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}
