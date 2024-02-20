package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const LOGIN types.PacketServerType = 4

// LOGIN(4)
type Login struct {
	Result  byte
	HiroID  uint64
	Balance uint64
}

func (p Login) String() string {
	return "LOGIN(4)"
}

func (login *Login) Unmarshal(r *bytes.Reader) error {
	var err error

	login.Result, err = leb128.ReadByte(r)

	if login.Result == 0 {

		login.HiroID, err = leb128.ReadUInt64(r)
		if err != nil {
			return nil
		}

		login.Balance, err = leb128.ReadUInt64(r)
		if err != nil {
			return nil
		}
	}

	return err
}
