package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) CollectionsPoints(reader io.Reader) (interface{}, error) {
	collectionsPoints := &server.CollectionsPoints{}

	err := leb128.Unmarshal(reader, collectionsPoints)
	if err != nil {
		return collectionsPoints, err
	}

	game.bot.CollectionsPoints = collectionsPoints.Points

	game.socket.Logger.Printf("Read [%T] %+v\n", collectionsPoints, collectionsPoints)

	return collectionsPoints, nil
}
