package candy_stdlib

// Modules is the starter stdlib package source map.
// Runtime can import these by module name before disk resolution.
var Modules = map[string]string{
	// Name-only imports (e.g. `import "math";`) — host bindings are injected by registerPrelude.
	"math":   `// Host stdlib: see prelude in candy_evaluator`,
	"file":   `// Host stdlib: file.*, read_file, …`,
	"fs":     `// Alias of file`,
	"json":   `// Host stdlib: json.*`,
	"random": `// Host stdlib: random.*`,
	"time":   `// Host stdlib: time.*`,
	"rand":   `// Alias of random`,
	"candy.math": `
fun sin(x: float): float { return math.sin(x); };
fun cos(x: float): float { return math.cos(x); };
fun tan(x: float): float { return math.tan(x); };
fun atan2(y: float, x: float): float { return math.atan2(y, x); };
fun sqrt(x: float): float { return math.sqrt(x); };
fun floor(x: float): float { return math.floor(x); };
fun ceil(x: float): float { return math.ceil(x); };
fun abs(x: float): float { return math.abs(x); };
fun max(a: float, b: float): float { return math.max(a, b); };
fun min(a: float, b: float): float { return math.min(a, b); };
fun random(min: float, max: float): float { return rand.float(min, max); };
`,
	"std/strings": `
fun isEmpty(s: String): Bool { return len(s) == 0; };
`,
	"std/fs": `
fun readText(path: String): String { return readFile(path); };
fun writeText(path: String, content: String): String? { return writeFile(path, content); };
`,
	"std/path": `
fun join2(a: String, b: String): String { return joinPath(a, b); };
`,
	"std/process": `
fun cwdNow(): String { return cwd(); };
`,
	"std/env": `
fun get(name: String): String? { return getEnv(name); };
`,
	"std/io": `
fun printLine(msg) { print(msg); };
fun printLines(lines) {
    for v in lines { print(v); }
};
`,
	"std/time": `
fun nowMillis() { return time.now(); };
fun sleepMs(ms) { time.sleep(ms); };
fun sleepSec(sec) { time.sleep_sec(sec); };
`,
	"std/json": `
fun parse(text: String) { return json.parse(text); };
fun stringify(v) { return json.stringify(v); };
fun load(path: String) { return json.load(path); };
fun save(path: String, v) { return json.save(path, v); };
`,
	"std/collections": `
fun makeSet(items) {
    var m = {}
    for v in items { m[v] = true; }
    return m
};
fun queue(items) { return items; };
fun stack(items) { return items; };
fun deque(items) { return items; };
`,
	"std/concurrent": `
fun start(fn) { return fn(); };
fun delay(ms, fn) { time.sleep(ms); return fn(); };
`,
	"std/result": `
fun isOk(r: Any): Bool {
    if r == null { return false; }
    if r.ok != null { return r.ok; }
    if r.error != null { return false; }
    return true
};
fun isErr(r: Any): Bool { return !isOk(r); };
`,
	"std/ffi": `
fun available(): Bool { return false; };
fun call(name: String, args) {
    print("ffi.call unavailable in interpreter: {name}")
    return null
};
`,
}

// Lookup returns stdlib source for a module path.
func Lookup(path string) (string, bool) {
	v, ok := Modules[path]
	return v, ok
}
