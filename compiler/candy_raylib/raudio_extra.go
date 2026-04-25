package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- helpers ----

func waveByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.Wave, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return 0, rl.Wave{}, err
	}
	w, ok := waves[id]
	if !ok {
		return 0, rl.Wave{}, fmt.Errorf("%s: invalid wave handle %d", name, id)
	}
	return id, w, nil
}

func audioStreamByID(name string, args []*candy_evaluator.Value, i int) (int64, rl.AudioStream, error) {
	id, err := argInt(name, args, i)
	if err != nil {
		return 0, rl.AudioStream{}, err
	}
	s, ok := audioStreams[id]
	if !ok {
		return 0, rl.AudioStream{}, fmt.Errorf("%s: invalid audioStream handle %d", name, id)
	}
	return id, s, nil
}

// ---- Audio device extras ----

func builtinIsAudioDeviceReady(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vBool(rl.IsAudioDeviceReady()), nil
}

func builtinSetMasterVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setMasterVolume", args, 1); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setMasterVolume", args, 0)
	rl.SetMasterVolume(float32(v))
	return null(), nil
}

func builtinGetMasterVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return vFloat(float64(rl.GetMasterVolume())), nil
}

// ---- Wave loading ----

func builtinLoadWave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadWave", args, 1); err != nil {
		return nil, err
	}
	path, _ := argString("loadWave", args, 0)
	w := rl.LoadWave(path)
	id := nextWaveID
	nextWaveID++
	waves[id] = w
	return vInt(id), nil
}

func builtinLoadWaveFromMemory(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("loadWaveFromMemory expects fileType, bytesArr")
	}
	ft, _ := argString("loadWaveFromMemory", args, 0)
	if args[1] == nil || args[1].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("loadWaveFromMemory: arg 2 must be byte array")
	}
	data := make([]byte, len(args[1].Elems))
	for i, e := range args[1].Elems {
		data[i] = byte(e.I64)
	}
	w := rl.LoadWaveFromMemory(ft, data, int32(len(data)))
	id := nextWaveID
	nextWaveID++
	waves[id] = w
	return vInt(id), nil
}

func builtinIsWaveValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isWaveValid", args, 1); err != nil {
		return nil, err
	}
	_, w, err := waveByID("isWaveValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsWaveValid(w)), nil
}

func builtinUnloadWave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadWave", args, 1); err != nil {
		return nil, err
	}
	id, w, err := waveByID("unloadWave", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadWave(w)
	delete(waves, id)
	return null(), nil
}

func builtinExportWave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("exportWave", args, 2); err != nil {
		return nil, err
	}
	_, w, err := waveByID("exportWave", args, 0)
	if err != nil {
		return nil, err
	}
	path, _ := argString("exportWave", args, 1)
	rl.ExportWave(w, path)
	return vBool(true), nil
}

// ---- Sound extras ----

func builtinLoadSoundFromWave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadSoundFromWave", args, 1); err != nil {
		return nil, err
	}
	_, w, err := waveByID("loadSoundFromWave", args, 0)
	if err != nil {
		return nil, err
	}
	s := rl.LoadSoundFromWave(w)
	id := nextSoundID
	nextSoundID++
	sounds[id] = s
	return vInt(id), nil
}

func builtinLoadSoundAlias(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadSoundAlias", args, 1); err != nil {
		return nil, err
	}
	_, src, err := soundByID("loadSoundAlias", args, 0)
	if err != nil {
		return nil, err
	}
	s := rl.LoadSoundAlias(src)
	id := nextSoundID
	nextSoundID++
	sounds[id] = s
	return vInt(id), nil
}

func builtinIsSoundValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isSoundValid", args, 1); err != nil {
		return nil, err
	}
	_, s, err := soundByID("isSoundValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsSoundValid(s)), nil
}

func builtinUnloadSoundAlias(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadSoundAlias", args, 1); err != nil {
		return nil, err
	}
	id, s, err := soundByID("unloadSoundAlias", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadSoundAlias(s)
	delete(sounds, id)
	return null(), nil
}

func builtinPauseSound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := soundByID("pauseSound", args, 0)
	if err != nil {
		return nil, err
	}
	rl.PauseSound(s)
	return null(), nil
}

func builtinResumeSound(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := soundByID("resumeSound", args, 0)
	if err != nil {
		return nil, err
	}
	rl.ResumeSound(s)
	return null(), nil
}

