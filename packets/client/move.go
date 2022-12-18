package client

const MOVE PacketClientType = 21

// MOVE(21) "IB"
type Move struct {
	PlayerID  int64
	ByteField byte
}
