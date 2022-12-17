package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Info(reader io.Reader) {

	info := &server.Info{}

	err := leb128.Unmarshal(reader, info)
	if err != nil {
		game.LogErrorPacket(info, err)
		return
	}

	game.bot.Name = info.Name

	game.bot.Sex = info.Sex

	game.bot.Avatar = info.Avatar

	game.bot.AvatarID = info.AvatarID

	game.bot.Profile = info.Profile

	game.bot.Status = info.Status

	game.LogReadPacket(info)
}
