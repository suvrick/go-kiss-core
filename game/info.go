package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) Info(reader io.Reader) (interface{}, error) {
	info := &server.Info{}

	err := leb128.Unmarshal(reader, info)
	if err != nil {
		return info, err
	}

	game.bot.Name = info.Name
	game.bot.Sex = info.Sex
	game.bot.Avatar = info.Avatar
	game.bot.Profile = info.Profile
	game.bot.Status = info.Status

	game.socket.Logger.Printf("Read [%T] %+v\n", info, info)

	return info, nil
}
