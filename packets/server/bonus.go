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

func (packet *Bonus) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	hiro.CanCollect = packet.CanCollect

	hiro.BonusDay = packet.Day

	if hiro.CanCollect == 1 {
		game.Send(client.BONUS, &client.Bonus{})
	}

	return nil
}
