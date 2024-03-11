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
			"Accept": {
				"*/*",
			},
		},
		Logger:        log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		TimeInTheGame: 7775,
	}
}

/*

Accept-Encoding:
gzip, deflate, br
Accept-Language:
ru,en;q=0.9,en-GB;q=0.8,en-US;q=0.7
Cache-Control:
no-cache
Connection:
Upgrade
Host:
bottlews.itsrealgames.com
Origin:
https://inspin.me
Pragma:
no-cache
Sec-Websocket-Extensions:
permessage-deflate; client_max_window_bits
Sec-Websocket-Key:
7hz52VOx48CkS33llBmrkQ==
Sec-Websocket-Version:
13
Upgrade:
websocket
User-Agent:
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0

*/
