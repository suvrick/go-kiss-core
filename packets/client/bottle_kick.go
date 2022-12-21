package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_KICK PacketClientType = 31

// BOTTLE_KICK(31) "I"
type BottleKick struct {
	IntField types.I
}
