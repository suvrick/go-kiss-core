package server

const BALANCE_ITEMS PacketServerType = 310

type BalanceItem struct {
	BalanceType byte
	A           int
	B           int
}

// BALANCE_ITEMS(310) "bottles:I, reason:B"
type BalanceItems struct {
	Items []BalanceItem
}

func GetBalanceItemName(t int64) string {
	switch t {
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
