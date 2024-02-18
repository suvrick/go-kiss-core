package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_JOIN types.PacketServerType = 26

// BOTTLE_JOIN(26) "IB"
type BottleJoin struct {
	PlayerID uint64
	Position byte
}

func (bottleJoin *BottleJoin) Unmarshal(r *bytes.Reader) error {
	var err error

	bottleJoin.PlayerID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	bottleJoin.Position, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	return err
}
