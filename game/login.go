package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Login(reader io.Reader) {

	login := &server.Login{}

	err := leb128.Unmarshal(reader, login)
	if err != nil {
		game.LogErrorPacket(login, err)
		return
	}

	game.LogReadPacket(*login)

	game.bot.Result = uint16(login.Result)

	game.bot.ResultString = login.Result.String()

	switch login.Result {
	case server.Success:

		game.bot.GameID = login.GameID

		game.bot.Balance = login.Balance

		game.bot.BalanceHistory = make([]uint, 0)

		game.bot.BalanceHistory = append(game.bot.BalanceHistory, login.Balance)

		game.BuySend()

	case server.Exist:
		game.LoginSend(nil)
	default:
		game.GameOver()
	}
}
