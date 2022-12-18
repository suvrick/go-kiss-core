package server

const BOTTLE_ROOM PacketServerType = 25

// BOTTLE_ROOM(25) "III[I][I]"
type BottleRoom struct {
	IntField1 int64
	IntField2 int64
	IntField3 int64
	Players   []player
	IntArray2 []player
}

type player struct {
	PlayerID int64
}
