package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_LEAVE PacketServerType = 27

// BOTTLE_LEAVE(27) "I"
type BottleLeave struct {
	PlayerID types.I
}
