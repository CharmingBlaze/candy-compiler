package candy_evaluator

import (
	"candy/candy_enet"
	"fmt"
)

var enetRuntime = candy_enet.NewRuntime()

func mapString(v *Value, key string) string {
	if v == nil || v.Kind != ValMap || v.StrMap == nil {
		return ""
	}
	x, ok := v.StrMap[key]
	if !ok {
		return ""
	}
	if x.Kind == ValString {
		return x.Str
	}
	return x.String()
}

func mapInt(v *Value, key string, d int) int {
	if v == nil || v.Kind != ValMap || v.StrMap == nil {
		return d
	}
	x, ok := v.StrMap[key]
	if !ok {
		return d
	}
	i, err := i64Arg(&x)
	if err != nil {
		return d
	}
	return int(i)
}

func addressFromValue(v *Value) (*candy_enet.Address, error) {
	if v == nil || v.Kind == ValNull {
		return nil, nil
	}
	if v.Kind != ValMap {
		return nil, fmt.Errorf("address must be map {host,port} or null")
	}
	host := mapString(v, "host")
	port := mapInt(v, "port", 0)
	return &candy_enet.Address{Host: host, Port: port}, nil
}

func addressValue(addr candy_enet.Address) *Value {
	return &Value{
		Kind: ValMap,
		StrMap: map[string]Value{
			"host": {Kind: ValString, Str: addr.Host},
			"port": {Kind: ValInt, I64: int64(addr.Port)},
		},
	}
}

func packetValue(pkt *candy_enet.Packet) *Value {
	if pkt == nil {
		return &Value{Kind: ValNull}
	}
	return &Value{
		Kind: ValMap,
		StrMap: map[string]Value{
			"id":    {Kind: ValInt, I64: int64(pkt.ID)},
			"flags": {Kind: ValInt, I64: int64(pkt.Flags)},
			"data":  {Kind: ValString, Str: string(pkt.Data)},
		},
	}
}

func eventValue(ev *candy_enet.Event) *Value {
	if ev == nil {
		return &Value{Kind: ValNull}
	}
	return &Value{
		Kind: ValMap,
		StrMap: map[string]Value{
			"type":    {Kind: ValString, Str: string(ev.Type)},
			"hostId":  {Kind: ValInt, I64: int64(ev.HostID)},
			"peerId":  {Kind: ValInt, I64: int64(ev.PeerID)},
			"channel": {Kind: ValInt, I64: int64(ev.Channel)},
			"data":    {Kind: ValInt, I64: int64(ev.Data)},
			"packet":  valueToValue(packetValue(ev.Packet)),
			"address": valueToValue(addressValue(ev.Address)),
		},
	}
}

