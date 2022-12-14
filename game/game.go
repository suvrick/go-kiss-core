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
	"github.com/suvrick/go-leb128"
)

const host = "wss://bottlews.itsrealgames.com"

type Game struct {
	socket       *websocket.Conn
	msgID        uint32
	header       http.Header
	end          chan struct{}
	proxy        *url.URL
	stop_packets []int

	s_buffer *bytes.Buffer
	r_buffer *bytes.Buffer

	cancel context.CancelFunc
	ctx    context.Context
}

func NewGame() *Game {
	g := Game{
		end:          make(chan struct{}),
		stop_packets: make([]int, 0),
		s_buffer:     new(bytes.Buffer),
		r_buffer:     new(bytes.Buffer),
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

func (g *Game) Send(id packets.ClientPacketType, data []any) error {
	b, err := packets.NewClientPacket(id, data)
	if err != nil {
		return err
	}
	_ = b
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

	leb128.Write(g.s_buffer, g.msgID)
	b2 := make([]byte, len(g.s_buffer.Bytes()))
	copy(b2, g.s_buffer.Bytes()) // msgID
	g.s_buffer.Reset()

	leb128.Write(g.s_buffer, len(b1)+len(b2)) // len
	g.s_buffer.Write(b2)
	g.s_buffer.Write(b1)

	log.Printf("Send >> data: %x\n", g.s_buffer.Bytes())

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
	var packetID int64
	var pack []any
	var msg []byte

	for g.socket != nil {

		_, msg, err = g.socket.ReadMessage()
		if err != nil {
			//s.emitErrorHandler(err)
			break
		}

		g.r_buffer.Reset()
		g.r_buffer.Write(msg)

		// read length packet
		_, err = leb128.ReadInt(g.r_buffer, 32)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		//read massege id
		_, err = leb128.ReadInt(g.r_buffer, 32)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		packetID, err = leb128.ReadInt(g.r_buffer, 8)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		pack, err = packets.NewServerPacket(packets.ServerPacketType(packetID), g.r_buffer)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		fmt.Printf("packetID: %d, data: %v, error: %v\n", packetID, pack, err)
	}

	fmt.Println("Close")
}
