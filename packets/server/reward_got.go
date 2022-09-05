package server

const REWARD_GOT PacketServerType = 315

// REWARD_GOT(315) "II"
type RewardGot struct {
	UserID   uint64
	RewardID int
}
