package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_PLAY types.PacketClientType = 26

// BOTTLE_PLAY(26) "B,B"
type BottlePlay struct {
	RoomID byte
	// LangID *byte
}

func (bottlePlay *BottlePlay) Marshal() ([]byte, error) {
	return leb128.WriteByte(nil, bottlePlay.RoomID)
}
