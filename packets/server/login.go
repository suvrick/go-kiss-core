package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const LOGIN types.PacketServerType = 4

// LOGIN(4) "B,II"
type Login struct {
	Result  models.LoginResultType
	GameID  types.I `pack:"optional"`
	Balance types.I `pack:"optional"`
}

func (packet *Login) Use(self *models.Bot, game interfaces.IGame) error {

	self.Result = packet.Result
	self.SelfID = packet.GameID
	self.Balance = packet.Balance

	game.UpdateSelfEmit()

	switch self.Result {
	case 0:
		game.Send(client.REQUEST, client.Request{
			Players: []types.I{
				self.SelfID,
			},
			Mask: 328588,
		})
	default:
		game.Close()
	}

	return nil
}
