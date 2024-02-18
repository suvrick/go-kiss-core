package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BALANCE types.PacketServerType = 7

// BALANCE(7) "bottles:I, reason:B"
type Balance struct {
	Coins uint64
}

func (balance *Balance) Unmarshal(r *bytes.Reader) error {
	coins, err := leb128.ReadUInt64(r)
	if err == nil {
		balance.Coins = coins
	}
	return err
}
