package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const PHOTOS_LIKE types.PacketClientType = 207

// PHOTOS_LIKE(207)
type PhotoLike struct {
	PlayerID uint64
	PhotoID  byte
	IsLike   byte
}

func (photoLike PhotoLike) String() string {
	return "PHOTOS_LIKE(207)"
}

func (photoLike *PhotoLike) Marshal() ([]byte, error) {

	data, err := leb128.WriteUInt64(nil, photoLike.PlayerID)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, photoLike.PhotoID)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, photoLike.IsLike)
	if err != nil {
		return nil, err
	}

	return data, nil
}
