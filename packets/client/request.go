package client

import "github.com/suvrick/go-kiss-core/types"

const REQUEST PacketClientType = 8

// REQUEST (8) ""
type Request struct {
	Players []types.I
	Mask    types.I
}
