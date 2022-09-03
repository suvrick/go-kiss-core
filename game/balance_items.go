package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) BalanceItems(reader io.Reader) (interface{}, error) {
	balanceItems := &server.BalanceItems{}

	err := leb128.Unmarshal(reader, balanceItems)
	if err != nil {
		return balanceItems, err
	}

	return balanceItems, nil
}
