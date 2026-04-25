package candy_raylib

import (
	"candy/candy_evaluator"
	"fmt"
	"math"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type helperCameraState struct {
	x, y       float64
	targetX    float64
	targetY    float64
	smoothness float64
	shake      float64
	zoom       float64
	minX       float64
	minY       float64
	maxX       float64
	maxY       float64
}

type helperTween struct {
	start    float64
	end      float64
	duration float64
	elapsed  float64
	easing   string
	done     bool
}

type helperTimer struct {
	duration float64
	elapsed  float64
	running  bool
}

type helperAnimation struct {
	name      string
	frames    []int64
	frameTime float64
	loop      bool
	current   int
	elapsed   float64
	playing   bool
}

type helperAnimController struct {
	anims   map[string]int64
	current string
}

type helperTask struct {
	delay    float64
	interval float64
	elapsed  float64
	repeat   bool
	fn       *candy_evaluator.Value
	done     bool
}

type helperParticle struct {
	x, y   float64
	vx, vy float64
	life   float64
	max    float64
	size   float64
}

type helperEmitter struct {
	x, y                     float64
	spread                   float64
	speedMin, speedMax       float64
	lifeMin, lifeMax         float64
	sizeStart, sizeEnd       float64
	gravityX, gravityY       float64
	startColor, endColorName string
	parts                    []*helperParticle
}

type helperPathGrid struct {
	w, h, tile int
	blocked    map[int]bool
}

type helperSpatialGrid struct {
	cell   float64
	bucket map[string][]candy_evaluator.Value
}

type helperScene struct {
	initFn   *candy_evaluator.Value
	updateFn *candy_evaluator.Value
	drawFn   *candy_evaluator.Value
	inited   bool
}

var gameCamera = helperCameraState{
	smoothness: 5.0,
	zoom:       1.0,
	minX:       -999999,
	minY:       -999999,
	maxX:       999999,
	maxY:       999999,
}

var (
	nextTweenID int64 = 1
	tweens            = map[int64]*helperTween{}

	nextTimerID      int64 = 1
	timers                 = map[int64]*helperTimer{}
	nextAnimID       int64 = 1
	animations             = map[int64]*helperAnimation{}
	nextCtrlID       int64 = 1
	controllers            = map[int64]*helperAnimController{}
	tasks                  = []*helperTask{}
	nextEmitterID    int64 = 1
	emitters               = map[int64]*helperEmitter{}
	nextPathGridID   int64 = 1
	pathGrids              = map[int64]*helperPathGrid{}
	nextSpatialID    int64 = 1
	spatialGrids           = map[int64]*helperSpatialGrid{}
	scenes                 = map[string]*helperScene{}
	currentScene           = ""
	scenePaused            = false
	profileStart           = map[string]float64{}
	profileAccum           = map[string]float64{}
	projectiles            = map[int64]*candy_evaluator.Value{}
	nextProjectileID int64 = 1
	helperPools            = map[int64][]candy_evaluator.Value{}
	nextPoolID       int64 = 1
	enemyWaves             = map[int64][]candy_evaluator.Value{}
	nextEnemyWaveID  int64 = 1
	stateMachines          = map[int64]map[string]*candy_evaluator.Value{}
	activeState            = map[int64]string{}
	nextSMID         int64 = 1
	quests                 = map[string]map[string]candy_evaluator.Value{}
	tags                   = map[string][]candy_evaluator.Value{}
	entityNetID            = map[*candy_evaluator.Value]int64{}
	nextNetID        int64 = 1
	helperSoundIDs         = map[string]int64{}
	helperMusicIDs         = map[string]int64{}
	currentHelperMusicID int64 = 0

	cooldowns = map[string]float64{}
	saveData  = map[string]candy_evaluator.Value{}
)

func registerGameHelperBuiltins() {
	// camera helpers
	candy_evaluator.RegisterBuiltin("cameraFollow", builtinCameraFollow)
	candy_evaluator.RegisterBuiltin("cameraSnapTo", builtinCameraSnapTo)
	candy_evaluator.RegisterBuiltin("cameraBounds", builtinCameraBounds)
	candy_evaluator.RegisterBuiltin("cameraShake", builtinCameraShake)
	candy_evaluator.RegisterBuiltin("cameraZoom", builtinCameraZoom)
	candy_evaluator.RegisterBuiltin("cameraPosition", builtinCameraPosition)
	candy_evaluator.RegisterBuiltin("screenToWorld", builtinScreenToWorld)
	candy_evaluator.RegisterBuiltin("worldToScreen", builtinWorldToScreen)
	candy_evaluator.RegisterBuiltin("cameraOrbit", builtinCameraOrbit)
	candy_evaluator.RegisterBuiltin("cameraLookAt", builtinCameraLookAt)
	candy_evaluator.RegisterBuiltin("camera3DFollow", builtinCamera3DFollow)
	candy_evaluator.RegisterBuiltin("cameraOrbitInput", builtinCameraOrbitInput)
	candy_evaluator.RegisterBuiltin("cameraYaw", builtinCameraYaw)
	candy_evaluator.RegisterBuiltin("cameraPitch", builtinCameraPitch)
	candy_evaluator.RegisterBuiltin("cameraRoll", builtinCameraRoll)
	candy_evaluator.RegisterBuiltin("cameraForward", builtinCameraForward)
	candy_evaluator.RegisterBuiltin("cameraRight", builtinCameraRight)
	candy_evaluator.RegisterBuiltin("cameraUp", builtinCameraUp)
	candy_evaluator.RegisterBuiltin("screenToWorld3D", builtinScreenToWorld3D)

	// input helpers
	candy_evaluator.RegisterBuiltin("keyDown", builtinKeyDown)
	candy_evaluator.RegisterBuiltin("keyPressed", builtinKeyPressed)
	candy_evaluator.RegisterBuiltin("keyReleased", builtinKeyReleased)
	candy_evaluator.RegisterBuiltin("axis", builtinAxis)
	candy_evaluator.RegisterBuiltin("mousePressed", builtinMousePressed)
	candy_evaluator.RegisterBuiltin("mouseX", builtinMouseX)
	candy_evaluator.RegisterBuiltin("mouseY", builtinMouseY)
	candy_evaluator.RegisterBuiltin("mousePosition", builtinMousePosition)
	candy_evaluator.RegisterBuiltin("gamepadButton", builtinGamepadButton)
	candy_evaluator.RegisterBuiltin("gamepadAxis", builtinGamepadAxis)
	candy_evaluator.RegisterBuiltin("combo", builtinCombo)

	// animation/tween/timer helpers
	candy_evaluator.RegisterBuiltin("animation", builtinAnimation)
	candy_evaluator.RegisterBuiltin("playAnimation", builtinPlayAnimation)
	candy_evaluator.RegisterBuiltin("updateAnimation", builtinUpdateAnimation)
	candy_evaluator.RegisterBuiltin("animationFrame", builtinAnimationFrame)
	candy_evaluator.RegisterBuiltin("animationDone", builtinAnimationDone)
	candy_evaluator.RegisterBuiltin("animController", builtinAnimController)
	candy_evaluator.RegisterBuiltin("addAnimation", builtinAddAnimation)
	candy_evaluator.RegisterBuiltin("setAnimation", builtinSetAnimation)
	candy_evaluator.RegisterBuiltin("updateAnimController", builtinUpdateAnimController)
	candy_evaluator.RegisterBuiltin("controllerFrame", builtinControllerFrame)

	candy_evaluator.RegisterBuiltin("tweenCreate", builtinTweenCreate)
	candy_evaluator.RegisterBuiltin("updateTween", builtinUpdateTween)
	candy_evaluator.RegisterBuiltin("tweenValue", builtinTweenValue)
	candy_evaluator.RegisterBuiltin("tweenDone", builtinTweenDone)
	candy_evaluator.RegisterBuiltin("tweenTo", builtinTweenTo)

	candy_evaluator.RegisterBuiltin("createTimer", builtinCreateTimer)
	candy_evaluator.RegisterBuiltin("updateTimer", builtinUpdateTimer)
	candy_evaluator.RegisterBuiltin("timerDone", builtinTimerDone)
	candy_evaluator.RegisterBuiltin("timerRemaining", builtinTimerRemaining)
	candy_evaluator.RegisterBuiltin("resetTimer", builtinResetTimer)
	candy_evaluator.RegisterBuiltin("cooldownReady", builtinCooldownReady)
	candy_evaluator.RegisterBuiltin("cooldownStart", builtinCooldownStart)
	candy_evaluator.RegisterBuiltin("after", builtinAfter)
	candy_evaluator.RegisterBuiltin("every", builtinEvery)

	// collision helpers
	candy_evaluator.RegisterBuiltin("boxCollision", builtinBoxCollision)
	candy_evaluator.RegisterBuiltin("circleCollision", builtinCircleCollision)
	candy_evaluator.RegisterBuiltin("pointInBox", builtinPointInBox)
	candy_evaluator.RegisterBuiltin("raycast", builtinRaycast)
	candy_evaluator.RegisterBuiltin("inRadius", builtinInRadius)

	// path/spatial
	candy_evaluator.RegisterBuiltin("pathfindGrid", builtinPathfindGrid)
	candy_evaluator.RegisterBuiltin("blockTile", builtinBlockTile)
	candy_evaluator.RegisterBuiltin("findPath", builtinFindPath)
	candy_evaluator.RegisterBuiltin("spatialGrid", builtinSpatialGrid)
	candy_evaluator.RegisterBuiltin("insert", builtinGridInsert)
	candy_evaluator.RegisterBuiltin("queryGrid", builtinQueryGrid)
	candy_evaluator.RegisterBuiltin("clearGrid", builtinClearGrid)

	// particles
	candy_evaluator.RegisterBuiltin("particles", builtinParticles)
	candy_evaluator.RegisterBuiltin("particleSpread", builtinParticleSpread)
	candy_evaluator.RegisterBuiltin("particleSpeed", builtinParticleSpeed)
	candy_evaluator.RegisterBuiltin("particleLife", builtinParticleLife)
	candy_evaluator.RegisterBuiltin("particleColor", builtinParticleColor)
	candy_evaluator.RegisterBuiltin("particleSize", builtinParticleSize)
	candy_evaluator.RegisterBuiltin("particleGravity", builtinParticleGravity)
	candy_evaluator.RegisterBuiltin("emit", builtinEmit)
	candy_evaluator.RegisterBuiltin("updateParticles", builtinUpdateParticles)
	candy_evaluator.RegisterBuiltin("drawParticles", builtinDrawParticles)
	candy_evaluator.RegisterBuiltin("explosion", builtinExplosion)
	candy_evaluator.RegisterBuiltin("smoke", builtinSmoke)
	candy_evaluator.RegisterBuiltin("sparkles", builtinSparkles)
	candy_evaluator.RegisterBuiltin("blood", builtinBlood)

	// scene/debug/audio helpers
	candy_evaluator.RegisterBuiltin("scene", builtinScene)
	candy_evaluator.RegisterBuiltin("startScene", builtinStartScene)
	candy_evaluator.RegisterBuiltin("switchScene", builtinSwitchScene)
	candy_evaluator.RegisterBuiltin("pauseScene", builtinPauseScene)
	candy_evaluator.RegisterBuiltin("resumeScene", builtinResumeScene)
	candy_evaluator.RegisterBuiltin("updateScene", builtinUpdateScene)
	candy_evaluator.RegisterBuiltin("drawScene", builtinDrawScene)
	candy_evaluator.RegisterBuiltin("debugText", builtinDebugText)
	candy_evaluator.RegisterBuiltin("debugBox", builtinDebugBox)
	candy_evaluator.RegisterBuiltin("debugCircle", builtinDebugCircle)
	candy_evaluator.RegisterBuiltin("debugPath", builtinDebugPath)
	candy_evaluator.RegisterBuiltin("startProfile", builtinStartProfile)
	candy_evaluator.RegisterBuiltin("endProfile", builtinEndProfile)
	candy_evaluator.RegisterBuiltin("printProfiles", builtinPrintProfiles)
	candy_evaluator.RegisterBuiltin("playMusic", builtinPlayMusic)
	candy_evaluator.RegisterBuiltin("music", builtinPlayMusic)
	candy_evaluator.RegisterBuiltin("pauseMusic", builtinPauseMusic)
	candy_evaluator.RegisterBuiltin("resumeMusic", builtinResumeMusic)
	candy_evaluator.RegisterBuiltin("stopMusic", builtinStopMusic)
	candy_evaluator.RegisterBuiltin("fadeMusic", builtinFadeMusic)
	candy_evaluator.RegisterBuiltin("fadeMusicIn", builtinFadeMusicIn)
	candy_evaluator.RegisterBuiltin("volume", builtinVolume)
	candy_evaluator.RegisterBuiltin("playSound3D", builtinPlaySound3D)

	// transform/movement helpers
	candy_evaluator.RegisterBuiltin("rotateX", builtinRotateX)
	candy_evaluator.RegisterBuiltin("rotateY", builtinRotateY)
	candy_evaluator.RegisterBuiltin("rotateZ", builtinRotateZ)
	candy_evaluator.RegisterBuiltin("lookAt", builtinLookAt)
	candy_evaluator.RegisterBuiltin("orbit", builtinOrbit)
	candy_evaluator.RegisterBuiltin("moveTowards", builtinMoveTowards)
	candy_evaluator.RegisterBuiltin("rotateTo", builtinRotateTo)
	candy_evaluator.RegisterBuiltin("moveForward", builtinMoveForward)
	candy_evaluator.RegisterBuiltin("moveRight", builtinMoveRight)
	candy_evaluator.RegisterBuiltin("moveUp", builtinMoveUp)
	candy_evaluator.RegisterBuiltin("faceVelocity", builtinFaceVelocity)
	candy_evaluator.RegisterBuiltin("clampPosition", builtinClampPosition)

	// save/load helpers (in-memory runtime store)
	candy_evaluator.RegisterBuiltin("save", builtinSave)
	candy_evaluator.RegisterBuiltin("load", builtinLoad)
	candy_evaluator.RegisterBuiltin("deleteSave", builtinDeleteSave)
	candy_evaluator.RegisterBuiltin("saveExists", builtinSaveExists)
	candy_evaluator.RegisterBuiltin("playSound", builtinPlaySoundHelper)
	candy_evaluator.RegisterBuiltin("play", builtinPlaySoundHelper)

	// extra system helpers
	candy_evaluator.RegisterBuiltin("projectile", builtinProjectile)
	candy_evaluator.RegisterBuiltin("updateProjectiles", builtinUpdateProjectiles)
	candy_evaluator.RegisterBuiltin("drawProjectiles", builtinDrawProjectiles)
	candy_evaluator.RegisterBuiltin("projectileHit", builtinProjectileHit)
	candy_evaluator.RegisterBuiltin("hitscan", builtinHitscan)
	candy_evaluator.RegisterBuiltin("lockOnNearest", builtinLockOnNearest)
	candy_evaluator.RegisterBuiltin("damage", builtinDamage)
	candy_evaluator.RegisterBuiltin("heal", builtinHeal)
	candy_evaluator.RegisterBuiltin("alive", builtinAlive)
	candy_evaluator.RegisterBuiltin("team", builtinTeam)
	candy_evaluator.RegisterBuiltin("setTeam", builtinSetTeam)
	candy_evaluator.RegisterBuiltin("spawn", builtinSpawn)
	candy_evaluator.RegisterBuiltin("despawn", builtinDespawn)
	candy_evaluator.RegisterBuiltin("spawnAtMarker", builtinSpawnAtMarker)
	candy_evaluator.RegisterBuiltin("poolCreate", builtinPoolCreate)
	candy_evaluator.RegisterBuiltin("poolGet", builtinPoolGet)
	candy_evaluator.RegisterBuiltin("poolRelease", builtinPoolRelease)
	candy_evaluator.RegisterBuiltin("waveCreate", builtinWaveCreate)
	candy_evaluator.RegisterBuiltin("waveAdd", builtinWaveAdd)
	candy_evaluator.RegisterBuiltin("waveStart", builtinWaveStart)
	candy_evaluator.RegisterBuiltin("waveDone", builtinWaveDone)
	candy_evaluator.RegisterBuiltin("stateMachine", builtinStateMachine)
	candy_evaluator.RegisterBuiltin("stateAdd", builtinStateAdd)
	candy_evaluator.RegisterBuiltin("stateSet", builtinStateSet)
	candy_evaluator.RegisterBuiltin("stateUpdate", builtinStateUpdate)
	candy_evaluator.RegisterBuiltin("patrol", builtinPatrol)
	candy_evaluator.RegisterBuiltin("chase", builtinChase)
	candy_evaluator.RegisterBuiltin("flee", builtinFlee)
	candy_evaluator.RegisterBuiltin("wander", builtinWander)
	candy_evaluator.RegisterBuiltin("lineOfSight", builtinLineOfSight)
	candy_evaluator.RegisterBuiltin("canSee", builtinCanSee)
	candy_evaluator.RegisterBuiltin("button", builtinButton)
	candy_evaluator.RegisterBuiltin("slider", builtinSlider)
	candy_evaluator.RegisterBuiltin("healthBar", builtinHealthBar)
	candy_evaluator.RegisterBuiltin("floatingText", builtinFloatingText)
	candy_evaluator.RegisterBuiltin("minimap", builtinMinimap)
	candy_evaluator.RegisterBuiltin("questAdd", builtinQuestAdd)
	candy_evaluator.RegisterBuiltin("questComplete", builtinQuestComplete)
	candy_evaluator.RegisterBuiltin("questState", builtinQuestState)
	candy_evaluator.RegisterBuiltin("questStep", builtinQuestStep)
	candy_evaluator.RegisterBuiltin("tag", builtinTag)
	candy_evaluator.RegisterBuiltin("untag", builtinUntag)
	candy_evaluator.RegisterBuiltin("withTag", builtinWithTag)
	candy_evaluator.RegisterBuiltin("distance2D", builtinDistance2D)
	candy_evaluator.RegisterBuiltin("distance3D", builtinDistance3D)
	candy_evaluator.RegisterBuiltin("angleTo", builtinAngleTo)
	candy_evaluator.RegisterBuiltin("lerp", builtinLerp)
	candy_evaluator.RegisterBuiltin("remap", builtinRemap)
	candy_evaluator.RegisterBuiltin("chance", builtinChance)
	candy_evaluator.RegisterBuiltin("randomPointInCircle", builtinRandomPointInCircle)
	candy_evaluator.RegisterBuiltin("randomPointInSphere", builtinRandomPointInSphere)
	candy_evaluator.RegisterBuiltin("netId", builtinNetID)
	candy_evaluator.RegisterBuiltin("setNetOwner", builtinSetNetOwner)
	candy_evaluator.RegisterBuiltin("snapshot", builtinSnapshot)
	candy_evaluator.RegisterBuiltin("interpolateRemote", builtinInterpolateRemote)
	candy_evaluator.RegisterBuiltin("predict", builtinPredict)
	candy_evaluator.RegisterBuiltin("reconcile", builtinReconcile)
}

func stepHelperCamera() {
	dt := float64(rl.GetFrameTime())
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	if gameCamera.smoothness < 0 {
		gameCamera.smoothness = 0
	}
	gameCamera.x += (gameCamera.targetX - gameCamera.x) * gameCamera.smoothness * dt
	gameCamera.y += (gameCamera.targetY - gameCamera.y) * gameCamera.smoothness * dt
	if gameCamera.shake > 0 {
		gameCamera.shake -= 6.0 * dt
		if gameCamera.shake < 0 {
			gameCamera.shake = 0
		}
	}
	gameCamera.x = math.Max(gameCamera.minX, math.Min(gameCamera.maxX, gameCamera.x))
	gameCamera.y = math.Max(gameCamera.minY, math.Min(gameCamera.maxY, gameCamera.y))
}

func currentHelperCamera2D() rl.Camera2D {
	stepHelperCamera()
	offX := float32(rl.GetScreenWidth()) * 0.5
	offY := float32(rl.GetScreenHeight()) * 0.5
	shakeX := float32(float64(rl.GetRandomValue(-1000, 1000)) / 1000.0 * gameCamera.shake)
	shakeY := float32(float64(rl.GetRandomValue(-1000, 1000)) / 1000.0 * gameCamera.shake)
	return rl.Camera2D{
		Offset:   rl.NewVector2(offX, offY),
		Target:   rl.NewVector2(float32(gameCamera.x)+shakeX, float32(gameCamera.y)+shakeY),
		Rotation: 0,
		Zoom:     float32(gameCamera.zoom),
	}
}

func builtinCameraFollow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraFollow", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("cameraFollow", args, 0)
	y, _ := getArgFloat("cameraFollow", args, 1)
	s, _ := getArgFloat("cameraFollow", args, 2)
	gameCamera.targetX, gameCamera.targetY, gameCamera.smoothness = x, y, s
	stepHelperCamera()
	return null(), nil
}

