package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEADER types.PacketServerType = 28

// BOTTLE_LEADER(28) "I"
type BottleLeader struct {
	LeaderID types.I
}

func (packet *BottleLeader) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
