package packets

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

/*

   BalanceType[BalanceType["Kicks"] = 0] = "Kicks";
   BalanceType[BalanceType["Saves"] = 1] = "Saves";
   BalanceType[BalanceType["KissPriority"] = 2] = "KissPriority";
   BalanceType[BalanceType["Video"] = 3] = "Video";
   BalanceType[BalanceType["Gifts"] = 4] = "Gifts";
   BalanceType[BalanceType["Hearts"] = 5] = "Hearts";
   BalanceType[BalanceType["Stickers"] = 6] = "Stickers";
   BalanceType[BalanceType["UniqueGifts"] = 7] = "UniqueGifts";
   BalanceType[BalanceType["PlayersTape"] = 8] = "PlayersTape";
   BalanceType[BalanceType["AdventCell"] = 9] = "AdventCell";
   BalanceType[BalanceType["Admire"] = 10] = "Admire";
   BalanceType[BalanceType["RouletteRoll"] = 11] = "RouletteRoll";
   BalanceType[BalanceType["DeclineAdmire"] = 14] = "DeclineAdmire";
   BalanceType[BalanceType["FestivalPoints"] = 15] = "FestivalPoints";
   BalanceType[BalanceType["MaxType"] = 16] = "MaxType";

*/