func builtinCameraSnapTo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraSnapTo", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("cameraSnapTo", args, 0)
	y, _ := getArgFloat("cameraSnapTo", args, 1)
	gameCamera.x, gameCamera.y, gameCamera.targetX, gameCamera.targetY = x, y, x, y
	return null(), nil
}

func builtinCameraBounds(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraBounds", args, 4); err != nil {
		return nil, err
	}
	gameCamera.minX, _ = getArgFloat("cameraBounds", args, 0)
	gameCamera.minY, _ = getArgFloat("cameraBounds", args, 1)
	gameCamera.maxX, _ = getArgFloat("cameraBounds", args, 2)
	gameCamera.maxY, _ = getArgFloat("cameraBounds", args, 3)
	return null(), nil
}

func builtinCameraShake(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraShake", args, 1); err != nil {
		return nil, err
	}
	gameCamera.shake, _ = getArgFloat("cameraShake", args, 0)
	return null(), nil
}

func builtinCameraZoom(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraZoom", args, 1); err != nil {
		return nil, err
	}
	gameCamera.zoom, _ = getArgFloat("cameraZoom", args, 0)
	if gameCamera.zoom <= 0 {
		gameCamera.zoom = 1.0
	}
	return null(), nil
}

func builtinCameraPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	stepHelperCamera()
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: gameCamera.x}, "y": {Kind: candy_evaluator.ValFloat, F64: gameCamera.y}}), nil
}

