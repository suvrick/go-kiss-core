package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const CHAT_WHISPER types.PacketServerType = 38

// CHAT_WHISPER(38) "ISBIB"
type ChatWhisper struct {
	WriterID   types.I
	Message    types.S
	ByteField  types.B
	IntField2  types.I
	ByteField2 types.B
}

func (packet *ChatWhisper) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	return nil
}
