package main

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-core/frame"
	"github.com/suvrick/go-kiss-core/game"
	"github.com/suvrick/go-kiss-core/packets"
)

//103786258
//sessionKey=5d09db98a83f25ff3885114f725c651022ee76138454ff
//dc93c8e0c365ca792cf1198ab71c73e7

var urls = []string{
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=83908794&viewerId=83908794&isAppUser=1&isAppWidgetUser=0&sessionKey=5d28ffe5b89a5030be1188c910d74503ebaa1f2a84e50c&authKey=8ad52d5e75a3a16bff2a7b839d977d6c&apiSettings=227&silentBilling=1&lang=ru&fromServiceBlock=0&ls=0&pos=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1571560169&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"https://m.inspin.me/build/v431/?type=vk&user_id=205226070&api_url=https%3A%2F%2Fapi.vk.com%2Fapi.php&api_id=1930071&api_settings=8207&viewer_id=205226070&viewer_type=2&access_token=e5d59fb7b7a476dc0b129724cc2040f330d6a0014294442f2804acbac9437124f490ed9de1fd76ea25dad&is_app_user=1&auth_key=437455c7c77d0506af6ac9c48319d46f&language=0&parent_language=0&is_secure=1&sid=c954a62fd780f73bc7eb67d8e6d00bf3bd849226fa1bbac1ae155d07a150cc801bba5c20b95c8d7e2efa3&secret=c06f122141&stats_hash=882761eec0e774f10e&lc_name=eddb5306&api_script=https%3A%2F%2Fapi.vk.com%2Fapi.php&referrer=unknown&ads_app_id=1930071_d0d503a623afabab5b&platform=html5_android&hash=",
	"https://bottle2.itsrealgames.com/www/ok.html?api=ok&5&container=true&web_server=https%3A%2F%2Fok.ru&first_start=0&logged_user_id=917434085754&sig=76f0970b24b99c27f84147fbe4148f6a&refplace=vitrine_user_apps_portlet&new_sig=1&apiconnection=83735040_1563210992233&authorized=1&session_key=-s-4d2fRuOe622bNtSeb3e2z.q142c2NwRzc265xyP2-466Tvs6b343LURe75aeQvw20755MxRda2fbOvP6z2d3SVS35a4-Nvs60a37Pxa&clientLog=0&session_secret_key=01889d7e1c81173aeb74b34580a81c60&auth_sig=3324c256fc7578c5092335f7dcfa7008&api_server=https%3A%2F%2Fapi.ok.ru%2F&ip_geo_location=RU%2C90%2CPerm&application_key=CBADLOPFABABABABA",
	"https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593",
	"https://m.inspin.me/build/v431/?type=vk&user_id=292003911&api_url=https%3A%2F%2Fapi.vk.com%2Fapi.php&api_id=1930071&api_settings=8207&viewer_id=292003911&viewer_type=2&access_token=a0ce925b6322055cd7c291e7577bb363fb21ddd1c1026076d2ae71d1dd7e0e1416b68617869e6d20d6078&is_app_user=1&auth_key=2ff87aebac3ec78d0dc0fa5c55efda33&language=0&parent_language=0&is_secure=1&sid=e2048d62a447474d27fa6c5b862035e9d87cce7c8aba0affd06f06353c91280e416f39adf2f5d62abf77c&secret=46f45eb797&stats_hash=f1304753fffaf8bec8&lc_name=9791cbb4&api_script=https%3A%2F%2Fapi.vk.com%2Fapi.php&referrer=unknown&ads_app_id=1930071_7f55035857df794ec1&platform=html5_android&hash=",
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100046693&viewerId=100046693&isAppUser=1&isAppWidgetUser=0&sessionKey=5d121ddedbef9721fc0fc02d33a2011a6938773f38a853&authKey=dd52b12107363624100e77b8b5160b02&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=0&ls=0&pos=2&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563080077&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100088538&viewerId=100088538&isAppUser=1&isAppWidgetUser=0&sessionKey=5d5a438a024349b54f24de4e2900ed26a89089f36d4edd&authKey=0f649b5a99bcd94ee913839afe100e75&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=1&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563888983&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100114300&viewerId=100114300&isAppUser=1&isAppWidgetUser=0&sessionKey=5d1153d086c9b167f5aa239744f92226d35da8283b02c7&authKey=cdb75b921e5797b6a34d21dc1188003c&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=1&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563054359&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100215745&viewerId=100215745&isAppUser=1&isAppWidgetUser=0&sessionKey=5de7c2e21b149f78d195e60e048e42b25c4906f346a238&authKey=0ed6709b750a47e518e46e0e16e5c265&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=10&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1570857178&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100221042&viewerId=100221042&isAppUser=1&isAppWidgetUser=0&sessionKey=5d56e7867437f91f4dcabba819421670e3c3575d7c5440&authKey=c5f9af3ccf00d7922067e0e3a398e56d&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=0&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563079622&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786125&viewerId=103786125&isAppUser=1&isAppWidgetUser=0&sessionKey=5da6d124c7f823d753f582f7a5b39fbd85919e9d2b5cb6&authKey=65ea9ab0e227437a66d73a010bff42c2&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569557906&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786202&viewerId=103786202&isAppUser=1&isAppWidgetUser=0&sessionKey=5d063a76c0ea63fd48b4d06f8cc85e45122d9abe8bb4f0&authKey=b508c5dbc13ee88cf28affbb76d5fe34&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569558150&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786258&viewerId=103786258&isAppUser=1&isAppWidgetUser=0&sessionKey=5d09db98a83f25ff3885114f725c651022ee76138454ff&authKey=dc93c8e0c365ca792cf1198ab71c73e7&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569558375&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
}

