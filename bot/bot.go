package bot

type LogItem string

type BalanceItem struct {
	Name   string `json:"name"`
	Common uint64 `json:"common"`
	Daily  uint64 `json:"daily"`
}

type Bot struct {
	Status            int64         `json:"status"`
	GameID            uint64        `json:"game_id"`
	Balance           uint64        `json:"balance"`
	BonusDay          int64         `json:"bonus_day"`
	IsBonus           int64         `json:"is_bonus"`
	Name              string        `json:"name"`
	Avatar            string        `json:"avatar"`
	VipDays           float64       `json:"vip_days"`
	Rewards           []string      `json:"rewards"`
	CollectionsPoints uint64        `json:"collections_points"`
	BalanceItems      []BalanceItem `json:"balance_items"`
	Log               []LogItem     `json:"log"`
}
