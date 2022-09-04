package bot

import (
	"time"

	"github.com/suvrick/go-kiss-core/packets/server"
)

type Bot struct {
	ID   uint64
	Time time.Time

	//Login
	GameID int
	Result byte

	// Info
	Name    string
	Sex     byte
	Avatar  string
	Profile string
	Status  string
	//Bonus
	CanCollect bool
	BonusDay   int

	//Balance
	Balance        int
	BalanceHistory []int

	//RewardGot
	RewardGot []int

	//Rewards
	Rewards []server.Reward

	// BalanceItems
	BalanceItems []server.BalanceItem

	CollectionsPoints uint16
}
