package client

import "github.com/suvrick/go-kiss-core/types"

const LOGIN types.PacketClientType = 4

// LOGIN(4) "IBBS,BSIIBSBSBS"
type Login struct {
	ID           types.L
	NetType      types.I
	DeviceType   types.I
	Key          types.S
	OAuth        types.B `pack:"optional"`
	AccessToken  types.S `pack:"optional"`
	Referrer     types.I `pack:"optional"`
	Tag          types.I `pack:"optional"`
	FieldInt     types.B `pack:"optional"`
	FieldString  types.S `pack:"optional"`
	RoomLanguage types.B `pack:"optional"`
	FieldString2 types.S `pack:"optional"`
	Gender       types.B `pack:"optional"`
	Captcha      types.S `pack:"optional"`
}
