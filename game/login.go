package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Login(reader io.Reader) (interface{}, error) {

	login := &server.Login{}

	err := leb128.Unmarshal(reader, login)
	if err != nil {
		return login, err
	}

	game.bot.Result = login.Result
	game.bot.GameID = login.GameID
	game.bot.Balance = login.Balance

	game.bot.BalanceHistory = make([]int, 0)
	game.bot.BalanceHistory = append(game.bot.BalanceHistory, login.Balance)

	game.socket.Logger.Printf("Read [%T] %+v\n", login, login)

	switch login.Result {
	case 0:
	default:
		game.GameOver()
	}

	return login, nil
}
