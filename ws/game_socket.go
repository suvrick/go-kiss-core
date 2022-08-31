package ws

import (
	"io"
)

type GameSocket struct {
	socket *Socket
	msgID  int64
	Done   chan struct{}
}

func NewGameSocket(config *SocketConfig) *GameSocket {

	gs := GameSocket{
		Done:  make(chan struct{}),
		msgID: 0,
	}

	socket := NewSocket(config)

	socket.SetErrorHandler(gs.ErrorHandler)

	socket.SetOpenHandler(gs.OpenHandler)

	socket.SetCloseHandler(gs.CloseHandler)

	socket.SetReadHandler(gs.ReadHandler)

	gs.socket = socket

	return &gs
}

func (gs *GameSocket) Run() {
	gs.socket.Logger.Println("connection socket")
	gs.socket.Go()
}

func (gs *GameSocket) OpenHandler() {
	gs.socket.Logger.Println("open socket")
}

func (gs *GameSocket) CloseHandler(rule byte, msg string) {
	gs.socket.Logger.Printf("close socket. %s\n", msg)
	gs.Done <- struct{}{}
}

func (gs *GameSocket) ErrorHandler(err error) {
	gs.socket.Logger.Printf("error socket. %s\n", err.Error())
}

func (gs *GameSocket) ReadHandler(reader io.Reader) {
	gs.socket.Logger.Println("read socket")
	gs.socket.Close()
}

func (gs *GameSocket) Send(t uint16, data []interface{}) {
}
