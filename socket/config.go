package socket

import (
	"log"
	"net/http"
	"os"
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
		ConnectTimeout: time.Second * 3000,
		Head: http.Header{
			"Origin": {
				"https://inspin.me",
			},
			"User-Agent": {
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			},
		},
		Logger:        log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		TimeInTheGame: 7775,
	}
}
