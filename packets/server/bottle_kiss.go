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

func (packet *BottleKiss) Use(self *models.Bot, game interfaces.IGame) error {
	if packet.PlayerID == self.Room.LeaderID {
		self.Room.KissAnswerLeader = packet.Answer
	} else if packet.PlayerID == self.Room.RollerID {
		self.Room.KissAnswerRoller = packet.Answer
	}

	if currentPlayer, ok := self.Room.Players[packet.PlayerID]; ok {
		if packet.Answer == 1 {
			currentPlayer.KissedRoom += 1
			currentPlayer.KissedDay += 1
			currentPlayer.Kissed += 1
		}
	}

	game.UpdateSelfEmit()
	return nil
}
