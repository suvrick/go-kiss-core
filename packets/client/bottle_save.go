package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_SAVE types.PacketClientType = 30

// BOTTLE_SAVE(30) "I"
type BottleSave struct {
	PlayerID uint64
}

func (bottleSave *BottleSave) Marshal() ([]byte, error) {
	return leb128.WriteUInt64(nil, bottleSave.PlayerID)
}
