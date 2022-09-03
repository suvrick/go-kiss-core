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

	return balance, nil
}
