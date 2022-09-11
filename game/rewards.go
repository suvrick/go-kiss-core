package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Rewards(reader io.Reader) {
	rewards := &server.Rewards{}

	err := leb128.Unmarshal(reader, rewards)
	if err != nil {
		game.LogErrorPacket(rewards, err)
		return
	}

	game.LogReadPacket(*rewards)

	for i := range rewards.Rewards {
		reward := rewards.Rewards[i]
		if reward.Count > 0 {
			if game.bot.Rewards == nil {
				game.bot.Rewards = make([]server.Reward, 0)
			}

			game.bot.Rewards = append(game.bot.Rewards, reward)

			game.Send(client.GAME_REWARDS_GET,
				&client.GameRewardsGet{RewardID: reward.ID})

			game.bot.Live++

			return
		}
	}

	if EmptyRewars(rewards) {
		game.bot.Live--
	}
}

func EmptyRewars(rewards *server.Rewards) bool {

	for _, v := range rewards.Rewards {
		if v.Count > 0 {
			return false
		}
	}

	return true
}
