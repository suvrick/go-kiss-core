package server

import "github.com/suvrick/go-kiss-core/types"

const BOTTLE_KISS PacketServerType = 30

// BOTTLE_KISS(30) "IB"
type BottleKiss struct {
	PlayerID types.I
	Answer   KissAnswer
}

type KissAnswer types.B

func (kiss KissAnswer) String() string {
	switch kiss {
	case 0:
		return "NO"
	case 1:
		return "YES"
	case 2:
		return "SKIP"
	default:
		return "NO NAME FOR ACTION"
	}
}
