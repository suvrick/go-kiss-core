package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_LEADER PacketServerType = 28

// BOTTLE_LEADER(28) "I"
type BottleLeader struct {
	LeaderID types.I
}
