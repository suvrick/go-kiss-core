package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const GAME_REWARDS_GET types.PacketClientType = 11

// GAME_REWARDS_GET(11)
type GameRewardsGet struct {
	RewardID uint64
}

func (gameRewardsGet GameRewardsGet) String() string {
	return "GAME_REWARDS_GET(11)"
}

func (gameRewardsGet *GameRewardsGet) Marshal() ([]byte, error) {
	return leb128.WriteUInt64(nil, gameRewardsGet.RewardID)
}
