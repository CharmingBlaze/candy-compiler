package candy_stdlib

func init() {
	Modules["candy.network"] = `
import enet
import json

fun _netEncode(kind, name, payload) {
    return json.encode({"kind": kind, "name": name, "payload": payload})
}

fun _netDecode(text) {
    return json.decode(text)
}

class NetworkServer {
    var hostId = -1
    var peers = []
    var rpcCallbacks = {}
    var onJoin = null
    var onLeave = null
    var serverPort = 1234
    var serverMaxPeers = 32
    var serverChannels = 2

    fun init(port = 1234, maxPeers = 32, channels = 2) {
        serverPort = port
        serverMaxPeers = 32
        serverChannels = 2
        hostId = -1
    }

    fun start() {
        if hostId >= 0 {
            return
        }
        enet.init()
        var addr = enet.address("0.0.0.0", serverPort)
        hostId = enet.host_create(addr, serverMaxPeers, serverChannels, 0, 0)
        if hostId < 0 {
            print("Failed to create server host")
            return
        }
        print("Server started on port " + toString(serverPort))
    }

    fun stop() {
        if hostId >= 0 {
            enet.host_destroy(hostId)
        }
        hostId = -1
        enet.deinit()
    }

    fun _addPeer(peerId) {
        for p in peers {
            if p == peerId {
                return
            }
        }
        peers.add(peerId)
    }

    fun _removePeer(peerId) {
        var peersNext = []
        for p in peers {
            if p != peerId { peersNext.add(p); }
        }
        peers = peersNext
    }

    fun update(timeoutMs = 0) {
        if hostId < 0 {
            return
        }
        var ev = enet.host_service(hostId, timeoutMs)
        while ev != null and ev.type != enet.EVENT_NONE {
            if ev.type == enet.EVENT_CONNECT {
                _addPeer(ev.peerId)
                if onJoin != null { onJoin(ev.peerId, ev.address); }
            } else if ev.type == enet.EVENT_DISCONNECT {
                _removePeer(ev.peerId)
                if onLeave != null { onLeave(ev.peerId); }
            } else if ev.type == enet.EVENT_RECEIVE {
                handlePacket(ev.peerId, ev.packet.data)
            }
            ev = enet.host_service(hostId, 0)
        }
    }

    fun handlePacket(peerId, data) {
        var msg = _netDecode(data)
        if msg == null {
            return
        }
        if msg.kind == "rpc" {
            var method = msg.name
            var cb = rpcCallbacks[method]
            if cb != null {
                cb(peerId, msg.payload)
            }
        }
    }

    fun on(method, callback) {
        rpcCallbacks[method] = callback
    }

    fun send(peerId, method, payload) {
        if hostId < 0 {
            return
        }
        var data = _netEncode("rpc", method, payload)
        var pkt = enet.packet_create(data, enet.PACKET_RELIABLE)
        enet.peer_send(peerId, 0, pkt)
        enet.host_flush(hostId)
    }

    fun broadcast(method, payload) {
        for p in peers {
            send(p, method, payload)
        }
    }

    fun peerCount() {
        return len(peers)
    }
}

class NetworkClient {
    var hostId = -1
    var peerId = -1
    var rpcCallbacks = {}
    var connected = false
    var onConnect = null
    var onDisconnect = null
    var clientChannels = 2

    fun init(channels = 2) {
        clientChannels = 2
        hostId = -1
        peerId = -1
        connected = false
    }

    fun start() {
        if hostId >= 0 {
            return
        }
        enet.init()
        hostId = enet.host_create(null, 1, clientChannels, 0, 0)
        if hostId < 0 {
            print("Failed to create client host")
            return
        }
    }

    fun connect(host, port, channels = 2) {
        clientChannels = channels
        if hostId < 0 {
            return
        }
        var addr = enet.address(host, port)
        peerId = enet.host_connect(hostId, addr, clientChannels, 0)
        if peerId < 0 {
            print("Failed to initiate connection")
        }
    }

    fun disconnect(data = 0) {
        if peerId >= 0 {
            enet.peer_disconnect(peerId, data)
        }
    }

    fun stop() {
        if hostId >= 0 {
            enet.host_destroy(hostId)
        }
        hostId = -1
        peerId = -1
        connected = false
        enet.deinit()
    }

    fun update(timeoutMs = 0) {
        if hostId < 0 {
            return
        }
        var ev = enet.host_service(hostId, timeoutMs)
        while ev != null and ev.type != enet.EVENT_NONE {
            if ev.type == enet.EVENT_CONNECT {
                connected = true
                if onConnect != null { onConnect(); }
            } else if ev.type == enet.EVENT_DISCONNECT {
                connected = false
                if onDisconnect != null { onDisconnect(); }
            } else if ev.type == enet.EVENT_RECEIVE {
                handlePacket(ev.packet.data)
            }
            ev = enet.host_service(hostId, 0)
        }
    }

    fun handlePacket(data) {
        var msg = _netDecode(data)
        if msg == null {
            return
        }
        if msg.kind == "rpc" {
            var method = msg.name
            var cb = rpcCallbacks[method]
            if cb != null {
                cb(msg.payload)
            }
        }
    }

    fun on(method, callback) {
        rpcCallbacks[method] = callback
    }

    fun send(method, payload) {
        if peerId < 0 {
            return
        }
        var data = _netEncode("rpc", method, payload)
        var pkt = enet.packet_create(data, enet.PACKET_RELIABLE)
        enet.peer_send(peerId, 0, pkt)
        enet.host_flush(hostId)
    }
}
`
}
