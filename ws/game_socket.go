package ws

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"github.com/suvrick/go-kiss-core/packets"
)

type GameSocket struct {
	socket     *Socket
	msgID      int64
	Done       chan struct{}
	CloseEvent func()
	botID      uint32
	bot        packets.Bot
}

func NewGameSocket(config *GameSocketConfig) *GameSocket {

	gs := GameSocket{
		Done:  make(chan struct{}),
		msgID: 0,
		bot:   packets.Bot{},
	}

	socket := NewSocket(config.SocketConfig)
	socket.SetErrorHandler(gs.ErrorHandler)
	socket.SetOpenHandler(gs.OpenHandler)
	socket.SetCloseHandler(gs.CloseHandler)
	socket.SetReadHandler(gs.ReadHandler)

	gs.socket = socket

	return &gs
}

func (gs *GameSocket) SetBotID(id uint32) {
	gs.botID = id
}

func (gs *GameSocket) OpenHandler() {
	log.Printf("[OPEN(%d)] socket open\n", gs.botID)
}

func (gs *GameSocket) CloseHandler(rule byte, msg string) {
	log.Printf("[CLOSE(%d)] socket close: %s\n", gs.botID, msg)
	gs.CloseEvent()
	gs.Done <- struct{}{}
}

func (gs *GameSocket) ErrorHandler(err error) {
	log.Printf("[ERROR(%d)] socket error: %s\n", gs.botID, err.Error())

	if err == ErrProxyConnectionFail {
		//gs.reconnect()
		return
	}

	gs.GameOver()
}

func (gs *GameSocket) ReadHandler(reader io.Reader) {

	p := packets.CreateServerPacket(reader)
	if p.Error != nil {
		if p.Type == 306 {
			gs.GameOver()
			return
		}

		log.Println(fmt.Errorf("[error(%d)] %s. packetType: %d", gs.botID, p.Error.Error(), p.Type))
		return
	}

	for _, t := range []uint16{4, 5, 7, 13, 17, 130, 310} {

		if t != p.Type {
			continue
		}

		log.Printf("[READ(%d)] %s(%d) Format: %#v, Data: %v\n", gs.botID, p.Name, p.Type, p.Format, p.Params)
	}

	p.Fill(&gs.bot)

	gs.parse(p)

	p = nil
	reader = nil
	//runtime.GC()
}

func (gs *GameSocket) GameOver() {
	gs.socket.Close()
}

func (gs *GameSocket) Send(t uint16, data []interface{}) {

	p := packets.CreateClientPacket(t, data...)

	if p.Error != nil {
		gs.ErrorHandler(p.Error)
		return
	}

	p.GetBuffer(gs.msgID)
	if p.Error != nil {
		gs.ErrorHandler(p.Error)
		return
	}

	log.Printf("[SEND(%d)] %s(%d) Format: %#v, Data: %v, Buffer: %v\n",
		gs.botID,
		p.Name,
		p.Type,
		p.Format,
		p.Params,
		p.Buffer)

	atomic.AddInt64(&gs.msgID, 1)

	gs.socket.Send(p.Buffer)
}

func (gs *GameSocket) Run() {
	gs.socket.Go()
}

func (gs *GameSocket) parse(p *packets.Packet) {

	//log.Println(p)

	switch p.Type {
	case 4:
		switch gs.bot.Result {
		case 0:
			gs.Send(61, []interface{}{})
		default:
			gs.ErrorHandler(fmt.Errorf("bad auth. Result: %d", gs.bot.Result))
		}
	case 9:
		//gs.GameOver()
	case 13:

		if gs.bot.RewardsHistory == nil {
			gs.bot.RewardsHistory = make([]packets.Reward, 0)
		}

		for _, v := range gs.bot.Rewards {

			gs.Send(11, []interface{}{v.ID})
			gs.bot.RewardsHistory = append(gs.bot.RewardsHistory, v)
		}

		gs.bot.Rewards = make([]packets.Reward, 0)
	default:
	}
}
