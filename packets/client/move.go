package client

import "github.com/suvrick/go-kiss-core/types"

const MOVE PacketClientType = 21

// MOVE(21) "IB"
type Move struct {
	PlayerID  types.I
	ByteField types.B
}
