package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const MOVE types.PacketClientType = 21

// MOVE(21) "IB"
type Move struct {
	PlayerID  uint64
	ByteField *byte
}

func (move *Move) Marshal() ([]byte, error) {

	data, err := leb128.WriteUInt64(nil, move.PlayerID)
	if err != nil {
		return nil, err
	}

	if move.ByteField != nil {
		data2, err := leb128.WriteByte(nil, *move.ByteField)
		if err == nil {
			data = append(data, data2...)
		}
	}

	return data, nil
}
