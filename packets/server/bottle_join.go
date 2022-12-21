package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_JOIN PacketServerType = 26

// BOTTLE_JOIN(26) "IB"
type BottleJoin struct {
	PlayerID types.I
	Position types.B
}
