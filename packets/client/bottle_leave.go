package client

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_LEAVE types.PacketClientType = 27

// BOTTLE_LEAVE(27) ""
type BottleLeave struct {
}

func (bottleLeave BottleLeave) String() string {
	return "BOTTLE_LEAVE(27)"
}

func (bottleLeave *BottleLeave) Marshal() ([]byte, error) {
	return nil, nil
}
