package client

const BOTTLE_PLAY PacketClientType = 26

// BOTTLE_PLAY(26) "B,B"
type BottlePlay struct {
	RoomID byte
	LangID byte
}
