package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"strings"
	"testing"
)

func TestValStatements(t *testing.T) {
	input := `
val x = 5;
val y = 10;
val foobar = 838383;
`
	l := candy_lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testValStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testValStatement(t *testing.T, s candy_ast.Statement, name string) bool {
	if s.TokenLiteral() != "val" {
		t.Errorf("s.TokenLiteral not 'val'. got=%q", s.TokenLiteral())
		return false
	}

	valStmt, ok := s.(*candy_ast.ValStatement)
	if !ok {
		t.Errorf("s not *candy_ast.ValStatement. got=%T", s)
		return false
	}

	if valStmt.Name.Value != name {
		t.Errorf("valStmt.Name.Value not '%s'. got=%s", name, valStmt.Name.Value)
		return false
	}

	if valStmt.Name.TokenLiteral() != name {
		t.Errorf("valStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, valStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y = 10;
var foobar = 838383;
`
	l := candy_lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestVarStatements_WithoutInitializer(t *testing.T) {
	input := `
var a
var b: Int
int c
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}
	for i, st := range program.Statements {
		vs, ok := st.(*candy_ast.VarStatement)
		if !ok {
			t.Fatalf("stmt[%d] expected *VarStatement, got %T", i, st)
		}
		if vs.Value != nil {
			t.Fatalf("stmt[%d] expected nil initializer, got %#v", i, vs.Value)
		}
	}
}

func testVarStatement(t *testing.T, s candy_ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*candy_ast.VarStatement)
	if !ok {
		t.Errorf("s not *candy_ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := candy_lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*candy_ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *candy_ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

// Documents LANGUAGE.md: type-first var, sub, and fun with return type (teaching “equivalents” smoke test).
func TestParse_EquivalentVariablesAndFunctions(t *testing.T) {
	input := `
int x = 0
sub f() { }
fun g(): int { return 1; }
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*candy_ast.VarStatement); !ok {
		t.Fatalf("stmt0: want *VarStatement, got %T", program.Statements[0])
	}
	if sub, ok := program.Statements[1].(*candy_ast.FunctionStatement); !ok {
		t.Fatalf("stmt1: want *FunctionStatement (sub), got %T", program.Statements[1])
	} else if sub.Name == nil || sub.Name.Value != "f" {
		t.Fatalf("sub name: got %v", sub.Name)
	}
	if fun, ok := program.Statements[2].(*candy_ast.FunctionStatement); !ok {
		t.Fatalf("stmt2: want *FunctionStatement (fun), got %T", program.Statements[2])
	} else if fun.Name == nil || fun.Name.Value != "g" {
		t.Fatalf("fun name: got %v", fun.Name)
	} else if fun.ReturnType == nil {
		t.Fatalf("fun: expected return type")
	}
}

func TestTypeAnnotatedVariables(t *testing.T) {
	input := `
val x: Int = 5;
var y: String = "hi";
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	v1 := program.Statements[0].(*candy_ast.ValStatement)
	if v1.TypeName == nil || candy_ast.ExprAsSimpleTypeName(v1.TypeName) != "int" {
		t.Errorf("expected type 'int', got %v", v1.TypeName)
	}

	v2 := program.Statements[1].(*candy_ast.VarStatement)
	if v2.TypeName == nil || candy_ast.ExprAsSimpleTypeName(v2.TypeName) != "string" {
		t.Errorf("expected type 'string', got %v", v2.TypeName)
	}
}

func TestOptionalSemicolons(t *testing.T) {
	input := `
val x = 5
var y = 10
fun add(a: Int, b: Int): Int {
  return a + b
}
val z = add(x, y)
`
	l := candy_lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestNullableTypeSuffix(t *testing.T) {
	input := `
fun maybe(n: Int?): Int? { return null; };
struct Box { value: Int?; };
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	fn, ok := program.Statements[0].(*candy_ast.FunctionStatement)
	if !ok {
		t.Fatalf("stmt[0] not function, got %T", program.Statements[0])
	}
	if candy_ast.ExprAsSimpleTypeName(fn.Parameters[0].TypeName) != "int?" {
		t.Fatalf("param type = %q, want int?", candy_ast.ExprAsSimpleTypeName(fn.Parameters[0].TypeName))
	}
	if fn.ReturnType == nil || candy_ast.ExprAsSimpleTypeName(fn.ReturnType) != "int?" {
		t.Fatalf("return type = %v, want int?", fn.ReturnType)
	}
	st, ok := program.Statements[1].(*candy_ast.StructStatement)
	if !ok {
		t.Fatalf("stmt[1] not struct, got %T", program.Statements[1])
	}
	if len(st.Fields) != 1 || candy_ast.ExprAsSimpleTypeName(st.Fields[0].TypeName) != "int?" {
		t.Fatalf("struct field type mismatch: %#v", st.Fields)
	}
}

func TestImportAliasAndFromImportParsing(t *testing.T) {
	input := `
import math as m
from math import sin, cos
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	i1, ok := program.Statements[0].(*candy_ast.ImportStatement)
	if !ok {
		t.Fatalf("stmt0 expected import, got %T", program.Statements[0])
	}
	if i1.Path != "math" || i1.Alias != "m" {
		t.Fatalf("unexpected import alias parse: %+v", i1)
	}
	i2, ok := program.Statements[1].(*candy_ast.ImportStatement)
	if !ok {
		t.Fatalf("stmt1 expected import, got %T", program.Statements[1])
	}
	if i2.From != "math" || len(i2.Symbols) != 2 || i2.Symbols[0] != "sin" || i2.Symbols[1] != "cos" {
		t.Fatalf("unexpected from import parse: %+v", i2)
	}
}

func TestNamedCallArgumentParsing(t *testing.T) {
	input := `drawCube(x: 1, y: 2, z: 3)`
	l := candy_lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	es, ok := prog.Statements[0].(*candy_ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected expression statement, got %T", prog.Statements[0])
	}
	call, ok := es.Expression.(*candy_ast.CallExpression)
	if !ok || len(call.Arguments) != 3 {
		t.Fatalf("expected call with 3 args")
	}
	if _, ok := call.Arguments[0].(*candy_ast.NamedArgumentExpression); !ok {
		t.Fatalf("first arg should be named argument")
	}
}

func TestParse_NewStuff12CoreHelperCalls(t *testing.T) {
	input := `
pw = PhysicsWorld(vec3(0, -28, 0))
imap = InputMap()
imap.bind("jump", "space")
imap.bindAxis2D("move", "w", "s", "a", "d")
cam = OrbitCamera()
cc = CharacterController()
ents = EntityList()
hud = HUD()
ui = UILayout()
sm = StateMachine("playing")
tw = Tween()
tf = Transform()
`
	l := candy_lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 12 {
		t.Fatalf("expected 12 statements, got %d", len(prog.Statements))
	}
}

func TestParse_MultilineStructLiteralAsCallArg(t *testing.T) {
	input := `
add(Platform {
    position: vec3(0, 0, 0),
    size: vec3(30, 1, 30),
    color: "darkgray"
})
`
	l := candy_lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	checkParserErrors(t, p)
}

func TestNullableTypeSuffix_MixedCaseInputCanonicalized(t *testing.T) {
	input := `
fun maybe(n: INT?): StRiNg? { return null; };
struct Box { value: BoOl?; };
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	fn, ok := program.Statements[0].(*candy_ast.FunctionStatement)
	if !ok {
		t.Fatalf("stmt[0] not function, got %T", program.Statements[0])
	}
	if candy_ast.ExprAsSimpleTypeName(fn.Parameters[0].TypeName) != "int?" {
		t.Fatalf("param type = %q, want int?", candy_ast.ExprAsSimpleTypeName(fn.Parameters[0].TypeName))
	}
	if fn.ReturnType == nil || candy_ast.ExprAsSimpleTypeName(fn.ReturnType) != "string?" {
		t.Fatalf("return type = %v, want string?", fn.ReturnType)
	}
	st, ok := program.Statements[1].(*candy_ast.StructStatement)
	if !ok {
		t.Fatalf("stmt[1] not struct, got %T", program.Statements[1])
	}
	if len(st.Fields) != 1 || candy_ast.ExprAsSimpleTypeName(st.Fields[0].TypeName) != "bool?" {
		t.Fatalf("struct field type mismatch: %#v", st.Fields)
	}
}

func TestReceiverSyntaxFunction(t *testing.T) {
	input := `fun (self: Box) size(): Int { return 1; };`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	fn, ok := program.Statements[0].(*candy_ast.FunctionStatement)
	if !ok {
		t.Fatalf("stmt not function: %T", program.Statements[0])
	}
	if fn.Name == nil || fn.Name.Value != "size" {
		t.Fatalf("method name mismatch: %#v", fn.Name)
	}
	if fn.Receiver == nil || fn.Receiver.Name == nil || fn.Receiver.TypeName == nil {
		t.Fatalf("receiver not captured: %#v", fn.Receiver)
	}
	if fn.Receiver.Name.Value != "self" || candy_ast.ExprAsSimpleTypeName(fn.Receiver.TypeName) != "box" {
		t.Fatalf("receiver mismatch: %#v", fn.Receiver)
	}
}

func TestMalformedInput_BoundedErrorsAndRecovery(t *testing.T) {
	input := `
fun bad(x: Int) Int { return ;   // malformed function signature/body
val ok = 7;
return ok;
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatal("expected parser errors for malformed input")
	}
	if len(p.Errors()) > 30 {
		t.Fatalf("too many parser errors after recovery hardening: %d", len(p.Errors()))
	}

	foundVal := false
	for _, st := range program.Statements {
		if vs, ok := st.(*candy_ast.ValStatement); ok && vs.Name != nil && vs.Name.Value == "ok" {
			foundVal = true
			break
		}
	}
	if !foundVal {
		t.Fatalf("parser did not recover to following valid statement; statements=%#v errors=%v", program.Statements, p.Errors())
	}
}

func TestMalformedInput_DeduplicatesRepeatedMessages(t *testing.T) {
	input := `fun f(x: Int) Int { return ;`
	l := candy_lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatal("expected parse errors")
	}
	seen := map[string]int{}
	for _, diag := range p.Errors() {
		seen[diag.Message]++
		if seen[diag.Message] > 1 {
			t.Fatalf("duplicate error message found after dedupe: %q", diag.Message)
		}
	}
}

func TestMalformedInput_NoInternalRecoveryNoise(t *testing.T) {
	input := `fun f(x: Int) Int { return ;`
	l := candy_lexer.New(input)
	p := New(l)
	p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Fatal("expected parse errors")
	}
	for _, diag := range p.Errors() {
		if strings.Contains(diag.Message, "parser recovery stalled") {
			t.Fatalf("internal recovery message leaked to user-facing errors: %q", diag.Message)
		}
	}
}

func TestWhenExpressionParsing(t *testing.T) {
	input := `return when { true: 1; else: 2; };`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	ret, ok := program.Statements[0].(*candy_ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected return statement, got %T", program.Statements[0])
	}
	w, ok := ret.ReturnValue.(*candy_ast.WhenExpression)
	if !ok {
		t.Fatalf("expected when expression, got %T", ret.ReturnValue)
	}
	if len(w.Arms) != 1 {
		t.Fatalf("expected 1 when arm, got %d", len(w.Arms))
	}
	if w.ElseV == nil {
		t.Fatal("expected else arm to be parsed")
	}
}


func TestInlineObjectLiteralParsing(t *testing.T) {
	input := `point = {x: 10, y: 20};`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestMapLiteral_AllowsParenWrapperAndNoSemicolons(t *testing.T) {
	input := `
obj1 = map(
  "a": 1
  "b": 2
)
obj2 = {
  x: 10
  y: 20
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestSwitchCaseColonStyleParsing(t *testing.T) {
	input := `
switch value {
  case 1:
    doOne()
  case 2:
    doTwo()
  case 3:
    doThree()
  default:
    doDefault()
}
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*candy_ast.SwitchStatement); !ok {
		t.Fatalf("expected switch statement, got %T", program.Statements[0])
	}
}

func TestExclusiveRangeAndSliceParsing(t *testing.T) {
	input := `
r = 1..<10
items = [10, 20, 30, 40, 50]
subset = items[1..3]
subset2 = items[1..<3]
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestInOperatorParsing(t *testing.T) {
	input := `
ok1 = 3 in [1, 2, 3]
obj = {name: "alice"}
ok2 = "name" in obj
ok3 = "el" in "hello"
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestNotInOperatorParsing(t *testing.T) {
	input := `
ok = 9 not in [1, 2, 3]
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestTernaryOperatorParsing(t *testing.T) {
	input := `
msg = score > 100 ? "High Score!" : "Keep trying"
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestKotlinStyleDeclarationsParse(t *testing.T) {
	input := `
package demo;
interface Shape<T> { };
trait Renderable { };
class Base { };
sealed class Option<T> { };
class Box<T> extends Base { };
object App { };
extern fun native_add(x: Int, y: Int): Int { };
suspend fun fetch(): Int { return 1; };
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 9 {
		t.Fatalf("expected 9 statements, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*candy_ast.PackageStatement); !ok {
		t.Fatalf("stmt0 expected package, got %T", program.Statements[0])
	}
	if s, ok := program.Statements[4].(*candy_ast.ClassStatement); !ok || !s.Sealed {
		t.Fatalf("stmt4 expected sealed class, got %T %#v", program.Statements[4], program.Statements[4])
	}
	if s, ok := program.Statements[7].(*candy_ast.ExternFunctionStatement); !ok || s.Function == nil {
		t.Fatalf("stmt7 expected extern function, got %T", program.Statements[7])
	}
	if f, ok := program.Statements[8].(*candy_ast.FunctionStatement); !ok || !f.Suspend {
		t.Fatalf("stmt8 expected suspend function, got %T %#v", program.Statements[8], program.Statements[8])
	}
}

func TestKotlinStyleDeclarationsNoSemicolons(t *testing.T) {
	input := `
package demo
interface Shape<T> { }
trait Renderable { }
class Base { }
sealed class Option<T> { }
class Box<T> extends Base { }
object App { }
extern fun native_add(x: Int, y: Int): Int { }
suspend fun fetch(): Int { return 1 }
`
	l := candy_lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 9 {
		t.Fatalf("expected 9 statements, got %d", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, diag := range errors {
		t.Errorf("parser error at %d:%d: %q", diag.Line, diag.Col, diag.Message)
	}
	t.FailNow()
}
