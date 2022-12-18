package server

const BOTTLE_LEADER PacketServerType = 28

// BOTTLE_LEADER(28) "I"
type BottleLeader struct {
	LeaderID int64
}
