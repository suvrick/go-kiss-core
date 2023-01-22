package server

import (
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const COLLECTIONS_POINTS types.PacketServerType = 130

// COLLECTIONS_POINTS(130) "I"
type CollectionsPoints struct {
	Points types.I
}

func (packet *CollectionsPoints) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {
	hiro.CollectionPoint = packet.Points
	return nil
}
