package candy_llvm

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_typecheck"
	"fmt"
	"testing"
)

func TestDynamicSimpleC(t *testing.T) {
	input := `
// LEVEL 1: Dynamic Basics
name = "Alice"
score = 100
print "Hello, {name}! Your score is {score}"

// LEVEL 2: Structured
struct Player {
    name
    health = 100
    
    func info() {
        print "{name} HP: {health}"
    }
}

p = Player { name: "Bob" }
p.info()

// LEVEL 3: Performance (Static)
int count = 5
print "Count is {count}"
`
	l := candy_lexer.New(input)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	fmt.Printf("Parser errors: %d\n", len(p.Errors()))
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Printf("Parser error: %s\n", err.Message)
			t.Errorf("Parser error: %s", err.Message)
		}
		t.FailNow()
	}

	issues := candy_typecheck.CheckProgram(prog)
	fmt.Printf("TypeChecker issues: %d\n", len(issues))
	for _, is := range issues {
		t.Logf("TypeChecker issue: %s", is.Message)
	}

	comp := New()
	ir, err := comp.GenerateIR(prog)
	if err != nil {
		t.Fatalf("Codegen error: %s", err)
	}

	fmt.Println("--- GENERATED IR ---")
	fmt.Println(ir)
	fmt.Println("--------------------")
}
