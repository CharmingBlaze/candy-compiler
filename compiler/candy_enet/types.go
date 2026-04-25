package candy_enet

type EventType string

const (
	EventNone       EventType = "none"
	EventConnect    EventType = "connect"
	EventDisconnect EventType = "disconnect"
	EventReceive    EventType = "receive"
)

const (
	PacketFlagReliable    = 1
	PacketFlagUnsequenced = 2
	PacketFlagNoAllocate  = 4
	PacketFlagUnreliable  = 8
)

type Address struct {
	Host string
	Port int
}

type Packet struct {
	ID    int
	Data  []byte
	Flags int
}

type Event struct {
	Type    EventType
	HostID  int
	PeerID  int
	Channel int
	Data    int
	Packet  *Packet
	Address Address
}

type Host struct {
	ID           int
	Address      Address
	PeerLimit    int
	ChannelLimit int
	InBandwidth  int
	OutBandwidth int
}

type Peer struct {
	ID      int
	HostID  int
	Address Address
	Data    int
}

