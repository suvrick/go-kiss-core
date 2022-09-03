package ws

import (
	"bytes"
	"errors"
	"io"
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
	*SocketConfig

	rule_close byte
	client     *websocket.Conn

	wg sync.WaitGroup

	send_chan  chan []byte
	read_chan  chan []byte
	error_chan chan error

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
		SocketConfig: config,
		wg:           sync.WaitGroup{},
		rule_close:   255,
		send_chan:    make(chan []byte),
		read_chan:    make(chan []byte),
		error_chan:   make(chan error),
	}
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

func (socket *Socket) Go() {
	socket.connect()
	socket.wg.Add(1) // for read
	go socket.read()
	go socket.done()
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

func (socket *Socket) connect() {

	dialer := websocket.Dialer{
		HandshakeTimeout: (socket.Timeout),
	}

	client, resp, err := dialer.Dial(socket.Host, socket.Head)

	socket.client = client

	if err != nil {
		socket.setClosedRule(ERROR_CONNECT_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(err)
		}
		return
	}

	if resp != nil {
		if resp.StatusCode != 101 {
			socket.setClosedRule(ERROR_CONNECT_CLOSE)
			// При корректном открытии ws соединения, код должен быть 101
			if socket.errorHandle != nil {
				socket.errorHandle(err)
			}
			return
		}
	}

	go socket.timeout()

	if socket.openHandle != nil {
		socket.openHandle()
	}
}

func (socket *Socket) timeout() {
	<-time.After(time.Duration(socket.TimeInTheGame) * time.Second)
	socket.setClosedRule(ERROR_TIMEOUT_CLOSE)
	socket.close_connection()
	//наверное закрытие по таймауту не является ошибкой
	// if socket.errorHandle != nil {
	// 	socket.errorHandle(ErrTimeoutTheGame)
	// }
}

func (socket *Socket) done() {
	socket.wg.Wait()
	socket.close_chan()
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

func (socket *Socket) close_chan() {
	close(socket.error_chan)
	socket.error_chan = nil

	close(socket.send_chan)
	socket.send_chan = nil

	close(socket.read_chan)
	socket.read_chan = nil
}

func (socket *Socket) setClosedRule(rule byte) {
	if socket.rule_close == 255 {
		socket.rule_close = rule
	}
}
