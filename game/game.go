package game

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets"
	"github.com/suvrick/go-kiss-core/schemes"
)

const host = "wss://bottlews.itsrealgames.com"

type Game struct {
	socket *websocket.Conn
	msgID  uint32
	header http.Header
	end    chan struct{}
	proxy  *url.URL

	s_buffer *bytes.Buffer
	r_buffer *bytes.Buffer

	cancel context.CancelFunc
	ctx    context.Context

	GameID int

	actions map[int]func(self *Game, packet map[string]interface{})
}

func NewGame() *Game {
	g := Game{
		end:      make(chan struct{}),
		s_buffer: bytes.NewBuffer(make([]byte, 1024)),
		r_buffer: bytes.NewBuffer(make([]byte, 1024)),
		actions:  make(map[int]func(self *Game, packet map[string]interface{}), 0),
	}

	g.ctx, g.cancel = context.WithCancel(context.Background())

	g.header = http.Header{
		"Origin": {
			"https://inspin.me",
		},
		"User-Agent": {
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		},
	}

	go func() {
		<-time.After(5 * time.Second)
		g.cancel()
		log.Println("call cancel")
	}()

	return &g
}

func (g *Game) Send(id int, data map[string]interface{}) error {

	g.s_buffer.Reset()

	err := packets.NewClientPacket(g.s_buffer, id, data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	g.send()
	return nil
}

func (g *Game) End() chan struct{} {
	return g.end
}

func (g *Game) Close() {
	if g.socket != nil {
		g.socket.Close()
	}
	close(g.end)
}

func (g *Game) Connect(proxy *url.URL) error {

	dialer := websocket.Dialer{
		HandshakeTimeout: (time.Second * 60),
	}

	if proxy != nil {
		g.proxy = proxy
		dialer.Proxy = http.ProxyURL(proxy)
	}

	s, _, err := dialer.DialContext(g.ctx, host, g.header)

	g.socket = s

	if err != nil {
		//g.emitErrorHandler(err)
		return err
	}

	go func() {
		g.loop()
	}()

	return nil
}

func (g *Game) send() {

	b1 := make([]byte, len(g.s_buffer.Bytes()))
	copy(b1, g.s_buffer.Bytes()) //пакет

	g.s_buffer.Reset()

	b2, _ := leb128.WriteLong(g.msgID)
	b4, _ := leb128.WriteLong(len(b1) + len(b2)) // len

	g.s_buffer.Write(b4)
	g.s_buffer.Write(b2)
	g.s_buffer.Write(b1)

	log.Printf("[Send] data: %v\n", g.s_buffer.Bytes())

	err := g.socket.WriteMessage(websocket.BinaryMessage, g.s_buffer.Bytes())

	g.msgID++

	if err != nil {
		//gs.emitErrorHandler(err)
		return
	}
}

func (g *Game) loop() {

	defer func() {
		// g.Close()
	}()

	var err error
	var packetID int
	var pack map[string]interface{}
	var msg []byte

	for g.socket != nil {

		_, msg, err = g.socket.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			//s.emitErrorHandler(err)
			break
		}

		g.r_buffer.Reset()
		g.r_buffer.Write(msg)

		// read length packet
		_, err = leb128.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		//read massege id
		_, err = leb128.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		packetID, err = leb128.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		pack, err = packets.NewServerPacket(g.r_buffer, packetID)
		if err != nil {
			fmt.Printf("[read] %s\n", err.Error())
			continue
		}

		scheme := schemes.FindScheme(0, packetID)
		if scheme != nil {
			fmt.Printf("[read] %s(%d), format: %#v, data: %v, error: %v\n", scheme.PacketName, scheme.PacketID, scheme.PacketFormat, pack, err)
			p, _ := json.MarshalIndent(pack, "", " ")
			fmt.Println(string(p))
		}

		g.use(packetID, pack)
	}

	fmt.Println("loop close")
}

func (g *Game) use(packetID int, packet map[string]interface{}) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Use] error do action for server packetID: %d, packet: %v\n", packetID, packet)
		}
	}()

	if act, ok := g.actions[packetID]; ok {
		act(g, packet)
	}
}

func (g *Game) AddListen(packetID int, do func(self *Game, packet map[string]interface{})) {
	g.actions[packetID] = do
}
