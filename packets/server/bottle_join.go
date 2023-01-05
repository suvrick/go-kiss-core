package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_JOIN types.PacketServerType = 26

// BOTTLE_JOIN(26) "IB"
type BottleJoin struct {
	PlayerID types.I
	Position types.B
}

func (packet *BottleJoin) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
