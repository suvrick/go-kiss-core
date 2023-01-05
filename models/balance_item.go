package models

import "github.com/suvrick/go-kiss-core/types"

type BalanceItem struct {
	BalanceType BalanceType
	A           types.I
	B           types.I
}

type BalanceType types.B

func (b BalanceType) String() string {
	switch b {
	case 0:
		return "Kicks"
	case 1:
		return "Saves"
	case 2:
		return "KissPriority"
	case 3:
		return "Video"
	case 4:
		return "Gifts"
	case 5:
		return "Hearts"
	case 6:
		return "Stickers"
	case 7:
		return "UniqueGifts"
	case 8:
		return "PlayersTape"
	case 9:
		return "AdventCell"
	case 10:
		return "Admire"
	case 11:
		return "RouletteRoll"
	case 14:
		return "DeclineAdmire"
	case 15:
		return "FestivalPoints"
	case 16:
		return "MaxType"
	default:
		return "None"
	}
}
