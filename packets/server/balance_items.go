package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BALANCE_ITEMS types.PacketServerType = 310

// BALANCE_ITEMS(310) "[BII]"
type BalanceItems struct {
	Items []BalanceItem
}

type BalanceItem struct {
	Type   byte
	Count1 uint64
	Count2 uint64
}

func (p BalanceItems) String() string {
	return "BALANCE_ITEMS(310)"
}

func (balanceItems *BalanceItems) Unmarshal(r *bytes.Reader) error {

	len, err := leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	balanceItems.Items = make([]BalanceItem, len)

	for len > 0 {
		var err error
		item := BalanceItem{}
		item.Type, err = leb128.ReadByte(r)
		if err != nil {
			return err
		}

		item.Count1, err = leb128.ReadUInt64(r)
		if err != nil {
			return err
		}

		item.Count2, err = leb128.ReadUInt64(r)
		if err != nil {
			return err
		}

		balanceItems.Items = append(balanceItems.Items, item)

		len--
	}

	return nil
}
