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

func (packet *BottleLeave) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	if room != nil && len(room.Players) > 0 {
		delete(room.Players, packet.PlayerID)
	}

	return nil
}
