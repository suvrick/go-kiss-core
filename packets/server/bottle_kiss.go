package server

const BOTTLE_KISS PacketServerType = 30

// BOTTLE_KISS(30) "IB"
type BottleKiss struct {
	PlayerID  int64
	ByteField KissAnswer
}

type KissAnswer byte

func (kiss KissAnswer) String() string {
	if kiss == 0 {
		return "NO"
	} else {
		return "YES"
	}
}
