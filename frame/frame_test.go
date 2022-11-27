package frame

import (
	"fmt"
	"testing"
)

var frames = []string{
	//"",
	//"invalid frame",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1571556651549&locale=RU&viewer_id=96172620&userId=96172620&access_token=2e39e1384066cc2c37665f53eff40409adf22061c804db05b80d05a1cbe2d22a41fe282bc264a14405b29&api_url=https://api.vk.com/method/&net_type=0&api_id=4432260&useApiType=vk&OAuth=true&auth_key=&email=pripevochka8989@mail.ru",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=104015872&viewerId=104015872&isAppUser=1&isAppWidgetUser=0&sessionKey=5db68d0bc9be09339317d785fd88ec42ae7976d721bf8e&authKey=3d6ac3df40e40ffc4243b650efae46f5&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1570189957&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?161&width=1000&height=690&sprites_version=87&useApiType=fs&",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1571552893425&&userId=114941701&sessionKey=67f5e4f7a90144c5eba1b91694132904&authKey=33513e2ce85cabfd6ec59d827aa28cea&net_type=32&useApiType=sa&email=tsd_8326@mail.ru&locale=RU",
	"https://m.inspin.me/build/v431/?api=ok&5&sig=13c32cf30efbedcbfd7afa689f9fa5b6&mob=true&refplace=user_apps_portlet&session_key=-s-6622x0r.b2d-yvr6022bQSS26241xUv3433cvvq4c232NwOZ376bTuxd62-4TxN41178Msr1a63fMXr45827NWo277c0TUq8b265Ne&session_secret_key=7f16a539d8d56e2fca91be30ac68c7e3&auth_sig=d04204aa3632ab2dcf9e6f2b41084b57&api_server=https%3A%2F%2Fapi.ok.ru%2F&payment_server=https%3A%2F%2Fmpay.ok.ru&lang=ru&application_key=CBADLOPFABABABABA&mob_platform=androidweb&logged_user_id=952642558996&type=ok",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000937&sessionKey=6c92d6b72649d11dc6a436777c8f3fd&authKey=3109ee5271d550293cf9b2fb494556f1&net_type=32&useApiType=sa",
	"https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=95ac81be3e64f73fd56e4b12e004c9d0&vid=9914552860808032409&oid=9914552860808032409&app_id=543574&authentication_key=27861583a5de8bd7c1aa5955dfe757b8&session_expire=1623244207&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=2add120d8f235755ee8b34801798f775&window_id=CometName_659cc22f549d0d25dc769cb39daa105f&referer_type=left_menu&version=1593",
	"https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593",
	"https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593",
}

func TestParseFrames(t *testing.T) {
	for _, v := range frames {
		r := Parse2(v)
		fmt.Printf("%#v\n", r)

		b, ok := r["frame_type"].(int)
		if !ok {
			fmt.Println(b)
		}

		fmt.Println(b)
	}
}
