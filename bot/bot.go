package bot

import (
	"time"

	"github.com/suvrick/go-kiss-core/packets/server"
)

type Bot struct {
	ID   uint64    `json:"id,omitempty"`
	Time time.Time `json:"-"`

	//Login
	GameID int  `json:"game_id,omitempty"`
	Result byte `json:"result,omitempty"`

	// Info
	Name    string `json:"name,omitempty"`
	Sex     byte   `json:"sex,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Profile string `json:"profile,omitempty"`
	Status  string `json:"status,omitempty"`
	//Bonus
	CanCollect bool `json:"can_collect,omitempty"`
	BonusDay   int  `json:"bonus_day,omitempty"`

	//Balance
	Balance        int   `json:"balance,omitempty"`
	BalanceHistory []int `json:"balance_history,omitempty"`

	//RewardGot
	RewardGot []int `json:"reward_got,omitempty"`

	//Rewards
	Rewards []server.Reward `json:"rewards,omitempty"`

	// BalanceItems
	BalanceItems []server.BalanceItem `json:"balance_items,omitempty"`

	CollectionsPoints uint16 `json:"collections_points,omitempty"`
}
