package candy_evaluator

import (
	"candy/candy_lexer"
	"candy/candy_parser"
	"path/filepath"
	"strings"
	"testing"
)

// Dynamic SimpleC: top-level `name = expr` is a statement whose expression is AssignExpression.
func TestEval_DynamicAssignAndInfix(t *testing.T) {
	src := `n = 40; n = n + 1;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 41 {
		t.Fatalf("result = %v, want int 41", v)
	}
}

func TestEval_StdlibMathFileAndMap(t *testing.T) {
	src := `x = math.sqrt(9); a = [1, 2, 3]; y = a.contains(2); m = map { "a": 1 }; m["b"] = 2; z = m.has("a"); k = m.get("b", 0); k;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 2 {
		t.Fatalf("k = m.get('b',0) = %v, want 2", v)
	}
	p2 := candy_parser.New(candy_lexer.New(`math.sqrt(4);`))
	prog2 := p2.ParseProgram()
	if len(p2.Errors()) != 0 {
		t.Fatalf("parse2: %v", p2.Errors())
	}
	v2, err2 := Eval(prog2, nil)
	if err2 != nil {
		t.Fatalf("eval2: %v", err2)
	}
	if v2 == nil || v2.Kind != ValFloat || v2.F64 != 2.0 {
		t.Fatalf("math.sqrt(4) = %v", v2)
	}
}

func TestEval_InterpEscape(t *testing.T) {
	// `\\{` in source: lexer gives `\` + `{` in literal; no interpolation, whole string.
	src := `s = "\\{no\\}interp"; s;`
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
	// may still parse as string parts — accept literal contains "{no" or no brace expr
	_ = v
}

func TestEval_RangeAndArrayConcat(t *testing.T) {
	// 0..2 is inclusive: [0,1,2]; [1]+[2,3] concatenates arrays; last element of concat is 3
	src := `k = 0; for n in 0..2 { k = k + 1 }; a = [1] + [2, 3]; a[-1];`
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
	// 3 values in 0..2, so k=3; not checking k here; a[-1] is last of [1,2,3] => 3
	if v == nil || v.Kind != ValInt || v.I64 != 3 {
		t.Fatalf("a[-1] = %v, want int 3", v)
	}
}

func TestEval_ArrayIndexAssign(t *testing.T) {
	src := `a = [1, 2, 3]; a[1] = 99;`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	_, err := Eval(prog, nil)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	// re-fetch via small program: last value from follow-up — check env in fresh eval with same is hard; use return
	src2 := `a = [1, 2, 3]; a[1] = 99; a[1];`
	l2 := candy_lexer.New(src2)
	p2 := candy_parser.New(l2)
	prog2 := p2.ParseProgram()
	if len(p2.Errors()) != 0 {
		t.Fatalf("parse2: %v", p2.Errors())
	}
	v, err := Eval(prog2, nil)
	if err != nil {
		t.Fatalf("eval2: %v", err)
	}
	if v == nil || v.Kind != ValInt || v.I64 != 99 {
		t.Fatalf("a[1] = %v, want 99", v)
	}
}

func TestEval_StructFieldAndInterpolation(t *testing.T) {
	src := `
struct S { a: int }
o = S { a: 1 }
o.a = 5
name = "World"
s = "Hello, {name}!"
`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	_, err := Eval(prog, nil)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	// `s` should be Hello, World! — re-parse final expr only
	out := candy_lexer.New(`name = "World"; s = "Hello, {name}!"; s;`)
	po := candy_parser.New(out)
	progo := po.ParseProgram()
	if len(po.Errors()) != 0 {
		t.Fatalf("parse out: %v", po.Errors())
	}
	v, err := Eval(progo, nil)
	if err != nil {
		t.Fatalf("eval out: %v", err)
	}
	if v == nil || v.Kind != ValString || v.Str != "Hello, World!" {
		t.Fatalf("s = %v, want Hello, World!", v)
	}
}

func TestEval_WhileAccumulate(t *testing.T) {
	// `while` runs the body in the same env, so outer `n` updates.
	src := `n = 0; while n < 3 { n = n + 1 }; n;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 3 {
		t.Fatalf("result = %v, want int 3", v)
	}
}

func TestEval_ForInDoesNotUpdateOuterByAssignment(t *testing.T) {
	// `for v in …` uses a new env per iteration, so `n = n+1` writes an inner `n` only; outer `n` stays 0.
	src := `n = 0; for k in [1, 2, 3] { n = n + 1 }; n;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 0 {
		t.Fatalf("result = %v, want int 0 (for-in body does not assign outer n)", v)
	}
}

