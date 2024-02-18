package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const GAME_REWARDS_GET types.PacketClientType = 11

type GameRewardsGet struct {
	RewardID byte
}

func (gameRewardsGet *GameRewardsGet) Marshal() ([]byte, error) {
	return leb128.WriteByte(nil, gameRewardsGet.RewardID)
}
