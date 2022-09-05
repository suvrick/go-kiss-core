package server

const BALANCE PacketServerType = 7

// BALANCE(7) "bottles:I, reason:B"
type Balance struct {
	Bottles uint
}
