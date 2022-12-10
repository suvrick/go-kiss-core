package game

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-leb128"
)

type Chunk struct {
	ID        int
	MetaID    int
	Index     int
	Type      rune
	Name      string
	IsRequred bool

	Parent *Chunk
}

type Meta struct {
	ID           int
	PacketID     uint64
	PacketFormat string
	PacketName   string
	PacketType   int
}

var metes = []*Meta{
	{
		ID:           1,
		PacketID:     4,
		PacketFormat: "IBBS,BSIIBSBSBS",
		PacketName:   "LOGIN",
		PacketType:   1,
	},
	{
		ID:           2,
		PacketID:     4,
		PacketFormat: "BII",
		PacketName:   "LOGIN",
		PacketType:   0,
	},
}

var chunks = []*Chunk{
	{
		ID:        1,
		MetaID:    1,
		Index:     0,
		Type:      'I',
		Name:      "packet_id",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        2,
		MetaID:    1,
		Index:     1,
		Type:      'B',
		Name:      "device",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        3,
		MetaID:    1,
		Index:     2,
		Type:      'I',
		Name:      "login_id",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        4,
		MetaID:    1,
		Index:     3,
		Type:      'I',
		Name:      "frame_type",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        6,
		MetaID:    1,
		Index:     5,
		Type:      'S',
		Name:      "key",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        5,
		MetaID:    1,
		Index:     4,
		Type:      'B',
		Name:      "device",
		IsRequred: true,
		Parent:    nil,
	},

	// {
	// 	ID:        7,
	// 	MetaID:    2,
	// 	Index:     0,
	// 	Type:      'I',
	// 	Name:      "packet_id",
	// 	IsRequred: true,
	// 	Parent:    nil,
	// },
	// {
	// 	ID:        8,
	// 	MetaID:    2,
	// 	Index:     1,
	// 	Type:      'B',
	// 	Name:      "device",
	// 	IsRequred: true,
	// 	Parent:    nil,
	// },
	{
		ID:        9,
		MetaID:    2,
		Index:     2,
		Type:      'B',
		Name:      "result",
		IsRequred: true,
		Parent:    nil,
	},
	{
		ID:        10,
		MetaID:    2,
		Index:     3,
		Type:      'I',
		Name:      "game_id",
		IsRequred: false,
		Parent:    nil,
	},
	{
		ID:        11,
		MetaID:    2,
		Index:     4,
		Type:      'I',
		Name:      "balance",
		IsRequred: false,
		Parent:    nil,
	},
}

type SortByChunk []*Chunk

func (a SortByChunk) Len() int           { return len(a) }
func (a SortByChunk) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByChunk) Less(i, j int) bool { return a[i].Index < a[j].Index }

func GetMeta(packetType int, packetID uint64) *Meta {
	var m *Meta
	for _, v := range metes {
		if v.PacketType == packetType && v.PacketID == packetID {
			m = v
			break
		}
	}
	return m
}

func GetScheme(metaID int) []*Chunk {
	result := make([]*Chunk, 0)
	for _, chunk := range chunks {
		if chunk.MetaID == metaID {
			result = append(result, chunk)
		}
	}
	return result
}

func GetSubScheme(chunks []*Chunk, parent *Chunk) []*Chunk {
	result := make([]*Chunk, 0)
	for _, chunk := range chunks {
		if chunk.Parent == parent {
			result = append(result, chunk)
		}
	}
	return result
}

const host = "wss://bottlews.itsrealgames.com"

type Game struct {
	socket       *websocket.Conn
	msgID        uint32
	header       http.Header
	end          chan struct{}
	stop_packets []int

	s_buffer *bytes.Buffer
	r_buffer *bytes.Buffer
	cancel   context.CancelFunc
	ctx      context.Context

	s_packet map[string]interface{}
	c_packet map[string]interface{}
}

