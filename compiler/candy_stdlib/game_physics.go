package candy_stdlib

func init() {
	Modules["candy.physics2d"] = `
import candy.math

fun _radius2d(body) {
    if body == null { return 0.5; }
    if body.collider == null { return 0.5; }
    if body.collider.radius != null { return body.collider.radius; }
    if body.collider.size != null {
        var sx = body.collider.size.x
        var sy = body.collider.size.y
        return math.max(sx, sy) * 0.5
    }
    if body.collider.width != null and body.collider.height != null {
        return math.max(body.collider.width, body.collider.height) * 0.5
    }
    return 0.5
}

fun _findCollision2d(world, selfBody) {
    if world == null or world.bodies == null { return null; }
    for other in world.bodies {
        if other == null or other == selfBody { continue; }
        if other.position == null { continue; }
        var dx = selfBody.position.x - other.position.x
        var dy = selfBody.position.y - other.position.y
        var rr = _radius2d(selfBody) + _radius2d(other)
        var d2 = dx * dx + dy * dy
        if d2 <= rr * rr {
            var len = math.sqrt(d2)
            var nx = 1.0
            var ny = 0.0
            if len > 0.00001 {
                nx = dx / len
                ny = dy / len
            }
            return {
                "body": other,
                "normal": vec2(nx, ny),
                "point": vec2(
                    selfBody.position.x - nx * _radius2d(selfBody),
                    selfBody.position.y - ny * _radius2d(selfBody)
                )
            }
        }
    }
    return null
}

class BoxCollider {
    var size = vec2(0, 0)
    var offset = vec2(0, 0)
    fun init(size, offset = vec2(0,0)) {
        this.size = size
        this.offset = offset
    }
}

class CircleCollider {
    var radius = 0.0
    var offset = vec2(0, 0)
    fun init(radius, offset = vec2(0,0)) {
        this.radius = radius
        this.offset = offset
    }
}

class CapsuleCollider {
    var width = 0.0
    var height = 0.0
    fun init(width, height) {
        this.width = width
        this.height = height
    }
}

class PolygonCollider {
    var points = []
    fun init(points) { this.points = points; }
}

class EdgeCollider {
    var points = []
    fun init(points) { this.points = points; }
}

class KinematicBody2D {
    var position = vec2(0, 0)
    var velocity = vec2(0, 0)
    var collider = null
    var world = null

    fun moveAndCollide(motion) {
        position.x = position.x + motion.x
        position.y = position.y + motion.y
        var hit = _findCollision2d(world, this)
        if hit == null {
            return {"position": position, "collided": false}
        }
        return {
            "position": position,
            "collided": true,
            "collider": hit.body,
            "normal": hit.normal,
            "point": hit.point
        }
    }

    fun moveAndSlide(motion) {
        var result = moveAndCollide(motion)
        if result.collided {
            var n = result.normal
            var dot = motion.x * n.x + motion.y * n.y
            var sx = motion.x - n.x * dot
            var sy = motion.y - n.y * dot
            position.x = position.x + sx
            position.y = position.y + sy
        }
        return result
    }
}

class StaticBody2D {
    var position = vec2(0, 0)
    var collider = null
}

class RigidBody2D {
    var position = vec2(0, 0)
    var velocity = vec2(0, 0)
    var mass = 1.0
    var gravity = vec2(0, 980)
    var collider = null

    fun applyForce(force) {
        velocity.x = velocity.x + force.x / mass
        velocity.y = velocity.y + force.y / mass
    }
}

class Physics2D {
    var bodies = []
    
    fun add(body) {
        bodies.add(body)
        if body != null { body.world = this; }
    }
    
    fun remove(body) {
        bodies.remove(body)
        if body != null and body.world == this { body.world = null; }
    }
    
    fun update(dt) {
        for b in bodies {
            if b.velocity != null {
                b.position.x = b.position.x + b.velocity.x * dt
                b.position.y = b.position.y + b.velocity.y * dt
            }
        }
    }
    
    fun raycast(startPos, endPos) {
        var hit = null
        var bestDist = infinity
        var dx = endPos.x - startPos.x
        var dy = endPos.y - startPos.y
        var segLen2 = dx * dx + dy * dy
        if segLen2 <= 0.000001 { return null; }
        for b in bodies {
            if b == null { continue; }
            if b.position == null { continue; }
            var bx = b.position.x
            var by = b.position.y
            var t = ((bx - startPos.x) * dx + (by - startPos.y) * dy) / segLen2
            if t < 0 { t = 0; }
            if t > 1 { t = 1; }
            var px = startPos.x + dx * t
            var py = startPos.y + dy * t
            var ex = bx - px
            var ey = by - py
            var dist2 = ex * ex + ey * ey
            var radius = _radius2d(b)
            if dist2 <= radius * radius {
                var along2 = (px - startPos.x) * (px - startPos.x) + (py - startPos.y) * (py - startPos.y)
                if along2 < bestDist {
                    bestDist = along2
                    hit = {"body": b, "point": vec2(px, py), "distance": math.sqrt(along2)}
                }
            }
        }
        return hit
    }

    fun overlapPoint(point) {
        for b in bodies {
            if b == null or b.position == null { continue; }
            var dx = b.position.x - point.x
            var dy = b.position.y - point.y
            var r = _radius2d(b)
            if dx * dx + dy * dy <= r * r { return b; }
        }
        return null
    }

    fun overlapArea(center, radius = 1.0) {
        var out = []
        var r2 = radius * radius
        for b in bodies {
            if b == null or b.position == null { continue; }
            var dx = b.position.x - center.x
            var dy = b.position.y - center.y
            if dx * dx + dy * dy <= r2 {
                out.add(b)
            }
        }
        return out
    }
}
`
	Modules["candy.physics3d"] = `
import candy.math

fun _radius3d(body) {
    if body == null { return 0.6; }
    if body.collider == null { return 0.6; }
    if body.collider.radius != null { return body.collider.radius; }
    if body.collider.size != null {
        return math.max(body.collider.size.x, math.max(body.collider.size.y, body.collider.size.z)) * 0.5
    }
    return 0.6
}

class BoxCollider3D {
    var size = vec3(1, 1, 1)
    fun init(size) { this.size = size; }
}

class SphereCollider3D {
    var radius = 1.0
    fun init(radius) { this.radius = radius; }
}

class CapsuleCollider3D {
    var radius = 0.5
    var height = 2.0
    fun init(radius, height) { this.radius = radius; this.height = height; }
}

class BoxCollider extends BoxCollider3D {}
class SphereCollider extends SphereCollider3D {}
class CapsuleCollider extends CapsuleCollider3D {}

class MeshCollider {
    var mesh = null
    fun init(mesh) { this.mesh = mesh; }
}

class ConvexCollider {
    var points = []
    fun init(points) { this.points = points; }
}

class KinematicBody3D {
    var position = vec3(0, 0, 0)
    var velocity = vec3(0, 0, 0)
    var collider = null
    var world = null

    fun moveAndCollide(motion) {
        position.x = position.x + motion.x
        position.y = position.y + motion.y
        position.z = position.z + motion.z
        if world == null {
            return {"position": position, "collided": false}
        }
        for b in world.bodies {
            if b == null or b == this or b.position == null { continue; }
            var dx = position.x - b.position.x
            var dy = position.y - b.position.y
            var dz = position.z - b.position.z
            var rr = _radius3d(this) + _radius3d(b)
            var d2 = dx * dx + dy * dy + dz * dz
            if d2 <= rr * rr {
                var len = math.sqrt(d2)
                var nx = 1.0; var ny = 0.0; var nz = 0.0
                if len > 0.00001 {
                    nx = dx / len; ny = dy / len; nz = dz / len
                }
                return {
                    "position": position,
                    "collided": true,
                    "collider": b,
                    "normal": vec3(nx, ny, nz),
                    "point": vec3(position.x - nx * _radius3d(this), position.y - ny * _radius3d(this), position.z - nz * _radius3d(this))
                }
            }
        }
        return {"position": position, "collided": false}
    }

    fun moveAndSlide(motion) {
        var result = moveAndCollide(motion)
        if result.collided {
            var n = result.normal
            var dot = motion.x * n.x + motion.y * n.y + motion.z * n.z
            position.x = position.x + (motion.x - n.x * dot)
            position.y = position.y + (motion.y - n.y * dot)
            position.z = position.z + (motion.z - n.z * dot)
        }
        return result
    }
}

class CharacterController3D {
    var position = vec3(0, 0, 0)
    var velocity = vec3(0, 0, 0)
    var gravity = 20.0
    var jumpPower = 8.0
    var onGround = false

    fun move(dir: vec3, speed: float, dt: float) {
        // Accept either native vec values or raylib-style {x,y,z} maps.
        var dx = dir.x
        var dz = dir.z
        var len = math.sqrt(dx * dx + dz * dz)
        if len > 0.0001 {
            dx = dx / len
            dz = dz / len
        } else {
            dx = 0
            dz = 0
        }
        velocity.x = dx * speed
        velocity.z = dz * speed
        velocity.y = velocity.y - gravity * dt

        position.x = position.x + velocity.x * dt
        position.y = position.y + velocity.y * dt
        position.z = position.z + velocity.z * dt

        // Simple floor clamp so demo scripts remain playable.
        if position.y <= 0.5 {
            position.y = 0.5
            velocity.y = 0
            onGround = true
        } else {
            onGround = false
        }
    }

    fun jump() {
        if onGround {
            velocity.y = jumpPower
            onGround = false
        }
    }
}

class StaticBody3D {
    var position = vec3(0, 0, 0)
    var collider = null
}

class RigidBody3D {
    var position = vec3(0, 0, 0)
    var velocity = vec3(0, 0, 0)
    var rotation = vec3(0, 0, 0)
    var angularVelocity = vec3(0, 0, 0)
    var mass = 1.0
    var collider = null
    
    fun applyForce(force) {
        velocity.x = velocity.x + force.x / mass
        velocity.y = velocity.y + force.y / mass
        velocity.z = velocity.z + force.z / mass
    }
    fun applyTorque(torque) {
        angularVelocity.x = angularVelocity.x + torque.x / mass
        angularVelocity.y = angularVelocity.y + torque.y / mass
        angularVelocity.z = angularVelocity.z + torque.z / mass
    }
}

class Physics3D {
    var gravity = vec3(0, -20, 0)
    var bodies = []

    fun add(body) {
        bodies.add(body)
        if body != null { body.world = this; }
    }

    fun remove(body) {
        bodies.remove(body)
        if body != null and body.world == this { body.world = null; }
    }

    fun update(dt) {
        for b in bodies {
            if b.velocity != null {
                b.position.x = b.position.x + b.velocity.x * dt
                b.position.y = b.position.y + b.velocity.y * dt
                b.position.z = b.position.z + b.velocity.z * dt
            }
            if b.angularVelocity != null and b.rotation != null {
                b.rotation.x = b.rotation.x + b.angularVelocity.x * dt
                b.rotation.y = b.rotation.y + b.angularVelocity.y * dt
                b.rotation.z = b.rotation.z + b.angularVelocity.z * dt
            }
        }
    }

    fun raycast(startPos, endPos, ignoreBody = null) {
        var hit = null
        var bestDist = infinity
        var dx = endPos.x - startPos.x
        var dy = endPos.y - startPos.y
        var dz = endPos.z - startPos.z
        var segLen2 = dx * dx + dy * dy + dz * dz
        if segLen2 <= 0.000001 { return null; }
        for b in bodies {
            if b == null { continue; }
            if b == ignoreBody { continue; }
            if b.position == null { continue; }
            var bx = b.position.x
            var by = b.position.y
            var bz = b.position.z
            var t = ((bx - startPos.x) * dx + (by - startPos.y) * dy + (bz - startPos.z) * dz) / segLen2
            if t < 0 { t = 0; }
            if t > 1 { t = 1; }
            var px = startPos.x + dx * t
            var py = startPos.y + dy * t
            var pz = startPos.z + dz * t
            var ex = bx - px
            var ey = by - py
            var ez = bz - pz
            var dist2 = ex * ex + ey * ey + ez * ez
            var radius = _radius3d(b)
            if dist2 <= radius * radius {
                var along2 = (px - startPos.x) * (px - startPos.x) + (py - startPos.y) * (py - startPos.y) + (pz - startPos.z) * (pz - startPos.z)
                if along2 < bestDist {
                    bestDist = along2
                    hit = {"body": b, "point": vec3(px, py, pz), "distance": math.sqrt(along2)}
                }
            }
        }
        return hit
    }

    fun shapeCast(shape, startPos, endPos, ignoreBody = null) {
        // Minimal compatibility implementation uses raycast centerline.
        return raycast(startPos, endPos, ignoreBody)
    }
}
`
}
