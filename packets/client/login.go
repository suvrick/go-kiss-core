package client

const LOGIN PacketClientType = 4

// LOGIN(4) "IBBS,BSIIBSBSBS"
type Login struct {
	ID           uint64
	NetType      uint16
	DeviceType   byte
	Key          string
	OAuth        byte
	AccessToken  string
	Referrer     int
	Tag          int
	FieldInt     byte
	FieldString  string
	RoomLanguage byte
	FieldString2 string
	Gender       byte
	Captcha      string
}