func builtinScreenToWorld(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("screenToWorld", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("screenToWorld", args, 0)
	y, _ := getArgFloat("screenToWorld", args, 1)
	p := rl.GetScreenToWorld2D(rl.NewVector2(float32(x), float32(y)), currentHelperCamera2D())
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(p.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(p.Y)}}), nil
}

func builtinWorldToScreen(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("worldToScreen", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("worldToScreen", args, 0)
	y, _ := getArgFloat("worldToScreen", args, 1)
	p := rl.GetWorldToScreen2D(rl.NewVector2(float32(x), float32(y)), currentHelperCamera2D())
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(p.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(p.Y)}}), nil
}

func builtinCameraOrbit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraOrbit", args, 5); err != nil {
		return nil, err
	}
	tx, _ := getArgFloat("cameraOrbit", args, 0)
	ty, _ := getArgFloat("cameraOrbit", args, 1)
	tz, _ := getArgFloat("cameraOrbit", args, 2)
	radius, _ := getArgFloat("cameraOrbit", args, 3)
	speed, _ := getArgFloat("cameraOrbit", args, 4)
	t := rl.GetTime() * speed
	activeCamera3D.Position = rl.NewVector3(float32(tx+math.Cos(t)*radius), activeCamera3D.Position.Y, float32(tz+math.Sin(t)*radius))
	activeCamera3D.Target = rl.NewVector3(float32(tx), float32(ty), float32(tz))
	return null(), nil
}

func builtinCameraLookAt(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraLookAt", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("cameraLookAt", args, 0)
	y, _ := getArgFloat("cameraLookAt", args, 1)
	z, _ := getArgFloat("cameraLookAt", args, 2)
	activeCamera3D.Target = rl.NewVector3(float32(x), float32(y), float32(z))
	return null(), nil
}
func builtinCamera3DFollow(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("camera3DFollow", args, 6); err != nil {
		return nil, err
	}
	tx, _ := getArgFloat("camera3DFollow", args, 0)
	ty, _ := getArgFloat("camera3DFollow", args, 1)
	tz, _ := getArgFloat("camera3DFollow", args, 2)
	dist, _ := getArgFloat("camera3DFollow", args, 3)
	height, _ := getArgFloat("camera3DFollow", args, 4)
	smooth, _ := getArgFloat("camera3DFollow", args, 5)
	target := rl.NewVector3(float32(tx), float32(ty), float32(tz))
	desired := rl.NewVector3(float32(tx-dist), float32(ty+height), float32(tz-dist))
	activeCamera3D.Position = rl.Vector3Lerp(activeCamera3D.Position, desired, float32(math.Max(0.01, math.Min(1, smooth*float64(rl.GetFrameTime())))))
	activeCamera3D.Target = target
	return null(), nil
}
func builtinCameraOrbitInput(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cameraOrbitInput", args, 5); err != nil {
		return nil, err
	}
	tx, _ := getArgFloat("cameraOrbitInput", args, 0)
	ty, _ := getArgFloat("cameraOrbitInput", args, 1)
	tz, _ := getArgFloat("cameraOrbitInput", args, 2)
	radius, _ := getArgFloat("cameraOrbitInput", args, 3)
	sens, _ := getArgFloat("cameraOrbitInput", args, 4)
	delta := rl.GetMouseDelta()
	yaw := math.Atan2(float64(activeCamera3D.Position.Z)-tz, float64(activeCamera3D.Position.X)-tx) - float64(delta.X)*sens*0.01
	activeCamera3D.Position = rl.NewVector3(float32(tx+math.Cos(yaw)*radius), activeCamera3D.Position.Y, float32(tz+math.Sin(yaw)*radius))
	activeCamera3D.Target = rl.NewVector3(float32(tx), float32(ty), float32(tz))
	return null(), nil
}
func rotateCameraAroundUp(deg float64) {
	dir := rl.Vector3Subtract(activeCamera3D.Target, activeCamera3D.Position)
	rad := deg * math.Pi / 180.0
	c := float32(math.Cos(rad))
	s := float32(math.Sin(rad))
	rot := rl.NewVector3(dir.X*c-dir.Z*s, dir.Y, dir.X*s+dir.Z*c)
	activeCamera3D.Target = rl.Vector3Add(activeCamera3D.Position, rot)
}
func builtinCameraYaw(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { a, _ := getArgFloat("cameraYaw", args, 0); rotateCameraAroundUp(a); return null(), nil }
func builtinCameraPitch(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	a, _ := getArgFloat("cameraPitch", args, 0)
	dir := rl.Vector3Subtract(activeCamera3D.Target, activeCamera3D.Position)
	dir.Y += float32(a * math.Pi / 180.0)
	activeCamera3D.Target = rl.Vector3Add(activeCamera3D.Position, dir)
	return null(), nil
}
func builtinCameraRoll(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinCameraForward(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	d := rl.Vector3Normalize(rl.Vector3Subtract(activeCamera3D.Target, activeCamera3D.Position))
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(d.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(d.Y)}, "z": {Kind: candy_evaluator.ValFloat, F64: float64(d.Z)}}), nil
}
func builtinCameraRight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	f := rl.Vector3Normalize(rl.Vector3Subtract(activeCamera3D.Target, activeCamera3D.Position))
	r := rl.Vector3Normalize(rl.Vector3CrossProduct(f, activeCamera3D.Up))
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(r.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(r.Y)}, "z": {Kind: candy_evaluator.ValFloat, F64: float64(r.Z)}}), nil
}
func builtinCameraUp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	u := activeCamera3D.Up
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(u.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(u.Y)}, "z": {Kind: candy_evaluator.ValFloat, F64: float64(u.Z)}}), nil
}
func builtinScreenToWorld3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("screenToWorld3D", args, 3); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("screenToWorld3D", args, 0)
	y, _ := getArgFloat("screenToWorld3D", args, 1)
	planeY, _ := getArgFloat("screenToWorld3D", args, 2)
	r := rl.GetMouseRay(rl.NewVector2(float32(x), float32(y)), activeCamera3D)
	if math.Abs(float64(r.Direction.Y)) < 1e-6 {
		return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(r.Position.X)}, "y": {Kind: candy_evaluator.ValFloat, F64: planeY}, "z": {Kind: candy_evaluator.ValFloat, F64: float64(r.Position.Z)}}), nil
	}
	t := (planeY - float64(r.Position.Y)) / float64(r.Direction.Y)
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(r.Position.X) + float64(r.Direction.X)*t}, "y": {Kind: candy_evaluator.ValFloat, F64: planeY}, "z": {Kind: candy_evaluator.ValFloat, F64: float64(r.Position.Z) + float64(r.Direction.Z)*t}}), nil
}

func builtinKeyDown(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsKeyDown(args)
}
func builtinKeyPressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsKeyPressed(args)
}
func builtinKeyReleased(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsKeyReleased(args)
}
func builtinMousePressed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsMouseButtonPressed(args)
}
func builtinMouseX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinGetMouseX(args)
}
func builtinMouseY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinGetMouseY(args)
}
func builtinMousePosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinGetMousePosition(args)
}
func builtinGamepadButton(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinIsGamepadButtonDown(args)
}

func builtinGamepadAxis(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("gamepadAxis", args, 2); err != nil {
		return nil, err
	}
	idx, _ := argInt("gamepadAxis", args, 0)
	name, err := argString("gamepadAxis", args, 1)
	if err != nil {
		return nil, err
	}
	axisX := int32(rl.GamepadAxisLeftX)
	axisY := int32(rl.GamepadAxisLeftY)
	switch strings.ToLower(name) {
	case "right":
		axisX = int32(rl.GamepadAxisRightX)
		axisY = int32(rl.GamepadAxisRightY)
	}
	x := rl.GetGamepadAxisMovement(int32(idx), axisX)
	y := rl.GetGamepadAxisMovement(int32(idx), axisY)
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: float64(x)}, "y": {Kind: candy_evaluator.ValFloat, F64: float64(y)}}), nil
}

