package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_SAVE types.PacketClientType = 30

// BOTTLE_SAVE(30) "I"
type BottleSave struct {
	IntField types.I
}
