package server

import (
	"log"
	"time"

	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_ROLL types.PacketServerType = 29

// BOTTLE_ROLL(29) "II,II"
type BottleRoll struct {
	LeaderID  types.I
	RollerID  types.I
	IntField  types.I `pack:"optional"`
	IntField2 types.I `pack:"optional"`
}

func (packet *BottleRoll) Use(self *models.Bot, game interfaces.IGame) error {
	self.Room.LeaderID = packet.LeaderID
	self.Room.RollerID = packet.RollerID
	if packet.LeaderID == self.SelfID || packet.RollerID == self.SelfID {
		go func() {
			if packet.RollerID == self.SelfID {
				log.Println("I am kissed as roller!")
			} else {
				log.Println("I am kissed as leader!")
			}

			<-time.After(time.Second * time.Duration(5))

			game.Send(client.BOTTLE_KISS, &client.BottleKiss{
				Answer: 1,
			})
		}()
	}
	// 3.2 MB
	game.UpdateSelfEmit()
	return nil
}
