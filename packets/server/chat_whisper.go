package server

import "github.com/suvrick/go-kiss-core/types"

const CHAT_WHISPER PacketServerType = 38

// CHAT_WHISPER(38) "ISBIB"
type ChatWhisper struct {
	WriterID   types.I
	Message    types.S
	ByteField  types.B
	IntField2  types.I
	ByteField2 types.B
}
