package server

const REWARDS PacketServerType = 13

type Reward struct {
	ID    int
	Count int
}

// REWARDS(13) "II"
type Rewards struct {
	Rewards []Reward
}
