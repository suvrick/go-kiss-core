package main

import (
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
	// "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=f53f650cd57b6bc75da0b65af0d0c028&vid=13402962412699287699&oid=13402962412699287699&app_id=543574&authentication_key=e1de7d6b1b9a18e124331d1a8e7a6709&session_expire=1623248257&ext_perm=notifications%2Cemails%2Cpayments&sig=d38fca257b4651d5fc2bbc3e2531842f&window_id=CometName_74be9f9e99659ab7f65e85f2a31d3d3b&referer_type=left_menu&version=1593",
	// "https://m.inspin.me/build/v431/?type=vk&user_id=292003911&api_url=https%3A%2F%2Fapi.vk.com%2Fapi.php&api_id=1930071&api_settings=8207&viewer_id=292003911&viewer_type=2&access_token=a0ce925b6322055cd7c291e7577bb363fb21ddd1c1026076d2ae71d1dd7e0e1416b68617869e6d20d6078&is_app_user=1&auth_key=2ff87aebac3ec78d0dc0fa5c55efda33&language=0&parent_language=0&is_secure=1&sid=e2048d62a447474d27fa6c5b862035e9d87cce7c8aba0affd06f06353c91280e416f39adf2f5d62abf77c&secret=46f45eb797&stats_hash=f1304753fffaf8bec8&lc_name=9791cbb4&api_script=https%3A%2F%2Fapi.vk.com%2Fapi.php&referrer=unknown&ads_app_id=1930071_7f55035857df794ec1&platform=html5_android&hash=",
	// "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100046693&viewerId=100046693&isAppUser=1&isAppWidgetUser=0&sessionKey=5d121ddedbef9721fc0fc02d33a2011a6938773f38a853&authKey=dd52b12107363624100e77b8b5160b02&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=0&ls=0&pos=2&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563080077&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	// "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100088538&viewerId=100088538&isAppUser=1&isAppWidgetUser=0&sessionKey=5d5a438a024349b54f24de4e2900ed26a89089f36d4edd&authKey=0f649b5a99bcd94ee913839afe100e75&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=1&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563888983&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	// "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100114300&viewerId=100114300&isAppUser=1&isAppWidgetUser=0&sessionKey=5d1153d086c9b167f5aa239744f92226d35da8283b02c7&authKey=cdb75b921e5797b6a34d21dc1188003c&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=1&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563054359&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	// "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100215745&viewerId=100215745&isAppUser=1&isAppWidgetUser=0&sessionKey=5de7c2e21b149f78d195e60e048e42b25c4906f346a238&authKey=0ed6709b750a47e518e46e0e16e5c265&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=10&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1570857178&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	// "https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=100221042&viewerId=100221042&isAppUser=1&isAppWidgetUser=0&sessionKey=5d56e7867437f91f4dcabba819421670e3c3575d7c5440&authKey=c5f9af3ccf00d7922067e0e3a398e56d&apiSettings=743&silentBilling=1&lang=ru&fromServiceBlock=1&ls=0&pos=1&isFavlb=1&is_global=1&from_id=left.menu.service&from=left.menu.service&hasNotifications=1&_v=1&isOfferWallEnabled=0&appManage=0&connId=1563079622&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786125&viewerId=103786125&isAppUser=1&isAppWidgetUser=0&sessionKey=5da6d124c7f823d753f582f7a5b39fbd85919e9d2b5cb6&authKey=65ea9ab0e227437a66d73a010bff42c2&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569557906&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786202&viewerId=103786202&isAppUser=1&isAppWidgetUser=0&sessionKey=5d063a76c0ea63fd48b4d06f8cc85e45122d9abe8bb4f0&authKey=b508c5dbc13ee88cf28affbb76d5fe34&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569558150&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
	"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=103786258&viewerId=103786258&isAppUser=1&isAppWidgetUser=0&sessionKey=5d09db98a83f25ff3885114f725c651022ee76138454ff&authKey=dc93c8e0c365ca792cf1198ab71c73e7&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1569558375&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?158&width=1000&height=690&sprites_version=83&useApiType=fs&",
}

