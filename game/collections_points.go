package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) CollectionsPoints(reader io.Reader) {

	collectionsPoints := &server.CollectionsPoints{}

	err := leb128.Unmarshal(reader, collectionsPoints)
	if err != nil {
		game.LogErrorPacket(collectionsPoints, err)
		return
	}

	game.bot.CollectionsPoints = collectionsPoints.Points

	game.LogReadPacket(*collectionsPoints)
}
