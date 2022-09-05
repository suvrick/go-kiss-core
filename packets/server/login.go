package server

const LOGIN PacketServerType = 4

// LOGIN(4) "B,II"
type Login struct {
	Result  LoginResultType
	GameID  uint64
	Balance uint
}

// Login result response
const (
	Success       LoginResultType = 0
	Failed        LoginResultType = 1
	Exist         LoginResultType = 2
	Blocked       LoginResultType = 3
	WronngVersion LoginResultType = 4
	NoSex         LoginResultType = 5
	Captcha       LoginResultType = 6
	BlockedByAge  LoginResultType = 7
	NeedVerify    LoginResultType = 8
	Deleted       LoginResultType = 9
)

func (r LoginResultType) String() string {
	switch r {
	case 0:
		return "Success"
	case 1:
		return "Failed"
	case 2:
		return "Exist"
	case 3:
		return "Blocked"
	case 4:
		return "WronngVersion"
	case 5:
		return "NoSex"
	case 6:
		return "Captcha"
	case 7:
		return "BlockedByAge"
	case 8:
		return "NeedVerify"
	case 9:
		return "Deleted"
	default:
		return "Error"
	}
}
