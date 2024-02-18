package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_KICK types.PacketClientType = 31

// BOTTLE_KICK(31) "I"
type BottleKick struct {
	IntField uint64
}

func (bottleKick *BottleKick) Marshal() ([]byte, error) {
	return leb128.WriteUInt64(nil, bottleKick.IntField)
}
