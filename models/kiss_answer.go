package models

import "github.com/suvrick/go-kiss-core/types"

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