func builtinAxis(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("axis", args, 1); err != nil {
		return nil, err
	}
	name, err := argString("axis", args, 0)
	if err != nil {
		return nil, err
	}
	val := 0.0
	switch strings.ToLower(name) {
	case "horizontal":
		if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
			val -= 1
		}
		if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
			val += 1
		}
	case "vertical":
		if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
			val -= 1
		}
		if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			val += 1
		}
	}
	return vFloat(val), nil
}

func builtinTweenCreate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("tweenCreate", args, 4); err != nil {
		return nil, err
	}
	start, _ := getArgFloat("tweenCreate", args, 0)
	end, _ := getArgFloat("tweenCreate", args, 1)
	duration, _ := getArgFloat("tweenCreate", args, 2)
	easing, err := argString("tweenCreate", args, 3)
	if err != nil {
		return nil, err
	}
	if duration <= 0 {
		duration = 0.0001
	}
	id := nextTweenID
	nextTweenID++
	tweens[id] = &helperTween{start: start, end: end, duration: duration, easing: easing}
	return vInt(id), nil
}

func builtinUpdateTween(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateTween", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("updateTween", args, 0)
	dt, _ := getArgFloat("updateTween", args, 1)
	if tw, ok := tweens[id]; ok && !tw.done {
		tw.elapsed += dt
		if tw.elapsed >= tw.duration {
			tw.elapsed = tw.duration
			tw.done = true
		}
	}
	return null(), nil
}

func tweenEase(t float64, easing string) float64 {
	switch strings.ToLower(easing) {
	case "easein":
		return t * t
	case "easeout":
		return 1.0 - (1.0-t)*(1.0-t)
	case "easeinout":
		if t < 0.5 {
			return 2.0 * t * t
		}
		return 1 - 2*(1-t)*(1-t)
	default:
		return t
	}
}

func builtinTweenValue(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("tweenValue", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("tweenValue", args, 0)
	tw, ok := tweens[id]
	if !ok {
		return vFloat(0), nil
	}
	t := tw.elapsed / tw.duration
	if t > 1 {
		t = 1
	}
	t = tweenEase(t, tw.easing)
	return vFloat(tw.start + (tw.end-tw.start)*t), nil
}

func builtinTweenDone(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("tweenDone", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("tweenDone", args, 0)
	tw, ok := tweens[id]
	return vBool(ok && tw.done), nil
}

func builtinCreateTimer(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("createTimer", args, 1); err != nil {
		return nil, err
	}
	d, _ := getArgFloat("createTimer", args, 0)
	id := nextTimerID
	nextTimerID++
	timers[id] = &helperTimer{duration: d, running: true}
	return vInt(id), nil
}
func builtinUpdateTimer(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateTimer", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("updateTimer", args, 0)
	dt, _ := getArgFloat("updateTimer", args, 1)
	if tm, ok := timers[id]; ok && tm.running {
		tm.elapsed += dt
	}
	return null(), nil
}
func builtinTimerDone(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("timerDone", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("timerDone", args, 0)
	if tm, ok := timers[id]; ok {
		return vBool(tm.elapsed >= tm.duration), nil
	}
	return vBool(false), nil
}
func builtinTimerRemaining(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("timerRemaining", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("timerRemaining", args, 0)
	if tm, ok := timers[id]; ok {
		return vFloat(math.Max(0, tm.duration-tm.elapsed)), nil
	}
	return vFloat(0), nil
}
func builtinResetTimer(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("resetTimer", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("resetTimer", args, 0)
	d, _ := getArgFloat("resetTimer", args, 1)
	if tm, ok := timers[id]; ok {
		tm.duration = d
		tm.elapsed = 0
		tm.running = true
	}
	return null(), nil
}

func builtinCooldownReady(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cooldownReady", args, 1); err != nil {
		return nil, err
	}
	name, err := argString("cooldownReady", args, 0)
	if err != nil {
		return nil, err
	}
	dt := float64(rl.GetFrameTime())
	for k := range cooldowns {
		cooldowns[k] = math.Max(0, cooldowns[k]-dt)
	}
	return vBool(cooldowns[name] <= 0), nil
}
func builtinCooldownStart(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("cooldownStart", args, 2); err != nil {
		return nil, err
	}
	name, err := argString("cooldownStart", args, 0)
	if err != nil {
		return nil, err
	}
	d, _ := getArgFloat("cooldownStart", args, 1)
	cooldowns[name] = d
	return null(), nil
}

func builtinBoxCollision(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("boxCollision", args, 8); err != nil {
		return nil, err
	}
	x1, _ := getArgFloat("boxCollision", args, 0)
	y1, _ := getArgFloat("boxCollision", args, 1)
	w1, _ := getArgFloat("boxCollision", args, 2)
	h1, _ := getArgFloat("boxCollision", args, 3)
	x2, _ := getArgFloat("boxCollision", args, 4)
	y2, _ := getArgFloat("boxCollision", args, 5)
	w2, _ := getArgFloat("boxCollision", args, 6)
	h2, _ := getArgFloat("boxCollision", args, 7)
	return vBool(rl.CheckCollisionRecs(rl.NewRectangle(float32(x1), float32(y1), float32(w1), float32(h1)), rl.NewRectangle(float32(x2), float32(y2), float32(w2), float32(h2)))), nil
}
func builtinCircleCollision(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("circleCollision", args, 6); err != nil {
		return nil, err
	}
	x1, _ := getArgFloat("circleCollision", args, 0)
	y1, _ := getArgFloat("circleCollision", args, 1)
	r1, _ := getArgFloat("circleCollision", args, 2)
	x2, _ := getArgFloat("circleCollision", args, 3)
	y2, _ := getArgFloat("circleCollision", args, 4)
	r2, _ := getArgFloat("circleCollision", args, 5)
	return vBool(rl.CheckCollisionCircles(rl.NewVector2(float32(x1), float32(y1)), float32(r1), rl.NewVector2(float32(x2), float32(y2)), float32(r2))), nil
}
func builtinPointInBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("pointInBox", args, 6); err != nil {
		return nil, err
	}
	px, _ := getArgFloat("pointInBox", args, 0)
	py, _ := getArgFloat("pointInBox", args, 1)
	x, _ := getArgFloat("pointInBox", args, 2)
	y, _ := getArgFloat("pointInBox", args, 3)
	w, _ := getArgFloat("pointInBox", args, 4)
	h, _ := getArgFloat("pointInBox", args, 5)
	return vBool(rl.CheckCollisionPointRec(rl.NewVector2(float32(px), float32(py)), rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)))), nil
}

func getMapNum(obj *candy_evaluator.Value, key string) float64 {
	if obj == nil || obj.Kind != candy_evaluator.ValMap {
		return 0
	}
	v, ok := obj.StrMap[key]
	if !ok {
		return 0
	}
	if v.Kind == candy_evaluator.ValInt {
		return float64(v.I64)
	}
	if v.Kind == candy_evaluator.ValFloat {
		return v.F64
	}
	return 0
}
func setMapNum(obj *candy_evaluator.Value, key string, val float64) {
	if obj == nil || obj.Kind != candy_evaluator.ValMap {
		return
	}
	obj.StrMap[key] = candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: val}
}

func builtinRotateX(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("rotateX", args, 2); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("rotateX", args, 1)
	setMapNum(args[0], "pitch", getMapNum(args[0], "pitch")+a)
	return null(), nil
}
func builtinRotateY(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("rotateY", args, 2); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("rotateY", args, 1)
	setMapNum(args[0], "yaw", getMapNum(args[0], "yaw")+a)
	return null(), nil
}
func builtinRotateZ(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("rotateZ", args, 2); err != nil {
		return nil, err
	}
	a, _ := getArgFloat("rotateZ", args, 1)
	setMapNum(args[0], "roll", getMapNum(args[0], "roll")+a)
	return null(), nil
}
func builtinLookAt(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("lookAt expects object, targetX, targetY, targetZ")
	}
	ox, oy, oz := getMapNum(args[0], "x"), getMapNum(args[0], "y"), getMapNum(args[0], "z")
	tx, _ := getArgFloat("lookAt", args, 1)
	ty, _ := getArgFloat("lookAt", args, 2)
	tz, _ := getArgFloat("lookAt", args, 3)
	dx, dy, dz := tx-ox, ty-oy, tz-oz
	yaw := math.Atan2(dz, dx) * 180 / math.Pi
	pitch := math.Atan2(dy, math.Sqrt(dx*dx+dz*dz)) * 180 / math.Pi
	setMapNum(args[0], "yaw", yaw)
	setMapNum(args[0], "pitch", pitch)
	return null(), nil
}
func builtinOrbit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("orbit expects object, cx, cy, cz, radius, speed")
	}
	cx, _ := getArgFloat("orbit", args, 1)
	cy, _ := getArgFloat("orbit", args, 2)
	cz, _ := getArgFloat("orbit", args, 3)
	r, _ := getArgFloat("orbit", args, 4)
	speed, _ := getArgFloat("orbit", args, 5)
	ang := getMapNum(args[0], "__orbitAngle") + speed*float64(rl.GetFrameTime())
	setMapNum(args[0], "__orbitAngle", ang)
	setMapNum(args[0], "x", cx+math.Cos(ang)*r)
	setMapNum(args[0], "y", cy)
	setMapNum(args[0], "z", cz+math.Sin(ang)*r)
	return null(), nil
}
func builtinMoveTowards(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("moveTowards expects object, targetX, targetY, speed")
	}
	tx, _ := getArgFloat("moveTowards", args, 1)
	ty, _ := getArgFloat("moveTowards", args, 2)
	speed, _ := getArgFloat("moveTowards", args, 3)
	x := getMapNum(args[0], "x")
	y := getMapNum(args[0], "y")
	dx, dy := tx-x, ty-y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0.00001 {
		step := speed * float64(rl.GetFrameTime())
		if step > dist {
			step = dist
		}
		setMapNum(args[0], "x", x+dx/dist*step)
		setMapNum(args[0], "y", y+dy/dist*step)
	}
	return null(), nil
}
func builtinRotateTo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 { return nil, fmt.Errorf("rotateTo expects object,pitch,yaw,roll,speed") }
	p, y, r, s := numOr0(args[1]), numOr0(args[2]), numOr0(args[3]), numOr0(args[4])
	curP, curY, curR := getMapNum(args[0], "pitch"), getMapNum(args[0], "yaw"), getMapNum(args[0], "roll")
	t := math.Max(0.01, math.Min(1.0, s*float64(rl.GetFrameTime())))
	setMapNum(args[0], "pitch", curP+(p-curP)*t); setMapNum(args[0], "yaw", curY+(y-curY)*t); setMapNum(args[0], "roll", curR+(r-curR)*t)
	return null(), nil
}
func builtinMoveForward(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 { return nil, fmt.Errorf("moveForward expects object,speed,dt") }
	speed, dt := numOr0(args[1]), numOr0(args[2]); yaw := getMapNum(args[0], "yaw") * math.Pi / 180.0
	setMapNum(args[0], "x", getMapNum(args[0], "x")+math.Cos(yaw)*speed*dt); setMapNum(args[0], "z", getMapNum(args[0], "z")+math.Sin(yaw)*speed*dt); return null(), nil
}
func builtinMoveRight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 { return nil, fmt.Errorf("moveRight expects object,speed,dt") }
	speed, dt := numOr0(args[1]), numOr0(args[2]); yaw := (getMapNum(args[0], "yaw")+90) * math.Pi / 180.0
	setMapNum(args[0], "x", getMapNum(args[0], "x")+math.Cos(yaw)*speed*dt); setMapNum(args[0], "z", getMapNum(args[0], "z")+math.Sin(yaw)*speed*dt); return null(), nil
}
func builtinMoveUp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 { return nil, fmt.Errorf("moveUp expects object,speed,dt") }
	setMapNum(args[0], "y", getMapNum(args[0], "y")+numOr0(args[1])*numOr0(args[2])); return null(), nil
}
func builtinFaceVelocity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 { return nil, fmt.Errorf("faceVelocity expects object") }
	vx, vz := getMapNum(args[0], "vx"), getMapNum(args[0], "vz")
	if math.Abs(vx)+math.Abs(vz) > 1e-6 { setMapNum(args[0], "yaw", math.Atan2(vz, vx)*180/math.Pi) }
	return null(), nil
}
func builtinClampPosition(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 { return nil, fmt.Errorf("clampPosition expects object,minX,minY,minZ,maxX,maxY,maxZ") }
	setMapNum(args[0], "x", math.Max(numOr0(args[1]), math.Min(numOr0(args[4]), getMapNum(args[0], "x"))))
	setMapNum(args[0], "y", math.Max(numOr0(args[2]), math.Min(numOr0(args[5]), getMapNum(args[0], "y"))))
	setMapNum(args[0], "z", math.Max(numOr0(args[3]), math.Min(numOr0(args[6]), getMapNum(args[0], "z"))))
	return null(), nil
}

