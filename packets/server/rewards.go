package server

import "github.com/suvrick/go-kiss-core/types"

const REWARDS PacketServerType = 13

// REWARDS(13) "[II]"
type Rewards struct {
	Rewards []Reward
}

type Reward struct {
	ID    types.I
	Count types.I
}
