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
	RewardID byte
	Count    byte
}

func (p Rewards) String() string {
	return "REWARDS(13)"
}

func (rewards *Rewards) Unmarshal(r *bytes.Reader) error {
	var err error
	var len uint64

	len, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	rewards.Items = make([]Reward, len)

	for len > 0 {
		var item = Reward{}

		item.RewardID, err = leb128.ReadByte(r)
		if err != nil {
			return err
		}

		item.Count, err = leb128.ReadByte(r)
		if err != nil {
			return err
		}

		rewards.Items = append(rewards.Items, item)
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
