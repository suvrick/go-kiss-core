package client

const INFO PacketClientType = 5

// INFO(5) ""
type Info struct {
	PlayerID uint64
	Mask     byte
}
