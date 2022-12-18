package server

const BOTTLE_ROLL PacketServerType = 29

// BOTTLE_ROLL(29) "II,II"
type BottleRoll struct {
	LeaderID  int64
	RollerID  int64
	IntField  int64
	IntField2 int64
}
