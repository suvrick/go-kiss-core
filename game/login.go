package game

import (
	"fmt"
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

	game.bot.Result = uint16(login.Result)
	game.bot.ResultString = login.Result.String()
	game.bot.GameID = login.GameID
	game.bot.Balance = login.Balance

	game.bot.BalanceHistory = make([]uint, 0)
	game.bot.BalanceHistory = append(game.bot.BalanceHistory, login.Balance)

	game.socket.Logger.Printf("Read [%T] %+v\n", login, login)

	fmt.Println(login.Result)

	switch login.Result {
	case server.Success:
	case server.Exist:
		// game.Send(server.LOGIN, pack)
	default:
		game.GameOver()
	}

	return login, nil
}
