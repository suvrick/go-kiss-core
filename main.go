package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/suvrick/go-kiss-core/frame"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/ws"
)

var urls = []string{
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

	f := frame.NewFrameDefault()

	id, r, _ := f.GetValue(urls[0])

	bb := make([]byte, 0)
	for _, v := range r {
		b := leb128.Write(v)
		bb = append(bb, b...)
	}

	bb = append(leb128.Write(4), bb...)
	bb = append(leb128.Write(0), bb...)
	l := leb128.Write(len(bb))
	bb = append(l, bb...)

	s := ws.NewSocket(ws.GetDefaultSocketConfig())
	s.SetErrorHandler(func(err error) {
		log.Println("socket error")
		log.Println(err.Error())
	})

	s.SetOpenHandler(func() {
		log.Println("socket open")
	})

	s.SetCloseHandler(func(rule byte, msg string) {
		log.Println("socket close")
	})

	s.SetReadHandler(func(reader io.Reader) {
		var b = make([]byte, 512)
		reader.Read(b)
		log.Println(b)
	})

	fmt.Printf("%d, %v, %v", id, r, bb)

	s.Go()
	s.Send(bb)

	<-time.After(1 * time.Minute)
}
