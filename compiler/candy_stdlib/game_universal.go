package candy_stdlib

func init() {
	Modules["candy.scene"] = `
class Scene {
    var entities = []
    var root = null

    fun init() {
        if root == null { root = {"name": "root"}; }
    }

    fun add(entity) {
        entities.add(entity)
    }

    fun remove(entity) {
        entities.remove(entity)
    }

    fun update(dt) {
        // Intentionally lightweight; entities can update themselves in user code.
    }

    fun draw() {
        // Intentionally lightweight; entities can draw themselves in user code.
    }

    fun findByTag(tag) {
        return []
    }

    fun findByType(typeName) {
        return []
    }

    fun findInRadius(position, radius) {
        return []
    }

    fun load(path) {
        var s = Scene()
        s.sourcePath = path
        return s
    }
}

object SceneManager {
    var currentScene = null
    var stack = []

    fun change(scene) {
        currentScene = scene
        if currentScene != null {
            currentScene.init()
        }
    }

    fun push(scene) {
        stack.add(currentScene)
        currentScene = scene
        if currentScene != null {
            currentScene.init()
        }
    }

    fun pop() {
        currentScene = stack.pop()
    }

    fun update(dt) {
        if currentScene != null {
            currentScene.update(dt)
        }
    }

    fun draw() {
        if currentScene != null {
            currentScene.draw()
        }
    }
}
`
	Modules["candy.audio"] = `
fun playSound(path: String, volume: float = 1.0) {
    play(path)
}

fun playMusic(path: String, shouldLoop: bool = true) {
    // shouldLoop is accepted for API parity; runtime helper handles playback.
    music(path)
}

class AudioSource {
    var position = vec3(0, 0, 0)
    var volume = 1.0
    var pitch = 1.0

    fun play(path: String, shouldLoop: bool = false) {
        if shouldLoop {
            music(path)
        } else {
            play(path)
        }
    }
}

fun stopMusic() {}
fun pauseMusic() {}
fun resumeMusic() {}

class AudioBus {
    var name = "Bus"
    var volume = 1.0
    fun init(name = "Bus") { this.name = name; }
}

object AudioListener {
    var position = vec3(0, 0, 0)
    var forward = vec3(0, 0, 1)
    var up = vec3(0, 1, 0)
}

object Audio {
    var masterVolume = 1.0
}
`
	Modules["candy.input"] = `
val Input = InputMap()

object Gamepad {
    fun isConnected(index = 0) { return false; }
    fun getLeftStick(index = 0) { return vec2(0, 0); }
    fun getRightStick(index = 0) { return vec2(0, 0); }
    fun getLeftTrigger(index = 0) { return 0.0; }
    fun vibrate(index = 0, leftMotor = 0.5, rightMotor = 0.5, duration = 0.1) {}
}

object Touch {
    var touching = false
    fun getTouch(index = 0) { return {"position": vec2(0,0), "delta": vec2(0,0)}; }
}
`
	Modules["candy.resources"] = `
object Resources {
    var cache = {}
    var groups = {}
    var loadedCallbacks = []
    var hotReload = false
    var asyncQueue = []

    fun preload(paths) {
        for p in paths {
            cache[p] = p
        }
        for cb in loadedCallbacks {
            cb()
        }
    }

    fun get(path: String) {
        return cache[path]
    }

    fun onLoaded(callback) {
        loadedCallbacks.add(callback)
    }

    fun enableHotReload() {
        hotReload = true
    }

    fun group(name, paths) {
        groups[name] = paths
    }

    fun loadGroup(name, onComplete = null) {
        var paths = groups[name]
        if paths != null { preload(paths); }
        if onComplete != null { onComplete(); }
    }

    fun unloadGroup(name) {
        var paths = groups[name]
        if paths == null { return; }
        for p in paths {
            cache[p] = null
        }
    }

    fun loadAsync(path, onComplete = null) {
        asyncQueue.add({"path": path, "onComplete": onComplete})
    }

    fun update() {
        if asyncQueue.length == 0 { return 0; }
        var req = asyncQueue.pop()
        cache[req.path] = req.path
        if req.onComplete != null { req.onComplete(cache[req.path]); }
        return asyncQueue.length
    }
}
`
	Modules["candy.save"] = `
object Save {
    var slot = 0
    var slots = {}

    fun _slotName() {
        return "slot_" + slot
    }

    fun _slotMap() {
        var name = _slotName()
        if slots[name] == null { slots[name] = {}; }
        return slots[name]
    }

    fun set(key: String, value) {
        var m = _slotMap()
        m[key] = value
    }

    fun get(key: String, fallback = null) {
        var m = _slotMap()
        if m[key] != null { return m[key]; }
        return fallback
    }

    fun setObject(key: String, value) { set(key, value); }
    fun getObject(key: String, fallback = null) { return get(key, fallback); }

    fun setSlot(id) { slot = id; }

    fun saveJSON(path: String, value) { saveJson(path, value); }
    fun loadJSON(path: String) { return loadJson(path); }
    fun saveBinary(path: String, value) { saveJson(path, value); }
    fun loadBinary(path: String) { return loadJson(path); }
}

object Cloud {
    var data = {}
    fun save(key, value) { data[key] = value; }
    fun load(key, onComplete = null) {
        var value = data[key]
        if onComplete != null { onComplete(value); }
        return value
    }
}
`
	Modules["candy.debug"] = `
object Debug {
    fun drawLine(start: vec3, finish: vec3, color) {
        drawLine3D(start, finish, color)
    }

    fun showFPS() {
        drawFPS(10, 10)
    }

    fun drawBox(position, size, color) {
        drawCubeWires(position, size.x, size.y, size.z, color)
    }

    fun drawSphere(position, radius, color) {
        drawSphereWires(position, radius, 8, 8, color)
    }

    fun drawText(text, position, color) {
        // 3D text fallback to console for interpreter mode.
        print(text)
    }

    fun showMemory() {}
    fun showDrawCalls() {}
}

object Profiler {
    var marks = {}
    var stats = {}
    fun begin(name) { marks[name] = time.now(); }
    fun finish(name) {
        var started = marks[name]
        if started == null { return; }
        var dt = time.now() - started
        if stats[name] == null { stats[name] = {"count": 0, "total": 0.0, "avgTime": 0.0}; }
        stats[name].count = stats[name].count + 1
        stats[name].total = stats[name].total + dt
        stats[name].avgTime = stats[name].total / stats[name].count
    }
    fun getStats() { return stats; }
}

object Console {
    var commands = {}
    fun log(msg) {
        print(msg)
    }

    fun warn(msg) { print(msg); }
    fun error(msg) { print(msg); }
    fun registerCommand(name, fn) { commands[name] = fn; }
    fun show() {}
}
`
	Modules["candy.state"] = `
class StateMachine {
    var states = {}
    var currentState = ""

    fun addState(name, state) {
        states[name] = state
    }

    fun change(name) {
        currentState = name
    }

    fun update(dt) {
        var st = states[currentState]
        if st != null and st["update"] != null {
            st["update"](dt)
        }
    }
}
`
	Modules["candy.camera"] = `
import candy.math

class Camera2D {
    var position = vec2(0, 0)
    var zoom = 1.0
    var rotation = 0.0
    var bounds = null
    var viewport = null
    var smoothing = 1.0

    fun follow(target, smooth = 0.1) {
        if target == null { return; }
        var lerp = smooth
        if lerp < 0 { lerp = 0; }
        if lerp > 1 { lerp = 1; }
        position.x = position.x + (target.x - position.x) * lerp
        position.y = position.y + (target.y - position.y) * lerp
        if bounds != null {
            if position.x < bounds.x { position.x = bounds.x; }
            if position.y < bounds.y { position.y = bounds.y; }
            if bounds.w != null and position.x > bounds.x + bounds.w { position.x = bounds.x + bounds.w; }
            if bounds.h != null and position.y > bounds.y + bounds.h { position.y = bounds.y + bounds.h; }
        }
    }

    fun shake(duration, intensity) {
        if duration <= 0 or intensity <= 0 { return; }
        position.x = position.x + rand.float(-intensity, intensity)
        position.y = position.y + rand.float(-intensity, intensity)
    }

    fun screenToWorld(mousePos) {
        return vec2(
            mousePos.x / zoom + position.x,
            mousePos.y / zoom + position.y
        )
    }

    fun worldToScreen(worldPos) {
        return vec2(
            (worldPos.x - position.x) * zoom,
            (worldPos.y - position.y) * zoom
        )
    }
}

class FollowCamera2D extends Camera2D {
    var target = null
    var deadzone = null
    var lookAhead = vec2(0, 0)
    
    fun init(props = {}) {
        this.target = props.target
        this.smoothing = props.smoothing
        if this.smoothing == null { this.smoothing = 0.1; }
        this.deadzone = props.deadzone
        this.lookAhead = props.lookAhead
        if this.lookAhead == null { this.lookAhead = vec2(0, 0); }
    }
}
`
	Modules["candy.ai"] = `
import candy.math

class SteeringAgent {
    var position = vec2(0, 0)
    var velocity = vec2(0, 0)
    var acceleration = vec2(0, 0)
    var maxSpeed = 120.0
    var maxForce = 40.0

    fun init(props = {}) {
        if props.position != null { position = props.position; }
        if props.velocity != null { velocity = props.velocity; }
        if props.maxSpeed != null { maxSpeed = props.maxSpeed; }
        if props.maxForce != null { maxForce = props.maxForce; }
    }

    fun clampVec(v, m) {
        var len = math.sqrt(v.x * v.x + v.y * v.y)
        if len <= m or len <= 0.000001 { return v; }
        return vec2(v.x / len * m, v.y / len * m)
    }

    fun seek(target) {
        var desired = vec2(target.x - position.x, target.y - position.y)
        var dlen = math.sqrt(desired.x * desired.x + desired.y * desired.y)
        if dlen <= 0.000001 { return vec2(0, 0); }
        desired.x = desired.x / dlen * maxSpeed
        desired.y = desired.y / dlen * maxSpeed
        return clampVec(vec2(desired.x - velocity.x, desired.y - velocity.y), maxForce)
    }

    fun flee(target) {
        var f = seek(target)
        return vec2(-f.x, -f.y)
    }

    fun arrive(target, radius = 64.0) {
        var desired = vec2(target.x - position.x, target.y - position.y)
        var dlen = math.sqrt(desired.x * desired.x + desired.y * desired.y)
        if dlen <= 0.000001 { return vec2(0, 0); }
        var speed = maxSpeed
        if dlen < radius { speed = maxSpeed * (dlen / radius); }
        desired.x = desired.x / dlen * speed
        desired.y = desired.y / dlen * speed
        return clampVec(vec2(desired.x - velocity.x, desired.y - velocity.y), maxForce)
    }

    fun wander(jitter = 20.0) {
        var w = vec2(rand.float(-jitter, jitter), rand.float(-jitter, jitter))
        return clampVec(w, maxForce)
    }

    fun pursuit(target) {
        if target == null or target.position == null { return vec2(0, 0); }
        var tv = vec2(0, 0)
        if target.velocity != null { tv = target.velocity; }
        var dx = target.position.x - position.x
        var dy = target.position.y - position.y
        var d = math.sqrt(dx * dx + dy * dy)
        var lookAhead = d / maxSpeed
        var predicted = vec2(target.position.x + tv.x * lookAhead, target.position.y + tv.y * lookAhead)
        return seek(predicted)
    }

    fun evasion(target) {
        var f = pursuit(target)
        return vec2(-f.x, -f.y)
    }

    fun applyForce(force) {
        acceleration.x = acceleration.x + force.x
        acceleration.y = acceleration.y + force.y
    }

    fun update(dt) {
        velocity.x = velocity.x + acceleration.x * dt
        velocity.y = velocity.y + acceleration.y * dt
        velocity = clampVec(velocity, maxSpeed)
        position.x = position.x + velocity.x * dt
        position.y = position.y + velocity.y * dt
        acceleration = vec2(0, 0)
    }
}

class BehaviorNode {
    fun tick(dt) { return "SUCCESS"; }
}

class Selector extends BehaviorNode {
    var children = []
    fun init(children) { this.children = children; }
    fun tick(dt) {
        for c in children {
            var res = c.tick(dt)
            if res != "FAILURE" { return res; }
        }
        return "FAILURE"
    }
}

class Sequence extends BehaviorNode {
    var children = []
    fun init(children) { this.children = children; }
    fun tick(dt) {
        for c in children {
            var res = c.tick(dt)
            if res != "SUCCESS" { return res; }
        }
        return "SUCCESS"
    }
}

class Action extends BehaviorNode {
    var fn = null
    fun init(fn) { this.fn = fn; }
    fun tick(dt) { return fn(); }
}

class Condition extends BehaviorNode {
    var fn = null
    fun init(fn) { this.fn = fn; }
    fun tick(dt) {
        if fn() { return "SUCCESS"; }
        return "FAILURE";
    }
}

class Parallel extends BehaviorNode {
    var children = []
    var successPolicy = "ALL" // or "ANY"
    fun init(children, policy = "ALL") { this.children = children; this.successPolicy = policy; }
    fun tick(dt) {
        var successes = 0
        for c in children {
            var res = c.tick(dt)
            if res == "SUCCESS" { successes = successes + 1; }
            if res == "FAILURE" and successPolicy == "ALL" { return "FAILURE"; }
        }
        if successPolicy == "ANY" and successes > 0 { return "SUCCESS"; }
        if successPolicy == "ALL" and successes == len(children) { return "SUCCESS"; }
        return "RUNNING"
    }
}

class Inverter extends BehaviorNode {
    var child = null
    fun init(child) { this.child = child; }
    fun tick(dt) {
        var res = child.tick(dt)
        if res == "SUCCESS" { return "FAILURE"; }
        if res == "FAILURE" { return "SUCCESS"; }
        return res
    }
}

class Repeater extends BehaviorNode {
    var child = null
    var count = -1
    var current = 0
    fun init(child, count = -1) { this.child = child; this.count = count; }
    fun tick(dt) {
        if count != -1 and current >= count { return "SUCCESS"; }
        var res = child.tick(dt)
        if res != "RUNNING" { current = current + 1; }
        return "RUNNING"
    }
}

class Blackboard {
    var data = {}
    fun set(key, value) { data[key] = value; }
    fun get(key, fallback = null) {
        if data[key] == null { return fallback; }
        return data[key]
    }
}

class BehaviorTree {
    var root = null
    var blackboard = null
    fun init(root, bb = null) { 
        this.root = root; 
        this.blackboard = bb;
        if this.blackboard == null { this.blackboard = Blackboard(); }
    }
    fun tick(dt) {
        if root != null { root.tick(dt); }
    }
}
`
	Modules["candy.game3d"] = `
import candy.math
import candy.3d
import candy.physics3d
import candy.input

class ThirdPersonRig extends OrbitCamera {
    var follow = null
    var followHeight = 1.0
    var sensitivity = 0.0035
    var desiredYaw = 0.0
    var desiredPitch = 0.35
    var minDistance = 3.0
    var maxDistance = 25.0
    var zoomSpeed = 1.2
    var targetSmoothing = 10.0
    var rotSmoothing = 16.0
    var positionSmoothing = 14.0
    var smoothTarget = vec3(0, 0, 0)

    fun init(target = null) {
        follow = target
        target = vec3(0, 0, 0)
        if follow != null { target = follow.position; }
        smoothTarget = vec3(target.x, target.y + followHeight, target.z)
        this.target = smoothTarget
        distance = 8.0
        yaw = 0.5
        pitch = 0.35
        desiredYaw = yaw
        desiredPitch = pitch
    }

    fun update(dt) {
        if dt > 0.05 { dt = 0.05; }
        var t = target
        if follow != null { t = follow.position; }

        var md = getMouseDelta()
        desiredYaw -= md.x * sensitivity
        desiredPitch += md.y * sensitivity
        if desiredPitch > 1.2 { desiredPitch = 1.2; }
        if desiredPitch < -1.0 { desiredPitch = -1.0; }

        var wheel = getMouseWheelMove()
        if wheel != 0 { distance -= wheel * zoomSpeed; }
        if distance < minDistance { distance = minDistance; }
        if distance > maxDistance { distance = maxDistance; }

        var rLerp = rotSmoothing * dt
        if rLerp > 1 { rLerp = 1; }
        yaw = yaw + (desiredYaw - yaw) * rLerp
        pitch = pitch + (desiredPitch - pitch) * rLerp

        var wantedTarget = vec3(t.x, t.y + followHeight, t.z)
        var tLerp = targetSmoothing * dt
        if tLerp > 1 { tLerp = 1; }
        smoothTarget.x = smoothTarget.x + (wantedTarget.x - smoothTarget.x) * tLerp
        smoothTarget.y = smoothTarget.y + (wantedTarget.y - smoothTarget.y) * tLerp
        smoothTarget.z = smoothTarget.z + (wantedTarget.z - smoothTarget.z) * tLerp
        this.target = smoothTarget

        var cp = math.cos(pitch)
        var wantedPos = vec3(
            this.target.x + math.sin(yaw) * cp * distance,
            this.target.y + math.sin(pitch) * distance + 1.0,
            this.target.z + math.cos(yaw) * cp * distance
        )
        var pLerp = positionSmoothing * dt
        if pLerp > 1 { pLerp = 1; }
        position.x = position.x + (wantedPos.x - position.x) * pLerp
        position.y = position.y + (wantedPos.y - position.y) * pLerp
        position.z = position.z + (wantedPos.z - position.z) * pLerp
    }
}

class PlatformWorld3D {
    var boxes = []
    var groundY = 0.0

    fun addBox(cx, cy, cz, w, h, d) {
        boxes.add([cx, cy, cz, w, h, d])
    }

    fun solveAABB(position, size, velocity) {
        var halfW = size.x * 0.5
        var halfH = size.y * 0.5
        var halfD = size.z * 0.5
        var onGround = false

        if position.y - halfH <= groundY {
            position.y = groundY + halfH
            if velocity.y < 0 { velocity.y = 0; }
            onGround = true
        }

        for b in boxes {
            var cx = b[0]; var cy = b[1]; var cz = b[2]
            var w = b[3]; var h = b[4]; var d = b[5]
            var bhw = w * 0.5; var bhh = h * 0.5; var bhd = d * 0.5

            var dx = position.x - cx
            var dy = position.y - cy
            var dz = position.z - cz
            var adx = dx; if adx < 0 { adx = -adx; }
            var ady = dy; if ady < 0 { ady = -ady; }
            var adz = dz; if adz < 0 { adz = -adz; }

            var ox = (halfW + bhw) - adx
            var oy = (halfH + bhh) - ady
            var oz = (halfD + bhd) - adz
            if ox <= 0 or oy <= 0 or oz <= 0 { continue; }

            if oy <= ox and oy <= oz {
                if dy >= 0 {
                    position.y = cy + bhh + halfH
                    if velocity.y < 0 { velocity.y = 0; }
                    onGround = true
                } else {
                    position.y = cy - bhh - halfH
                    if velocity.y > 0 { velocity.y = 0; }
                }
            } else if ox <= oz {
                if dx >= 0 { position.x = position.x + ox; } else { position.x = position.x - ox; }
                velocity.x = 0
            } else {
                if dz >= 0 { position.z = position.z + oz; } else { position.z = position.z - oz; }
                velocity.z = 0
            }
        }
        return onGround
    }
}
`
}
