package models

import "github.com/suvrick/go-kiss-core/types"

type Player struct {
	PlayerID   types.I
	Name       types.S
	Avatar     types.S
	Profile    types.S
	Sex        types.B
	Vip        types.B
	Kissed     types.I
	KissedDay  types.I
	KissedRoom types.I
}
