package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func soundByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Sound, error) {
	id, err := argInt(name, args, i)
	if err != nil { return 0, rl.Sound{}, err }
	s, ok := sounds[id]
	if !ok { return 0, rl.Sound{}, fmt.Errorf("%s: invalid sound handle %d", name, id) }
	return id, s, nil
}

func musicByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Music, error) {
	id, err := argInt(name, args, i)
	if err != nil { return 0, rl.Music{}, err }
	m, ok := musics[id]
	if !ok { return 0, rl.Music{}, fmt.Errorf("%s: invalid music handle %d", name, id) }
	return id, m, nil
}

func builtinInitAudioDevice(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.InitAudioDevice()
	return null(), nil
}

func builtinCloseAudioDevice(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	rl.CloseAudioDevice()
	return null(), nil
}

func builtinLoadSound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadSound", args, 1); err != nil { return nil, err }
	path, err := argString("loadSound", args, 0)
	if err != nil { return nil, err }
	s := rl.LoadSound(path)
	id := nextSoundID
	nextSoundID++
	sounds[id] = s
	return vInt(id), nil
}

func builtinUnloadSound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, s, err := soundByID("unloadSound", args, 0)
	if err != nil { return nil, err }
	rl.UnloadSound(s)
	delete(sounds, id)
	return null(), nil
}

func builtinPlaySound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := soundByID("playSound", args, 0)
	if err != nil { return nil, err }
	rl.PlaySound(s)
	return null(), nil
}

func builtinStopSound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, s, err := soundByID("stopSound", args, 0)
    if err != nil { return nil, err }
    rl.StopSound(s)
    return null(), nil
}

func builtinIsSoundPlaying(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, s, err := soundByID("isSoundPlaying", args, 0)
    if err != nil { return nil, err }
    return vBool(rl.IsSoundPlaying(s)), nil
}

func builtinSetSoundVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 2 { return nil, fmt.Errorf("setSoundVolume expects soundId, volume") }
    _, s, err := soundByID("setSoundVolume", args, 0)
    if err != nil { return nil, err }
    vol, _ := getArgFloat("setSoundVolume", args, 1)
    rl.SetSoundVolume(s, float32(vol))
    return null(), nil
}

func builtinLoadMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if err := expectArgs("loadMusicStream", args, 1); err != nil { return nil, err }
    path, err := argString("loadMusicStream", args, 0)
    if err != nil { return nil, err }
    m := rl.LoadMusicStream(path)
    id := nextMusicID
    nextMusicID++
    musics[id] = m
    return vInt(id), nil
}

func builtinUnloadMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    id, m, err := musicByID("unloadMusicStream", args, 0)
    if err != nil { return nil, err }
    rl.UnloadMusicStream(m)
    delete(musics, id)
    return null(), nil
}

func builtinPlayMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, m, err := musicByID("playMusicStream", args, 0)
    if err != nil { return nil, err }
    rl.PlayMusicStream(m)
    return null(), nil
}

func builtinUpdateMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, m, err := musicByID("updateMusicStream", args, 0)
    if err != nil { return nil, err }
    rl.UpdateMusicStream(m)
    return null(), nil
}

func builtinStopMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    _, m, err := musicByID("stopMusicStream", args, 0)
    if err != nil { return nil, err }
    rl.StopMusicStream(m)
    return null(), nil
}

func builtinSetMusicVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
    if len(args) < 2 { return nil, fmt.Errorf("setMusicVolume expects musicId, volume") }
    _, m, err := musicByID("setMusicVolume", args, 0)
    if err != nil { return nil, err }
    vol, _ := getArgFloat("setMusicVolume", args, 1)
    rl.SetMusicVolume(m, float32(vol))
    return null(), nil
}

func builtinBeep(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	fmt.Print("\a")
	return null(), nil
}
