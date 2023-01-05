package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_ROLL types.PacketClientType = 28

// BOTTLE_ROLL(28) "I"
type BottleRoll struct {
	IntField types.B
}
