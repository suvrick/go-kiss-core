package client

import "github.com/suvrick/go-kiss-core/types"

const GAME_REWARDS_GET types.PacketClientType = 11

type GameRewardsGet struct {
	RewardID types.B
}
