package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) RewardGot(reader io.Reader) (interface{}, error) {
	rewardGot := &server.RewardGot{}

	err := leb128.Unmarshal(reader, rewardGot)
	if err != nil {
		return rewardGot, err
	}

	if rewardGot.UserID == game.bot.GameID {
		if game.bot.RewardGot == nil {
			game.bot.RewardGot = make([]int, 0)
		}
		game.bot.RewardGot = append(game.bot.RewardGot, rewardGot.RewardID)
	}

	game.socket.Logger.Printf("Read [%T] %+v\n", rewardGot, rewardGot)

	return rewardGot, nil
}
