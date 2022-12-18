package client

const BOTTLE_PLAY PacketClientType = 26

// BOTTLE_PLAY(26) "B,B"
type BottlePlay struct {
	ByteField  byte
	ByteField2 byte
}
