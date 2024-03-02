package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/suvrick/go-kiss-core/frame"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/socket"
	"github.com/suvrick/go-kiss-core/types"
)

// 103786258
// sessionKey=5d09db98a83f25ff3885114f725c651022ee76138454ff
// dc93c8e0c365ca792cf1198ab71c73e7
const Tototo93 uint64 = 22132982

var urls = []string{
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

var urls2 = []string{
	//"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105345504&viewerId=105345504&isAppUser=1&isAppWidgetUser=0&sessionKey=5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38&authKey=7b0a077a088b9e5169bcfc0bf2ee9ae8&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573540656&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105345662&viewerId=105345662&isAppUser=1&isAppWidgetUser=0&sessionKey=5d896fb2ee4b69d7e910436a2f14b1ae33fd14ca433b78&authKey=5f93655e93d0a3ab6195dbf0656bab60&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573541066&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342583&viewerId=105342583&isAppUser=1&isAppWidgetUser=0&sessionKey=5df0a1002a5468b420eba13cddc23a79d6f994f83a92c9&authKey=fd02dd9d0285983e8a65f07a729a1193&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573532617&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342679&viewerId=105342679&isAppUser=1&isAppWidgetUser=0&sessionKey=5df2dbb3a1f8a8a701168d94755bcd98def3344a38ef49&authKey=a0dd2586c35a53319b3c53e04f2fd94d&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573532922&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342780&viewerId=105342780&isAppUser=1&isAppWidgetUser=0&sessionKey=5ded15002b6e8b4b03bb6be83436f838ed37b7c9272338&authKey=550ea3065c2c38c25e2bf8f50c893f1b&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573533240&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342864&viewerId=105342864&isAppUser=1&isAppWidgetUser=0&sessionKey=5d927abcb6d37d6a6aebe74e19e0a1b4e463c067584dd6&authKey=f79bd404e48fbd57a5a71d02b3ea9bb8&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573533546&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
}

var proxies = []string{
	"brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-178.171.38.41:hcx7fnqnph27",
}

var frameManager frame.IFrameManager

var games []*socket.Socket

var ids = []uint64{
	43591658, // Ж
	43591678, // Ж
	43591362, // Ж
	43591376, // Ж
	43591396, // Ж
	43591386, // Ж
}

var wg sync.WaitGroup

func main() {

	games = make([]*socket.Socket, 0)

	wg = sync.WaitGroup{}

	for _, url := range urls2 {

		wg.Add(1)

		login := getLoginPacket(url)

		config := socket.GetDefaultSocketConfig()

		homeDir, err := os.UserHomeDir()
		if err == nil {
			logDir := path.Join(homeDir, "logs")
			err := os.MkdirAll(logDir, 0700)
			if err == nil {
				fileName := fmt.Sprintf("%s%v.log", frame.GetFrameTypeName(login.NetType), login.LoginID)
				logPath := path.Join(logDir, fileName)
				file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0700)
				if err == nil {
					defer file.Close()
					config.Logger = log.New(file, "", log.Ltime|log.Lshortfile)
				}
			}
		}

		//config.TimeInTheGame = 5

		g := socket.NewSocket(config)
		g.SetOpenHandler(openHandle)
		g.SetCloseHandler(closeHandle)
		g.SetErrorHandler(errorHandle)
		g.SetRecvHandler(recvHandler)

		games = append(games, g)

		if err := g.Connection(); err != nil {
			// Error: Auth Failed (code: ip_forbidden)
			log.Fatalln(err.Error())
		}

		g.Send(client.LOGIN, login)
	}

	wg.Wait()
}

func openHandle(sender *socket.Socket) {
	sender.Log("Open connection")
}

func closeHandle(sender *socket.Socket, rule byte, msg string) {
	sender.Logf("Close connection. Rule: %v, %s\n", rule, msg)
	wg.Done()
}

func errorHandle(sender *socket.Socket, err error) {
	if err != nil {
		sender.Logf("Error: %v\n", err.Error())
		sender.Close()
	}
}

