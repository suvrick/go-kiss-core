package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROLL types.PacketClientType = 28

// BOTTLE_ROLL(28) "I"
type BottleRoll struct {
	IntField byte
}

func (bottleRoll *BottleRoll) Marshal() ([]byte, error) {
	return leb128.WriteByte(nil, bottleRoll.IntField)
}
