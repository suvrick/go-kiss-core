package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROOM types.PacketServerType = 25

// BOTTLE_ROOM(25) "III[I][I]"
type BottleRoom struct {
	RoomID    uint64
	IntField2 uint64
	IntField3 uint64
	Players   []uint64
	Positions []uint64
}

func (p BottleRoom) String() string {
	return "BOTTLE_ROOM(25)"
}

func (bottleRoom *BottleRoom) Unmarshal(r *bytes.Reader) error {
	var err error

	bottleRoom.RoomID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	bottleRoom.IntField2, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	bottleRoom.IntField3, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	var len uint64
	len, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}
	bottleRoom.Players = make([]uint64, len)
	var t uint64
	for len > 0 {
		t, err = leb128.ReadUInt64(r)
		if err != nil {
			return err
		}
		bottleRoom.Players = append(bottleRoom.Players, t)
		len--
	}

	len, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}
	bottleRoom.Positions = make([]uint64, len)
	for len > 0 {
		t, err = leb128.ReadUInt64(r)
		if err != nil {
			return err
		}
		bottleRoom.Positions = append(bottleRoom.Positions, t)
		len--
	}

	return err
}
