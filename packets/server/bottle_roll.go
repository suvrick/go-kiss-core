package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROLL types.PacketServerType = 29

// BOTTLE_ROLL(29) "II,II"
type BottleRoll struct {
	LeaderID  types.I
	RollerID  types.I
	IntField  types.I `pack:"optional"`
	IntField2 types.I `pack:"optional"`
}

func (packet *BottleRoll) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
