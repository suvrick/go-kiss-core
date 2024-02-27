package client

import (
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const LOGIN types.PacketClientType = 4

// LOGIN(4) "IBBS,BSIIBSBSBS"
type Login struct {
	LoginID      string
	NetType      byte
	DeviceType   byte
	Key          string
	OAuth        byte
	AccessToken  string
	StringField  string
	ByteField    byte
	ByteField1   byte
	ByteField2   byte
	StringField2 string
	Captcha      string
}

func (login Login) String() string {
	return "LOGIN(4)"
}

func (login *Login) Marshal() ([]byte, error) {
	data, err := leb128.WriteBigNumber(nil, login.LoginID)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.NetType)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.DeviceType)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteString(data, login.Key)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.OAuth)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteString(data, login.AccessToken)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteString(data, login.StringField)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.ByteField)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.ByteField1)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteByte(data, login.ByteField2)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteString(data, login.StringField2)
	if err != nil {
		return nil, err
	}

	data, err = leb128.WriteString(data, login.Captcha)
	if err != nil {
		return nil, err
	}

	return data, nil
}