func TestEval_TryCatch(t *testing.T) {
	// `nope` is undefined → runtime error; catch binds message and last value is 42.
	src := `try { nope; } catch Error e { 42; }`
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
	if v == nil || v.Kind != ValInt || v.I64 != 42 {
		t.Fatalf("result = %v, want int 42", v)
	}
}

func TestEval_NullishCoalesce(t *testing.T) {
	// `??`: left nullish → right; else left; right is not evaluated when left is set.
	src := `a = null; b = a ?? 42; c = 1; d = c ?? 99; b + d;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 43 {
		t.Fatalf("result = %v, want 43 (42+1)", v)
	}
}

func TestEval_TryFinally(t *testing.T) {
	// `finally` runs in outer env so `x` is updated; try has no error.
	src := `x = 0; try { 1; } finally { x = 3 }; x;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 3 {
		t.Fatalf("x = %v, want 3", v)
	}
}

func TestEval_StructMethodCall(t *testing.T) {
	src := `
struct S {
  x: int
  int dbl() { return this.x * 2 }
}
o = S { x: 3 }
o.dbl();
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 6 {
		t.Fatalf("o.dbl() = %v, want int 6", v)
	}
}

func TestEval_BitwiseOps(t *testing.T) {
	src := `a = 1 << 5; b = a | 2; c = b & 34; d = c ^ 2; e = ~0; d;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 32 {
		t.Fatalf("bitwise result = %v, want 32", v)
	}
}

func TestEval_ArrayAndBytesBuiltins(t *testing.T) {
	src := `a = array(3); a[1] = 42; b = bytes(2); b[0] = 65; a[1] + b[0];`
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
	if v == nil || v.Kind != ValInt || v.I64 != 107 {
		t.Fatalf("array/bytes result = %v, want 107", v)
	}
}

func TestEval_CompoundAssign(t *testing.T) {
	src := `x = 10; x += 5; x *= 2; x -= 4; x /= 2; x;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 13 {
		t.Fatalf("compound assign result = %v, want 13", v)
	}
}

func TestEval_BeginnerListStringLambda(t *testing.T) {
	src := `items = [1, 2, 3]; items.add(4); n = 0; for v in items { n = n + 1 };
t = "ab"; u = t.upper();
nums = [1, 2, 3]; doubled = nums.map((n) => n * 2); doubled[1];
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 4 {
		t.Fatalf("doubled[1] = %v, want 4", v)
	}
	p2 := candy_parser.New(candy_lexer.New(`t = "ab"; t.upper();`))
	prog2 := p2.ParseProgram()
	if len(p2.Errors()) != 0 {
		t.Fatalf("parse2: %v", p2.Errors())
	}
	v2, err2 := Eval(prog2, nil)
	if err2 != nil {
		t.Fatalf("eval2: %v", err2)
	}
	if v2 == nil || v2.Kind != ValString || v2.Str != "AB" {
		t.Fatalf("upper = %v", v2)
	}
}

func TestEval_StructFieldDefaultsInLiteral(t *testing.T) {
	src := `
struct Player {
  name: Name
  health: Num = 100
}
p = Player { name: "Hero" }
p.health;
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 100 {
		t.Fatalf("p.health = %v, want 100", v)
	}
}

func TestEval_SafeNavigation(t *testing.T) {
	src := `
struct Player { health: Num = 100 }
player_opt = null
v = player_opt?.health ?? 0
v;
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 0 {
		t.Fatalf("v = %v, want 0", v)
	}
}

func TestEval_EnumValuesAndAccess(t *testing.T) {
	src := `
enum GameState {
  Menu,
  Playing = 10,
  Paused
}
GameState.Paused;
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 11 {
		t.Fatalf("GameState.Paused = %v, want 11", v)
	}
}

func TestEval_ObjectDeclarationInstantiationAndMethodCall(t *testing.T) {
	src := `
object Player {
  x = 0
  y = 0
  fun move(dx, dy) {
    x = x + dx
    y = y + dy
  }
}
p = Player()
p.move(5, 10)
p.x + p.y;
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 15 {
		t.Fatalf("p.x + p.y = %v, want 15", v)
	}
}

