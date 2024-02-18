package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const COLLECTIONS_POINTS types.PacketServerType = 130

// COLLECTIONS_POINTS(130) "I"
type CollectionsPoints struct {
	Points uint64
}

func (collectionsPoints *CollectionsPoints) Unmarshal(r *bytes.Reader) error {
	var err error
	collectionsPoints.Points, err = leb128.ReadUInt64(r)
	return err
}
