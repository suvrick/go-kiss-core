package main

import (
	"bytes"
	"log"

	"github.com/go-interpreter/wagon/wasm/leb128"
	"github.com/gorilla/websocket"
	gokisscore "github.com/suvrick/go-kiss-core"
)

func main() {

	filler := gokisscore.NewFiller()
	filler.ParseScript()

	if filler.Error != nil {
		log.Fatalln(filler.Error.Error())
	}

	// connection.sendData(ClientPacketType_1.ClientPacketType.LOGIN, loginData.id, loginData.netType, deviceType, loginData.key, loginData.oauth || 0, loginData.accessToken || "", referrer, tag, 0, "", Localization_1.Localization.ROOM_LANGUAGE, "", gender, captcha);
	b := filler.CreateClientPacket(gokisscore.ClientPacketType(4), 114941701, 32, 4, "33513e2ce85cabfd6ec59d827aa28cea")

	ws := websocket.Dialer{}
	con, _, err := ws.Dial("wss://bottlews.itsrealgames.com", nil)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	msg := make([]byte, 0)
	msg = leb128.AppendUleb128(msg, uint64(len(b)+1))
	msg = append(msg, 0)
	msg = append(msg, b...)

	err = con.WriteMessage(2, msg)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for {

		_, msg, err := con.ReadMessage()

		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		if len(msg) < 3 {
			continue
		}

		reader := bytes.NewReader(msg)

		_, _ = leb128.ReadVarUint32(reader)
		_, _ = leb128.ReadVarUint32(reader)
		msgType, _ := leb128.ReadVarUint32(reader)

		//log.Printf("Recv >> msgType: %d,msgID: %d, msgLen: %d\n", msgType, msgID, msgLen)
		name, data := filler.CreateServerPacket(gokisscore.ServerPacketType(msgType), reader)
		log.Printf("RECV >> name: %s, data: %v\n", name, data)
	}

	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1571552893425&&userId=114941701&sessionKey=67f5e4f7a90144c5eba1b91694132904&authKey=33513e2ce85cabfd6ec59d827aa28cea&net_type=32&useApiType=sa&email=tsd_8326@mail.ru&locale=RU"
	//url := "https://m.inspin.me/build/v431/?api=ok&5&sig=13c32cf30efbedcbfd7afa689f9fa5b6&mob=true&refplace=user_apps_portlet&session_key=-s-6622x0r.b2d-yvr6022bQSS26241xUv3433cvvq4c232NwOZ376bTuxd62-4TxN41178Msr1a63fMXr45827NWo277c0TUq8b265Ne&session_secret_key=7f16a539d8d56e2fca91be30ac68c7e3&auth_sig=d04204aa3632ab2dcf9e6f2b41084b57&api_server=https%3A%2F%2Fapi.ok.ru%2F&payment_server=https%3A%2F%2Fmpay.ok.ru&lang=ru&application_key=CBADLOPFABABABABA&mob_platform=androidweb&logged_user_id=952642558996&type=ok"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000937&sessionKey=6c92d6b72649d11dc6a436777c8f3fd&authKey=3109ee5271d550293cf9b2fb494556f1&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=95ac81be3e64f73fd56e4b12e004c9d0&vid=9914552860808032409&oid=9914552860808032409&app_id=543574&authentication_key=27861583a5de8bd7c1aa5955dfe757b8&session_expire=1623244207&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=2add120d8f235755ee8b34801798f775&window_id=CometName_659cc22f549d0d25dc769cb39daa105f&referer_type=left_menu&version=1593"
	//url := "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593"

	// w, err := ws.New()

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// login := parser.NewLoginParams(url)
	// w.Go(login)
	// log.Printf("%+v", w.Bot)
}