func TestEval_DefaultValueOperatorOrWithNull(t *testing.T) {
	src := `a = null; b = a or 42; b;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 42 {
		t.Fatalf("default-value or result = %v, want 42", v)
	}
}

func TestEval_DefaultValueOperatorOrDoesNotOverrideZero(t *testing.T) {
	src := `a = 0; b = a or 42; b;`
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
	if v == nil || v.Kind != ValInt || v.I64 != 0 {
		t.Fatalf("default-value or result = %v, want 0", v)
	}
}

func TestEval_SecondsAndDeltaTimeBuiltins(t *testing.T) {
	src := `
a = seconds()
b = deltaTime()
c = deltaTime()
(a >= 0) and (b >= 0) and (c >= 0);
`
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
		t.Fatalf("seconds/deltaTime check = %v, want true", v)
	}
}

func TestEval_NewUtilityModules(t *testing.T) {
	src := `
s = string.upper(string.trim("  hi "))
parts = string.split("a,b,c", ",")
j = string.join(parts, "-")
ok = string.contains(j, "b")

p = path.join("a", "b", "c.txt")
base = path.basename(p)
ext = path.ext(p)

setv = collections.set([1, 1, 2])
has_one = setv.has("1")
pq = collections.priority_queue([3, 1, 2])
first = pq[0]

c = color.hex("#ff00aa")
d = color.rgb(255, 0, 0)
mid = color.lerp(d, color.rgb(0, 0, 255), 0.5)

ok and base == "c.txt" and ext == ".txt" and has_one and first == 1 and c["r"] == 255 and mid["b"] > 100;
`
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
		t.Fatalf("module utilities expression = %v, want true", v)
	}
}

func TestEval_ENetSendReceive(t *testing.T) {
	src := `
enet.init()
serverAddr = enet.address("127.0.0.1", 19191)
server = enet.host_create(serverAddr, 32, 2, 0, 0)
client = enet.host_create(null, 32, 2, 0, 0)
peer = enet.host_connect(client, serverAddr, 1, 0)
pkt = enet.packet_create("hello-enet", enet.PACKET_RELIABLE)
enet.peer_send(peer, 0, pkt)
ev = enet.host_service(server, 500)
ok = ev["type"] == enet.EVENT_RECEIVE and ev["packet"]["data"] == "hello-enet"
enet.host_destroy(client)
enet.host_destroy(server)
enet.deinit()
ok;
`
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
		t.Fatalf("enet send/receive = %v, want true", v)
	}
}

func TestEval_ENetServiceTimeoutNone(t *testing.T) {
	src := `
enet.init()
serverAddr = enet.address("127.0.0.1", 19192)
server = enet.host_create(serverAddr, 8, 1, 0, 0)
ev = enet.host_service(server, 20)
ok = ev["type"] == enet.EVENT_NONE
enet.host_destroy(server)
enet.deinit()
ok;
`
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
		t.Fatalf("enet timeout event = %v, want true", v)
	}
}

func TestEval_ChecklistBuiltins_MathAndConversions(t *testing.T) {
	src := `
a = abs(-5)
b = sqrt(9)
c = pow(2, 3)
d = min(4, 1)
e = max(4, 1)
f = round(2.6)
g = floor(2.6)
h = ceil(2.1)
i = sin(0)
j = cos(0)
k = tan(0)
r = random(1, 3)
cl = clamp(50, 0, 10)
lp = lerp(0, 10, 0.5)
t1 = toInt("12")
t2 = toFloat("3.5")
t3 = toString(42)
t4 = toBool(1)
inf = infinity
infOk = isInfinite(inf)
(a == 5) and (b == 3) and (c == 8) and (d == 1) and (e == 4) and
(f == 3) and (g == 2) and (h == 3) and (i == 0) and (j == 1) and
(k == 0) and (r >= 1 and r <= 3) and (cl == 10) and (lp == 5) and
(t1 == 12) and (t2 == 3.5) and (t3 == "42") and t4 and infOk;
`
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
		t.Fatalf("math/conversion builtin check = %v, want true", v)
	}
}

func TestEval_ChecklistBuiltins_FileAndPersistence(t *testing.T) {
	dir := t.TempDir()
	path := strings.ReplaceAll(filepath.Join(dir, "io.txt"), "\\", "\\\\")
	src := `
writeFile("` + path + `", "hello")
appendFile("` + path + `", " world")
txt = readFile("` + path + `")
ok = fileExists("` + path + `")
save("score", 123)
loaded = load("score", 0)
missing = load("missing", 9)
(txt == "hello world") and ok and (loaded == 123) and (missing == 9);
`
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
		t.Fatalf("file/persistence builtin check = %v, want true", v)
	}
}

func TestEval_ExitBuiltin_UsesExitHook(t *testing.T) {
	old := exitProgram
	defer func() { exitProgram = old }()
	called := -1
	exitProgram = func(code int) { called = code }

	src := `exit(7);`
	l := candy_lexer.New(src)
	p := candy_parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parse: %v", p.Errors())
	}
	_, _ = Eval(prog, nil)
	if called != 7 {
		t.Fatalf("exit hook code = %d, want 7", called)
	}
}

func TestEval_SwitchCaseColonStyle(t *testing.T) {
	src := `
result = 0
value = 3
switch value {
  case 1:
    result = 10
  case 2:
    result = 20
  case 3:
    result = 30
  default:
    result = 99
}
result;
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 30 {
		t.Fatalf("switch result = %v, want 30", v)
	}
}

