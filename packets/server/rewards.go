package server

const REWARDS PacketServerType = 13

type Reward struct {
	ID    uint16
	Count uint16
}

// REWARDS(13) "II"
type Rewards struct {
	Rewards []Reward
}
