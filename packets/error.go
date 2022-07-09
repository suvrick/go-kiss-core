package packets

import "errors"

var (
	ErrBadSignaturePacket = errors.New("bad signature packet")
	ErrNotFoundPacket     = errors.New("not found packet")
	ErrCreatePacket       = errors.New("error create packet")
	ErrInvalidType        = errors.New("invalid type")
)
