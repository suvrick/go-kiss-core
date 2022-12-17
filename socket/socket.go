package socket

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrTimeoutTheGame           = errors.New("time in the game success")
	ErrConnectionNot101         = errors.New("websocket connection fail")
	ErrConnectionFail           = errors.New("connection fail")
	ErrProxyConnectionFail      = errors.New("proxy connection fail")
	ErrCloseConnectionByTimeout = errors.New("close connection by timeout")
)

// GameSock ...
type Socket struct {
	config *SocketConfig

	rule_close byte
	client     *websocket.Conn

	proxy *url.URL

	wg sync.WaitGroup

	openHandle  func()
	closeHandle func(byte, string)
	readHandle  func(io.Reader)
	errorHandle func(error)
}

const (
	NORMAL_CLOSE        = 0x00
	ERROR_CONNECT_CLOSE = 0x01
	ERROR_READ_CLOSE    = 0x02
	ERROR_SEND_CLOSE    = 0x03
	ERROR_TIMEOUT_CLOSE = 0x04
)

var closed_rules = map[byte]string{
	0x00: "Normal close",
	0x01: "Connection error",
	0x02: "Read error",
	0x03: "Send error",
	0x04: "Timeout in the game",
}

func NewSocket(config *SocketConfig) *Socket {
	return &Socket{
		config:     config,
		wg:         sync.WaitGroup{},
		rule_close: 255,
	}
}

func (socket *Socket) Connection() error {

	dialer := websocket.Dialer{
		HandshakeTimeout: (socket.config.ConnectTimeout),
	}

	if socket.proxy != nil {
		dialer.Proxy = http.ProxyURL(socket.proxy)
	}

	client, resp, err := dialer.Dial(socket.config.Host, socket.config.Head)

	if err != nil {
		socket.setClosedRule(ERROR_CONNECT_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(err)
			return ErrConnectionFail
		}
	}

	socket.client = client

	if resp != nil {
		if resp.StatusCode != 101 {
			socket.setClosedRule(ERROR_CONNECT_CLOSE)
			// При корректном открытии ws соединения, код должен быть 101
			if socket.errorHandle != nil {
				socket.errorHandle(err)
			}
			return ErrConnectionNot101
		}
	}

	//Таймаут после которого игра закроется.
	if socket.openHandle != nil {
		socket.openHandle()
	}

	socket.wg.Add(1) // for read
	go socket.timeoutToGame()
	go socket.read()
	go socket.done()

	return nil
}

func (socket *Socket) SetOpenHandler(handler func()) {
	socket.openHandle = handler
}

func (socket *Socket) SetCloseHandler(handler func(rule byte, msg string)) {
	socket.closeHandle = handler
}

func (socket *Socket) SetReadHandler(handler func(reader io.Reader)) {
	socket.readHandle = handler
}

func (socket *Socket) SetErrorHandler(handler func(err error)) {
	socket.errorHandle = handler
}

func (socket *Socket) Send(packet []byte) {

	if socket.client == nil {
		socket.errorHandle(ErrConnectionFail)
		return
	}

	err := socket.client.WriteMessage(websocket.BinaryMessage, packet)

	if err != nil {
		socket.setClosedRule(ERROR_SEND_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(err)
		}
	}
}

func (socket *Socket) Close() {
	socket.setClosedRule(NORMAL_CLOSE)
	socket.close_connection()
}

func (socket *Socket) getCloseRuleMsg() string {
	return closed_rules[socket.rule_close]
}

func (socket *Socket) timeoutToGame() {
	<-time.After(time.Duration(socket.config.TimeInTheGame) * time.Second)
	socket.setClosedRule(ERROR_TIMEOUT_CLOSE)
	socket.close_connection()
}

func (socket *Socket) done() {
	socket.wg.Wait()
	if socket.closeHandle != nil {
		socket.closeHandle(socket.rule_close, socket.getCloseRuleMsg())
	}
}

func (socket *Socket) read() {

	defer func() {
		socket.wg.Done()
	}()

	for socket.client != nil {

		_, msg, err := socket.client.ReadMessage()

		if err != nil {
			socket.setClosedRule(ERROR_READ_CLOSE)
			if socket.errorHandle != nil {
				socket.errorHandle(err)
			}
			break
		}

		if socket.readHandle != nil {
			socket.readHandle(bytes.NewReader(msg))
		}
	}
}

func (socket *Socket) close_connection() {
	if socket.client != nil {
		socket.client.Close()
		socket.client = nil
	}
}

func (socket *Socket) setClosedRule(rule byte) {
	if socket.rule_close == 255 {
		socket.rule_close = rule
	}
}
