package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_KISS types.PacketClientType = 29

// BOTTLE_KISS(29) "B"
type BottleKiss struct {
	Answer types.B
}
