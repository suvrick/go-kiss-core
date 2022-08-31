package ws

import (
	"log"
	"net/http"
	"time"
)

type SocketConfig struct {
	Host          string
	Head          http.Header
	Timeout       time.Duration
	TimeInTheGame int
	Logger        *log.Logger
	MaskInfo      int
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
		TimeInTheGame: 30,
		MaskInfo:      1114252,
	}
}
