package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Bonus(reader io.Reader) (interface{}, error) {
	bonus := &server.Bonus{}

	err := leb128.Unmarshal(reader, bonus)
	if err != nil {
		return bonus, err
	}

	if bonus.CanCollect == 1 {
		game.Send(client.BONUS, client.Bonus{})
	}

	return bonus, nil
}
