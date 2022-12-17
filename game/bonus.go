package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Bonus(reader io.Reader) {

	bonus := &server.Bonus{}

	err := leb128.Unmarshal(reader, bonus)
	if err != nil {
		game.LogErrorPacket(bonus, err)
		return
	}

	game.bot.CanCollect = bonus.CanCollect

	game.bot.BonusDay = bonus.Day

	if game.bot.CanCollect {
		game.Send(client.BONUS, client.Bonus{})
	}

	game.LogReadPacket(*bonus)
}
