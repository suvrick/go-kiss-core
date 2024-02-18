package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BUY types.PacketClientType = 6

// BONUS(6) "IIIIBI,B"
type Buy struct {
	BuyType    uint64
	Coin       uint64
	PlayerID   uint64
	PrizeID    uint64
	ByteFiald  byte
	Count      uint64
	ByteFiald2 *byte
}

func (buy *Buy) Marshal() ([]byte, error) {

	data, err := leb128.WriteUInt64(nil, buy.BuyType)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteUInt64(data, buy.Coin)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteUInt64(data, buy.PlayerID)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteUInt64(data, buy.PrizeID)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, buy.ByteFiald)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteUInt64(data, buy.Count)
	if err != nil {
		return nil, err
	}

	if buy.ByteFiald2 != nil {
		data2, err := leb128.WriteByte(nil, (*buy.ByteFiald2))
		if err == nil {
			data = append(data, data2...)
		}
	}

	return data, nil
}

/*

[2, 30, 45066660, 10169, 0, 1, 5]  //free buy
[2, 30, 45066660, 10169, 0, 10, 5] //coin buy
[2, 20, 45066660, 9827, 0, 1, 5]
[2, 25, 44202850, 9996, 0, 1, 5]
[251, 10, 45657611, 10242, 0, 1, 6] // vip





*/
