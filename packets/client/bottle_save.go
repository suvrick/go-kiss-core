package client

const BOTTLE_SAVE PacketClientType = 30

// BOTTLE_SAVE(30) "I"
type BottleSave struct {
	IntField int64
}
