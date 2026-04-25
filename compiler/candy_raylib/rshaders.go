package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func shaderByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Shader, error) {
	id, err := argInt(name, args, i)
	if err != nil { return 0, rl.Shader{}, err }
	s, ok := shaders[id]
	if !ok { return 0, rl.Shader{}, fmt.Errorf("%s: invalid shader handle %d", name, id) }
	return id, s, nil
}

func builtinLoadShader(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadShader", args, 2); err != nil { return nil, err }
	vs, _ := argString("loadShader", args, 0)
	fs, _ := argString("loadShader", args, 1)
	s := rl.LoadShader(vs, fs)
	id := nextShaderID
	nextShaderID++
	shaders[id] = s
	return vInt(id), nil
}

func builtinUnloadShader(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, s, err := shaderByID("unloadShader", args, 0)
	if err != nil { return nil, err }
	rl.UnloadShader(s)
	delete(shaders, id)
	return null(), nil
}

func builtinBeginShaderMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := shaderByID("beginShaderMode", args, 0)
	if err != nil { return nil, err }
	rl.BeginShaderMode(s)
	return null(), nil
}

func builtinEndShaderMode(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.EndShaderMode()
	return null(), nil
}
