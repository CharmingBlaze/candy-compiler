package candy_enet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

type goHost struct {
	host   Host
	conn   *net.UDPConn
	events []Event
}

type goPeer struct {
	peer  Peer
	udp   *net.UDPAddr
	state string
}

type goBackend struct {
	mu       sync.Mutex
	inited   bool
	nextHost int
	nextPeer int
	nextPkt  int
	hosts    map[int]*goHost
	peers    map[int]*goPeer
	packets  map[int]*Packet
}

func newGoBackend() Backend {
	return &goBackend{
		nextHost: 1,
		nextPeer: 1,
		nextPkt:  1,
		hosts:    map[int]*goHost{},
		peers:    map[int]*goPeer{},
		packets:  map[int]*Packet{},
	}
}

func (g *goBackend) Name() string    { return "go-fallback" }
func (g *goBackend) Version() string { return "go-fallback-1.0" }

func (g *goBackend) Init() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.inited = true
	return nil
}

func (g *goBackend) Deinit() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, h := range g.hosts {
		if h.conn != nil {
			_ = h.conn.Close()
		}
	}
	g.hosts = map[int]*goHost{}
	g.peers = map[int]*goPeer{}
	g.packets = map[int]*Packet{}
	g.inited = false
	return nil
}

func (g *goBackend) HostCreate(addr *Address, peerLimit, channelLimit, inBandwidth, outBandwidth int) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.inited {
		return 0, fmt.Errorf("enet not initialized")
	}
	bind := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	hostAddr := Address{Host: "0.0.0.0", Port: 0}
	if addr != nil {
		ip := net.ParseIP(addr.Host)
		if ip == nil {
			ip = net.IPv4zero
		}
		bind = &net.UDPAddr{IP: ip, Port: addr.Port}
		hostAddr = *addr
	}
	conn, err := net.ListenUDP("udp", bind)
	if err != nil {
		return 0, err
	}
	if hostAddr.Port == 0 {
		hostAddr.Port = conn.LocalAddr().(*net.UDPAddr).Port
	}
	id := g.nextHost
	g.nextHost++
	g.hosts[id] = &goHost{
		host: Host{
			ID:           id,
			Address:      hostAddr,
			PeerLimit:    peerLimit,
			ChannelLimit: channelLimit,
			InBandwidth:  inBandwidth,
			OutBandwidth: outBandwidth,
		},
		conn: conn,
	}
	return id, nil
}

func (g *goBackend) HostDestroy(hostID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ok := g.hosts[hostID]
	if !ok {
		return fmt.Errorf("unknown host id: %d", hostID)
	}
	if h.conn != nil {
		_ = h.conn.Close()
	}
	delete(g.hosts, hostID)
	for id, p := range g.peers {
		if p.peer.HostID == hostID {
			delete(g.peers, id)
		}
	}
	return nil
}

func encodePacket(channel int, flags int, data []byte) []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(data)+10))
	_ = binary.Write(b, binary.BigEndian, uint16(channel))
	_ = binary.Write(b, binary.BigEndian, uint16(flags))
	_ = binary.Write(b, binary.BigEndian, uint32(len(data)))
	b.Write(data)
	return b.Bytes()
}

func decodePacket(payload []byte) (int, int, []byte, error) {
	if len(payload) < 8 {
		return 0, 0, nil, fmt.Errorf("short packet")
	}
	channel := int(binary.BigEndian.Uint16(payload[0:2]))
	flags := int(binary.BigEndian.Uint16(payload[2:4]))
	n := int(binary.BigEndian.Uint32(payload[4:8]))
	if len(payload) < 8+n {
		return 0, 0, nil, fmt.Errorf("invalid packet length")
	}
	return channel, flags, payload[8 : 8+n], nil
}

func (g *goBackend) HostService(hostID int, timeoutMS int) (*Event, error) {
	g.mu.Lock()
	h, ok := g.hosts[hostID]
	if !ok {
		g.mu.Unlock()
		return nil, fmt.Errorf("unknown host id: %d", hostID)
	}
	if len(h.events) > 0 {
		ev := h.events[0]
		h.events = h.events[1:]
		g.mu.Unlock()
		return &ev, nil
	}
	conn := h.conn
	g.mu.Unlock()

	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
		return nil, err
	}
	buf := make([]byte, 65535)
	n, from, err := conn.ReadFromUDP(buf)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return &Event{Type: EventNone, HostID: hostID}, nil
		}
		return nil, err
	}
	ch, flags, payload, derr := decodePacket(buf[:n])
	if derr != nil {
		return &Event{Type: EventNone, HostID: hostID}, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	peerID := 0
	for id, p := range g.peers {
		if p.peer.HostID == hostID && p.udp.IP.Equal(from.IP) && p.udp.Port == from.Port {
			peerID = id
			break
		}
	}
	if peerID == 0 {
		peerID = g.nextPeer
		g.nextPeer++
		p := &goPeer{
			peer: Peer{
				ID:      peerID,
				HostID:  hostID,
				Address: Address{Host: from.IP.String(), Port: from.Port},
			},
			udp:   from,
			state: "connected",
		}
		g.peers[peerID] = p
		gh := g.hosts[hostID]
		gh.events = append(gh.events, Event{
			Type:    EventConnect,
			HostID:  hostID,
			PeerID:  peerID,
			Address: p.peer.Address,
		})
	}
	packet := &Packet{
		ID:    g.nextPkt,
		Data:  payload,
		Flags: flags,
	}
	g.nextPkt++
	g.packets[packet.ID] = packet
	return &Event{
		Type:    EventReceive,
		HostID:  hostID,
		PeerID:  peerID,
		Channel: ch,
		Packet:  packet,
		Address: Address{Host: from.IP.String(), Port: from.Port},
	}, nil
}

