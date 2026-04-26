package candy_stdlib

func init() {
	Modules["candy.physics2d"] = `
import candy.math

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

    fun moveAndCollide(motion) {
        position.x = position.x + motion.x
        position.y = position.y + motion.y
        return {"position": position, "collided": false}
    }

    fun moveAndSlide(motion) {
        moveAndCollide(motion)
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
    }
    
    fun remove(body) {
        bodies.remove(body)
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
            var radius = 0.5
            if b.collider != null and b.collider.radius != null {
                radius = b.collider.radius
            }
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
            var r = 0.5
            if b.collider != null and b.collider.radius != null { r = b.collider.radius; }
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

    fun moveAndCollide(motion) {
        position.x = position.x + motion.x
        position.y = position.y + motion.y
        position.z = position.z + motion.z
        return {"position": position, "collided": false}
    }

    fun moveAndSlide(motion) {
        moveAndCollide(motion)
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
            var radius = 0.6
            if b.collider != null and b.collider.radius != null {
                radius = b.collider.radius
            }
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
