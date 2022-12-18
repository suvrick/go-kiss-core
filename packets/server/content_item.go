package server

const CONTEST_ITEMS PacketServerType = 10

// CONTEST_ITEMS(10) "B,B"
type ContestItems struct {
	ByteField  byte
	ByteField2 byte
}