func builtinSetSoundPitch(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setSoundPitch expects soundId, pitch")
	}
	_, s, err := soundByID("setSoundPitch", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setSoundPitch", args, 1)
	rl.SetSoundPitch(s, float32(v))
	return null(), nil
}

func builtinSetSoundPan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setSoundPan expects soundId, pan")
	}
	_, s, err := soundByID("setSoundPan", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setSoundPan", args, 1)
	rl.SetSoundPan(s, float32(v))
	return null(), nil
}

// ---- Wave management ----

func builtinWaveCopy(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("waveCopy", args, 1); err != nil {
		return nil, err
	}
	_, w, err := waveByID("waveCopy", args, 0)
	if err != nil {
		return nil, err
	}
	cpy := rl.WaveCopy(w)
	id := nextWaveID
	nextWaveID++
	waves[id] = cpy
	return vInt(id), nil
}

func builtinWaveCrop(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("waveCrop expects waveId, initFrame, finalFrame")
	}
	id, w, err := waveByID("waveCrop", args, 0)
	if err != nil {
		return nil, err
	}
	initF, _ := argInt("waveCrop", args, 1)
	finalF, _ := argInt("waveCrop", args, 2)
	rl.WaveCrop(&w, int32(initF), int32(finalF))
	waves[id] = w
	return null(), nil
}

func builtinWaveFormat(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("waveFormat expects waveId, sampleRate, sampleSize, channels")
	}
	id, w, err := waveByID("waveFormat", args, 0)
	if err != nil {
		return nil, err
	}
	sampleRate, _ := argInt("waveFormat", args, 1)
	sampleSize, _ := argInt("waveFormat", args, 2)
	channels, _ := argInt("waveFormat", args, 3)
	rl.WaveFormat(&w, int32(sampleRate), int32(sampleSize), int32(channels))
	waves[id] = w
	return null(), nil
}

func builtinLoadWaveSamples(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadWaveSamples", args, 1); err != nil {
		return nil, err
	}
	_, w, err := waveByID("loadWaveSamples", args, 0)
	if err != nil {
		return nil, err
	}
	samples := rl.LoadWaveSamples(w)
	elems := make([]candy_evaluator.Value, len(samples))
	for i, s := range samples {
		elems[i] = candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: float64(s)}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: elems}, nil
}

// ---- Music extras ----

func builtinLoadMusicStreamFromMemory(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("loadMusicStreamFromMemory expects fileType, bytesArr")
	}
	ft, _ := argString("loadMusicStreamFromMemory", args, 0)
	if args[1] == nil || args[1].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("loadMusicStreamFromMemory: arg 2 must be byte array")
	}
	data := make([]byte, len(args[1].Elems))
	for i, e := range args[1].Elems {
		data[i] = byte(e.I64)
	}
	m := rl.LoadMusicStreamFromMemory(ft, data, int32(len(data)))
	id := nextMusicID
	nextMusicID++
	musics[id] = m
	return vInt(id), nil
}

func builtinIsMusicValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isMusicValid", args, 1); err != nil {
		return nil, err
	}
	_, m, err := musicByID("isMusicValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsMusicValid(m)), nil
}

func builtinIsMusicStreamPlaying(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, m, err := musicByID("isMusicStreamPlaying", args, 0)
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsMusicStreamPlaying(m)), nil
}

func builtinPauseMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, m, err := musicByID("pauseMusicStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.PauseMusicStream(m)
	return null(), nil
}

func builtinResumeMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, m, err := musicByID("resumeMusicStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.ResumeMusicStream(m)
	return null(), nil
}

func builtinSeekMusicStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("seekMusicStream expects musicId, position")
	}
	_, m, err := musicByID("seekMusicStream", args, 0)
	if err != nil {
		return nil, err
	}
	pos, _ := getArgFloat("seekMusicStream", args, 1)
	rl.SeekMusicStream(m, float32(pos))
	return null(), nil
}

func builtinSetMusicPitch(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setMusicPitch expects musicId, pitch")
	}
	_, m, err := musicByID("setMusicPitch", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setMusicPitch", args, 1)
	rl.SetMusicPitch(m, float32(v))
	return null(), nil
}