var urls2 = []string{
	//"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105345504&viewerId=105345504&isAppUser=1&isAppWidgetUser=0&sessionKey=5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38&authKey=7b0a077a088b9e5169bcfc0bf2ee9ae8&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573540656&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	//"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105345662&viewerId=105345662&isAppUser=1&isAppWidgetUser=0&sessionKey=5d896fb2ee4b69d7e910436a2f14b1ae33fd14ca433b78&authKey=5f93655e93d0a3ab6195dbf0656bab60&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573541066&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	//"view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342583&viewerId=105342583&isAppUser=1&isAppWidgetUser=0&sessionKey=5df0a1002a5468b420eba13cddc23a79d6f994f83a92c9&authKey=fd02dd9d0285983e8a65f07a729a1193&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573532617&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342679&viewerId=105342679&isAppUser=1&isAppWidgetUser=0&sessionKey=5df2dbb3a1f8a8a701168d94755bcd98def3344a38ef49&authKey=a0dd2586c35a53319b3c53e04f2fd94d&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573532922&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342780&viewerId=105342780&isAppUser=1&isAppWidgetUser=0&sessionKey=5ded15002b6e8b4b03bb6be83436f838ed37b7c9272338&authKey=550ea3065c2c38c25e2bf8f50c893f1b&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573533240&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
	// "view-source:https://bottle2.itsrealgames.com/www/fs.html?5&apiUrl=https%3A%2F%2Fapi.fotostrana.ru%2Fapifs.php&apiId=bottle&userId=105342864&viewerId=105342864&isAppUser=1&isAppWidgetUser=0&sessionKey=5d927abcb6d37d6a6aebe74e19e0a1b4e463c067584dd6&authKey=f79bd404e48fbd57a5a71d02b3ea9bb8&apiSettings=743&silentBilling=1&lang=ru&forceInstall=1&from=app.popup&from_id=app.popup&hasNotifications=0&_v=1&isOfferWallEnabled=0&appManage=0&connId=1573533546&ourIp=0&lc_name=&fs_api=https://st.fotocdn.net/swf/api/__v1344942768.fs_api.swf&log=0&swfobject=https://st.fotocdn.net/js/__v1368780425_1_4.swfobject2.js&fsapi=https://st.fotocdn.net/app/app/js/__v1540476017_1_4.fsapi.js&xdm_e=https://fotostrana.ru&xdm_c=default0&xdm_p=1#api=fs&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?182&width=1000&height=690&sprites_version=96&useApiType=fs&",
}

