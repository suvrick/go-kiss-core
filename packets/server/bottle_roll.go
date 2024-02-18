package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROLL types.PacketServerType = 29

// BOTTLE_ROLL(29) "II,II"
type BottleRoll struct {
	LeaderID  uint64
	RollerID  uint64
	IntField  *uint64
	IntField2 *uint64
}

func (bottleRoll *BottleRoll) Unmarshal(r *bytes.Reader) error {
	var err error

	bottleRoll.LeaderID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	bottleRoll.RollerID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	t, err2 := leb128.ReadUInt64(r)
	if err2 != nil {
		return nil
	}

	bottleRoll.IntField = &t

	t, err2 = leb128.ReadUInt64(r)
	if err2 != nil {
		return nil
	}

	bottleRoll.IntField2 = &t

	return nil
}
