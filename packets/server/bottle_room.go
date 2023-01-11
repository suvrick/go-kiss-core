package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROOM types.PacketServerType = 25

// BOTTLE_ROOM(25) "III[I][I]"
type BottleRoom struct {
	RoomID    types.I
	IntField2 types.I
	IntField3 types.I
	Players   []types.I
	IntArray2 []types.I
}

func (packet *BottleRoom) Use(self *models.Bot, game interfaces.IGame) error {

	self.Room = &models.Room{
		RoomID:  packet.RoomID,
		Players: make(map[types.I]*models.Player),
	}

	for _, p := range packet.Players {
		self.Room.Players[p] = &models.Player{
			PlayerID: p,
		}

		game.Send(client.REQUEST, &client.Request{
			Players: []types.I{p},
			Mask:    INFOMASK,
		})
	}

	game.UpdateSelfEmit()
	return nil
}
