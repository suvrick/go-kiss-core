package interfaces

import (
	"net/url"

	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

type IGame interface {
	Send(packetID types.PacketClientType, packet interface{})
	UpdateSelfEmit()
	Connection() error
	ConnectionWithProxy(proxy *url.URL) error
	Close()
	//ConnectionWithProxy() error
}

type IServerPacket interface {
	Use(self *models.Bot, game IGame) error
}
