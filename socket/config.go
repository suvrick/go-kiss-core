package socket

import (
	"log"
	"net/http"
	"time"
)

type SocketConfig struct {
	Host           string
	Head           http.Header
	ConnectTimeout time.Duration
	TimeInTheGame  int
	Logger         *log.Logger
}

func GetDefaultSocketConfig() *SocketConfig {
	return &SocketConfig{
		Host:           "wss://bottlews.itsrealgames.com",
		ConnectTimeout: time.Second * 30,
		Head: http.Header{
			"Origin": {
				"https://inspin.me",
			},
		},
		Logger:        log.Default(),
		TimeInTheGame: 3,
	}
}
