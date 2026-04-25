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
	"rand":   `// Alias of random`,
	"time":   `// Host stdlib: time.*`,
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
// Placeholder for interpreter-side IO bindings.
`,
	"std/time": `
// Placeholder for interpreter-side time bindings.
`,
	"std/json": `
// Placeholder for interpreter-side json bindings.
`,
	"std/collections": `
// Kotlin-like collections surface (placeholder).
`,
	"std/concurrent": `
// Coroutine scheduler/channel helpers (placeholder).
`,
	"std/result": `
fun isOk(r: Any): Bool { return true; };
`,
	"std/ffi": `
// Native interop helpers (placeholder).
`,
}

// Lookup returns stdlib source for a module path.
func Lookup(path string) (string, bool) {
	v, ok := Modules[path]
	return v, ok
}
