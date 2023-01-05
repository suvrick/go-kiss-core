package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_KISS types.PacketServerType = 30

// BOTTLE_KISS(30) "IB"
type BottleKiss struct {
	PlayerID types.I
	Answer   models.KissAnswer
}

func (packet *BottleKiss) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
