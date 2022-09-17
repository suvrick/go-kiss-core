package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/suvrick/go-kiss-core/bot"
	"github.com/suvrick/go-kiss-core/frame"
	"github.com/suvrick/go-kiss-core/game"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/proxy"
)

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

var frameParser *frame.Frame

var writeLog = flag.Bool("l", false, "write result to 'log/bot_id'")
var url = flag.String("f", "", "frame by game")
var path = flag.String("p", "", "path from file frames")
var count = flag.Int("b", 3, "count max instance game")

func main() {

	// web.Run()

	// return

	flag.Parse()

	if *writeLog {
		fmt.Printf("write log: on\n")
	} else {
		fmt.Printf("write log: off\n")
	}

	uu := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&type=vk&record_first_session=1&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=8207&viewer_id=95680242&viewer_type=2&sid=4921eb0cd7c7889784f9833c6c4577fd07b6ddc15aa862d3fd307d59fbb6900c1fcc3282353b5d8ef9179&secret=5d754c0d65&access_token=20ba1ae06cbab26b14af359334e71e7226e341c70c1f0f0995faac9cc0e6399e862850ba826d748257958&user_id=95680242&group_id=0&is_app_user=1&auth_key=d3936a1b653517c641e2e9fdae093d27&language=0&parent_language=0&is_secure=1&stats_hash=6e744671f2cea9e1bd&ads_app_id=1930071_6b0928ee3c2621bd14&referrer=unknown&lc_name=67c77d2f&platform=web&hash="

	url = &uu

	Run()
}

func Run() {

	frameParser = frame.NewFrameDefault()
	if frameParser.Err != nil {
		fmt.Println(frameParser.Err.Error())
		return
	}

	if *url != "" {
		login, err := frameParser.Parse2([]byte(*url))
		if err != nil {
			fmt.Println("parse frame: FAIL")
			return
		}

		fmt.Println("parse frame: OK")

		p := proxy.GetDefaultProxy()

		game := game.NewGameWithProxyDefault(p.URL)
		game.Run()
		game.LoginSend(login)

		// game := game.NewGame(context *Context)
		// game.SetConfig(config *socket.Config)
		// game.SetBot(bot *bot.Bot)
		// game.SetProxyManager(proxyManager *IProxyManager)
		// game.SetDoneFunc(doneFunc func(bot *bot.Bot))
		// game.Run()

		// <-game.Done

		// buy := &client.Buy{
		// 	BuyType:  2,
		// 	Coin:     30,
		// 	PlayerID: 42870078,
		// 	PrizeID:  10169,
		// 	XZ:       0,
		// 	Count:    1,
		// 	XZ2:      5,
		// }

		// game.SetBuyPacket(buy)

		bot := <-game.Done

		if *writeLog {
			Log(bot)
		} else {
			js, err := MarshalBot(bot)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println(string(js))
		}
		return
	}

	if *path != "" {
		Do()
	}
}

func Do() {

	frames := Load(*path)

	fmt.Printf("load frames (%d) for file '%s'\n", len(frames), *path)

	wg := sync.WaitGroup{}

	balanser := make(chan struct{}, *count)

	fmt.Printf("set game max instanse: %d\n", *count)

	for _, login := range frames {
		balanser <- struct{}{}
		wg.Add(1)

		go func(l client.Login) {
			defer func() {
				<-balanser
				wg.Done()
			}()
			game := game.NewGameDefault()
			game.Run()
			game.LoginSend(&l)

			// buy := &client.Buy{
			// 	BuyType:  251,
			// 	Coin:     10,
			// 	PlayerID: 42870078,
			// 	PrizeID:  10242,
			// 	XZ:       0,
			// 	Count:    1,
			// 	XZ2:      6,
			// }

			//game.SetBuyPacket(buy)

			bot := <-game.Done

			if *writeLog {
				Log(bot)
			} else {
				js, err := MarshalBot(bot)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(js))
			}
		}(login)
	}

	wg.Wait()

	close(balanser)
}

func Load(path string) []client.Login {
	result := make([]client.Login, 0)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	frames := bytes.Split(data, []byte{'\n'})

	for _, v := range frames {

		login, err := frameParser.Parse2(v)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		result = append(result, *login)
	}

	return result
}

func MarshalBot(bot bot.Bot) ([]byte, error) {
	js, err := json.MarshalIndent(&bot, " ", "  ")
	if err != nil {
		return nil, err
	}

	return js, nil
}

func Log(bot bot.Bot) {

	js, err := MarshalBot(bot)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dir := filepath.Join(cDir, "log")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	f_name := filepath.Join(dir, bot.ID)

	f, err := os.OpenFile(f_name, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Fprint(f, string(js))

	defer f.Close()
}
