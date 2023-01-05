package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_PLAY types.PacketClientType = 26

// BOTTLE_PLAY(26) "B,B"
type BottlePlay struct {
	RoomID types.B
	LangID types.B `pack:"optional"`
}
