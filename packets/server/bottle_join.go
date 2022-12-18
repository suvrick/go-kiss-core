package server

const BOTTLE_JOIN PacketServerType = 26

// BOTTLE_JOIN(26) "IB"
type BottleJoin struct {
	IntField  int64
	ByteField byte
}
