package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Balance(reader io.Reader) (interface{}, error) {
	balance := &server.Balance{}

	err := leb128.Unmarshal(reader, balance)
	if err != nil {
		return balance, err
	}

	game.bot.Balance = balance.Bottles
	game.bot.BalanceHistory = append(game.bot.BalanceHistory, game.bot.Balance)

	game.socket.Logger.Printf("Read [%T] %+v\n", balance, balance)

	return balance, nil
}
