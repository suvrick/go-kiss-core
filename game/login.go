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

	switch login.Result {
	case 0:
	default:
		game.GameOver()
	}

	return login, nil
}
