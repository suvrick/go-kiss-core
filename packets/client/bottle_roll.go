package client

const BOTTLE_ROLL PacketClientType = 28

// BOTTLE_ROLL(28) "I"
type BottleRoll struct {
	IntField int64
}
