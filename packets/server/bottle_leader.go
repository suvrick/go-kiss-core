package server

import (
	"log"
	"time"

	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEADER types.PacketServerType = 28

// BOTTLE_LEADER(28) "I"
type BottleLeader struct {
	LeaderID types.I
}

func (packet *BottleLeader) Use(self *models.Bot, game interfaces.IGame) error {
	self.Room.LeaderID = packet.LeaderID
	if packet.LeaderID == self.SelfID {
		go func() {
			log.Println("I am leader!")
			<-time.After(time.Second * time.Duration(5))
			log.Println("I am rolled bottle!")
			game.Send(client.BOTTLE_ROLL, &client.BottleRoll{
				IntField: 0,
			})
		}()
	}

	game.UpdateSelfEmit()
	return nil
}
