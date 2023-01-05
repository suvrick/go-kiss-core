package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BALANCE types.PacketServerType = 7

// BALANCE(7) "bottles:I, reason:B"
type Balance struct {
	Bottles types.I
	Reason  types.B `pack:"optional"`
}

func (packet *Balance) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
