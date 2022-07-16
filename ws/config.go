package ws

import (
	"log"
	"net/http"
	"time"
)

type GameSocketConfig struct {
	*SocketConfig
	MaskInfo int
	Logger   *log.Logger
}

type SocketConfig struct {
	Host          string
	Head          http.Header
	Timeout       time.Duration
	TimeInTheGame int
	Logger        *log.Logger
	Load          int
}

func GetDefaultGameSocketConfig() *GameSocketConfig {
	logger := log.Default()
	sc := GetDefaultSocketConfig()
	sc.Logger = logger
	return &GameSocketConfig{
		SocketConfig: sc,
		MaskInfo:     1114252,
		Logger:       logger,
	}
}

func GetDefaultSocketConfig() *SocketConfig {
	return &SocketConfig{
		Host:    "wss://bottlews.itsrealgames.com",
		Timeout: time.Second * 30,
		Head: http.Header{
			"Origin": {
				"https://inspin.me",
			},
		},
		Logger:        log.Default(),
		TimeInTheGame: 3000,
		Load:          10,
	}
}