func builtinSave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("save", args, 2); err != nil {
		return nil, err
	}
	key, err := argString("save", args, 0)
	if err != nil {
		return nil, err
	}
	saveData[key] = *args[1]
	return null(), nil
}
func builtinLoad(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("load expects key, [default]")
	}
	key, err := argString("load", args, 0)
	if err != nil {
		return nil, err
	}
	if v, ok := saveData[key]; ok {
		vv := v
		return &vv, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return null(), nil
}
func builtinDeleteSave(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("deleteSave", args, 1); err != nil {
		return nil, err
	}
	key, err := argString("deleteSave", args, 0)
	if err != nil {
		return nil, err
	}
	delete(saveData, key)
	return null(), nil
}
func builtinSaveExists(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("saveExists", args, 1); err != nil {
		return nil, err
	}
	key, err := argString("saveExists", args, 0)
	if err != nil {
		return nil, err
	}
	_, ok := saveData[key]
	return vBool(ok), nil
}

func helperTick(dt float64) {
	if dt <= 0 {
		dt = float64(rl.GetFrameTime())
	}
	for _, t := range tasks {
		if t.done {
			continue
		}
		t.elapsed += dt
		if t.elapsed >= t.delay {
			_, _ = candy_evaluator.InvokeCallable(t.fn, nil)
			if t.repeat {
				t.elapsed = 0
				t.delay = t.interval
			} else {
				t.done = true
			}
		}
	}
}

func builtinCombo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) == 0 {
		return vBool(false), nil
	}
	last := args[len(args)-1]
	if last == nil || last.Kind != candy_evaluator.ValString {
		return vBool(false), nil
	}
	return vBool(rl.IsKeyPressed(keyCode(last.Str))), nil
}

func valueToFrameList(v *candy_evaluator.Value) []int64 {
	if v == nil || v.Kind != candy_evaluator.ValArray {
		return []int64{0}
	}
	out := make([]int64, 0, len(v.Elems))
	for i := range v.Elems {
		it := v.Elems[i]
		switch it.Kind {
		case candy_evaluator.ValInt:
			out = append(out, it.I64)
		case candy_evaluator.ValFloat:
			out = append(out, int64(it.F64))
		}
	}
	if len(out) == 0 {
		out = append(out, 0)
	}
	return out
}

func builtinAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("animation expects name, frames, frameTime, loop")
	}
	name, err := argString("animation", args, 0)
	if err != nil {
		return nil, err
	}
	ft, _ := getArgFloat("animation", args, 2)
	loop := false
	if args[3] != nil && args[3].Kind == candy_evaluator.ValBool {
		loop = args[3].B
	}
	id := nextAnimID
	nextAnimID++
	animations[id] = &helperAnimation{name: name, frames: valueToFrameList(args[1]), frameTime: math.Max(0.0001, ft), loop: loop, playing: true}
	return vInt(id), nil
}
func builtinPlayAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("playAnimation", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("playAnimation", args, 0)
	if a, ok := animations[id]; ok {
		a.playing = true
	}
	return null(), nil
}
func builtinUpdateAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateAnimation", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("updateAnimation", args, 0)
	dt, _ := getArgFloat("updateAnimation", args, 1)
	if a, ok := animations[id]; ok && a.playing {
		a.elapsed += dt
		if a.elapsed >= a.frameTime {
			a.elapsed -= a.frameTime
			a.current++
			if a.current >= len(a.frames) {
				if a.loop {
					a.current = 0
				} else {
					a.current = len(a.frames) - 1
					a.playing = false
				}
			}
		}
	}
	return null(), nil
}
func builtinAnimationFrame(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("animationFrame", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("animationFrame", args, 0)
	if a, ok := animations[id]; ok {
		return vInt(a.frames[a.current]), nil
	}
	return vInt(0), nil
}
func builtinAnimationDone(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("animationDone", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("animationDone", args, 0)
	if a, ok := animations[id]; ok {
		return vBool(!a.loop && !a.playing), nil
	}
	return vBool(true), nil
}
func builtinAnimController(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextCtrlID
	nextCtrlID++
	controllers[id] = &helperAnimController{anims: map[string]int64{}}
	return vInt(id), nil
}
func builtinAddAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("addAnimation expects controller, name, frames, frameTime, loop")
	}
	cid, _ := argInt("addAnimation", args, 0)
	c, ok := controllers[cid]
	if !ok {
		return null(), nil
	}
	anim, err := builtinAnimation(args[1:])
	if err != nil {
		return nil, err
	}
	name, _ := argString("addAnimation", args, 1)
	c.anims[name] = anim.I64
	if c.current == "" {
		c.current = name
	}
	return null(), nil
}
func builtinSetAnimation(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("setAnimation", args, 2); err != nil {
		return nil, err
	}
	cid, _ := argInt("setAnimation", args, 0)
	name, err := argString("setAnimation", args, 1)
	if err != nil {
		return nil, err
	}
	if c, ok := controllers[cid]; ok {
		if _, ok := c.anims[name]; ok {
			c.current = name
		}
	}
	return null(), nil
}
func builtinUpdateAnimController(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("updateAnimController", args, 2); err != nil {
		return nil, err
	}
	cid, _ := argInt("updateAnimController", args, 0)
	dt, _ := getArgFloat("updateAnimController", args, 1)
	if c, ok := controllers[cid]; ok && c.current != "" {
		if aid, ok2 := c.anims[c.current]; ok2 {
			_, _ = builtinUpdateAnimation([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: aid}, {Kind: candy_evaluator.ValFloat, F64: dt}})
		}
	}
	return null(), nil
}
func builtinControllerFrame(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("controllerFrame", args, 1); err != nil {
		return nil, err
	}
	cid, _ := argInt("controllerFrame", args, 0)
	if c, ok := controllers[cid]; ok && c.current != "" {
		if aid, ok2 := c.anims[c.current]; ok2 {
			return builtinAnimationFrame([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: aid}})
		}
	}
	return vInt(0), nil
}

func builtinTweenTo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("tweenTo expects object, field, target, duration [, callback]")
	}
	obj := args[0]
	field, err := argString("tweenTo", args, 1)
	if err != nil {
		return nil, err
	}
	target, _ := getArgFloat("tweenTo", args, 2)
	duration, _ := getArgFloat("tweenTo", args, 3)
	start := getMapNum(obj, field)
	idv, err := builtinTweenCreate([]*candy_evaluator.Value{
		{Kind: candy_evaluator.ValFloat, F64: start},
		{Kind: candy_evaluator.ValFloat, F64: target},
		{Kind: candy_evaluator.ValFloat, F64: duration},
		{Kind: candy_evaluator.ValString, Str: "easeOut"},
	})
	if err != nil {
		return nil, err
	}
	// apply current value immediately to keep API simple.
	setMapNum(obj, field, start)
	return idv, nil
}

func builtinAfter(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("after expects delay, callback")
	}
	delay, _ := getArgFloat("after", args, 0)
	tasks = append(tasks, &helperTask{delay: delay, fn: args[1]})
	return null(), nil
}
func builtinEvery(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("every expects interval, callback")
	}
	interval, _ := getArgFloat("every", args, 0)
	tasks = append(tasks, &helperTask{delay: interval, interval: interval, repeat: true, fn: args[1]})
	return null(), nil
}

