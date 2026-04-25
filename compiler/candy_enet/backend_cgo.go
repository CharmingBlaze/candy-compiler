//go:build enet && cgo

package candy_enet

/*
#cgo CFLAGS: -DENET_CANDY=1
*/
import "C"

// cgoBackend currently shares the same behavior contract as the Go backend
// while being selected under `-tags enet` for compatibility routing.
type cgoBackend struct {
	inner Backend
}

func (c *cgoBackend) Name() string { return "cgo-enet" }
func (c *cgoBackend) Init() error  { return c.inner.Init() }
func (c *cgoBackend) Deinit() error {
	return c.inner.Deinit()
}
func (c *cgoBackend) Version() string { return "cgo-enet-1.0" }
func (c *cgoBackend) HostCreate(addr *Address, peerLimit, channelLimit, inBandwidth, outBandwidth int) (int, error) {
	return c.inner.HostCreate(addr, peerLimit, channelLimit, inBandwidth, outBandwidth)
}
func (c *cgoBackend) HostDestroy(hostID int) error { return c.inner.HostDestroy(hostID) }
func (c *cgoBackend) HostService(hostID int, timeoutMS int) (*Event, error) {
	return c.inner.HostService(hostID, timeoutMS)
}
func (c *cgoBackend) HostFlush(hostID int) error { return c.inner.HostFlush(hostID) }
func (c *cgoBackend) HostBandwidthLimit(hostID int, inBandwidth, outBandwidth int) error {
	return c.inner.HostBandwidthLimit(hostID, inBandwidth, outBandwidth)
}
func (c *cgoBackend) HostChannelLimit(hostID int, channels int) error {
	return c.inner.HostChannelLimit(hostID, channels)
}
func (c *cgoBackend) HostCompressWithRangeCoder(hostID int) error {
	return c.inner.HostCompressWithRangeCoder(hostID)
}
func (c *cgoBackend) HostConnect(hostID int, addr Address, channels int, data int) (int, error) {
	return c.inner.HostConnect(hostID, addr, channels, data)
}
func (c *cgoBackend) PeerDisconnect(peerID int, data int) error {
	return c.inner.PeerDisconnect(peerID, data)
}
func (c *cgoBackend) PeerDisconnectNow(peerID int, data int) error {
	return c.inner.PeerDisconnectNow(peerID, data)
}
func (c *cgoBackend) PeerPing(peerID int) error { return c.inner.PeerPing(peerID) }
func (c *cgoBackend) PeerTimeout(peerID int, timeoutLimit, timeoutMinimum, timeoutMaximum int) error {
	return c.inner.PeerTimeout(peerID, timeoutLimit, timeoutMinimum, timeoutMaximum)
}
func (c *cgoBackend) PeerReset(peerID int) error { return c.inner.PeerReset(peerID) }
func (c *cgoBackend) PacketCreate(data []byte, flags int) (int, error) {
	return c.inner.PacketCreate(data, flags)
}
func (c *cgoBackend) PacketDestroy(packetID int) error { return c.inner.PacketDestroy(packetID) }
func (c *cgoBackend) PeerSend(peerID int, channel int, packetID int) error {
	return c.inner.PeerSend(peerID, channel, packetID)
}

func tryNewCgoBackend() Backend {
	return &cgoBackend{inner: newGoBackend()}
}
