package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const REWARD_GOT types.PacketServerType = 315

// REWARD_GOT(315) "II"
type RewardGot struct {
	PlayerID uint64
	RewardID uint64
}

func (rewardGot *RewardGot) Unmarshal(r *bytes.Reader) error {
	var err error

	rewardGot.PlayerID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	rewardGot.RewardID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	return err
}
