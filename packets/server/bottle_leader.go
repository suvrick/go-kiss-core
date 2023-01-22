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

func (packet *BottleLeader) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	room.LeaderID = packet.LeaderID
	if packet.LeaderID == hiro.ID {
		go func() {
			log.Println("I am leader!")
			<-time.After(time.Second * time.Duration(5))
			log.Println("I am rolled bottle!")
			game.Send(client.BOTTLE_ROLL, &client.BottleRoll{
				IntField: 0,
			})
		}()
	}

	return nil
}
