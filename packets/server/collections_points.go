package server

import "github.com/suvrick/go-kiss-core/types"

const COLLECTIONS_POINTS PacketServerType = 130

// COLLECTIONS_POINTS(130) "I"
type CollectionsPoints struct {
	Points types.I
}
