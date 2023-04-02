package game_test

import (
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-core/game"
	"github.com/suvrick/go-kiss-core/types"
)

const frameURL = "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593"

func TestGame(t *testing.T) {
	//loginData, err := frame.Parse3(frameURL)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	data := []interface{}{
		types.L(105345504),
		types.B(30),
		types.B(5),
		types.S("7b0a077a088b9e5169bcfc0bf2ee9ae8"),
		types.B(1),
		types.S("5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38"),
		types.I(0),
		types.I(0),
		types.B(0),
		types.S(""),
		types.B(0),
		types.S(""),
		types.B(0),
		types.S(""),
	}

	g := game.NewGame()
	err := g.Connect(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	g.Send(4, data)
	<-g.End()
}