func builtinRaycast(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("raycast expects x1, y1, x2, y2, obstacles")
	}
	x1, _ := getArgFloat("raycast", args, 0)
	y1, _ := getArgFloat("raycast", args, 1)
	x2, _ := getArgFloat("raycast", args, 2)
	y2, _ := getArgFloat("raycast", args, 3)
	if args[4] == nil || args[4].Kind != candy_evaluator.ValArray {
		return vBool(false), nil
	}
	for _, o := range args[4].Elems {
		if o.Kind != candy_evaluator.ValMap {
			continue
		}
		r := rl.NewRectangle(float32(getMapNum(&o, "x")), float32(getMapNum(&o, "y")), float32(getMapNum(&o, "w")), float32(getMapNum(&o, "h")))
		hit := rl.CheckCollisionLines(rl.NewVector2(float32(x1), float32(y1)), rl.NewVector2(float32(x2), float32(y2)),
			rl.NewVector2(r.X, r.Y), rl.NewVector2(r.X+r.Width, r.Y), nil)
		if hit {
			return vBool(true), nil
		}
	}
	return vBool(false), nil
}
func builtinInRadius(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("inRadius expects x, y, radius, objects")
	}
	x, _ := getArgFloat("inRadius", args, 0)
	y, _ := getArgFloat("inRadius", args, 1)
	r, _ := getArgFloat("inRadius", args, 2)
	out := []candy_evaluator.Value{}
	if args[3] != nil && args[3].Kind == candy_evaluator.ValArray {
		for i := range args[3].Elems {
			o := args[3].Elems[i]
			dx := getMapNum(&o, "x") - x
			dy := getMapNum(&o, "y") - y
			if math.Sqrt(dx*dx+dy*dy) <= r {
				out = append(out, o)
			}
		}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: out}, nil
}

func cellKey(x, y int) string { return fmt.Sprintf("%d:%d", x, y) }
func builtinPathfindGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("pathfindGrid", args, 3); err != nil {
		return nil, err
	}
	w, _ := argInt("pathfindGrid", args, 0)
	h, _ := argInt("pathfindGrid", args, 1)
	t, _ := argInt("pathfindGrid", args, 2)
	id := nextPathGridID
	nextPathGridID++
	pathGrids[id] = &helperPathGrid{w: int(w), h: int(h), tile: int(t), blocked: map[int]bool{}}
	return vInt(id), nil
}
func builtinBlockTile(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("blockTile", args, 3); err != nil {
		return nil, err
	}
	id, _ := argInt("blockTile", args, 0)
	x, _ := argInt("blockTile", args, 1)
	y, _ := argInt("blockTile", args, 2)
	if g, ok := pathGrids[id]; ok {
		g.blocked[int(y)*g.w+int(x)] = true
	}
	return null(), nil
}
func builtinFindPath(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("findPath", args, 5); err != nil {
		return nil, err
	}
	id, _ := argInt("findPath", args, 0)
	sx, _ := argInt("findPath", args, 1)
	sy, _ := argInt("findPath", args, 2)
	gx, _ := argInt("findPath", args, 3)
	gy, _ := argInt("findPath", args, 4)
	g, ok := pathGrids[id]
	if !ok {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: []candy_evaluator.Value{}}, nil
	}
	type node struct{ x, y int }
	q := []node{{int(sx), int(sy)}}
	prev := map[string]node{}
	seen := map[string]bool{cellKey(int(sx), int(sy)): true}
	dirs := []node{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	found := false
	for len(q) > 0 && !found {
		n := q[0]
		q = q[1:]
		if n.x == int(gx) && n.y == int(gy) {
			found = true
			break
		}
		for _, d := range dirs {
			nx, ny := n.x+d.x, n.y+d.y
			if nx < 0 || ny < 0 || nx >= g.w || ny >= g.h || g.blocked[ny*g.w+nx] {
				continue
			}
			k := cellKey(nx, ny)
			if seen[k] {
				continue
			}
			seen[k] = true
			prev[k] = n
			q = append(q, node{nx, ny})
		}
	}
	if !found {
		return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: []candy_evaluator.Value{}}, nil
	}
	path := []candy_evaluator.Value{}
	cur := node{int(gx), int(gy)}
	for !(cur.x == int(sx) && cur.y == int(sy)) {
		path = append(path, candy_evaluator.Value{Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{
			"x": {Kind: candy_evaluator.ValFloat, F64: float64(cur.x * g.tile)},
			"y": {Kind: candy_evaluator.ValFloat, F64: float64(cur.y * g.tile)},
		}})
		p, ok := prev[cellKey(cur.x, cur.y)]
		if !ok {
			break
		}
		cur = p
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: path}, nil
}

func builtinSpatialGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("spatialGrid", args, 3); err != nil {
		return nil, err
	}
	cell, _ := getArgFloat("spatialGrid", args, 2)
	id := nextSpatialID
	nextSpatialID++
	spatialGrids[id] = &helperSpatialGrid{cell: math.Max(1, cell), bucket: map[string][]candy_evaluator.Value{}}
	return vInt(id), nil
}
func builtinGridInsert(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("insert", args, 2); err != nil {
		return nil, err
	}
	id, _ := argInt("insert", args, 0)
	if g, ok := spatialGrids[id]; ok && args[1] != nil && args[1].Kind == candy_evaluator.ValMap {
		x := int(getMapNum(args[1], "x") / g.cell)
		y := int(getMapNum(args[1], "y") / g.cell)
		k := cellKey(x, y)
		g.bucket[k] = append(g.bucket[k], *args[1])
	}
	return null(), nil
}
func builtinQueryGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("queryGrid", args, 5); err != nil {
		return nil, err
	}
	id, _ := argInt("queryGrid", args, 0)
	x, _ := getArgFloat("queryGrid", args, 1)
	y, _ := getArgFloat("queryGrid", args, 2)
	if g, ok := spatialGrids[id]; ok {
		cx := int(x / g.cell)
		cy := int(y / g.cell)
		out := []candy_evaluator.Value{}
		for iy := cy - 1; iy <= cy+1; iy++ {
			for ix := cx - 1; ix <= cx+1; ix++ {
				out = append(out, g.bucket[cellKey(ix, iy)]...)
			}
		}
		return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: out}, nil
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: []candy_evaluator.Value{}}, nil
}
func builtinClearGrid(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("clearGrid", args, 1); err != nil {
		return nil, err
	}
	id, _ := argInt("clearGrid", args, 0)
	if g, ok := spatialGrids[id]; ok {
		g.bucket = map[string][]candy_evaluator.Value{}
	}
	return null(), nil
}

func builtinParticles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("particles", args, 2); err != nil {
		return nil, err
	}
	x, _ := getArgFloat("particles", args, 0)
	y, _ := getArgFloat("particles", args, 1)
	id := nextEmitterID
	nextEmitterID++
	emitters[id] = &helperEmitter{x: x, y: y, spread: 360, speedMin: 50, speedMax: 120, lifeMin: 0.5, lifeMax: 1.0, sizeStart: 4, sizeEnd: 1, startColor: "yellow", endColorName: "red"}
	return vInt(id), nil
}
func builtinParticleSpread(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleSpread", args, 0)
	v, _ := getArgFloat("particleSpread", args, 1)
	if e, ok := emitters[id]; ok {
		e.spread = v
	}
	return null(), nil
}
func builtinParticleSpeed(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleSpeed", args, 0)
	a, _ := getArgFloat("particleSpeed", args, 1)
	b, _ := getArgFloat("particleSpeed", args, 2)
	if e, ok := emitters[id]; ok {
		e.speedMin, e.speedMax = a, b
	}
	return null(), nil
}
func builtinParticleLife(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleLife", args, 0)
	a, _ := getArgFloat("particleLife", args, 1)
	b, _ := getArgFloat("particleLife", args, 2)
	if e, ok := emitters[id]; ok {
		e.lifeMin, e.lifeMax = a, b
	}
	return null(), nil
}
func builtinParticleColor(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleColor", args, 0)
	a, _ := argString("particleColor", args, 1)
	b, _ := argString("particleColor", args, 2)
	if e, ok := emitters[id]; ok {
		e.startColor, e.endColorName = a, b
	}
	return null(), nil
}
func builtinParticleSize(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleSize", args, 0)
	a, _ := getArgFloat("particleSize", args, 1)
	b, _ := getArgFloat("particleSize", args, 2)
	if e, ok := emitters[id]; ok {
		e.sizeStart, e.sizeEnd = a, b
	}
	return null(), nil
}
func builtinParticleGravity(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("particleGravity", args, 0)
	gx, _ := getArgFloat("particleGravity", args, 1)
	gy, _ := getArgFloat("particleGravity", args, 2)
	if e, ok := emitters[id]; ok {
		e.gravityX, e.gravityY = gx, gy
	}
	return null(), nil
}
func builtinEmit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("emit", args, 0)
	c, _ := argInt("emit", args, 1)
	e, ok := emitters[id]
	if !ok {
		return null(), nil
	}
	for i := 0; i < int(c); i++ {
		a := (float64(rl.GetRandomValue(0, 1000)) / 1000.0) * e.spread * math.Pi / 180.0
		s := e.speedMin + (float64(rl.GetRandomValue(0, 1000))/1000.0)*(e.speedMax-e.speedMin)
		l := e.lifeMin + (float64(rl.GetRandomValue(0, 1000))/1000.0)*(e.lifeMax-e.lifeMin)
		e.parts = append(e.parts, &helperParticle{x: e.x, y: e.y, vx: math.Cos(a) * s, vy: math.Sin(a) * s, max: l, life: l, size: e.sizeStart})
	}
	return null(), nil
}
func builtinUpdateParticles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("updateParticles", args, 0)
	dt, _ := getArgFloat("updateParticles", args, 1)
	e, ok := emitters[id]
	if !ok {
		return null(), nil
	}
	out := make([]*helperParticle, 0, len(e.parts))
	for _, p := range e.parts {
		p.vx += e.gravityX * dt
		p.vy += e.gravityY * dt
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.life -= dt
		if p.life > 0 {
			out = append(out, p)
		}
	}
	e.parts = out
	return null(), nil
}
func builtinDrawParticles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("drawParticles", args, 0)
	e, ok := emitters[id]
	if !ok {
		return null(), nil
	}
	for _, p := range e.parts {
		rl.DrawCircle(int32(p.x), int32(p.y), float32(p.size), colorFrom(e.startColor))
	}
	return null(), nil
}
func builtinExplosion(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	idv, _ := builtinParticles([]*candy_evaluator.Value{args[0], args[1]})
	_, _ = builtinEmit([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: idv.I64}, {Kind: candy_evaluator.ValInt, I64: 40}})
	return null(), nil
}
func builtinSmoke(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	idv, _ := builtinParticles([]*candy_evaluator.Value{args[0], args[1]})
	_, _ = builtinEmit([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: idv.I64}, {Kind: candy_evaluator.ValInt, I64: 20}})
	return null(), nil
}
func builtinSparkles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	idv, _ := builtinParticles([]*candy_evaluator.Value{args[0], args[1]})
	_, _ = builtinEmit([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: idv.I64}, {Kind: candy_evaluator.ValInt, I64: 15}})
	return null(), nil
}
func builtinBlood(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	idv, _ := builtinParticles([]*candy_evaluator.Value{args[0], args[1]})
	_, _ = builtinEmit([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: idv.I64}, {Kind: candy_evaluator.ValInt, I64: 25}})
	return null(), nil
}

func builtinScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("scene expects name, callbacks")
	}
	name, err := argString("scene", args, 0)
	if err != nil {
		return nil, err
	}
	if args[1] == nil || args[1].Kind != candy_evaluator.ValMap {
		return nil, fmt.Errorf("scene callbacks must be map")
	}
	cb := args[1].StrMap
	scenes[name] = &helperScene{initFn: valuePtr(cb["init"]), updateFn: valuePtr(cb["update"]), drawFn: valuePtr(cb["draw"])}
	return null(), nil
}
func valuePtr(v candy_evaluator.Value) *candy_evaluator.Value { vv := v; return &vv }
func runSceneInit(name string) {
	if s, ok := scenes[name]; ok && !s.inited {
		_, _ = candy_evaluator.InvokeCallable(s.initFn, nil)
		s.inited = true
	}
}
func builtinStartScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return builtinSwitchScene(args)
}
func builtinSwitchScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("switchScene", args, 1); err != nil {
		return nil, err
	}
	name, err := argString("switchScene", args, 0)
	if err != nil {
		return nil, err
	}
	currentScene = name
	scenePaused = false
	runSceneInit(name)
	return null(), nil
}
func builtinPauseScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	scenePaused = true
	return null(), nil
}
func builtinResumeScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	scenePaused = false
	return null(), nil
}
func builtinUpdateScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if currentScene == "" || scenePaused {
		return null(), nil
	}
	dt := float64(rl.GetFrameTime())
	if len(args) == 1 {
		dt, _ = getArgFloat("updateScene", args, 0)
	}
	// helperTick runs from show()/endDrawing/flip so every/after work in normal kid loops; avoid double-tick.
	if s, ok := scenes[currentScene]; ok {
		_, _ = candy_evaluator.InvokeCallable(s.updateFn, []*candy_evaluator.Value{{Kind: candy_evaluator.ValFloat, F64: dt}})
	}
	return null(), nil
}
func builtinDrawScene(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if currentScene == "" {
		return null(), nil
	}
	if s, ok := scenes[currentScene]; ok {
		_, _ = candy_evaluator.InvokeCallable(s.drawFn, nil)
	}
	return null(), nil
}

func builtinDebugText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("debugText expects x,y,text")
	}
	x, _ := getArgFloat("debugText", args, 0)
	y, _ := getArgFloat("debugText", args, 1)
	txt, _ := argString("debugText", args, 2)
	rl.DrawText(txt, int32(x), int32(y), 16, rl.Yellow)
	return null(), nil
}
func builtinDebugBox(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("debugBox expects x,y,w,h,color")
	}
	x, _ := getArgFloat("debugBox", args, 0)
	y, _ := getArgFloat("debugBox", args, 1)
	w, _ := getArgFloat("debugBox", args, 2)
	h, _ := getArgFloat("debugBox", args, 3)
	col, _ := argString("debugBox", args, 4)
	rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), colorFrom(col))
	return null(), nil
}
func builtinDebugCircle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("debugCircle expects x,y,r,color")
	}
	x, _ := getArgFloat("debugCircle", args, 0)
	y, _ := getArgFloat("debugCircle", args, 1)
	r, _ := getArgFloat("debugCircle", args, 2)
	col, _ := argString("debugCircle", args, 3)
	rl.DrawCircleLines(int32(x), int32(y), float32(r), colorFrom(col))
	return null(), nil
}
func builtinDebugPath(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 || args[0] == nil || args[0].Kind != candy_evaluator.ValArray {
		return null(), nil
	}
	col, _ := argString("debugPath", args, 1)
	c := colorFrom(col)
	pts := args[0].Elems
	for i := 0; i+1 < len(pts); i++ {
		a, b := pts[i], pts[i+1]
		rl.DrawLine(int32(getMapNum(&a, "x")), int32(getMapNum(&a, "y")), int32(getMapNum(&b, "x")), int32(getMapNum(&b, "y")), c)
	}
	return null(), nil
}
func builtinStartProfile(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	n, _ := argString("startProfile", args, 0)
	profileStart[n] = rl.GetTime()
	return null(), nil
}
func builtinEndProfile(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	n, _ := argString("endProfile", args, 0)
	profileAccum[n] += rl.GetTime() - profileStart[n]
	return null(), nil
}
func builtinPrintProfiles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	keys := make([]string, 0, len(profileAccum))
	for k := range profileAccum {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %.6fs\n", k, profileAccum[k])
	}
	return null(), nil
}

func builtinPlayMusic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("playMusic expects path")
	}
	path, err := argString("playMusic", args, 0)
	if err != nil {
		return nil, err
	}
	id, ok := helperMusicIDs[path]
	if !ok {
		v, e := builtinLoadMusicStream([]*candy_evaluator.Value{{Kind: candy_evaluator.ValString, Str: path}})
		if e != nil {
			return nil, e
		}
		id = v.I64
		helperMusicIDs[path] = id
	}
	currentHelperMusicID = id
	return builtinPlayMusicStream([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: id}})
}
func builtinPauseMusic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if currentHelperMusicID == 0 {
		return null(), nil
	}
	return builtinPauseMusicStream([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: currentHelperMusicID}})
}
func builtinResumeMusic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if currentHelperMusicID == 0 {
		return null(), nil
	}
	return builtinResumeMusicStream([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: currentHelperMusicID}})
}
func builtinStopMusic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if currentHelperMusicID == 0 {
		return null(), nil
	}
	return builtinStopMusicStream([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: currentHelperMusicID}})
}
func builtinFadeMusic(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return null(), nil
}
func builtinFadeMusicIn(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	return null(), nil
}
func builtinVolume(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if err := expectArgs("volume", args, 1); err != nil {
		return nil, err
	}
	v, _ := getArgFloat("volume", args, 0)
	rl.SetMasterVolume(float32(v))
	return null(), nil
}
func builtinPlaySound3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	// Fallback to 2D playback for now; 3D attenuation can be layered later.
	if len(args) < 1 {
		return nil, fmt.Errorf("playSound3D expects path")
	}
	return builtinPlaySoundHelper(args)
}
func builtinPlaySoundHelper(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("playSound expects path")
	}
	path, err := argString("playSound", args, 0)
	if err != nil {
		return nil, err
	}
	id, ok := helperSoundIDs[path]
	if !ok {
		v, e := builtinLoadSound([]*candy_evaluator.Value{{Kind: candy_evaluator.ValString, Str: path}})
		if e != nil {
			return nil, e
		}
		id = v.I64
		helperSoundIDs[path] = id
	}
	if len(args) >= 2 {
		vol := numOr0(args[1])
		_, _ = builtinSetSoundVolume([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: id}, {Kind: candy_evaluator.ValFloat, F64: vol}})
	}
	if len(args) >= 3 {
		pitch := numOr0(args[2])
		_, _ = builtinSetSoundPitch([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: id}, {Kind: candy_evaluator.ValFloat, F64: pitch}})
	}
	return builtinPlaySound([]*candy_evaluator.Value{{Kind: candy_evaluator.ValInt, I64: id}})
}

func numOr0(v *candy_evaluator.Value) float64 {
	if v == nil {
		return 0
	}
	if v.Kind == candy_evaluator.ValInt {
		return float64(v.I64)
	}
	if v.Kind == candy_evaluator.ValFloat {
		return v.F64
	}
	return 0
}

func builtinProjectile(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 8 {
		return nil, fmt.Errorf("projectile expects x,y,z,dirX,dirY,dirZ,speed,life")
	}
	id := nextProjectileID
	nextProjectileID++
	m := &candy_evaluator.Value{Kind: candy_evaluator.ValMap, StrMap: map[string]candy_evaluator.Value{}}
	setMapNum(m, "x", numOr0(args[0]))
	setMapNum(m, "y", numOr0(args[1]))
	setMapNum(m, "z", numOr0(args[2]))
	setMapNum(m, "dx", numOr0(args[3]))
	setMapNum(m, "dy", numOr0(args[4]))
	setMapNum(m, "dz", numOr0(args[5]))
	setMapNum(m, "speed", numOr0(args[6]))
	setMapNum(m, "life", numOr0(args[7]))
	projectiles[id] = m
	return vInt(id), nil
}
func builtinUpdateProjectiles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	dt := float64(rl.GetFrameTime())
	if len(args) == 1 {
		dt = numOr0(args[0])
	}
	for id, p := range projectiles {
		life := getMapNum(p, "life") - dt
		if life <= 0 {
			delete(projectiles, id)
			continue
		}
		setMapNum(p, "life", life)
		speed := getMapNum(p, "speed")
		setMapNum(p, "x", getMapNum(p, "x")+getMapNum(p, "dx")*speed*dt)
		setMapNum(p, "y", getMapNum(p, "y")+getMapNum(p, "dy")*speed*dt)
		setMapNum(p, "z", getMapNum(p, "z")+getMapNum(p, "dz")*speed*dt)
	}
	return null(), nil
}
func builtinDrawProjectiles(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	for _, p := range projectiles {
		rl.DrawCircle(int32(getMapNum(p, "x")), int32(getMapNum(p, "y")), 3, rl.White)
	}
	return null(), nil
}
func builtinProjectileHit(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return vBool(false), nil }
func builtinHitscan(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)       { return null(), nil }
func builtinLockOnNearest(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinDamage(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) >= 2 {
		setMapNum(args[0], "hp", getMapNum(args[0], "hp")-numOr0(args[1]))
	}
	return null(), nil
}
func builtinHeal(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) >= 2 {
		setMapNum(args[0], "hp", getMapNum(args[0], "hp")+numOr0(args[1]))
	}
	return null(), nil
}
func builtinAlive(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return vBool(false), nil
	}
	return vBool(getMapNum(args[0], "hp") > 0), nil
}
func builtinTeam(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 || args[0] == nil || args[0].Kind != candy_evaluator.ValMap {
		return null(), nil
	}
	if v, ok := args[0].StrMap["team"]; ok {
		vv := v
		return &vv, nil
	}
	return null(), nil
}
func builtinSetTeam(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 || args[0] == nil || args[0].Kind != candy_evaluator.ValMap {
		return null(), nil
	}
	args[0].StrMap["team"] = *args[1]
	return null(), nil
}

