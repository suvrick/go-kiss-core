package client

const REQUEST PacketClientType = 8

// REQUEST (8) ""
type Request struct {
	Players []uint64
	ID      uint64
}
