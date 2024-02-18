package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const CHAT_MESSAGE types.PacketServerType = 37

// CHAT_MESSAGE(37) "BIIS,II"
type ChatMessage struct {
	ByteField byte
	IntField  uint64
	WriterID  uint64
	Message   string
	//IntField3 types.I `pack:"optional"`
	//IntField4 types.I `pack:"optional"`
}

func (chatMessage *ChatMessage) Unmarshal(r *bytes.Reader) error {
	var err error

	chatMessage.ByteField, err = leb128.ReadByte(r)
	if err != nil {
		return err
	}

	chatMessage.IntField, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	chatMessage.WriterID, err = leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	chatMessage.Message, err = leb128.ReadString(r)
	if err != nil {
		return err
	}

	return err
}
