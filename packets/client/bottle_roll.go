package client

import (
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROLL types.PacketClientType = 28

// BOTTLE_ROLL(28) "I"
type BottleRoll struct {
	//IntField byte
}

func (bottleRoll BottleRoll) String() string {
	return "BOTTLE_ROLL(28)"
}

func (bottleRoll *BottleRoll) Marshal() ([]byte, error) {
	return []byte{}, nil //leb128.WriteByte(nil, bottleRoll.IntField)
}
