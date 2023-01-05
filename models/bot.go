package models

import (
	"github.com/suvrick/go-kiss-core/types"
)

type Bot struct {
	SelfID          types.I
	Result          LoginResultType
	Balance         types.I
	CanCollect      types.B
	BonusDay        types.B
	CollectionPoint types.I

	Player
}

// type Bot struct {
// 	Log []string `json:"_log"`

// 	ID   string    `json:"id"`
// 	Time time.Time `json:"time"`

// 	//Login
// 	GameID       uint64 `json:"game_id"`
// 	Result       uint16 `json:"result"`
// 	ResultString string `json:"result_string"`
// 	// Info
// 	Name     string `json:"name"`
// 	Sex      byte   `json:"sex"`
// 	Avatar   string `json:"avatar"`
// 	AvatarID byte   `json:"avatar_id"`
// 	Profile  string `json:"profile"`
// 	Status   string `json:"status"`
// 	//Bonus
// 	CanCollect bool `json:"can_collect"`
// 	BonusDay   int  `json:"bonus_day"`
// 	//Balance
// 	Balance        uint   `json:"balance"`
// 	BalanceHistory []uint `json:"balance_history"`
// 	//RewardGot
// 	RewardGot []int `json:"reward_got"`
// 	//Rewards
// 	Rewards []server.Reward `json:"rewards"`
// 	// BalanceItems
// 	BalanceItems []server.BalanceItem `json:"-"`

// 	CollectionsPoints uint16 `json:"collections_points"`

// 	ErrorHistory []string `json:"errors"`
// }
