package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BALANCE_ITEMS types.PacketServerType = 310

// BALANCE_ITEMS(310) "[BII]"
type BalanceItems struct {
	Items []models.BalanceItem
}

func (packet *BalanceItems) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
