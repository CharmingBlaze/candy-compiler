package candy_evaluator

// AllNameBindings collects bound names in this environment chain.
func (e *Env) AllNameBindings() []string {
	seen := make(map[string]struct{})
	var rec func(*Env)
	rec = func(x *Env) {
		if x == nil {
			return
		}
		if x.Parent != nil {
			rec(x.Parent)
		}
		for k := range x.Store {
			seen[k] = struct{}{}
		}
	}
	rec(e)
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
