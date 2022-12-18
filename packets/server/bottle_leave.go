package server

const BOTTLE_LEAVE PacketServerType = 27

// BOTTLE_LEAVE(27) "I"
type BottleLeave struct {
	PlayerID int64
}
