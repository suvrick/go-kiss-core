package client

import "github.com/suvrick/go-kiss-core/types"

const BUY PacketClientType = 6

// BONUS(6) "IIIIBI,B"
type Buy struct {
	BuyType    types.I
	Coin       types.I
	PlayerID   types.I
	PrizeID    types.I
	ByteFiald  types.B
	Count      types.I
	ByteFiald2 types.B `pack:"optional"`
}

/*

[2, 30, 45066660, 10169, 0, 1, 5]  //free buy
[2, 30, 45066660, 10169, 0, 10, 5] //coin buy
[2, 20, 45066660, 9827, 0, 1, 5]
[2, 25, 44202850, 9996, 0, 1, 5]
[251, 10, 45657611, 10242, 0, 1, 6] // vip





*/
