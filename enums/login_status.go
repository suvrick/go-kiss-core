package enums

type LoginStatus uint8

const (
	Success LoginStatus = iota
	Failed
	Exist
	Blocked
	WrongVersion
	NoSex
	Captcha
	BlockedByAge
	NeedVerify
	Deleted
)

func (status LoginStatus) String() string {
	switch status {
	case 0:
		return "Success"
	case 1:
		return "Failed"
	case 2:
		return "Exist"
	case 3:
		return "Blocked"
	case 4:
		return "WrongVersion"
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
		return "None"
	}
}