var urls3 = []string{
	"https://bottle2.itsrealgames.com/www/sta.html?time=1559375815493&&userId=116590694&sessionKey=a070d802fba9142968b956b9036aaa53&authKey=37b02e73975ffb44f2c60b67e85ec22f&net_type=32&useApiType=sa&email=nemiro-f@ya.ru&locale=RU#time=1559375815493&userId=116590694&sessionKey=a070d802fba9142968b956b9036aaa53&authKey=37b02e73975ffb44f2c60b67e85ec22f&net_type=32&useApiType=sa&email=nemiro-f@ya.ru&locale=RU&api=sa&packageName=bottlePackage&config=config_release.xml&protocol=https:&international=false&locale_url=../resources/locale/EN_All.lp?100&width=1000&height=690&sprites_version=54&",
	"https://bottle2.itsrealgames.com/www/ok.html?api=ok&5&container=true&web_server=https%3A%2F%2Fok.ru&first_start=0&logged_user_id=951449229841&sig=b43e71958dc877146f26aec36e98b424&refplace=vitrine_app_search_apps&new_sig=1&apiconnection=83735040_1558615128953&authorized=1&session_key=-s-7b2dNSqb5374zvq.-0b2.ut802ebO0P0.a72wvN67bf3y1N161a4PVr3Z662K.w11b66SUPeb339xuP30073M-v20b20NtO11f17PWO7&clientLog=0&session_secret_key=a510c8475cb8962fe148cb9ffac478b2&auth_sig=30b8f8f633bbb4630cb6960e826784a2&api_server=https%3A%2F%2Fapi.ok.ru%2F&ip_geo_location=RU%2C04%2CBarnaul&application_key=CBADLOPFABABABABA#api=ok&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?96&width=1000&height=690&sprites_version=53&useApiType=ok&",
	"https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=7c42054278baa4e7096e6748757c8954&vid=9156245585513337776&oid=9156245585513337776&app_id=543574&authentication_key=03e68868caf7d1b23a9a77ac096243d0&session_expire=1539077837&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=9fa9f6169c6b10194d2d8e08e718e31f&window_id=CometName_29e3ffd578035c0930a5f0e9149a702d&referer_type=left_menu#api=mm&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/EN_All.lp?28&width=1000&height=690&useApiType=mm&",
	"https://bottle2.itsrealgames.com/www/sta.html?time=1538994639991&&userId=116590694&sessionKey=b0bd10c0fa6c94d4c1fdb8d8265bb47b&authKey=37b02e73975ffb44f2c60b67e85ec22f&net_type=32&useApiType=sa&email=nemiro-f@ya.ru&locale=RU#time=1538994639991&userId=116590694&sessionKey=b0bd10c0fa6c94d4c1fdb8d8265bb47b&authKey=37b02e73975ffb44f2c60b67e85ec22f&net_type=32&useApiType=sa&email=nemiro-f@ya.ru&locale=RU&api=sa&packageName=bottlePackage&config=config_release.xml&protocol=https:&international=false&locale_url=../resources/locale/EN_All.lp?28&width=1000&height=690&",
	"https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=b42a24131eeab78ea3a31189aadc7199&vid=9947733686513067209&oid=9947733686513067209&app_id=543574&authentication_key=a13037b47258381aa4cc132f1a5f1c3b&session_expire=1528948440&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=068afb8f0ae3f6054c9645420f6e2c43&window_id=CometName_6bb080f37b95e979c9c47d182371c2a2&referer_type=left_menu#api=mm&packageName=bottlePackage&config=https://bottle2.itsrealgames.com/www/config_release.xml&protocol=https:&locale=RU&locale_url=&width=1000&height=690&useApiType=mm&",
	"https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=79a5e88f5c61f4e10abd4ce13b32402a&vid=14264633980334404830&oid=14264633980334404830&app_id=543574&authentication_key=189167eb16d0744f211c6003d7c334a0&session_expire=1528945364&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=ca6e1f0c84d394d97a68a401a0225a57&window_id=CometName_eaa61a0fd918f3f639a2b80590e03225&referer_type=left_menu#api=mm&packageName=bottlePackage&config=https://bottle2.itsrealgames.com/www/config_release.xml&protocol=https:&locale=RU&locale_url=&width=1000&height=690&useApiType=mm&",
	"https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=7556246c19559181c1bcd32682446ea2&vid=6611932282865570796&oid=6611932282865570796&app_id=543574&authentication_key=5a4b728af3bb1c239e6a0c7978b77361&session_expire=1530439148&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=dbece9235b7b4f42d80aaa7054e263c0&window_id=CometName_d8d8afe8a45d0f3fc9218526c97b438f&referer_type=left_menu#api=mm&packageName=bottlePackage&config=https://bottle2.itsrealgames.com/www/config_release.xml&protocol=https:&locale=RU&locale_url=https://bottle2.itsrealgames.com/resources/locale/EN_All_14.lp&width=1000&height=690&useApiType=mm&",
	"https://bottle2.itsrealgames.com/www/mРІ.html?is_app_user=1&session_key=7556246c19559181c1bcd32682446ea2&vid=6611932282865570796&oid=6611932282865570796&app_id=543574&authentication_key=5a4b728af3bb1c239e6a0c7978b77361&session_expire=1530439148&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=dbece9235b7b4f42d80aaa7054e263c0&window_id=CometName_d8d8afe8a45d0f3fc9218526c97b438f&referer_type=left_menu#api=mm&packageName=bottlePackage&config=https://bottle2.itsrealgames.com/www/config_release.xml&protocol=https:&locale=RU&locale_url=https://bottle2.itsrealgames.com/resources/locale/EN_All_14.lp&width=1000&height=690&useApiType=mm&",
}

