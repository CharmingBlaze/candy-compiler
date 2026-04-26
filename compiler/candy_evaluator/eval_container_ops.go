package candy_evaluator

import (
	"candy/candy_ast"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func lookupModFn(m *moduleVal, name string) (func([]*Value) (*Value, error), bool) {
	ln := strings.ToLower(name)
	if m.Fns != nil {
		for k, f := range m.Fns {
			if strings.ToLower(k) == ln {
				return f, true
			}
		}
	}
	return nil, false
}

func evalArgs(argExprs []candy_ast.Expression, e *Env) ([]*Value, error) {
	var args []*Value
	for _, a := range argExprs {
		av, err2 := evalExpression(a, e)
		if err2 != nil {
			return nil, err2
		}
		args = append(args, av)
	}
	return args, nil
}

func callArrayMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	elems := recv.Elems
	ln := strings.ToLower(name)
	switch ln {
	case "contains", "include":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "contains: one value"}
		}
		for _, e0 := range elems {
			if valueEqual(&e0, args[0]) {
				return &Value{Kind: ValBool, B: true}, nil
			}
		}
		return &Value{Kind: ValBool, B: false}, nil
	case "index_of", "indexof", "index":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "index_of: one value"}
		}
		for i := range elems {
			if valueEqual(&elems[i], args[0]) {
				return &Value{Kind: ValInt, I64: int64(i)}, nil
			}
		}
		return &Value{Kind: ValInt, I64: -1}, nil
	case "is_empty", "isempty", "empty":
		return &Value{Kind: ValBool, B: len(elems) == 0}, nil
	case "length", "count", "size":
		return &Value{Kind: ValInt, I64: int64(len(elems))}, nil
	case "first", "head":
		if len(elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		return ptrVal(elems[0]), nil
	case "last", "tail":
		if len(elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		return ptrVal(elems[len(elems)-1]), nil
	case "clear", "empty_all":
		recv.Elems = recv.Elems[:0] // in-place: receiver must be shared *Value
		return recv, nil
	case "shuffle":
		rg := rand.New(rand.NewSource(time.Now().UnixNano()))
		ix := make([]int, len(recv.Elems))
		for i := range ix {
			ix[i] = i
		}
		rg.Shuffle(len(ix), func(i, j int) { ix[i], ix[j] = ix[j], ix[i] })
		ne := make([]Value, len(recv.Elems))
		for to, from := range ix {
			ne[to] = recv.Elems[from]
		}
		recv.Elems = ne
		return recv, nil
	case "add", "push", "append":
		for _, a := range args {
			recv.Elems = append(recv.Elems, valueToValue(a))
		}
		return recv, nil
	case "insert":
		if len(args) != 2 {
			return nil, &RuntimeError{Msg: "insert: (index, value)"}
		}
		ix, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "insert: int index"}
		}
		i := int(ix)
		if i < 0 {
			i = len(recv.Elems) + i
		}
		if i < 0 || i > len(recv.Elems) {
			return nil, &RuntimeError{Msg: "insert: index out of range"}
		}
		v := valueToValue(args[1])
		recv.Elems = append(recv.Elems, Value{})
		copy(recv.Elems[i+1:], recv.Elems[i:])
		recv.Elems[i] = v
		return recv, nil
	case "remove", "delete":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove: one argument (index or value)"}
		}
		// If it's an integer, try removing by index first.
		if ix, ok := asInt64Value(args[0]); ok {
			i := int(ix)
			if i < 0 {
				i = len(recv.Elems) + i
			}
			if i >= 0 && i < len(recv.Elems) {
				recv.Elems = append(recv.Elems[:i], recv.Elems[i+1:]...)
				return recv, nil
			}
		}
		// Otherwise (or if index was out of bounds), try removing by value.
		for i := range recv.Elems {
			if valueEqual(&recv.Elems[i], args[0]) {
				recv.Elems = append(recv.Elems[:i], recv.Elems[i+1:]...)
				return recv, nil
			}
		}
		return recv, nil
	case "remove_first", "remove_value", "removevalue":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "removeFirst: one value to delete (first match)"}
		}
		for j := range recv.Elems {
			if valueEqual(&recv.Elems[j], args[0]) {
				recv.Elems = append(recv.Elems[:j], recv.Elems[j+1:]...)
				return recv, nil
			}
		}
		return recv, nil
	case "remove_last", "removelast", "pop":
		if len(recv.Elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		last := recv.Elems[len(recv.Elems)-1]
		recv.Elems = recv.Elems[:len(recv.Elems)-1]
		return ptrVal(last), nil
	case "remove_at", "delete_at", "splice1":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove_at: one index"}
		}
		ix, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "remove_at: int index"}
		}
		i := int(ix)
		if i < 0 {
			i = len(recv.Elems) + i
		}
		if i < 0 || i >= len(recv.Elems) {
			return nil, &RuntimeError{Msg: "remove_at: index out of range"}
		}
		recv.Elems = append(recv.Elems[:i], recv.Elems[i+1:]...)
		return recv, nil
	case "sort":
		sort.Slice(recv.Elems, func(i, j int) bool {
			return valueOrderLessForSort(&recv.Elems[i], &recv.Elems[j])
		})
		return recv, nil
	case "reverse":
		for l, r := 0, len(recv.Elems)-1; l < r; l, r = l+1, r-1 {
			recv.Elems[l], recv.Elems[r] = recv.Elems[r], recv.Elems[l]
		}
		return recv, nil
	case "sum":
		if len(recv.Elems) == 0 {
			return &Value{Kind: ValInt, I64: 0}, nil
		}
		allInt := true
		var isum int64
		var fsum float64
		for i := range recv.Elems {
			v := recv.Elems[i]
			switch v.Kind {
			case ValInt:
				if allInt {
					isum += v.I64
				} else {
					fsum += float64(v.I64)
				}
			case ValFloat:
				if allInt {
					fsum = float64(isum) + v.F64
					isum = 0
					allInt = false
				} else {
					fsum += v.F64
				}
			default:
				return nil, &RuntimeError{Msg: "sum: need numeric elements"}
			}
		}
		if allInt {
			return &Value{Kind: ValInt, I64: isum}, nil
		}
		return &Value{Kind: ValFloat, F64: fsum}, nil
	case "max", "min":
		if len(recv.Elems) == 0 {
			return &Value{Kind: ValNull}, nil
		}
		best := recv.Elems[0]
		isMax := strings.EqualFold(ln, "max")
		for i := 1; i < len(recv.Elems); i++ {
			cur := recv.Elems[i]
			if isMax {
				if valueOrderLessForSort(&best, &cur) {
					best = cur
				}
			} else {
				if valueOrderLessForSort(&cur, &best) {
					best = cur
				}
			}
		}
		return ptrVal(best), nil
	case "join":
		if len(args) < 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "join: separator string"}
		}
		sep := args[0].Str
		var b strings.Builder
		for j, v := range recv.Elems {
			if j > 0 {
				b.WriteString(sep)
			}
			p := ptrVal(v)
			if p == nil {
				b.WriteString("null")
			} else {
				b.WriteString(p.String())
			}
		}
		return &Value{Kind: ValString, Str: b.String()}, nil
	case "map":
		if len(args) != 1 || args[0] == nil {
			return nil, &RuntimeError{Msg: "map: one function"}
		}
		if args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "map: need function (e.g. (x) => x + 1)"}
		}
		out := make([]Value, 0, len(recv.Elems))
		for i := range recv.Elems {
			ev := &recv.Elems[i]
			mv, err := evalUserFunction(args[0], []*Value{ptrVal(*ev)})
			if err != nil {
				return nil, err
			}
			out = append(out, valueToValue(mv))
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "filter":
		if len(args) != 1 || args[0] == nil {
			return nil, &RuntimeError{Msg: "filter: one function"}
		}
		if args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "filter: need function (e.g. (n) => n > 0)"}
		}
		var out []Value
		for i := range recv.Elems {
			ev := &recv.Elems[i]
			fv, err := evalUserFunction(args[0], []*Value{ptrVal(*ev)})
			if err != nil {
				return nil, err
			}
			if fv.Truthy() {
				out = append(out, recv.Elems[i])
			}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "reduce":
		if len(args) != 2 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "reduce: (function, initial)"}
		}
		acc := args[1]
		for i := range recv.Elems {
			nv, err := evalUserFunction(args[0], []*Value{acc, ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			acc = nv
		}
		return acc, nil
	case "find":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "find: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if okv.Truthy() {
				return ptrVal(recv.Elems[i]), nil
			}
		}
		return &Value{Kind: ValNull}, nil
	case "all":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "all: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if !okv.Truthy() {
				return &Value{Kind: ValBool, B: false}, nil
			}
		}
		return &Value{Kind: ValBool, B: true}, nil
	case "any":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValFunction {
			return nil, &RuntimeError{Msg: "any: one predicate function"}
		}
		for i := range recv.Elems {
			okv, err := evalUserFunction(args[0], []*Value{ptrVal(recv.Elems[i])})
			if err != nil {
				return nil, err
			}
			if okv.Truthy() {
				return &Value{Kind: ValBool, B: true}, nil
			}
		}
		return &Value{Kind: ValBool, B: false}, nil
	case "unique":
		out := make([]Value, 0, len(recv.Elems))
		for i := range recv.Elems {
			dup := false
			for j := range out {
				if valueEqual(&out[j], &recv.Elems[i]) {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, recv.Elems[i])
			}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	default:
		return nil, &RuntimeError{Msg: "array has no method: " + name}
	}
}

func callMapMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	if recv.StrMap == nil {
		recv.StrMap = make(map[string]Value)
	}
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	m := recv.StrMap
	ln := strings.ToLower(name)
	switch ln {
	case "addstatic":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "physicsworld") {
			if len(args) != 1 {
				return nil, &RuntimeError{Msg: "addStatic: one collider"}
			}
			st := m["static"]
			st.Elems = append(st.Elems, valueToValue(args[0]))
			m["static"] = st
			return recv, nil
		}
	case "adddynamic":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "physicsworld") {
			if len(args) != 1 {
				return nil, &RuntimeError{Msg: "addDynamic: one body"}
			}
			dy := m["dynamic"]
			dy.Elems = append(dy.Elems, valueToValue(args[0]))
			m["dynamic"] = dy
			return recv, nil
		}
	case "resolvecollision":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "physicsworld") {
			if len(args) < 2 {
				return nil, &RuntimeError{Msg: "resolveCollision: dynamicBox, staticBox [, velVec3]"}
			}
			if args[0] == nil || args[1] == nil || args[0].Kind != ValMap || args[1].Kind != ValMap {
				return nil, &RuntimeError{Msg: "resolveCollision: requires boxes"}
			}
			return aabbOverlapsMap(args[0], args[1])
		}
	case "update":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString {
			switch strings.ToLower(mtype.Str) {
			case "physicsworld":
				dt := 0.016
				if len(args) > 0 && args[0] != nil {
					if f, err := f64Arg(args[0]); err == nil {
						dt = f
					}
				}
				dynamic := m["dynamic"]
				static := m["static"]
				gravity := m["gravity"]
				gVec := []float64{0, -9.8, 0}
				if gravity.Kind == ValVec && len(gravity.Vec) == 3 {
					gVec = gravity.Vec
				}

				for i := range dynamic.Elems {
					body := &dynamic.Elems[i]
					if body.Kind != ValMap {
						continue
					}
					posVal, okP := body.StrMap["position"]
					velVal, okV := body.StrMap["velocity"]
					if !okP || !okV || posVal.Kind != ValVec || velVal.Kind != ValVec {
						continue
					}
					// Apply gravity
					velVal.Vec[0] += gVec[0] * dt
					velVal.Vec[1] += gVec[1] * dt
					velVal.Vec[2] += gVec[2] * dt

					// Integrate position
					posVal.Vec[0] += velVal.Vec[0] * dt
					posVal.Vec[1] += velVal.Vec[1] * dt
					posVal.Vec[2] += velVal.Vec[2] * dt

					// Simple ground collision
					if posVal.Vec[1] < 0 {
						posVal.Vec[1] = 0
						velVal.Vec[1] = 0
					}

					// Check static collisions (AABB)
					for j := range static.Elems {
						stat := &static.Elems[j]
						if stat.Kind != ValMap {
							continue
						}
						
						sAABB, okS := stat.StrMap["aabb"]
						if !okS || sAABB.Kind != ValMap {
							continue
						}

						// Calculate player AABB (approximate)
						bAABB := &Value{
							Kind: ValMap,
							StrMap: map[string]Value{
								"center":      {Kind: ValVec, Vec: []float64{posVal.Vec[0], posVal.Vec[1] + 1.0, posVal.Vec[2]}},
								"halfExtents": {Kind: ValVec, Vec: []float64{0.5, 1.0, 0.5}},
							},
						}

						if overlaps, _ := aabbOverlapsMap(bAABB, &sAABB); overlaps.Truthy() {
							// Basic resolution: revert to previous horizontal position
							posVal.Vec[0] -= velVal.Vec[0] * dt
							posVal.Vec[2] -= velVal.Vec[2] * dt
							velVal.Vec[0] = 0
							velVal.Vec[2] = 0
						}
					}

					body.StrMap["position"] = posVal
					body.StrMap["velocity"] = velVal
				}
				m["dynamic"] = dynamic
				return recv, nil
			case "orbitcamera":
				// Get current state from map
				target := m["target"]
				dist := m["distance"]
				yaw := m["yaw"]
				pitch := m["pitch"]
				sens := m["sensitivity"]
				zoomSpd := m["zoomSpeed"]

				if target.Kind != ValVec || dist.Kind != ValFloat || yaw.Kind != ValFloat || pitch.Kind != ValFloat {
					return recv, nil
				}

				// Input handling via builtins
				if fn, ok := Builtins["ismousebuttondown"]; ok {
					// 0 is left, 1 is right
					mb := Value{Kind: ValInt, I64: 1}
					if down, _ := fn([]*Value{&mb}); down != nil && down.Truthy() {
						if deltaFn, okD := Builtins["getmousedelta"]; okD {
							d, _ := deltaFn(nil)
							if d != nil && d.Kind == ValVec && len(d.Vec) >= 2 {
								yaw.F64 += d.Vec[0] * sens.F64
								pitch.F64 += d.Vec[1] * sens.F64
								// Clamp pitch
								if pitch.F64 > 1.5 {
									pitch.F64 = 1.5
								}
								if pitch.F64 < -1.5 {
									pitch.F64 = -1.5
								}
							}
						}
					}
				}

				if wheelFn, okW := Builtins["getmousewheelmove"]; okW {
					w, _ := wheelFn(nil)
					if w != nil && w.Kind == ValFloat {
						dist.F64 -= w.F64 * zoomSpd.F64
						if dist.F64 < 1 {
							dist.F64 = 1
						}
					}
				}

				// Update map values
				m["yaw"] = yaw
				m["pitch"] = pitch
				m["distance"] = dist

				// Sync with host camera if possible
				// We need to calculate the camera position from target, dist, yaw, pitch
				camPos := []float64{
					target.Vec[0] + dist.F64*math.Cos(yaw.F64)*math.Cos(pitch.F64),
					target.Vec[1] + dist.F64*math.Sin(pitch.F64),
					target.Vec[2] + dist.F64*math.Sin(yaw.F64)*math.Cos(pitch.F64),
				}
				
				if fn, ok := Builtins["setcamera"]; ok {
					// setCamera(pos, target, up)
					p := Value{Kind: ValVec, Vec: camPos}
					up := Value{Kind: ValVec, Vec: []float64{0, 1, 0}}
					_, _ = fn([]*Value{&p, &target, &up})
				}

				return recv, nil
			case "firstpersoncamera":
				return recv, nil
			case "statemachine":
				// optional callback map: states[current].onUpdate
				cur, okCur := m["current"]
				states, okStates := m["states"]
				if okCur && okStates && cur.Kind == ValString && states.Kind == ValMap {
					if sv, okS := states.StrMap[cur.Str]; okS && sv.Kind == ValMap {
						if upd, okU := sv.StrMap["onUpdate"]; okU {
							fn := upd
							_, _ = InvokeCallable(&fn, args)
						}
					}
				}
				return recv, nil
			case "entitylist":
				ents := m["entities"]
				for i := range ents.Elems {
					ev := ents.Elems[i]
					if ev.Kind == ValMap {
						tmp := ev
						_, _ = callMapMethod(&tmp, "update", nil, e)
						ents.Elems[i] = tmp
					}
				}
				m["entities"] = ents
				return recv, nil
			case "tween":
				dt := 0.016
				if len(args) > 0 && args[0] != nil {
					if f, err := f64Arg(args[0]); err == nil {
						dt = f
					}
				}
				tv := m["time"]
				dur := m["duration"]
				dir := m["direction"]
				if tv.Kind != ValFloat || dur.Kind != ValFloat || dir.Kind != ValFloat {
					return recv, nil
				}
				tv.F64 += dt * dir.F64
				if tv.F64 > dur.F64 {
					if lp, okLP := m["loop"]; okLP && lp.Kind == ValBool && lp.B {
						if pp, okPP := m["pingpong"]; okPP && pp.Kind == ValBool && pp.B {
							dir.F64 = -1
							tv.F64 = dur.F64
						} else {
							tv.F64 = 0
						}
					} else {
						tv.F64 = dur.F64
					}
				}
				if tv.F64 < 0 {
					tv.F64 = 0
					dir.F64 = 1
				}
				m["time"] = tv
				m["direction"] = dir
				return recv, nil
			}
		}
	case "bind", "map":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "inputmap") {
			if len(args) != 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValString {
				return nil, &RuntimeError{Msg: "bind(action, keyOrKeys)"}
			}
			// Accept either a single key string (action button) or a 4-key array
			// ([up, down, left, right]) for 2D movement actions.
			if args[1].Kind == ValArray {
				if len(args[1].Elems) >= 4 {
					axes := m["axes2d"]
					if axes.Kind != ValMap || axes.StrMap == nil {
						axes = Value{Kind: ValMap, StrMap: map[string]Value{}}
					}
					axes.StrMap[args[0].Str] = Value{
						Kind: ValArray,
						Elems: []Value{
							args[1].Elems[0],
							args[1].Elems[1],
							args[1].Elems[2],
							args[1].Elems[3],
						},
					}
					m["axes2d"] = axes
					return recv, nil
				}
				if len(args[1].Elems) >= 1 {
					acts := m["actions"]
					if acts.Kind != ValMap || acts.StrMap == nil {
						acts = Value{Kind: ValMap, StrMap: map[string]Value{}}
					}
					acts.StrMap[args[0].Str] = valueToValue(args[1])
					m["actions"] = acts
					return recv, nil
				}
				return recv, nil
			}
			if args[1].Kind != ValString && args[1].Kind != ValInt && args[1].Kind != ValFloat {
				return nil, &RuntimeError{Msg: "bind(action, keyOrKeys): key must be int, float, string, or key array"}
			}
			acts := m["actions"]
			if acts.Kind != ValMap || acts.StrMap == nil {
				acts = Value{Kind: ValMap, StrMap: map[string]Value{}}
			}
			acts.StrMap[args[0].Str] = *args[1]
			m["actions"] = acts
			return recv, nil
		}
	case "bindaxis2d":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "inputmap") {
			if len(args) != 5 || args[0] == nil || args[0].Kind != ValString {
				return nil, &RuntimeError{Msg: "bindAxis2D(name, up, down, left, right)"}
			}
			axes := m["axes2d"]
			if axes.Kind != ValMap || axes.StrMap == nil {
				axes = Value{Kind: ValMap, StrMap: map[string]Value{}}
			}
			axes.StrMap[args[0].Str] = Value{Kind: ValArray, Elems: []Value{valueToValue(args[1]), valueToValue(args[2]), valueToValue(args[3]), valueToValue(args[4])}}
			m["axes2d"] = axes
			return recv, nil
		}
	case "getaxis2d", "get2daxis":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "inputmap") {
			if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
				return nil, &RuntimeError{Msg: "getAxis2D(name)"}
			}
			axes := m["axes2d"]
			if axes.Kind != ValMap || axes.StrMap == nil {
				return &Value{Kind: ValVec, Vec: []float64{0, 0}}, nil
			}
			b, okB := axes.StrMap[args[0].Str]
			if !okB || b.Kind != ValArray || len(b.Elems) < 4 {
				return &Value{Kind: ValVec, Vec: []float64{0, 0}}, nil
			}
			getDown := func(keyVal Value) bool {
				if keyVal.Kind != ValString && keyVal.Kind != ValInt && keyVal.Kind != ValFloat {
					return false
				}
				if fn, ok := Builtins["iskeydown"]; ok {
					kv := keyVal
					v, err := fn([]*Value{&kv})
					return err == nil && v != nil && v.Truthy()
				}
				return false
			}
			x := 0.0
			y := 0.0
			if getDown(b.Elems[0]) {
				y -= 1
			}
			if getDown(b.Elems[1]) {
				y += 1
			}
			if getDown(b.Elems[2]) {
				x -= 1
			}
			if getDown(b.Elems[3]) {
				x += 1
			}
			return &Value{Kind: ValVec, Vec: []float64{x, y}}, nil
		}
	case "ispressed", "isdown", "justpressed", "justreleased":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "inputmap") {
			if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
				return nil, &RuntimeError{Msg: "isPressed(action) / justPressed(action) / justReleased(action)"}
			}
			acts := m["actions"]
			if acts.Kind != ValMap || acts.StrMap == nil {
				return &Value{Kind: ValBool, B: false}, nil
			}
			k, okK := acts.StrMap[args[0].Str]
			if !okK {
				return &Value{Kind: ValBool, B: false}, nil
			}
			bname := "iskeypressed"
			switch ln {
			case "isdown":
				bname = "iskeydown"
			case "justreleased":
				bname = "iskeyreleased"
			}
			if k.Kind == ValArray {
				for i := range k.Elems {
					kv := k.Elems[i]
					if kv.Kind != ValString && kv.Kind != ValInt && kv.Kind != ValFloat {
						continue
					}
					if fn, ok := Builtins[bname]; ok {
						v, err := fn([]*Value{&kv})
						if err == nil && v != nil && v.Truthy() {
							return &Value{Kind: ValBool, B: true}, nil
						}
					}
				}
				return &Value{Kind: ValBool, B: false}, nil
			}
			if k.Kind != ValString && k.Kind != ValInt && k.Kind != ValFloat {
				return &Value{Kind: ValBool, B: false}, nil
			}
			if fn, ok := Builtins[bname]; ok {
				kk := k
				return fn([]*Value{&kk})
			}
			return &Value{Kind: ValBool, B: false}, nil
		}
	case "move":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "charactercontroller") {
			// move(playerMap, inputVec2, dt)
			if len(args) != 3 || args[0] == nil || args[1] == nil || args[2] == nil || args[0].Kind != ValMap || args[1].Kind != ValVec || len(args[1].Vec) < 2 {
				return nil, &RuntimeError{Msg: "move(player, inputDirVec2, dt)"}
			}
			player := args[0]
			dt, err := f64Arg(args[2])
			if err != nil {
				return nil, err
			}
			acc := m["acceleration"].F64
			drag := m["drag"].F64
			maxSpeed := m["maxSpeed"].F64
			vel, okVel := player.StrMap["vel"]
			if !okVel || vel.Kind != ValVec || len(vel.Vec) < 3 {
				vel = Value{Kind: ValVec, Vec: []float64{0, 0, 0}}
			}
			vel.Vec[0] += args[1].Vec[0] * acc * dt
			vel.Vec[2] += args[1].Vec[1] * acc * dt
			vel.Vec[0] *= (1.0 - clampFloat(drag*dt, 0, 0.9))
			vel.Vec[2] *= (1.0 - clampFloat(drag*dt, 0, 0.9))
			vel.Vec[0] = clampFloat(vel.Vec[0], -maxSpeed, maxSpeed)
			vel.Vec[2] = clampFloat(vel.Vec[2], -maxSpeed, maxSpeed)
			player.StrMap["vel"] = vel
			return player, nil
		}
	case "jump":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "charactercontroller") {
			if len(args) != 1 || args[0] == nil || args[0].Kind != ValMap {
				return nil, &RuntimeError{Msg: "jump(player)"}
			}
			player := args[0]
			vel, okVel := player.StrMap["vel"]
			if !okVel || vel.Kind != ValVec || len(vel.Vec) < 3 {
				vel = Value{Kind: ValVec, Vec: []float64{0, 0, 0}}
			}
			vel.Vec[1] = m["jumpPower"].F64
			player.StrMap["vel"] = vel
			player.StrMap["onGround"] = Value{Kind: ValBool, B: false}
			return player, nil
		}
	case "add":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "entitylist") {
			if len(args) != 1 {
				return nil, &RuntimeError{Msg: "EntityList.add(entity)"}
			}
			es := m["entities"]
			es.Elems = append(es.Elems, valueToValue(args[0]))
			m["entities"] = es
			return recv, nil
		}
	case "where", "filter":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "entitylist") {
			if len(args) != 1 || args[0] == nil {
				return nil, &RuntimeError{Msg: "where(filterFn)"}
			}
			es := m["entities"]
			out := make([]Value, 0, len(es.Elems))
			for i := range es.Elems {
				it := es.Elems[i]
				okv, _ := InvokeCallable(args[0], []*Value{ptrVal(it)})
				if okv != nil && okv.Truthy() {
					out = append(out, it)
				}
			}
			return &Value{Kind: ValMap, StrMap: map[string]Value{
				"type":     {Kind: ValString, Str: "EntityList"},
				"entities": {Kind: ValArray, Elems: out},
			}}, nil
		}
	case "all", "any":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "entitylist") {
			if len(args) != 1 || args[0] == nil {
				return nil, &RuntimeError{Msg: "all/any(predicateFn)"}
			}
			es := m["entities"]
			wantAll := ln == "all"
			result := wantAll
			for i := range es.Elems {
				it := es.Elems[i]
				okv, _ := InvokeCallable(args[0], []*Value{ptrVal(it)})
				tr := okv != nil && okv.Truthy()
				if wantAll && !tr {
					result = false
					break
				}
				if !wantAll && tr {
					result = true
					break
				}
				if !wantAll {
					result = false
				}
			}
			return &Value{Kind: ValBool, B: result}, nil
		}
	case "draw":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString {
			switch strings.ToLower(mtype.Str) {
			case "entitylist":
				es := m["entities"]
				for i := range es.Elems {
					it := es.Elems[i]
					if it.Kind == ValMap {
						tmp := it
						_, _ = callMapMethod(&tmp, "draw", nil, e)
					}
				}
				return recv, nil
			case "uilayout", "hud":
				// hud.update(dt) could handle animations or auto-layouts
				return recv, nil
			}
		}
	case "text", "topleft", "topcenter", "center", "bottomleft":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString {
			if strings.EqualFold(mtype.Str, "uilayout") || strings.EqualFold(mtype.Str, "hud") {
				if len(args) >= 1 && args[0] != nil {
					txt := args[0]
					x := m["padding"]
					y := m["padding"]
					if x.Kind != ValInt { x = Value{Kind: ValInt, I64: 20} }
					if y.Kind != ValInt { y = Value{Kind: ValInt, I64: 20} }
					
					if fn, ok := Builtins["drawtext"]; ok {
						size := Value{Kind: ValInt, I64: 24}
						col := Value{Kind: ValString, Str: "white"}
						_, _ = fn([]*Value{txt, &x, &y, &size, &col})
					}
				}
				return recv, nil
			}
		}
	case "goto":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "statemachine") {
			if len(args) == 1 && args[0] != nil && args[0].Kind == ValString {
				m["current"] = *args[0]
			}
			return recv, nil
		}
	case "state":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "statemachine") {
			// state(name, mapWithCallbacks)
			if len(args) == 2 && args[0] != nil && args[0].Kind == ValString && args[1] != nil && args[1].Kind == ValMap {
				st := m["states"]
				if st.Kind != ValMap || st.StrMap == nil {
					st = Value{Kind: ValMap, StrMap: map[string]Value{}}
				}
				st.StrMap[args[0].Str] = valueToValue(args[1])
				m["states"] = st
			}
			return recv, nil
		}
	case "save", "restore", "reset":
		if len(args) >= 0 {
			if ln == "save" {
				return deepCloneValue(recv), nil
			}
			if ln == "restore" && len(args) == 1 && args[0] != nil {
				clone := deepCloneValue(args[0])
				*recv = *clone
				return recv, nil
			}
			if ln == "reset" {
				if initv, ok := m["_initial"]; ok {
					cv := initv
					*recv = cv
				}
				return recv, nil
			}
		}
	case "animate", "patrol":
		if len(args) >= 0 {
			// marker no-op behavior for now
			return recv, nil
		}
	case "rotatevector":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "transform") {
			if len(args) != 1 || args[0] == nil || args[0].Kind != ValVec {
				return nil, &RuntimeError{Msg: "rotateVector(vec)"}
			}
			ang := m["rotation"]
			if ang.Kind != ValFloat {
				return args[0], nil
			}
			if len(args[0].Vec) == 2 {
				x := args[0].Vec[0]
				y := args[0].Vec[1]
				c := mathCos(ang.F64)
				s := mathSin(ang.F64)
				return &Value{Kind: ValVec, Vec: []float64{x*c - y*s, x*s + y*c}}, nil
			}
			return args[0], nil
		}
	case "overlaps":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString {
			switch strings.ToLower(mtype.Str) {
			case "aabb":
				if len(args) != 1 || args[0] == nil || args[0].Kind != ValMap {
					return nil, &RuntimeError{Msg: "AABB.overlaps expects one AABB"}
				}
				return aabbOverlapsMap(recv, args[0])
			case "sphere":
				if len(args) != 1 || args[0] == nil || args[0].Kind != ValMap {
					return nil, &RuntimeError{Msg: "Sphere.overlaps expects one Sphere"}
				}
				return sphereOverlapsMap(recv, args[0])
			}
		}
		return nil, &RuntimeError{Msg: "map.overlaps is only available on physics types"}
	case "intersects":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "ray") {
			if len(args) != 1 || args[0] == nil || args[0].Kind != ValMap {
				return nil, &RuntimeError{Msg: "Ray.intersects expects one AABB"}
			}
			return rayIntersectsAABB(recv, args[0])
		}
		return nil, &RuntimeError{Msg: "map.intersects is only available on Ray"}
	case "keys", "key_list":
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		out := make([]Value, len(ks))
		for i, k := range ks {
			out[i] = Value{Kind: ValString, Str: k}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "values", "vals":
		// stable order: sort by key
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		out := make([]Value, len(ks))
		for i, k := range ks {
			out[i] = m[k]
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "has", "contains", "include":
		if mtype, ok := m["type"]; ok && mtype.Kind == ValString && strings.EqualFold(mtype.Str, "aabb") {
			if len(args) == 1 && args[0] != nil && args[0].Kind == ValVec {
				return aabbContainsPoint(recv, args[0])
			}
		}
		if len(args) < 1 {
			return nil, &RuntimeError{Msg: "has: one key"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		_, ok := m[ks]
		if !ok {
			// case-insensitive
			for km := range m {
				if strings.EqualFold(km, ks) {
					ok = true
					break
				}
			}
		}
		return &Value{Kind: ValBool, B: ok}, nil
	case "get", "at":
		if len(args) < 1 {
			return nil, &RuntimeError{Msg: "get: key, optional default"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		if v, ok := m[ks]; ok {
			return ptrVal(v), nil
		}
		for km, v := range m {
			if strings.EqualFold(km, ks) {
				return ptrVal(v), nil
			}
		}
		if len(args) >= 2 {
			return args[1], nil
		}
		return &Value{Kind: ValNull}, nil
	case "remove", "delete", "rm":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "remove: one key"}
		}
		ks, e1 := mapKeyString(args[0])
		if e1 != nil {
			return nil, e1
		}
		if _, ok := m[ks]; ok {
			delete(m, ks)
			return recv, nil
		}
		for k := range m {
			if strings.EqualFold(k, ks) {
				delete(m, k)
				return recv, nil
			}
		}
		return recv, nil
	case "clear", "empty":
		for k := range m {
			delete(m, k)
		}
		return recv, nil
	case "merge", "put_all", "putall":
		if len(args) != 1 {
			return nil, &RuntimeError{Msg: "merge: one map"}
		}
		if args[0] == nil || args[0].Kind != ValMap || args[0].StrMap == nil {
			return nil, &RuntimeError{Msg: "merge: need map"}
		}
		for k, v := range args[0].StrMap {
			m[k] = v
		}
		return recv, nil
	default:
		return nil, &RuntimeError{Msg: "map has no method: " + name}
	}
	return nil, &RuntimeError{Msg: "map has no method: " + name}
}

// asInt64Value coerces a value to int index (for list insert/remove).
func asInt64Value(v *Value) (int64, bool) {
	if v == nil {
		return 0, true
	}
	if v.Kind == ValInt {
		return v.I64, true
	}
	if v.Kind == ValFloat {
		return int64(v.F64), true
	}
	return 0, false
}

// valueOrderLessForSort is a total order for sort / max / min: numeric when possible, else String().
func valueOrderLessForSort(a, b *Value) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	lf, rf, useF := toFloats(a, b)
	if (a.Kind == ValInt || a.Kind == ValFloat) && (b.Kind == ValInt || b.Kind == ValFloat) {
		if useF {
			return lf < rf
		}
		if a.Kind == ValInt && b.Kind == ValInt {
			return a.I64 < b.I64
		}
		return lf < rf
	}
	if a.Kind == ValString && b.Kind == ValString {
		return a.Str < b.Str
	}
	return a.String() < b.String()
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func mathCos(v float64) float64 { return math.Cos(v) }
func mathSin(v float64) float64 { return math.Sin(v) }

func callStringMethod(recv *Value, name string, argExprs []candy_ast.Expression, e *Env) (*Value, error) {
	if recv == nil || recv.Kind != ValString {
		return nil, &RuntimeError{Msg: "string method: need string"}
	}
	args, err := evalArgs(argExprs, e)
	if err != nil {
		return nil, err
	}
	s := recv.Str
	ln := strings.ToLower(name)
	switch ln {
	case "upper", "toupper", "to_upper":
		return &Value{Kind: ValString, Str: strings.ToUpper(s)}, nil
	case "lower", "tolower", "to_lower":
		return &Value{Kind: ValString, Str: strings.ToLower(s)}, nil
	case "split":
		sep := " "
		if len(args) >= 1 && args[0] != nil && args[0].Kind == ValString {
			sep = args[0].Str
		}
		n := -1
		if len(args) >= 2 && args[1] != nil && args[1].Kind == ValInt {
			n = int(args[1].I64)
		}
		var parts []string
		if n < 0 {
			parts = strings.Split(s, sep)
		} else {
			parts = strings.SplitN(s, sep, n+1)
		}
		out := make([]Value, len(parts))
		for i, p := range parts {
			out[i] = Value{Kind: ValString, Str: p}
		}
		return &Value{Kind: ValArray, Elems: out}, nil
	case "contains", "includes":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "contains: substring string"}
		}
		return &Value{Kind: ValBool, B: strings.Contains(s, args[0].Str)}, nil
	case "starts_with", "startswith", "prefix":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "starts_with: string prefix"}
		}
		return &Value{Kind: ValBool, B: strings.HasPrefix(s, args[0].Str)}, nil
	case "ends_with", "endswith", "suffix":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "ends_with: string suffix"}
		}
		return &Value{Kind: ValBool, B: strings.HasSuffix(s, args[0].Str)}, nil
	case "trim", "strip", "trim_space":
		return &Value{Kind: ValString, Str: strings.TrimSpace(s)}, nil
	case "replace":
		if len(args) < 2 || args[0] == nil || args[1] == nil || args[0].Kind != ValString || args[1].Kind != ValString {
			return nil, &RuntimeError{Msg: "replace: (old, new) strings"}
		}
		return &Value{Kind: ValString, Str: strings.ReplaceAll(s, args[0].Str, args[1].Str)}, nil
	case "index_of", "indexof":
		if len(args) != 1 || args[0] == nil || args[0].Kind != ValString {
			return nil, &RuntimeError{Msg: "indexOf: substring string"}
		}
		return &Value{Kind: ValInt, I64: int64(strings.Index(s, args[0].Str))}, nil
	case "substring":
		if len(args) < 1 || len(args) > 2 {
			return nil, &RuntimeError{Msg: "substring: start, [end]"}
		}
		start, ok := asInt64Value(args[0])
		if !ok {
			return nil, &RuntimeError{Msg: "substring: start int"}
		}
		runes := []rune(s)
		st := int(start)
		if st < 0 {
			st = len(runes) + st
		}
		if st < 0 {
			st = 0
		}
		if st > len(runes) {
			st = len(runes)
		}
		en := len(runes)
		if len(args) == 2 {
			end, ok := asInt64Value(args[1])
			if !ok {
				return nil, &RuntimeError{Msg: "substring: end int"}
			}
			en = int(end)
			if en < 0 {
				en = len(runes) + en
			}
			if en < 0 {
				en = 0
			}
			if en > len(runes) {
				en = len(runes)
			}
		}
		if en < st {
			en = st
		}
		return &Value{Kind: ValString, Str: string(runes[st:en])}, nil
	default:
		return nil, &RuntimeError{Msg: "string has no method: " + name}
	}
}