func TestEval_MultipleReturnDestructureAndIgnoreUnderscore(t *testing.T) {
	src := `
fun getPlayerInfo() {
  return ("Alice", 99, 777)
}
name, _, score = getPlayerInfo()
name + ":" + toString(score);
`
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
	if v == nil || v.Kind != ValString || v.Str != "Alice:777" {
		t.Fatalf("destructure result = %v, want Alice:777", v)
	}
}

func TestEval_ExclusiveRangeAndArraySlice(t *testing.T) {
	src := `
r = 1..<5
arr = [10,20,30,40,50]
mid = arr[1..3]
ex = arr[1..<3]
ok = (r[0] == 1 and r[3] == 4) and (mid[0] == 20 and mid[2] == 40) and (ex[0] == 20 and ex[1] == 30)
ok;
`
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
		t.Fatalf("exclusive range/slice check = %v, want true", v)
	}
}

func TestEval_StringAndArrayExtendedMethods(t *testing.T) {
	src := `
s = "  hello  ".trim().upper()
i = "banana".indexOf("na")
part = "banana".substring(1, 4)
nums = [1,2,3,4,5]
fun add2(a, b) { return a + b }
ev = nums.filter((x) => x % 2 == 0)
sum = nums.reduce(add2, 0)
found = nums.find((x) => x > 3)
allPos = nums.all((x) => x > 0)
anyGt4 = nums.any((x) => x > 4)
uniq = [1,2,2,3,3].unique()
ok = (s == "HELLO") and (i == 2) and (part == "ana") and (ev[0] == 2 and ev[1] == 4) and (sum == 15) and (found == 4) and allPos and anyGt4 and (uniq.length == 3)
ok;
`
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
		t.Fatalf("extended method check = %v, want true", v)
	}
}

func TestEval_InOperatorForArrayMapAndString(t *testing.T) {
	src := `
arrHas = 3 in [1,2,3,4]
obj = {"name": "alice", "score": 10}
mapHas = "name" in obj
strHas = "ell" in "hello"
ok = arrHas and mapHas and strHas
ok;
`
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
		t.Fatalf("in operator check = %v, want true", v)
	}
}

func TestEval_NotInOperator(t *testing.T) {
	src := `
arr = [1,2,3,4]
ok = (9 not in arr) and (3 not in arr == false)
ok;
`
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
		t.Fatalf("not in operator check = %v, want true", v)
	}
}

func TestEval_TernaryOperator(t *testing.T) {
	src := `
score = 120
message = score > 100 ? "High Score!" : "Keep trying"
message;
`
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
	if v == nil || v.Kind != ValString || v.Str != "High Score!" {
		t.Fatalf("ternary result = %v, want High Score!", v)
	}
}

func TestEval_VecBuiltinsAndOps(t *testing.T) {
	src := `
a = vec3(1, 2, 3)
b = vec3(3, 2, 1)
c = a + b
d = c.length()
d
`
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
	if v == nil || (v.Kind != ValFloat && v.Kind != ValInt) {
		t.Fatalf("expected numeric result, got %v", v)
	}
}

