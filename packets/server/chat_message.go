package server

import "github.com/suvrick/go-kiss-core/types"

const CHAT_MESSAGE PacketServerType = 37

// CHAT_MESSAGE(37) "BIIS,II"
type ChatMessage struct {
	ByteField types.B
	IntField  types.I
	WriterID  types.I
	Message   types.S
	IntField3 types.I `pack:"optional"`
	IntField4 types.I `pack:"optional"`
}
