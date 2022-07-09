package main

import (
	"log"
	"net/http"
	"text/template"
)

func settingHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/setting" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/setting.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/setting",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func packetsHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/packets" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/packets.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/packets",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/home.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

type Menu struct {
	IsMenu   bool
	PathName string
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/packets", packetsHandler)

	mux.HandleFunc("/setting", settingHandler)

	fileServer := http.FileServer(http.Dir("./frontend/"))

	mux.Handle("/frontend/", http.StripPrefix("/frontend", fileServer))

	http.ListenAndServe(":8080", mux)

	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000995&sessionKey=fa085f019b006228b268d70717a6bef&authKey=596cef1cbd19f0eb1377e6c26677cb9b&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000994&sessionKey=6d99d2f4c1b0afe892081dbd8ca29a9&authKey=affb86fd9cfe75896ab66d76f15105e6&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000993&sessionKey=b40ae09ab601efaa4e5e35aaaa8d29c&authKey=cd302d09033bdd4bccd9c1cb149a4b16&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000991&sessionKey=ccb0092af206af82d0e104b89a35a6a&authKey=57c4f876706cf5b0d72d4658c5d37a5c&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000990&sessionKey=41452e5bf03e7549e97da02dc523e53&authKey=b5f90a924aef38115828cdeec51a66fc&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1571552893425&&userId=114941701&sessionKey=67f5e4f7a90144c5eba1b91694132904&authKey=33513e2ce85cabfd6ec59d827aa28cea&net_type=32&useApiType=sa&email=tsd_8326@mail.ru&locale=RU"
	//url := "https://m.inspin.me/build/v431/?api=ok&5&sig=13c32cf30efbedcbfd7afa689f9fa5b6&mob=true&refplace=user_apps_portlet&session_key=-s-6622x0r.b2d-yvr6022bQSS26241xUv3433cvvq4c232NwOZ376bTuxd62-4TxN41178Msr1a63fMXr45827NWo277c0TUq8b265Ne&session_secret_key=7f16a539d8d56e2fca91be30ac68c7e3&auth_sig=d04204aa3632ab2dcf9e6f2b41084b57&api_server=https%3A%2F%2Fapi.ok.ru%2F&payment_server=https%3A%2F%2Fmpay.ok.ru&lang=ru&application_key=CBADLOPFABABABABA&mob_platform=androidweb&logged_user_id=952642558996&type=ok"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000937&sessionKey=6c92d6b72649d11dc6a436777c8f3fd&authKey=3109ee5271d550293cf9b2fb494556f1&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=95ac81be3e64f73fd56e4b12e004c9d0&vid=9914552860808032409&oid=9914552860808032409&app_id=543574&authentication_key=27861583a5de8bd7c1aa5955dfe757b8&session_expire=1623244207&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=2add120d8f235755ee8b34801798f775&window_id=CometName_659cc22f549d0d25dc769cb39daa105f&referer_type=left_menu&version=1593"
	//url := "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593"
	//url := "https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=7556246c19559181c1bcd32682446ea2&vid=6611932282865570796&oid=6611932282865570796&app_id=543574&authentication_key=5a4b728af3bb1c239e6a0c7978b77361&session_expire=1530439148&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=dbece9235b7b4f42d80aaa7054e263c0&window_id=CometName_d8d8afe8a45d0f3fc9218526c97b438f&referer_type=left_menu#api=mm&packageName=bottlePackage&config=https://bottle2.itsrealgames.com/www/config_release.xml&protocol=https:&locale=RU&locale_url=https://bottle2.itsrealgames.com/resources/locale/EN_All_14.lp&width=1000&height=690&useApiType=mm&"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1563646914252&&userId=117513726&sessionKey=6d7119c40fea856164b86cea934048e5&authKey=0b0cfe8b4d5bf6ed0651bf8acff50a42&net_type=32&useApiType=sa&email=sancha5259@yandex.ru&locale=RU"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1563443602913&&userId=120302710&sessionKey=3934f9b4f51f53f85a40e1be8c1d27fd&authKey=54a9a75685cbd5074106ee5d85921e23&net_type=32&useApiType=sa&email=honey.alekse@bk.ru&locale=RU"
	//url := "https://bottle2.itsrealgames.com/www/ok.html?api=ok&5&container=true&web_server=https%3A%2F%2Fok.ru&first_start=0&logged_user_id=951449229841&sig=b43e71958dc877146f26aec36e98b424&refplace=vitrine_app_search_apps&new_sig=1&apiconnection=83735040_1558615128953&authorized=1&session_key=-s-7b2dNSqb5374zvq.-0b2.ut802ebO0P0.a72wvN67bf3y1N161a4PVr3Z662K.w11b66SUPeb339xuP30073M-v20b20NtO11f17PWO7&clientLog=0&session_secret_key=a510c8475cb8962fe148cb9ffac478b2&auth_sig=30b8f8f633bbb4630cb6960e826784a2&api_server=https%3A%2F%2Fapi.ok.ru%2F&ip_geo_location=RU%2C04%2CBarnaul&application_key=CBADLOPFABABABABA#api=ok&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?96&width=1000&height=690&sprites_version=53&useApiType=ok&"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000100&sessionKey=952572205dbf4c77c77144718b68c78&authKey=523f781485b25b8a663d3fb693a1e2dd&net_type=32&useApiType=sa"
	//url := "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=104015872&viewerId=104015872&isAppUser=1&isAppWidgetUser=0&sessionKey=5db68d0bc9be09339317d785fd88ec42ae7976d721bf8e&authKey=3d6ac3df40e40ffc4243b650efae46f5&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1570189957&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?161&width=1000&height=690&sprites_version=87&useApiType=fs&"
	//url := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&type=vk&record_first_session=1&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=2105359&viewer_id=565065176&viewer_type=2&sid=787a953c613f2ad17c1549da3cab52d2333bb50fe4196e3472b29cede2ba324d1ddd0f5c8141a4a349aea&secret=faf91f3930&access_token=c5b2de1bccd3965d89f70436e53e1fca0cfda3c9934f271c6f06d04894d116dc2134f128fcc6a2c425ac7&user_id=565065176&group_id=0&is_app_user=1&auth_key=5d89c796f75abda9172c6525f08fbc6f&language=0&parent_language=0&is_secure=1&stats_hash=6e92d3737ecd984f7d&ads_app_id=1930071_20986469aa936156b7&referrer=apps_installed_games&lc_name=f0aef739&platform=web&hash="
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000263&sessionKey=0767f51d8f7672849206fe612d74dcb&authKey=c98103fcc0e4d859026f28ea5eeb7ce9&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=106730666&viewerId=106730666&isAppUser=1&isAppWidgetUser=0&sessionKey=5dd6c1cccf4e6b36384a144bf4ad69d3951f43173cd33c&authKey=edc6ca6a0644ef7651a2fe28e2c53f51&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1575621360&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?193&width=1000&height=690&sprites_version=103&useApiType=fs&"
	//url := "https://bottle2.itsrealgames.com/www/ok.html?api=ok&5&container=true&web_server=https%3A%2F%2Fok.ru&first_start=0&logged_user_id=952315233325&sig=2fd168e6af2b022158d2ee2e83fd6b06&refplace=vitrine_user_apps_portlet&new_sig=1&apiconnection=83735040_1575626942063&authorized=1&session_key=-s-7f2aNSq85364Luq3-f62Q1sd02bcOyP06f62NvS80b43NyN168a4OVq7b682K.r4146aSuP0b33cw0P70475y-P2092eNtP11554Pxw7&clientLog=0&session_secret_key=0d1c183a1d6593ec1fec6dc49f0a9b80&auth_sig=f6e4cdb1535d23b50072ace0cd2c6ed5&api_server=https%3A%2F%2Fapi.ok.ru%2F&ip_geo_location=RU%2C05%2CBlagoveshchensk&application_key=CBADLOPFABABABABA#api=ok&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?193&width=1000&height=690&sprites_version=103&useApiType=ok&                        "
	//url := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&type=vk&record_first_session=1&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=2105359&viewer_id=95767995&viewer_type=2&sid=5d60ef9000861114a61983994d580fc0cbd1d09b5cfafaf668d6521443bcbf7cc8979500f911000d0a2a6&secret=322b91c396&access_token=161f874f04e39345b9ba3156fc0ad9b0f76fbcd273510636c9aec6016792d58c651a187e39f33700a6f3a&user_id=95767995&group_id=0&is_app_user=1&auth_key=3899e4f723110882fcf562a4364109bc&language=0&parent_language=0&is_secure=1&stats_hash=202c2e0e89e7622264&ads_app_id=1930071_f6ea1f58620f95d827&referrer=apps_installed_games&lc_name=6dae92ed&platform=web&hash="
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1570856185867&&userId=120319455&sessionKey=d3c1478725153626c21e55ecf3a77418&authKey=b4da2a18078edbbdbd0b74357620c416&net_type=32&useApiType=sa&email=zoro-7z@mail.ru&locale=RU"
	//url := "https://m.inspin.me/build/v1854/?type=vk&user_id=115111645&api_url=https%3A%2F%2Fapi.vk.com%2Fapi.php&api_id=1930071&api_settings=8207&viewer_id=115111645&viewer_type=2&access_token=6cce3eef8753639b36b3a37b3ee61982c10795f77b51f5de67201a8f7547fe02c948e5324034d1eaa1453&is_app_user=1&auth_key=c411ffa49f3cd61fa4c9f14f1e6e2443&language=0&parent_language=0&is_secure=1&sid=c4f3cbf8a2af443d9ca85d8ee23850359c0bdaaa9f8c5321fcf9e179659d02187f7b98e74c4af7b11a42b&secret=10062827e3&stats_hash=951e51cb625a4394e1&lc_name=b758426b&api_script=https%3A%2F%2Fapi.vk.com%2Fapi.php&referrer=unknown&ads_app_id=1930071_7d6938509e27c7a4ae&platform=html5_android&hash=&version=575"
	//url := "https://inspin.me/?type=vk&api=vk&net_type=10&record_first_session=1&api_url=https://api.vk.com/api.php&api_id=6926012&api_settings=9247&viewer_id=395993763&viewer_type=0&sid=76c752b06f04b4ed765b9d2efcd22c9431fed4da85bdf4e975f8b212dca6b1b97e4807469fee2abec841f&secret=192ae157c4&access_token=48b5be06f848847b33d7ec5902d991997abaf687c8aa806054387e35a76e8fc178d5b674fe88fc7f4d280&user_id=0&is_app_user=1&language=0&parent_language=0&is_secure=1&stats_hash=7e7beb6e27aa4afe67&is_favorite=0&group_id=0&ads_app_id=6926012_6d11e30be48738e51d&access_token_settings=notify,friends,photos,audio,video,status,wall&referrer=apps_installed_games&lc_name=d87f4daf&platform=web&is_widescreen=0&whitelist_scopes=friends,photos,video,stories,pages,status,notes,wall,docs,groups,stats,market,ads,notifications&group_whitelist_scopes=stories,photos,app_widget,messages,wall,docs,manage&auth_key=86af7f8b3a199db98dec53d56238519e&timestamp=1641735341&sign=imYaNCj4FPFHY8sU0OPeWzvDny6jtuiIMpFXhl-ONyQ&sign_keys=access_token,access_token_settings,ads_app_id,api_id,api_settings,api_url,auth_key,group_id,group_whitelist_scopes,is_app_user,is_favorite,is_secure,is_widescreen,language,lc_name,parent_language,platform,referrer,secret,sid,stats_hash,timestamp,user_id,viewer_id,viewer_type,whitelist_scopes&hash="
	//url := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=8206&viewer_id=101595089&viewer_type=2&sid=e4d3a95c860dc1fa7b89a51b7f41cb42a476f87dfabdda68e9a930f573bab03962246910e605846f6326d&secret=b54e5afef8&access_token=3002be5efa8f0426fcf62b9b9dcba0db0cb1ed593ae1b6959dbcc3e66e5a37e261586b0c6dc08c072d58f&user_id=101595089&group_id=0&is_app_user=1&auth_key=e4f2fd67d7a93326b0a7d1d066150aaf&language=0&parent_language=0&is_secure=1&stats_hash=4397c8480f5251f767&ads_app_id=1930071_527d26e6597a93541d&referrer=unknown&lc_name=875f1eaf&platform=web&hash="
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000010&sessionKey=0d553183c67145f2b218e41f1c2c4ef&authKey=f93c2c8a65649db70e66c2560409022c&net_type=32&useApiType=sa"
	//url := "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105006740&viewerId=105006740&isAppUser=1&isAppWidgetUser=0&sessionKey=5dc02d6697b7bc435d9ebd74f9f1878ee5443d547e0afd&authKey=92112d12e2aea25317c527fabb027cf3&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1572743067&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?172&width=1000&height=690&sprites_version=93&useApiType=fs&"
	//url := "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105006807&viewerId=105006807&isAppUser=1&isAppWidgetUser=0&sessionKey=5d29e4421477a7a3c7c00612962331d51d07528297cd58&authKey=2391e221fce44c0b2625e97f77d3331f&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1572743450&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?172&width=1000&height=690&sprites_version=93&useApiType=fs&"
	//url := "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105189479&viewerId=105189479&isAppUser=1&isAppWidgetUser=0&sessionKey=5d49a3ce1b91be4574ec0fe85f597427f651c36f8cbf3f&authKey=12982ebe43036d53a6ea3f35595d54ac&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573199089&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?178&width=1000&height=690&sprites_version=95&useApiType=fs&"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000996&sessionKey=9b48038c96a1fb08f93a058214ad2a1&authKey=c0dc84721d911684552bfa9e5b9b095c&net_type=32&useApiType=sa"
	//url := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000997&sessionKey=a56e72219f9d53fe72c733e4a9e40c9&authKey=aadce5a5f4604f10701c0b8c2871881b&net_type=32&useApiType=sa"
	//url := "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105006496&viewerId=105006496&isAppUser=1&isAppWidgetUser=0&sessionKey=5df4e1cdfb85ce645125974f6377fe68122d2c574ac047&authKey=3526dceb8566fc68733ce0e2ad51704f&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1572741514&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?172&width=1000&height=690&sprites_version=93&useApiType=fs&"
	//url := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=8207&viewer_id=548071843&viewer_type=2&sid=babc3204328374be5db0de84de1d9370d36d4c3eeb00ceb6fcfa4f0f7814f542feaeadba82907631f9678&secret=47698d7a00&access_token=e6ad22113d51534b7a3ac08ca5a6a4c7a2e7c93d05897fe79dd77553baf753231d7031d2c6bdd04fd2782&user_id=548071843&group_id=0&is_app_user=1&auth_key=548e16b40805ec1ddaedaf54cf3ae807&language=0&parent_language=0&is_secure=1&stats_hash=17342c2962c0b6892c&ads_app_id=1930071_e275bbf2fa6b9dded2&referrer=apps_installed_games&lc_name=3a18b421&platform=web&hash="
	// http.HandleFunc("/", serveFiles)
	// http.HandleFunc("/ws", acceptWS)
	// http.ListenAndServe(":8080", nil)
}

// [113594657, 32, 4, "7a2b140e7b42935768c040a54ade4cfc", 0, "8c9991f3e49ef7d20d33432d1534e378"]

// [104015872, 30, 4, "3d6ac3df40e40ffc4243b650efae46f5"
