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

func (packet *Login) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	hiro.Result = packet.Result
	hiro.ID = packet.GameID
	hiro.Balance = packet.Balance

	switch hiro.Result {
	case 0:
		game.Send(client.REQUEST, client.Request{
			Players: []types.I{
				hiro.ID,
			},
			Mask: 328588,
		})
	default:
		game.Close()
	}

	return nil
}