func recvHandler(sender *socket.Socket, packetID types.PacketServerType, packet any) {
	switch packetID {
	case server.LOGIN:
		login, ok := packet.(*server.Login)
		if ok {
			switch login.Result {
			case 0:
				sender.HiroID = login.HiroID
				sender.Send(client.TOGGLE_TAPE, &client.ToggleType{})
				sender.Send(client.BOTTLE_PLAY, &client.BottlePlay{
					RoomID: 0,
				})
				sender.Send(client.MOVE, &client.Move{
					PlayerID: Tototo93,
				})

				go func() {
					<-time.After(time.Second * 5)
					photoLike := client.PhotoLike{
						PlayerID: Tototo93,
						PhotoID:  1,
						IsLike:   1,
					}
					sender.Send(client.PHOTOS_LIKE, &photoLike)
				}()

				// go func() {
				// 	<-time.After(time.Second * 5)
				// 	buy := client.Buy{
				// 		BuyType:    2,
				// 		Coin:       260,
				// 		PlayerID:   Tototo93,
				// 		PrizeID:    10321,
				// 		ByteFiald:  0,
				// 		Count:      1,
				// 		ByteFiald2: 6,
				// 	}
				// 	sender.Send(client.BUY, &buy)
				// }()
			default:
				sender.Close()
			}
		}
	case server.BONUS:
		bonus, ok := packet.(*server.Bonus)
		if ok {
			if bonus.Day > 0 {
				sender.Send(client.BONUS, &client.Bonus{})
			}
		}
	case server.BOTTLE_LEADER:
		bottleLeader, ok := packet.(*server.BottleLeader)
		if ok {
			if bottleLeader.LeaderID == sender.HiroID {
				sender.Logf("[use] %s -> i am roll bottle %d", bottleLeader, bottleLeader.LeaderID)
				go func() {
					<-time.After(time.Second * 7)
					sender.Send(client.BOTTLE_ROLL, &client.BottleRoll{})
				}()
			}
		}
	case server.BOTTLE_ROLL:
		bottleRoll, ok := packet.(*server.BottleRoll)
		if ok {

			needSendKiss := false

			if bottleRoll.LeaderID == sender.HiroID {
				sender.Logf("[use] %s -> i am kiss as leader %d", bottleRoll, bottleRoll.LeaderID)
				needSendKiss = true
			} else if bottleRoll.RollerID == sender.HiroID {
				sender.Logf("[use] %s -> i am kiss as roller %d", bottleRoll, bottleRoll.RollerID)
				needSendKiss = true
			}

			if needSendKiss {
				go func() {
					<-time.After(time.Second * 6)
					sender.Send(client.BOTTLE_KISS, &client.BottleKiss{
						Answer: 1,
					})
				}()
			}
		}
	case server.REWARDS:
		rewards, ok := packet.(*server.Rewards)
		if ok {
			for _, reward := range rewards.Items {
				sender.Send(client.GAME_REWARDS_GET, &client.GameRewardsGet{
					RewardID: reward.RewardID,
				})
				break
			}
		}
	case server.COLLECTIONS_POINTS:
		points, ok := packet.(*server.CollectionsPoints)
		if ok {
			sender.Logf("[use] %s -> collections points %d", points, points.Points)
		}
	case server.BOTTLE_JOIN:
		join, ok := packet.(*server.BottleJoin)
		if ok {
			sender.Send(client.REQUEST, &client.Request{
				Players: []uint64{join.PlayerID},
				Mask:    262143,
			})
		}
	}

	if pool, ok := packet.(server.IPooling); ok {
		pool.Reset()
	}
}

func getLoginPacket(url string) *client.Login {

	if frameManager == nil {
		frameManager = frame.New()

	}

	frameDTO, err := frameManager.Parse(url)
	if err != nil {
		return nil
	}

	auth := byte(0)
	if len(frameDTO.AccessToken) > 0 {
		auth = byte(1)
	}

	return &client.Login{
		LoginID:     frameDTO.ID,
		NetType:     frameDTO.NetType,
		DeviceType:  5,
		Key:         frameDTO.Key,
		OAuth:       auth,
		AccessToken: frameDTO.AccessToken,
		StringField: frameDTO.StringField,
		Captcha:     captcha,
	}
}

const captcha = "BM-2763"
