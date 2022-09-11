package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) BalanceItems(reader io.Reader) {

	balanceItems := &server.BalanceItems{}

	err := leb128.Unmarshal(reader, balanceItems)
	if err != nil {
		game.LogErrorPacket(balanceItems, err)
		return
	}

	game.bot.BalanceItems = balanceItems.Items

	game.LogReadPacket(*balanceItems)

}
