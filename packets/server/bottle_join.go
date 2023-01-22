package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_JOIN types.PacketServerType = 26

// BOTTLE_JOIN(26) "IB"
type BottleJoin struct {
	PlayerID types.I
	Position types.B
}

func (packet *BottleJoin) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	room.Players[packet.PlayerID] = &models.Player{
		PlayerID: packet.PlayerID,
	}

	game.Send(client.REQUEST, &client.Request{
		Players: []types.I{packet.PlayerID},
		Mask:    INFOMASK,
	})

	return nil
}
