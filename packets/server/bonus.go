package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BONUS types.PacketServerType = 17

// BONUS(17) "B"
type Bonus struct {
	// CanCollect byte
	Day byte
}

func (p Bonus) String() string {
	return "BONUS(17)"
}

func (bonus *Bonus) Unmarshal(r *bytes.Reader) error {
	var err error
	bonus.Day, err = leb128.ReadByte(r)
	return err
}
