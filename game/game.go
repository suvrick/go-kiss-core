package game

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/packets"
	"github.com/suvrick/go-kiss-core/types"
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
}

func NewGame() *Game {
	g := Game{
		end:      make(chan struct{}),
		s_buffer: bytes.NewBuffer(make([]byte, 1024)),
		r_buffer: bytes.NewBuffer(make([]byte, 1024)),
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

	return &g
}

func (g *Game) Send(id packets.ClientPacketType, data any) error {

	g.s_buffer.Reset()

	err := packets.NewClientPacket(id, data, g.s_buffer)
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

	packets.WriteLong(g.s_buffer, types.L(g.msgID))
	b2 := make([]byte, len(g.s_buffer.Bytes()))
	copy(b2, g.s_buffer.Bytes()) // msgID
	g.s_buffer.Reset()

	packets.WriteLong(g.s_buffer, types.L(len(b1)+len(b2))) // len
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
	var packetID types.I
	var pack []any
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
		_, err = packets.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		//read massege id
		_, err = packets.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		packetID, err = packets.ReadInt(g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		pack, err = packets.NewServerPacket(packets.ServerPacketType(packetID), g.r_buffer)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("packetID: %d, data: %v, error: %v\n", packetID, pack, err)
	}

	fmt.Println("Close")
}
