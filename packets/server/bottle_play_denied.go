package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_PLAY_DENIED types.PacketServerType = 24

// BOTTLE_PLAY_DENIED(24) "B"
type BottlePlayDenied struct {
	ByteField types.B
}

func (packet *BottlePlayDenied) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	return nil
}
