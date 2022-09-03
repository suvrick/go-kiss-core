package client

const GAME_REWARDS_GET PacketClientType = 11

type GameRewardsGet struct {
	RewardID int
}
