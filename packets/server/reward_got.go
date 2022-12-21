package server

import "github.com/suvrick/go-kiss-core/types"

const REWARD_GOT PacketServerType = 315

// REWARD_GOT(315) "II"
type RewardGot struct {
	UserID   types.I
	RewardID types.I
}
