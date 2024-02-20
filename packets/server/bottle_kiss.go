package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_KISS types.PacketServerType = 30

// BOTTLE_KISS(30) "IB"
type BottleKiss struct {
	PlayerID uint64
	Answer   byte
}

func (p BottleKiss) String() string {
	return "BOTTLE_KISS(30)"
}

func (bottleKiss *BottleKiss) Unmarshal(r *bytes.Reader) error {
	var err error

	bottleKiss.PlayerID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	bottleKiss.Answer, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	return err
}
