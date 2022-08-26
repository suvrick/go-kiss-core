package ws_test

import (
	"log"
	"testing"
)

func TestGameSocket(t *testing.T) {

	p := Login{
		PacketBase: PacketBase{
			0,
			4,
		},
	}

	p.Load([]interface{}{1, 2, 3})

	p2 := Bonus{
		PacketBase: PacketBase{
			0,
			61,
		},
	}

	p2.Load([]interface{}{1, 5})

}

type IPacket interface {
	Load([]interface{})
}

type PacketBase struct {
	PacketType byte
	PacketID   uint16
}

func (p *PacketBase) load(params []interface{}) {
	for _, v := range params {
		log.Println(v)
	}
}

type Login struct {
	PacketBase
	Result  int
	UserID  uint64
	Balance int
}

func (l *Login) Load(params []interface{}) {
	l.PacketBase.load(params)
}

type Bonus struct {
	PacketBase
	IsCan int
	Day   int
}

func (l *Bonus) Load(params []interface{}) {
	l.PacketBase.load(params)
}
