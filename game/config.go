package game

import (
	"log"

	"github.com/suvrick/go-kiss-core/socket"
)

type GameConfig struct {
	*socket.SocketConfig
	MaskInfo int
	Logger   *log.Logger
}

func GetDefaultGameConfig() *GameConfig {
	logger := log.Default()
	sc := socket.GetDefaultSocketConfig()
	sc.Logger = logger
	return &GameConfig{
		SocketConfig: sc,
		Logger:       logger,
	}
}
