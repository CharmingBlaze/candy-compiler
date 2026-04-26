package candy_stdlib

func init() {
	Modules["candy.3d"] = `
import candy.math

class Entity3D {
    var position = vec3(0, 0, 0)
    var rotation = vec3(0, 0, 0)
    var scale = vec3(1, 1, 1)
    var parent = null
    var children = []
    
    fun addChild(child) {
        child.parent = this
        children.add(child)
    }

    fun removeChild(child) {
        children.remove(child)
        child.parent = null
    }

    fun lookAt(target: vec3) {
        var dx = target.x - position.x
        var dz = target.z - position.z
        if dx != 0 or dz != 0 {
            rotation.y = math.atan2(dx, dz)
        }
    }

    fun translate(offset: vec3) {
        position = vector3Add(position, offset)
    }

    fun rotate(axis: vec3, angle: float) {
        rotation.x = rotation.x + axis.x * angle
        rotation.y = rotation.y + axis.y * angle
        rotation.z = rotation.z + axis.z * angle
    }

    fun update(dt) {
        for c in children { c.update(dt); }
    }
    
    fun draw() {
        for c in children { c.draw(); }
    }
}

class Model extends Entity3D {
    var model = null
    var tint = COLOR_WHITE

    fun init(path: String = "") {
        if path != "" and fileExists(path) {
            model = loadModel(path)
        }
    }

    fun draw() {
        if model == null {
            drawCube(position.x, position.y, position.z, scale.x, scale.y, scale.z, "red")
            return
        }
        drawModelEx(model, position, vec3(0, 1, 0), rotation.y, scale, tint)
    }
}

class Mesh extends Entity3D {
    var vertices = []
    var normals = []
    var uvs = []
    var indices = []
    var material = null

    fun build() {
        // Upload to GPU logic
    }

    object sphere {
        fun create(radius: float, segments: int) {
            // genMeshSphere
        }
    }
    
    object cube {
        fun create(size: vec3) {
            // genMeshCube
        }
    }
}

class Camera3D {
    var position = vec3(0, 10, 10)
    var target = vec3(0, 0, 0)
    var up = vec3(0, 1, 0)
    var fov = 60.0
    var near = 0.1
    var far = 1000.0

    fun lookAt(t: vec3) {
        target = t
    }

    fun shake(duration: float, intensity: float) {
        if duration <= 0 or intensity <= 0 { return; }
        var jx = rand.float(-intensity, intensity)
        var jy = rand.float(-intensity, intensity)
        var jz = rand.float(-intensity, intensity)
        position.x = position.x + jx
        position.y = position.y + jy
        position.z = position.z + jz
    }
}

class FirstPersonCamera extends Camera3D {
    var sensitivity = 0.1
    var speed = 5.0
    
    fun update(dt) {
        var md = getMouseDelta()
        rotation.y = rotation.y - md.x * sensitivity
        rotation.x = rotation.x + md.y * sensitivity
        if rotation.x > 1.2 { rotation.x = 1.2; }
        if rotation.x < -1.2 { rotation.x = -1.2; }

        var fx = math.sin(rotation.y)
        var fz = math.cos(rotation.y)
        var rx = fz
        var rz = -fx
        var moveX = 0.0
        var moveZ = 0.0
        if isKeyDown(KEY_W) { moveX = moveX + fx; moveZ = moveZ + fz; }
        if isKeyDown(KEY_S) { moveX = moveX - fx; moveZ = moveZ - fz; }
        if isKeyDown(KEY_D) { moveX = moveX + rx; moveZ = moveZ + rz; }
        if isKeyDown(KEY_A) { moveX = moveX - rx; moveZ = moveZ - rz; }

        var mlen = math.sqrt(moveX * moveX + moveZ * moveZ)
        if mlen > 0.0001 {
            moveX = moveX / mlen
            moveZ = moveZ / mlen
        }

        position.x = position.x + moveX * speed * dt
        position.z = position.z + moveZ * speed * dt
        target.x = position.x + math.sin(rotation.y) * math.cos(rotation.x)
        target.y = position.y + math.sin(rotation.x)
        target.z = position.z + math.cos(rotation.y) * math.cos(rotation.x)
    }
}

class OrbitCamera extends Camera3D {
    var distance = 10.0
    var yaw = 0.0
    var pitch = 0.35
    var sensitivity = 0.004
    var zoomSpeed = 1.2
    var minDistance = 3.0
    var maxDistance = 30.0

    fun update(dt) {
        var d = getMouseDelta()
        yaw -= d.x * sensitivity
        pitch += d.y * sensitivity
        if pitch > 1.2 { pitch = 1.2; }
        if pitch < -1.0 { pitch = -1.0; }

        var wheel = getMouseWheelMove()
        if wheel != 0 {
            distance -= wheel * zoomSpeed
        }
        if distance < minDistance { distance = minDistance; }
        if distance > maxDistance { distance = maxDistance; }

        var cp = math.cos(pitch)
        position.x = target.x + math.sin(yaw) * cp * distance
        position.y = target.y + math.sin(pitch) * distance + 1.0
        position.z = target.z + math.cos(yaw) * cp * distance
    }
}

class ThirdPersonCamera extends Camera3D {
    var offset = vec3(0, 2, -5)
    var smoothing = 8.0
    var follow = null
    
    fun update(dt) {
        var t = target
        if follow != null {
            t = follow.position
        }
        var lerp = smoothing * dt
        if lerp > 1 { lerp = 1; }
        position.x = position.x + (t.x + offset.x - position.x) * lerp
        position.y = position.y + (t.y + offset.y - position.y) * lerp
        position.z = position.z + (t.z + offset.z - position.z) * lerp
        target = t
    }
}

class Light extends Entity3D {
    var color = COLOR_WHITE
    var intensity = 1.0
}

class DirectionalLight extends Light {
    var direction = vec3(-1, -1, -1)
}

class PointLight extends Light {
    var radius = 10.0
}

class SpotLight extends Light {
    var direction = vec3(0, -1, 0)
    var angle = 30.0
}

class AmbientLight extends Light {}

class Animation3D {
    var model = null
    var anims = []
    var current = 0
    var frame = 0

    fun init(path: String = "") {
        if path != "" and fileExists(path) {
            anims = loadModelAnimations(path)
        } else {
            anims = []
        }
    }

    fun play(name = "") {
        if len(anims) == 0 { return; }
        // Name-based clip lookup can be added when metadata is exposed.
        current = 0
        frame = 0
    }

    fun update(dt) {
        if len(anims) > 0 {
            updateModelAnimation(model, anims[current], frame)
            frame += 1
        }
    }
}

class Animator {
    var target = null
    var clips = {}
    var current = ""
    var speed = 1.0
    var paused = false
    var events = {}

    fun init(target = null) {
        this.target = target
    }

    fun play(name, loop = true) {
        current = name
        paused = false
    }

    fun stop() {
        current = ""
    }

    fun pause() { paused = true; }
    fun resume() { paused = false; }

    fun crossFade(from, to, duration = 0.2) {
        current = to
    }

    fun onEvent(name, callback) {
        events[name] = callback
    }
}

class AnimationStateMachine {
    var animator = null
    var states = {}
    var transitions = []
    var current = ""

    fun init(animator) {
        this.animator = animator
    }

    fun addState(name, clipName) {
        states[name] = clipName
        if current == "" { current = name; }
    }

    fun addTransition(from, to, condition) {
        transitions.add({"from": from, "to": to, "condition": condition})
    }

    fun update(dt) {
        for t in transitions {
            if t.from == current or t.from == "*" {
                if t.condition() {
                    current = t.to
                    if animator != null and states[current] != null {
                        animator.play(states[current], true)
                    }
                    return
                }
            }
        }
    }
}
`
}
