package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const CHAT_WHISPER types.PacketServerType = 38

// CHAT_WHISPER(38) "ISBIB"
type ChatWhisper struct {
	WriterID   uint64
	Message    string
	ByteField  byte
	IntField2  uint64
	ByteField2 byte
}

func (p ChatWhisper) String() string {
	return "CHAT_WHISPER(38)"
}

func (chatWhisper *ChatWhisper) Unmarshal(r *bytes.Reader) error {
	var err error

	chatWhisper.WriterID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	chatWhisper.Message, err = leb128.ReadString(r)
	if err != nil {
		return err
	}

	chatWhisper.ByteField, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	chatWhisper.IntField2, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	chatWhisper.ByteField2, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	return err
}