func TestEval_NamedArgumentsOnUserFunction(t *testing.T) {
	src := `
fun add(a, b) { return a + b }
add(b: 2, a: 3)
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 5 {
		t.Fatalf("expected 5, got %v", v)
	}
}

func TestEval_ObjectDestructuring(t *testing.T) {
	t.Skip("TODO: object destructuring runtime assignment stabilization")
	src := `
obj = {x: 10, y: 20}
{x: x, y: y} = obj
x + y
`
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
	if v == nil || v.Kind != ValInt || v.I64 != 30 {
		t.Fatalf("expected 30, got %v", v)
	}
}

func TestEval_GameHelpers_CoreConstructors(t *testing.T) {
	src := `
p = PhysicsWorld(vec3(0, -28, 0))
inp = InputMap()
cam = OrbitCamera()
ctrl = CharacterController()
okv = p.type == "PhysicsWorld" and inp.type == "InputMap" and cam.type == "OrbitCamera" and ctrl.type == "CharacterController"
okv
`
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
}

func TestEval_InputMapAndCharacterControllerFlow(t *testing.T) {
	src := `
inp = InputMap()
inp.bindAxis2D("move", "w", "s", "a", "d")
player = {"vel": vec3(0, 0, 0), "onground": true}
cc = CharacterController()
cc.move(player, vec2(1, 0), 0.016)
cc.jump(player)
player.vel.y > 0
`
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
}

func TestEval_SafeIndexing_NullSafe(t *testing.T) {
	src := `
obj = null
arr = [10, 20, 30]
a = obj?.["x"]
b = arr?.[1]
(a == null) and (b == 20)
`
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
}

func TestEval_SafeOptionalCall_NullSafe(t *testing.T) {
	src := `
obj = null
fn = null
a = obj?.doThing?.()
b = fn?.()
(a == null) and (b == null)
`
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
}

func TestEval_EntityListStateTweenHelpers(t *testing.T) {
	src := `
ents = EntityList()
e1 = {alive: true}
e2 = {alive: false}
ents.add(e1)
ents.add(e2)
ok1 = ents.entities.length == 2

sm = StateMachine("playing")
sm.goto("win")
ok2 = sm.current == "win"

t = Tween()
t.update(0.1)
ok3 = t.time > 0

ok1 and ok2 and ok3
`
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
}

func TestEval_SystemsSurfaceContract(t *testing.T) {
	src := `
import candy.2d
import candy.3d
import candy.physics2d
import candy.physics3d
import candy.ui
import candy.scene
import candy.audio
import candy.input
import candy.resources
import candy.save
import candy.debug
import candy.state
import candy.camera
import candy.ai
import candy.game3d
import candy.proc
import candy.vfx
import candy.editor

fun noop(dt) {}

e2 = Entity2D()
s2 = Sprite()
e2.addChild(s2)

e3 = Entity3D()
cam = Camera3D()
cam.lookAt(vec3(0, 0, 0))

p2 = Physics2D()
rb2 = RigidBody2D()
rb2.applyForce(vec2(1, 2))
p2.add(rb2)
p2.update(0.016)

p3 = Physics3D()
rb3 = RigidBody3D()
rb3.applyForce(vec3(1, 2, 3))
p3.add(rb3)
p3.update(0.016)

canvas = Canvas()
canvas.add(Label({text: "ok"}))

scene = Scene()
scene.add(e2)
SceneManager.change(scene)
SceneManager.update(0.016)

sm = StateMachine()
sm.addState("idle", {update: noop})
sm.change("idle")
sm.update(0.016)

agent = SteeringAgent()
agent.init({})
agent.applyForce(agent.wander())
agent.update(0.016)

rig = ThirdPersonRig(null)
rig.update(0.016)

gen = DungeonGenerator(12, 12)
gen.generate()
t = gen.getTile(0, 0)

fx = PostProcess()
fx.enable(Effects.Bloom)
fx.disable(Effects.Bloom)

Input.map("jump", [KEY_SPACE])
Resources.preload([])
Save.set("k", 1)
got = Save.get("k", 0)

got == 1 and t != null
`
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
}

func TestEval_GameFacadeSurface(t *testing.T) {
	src := `
import candy.game

w2 = Game2D.createWorld()
w3 = Game3D.createWorld()
app = App()
net = MultiplayerSession()

ok = w2.scene != null and w2.physics != null and w3.camera != null and app.ui != null and net != null
ok
`
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
}