func (g *goBackend) HostFlush(hostID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.hosts[hostID]; !ok {
		return fmt.Errorf("unknown host id: %d", hostID)
	}
	return nil
}

func (g *goBackend) HostBandwidthLimit(hostID int, inBandwidth, outBandwidth int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ok := g.hosts[hostID]
	if !ok {
		return fmt.Errorf("unknown host id: %d", hostID)
	}
	h.host.InBandwidth = inBandwidth
	h.host.OutBandwidth = outBandwidth
	return nil
}

func (g *goBackend) HostChannelLimit(hostID int, channels int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ok := g.hosts[hostID]
	if !ok {
		return fmt.Errorf("unknown host id: %d", hostID)
	}
	h.host.ChannelLimit = channels
	return nil
}

func (g *goBackend) HostCompressWithRangeCoder(hostID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.hosts[hostID]; !ok {
		return fmt.Errorf("unknown host id: %d", hostID)
	}
	return nil
}

func (g *goBackend) HostConnect(hostID int, addr Address, channels int, data int) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.hosts[hostID]; !ok {
		return 0, fmt.Errorf("unknown host id: %d", hostID)
	}
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr.Host, addr.Port))
	if err != nil {
		return 0, err
	}
	id := g.nextPeer
	g.nextPeer++
	g.peers[id] = &goPeer{
		peer: Peer{
			ID:      id,
			HostID:  hostID,
			Address: addr,
			Data:    data,
		},
		udp:   udpAddr,
		state: "connected",
	}
	g.hosts[hostID].events = append(g.hosts[hostID].events, Event{
		Type:    EventConnect,
		HostID:  hostID,
		PeerID:  id,
		Data:    data,
		Address: addr,
	})
	_ = channels
	return id, nil
}

func (g *goBackend) PeerDisconnect(peerID int, data int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.peers[peerID]
	if !ok {
		return fmt.Errorf("unknown peer id: %d", peerID)
	}
	if h, ok := g.hosts[p.peer.HostID]; ok {
		h.events = append(h.events, Event{
			Type:    EventDisconnect,
			HostID:  p.peer.HostID,
			PeerID:  peerID,
			Data:    data,
			Address: p.peer.Address,
		})
	}
	delete(g.peers, peerID)
	return nil
}

func (g *goBackend) PeerDisconnectNow(peerID int, data int) error {
	return g.PeerDisconnect(peerID, data)
}

func (g *goBackend) PeerPing(peerID int) error {
	g.mu.Lock()
	p, ok := g.peers[peerID]
	if !ok {
		g.mu.Unlock()
		return fmt.Errorf("unknown peer id: %d", peerID)
	}
	h := g.hosts[p.peer.HostID]
	conn := h.conn
	addr := p.udp
	g.mu.Unlock()
	_, err := conn.WriteToUDP(encodePacket(0, 0, []byte("PING")), addr)
	return err
}

func (g *goBackend) PeerTimeout(peerID int, timeoutLimit, timeoutMinimum, timeoutMaximum int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.peers[peerID]; !ok {
		return fmt.Errorf("unknown peer id: %d", peerID)
	}
	_, _, _ = timeoutLimit, timeoutMinimum, timeoutMaximum
	return nil
}

func (g *goBackend) PeerReset(peerID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.peers[peerID]; !ok {
		return fmt.Errorf("unknown peer id: %d", peerID)
	}
	delete(g.peers, peerID)
	return nil
}

func (g *goBackend) PacketCreate(data []byte, flags int) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id := g.nextPkt
	g.nextPkt++
	payload := make([]byte, len(data))
	copy(payload, data)
	g.packets[id] = &Packet{ID: id, Data: payload, Flags: flags}
	return id, nil
}

func (g *goBackend) PacketDestroy(packetID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.packets[packetID]; !ok {
		return fmt.Errorf("unknown packet id: %d", packetID)
	}
	delete(g.packets, packetID)
	return nil
}

func (g *goBackend) PeerSend(peerID int, channel int, packetID int) error {
	g.mu.Lock()
	p, ok := g.peers[peerID]
	if !ok {
		g.mu.Unlock()
		return fmt.Errorf("unknown peer id: %d", peerID)
	}
	h, ok := g.hosts[p.peer.HostID]
	if !ok {
		g.mu.Unlock()
		return fmt.Errorf("missing host for peer id: %d", peerID)
	}
	pkt, ok := g.packets[packetID]
	if !ok {
		g.mu.Unlock()
		return fmt.Errorf("unknown packet id: %d", packetID)
	}
	addr := p.udp
	conn := h.conn
	body := encodePacket(channel, pkt.Flags, pkt.Data)
	g.mu.Unlock()
	_, err := conn.WriteToUDP(body, addr)
	return err
}