func main() {
	CreateGame()
}

func CreateGame() {
	data := frame.Parse(urls[0])

	if strError, ok := data["error"]; ok {
		fmt.Println(strError)
		return
	}

	g := game.NewGame()
	err := g.Connect(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// LOGIN
	g.AddListen(4, func(self *game.Game, packet map[string]interface{}) {

		status, ok := packets.GetByte("status", packet)
		if !ok {
			fmt.Println("Fail get fiels \"status\"")
			self.Close()
		}

		game_id, ok := packets.GetInt("game_id", packet)
		if !ok {
			fmt.Println("Fail get fiels \"game_id\"")
			self.Close()
		}

		self.GameID = game_id

		self.Send(61, nil)

		self.Send(26, map[string]interface{}{
			"type": 0,
		})

		fmt.Printf("Get auth result: %d\n", status)
	})

	// REWARDS
	g.AddListen(13, func(self *game.Game, packet map[string]interface{}) {

		rewards, ok := packets.GetMapArray("rewards", packet)
		if !ok {
			fmt.Println("Fail get fields \"rewards\"")
			self.Close()
		}

		for _, v := range rewards {

			reward_id, _ := packets.GetInt("id", v)

			count, _ := packets.GetInt("count", v)

			self.Send(11, map[string]interface{}{
				"reward_id": reward_id,
				"count":     count,
			})
		}

		fmt.Printf("Get auth result: %v\n", rewards)
	})

	// BOTTLE_ROOM
	g.AddListen(25, func(self *game.Game, packet map[string]interface{}) {

		players, ok := packets.GetMapArray("players", packet)
		if !ok {
			fmt.Println("Fail get fields \"players\"")
			self.Close()
		}

		// REQUEST
		// ERROR!!!
		// self.Send(8, map[string]interface{}{
		// 	"players": players,
		// 	"mask1":   math.MaxInt32,
		// })
	})

	// BOTTLE_LEADER
	g.AddListen(28, func(self *game.Game, packet map[string]interface{}) {

		leader_id, ok := packets.GetInt("leader_id", packet)

		if !ok {
			fmt.Println("Fail get fields \"leader_id\"")
			return
		}

		if leader_id == self.GameID {
			go func() {
				fmt.Printf("I am leader!\n")
				<-time.After(time.Microsecond * 3)
				self.Send(28, map[string]interface{}{
					"speed": 0,
				})
			}()
		}

	})

	// BOTTLE_ROLL
	g.AddListen(29, func(self *game.Game, packet map[string]interface{}) {

		leader_id, ok := packets.GetInt("leader_id", packet)

		if !ok {
			fmt.Println("Fail get fields \"leader_id\"")
			return
		}

		rolled_id, ok := packets.GetInt("rolled_id", packet)

		if !ok {
			fmt.Println("Fail get fields \"rolled_id\"")
			return
		}

		msg := ""

		if leader_id == self.GameID {
			msg = fmt.Sprintln("I am has kissed as leader")
		}

		if rolled_id == self.GameID {
			msg = fmt.Sprintln("I am has kissed as rolled")
		}

		if len(msg) > 0 {
			go func() {
				fmt.Println(msg)
				<-time.After(time.Microsecond * 3)
				self.Send(29, map[string]interface{}{
					"answer": 1,
				})
			}()
		}
	})

	// LOGIN
	g.Send(4, data)

	// g.Send(202, []interface{}{types.I(5)})

	<-g.End()
}

// ['113594657', 32, 6,
// '7a2b140e7b42935768c040a54ade4cfc',
// 0,
// 'c5fdbab656b5cd5551d4f3e13ee20e41',
// 0, 0, 0, '0', ”]
func GetLoginPacket() map[string]interface{} {
	return map[string]interface{}{
		"login_id":   105345504,
		"frame_type": 30,
		"device":     5,
		"auth_key":   "7b0a077a088b9e5169bcfc0bf2ee9ae8",
	}
}
