package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const TOGGLE_TAPE types.PacketClientType = 10

// TOGGLE_TAPE(10) ""
type ToggleType struct {
	Off byte
}

func (toggleType ToggleType) String() string {
	return "TOGGLE_TAPE(10)"
}

func (toggleType *ToggleType) Marshal() ([]byte, error) {
	return leb128.WriteByte(nil, toggleType.Off)
}
