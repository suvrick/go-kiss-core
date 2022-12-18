package client

const BOTTLE_KICK PacketClientType = 31

// BOTTLE_KICK(31) "I"
type BottleKick struct {
	IntField int64
}
