package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BONUS types.PacketServerType = 17

// BONUS(17) "BB"
type Bonus struct {
	CanCollect byte
	Day        byte
}

func (bonus *Bonus) Unmarshal(r *bytes.Reader) error {
	var err error

	bonus.CanCollect, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	bonus.Day, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	return err
}
