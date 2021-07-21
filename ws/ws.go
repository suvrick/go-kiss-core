package ws

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-interpreter/wagon/wasm/leb128"
	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/packets"
)

type D map[string]interface{}

// GameSock ...
type game_socket struct {
	client    *websocket.Conn
	msgID     uint32
	Bot       D
	waitBonus bool
}

const host = "wss://bottlews.itsrealgames.com"

func New() (*game_socket, error) {

	gs := &game_socket{
		client: nil,
		msgID:  0,
		Bot:    make(map[string]interface{}),
	}

	gs.Bot["box"] = make([]string, 0)
	gs.Bot["logger"] = make([]string, 0)

	dialer := websocket.Dialer{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "zproxy.lum-superproxy.io:22225",
			User:   url.UserPassword("lum-customer-c_07f044e7-zone-static", "yg3h8lqi3noq"),
		}),
		HandshakeTimeout: (time.Second * 60),
	}

	con, _, err := dialer.Dial(host, nil)
	gs.client = con

	if err != nil {
		return nil, err
	}

	gs.log(INFO, "Open connetion")

	return gs, nil
}

func (gs *game_socket) Go(login map[string]interface{}) {

	if login == nil {
		gs.log(ERROR, "invalid login params")
		gs.client.Close()
		return
	}

	gs.Bot = gs.merge(gs.Bot, login)

	data := packets.NewLoginPacketClient(login)
	gs.send(data)
	gs.read()
}

func (gs *game_socket) read() {

	defer gs.client.Close()

	for {

		_, msg, err := gs.client.ReadMessage()

		if err != nil {
			gs.log(ERROR, err.Error())
			break
		}

		if len(msg) < 3 {
			continue
		}

		reader := bytes.NewReader(msg)

		msgLen, _ := leb128.ReadVarUint32(reader)
		msgID, _ := leb128.ReadVarUint32(reader)
		msgType, _ := leb128.ReadVarUint32(reader)

		log.Printf("Recv >> msgType: %d,msgID: %d, msgLen: %d\n", msgType, msgID, msgLen)

	}

}

func (gs *game_socket) send(buffer []byte) {

	msgIDBuffer := make([]byte, 0)
	msgIDBuffer = leb128.AppendUleb128(msgIDBuffer, uint64(gs.msgID))

	data := make([]byte, 0)
	data = append(data, leb128.AppendUleb128(data, uint64(len(buffer)+len(msgIDBuffer)))...)
	data = append(data, msgIDBuffer...)
	data = append(data, buffer...)

	log.Printf("Send >> data: %v\n", data)

	err := gs.client.WriteMessage(websocket.BinaryMessage, data)
	gs.msgID++

	if err != nil {
		gs.log(ERROR, err.Error())
		return
	}
}

func (gs *game_socket) merge(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {

	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}

func (gs *game_socket) log(lvl log_level, msg string) {
	gs.Bot["logger"] = append(gs.Bot["logger"].([]string), fmt.Sprintf("%s >> %s", lvl, msg))
}

type log_level string

const (
	INFO  log_level = "INFO"
	ERROR log_level = "ERROR"
)
