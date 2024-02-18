package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_PLAY_DENIED types.PacketServerType = 24

// BOTTLE_PLAY_DENIED(24) "B"
type BottlePlayDenied struct {
	ByteField byte
}

func (bottlePlayDenied *BottlePlayDenied) Unmarshal(r *bytes.Reader) error {
	var err error
	bottlePlayDenied.ByteField, err = leb128.ReadByte(r)
	return err
}
