package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const CHAT_MESSAGE types.PacketServerType = 37

// CHAT_MESSAGE(37) "BIIS,II"
type ChatMessage struct {
	ByteField types.B
	IntField  types.I
	WriterID  types.I
	Message   types.S
	IntField3 types.I `pack:"optional"`
	IntField4 types.I `pack:"optional"`
}

func (packet *ChatMessage) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
