package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
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

	game.bot.IsNeedSendBonus = game.bot.CanCollect

	game.LogReadPacket(*bonus)
}