func registerENetModule(e *Env) {
	if e == nil {
		return
	}
	fn := map[string]func(args []*Value) (*Value, error){
		"init": func(args []*Value) (*Value, error) {
			if err := arg0(args); err != nil {
				return nil, err
			}
			if err := enetRuntime.Init(); err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, nil
		},
		"deinit": func(args []*Value) (*Value, error) {
			if err := arg0(args); err != nil {
				return nil, err
			}
			if err := enetRuntime.Deinit(); err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, nil
		},
		"version": func(args []*Value) (*Value, error) {
			if err := arg0(args); err != nil {
				return nil, err
			}
			return &Value{Kind: ValString, Str: enetRuntime.Version()}, nil
		},
		"backend": func(args []*Value) (*Value, error) {
			if err := arg0(args); err != nil {
				return nil, err
			}
			return &Value{Kind: ValString, Str: enetRuntime.BackendName()}, nil
		},
		"address": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("enet.address(host, port)")
			}
			if args[0] == nil || args[0].Kind != ValString {
				return nil, fmt.Errorf("address host must be string")
			}
			port, err := i64Arg(args[1])
			if err != nil {
				return nil, err
			}
			return addressValue(candy_enet.Address{Host: args[0].Str, Port: int(port)}), nil
		},
		"host_create": func(args []*Value) (*Value, error) {
			if len(args) != 5 {
				return nil, fmt.Errorf("enet.host_create(address|null, peers, channels, inBandwidth, outBandwidth)")
			}
			addr, err := addressFromValue(args[0])
			if err != nil {
				return nil, err
			}
			peers, e1 := i64Arg(args[1])
			channels, e2 := i64Arg(args[2])
			inBW, e3 := i64Arg(args[3])
			outBW, e4 := i64Arg(args[4])
			if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
				return nil, fmt.Errorf("host_create: numeric arguments expected")
			}
			id, err := enetRuntime.HostCreate(addr, int(peers), int(channels), int(inBW), int(outBW))
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValInt, I64: int64(id)}, nil
		},
		"host_destroy": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("host_destroy(hostId)")
			}
			hostID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.HostDestroy(int(hostID))
		},
		"host_service": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("host_service(hostId, timeoutMs)")
			}
			hostID, e1 := i64Arg(args[0])
			timeoutMS, e2 := i64Arg(args[1])
			if e1 != nil || e2 != nil {
				return nil, fmt.Errorf("host_service: numeric arguments expected")
			}
			ev, err := enetRuntime.HostService(int(hostID), int(timeoutMS))
			if err != nil {
				return nil, err
			}
			return eventValue(ev), nil
		},
		"host_flush": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("host_flush(hostId)")
			}
			hostID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.HostFlush(int(hostID))
		},
		"host_bandwidth_limit": func(args []*Value) (*Value, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("host_bandwidth_limit(hostId, inBandwidth, outBandwidth)")
			}
			hostID, e1 := i64Arg(args[0])
			inBW, e2 := i64Arg(args[1])
			outBW, e3 := i64Arg(args[2])
			if e1 != nil || e2 != nil || e3 != nil {
				return nil, fmt.Errorf("host_bandwidth_limit: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.HostBandwidthLimit(int(hostID), int(inBW), int(outBW))
		},
		"host_channel_limit": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("host_channel_limit(hostId, channels)")
			}
			hostID, e1 := i64Arg(args[0])
			ch, e2 := i64Arg(args[1])
			if e1 != nil || e2 != nil {
				return nil, fmt.Errorf("host_channel_limit: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.HostChannelLimit(int(hostID), int(ch))
		},
		"host_compress_range_coder": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("host_compress_range_coder(hostId)")
			}
			hostID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.HostCompressWithRangeCoder(int(hostID))
		},
		"set_range_coder": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("set_range_coder(hostId)")
			}
			hostID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.HostCompressWithRangeCoder(int(hostID))
		},
		"host_connect": func(args []*Value) (*Value, error) {
			if len(args) != 4 {
				return nil, fmt.Errorf("host_connect(hostId, address, channels, data)")
			}
			hostID, e1 := i64Arg(args[0])
			addr, e2 := addressFromValue(args[1])
			channels, e3 := i64Arg(args[2])
			data, e4 := i64Arg(args[3])
			if e1 != nil || e2 != nil || e3 != nil || e4 != nil || addr == nil {
				return nil, fmt.Errorf("host_connect: invalid arguments")
			}
			peerID, err := enetRuntime.HostConnect(int(hostID), *addr, int(channels), int(data))
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValInt, I64: int64(peerID)}, nil
		},
		"peer_disconnect": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("peer_disconnect(peerId, data)")
			}
			peerID, e1 := i64Arg(args[0])
			data, e2 := i64Arg(args[1])
			if e1 != nil || e2 != nil {
				return nil, fmt.Errorf("peer_disconnect: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerDisconnect(int(peerID), int(data))
		},
		"peer_disconnect_now": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("peer_disconnect_now(peerId, data)")
			}
			peerID, e1 := i64Arg(args[0])
			data, e2 := i64Arg(args[1])
			if e1 != nil || e2 != nil {
				return nil, fmt.Errorf("peer_disconnect_now: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerDisconnectNow(int(peerID), int(data))
		},
		"peer_ping": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("peer_ping(peerId)")
			}
			peerID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerPing(int(peerID))
		},
		"peer_timeout": func(args []*Value) (*Value, error) {
			if len(args) != 4 {
				return nil, fmt.Errorf("peer_timeout(peerId, timeoutLimit, timeoutMinimum, timeoutMaximum)")
			}
			peerID, e1 := i64Arg(args[0])
			a, e2 := i64Arg(args[1])
			b, e3 := i64Arg(args[2])
			c, e4 := i64Arg(args[3])
			if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
				return nil, fmt.Errorf("peer_timeout: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerTimeout(int(peerID), int(a), int(b), int(c))
		},
		"peer_reset": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("peer_reset(peerId)")
			}
			peerID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerReset(int(peerID))
		},
		"packet_create": func(args []*Value) (*Value, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("packet_create(data, flags)")
			}
			if args[0] == nil || args[0].Kind != ValString {
				return nil, fmt.Errorf("packet_create data must be string")
			}
			flags, err := i64Arg(args[1])
			if err != nil {
				return nil, err
			}
			packetID, err := enetRuntime.PacketCreate([]byte(args[0].Str), int(flags))
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValInt, I64: int64(packetID)}, nil
		},
		"packet_destroy": func(args []*Value) (*Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("packet_destroy(packetId)")
			}
			packetID, err := i64Arg(args[0])
			if err != nil {
				return nil, err
			}
			return &Value{Kind: ValNull}, enetRuntime.PacketDestroy(int(packetID))
		},
		"peer_send": func(args []*Value) (*Value, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("peer_send(peerId, channel, packetId)")
			}
			peerID, e1 := i64Arg(args[0])
			channel, e2 := i64Arg(args[1])
			packetID, e3 := i64Arg(args[2])
			if e1 != nil || e2 != nil || e3 != nil {
				return nil, fmt.Errorf("peer_send: numeric arguments expected")
			}
			return &Value{Kind: ValNull}, enetRuntime.PeerSend(int(peerID), int(channel), int(packetID))
		},
	}
	e.Set("enet", newModule("enet", fn, map[string]*Value{
		"EVENT_NONE":         {Kind: ValString, Str: string(candy_enet.EventNone)},
		"EVENT_CONNECT":      {Kind: ValString, Str: string(candy_enet.EventConnect)},
		"EVENT_DISCONNECT":   {Kind: ValString, Str: string(candy_enet.EventDisconnect)},
		"EVENT_RECEIVE":      {Kind: ValString, Str: string(candy_enet.EventReceive)},
		"PACKET_RELIABLE":    {Kind: ValInt, I64: candy_enet.PacketFlagReliable},
		"PACKET_UNSEQUENCED": {Kind: ValInt, I64: candy_enet.PacketFlagUnsequenced},
		"PACKET_NO_ALLOCATE": {Kind: ValInt, I64: candy_enet.PacketFlagNoAllocate},
		"PACKET_UNRELIABLE":  {Kind: ValInt, I64: candy_enet.PacketFlagUnreliable},
	}))
}

