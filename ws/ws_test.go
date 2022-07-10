package ws_test

import (
	"log"
	"testing"

	"github.com/suvrick/go-kiss-core/packets/net"
	"github.com/suvrick/go-kiss-core/ws"
)

func TestGameSocket(t *testing.T) {

	parser := net.NewParser()
	parser.Initialize()
	if parser.Error != nil {
		log.Fatalln(parser.Error.Error())
	}

	config := ws.GetDefaultGameSocketConfig()

	gs := ws.NewGameSocket(config)
	gs.Run()
	login_params := []interface{}{1000015, 32, 4, "200514254f3678c2f79cb18760ba048d", 0, ""}
	gs.Send(4, login_params)
	<-gs.Done
}
