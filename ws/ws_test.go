package ws_test

import (
	"log"
	"testing"

	"github.com/suvrick/go-kiss-core/packets/meta"
	"github.com/suvrick/go-kiss-core/ws"
)

func TestGameSocket(t *testing.T) {

	meta := meta.NewMeta()
	if meta.Error != nil {
		log.Fatalln(meta.Error.Error())
	}

	config := ws.GetDefaultGameSocketConfig()

	gs := ws.NewGameSocket(config)
	gs.Run()
	login_params := []interface{}{1000015, 32, 4, "200514254f3678c2f79cb18760ba048d", 0, ""}
	gs.Send(4, login_params)

	gs2 := ws.NewGameSocket(config)
	gs2.Run()
	login_params2 := []interface{}{1000015, 32, 4, "200514254f3678c2f79cb18760ba048d", 0, ""}
	gs2.Send(4, login_params2)

	<-gs.Done
	<-gs2.Done
}
