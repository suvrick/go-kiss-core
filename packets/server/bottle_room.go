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

func (packet *BottleRoom) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	room.RoomID = packet.RoomID
	room.Players = make(map[types.I]*models.Player)

	for _, p := range packet.Players {
		room.Players[p] = &models.Player{
			PlayerID: p,
		}

		game.Send(client.REQUEST, &client.Request{
			Players: []types.I{p},
			Mask:    INFOMASK,
		})
	}

	return nil
}
