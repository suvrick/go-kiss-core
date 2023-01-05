package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BONUS types.PacketServerType = 17

// BONUS(17) "BB"
type Bonus struct {
	CanCollect types.B
	Day        types.B
}

func (packet *Bonus) Use(self *models.Bot, game interfaces.IGame) error {

	self.CanCollect = packet.CanCollect
	self.BonusDay = packet.Day

	game.UpdateSelfEmit()

	if self.CanCollect == 1 {
		game.Send(client.BONUS, &client.Bonus{})
	}

	return nil
}
