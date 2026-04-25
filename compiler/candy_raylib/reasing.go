package candy_raylib

import (
	"candy/candy_evaluator"
	"math"
)

// ---- Easing functions (Pure Go implementations of reasings.h) ----

func builtinEaseLinear(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("easeLinear", args, 4); err != nil { return nil, err }
    t, _ := getArgFloat("easeLinear", args, 0)
    b, _ := getArgFloat("easeLinear", args, 1)
    c, _ := getArgFloat("easeLinear", args, 2)
    d, _ := getArgFloat("easeLinear", args, 3)
    return vFloat(c*t/d + b), nil
}

func builtinEaseSineIn(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeSineIn", args)
    return vFloat(-c*math.Cos(t/d*(math.Pi/2)) + c + b), nil
}

func builtinEaseSineOut(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeSineOut", args)
    return vFloat(c*math.Sin(t/d*(math.Pi/2)) + b), nil
}

func builtinEaseSineInOut(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeSineInOut", args)
    return vFloat(-c/2*(math.Cos(math.Pi*t/d)-1) + b), nil
}

// Helpers for easing
func parseEaseArgs(name string, args []*candy_evaluator.Value) (float64, float64, float64, float64, error) {
    if err := expectArgs(name, args, 4); err != nil { return 0, 0, 0, 0, err }
    t, _ := getArgFloat(name, args, 0)
    b, _ := getArgFloat(name, args, 1)
    c, _ := getArgFloat(name, args, 2)
    d, _ := getArgFloat(name, args, 3)
    return t, b, c, d, nil
}

func builtinEaseCubicIn(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeCubicIn", args)
    t /= d
    return vFloat(c*t*t*t + b), nil
}

func builtinEaseCubicOut(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeCubicOut", args)
    t = t/d - 1
    return vFloat(c*(t*t*t+1) + b), nil
}

func builtinEaseBounceOut(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    t, b, c, d, _ := parseEaseArgs("easeBounceOut", args)
    t /= d
    if t < (1 / 2.75) {
        return vFloat(c*(7.5625*t*t) + b), nil
    } else if t < (2 / 2.75) {
        t -= (1.5 / 2.75)
        return vFloat(c*(7.5625*t*t+0.75) + b), nil
    } else if t < (2.5 / 2.75) {
        t -= (2.25 / 2.75)
        return vFloat(c*(7.5625*t*t+0.9375) + b), nil
    } else {
        t -= (2.625 / 2.75)
        return vFloat(c*(7.5625*t*t+0.984375) + b), nil
    }
}
