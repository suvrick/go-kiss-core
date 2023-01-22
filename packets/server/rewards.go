package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const REWARDS types.PacketServerType = 13

// REWARDS(13) "[BB]"
type Rewards struct {
	Rewards []models.Reward
}

func (packet *Rewards) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	return nil
}
