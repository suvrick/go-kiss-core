package server

const BONUS PacketServerType = 17

// BONUS(17) "BB"
type Bonus struct {
	CanCollect bool
	Day        int
}