func builtinSetMusicPan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setMusicPan expects musicId, pan")
	}
	_, m, err := musicByID("setMusicPan", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setMusicPan", args, 1)
	rl.SetMusicPan(m, float32(v))
	return null(), nil
}

func builtinGetMusicTimeLength(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, m, err := musicByID("getMusicTimeLength", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.GetMusicTimeLength(m))), nil
}

func builtinGetMusicTimePlayed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, m, err := musicByID("getMusicTimePlayed", args, 0)
	if err != nil {
		return nil, err
	}
	return vFloat(float64(rl.GetMusicTimePlayed(m))), nil
}

// ---- AudioStream management ----

func builtinLoadAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("loadAudioStream", args, 3); err != nil {
		return nil, err
	}
	sampleRate, _ := argInt("loadAudioStream", args, 0)
	sampleSize, _ := argInt("loadAudioStream", args, 1)
	channels, _ := argInt("loadAudioStream", args, 2)
	s := rl.LoadAudioStream(uint32(sampleRate), uint32(sampleSize), uint32(channels))
	id := nextAudioStreamID
	nextAudioStreamID++
	audioStreams[id] = s
	return vInt(id), nil
}

func builtinIsAudioStreamValid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("isAudioStreamValid", args, 1); err != nil {
		return nil, err
	}
	_, s, err := audioStreamByID("isAudioStreamValid", args, 0)
	if err != nil {
		return vBool(false), nil
	}
	return vBool(rl.IsAudioStreamValid(s)), nil
}

func builtinUnloadAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("unloadAudioStream", args, 1); err != nil {
		return nil, err
	}
	id, s, err := audioStreamByID("unloadAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.UnloadAudioStream(s)
	delete(audioStreams, id)
	return null(), nil
}

func builtinUpdateAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("updateAudioStream expects streamId, dataArr (float32 samples)")
	}
	_, s, err := audioStreamByID("updateAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	if args[1] == nil || args[1].Kind != candy_evaluator.ValArray {
		return nil, fmt.Errorf("updateAudioStream: arg 2 must be float array")
	}
	samples := make([]float32, len(args[1].Elems))
	for i, e := range args[1].Elems {
		samples[i] = float32(e.F64)
	}
	rl.UpdateAudioStream(s, samples)
	return null(), nil
}

func builtinIsAudioStreamProcessed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("isAudioStreamProcessed", args, 0)
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsAudioStreamProcessed(s)), nil
}

func builtinPlayAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("playAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.PlayAudioStream(s)
	return null(), nil
}

func builtinPauseAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("pauseAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.PauseAudioStream(s)
	return null(), nil
}

func builtinResumeAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("resumeAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.ResumeAudioStream(s)
	return null(), nil
}

func builtinIsAudioStreamPlaying(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("isAudioStreamPlaying", args, 0)
	if err != nil {
		return nil, err
	}
	return vBool(rl.IsAudioStreamPlaying(s)), nil
}

func builtinStopAudioStream(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	_, s, err := audioStreamByID("stopAudioStream", args, 0)
	if err != nil {
		return nil, err
	}
	rl.StopAudioStream(s)
	return null(), nil
}

func builtinSetAudioStreamVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setAudioStreamVolume expects streamId, volume")
	}
	_, s, err := audioStreamByID("setAudioStreamVolume", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setAudioStreamVolume", args, 1)
	rl.SetAudioStreamVolume(s, float32(v))
	return null(), nil
}

func builtinSetAudioStreamPitch(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setAudioStreamPitch expects streamId, pitch")
	}
	_, s, err := audioStreamByID("setAudioStreamPitch", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setAudioStreamPitch", args, 1)
	rl.SetAudioStreamPitch(s, float32(v))
	return null(), nil
}

func builtinSetAudioStreamPan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("setAudioStreamPan expects streamId, pan")
	}
	_, s, err := audioStreamByID("setAudioStreamPan", args, 0)
	if err != nil {
		return nil, err
	}
	v, _ := getArgFloat("setAudioStreamPan", args, 1)
	rl.SetAudioStreamPan(s, float32(v))
	return null(), nil
}

func builtinSetAudioStreamBufferSizeDefault(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setAudioStreamBufferSizeDefault", args, 1); err != nil {
		return nil, err
	}
	sz, _ := argInt("setAudioStreamBufferSizeDefault", args, 0)
	rl.SetAudioStreamBufferSizeDefault(int32(sz))
	return null(), nil
}
