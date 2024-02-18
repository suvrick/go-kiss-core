package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ENTER types.PacketServerType = 35

// BOTTLE_ENTER(35) ""
type BottleEnter struct {
}

func (bottleEnter *BottleEnter) Unmarshal(r *bytes.Reader) error {
	return nil
}
