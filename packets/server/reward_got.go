package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const REWARD_GOT types.PacketServerType = 315

// REWARD_GOT(315) "II"
type RewardGot struct {
	UserID   types.I
	RewardID types.I
}

func (packet *RewardGot) Use(self *models.Bot, game interfaces.IGame) error {
	return nil
}
