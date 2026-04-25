package candy_evaluator

import "candy/candy_ast"

// Env is a name→value map with an optional parent (closure chain).
type Env struct {
	Store    map[string]*Value
	Parent   *Env
	Cwd      string
	Imported map[string]bool
	Defers   []*candy_ast.DeferStatement
}

// NewEnclosed makes a new env linked to the parent.
func (e *Env) NewEnclosed() *Env {
	return &Env{Store: make(map[string]*Value), Parent: e, Cwd: e.Cwd, Imported: e.Imported}
}

// Get a binding.
func (e *Env) Get(name string) (*Value, bool) {
	if e == nil {
		return nil, false
	}
	if v, ok := e.Store[name]; ok {
		return v, true
	}
	return e.Parent.Get(name)
}

// Set a binding in this environment only.
func (e *Env) Set(name string, v *Value) {
	if e.Store == nil {
		e.Store = make(map[string]*Value)
	}
	e.Store[name] = v
}

// Update searches parent scopes for an existing binding and updates it.
// Returns false if the binding was not found in any scope.
func (e *Env) Update(name string, v *Value) bool {
	if e == nil {
		return false
	}
	if _, ok := e.Store[name]; ok {
		e.Store[name] = v
		return true
	}
	if e.Parent != nil {
		return e.Parent.Update(name, v)
	}
	return false
}

// BindFunc registers a K-Go function in the current env.
func (e *Env) BindFunc(fs *candy_ast.FunctionStatement) {
	built := &Value{Kind: ValFunction, Fn: &functionVal{Stmt: fs, Env: e.NewEnclosed(), Outer: e}}
	built.Fn.Env.Set(fs.Name.Value, built)
}
