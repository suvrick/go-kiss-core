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

func (packet *CollectionsPoints) Use(self *models.Bot, game interfaces.IGame) error {
	self.CollectionPoint = packet.Points
	return nil
}
