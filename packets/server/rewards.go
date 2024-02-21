package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const REWARDS types.PacketServerType = 13

// REWARDS(13) "[BB]"
type Rewards struct {
	Items []Reward
}

type Reward struct {
	RewardID uint64
	Count    byte
}

func (p Rewards) String() string {
	return "REWARDS(13)"
}

/*
Read packet REWARDS with data [[[252,1]]]
GAME_REWARDS_GET type 11 id 425 data: [252,1]
Read packet RACE_LEADERS with data [0,0,1,[[40305073,0,0]]]
Read packet REWARDS with data [[]]
Read packet BALANCE with data [80]
*/

func (rewards *Rewards) Unmarshal(r *bytes.Reader) error {
	var err error
	var len uint64

	len, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	rewards.Items = make([]Reward, len)

	for len > 0 {

		rewards.Items[len-1].RewardID, err = leb128.ReadUInt64(r)
		if err != nil {
			return err
		}

		rewards.Items[len-1].Count, err = leb128.ReadByte(r)
		if err != nil {
			return err
		}

		len--
	}

	return err
}

// func (packet *Rewards) getRewards() (rewards []Reward) {
// 	for _, v := range rewards {
// 		if v.RewardID > 0 && v.Count > 0 {
// 			rewards = append(rewards, v)
// 		}
// 	}
// 	return
// }
