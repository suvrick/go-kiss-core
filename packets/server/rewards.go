package server

import "github.com/suvrick/go-kiss-core/types"

const REWARDS PacketServerType = 13

// REWARDS(13) "[BB]"
type Rewards struct {
	Rewards []Reward
}

type Reward struct {
	ID    types.B
	Count types.B
}