func NewGame() *Game {
	g := Game{
		end:          make(chan struct{}),
		stop_packets: make([]int, 0),
		s_buffer:     new(bytes.Buffer),
		r_buffer:     new(bytes.Buffer),
		s_packet:     make(map[string]interface{}),
		c_packet:     make(map[string]interface{}),
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

func (g *Game) Send(p map[string]interface{}) error {
	err := g.Marshal(p)
	if err != nil {
		return err
	}
	g.send()
	return nil
}

func (g *Game) End() chan struct{} {
	return g.end
}

func (g *Game) Close() {
	g.socket.Close()
	close(g.end)
}

func (g *Game) SetStopPacketsID(ids []int) {
	g.stop_packets = append(g.stop_packets, ids...)
}

func (g *Game) Connect() error {

	dialer := websocket.Dialer{
		HandshakeTimeout: (time.Second * 60),
	}

	s, _, err := dialer.DialContext(g.ctx, host, g.header)

	g.socket = s

	if err != nil {
		//g.emitErrorHandler(err)
		return err
	}

	go func() {
		g.read()
	}()

	return nil
}

func (g *Game) ConnectWithProxy(proxy *url.URL) error {

	dialer := websocket.Dialer{
		HandshakeTimeout: (time.Second * 60),
		Proxy:            http.ProxyURL(proxy),
	}

	s, _, err := dialer.Dial(host, g.header)
	g.socket = s

	if err != nil {
		//gs.emitErrorHandler(err)
		return err
	}

	go func() {
		g.read()
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

func (g *Game) clearServerPacket() {
	g.s_packet = make(map[string]interface{}, 0)
}

func (g *Game) read() {

	defer g.socket.Close()

	for g.socket != nil {

		_, msg, err := g.socket.ReadMessage()

		if err != nil {
			//s.emitErrorHandler(err)
			break
		}

		g.r_buffer.Reset()
		g.r_buffer.Write(msg)

		msgLen, err := leb128.ReadUint(g.r_buffer, 32)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		msgID, err := leb128.ReadUint(g.r_buffer, 32)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		packetID, err := leb128.ReadUint(g.r_buffer, 8)
		if err != nil {
			//s.emitErrorHandler(err)
			continue
		}

		g.loop(packetID, g.r_buffer)

		fmt.Printf("Recv >> packetID: %d,msgID: %d, msgLen: %d\n", packetID, msgID, msgLen)
	}

	fmt.Println("Close")
}

func (g *Game) loop(packetID uint64, b *bytes.Buffer) {

	for _, v := range g.stop_packets {
		if uint64(v) == packetID {
			g.Close()
			return
		}
	}

	err := g.Unmarshal(packetID, b)
	if err != nil {
		// logger
		return
	}

	switch packetID {
	case 4:
		fmt.Printf("%v\n", g.s_packet)
		//g.Close()
	}
}

func (g *Game) Unmarshal2(r io.ByteReader, schemes []*Chunk) error {

	var value interface{}
	var err error

	ss := SortByChunk(schemes)
	sort.Sort(ss)

	for _, chunk := range schemes {
		switch chunk.Type {
		case 'B':
			value, err = leb128.Read(r, chunk.Type)
		case 'I':
			value, err = leb128.Read(r, chunk.Type)
		case 'L':
			value, err = leb128.Read(r, chunk.Type)
		case 'S':
			value, err = leb128.Read(r, chunk.Type)
		default:
			err = fmt.Errorf("unsupported type: v", chunk.Type)
		}

		if err != nil {
			if chunk.IsRequred {
				return err
			}
		} else {
			g.s_packet[chunk.Name] = value
		}
	}

	return nil
}

func (g *Game) Unmarshal(packetID uint64, r io.ByteReader) error {

	g.clearServerPacket()

	g.s_packet["packet_type"] = 0
	g.s_packet["packet_id"] = packetID

	m := GetMeta(0, packetID)
	if m == nil {
		return fmt.Errorf("not found meta for packetID: %d", packetID)
	}

	schemes := GetScheme(m.ID)
	if schemes == nil {
		return fmt.Errorf("not found meta for packetID: %d", packetID)
	}

	return g.Unmarshal2(r, schemes)
}

func (g *Game) Marshal2(schemes []*Chunk, values map[string]interface{}) error {

	ss := SortByChunk(schemes)
	sort.Sort(ss)

	for _, chunk := range schemes {

		val, ok := values[chunk.Name]
		if !ok {
			return fmt.Errorf("empty value by %s field", chunk.Name)
		}

		switch chunk.Type {
		case 'B':
			leb128.Write(g.s_buffer, val)
		case 'I':
			leb128.Write(g.s_buffer, val)
		case 'S':
			leb128.Write(g.s_buffer, val)
		case 'A':
			arr, ok := val.([]map[string]interface{})
			if !ok {
				return fmt.Errorf("bad sengnature")
			}

			schemes2 := GetSubScheme(schemes, chunk)
			for _, v := range arr {
				err := g.Marshal2(schemes2, v)
				if err != nil {
					return fmt.Errorf("bad sengnature")
				}
			}

		default:
			return fmt.Errorf("unsupported code %v", chunk.Type)
		}
	}

	return nil
}

func (g *Game) Marshal(values map[string]interface{}) error {

	g.c_packet = values
	g.s_buffer.Reset()

	// получить мета данные пакета
	packetID, ok := g.c_packet["packet_id"].(uint64)
	if !ok {
		return fmt.Errorf("empty packet_id")
	}

	m := GetMeta(1, packetID)
	if m == nil {
		return fmt.Errorf("not found meta by packet_id: %d", packetID)
	}

	schemes := GetScheme(m.ID)
	if schemes == nil {
		return fmt.Errorf("not found scheme by packet_id: %d, meta_id: %d", packetID, m.ID)
	}

	// TODO: need sort
	return g.Marshal2(schemes, g.c_packet)
}
