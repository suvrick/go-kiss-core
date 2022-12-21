package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_PLAY_DENIED PacketServerType = 24

// BOTTLE_PLAY_DENIED(24) "B"
type BottlePlayDenied struct {
	ByteField types.B
}
