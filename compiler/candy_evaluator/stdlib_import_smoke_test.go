package candy_evaluator

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"testing"
)

func TestEval_StdlibImportSmoke_AllModules(t *testing.T) {
	modules := []string{
		"candy.2d",
		"candy.3d",
		"candy.audio",
		"candy.ai",
		"candy.camera",
		"candy.debug",
		"candy.editor",
		"candy.game",
		"candy.game3d",
		"candy.input",
		"candy.network",
		"candy.physics2d",
		"candy.physics3d",
		"candy.proc",
		"candy.resources",
		"candy.save",
		"candy.scene",
		"candy.state",
		"candy.ui",
		"candy.vfx",
	}

	for _, mod := range modules {
		mod := mod
		t.Run(mod, func(t *testing.T) {
			src := "import " + mod + "\ntrue"
			l := candy_lexer.New(src)
			p := candy_parser.New(l)
			prog := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("parse: %v", p.Errors())
			}
			v, err := Eval(prog, nil)
			if err != nil {
				t.Fatalf("eval: %v", err)
			}
			if v == nil || v.Kind != ValBool || !v.B {
				t.Fatalf("expected true, got %v", v)
			}
		})
	}
}
