package packets

type Packet struct {
	Len    uint64
	ID     int64
	Type   uint64
	Format string
	Name   string
	Error  error
	Buffer []byte
	Params []interface{}
}

type Reward struct {
	ID    int64
	Count int64
}

type Bot struct {
	Result         int64
	GameID         int64
	Balance        int64
	BalanceHistory []int64
	Name           string
	CanCollect     int64
	BonusDay       int64
	Rewards        []Reward
	RewardsHistory []Reward
	RewardsGot     []int64
}

type Mask struct {
	PacketID uint64
	Index    int
	Type     string
	Name     string
}

// LOGIN(4);  status:B, inner_id:I, balance:I
// BALANCE(7); bottles:I, reason:B
// GAME_REWARDS(13); [id:I, count:I]
// BONUS(17); can_collect:B, day:B
// REWARD_GOT(315); owner_id:I, reward_id:I
var masks = []Mask{
	{Name: "Result", PacketID: 4, Index: 0},
	{Name: "GameID", PacketID: 4, Index: 1},
	{Name: "Balance", PacketID: 4, Index: 2},
	{Name: "Balance", PacketID: 7, Index: 0},
	{Name: "BalanceHistory", PacketID: 4, Index: 2},
	{Name: "BalanceHistory", PacketID: 7, Index: 0},
	{Name: "Rewards", PacketID: 13, Index: 0},
	{Name: "RewardsGot", PacketID: 315, Index: 1},
	{Name: "CanCollect", PacketID: 17, Index: 0},
	{Name: "BonusDay", PacketID: 17, Index: 1},
}

func (p *Packet) Fill(bot *Bot) {

	for _, mask := range masks {
		if p.Type == mask.PacketID {

			if mask.Index >= len(p.Params) {
				break
			}

			stuff := p.Params[mask.Index]

			switch mask.Name {
			case "Result":
				if value, ok := stuff.(int64); ok {
					bot.Result = value
				}
			case "GameID":
				if value, ok := stuff.(int64); ok {
					bot.GameID = value
				}
			case "Balance":
				if value, ok := stuff.(int64); ok {
					bot.Balance = value
				}
			case "BalanceHistory":

				if bot.BalanceHistory == nil {
					bot.BalanceHistory = make([]int64, 0)
				}

				if value, ok := stuff.(int64); ok {
					bot.BalanceHistory = append(bot.BalanceHistory, value)
				}
			case "Rewards":
				if bot.Rewards == nil {
					bot.Rewards = make([]Reward, 0)
				}

				if arr, ok := stuff.([]interface{}); ok {
					for _, r := range arr {
						if r2, ok := r.([]interface{}); ok {

							if len(r2) != 2 {
								continue
							}

							reward := Reward{}

							if item, ok := r2[0].(int64); ok {
								reward.ID = item
							} else {
								continue
							}

							if item, ok := r2[1].(int64); ok {
								reward.Count = item
							} else {
								continue
							}

							bot.Rewards = append(bot.Rewards, reward)
						}
					}
				}
			case "RewardsGot":

				if bot.RewardsGot == nil {
					bot.RewardsGot = make([]int64, 0)
				}

				if value, ok := stuff.(int64); ok {
					bot.RewardsGot = append(bot.RewardsGot, value)
				}
			case "CanCollect":
				if value, ok := stuff.(int64); ok {
					bot.CanCollect = value
				}
			case "BonusDay":
				if value, ok := stuff.(int64); ok {
					bot.BonusDay = value
				}
			}
		}
	}
}
