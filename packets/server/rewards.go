package server

import (
	"time"

	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/types"
)

const REWARDS types.PacketServerType = 13

// REWARDS(13) "[BB]"
type Rewards struct {
	Rewards []models.Reward
}

func (packet *Rewards) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	for _, reward := range packet.getRewards() {
		if reward.Count > 0 {
			time.Sleep(time.Microsecond * 100)
			game.Send(client.GAME_REWARDS_GET, &client.GameRewardsGet{
				RewardID: reward.ID,
			})
		}
	}

	return nil
}

func (packet *Rewards) getRewards() (rewards []models.Reward) {
	for _, v := range packet.Rewards {
		if v.ID > 0 && v.Count > 0 {
			rewards = append(rewards, v)
		}
	}
	return
}
