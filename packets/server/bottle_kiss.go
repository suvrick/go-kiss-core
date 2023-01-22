package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_KISS types.PacketServerType = 30

// BOTTLE_KISS(30) "IB"
type BottleKiss struct {
	PlayerID types.I
	Answer   models.KissAnswer
}

func (packet *BottleKiss) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	if packet.PlayerID == room.LeaderID {
		room.KissAnswerLeader = packet.Answer
	} else if packet.PlayerID == room.RollerID {
		room.KissAnswerRoller = packet.Answer
	}

	if currentPlayer, ok := room.Players[packet.PlayerID]; ok {
		if packet.Answer == 1 {
			currentPlayer.KissedRoom += 1
			currentPlayer.KissedDay += 1
			currentPlayer.Kissed += 1
		}
	}

	return nil
}
