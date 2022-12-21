package server

import "github.com/suvrick/go-kiss-core/types"

const BALANCE PacketServerType = 7

// BALANCE(7) "bottles:I, reason:B"
type Balance struct {
	Bottles types.I
	Reason  types.B `pack:"optional"`
}
