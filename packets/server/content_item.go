package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const CONTEST_ITEMS types.PacketServerType = 10

// CONTEST_ITEMS(10) "B,B"
type ContestItems struct {
	ByteField  byte
	ByteField2 *byte
}

func (contestItems *ContestItems) Unmarshal(r *bytes.Reader) error {
	var err error

	contestItems.ByteField, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	t, err2 := leb128.ReadByte(r)
	if err2 != nil {
		return nil
	}

	contestItems.ByteField = t

	return err
}
