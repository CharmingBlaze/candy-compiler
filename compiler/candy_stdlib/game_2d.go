package candy_stdlib

func init() {
	Modules["candy.2d"] = `
import candy.math

class Entity2D {
    var position = vec2(0, 0)
    var rotation = 0.0
    var scale = vec2(1, 1)
    var parent = null
    var children = []
    var visible = true

    fun addChild(child) {
        child.parent = this
        children.add(child)
    }

    fun removeChild(child) {
        children.remove(child)
        child.parent = null
    }

    fun update(dt) {}
    fun draw() {}
    
    fun destroy() {
        if parent != null {
            parent.removeChild(this)
        }
    }
}

class Sprite extends Entity2D {
    var texture = null
    var origin = vec2(0.5, 0.5)
    var flipX = false
    var flipY = false
    var flip = {"x": false, "y": false}
    var tint = COLOR_WHITE
    var alpha = 1.0

    fun init(path: String = "") {
        if path == "" { return; }
        if fileExists(path) {
            texture = loadTexture(path)
        }
    }

    fun draw() {
        if !visible or texture == null { return; }
        drawTextureEx(texture, position, rotation, scale.x, tint)
    }
}

class AnimatedSprite extends Sprite {
    var animations = {}
    var currentAnim = ""
    var frame = 0
    var timer = 0.0
    var speed = 0.1

    fun play(name: String) {
        currentAnim = name
        frame = 0
        timer = 0.0
    }

    fun stop() {
        currentAnim = ""
    }

    fun update(dt) {
        if currentAnim == "" { return; }
        var frames = animations[currentAnim]
        if frames == null or len(frames) == 0 { return; }
        timer += dt
        if timer >= speed {
            timer = 0.0
            frame += 1
            if frame >= len(frames) { frame = 0; }
        }
    }
}

class Tilemap extends Entity2D {
    var tileset = null
    var tileSize = 32
    var data = []
    var width = 0
    var height = 0
    var layers = {}

    fun init(path: String = "", w = 0, h = 0, ts = 32) {
        if path != "" and fileExists(path) {
            tileset = loadTexture(path)
        }
        width = w
        height = h
        tileSize = ts
        var baseLayer = {}
        baseLayer["visible"] = true
        baseLayer["data"] = []
        layers["base"] = baseLayer
    }

    fun draw() {
        if tileset == null { return; }
        for y in 0..height-1 {
            for x in 0..width-1 {
                var tid = data[y * width + x]
                if tid > 0 {
                    // draw tileset part
                }
            }
        }
    }

    fun setTile(x, y, tileID) {
        if x < 0 or x >= width or y < 0 or y >= height { return; }
        var i = y * width + x
        data[i] = tileID
    }

    fun getTile(x, y) {
        if x < 0 or x >= width or y < 0 or y >= height { return 0; }
        var i = y * width + x
        return data[i]
    }

    fun generateCollision() {
        var bodies = []
        for y in 0..height-1 {
            for x in 0..width-1 {
                if getTile(x, y) > 0 {
                    bodies.add({
                        "position": vec2((x + 0.5) * tileSize, (y + 0.5) * tileSize),
                        "collider": {"size": vec2(tileSize, tileSize)}
                    })
                }
            }
        }
        return bodies
    }

    fun worldToGrid(worldPos) {
        return vec2(math.floor(worldPos.x / tileSize), math.floor(worldPos.y / tileSize))
    }

    fun gridToWorld(x, y) {
        return vec2(x * tileSize, y * tileSize)
    }
}

class ParticleSystem2D extends Entity2D {
    var particles = []
    var maxParticles = 100
    var emissionRate = 10
    var lifetime = 1.0
    var tint = COLOR_WHITE

    fun emit() {
        if len(particles) >= maxParticles { return; }
        particles.add({
            "pos": vec2(position.x, position.y),
            "vel": vec2(rand.float(-30, 30), rand.float(-80, -20)),
            "life": lifetime
        })
    }

    fun emitBurst(cfg = {}) {
        var count = cfg.count
        if count == null { count = 12; }
        var i = 0
        while i < count {
            emit()
            i += 1
        }
    }

    fun update(dt) {
        var toEmit = emissionRate * dt
        var i = 0
        while i < toEmit {
            emit()
            i += 1
        }
        var particlesNext = []
        for p in particles {
            p.life = p.life - dt
            if p.life <= 0 { continue; }
            p.pos.x = p.pos.x + p.vel.x * dt
            p.pos.y = p.pos.y + p.vel.y * dt
            particlesNext.add(p)
        }
        particles = particlesNext
    }

    fun draw() {
        for p in particles {
            drawCircle(p.pos.x, p.pos.y, 2, tint)
        }
    }
}

class ParticleEmitter2D extends ParticleSystem2D {
    var running = false
    var rate = 20

    fun start() { running = true; }
    fun stop() { running = false; }

    fun update(dt) {
        if running {
            emissionRate = rate
        } else {
            emissionRate = 0
        }
        super.update(dt)
    }
}

object Particles {
    fun explosion(position) {
        var p = ParticleSystem2D()
        p.position = position
        p.emissionRate = 0
        p.emitBurst({"count": 40})
        return p
    }

    fun smoke(position) {
        var p = ParticleEmitter2D()
        p.position = position
        p.tint = COLOR_GRAY
        p.rate = 12
        p.start()
        return p
    }

    fun sparkles(position) {
        var p = ParticleSystem2D()
        p.position = position
        p.tint = COLOR_YELLOW
        p.emitBurst({"count": 24})
        return p
    }

    fun fire(position) {
        var p = ParticleEmitter2D()
        p.position = position
        p.tint = COLOR_ORANGE
        p.rate = 25
        p.start()
        return p
    }
}
`
}
