package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const CONTEST_ITEMS types.PacketServerType = 10

// CONTEST_ITEMS(10) "B,B"
type ContestItems struct {
	ByteField  types.B
	ByteField2 types.B `pack:"optional"`
}

func (packet *ContestItems) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	return nil
}
