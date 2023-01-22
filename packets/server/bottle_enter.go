package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ENTER types.PacketServerType = 35

// BOTTLE_ENTER(35) ""
type BottleEnter struct {
}

func (packet *BottleEnter) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	return nil
}
