package candy_stdlib

func init() {
	Modules["candy.game"] = `
import candy.2d
import candy.3d
import candy.physics2d
import candy.physics3d
import candy.ui
import candy.scene
import candy.audio
import candy.input
import candy.resources
import candy.save
import candy.debug
import candy.state
import candy.camera
import candy.ai
import candy.game3d
import candy.network

// High-level batteries-included facade so projects can start quickly
// while still exposing all lower-level systems for advanced users.
object Game2D {
    // One-line map return: multiline return { } has been misparsed as a block in some method bodies.
    fun createWorld() {
        return { "scene": Scene(), "physics": Physics2D(), "camera": Camera2D(), "ui": Canvas(), "state": StateMachine() }
    }

    fun tick(world, dt) {
        if world.scene != null { world.scene.update(dt); }
        if world.physics != null { world.physics.update(dt); }
    }
}

object Game3D {
    fun createWorld() {
        return { "scene": Scene(), "physics": Physics3D(), "camera": ThirdPersonRig(null), "ui": Canvas(), "state": StateMachine() }
    }

    fun tick(world, dt) {
        if world.scene != null { world.scene.update(dt); }
        if world.physics != null { world.physics.update(dt); }
        if world.camera != null and world.camera.update != null { world.camera.update(dt); }
    }
}

class MultiplayerSession {
    var isServer = false
    var server = null
    var client = null
    var handlers = {}

    fun host(port = 20000, maxPeers = 32, channels = 2) {
        isServer = true
        server = NetworkServer(port, maxPeers, channels)
    }

    fun join(host = "127.0.0.1", port = 20000, channels = 2) {
        isServer = false
        client = NetworkClient(channels)
        client.connect(host, port, channels)
    }

    fun on(eventName, callback) {
        handlers[eventName] = callback
        if server != null { server.on(eventName, callback); }
        if client != null { client.on(eventName, callback); }
    }

    fun send(eventName, payload) {
        if isServer {
            if server != null { server.broadcast(eventName, payload); }
        } else {
            if client != null { client.send(eventName, payload); }
        }
    }

    fun update(timeoutMs = 0) {
        if server != null { server.update(timeoutMs); }
        if client != null { client.update(timeoutMs); }
    }

    fun stop() {
        if server != null { server.stop(); server = null; }
        if client != null { client.stop(); client = null; }
    }
}

class App {
    var scene = Scene()
    var ui = Canvas()
    var state = StateMachine()
    var resources = Resources

    fun update(dt) {
        scene.update(dt)
    }

    fun draw() {
        scene.draw()
        ui.draw()
    }
}
`
}
