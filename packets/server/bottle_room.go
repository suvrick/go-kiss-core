package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_ROOM PacketServerType = 25

// BOTTLE_ROOM(25) "III[I][I]"
type BottleRoom struct {
	RoomID    types.I
	IntField2 types.I
	IntField3 types.I
	Players   []types.I
	IntArray2 []types.I
}
