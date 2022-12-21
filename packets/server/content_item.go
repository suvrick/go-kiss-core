package server

import "github.com/suvrick/go-kiss-core/types"

const CONTEST_ITEMS PacketServerType = 10

// CONTEST_ITEMS(10) "B,B"
type ContestItems struct {
	ByteField  types.B
	ByteField2 types.B `pack:"optional"`
}