var urls4 = []string{
	"https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=113180420&sessionKey=44d5464efcc34fb04e950ada31818da&authKey=f6cafa6c48ea05512cb5b08cee1568b2&net_type=32&useApiType=sa",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563079220719&&userId=113377129&sessionKey=d10b348cb0f2f8f71d7222b5ebbd77bd&authKey=2450d7fefa3f05bae3773244f54d24f6&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563080512759&&userId=113556057&sessionKey=b06761b682dc102272d4445353e47bf2&authKey=ab37502366a3a0b1cbd4b829ce779f76&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563080582929&&userId=113556072&sessionKey=08274282abadf20dcdc2b1f2e04796de&authKey=93e230f79e0024a07968739d8db4658f&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563454016409&&userId=114025257&sessionKey=16ea465a380afe08392ac92122852360&authKey=f2e10c36a7109b7bf6ba460c5c74b910&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=114029836&sessionKey=ea77682f0335546907b02c9406b5550&authKey=1c069f0c2b03684b7faca31be073a5f8&net_type=32&useApiType=sa",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563073122862&&userId=114150349&sessionKey=8b8abc6941525c8d4e4c379685f7d2ad&authKey=3264d360f3f850ba75381b9a46b433ca&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563075877774&&userId=114393355&sessionKey=a0243d0c686684f4c7c564f36820e09a&authKey=a989990d0566e783158e67cbe1764cee&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563076057275&&userId=114408341&sessionKey=e30fdfc53b999b7586fe1709e35d15b9&authKey=20fa44e5d14b4d097537723353e6522c&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563457171449&&userId=114714220&sessionKey=842533d864e6692e83b3d73e6e532b90&authKey=8ca262d35a4eb01f0deac9f94da01b41&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563080660475&&userId=114822291&sessionKey=2ed15743bf9f133e2bb08cb88ca88121&authKey=f0af3bf882f2bd21111a546f453c058a&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=114831235&sessionKey=aa0c9b114a00209952ed92a84451f49&authKey=a7655f8fb6b2e477e2fda88fd570951b&net_type=32&useApiType=sa",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563449618648&&userId=114880893&sessionKey=b7eaa8623a2f6d533a8e1dda00347ca0&authKey=8203a700274b07c6352f378c5382e6f3&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563648118505&&userId=115075281&sessionKey=87fa39080ed4180eb04090ffde387c7e&authKey=238c555b12c5e3388b233044f4e4c002&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563454536465&&userId=115397550&sessionKey=6e765eac1762b50d2d1990bd5e15ac32&authKey=ee191cb9d9c9c644116831f851b377a3&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563648172687&&userId=115491156&sessionKey=12f792053150c821a8a17ebcadfc7eb1&authKey=c6ca601c9bd4e8af5395da380b1b1eb4&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563648229502&&userId=115491166&sessionKey=236aa4a850e9af14525e6535e63d7bda&authKey=b95f6956730cec56e1b97aa435bb0709&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563096977496&&userId=115504841&sessionKey=730488c27f3265f5c632f44c2c4febac&authKey=5bc2cbb9b3f1bb99ddbb05c81b72b06c&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563124387068&&userId=115504847&sessionKey=de689dd99b1d0044bd36ad66b60f3c98&authKey=4eaf80a7bee37ad611dc14481c59f427&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563457201210&&userId=115504868&sessionKey=2f87212213a4f7cb2be1b25bf64be5cf&authKey=174f7acfd520229df24dc4cdd4626f71&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563648079142&&userId=115504890&sessionKey=07a816381ca1d39326cabc82b3d44e6d&authKey=b628a1369ea5bcc1d0ffd04c3f0fc1df&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563452497556&&userId=115647615&sessionKey=e62f52a12fbad6a3cf8ab579991e848c&authKey=eea7168c0934af302dfb15c5f1b191fd&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563458625357&&userId=117060960&sessionKey=6bbae1e2ebe4ca773f7f697dfbb63376&authKey=9a7464abe6e7f55276feace4451b58dc&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563456299593&&userId=117177330&sessionKey=8bf1314896492c131ee408eb6c73c4c3&authKey=472ff241f06fb6acc32cfab308f3b017&net_type=32&useApiType=s",
	"https://bottle2.itsrealgames.com/www/sa.html?time=1563075053056&&userId=117390712&sessionKey=197a15901c03f930bbd9117146291bed&authKey=afe06e0e238d5c39836f4973ac95b0d1&net_type=32&useApiType=s",
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

	for _, url := range urls4[6:9] {

		wg.Add(1)

		login := getLoginPacket(url)

		config := socket.GetDefaultSocketConfig()

		// homeDir, err := os.UserHomeDir()
		// if err == nil {
		// 	logDir := path.Join(homeDir, "logs")
		// 	err := os.MkdirAll(logDir, 0700)
		// 	if err == nil {
		// 		fileName := fmt.Sprintf("%s%v.log", frame.GetFrameTypeName(login.NetType), login.LoginID)
		// 		logPath := path.Join(logDir, fileName)
		// 		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0700)
		// 		if err == nil {
		// 			defer file.Close()
		// 			config.Logger = log.New(file, "", log.Ltime|log.Lshortfile)
		// 		}
		// 	}
		// }

		//config.TimeInTheGame = 5

		g := socket.NewSocket(config)
		g.SetOpenHandler(openHandle)
		g.SetCloseHandler(closeHandle)
		g.SetErrorHandler(errorHandle)
		g.SetRecvHandler(recvHandler)

		games = append(games, g)

		if err := g.Connection(); err != nil {
			// Error: Auth Failed (code: ip_forbidden)
			g.Logf("Conn fail %v", err)
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
				// sender.Send(client.BOTTLE_PLAY, &client.BottlePlay{
				// 	RoomID: 0,
				// })
				// sender.Send(client.MOVE, &client.Move{
				// 	PlayerID: Tototo93,
				// })

				// go func() {
				// 	<-time.After(time.Second * 5)
				// 	photoLike := client.PhotoLike{
				// 		PlayerID: Tototo93,
				// 		PhotoID:  1,
				// 		IsLike:   1,
				// 	}
				// 	sender.Send(client.PHOTOS_LIKE, &photoLike)
				// }()

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

const captcha = ""
