package candy_enet

import "fmt"

type Backend interface {
	Name() string
	Init() error
	Deinit() error
	Version() string

	HostCreate(addr *Address, peerLimit, channelLimit, inBandwidth, outBandwidth int) (int, error)
	HostDestroy(hostID int) error
	HostService(hostID int, timeoutMS int) (*Event, error)
	HostFlush(hostID int) error
	HostBandwidthLimit(hostID int, inBandwidth, outBandwidth int) error
	HostChannelLimit(hostID int, channels int) error
	HostCompressWithRangeCoder(hostID int) error

	HostConnect(hostID int, addr Address, channels int, data int) (int, error)
	PeerDisconnect(peerID int, data int) error
	PeerDisconnectNow(peerID int, data int) error
	PeerPing(peerID int) error
	PeerTimeout(peerID int, timeoutLimit, timeoutMinimum, timeoutMaximum int) error
	PeerReset(peerID int) error

	PacketCreate(data []byte, flags int) (int, error)
	PacketDestroy(packetID int) error
	PeerSend(peerID int, channel int, packetID int) error
}

type Runtime struct {
	backend Backend
}

func NewRuntime() *Runtime {
	var b Backend
	if cb := tryNewCgoBackend(); cb != nil {
		b = cb
	} else {
		b = newGoBackend()
	}
	return &Runtime{backend: b}
}

func (r *Runtime) BackendName() string {
	return r.backend.Name()
}

func (r *Runtime) Init() error                                 { return r.backend.Init() }
func (r *Runtime) Deinit() error                               { return r.backend.Deinit() }
func (r *Runtime) Version() string                             { return r.backend.Version() }
func (r *Runtime) HostCreate(addr *Address, a, b, c, d int) (int, error) {
	return r.backend.HostCreate(addr, a, b, c, d)
}
func (r *Runtime) HostDestroy(hostID int) error                                   { return r.backend.HostDestroy(hostID) }
func (r *Runtime) HostService(hostID int, timeoutMS int) (*Event, error)          { return r.backend.HostService(hostID, timeoutMS) }
func (r *Runtime) HostFlush(hostID int) error                                      { return r.backend.HostFlush(hostID) }
func (r *Runtime) HostBandwidthLimit(hostID int, inBandwidth, outBandwidth int) error {
	return r.backend.HostBandwidthLimit(hostID, inBandwidth, outBandwidth)
}
func (r *Runtime) HostChannelLimit(hostID int, channels int) error {
	return r.backend.HostChannelLimit(hostID, channels)
}
func (r *Runtime) HostCompressWithRangeCoder(hostID int) error {
	return r.backend.HostCompressWithRangeCoder(hostID)
}
func (r *Runtime) HostConnect(hostID int, addr Address, channels int, data int) (int, error) {
	return r.backend.HostConnect(hostID, addr, channels, data)
}
func (r *Runtime) PeerDisconnect(peerID int, data int) error {
	return r.backend.PeerDisconnect(peerID, data)
}
func (r *Runtime) PeerDisconnectNow(peerID int, data int) error {
	return r.backend.PeerDisconnectNow(peerID, data)
}
func (r *Runtime) PeerPing(peerID int) error { return r.backend.PeerPing(peerID) }
func (r *Runtime) PeerTimeout(peerID int, a, b, c int) error {
	return r.backend.PeerTimeout(peerID, a, b, c)
}
func (r *Runtime) PeerReset(peerID int) error { return r.backend.PeerReset(peerID) }
func (r *Runtime) PacketCreate(data []byte, flags int) (int, error) {
	return r.backend.PacketCreate(data, flags)
}
func (r *Runtime) PacketDestroy(packetID int) error { return r.backend.PacketDestroy(packetID) }
func (r *Runtime) PeerSend(peerID int, channel int, packetID int) error {
	return r.backend.PeerSend(peerID, channel, packetID)
}

func (r *Runtime) Ensure() error {
	if r == nil || r.backend == nil {
		return fmt.Errorf("enet runtime unavailable")
	}
	return nil
}

