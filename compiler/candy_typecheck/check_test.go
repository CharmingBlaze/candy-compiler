package candy_typecheck

import (
	"strings"
	"testing"

	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_parser"
	"candy/candy_token"
)

func TestCheckProgram_Nil(t *testing.T) {
	if len(CheckProgram(nil)) != 0 {
		t.Fatal("expected no issues for nil program")
	}
}

func TestCheckProgram_empty(t *testing.T) {
	issues := CheckProgram(&candy_ast.Program{})
	if len(issues) != 0 {
		t.Fatal("expected no issues for empty program")
	}
}

func TestCheckInfix_LiteralBooleanPlus(t *testing.T) {
	prog := &candy_ast.Program{
		Statements: []candy_ast.Statement{
			&candy_ast.ExpressionStatement{
				Token: candy_token.Token{Type: candy_token.RETURN},
				Expression: &candy_ast.InfixExpression{
					Token:    candy_token.Token{Type: candy_token.PLUS, Literal: "+", Line: 1, Col: 1},
					Operator: "+",
					Left:     &candy_ast.Boolean{Token: candy_token.Token{Type: candy_token.TRUE, Line: 1, Col: 1}, Value: true},
					Right:    &candy_ast.IntegerLiteral{Token: candy_token.Token{Type: candy_token.INT, Literal: "1", Line: 1, Col: 1}, Value: 1},
				},
			},
		},
	}
	issues := CheckProgram(prog)
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1: %#v", len(issues), issues)
	}
	if issues[0].Line == 0 {
		t.Fatalf("expected Line > 0 in issue, got %d", issues[0].Line)
	}
}

func TestCheckFunctionNonNullableReturnsNull(t *testing.T) {
	prog := &candy_ast.Program{
		Statements: []candy_ast.Statement{
			&candy_ast.FunctionStatement{
				Token:      candy_token.Token{Type: candy_token.FUNCTION, Line: 1, Col: 1},
				Name:       &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "f"}, Value: "f"},
				ReturnType: &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "Int"}, Value: "Int"},
				Body: &candy_ast.BlockStatement{
					Token: candy_token.Token{Type: candy_token.LBRACE},
					Statements: []candy_ast.Statement{
						&candy_ast.ReturnStatement{
							Token:       candy_token.Token{Type: candy_token.RETURN},
							ReturnValue: &candy_ast.NullLiteral{Token: candy_token.Token{Type: candy_token.NULL}},
						},
					},
				},
			},
		},
	}
	issues := CheckProgram(prog)
	if len(issues) == 0 {
		t.Fatal("expected non-nullable return/null issue")
	}
	if !strings.Contains(issues[0].Message, "non-nullable return type") {
		t.Fatalf("unexpected issue: %q", issues[0].Message)
	}
}

func TestCheckFunctionNonNullableReturnsNull_MixedCaseTypeName(t *testing.T) {
	prog := &candy_ast.Program{
		Statements: []candy_ast.Statement{
			&candy_ast.FunctionStatement{
				Token:      candy_token.Token{Type: candy_token.FUNCTION, Line: 1, Col: 1},
				Name:       &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "f"}, Value: "f"},
				ReturnType: &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "STRING"}, Value: "STRING"},
				Body: &candy_ast.BlockStatement{
					Token: candy_token.Token{Type: candy_token.LBRACE},
					Statements: []candy_ast.Statement{
						&candy_ast.ReturnStatement{
							Token:       candy_token.Token{Type: candy_token.RETURN},
							ReturnValue: &candy_ast.NullLiteral{Token: candy_token.Token{Type: candy_token.NULL}},
						},
					},
				},
			},
		},
	}
	issues := CheckProgram(prog)
	if len(issues) == 0 {
		t.Fatal("expected non-nullable return/null issue")
	}
	if !strings.Contains(strings.ToLower(issues[0].Message), "non-nullable return type") {
		t.Fatalf("unexpected issue: %q", issues[0].Message)
	}
}

func TestTypeAssignable_CaseInsensitive(t *testing.T) {
	c := &Checker{}
	if !c.typeAssignable("INT?", "Null") {
		t.Fatal("expected INT? assignable from Null")
	}
	if !c.typeAssignable("STRING", "string") {
		t.Fatal("expected STRING assignable from string")
	}
	if c.typeAssignable("BOOL", "string") {
		t.Fatal("did not expect BOOL assignable from string")
	}
}

func TestStructInheritanceTypecheck(t *testing.T) {
	src := `
struct Entity { name: String };
struct Player : Entity { health: Int };
val e: Entity = Player { name: "Hero", health: 100 };
val n = e.name;
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestWhenTypeConsistency(t *testing.T) {
	w := &candy_ast.WhenExpression{
		Token: candy_token.Token{Type: candy_token.WHEN, Line: 1, Col: 1},
		Arms: []candy_ast.WhenArm{
			{Cond: &candy_ast.Boolean{Token: candy_token.Token{Type: candy_token.TRUE}, Value: true}, Body: &candy_ast.IntegerLiteral{Token: candy_token.Token{Type: candy_token.INT}, Value: 1}},
			{Cond: &candy_ast.Boolean{Token: candy_token.Token{Type: candy_token.FALSE}, Value: false}, Body: &candy_ast.StringLiteral{Token: candy_token.Token{Type: candy_token.STR}, Value: "x"}},
		},
	}
	issues := CheckProgram(&candy_ast.Program{Statements: []candy_ast.Statement{
		&candy_ast.ExpressionStatement{Token: candy_token.Token{Type: candy_token.WHEN}, Expression: w},
	}})
	if len(issues) == 0 {
		t.Fatal("expected type consistency issue for when")
	}
}

func TestBuiltinArity(t *testing.T) {
	prog := &candy_ast.Program{Statements: []candy_ast.Statement{
		&candy_ast.ExpressionStatement{
			Token: candy_token.Token{Type: candy_token.IDENT},
			Expression: &candy_ast.CallExpression{
				Token:    candy_token.Token{Type: candy_token.LPAREN, Line: 1, Col: 1},
				Function: &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT}, Value: "len"},
			},
		},
	}}
	issues := CheckProgram(prog)
	if len(issues) == 0 {
		t.Fatal("expected builtin arity issue")
	}
}

func TestCheckVariableTypeMismatch(t *testing.T) {
	prog := &candy_ast.Program{
		Statements: []candy_ast.Statement{
			&candy_ast.ValStatement{
				Token:    candy_token.Token{Type: candy_token.VAL, Line: 1, Col: 1},
				Name:     &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "x"}, Value: "x"},
				TypeName: &candy_ast.Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "Int"}, Value: "Int"},
				Value:    &candy_ast.StringLiteral{Token: candy_token.Token{Type: candy_token.STR, Literal: "hi"}, Value: "hi"},
			},
		},
	}
	issues := CheckProgram(prog)
	if len(issues) == 0 {
		t.Fatal("expected type mismatch issue")
	}
	if !strings.Contains(issues[0].Message, "cannot assign string to int") {
		t.Fatalf("unexpected issue: %q", issues[0].Message)
	}
}

func TestWhenTypeConsistency_FromParserInput(t *testing.T) {
	src := `return when { true: 1; false: "x"; else: 2; };`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	if len(issues) == 0 {
		t.Fatal("expected type consistency issue for parsed when expression")
	}
}

func TestKotlinStyleDeclsTypecheck(t *testing.T) {
	src := `
class Base { };
class Child extends Base { };
interface I<T> { };
trait Tr { };
object App { };
extern fun native_add(x: Int, y: Int): Int { };
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	issues := CheckProgram(prog)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}
