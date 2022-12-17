package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Balance(reader io.Reader) {

	balance := &server.Balance{}

	err := leb128.Unmarshal(reader, balance)
	if err != nil {
		game.LogErrorPacket(balance, err)
		return
	}

	game.bot.Balance = balance.Bottles

	game.bot.BalanceHistory = append(game.bot.BalanceHistory, game.bot.Balance)

	game.LogReadPacket(*balance)
}
