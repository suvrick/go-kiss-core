package ws

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sync/atomic"

	"github.com/suvrick/go-kiss-core/packets"
)

type GameSocket struct {
	socket *Socket
	msgID  int64
	Done   chan struct{}
	bot    packets.Bot
}

func NewGameSocket(config *GameSocketConfig) *GameSocket {

	gs := GameSocket{
		Done:  make(chan struct{}),
		msgID: 0,
		bot:   packets.Bot{},
	}

	sock := NewSocket(config.SocketConfig)
	sock.SetErrorHandler(gs.onError)
	sock.SetOpenHandler(gs.onOpen)
	sock.SetCloseHandler(gs.onClose)
	sock.SetReadHandler(gs.onRead)

	gs.socket = sock

	return &gs
}

func (gs *GameSocket) onOpen() {
	log.Println("socket open")
}

func (gs *GameSocket) onClose(rule byte, msg string) {
	log.Printf("socket close: %s\n", msg)
	log.Printf("%#v", gs.bot)
	gs.Done <- struct{}{}
}

func (gs *GameSocket) onError(err error) {
	log.Printf("socket error: %s\n", err.Error())

	if err == ErrProxyConnectionFail {
		//gs.reconnect()
		return
	}
}

func (gs *GameSocket) onRead(reader io.Reader) {
	p := packets.CreateServerPacket(reader)
	if p.Error != nil {
		gs.onError(fmt.Errorf("%s. packetType: %d", p.Error.Error(), p.Type))
	}

	log.Printf("%s(%d) Format: %#v, Data: %v\n", p.Name, p.Type, p.Format, p.Params)

	p.Fill(&gs.bot)

	gs.parse(p)

	memprofile := "C:\\Users\\suvrick\\Desktop\\log.out"
	f, _ := os.Create(memprofile)
	//runtime.GC() // get up-to-date statistics
	pprof.WriteHeapProfile(f)
	f.Close()
}

func (gs *GameSocket) GameOver() {
	gs.socket.Close()
}

func (gs *GameSocket) Send(t uint16, data []interface{}) {

	p := packets.CreateClientPacket(t, data...)

	if p.Error != nil {
		gs.onError(p.Error)
		return
	}

	buf, err := p.GetBuffer(gs.msgID)
	if err != nil {
		gs.onError(err)
		return
	}

	atomic.AddInt64(&gs.msgID, 1)

	gs.socket.Send(buf)
}

func (gs *GameSocket) Run() {
	gs.connect()
}

func (gs *GameSocket) connect() {
	// proxy_str := "zproxy.lum-superproxy.io:22225:lum-customer-c_07f044e7-zone-static:hcx7fnqnph27"
	// proxy := proxy.Parse(proxy_str, ":")
	// proxy = nil
	//gs.socket.SetProxy(proxy)

	gs.socket.Go()
}

func (gs *GameSocket) reconnect() {
	gs.socket.connect()
}

func (gs *GameSocket) parse(p *packets.Packet) {

	//log.Println(p)

	switch p.Type {
	case 4:
		switch gs.bot.Result {
		case 0:
			gs.Send(61, []interface{}{})
		default:
			gs.onError(fmt.Errorf("bad auth. Result: %d", gs.bot.Result))
			gs.GameOver()
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
