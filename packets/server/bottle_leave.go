package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEAVE types.PacketServerType = 27

// BOTTLE_LEAVE(27) "I"
type BottleLeave struct {
	PlayerID types.I
}

func (packet *BottleLeave) Use(self *models.Bot, game interfaces.IGame) error {
	delete(self.Room.Players, packet.PlayerID)
	game.UpdateSelfEmit()
	return nil
}