func builtinSpawn(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("spawn expects prefab,x,y[,z]")
	}
	m := map[string]candy_evaluator.Value{
		"prefab": *args[0],
		"x":      {Kind: candy_evaluator.ValFloat, F64: numOr0(args[1])},
		"y":      {Kind: candy_evaluator.ValFloat, F64: numOr0(args[2])},
	}
	if len(args) > 3 {
		m["z"] = candy_evaluator.Value{Kind: candy_evaluator.ValFloat, F64: numOr0(args[3])}
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValMap, StrMap: m}, nil
}
func builtinDespawn(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinSpawnAtMarker(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("spawnAtMarker expects marker,prefab")
	}
	return builtinSpawn([]*candy_evaluator.Value{args[1], {Kind: candy_evaluator.ValFloat, F64: 0}, {Kind: candy_evaluator.ValFloat, F64: 0}})
}
func builtinPoolCreate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextPoolID
	nextPoolID++
	helperPools[id] = []candy_evaluator.Value{}
	return vInt(id), nil
}
func builtinPoolGet(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("poolGet expects poolId,x,y[,z]")
	}
	return builtinSpawn([]*candy_evaluator.Value{{Kind: candy_evaluator.ValString, Str: "pool"}, args[1], args[2], args[3]})
}
func builtinPoolRelease(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinWaveCreate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextEnemyWaveID
	nextEnemyWaveID++
	enemyWaves[id] = []candy_evaluator.Value{}
	return vInt(id), nil
}
func builtinWaveAdd(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)   { return null(), nil }
func builtinWaveStart(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinWaveDone(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)  { return vBool(true), nil }

func builtinStateMachine(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id := nextSMID
	nextSMID++
	stateMachines[id] = map[string]*candy_evaluator.Value{}
	return vInt(id), nil
}
func builtinStateAdd(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("stateAdd expects sm,name,callbackMap")
	}
	id, _ := argInt("stateAdd", args, 0)
	name, _ := argString("stateAdd", args, 1)
	stateMachines[id][name] = args[2]
	return null(), nil
}
func builtinStateSet(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argInt("stateSet", args, 0)
	name, _ := argString("stateSet", args, 1)
	activeState[id] = name
	return null(), nil
}
func builtinStateUpdate(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinPatrol(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)     { return vInt(1), nil }
func builtinChase(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)      { return vInt(1), nil }
func builtinFlee(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)       { return vInt(1), nil }
func builtinWander(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)     { return vInt(1), nil }
func builtinLineOfSight(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return vBool(true), nil }
func builtinCanSee(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)      { return vBool(true), nil }

func builtinButton(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("button expects id,x,y,w,h,label")
	}
	x, y, w, h := numOr0(args[1]), numOr0(args[2]), numOr0(args[3]), numOr0(args[4])
	hover := rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)))
	if hover {
		rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), rl.Green)
	} else {
		rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), rl.Gray)
	}
	return vBool(hover && rl.IsMouseButtonPressed(rl.MouseButtonLeft)), nil
}
func builtinSlider(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 7 {
		return nil, fmt.Errorf("slider expects id,x,y,w,min,max,value")
	}
	return args[6], nil
}
func builtinHealthBar(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return nil, fmt.Errorf("healthBar expects x,y,w,h,current,max")
	}
	x, y, w, h, cur, maxv := numOr0(args[0]), numOr0(args[1]), numOr0(args[2]), numOr0(args[3]), numOr0(args[4]), math.Max(1, numOr0(args[5]))
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.DarkGray)
	rl.DrawRectangle(int32(x), int32(y), int32(w*(cur/maxv)), int32(h), rl.Green)
	return null(), nil
}
func builtinFloatingText(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinMinimap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)       { return null(), nil }
func builtinQuestAdd(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("questAdd expects id,title,description")
	}
	id, _ := argString("questAdd", args, 0)
	quests[id] = map[string]candy_evaluator.Value{"title": *args[1], "description": *args[2], "state": {Kind: candy_evaluator.ValString, Str: "active"}}
	return null(), nil
}
func builtinQuestComplete(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argString("questComplete", args, 0)
	if q, ok := quests[id]; ok {
		q["state"] = candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: "completed"}
	}
	return null(), nil
}
func builtinQuestState(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	id, _ := argString("questState", args, 0)
	if q, ok := quests[id]; ok {
		v := q["state"]
		return &v, nil
	}
	return &candy_evaluator.Value{Kind: candy_evaluator.ValString, Str: "missing"}, nil
}
func builtinQuestStep(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }

func builtinTag(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 2 {
		return null(), nil
	}
	n, _ := argString("tag", args, 1)
	tags[n] = append(tags[n], *args[0])
	return null(), nil
}
func builtinUntag(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinWithTag(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	n, _ := argString("withTag", args, 0)
	return &candy_evaluator.Value{Kind: candy_evaluator.ValArray, Elems: tags[n]}, nil
}
func builtinDistance2D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return vFloat(0), nil
	}
	dx := numOr0(args[2]) - numOr0(args[0])
	dy := numOr0(args[3]) - numOr0(args[1])
	return vFloat(math.Sqrt(dx*dx + dy*dy)), nil
}
func builtinDistance3D(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 6 {
		return vFloat(0), nil
	}
	dx := numOr0(args[3]) - numOr0(args[0])
	dy := numOr0(args[4]) - numOr0(args[1])
	dz := numOr0(args[5]) - numOr0(args[2])
	return vFloat(math.Sqrt(dx*dx + dy*dy + dz*dz)), nil
}
func builtinAngleTo(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return vFloat(0), nil
	}
	return vFloat(math.Atan2(numOr0(args[3])-numOr0(args[1]), numOr0(args[2])-numOr0(args[0])) * 180 / math.Pi), nil
}
func builtinLerp(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return vFloat(0), nil
	}
	a, b, t := numOr0(args[0]), numOr0(args[1]), numOr0(args[2])
	return vFloat(a + (b-a)*t), nil
}
func builtinRemap(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 5 {
		return vFloat(0), nil
	}
	v, inA, inB, outA, outB := numOr0(args[0]), numOr0(args[1]), numOr0(args[2]), numOr0(args[3]), numOr0(args[4])
	if inB == inA {
		return vFloat(outA), nil
	}
	t := (v - inA) / (inB - inA)
	return vFloat(outA + (outB-outA)*t), nil
}
func builtinChance(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 {
		return vBool(false), nil
	}
	return vBool(float64(rl.GetRandomValue(0, 10000))/100.0 <= numOr0(args[0])), nil
}
func builtinRandomPointInCircle(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 3 {
		return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: 0}, "y": {Kind: candy_evaluator.ValFloat, F64: 0}}), nil
	}
	cx, cy, r := numOr0(args[0]), numOr0(args[1]), numOr0(args[2])
	ang := float64(rl.GetRandomValue(0, 1000)) / 1000.0 * 2 * math.Pi
	rad := math.Sqrt(float64(rl.GetRandomValue(0, 1000))/1000.0) * r
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: cx + math.Cos(ang)*rad}, "y": {Kind: candy_evaluator.ValFloat, F64: cy + math.Sin(ang)*rad}}), nil
}
func builtinRandomPointInSphere(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 4 {
		return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: 0}, "y": {Kind: candy_evaluator.ValFloat, F64: 0}, "z": {Kind: candy_evaluator.ValFloat, F64: 0}}), nil
	}
	cx, cy, cz, r := numOr0(args[0]), numOr0(args[1]), numOr0(args[2]), numOr0(args[3])
	u := float64(rl.GetRandomValue(0, 1000)) / 1000.0
	v := float64(rl.GetRandomValue(0, 1000)) / 1000.0
	w := float64(rl.GetRandomValue(0, 1000)) / 1000.0
	theta := 2 * math.Pi * u
	phi := math.Acos(2*v - 1)
	rr := math.Cbrt(w) * r
	return vMap(map[string]candy_evaluator.Value{"x": {Kind: candy_evaluator.ValFloat, F64: cx + rr*math.Sin(phi)*math.Cos(theta)}, "y": {Kind: candy_evaluator.ValFloat, F64: cy + rr*math.Sin(phi)*math.Sin(theta)}, "z": {Kind: candy_evaluator.ValFloat, F64: cz + rr*math.Cos(phi)}}), nil
}
func builtinNetID(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) < 1 || args[0] == nil {
		return vInt(0), nil
	}
	if id, ok := entityNetID[args[0]]; ok {
		return vInt(id), nil
	}
	id := nextNetID
	nextNetID++
	entityNetID[args[0]] = id
	return vInt(id), nil
}
func builtinSetNetOwner(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) >= 2 {
		setMapNum(args[0], "netOwner", numOr0(args[1]))
	}
	return null(), nil
}
func builtinSnapshot(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) {
	if len(args) >= 1 {
		return args[0], nil
	}
	return null(), nil
}
func builtinInterpolateRemote(args []*candy_evaluator.Value) (*candy_evaluator.Value, error) { return null(), nil }
func builtinPredict(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)           { return null(), nil }
func builtinReconcile(args []*candy_evaluator.Value) (*candy_evaluator.Value, error)         { return null(), nil }
