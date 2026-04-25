package candy_evaluator

import "os"

// ReplEnv returns an environment prepared for a multi-line session (REPL), with prelude
// bindings, working directory, and an import de-dupe map. Reuse the same *Env for each
// Eval(…, e) to keep top-level state across input lines.
func ReplEnv() *Env {
	e := &Env{Store: make(map[string]*Value), Imported: make(map[string]bool)}
	if wd, err := os.Getwd(); err == nil {
		e.Cwd = wd
	}
	registerPrelude(e)
	return e
}
