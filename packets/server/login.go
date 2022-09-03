package server

const LOGIN PacketServerType = 4

// LOGIN(4) "B,II"
type Login struct {
	Result  byte
	UserID  uint64
	Balance int32
}
