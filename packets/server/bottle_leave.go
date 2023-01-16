package server

import (
	"log"
	"time"

	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEAVE types.PacketServerType = 27

// BOTTLE_LEAVE(27) "I"
type BottleLeave struct {
	PlayerID types.I
}

func (packet *BottleLeave) Use(self *models.Bot, game interfaces.IGame) error {
	delete(self.Room.Players, packet.PlayerID)

	if !self.Find() {
		time.Sleep(1000)

		game.Send(client.MOVE, client.Move{
			PlayerID: self.HiroID,
		})

		log.Printf("[%d] I am start search hiro %d\n", self.SelfID, self.HiroID)
	}

	game.UpdateSelfEmit()
	return nil
}
