package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEAVE types.PacketServerType = 27

// BOTTLE_LEAVE(27) "I"
type BottleLeave struct {
	PlayerID uint64
}

func (bottleLeave *BottleLeave) Unmarshal(r *bytes.Reader) error {
	var err error
	bottleLeave.PlayerID, err = leb128.ReadUInt64(r)
	return err
}
