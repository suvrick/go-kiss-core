package socket

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/types"
)

type IClientPacket interface {
	Marshal() ([]byte, error)
}

type IServerPacket interface {
	Unmarshal(r *bytes.Reader) error
}

var (
	ErrTimeoutTheGame           = errors.New("time in the game success")
	ErrConnectionNot101         = errors.New("websocket connection fail")
	ErrConnectionFail           = errors.New("connection fail")
	ErrProxyConnectionFail      = errors.New("proxy connection fail")
	ErrCloseConnectionByTimeout = errors.New("close connection by timeout")
)

// GameSock ...
type Socket struct {
	config     *SocketConfig
	messageID  uint64
	rule_close byte
	client     *websocket.Conn
	logger     *log.Logger
	proxy      *url.URL
	done       chan struct{}
	Role       byte

	openHandle  func(sender *Socket)
	closeHandle func(sender *Socket, rule byte, caption string)
	readHandle  func(sender *Socket, packetID types.PacketServerType, packet interface{})
	errorHandle func(sender *Socket, err error)
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
		logger:     config.Logger,
		done:       make(chan struct{}),
		rule_close: 255,
	}
}

func (s *Socket) Log(msg string) {
	s.logger.Println(msg)
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
			socket.errorHandle(socket, err)
			return ErrConnectionFail
		}
	}

	socket.client = client

	if resp != nil {
		if resp.StatusCode != 101 {
			socket.setClosedRule(ERROR_CONNECT_CLOSE)
			// При корректном открытии ws соединения, код должен быть 101
			if socket.errorHandle != nil {
				socket.errorHandle(socket, err)
			}
			return ErrConnectionNot101
		}
	}

	if socket.openHandle != nil {
		socket.openHandle(socket)
	}

	go socket.timeoutToGame()
	go socket.loop()
	go socket.Wait()

	return nil
}

func (socket *Socket) ConnectionWithProxy(proxy *url.URL) error {
	socket.proxy = proxy
	return socket.Connection()
}

func (socket *Socket) SetOpenHandler(handler func(sender *Socket)) {
	socket.openHandle = handler
}

func (socket *Socket) SetCloseHandler(handler func(sender *Socket, rule byte, msg string)) {
	socket.closeHandle = handler
}

func (socket *Socket) SetErrorHandler(handler func(sender *Socket, err error)) {
	socket.errorHandle = handler
}

func (socket *Socket) Send(packetID types.PacketClientType, packet interface{}) {

	socket.Log(fmt.Sprintf("[send] %#v", packet))

	if socket.client == nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, ErrConnectionFail)
		}
		return
	}

	pckt, ok := packet.(IClientPacket)

	if !ok {
		return
	}

	pack, err := pckt.Marshal()
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	// messageID
	data, err := leb128.WriteUInt64(nil, uint64(socket.messageID))

	socket.messageID++

	// packetID
	data, err = leb128.WriteUInt64(data, uint64(packetID))

	//device
	data, err = leb128.WriteByte(data, 5)

	data = append(data, pack...)

	// packet len
	data_len, err := leb128.WriteUInt64(nil, uint64(len(data)))

	data_len = append(data_len, data...)

	socket.Log(fmt.Sprintf("[send (%d)] bytes: %#v", packetID, data_len))

	err = socket.client.WriteMessage(websocket.BinaryMessage, data_len)
	if err != nil {
		socket.setClosedRule(ERROR_SEND_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
	}
}

func (socket *Socket) Close() {
	socket.setClosedRule(NORMAL_CLOSE)
	socket.close_connection()
}

func (socket *Socket) Wait() {
	defer func() {
		if socket.closeHandle != nil {
			socket.closeHandle(socket, socket.rule_close, socket.getCloseRuleMsg())
		}
	}()
	<-socket.done
}

func (socket *Socket) loop() {

	defer func() {
		if r := recover(); r != nil {
			if socket.errorHandle != nil {
				socket.errorHandle(socket, fmt.Errorf("catch recover from loop. %v", r))
			}
		}
	}()

	for socket.client != nil {

		_, msg, err := socket.client.ReadMessage()

		if err != nil {
			socket.setClosedRule(ERROR_READ_CLOSE)
			if socket.errorHandle != nil {
				socket.errorHandle(socket, err)
			}
			break
		}

		socket.read(bytes.NewReader(msg))
	}
}

func (socket *Socket) read(reader *bytes.Reader) {

	//read packetLen
	_, err := leb128.ReadUInt64(reader)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	//read packetIndex
	_, err = leb128.ReadUInt64(reader)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	// read packetID
	ID, err := leb128.ReadUInt64(reader)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	// if ID == 7 {
	// 	fmt.Printf("%#v\n", reader)
	// }

	var packet interface{}

	packetID := types.PacketServerType(ID)

	switch packetID {
	case server.LOGIN:
		packet = &server.Login{}
	case server.INFO:
		packet = &server.Info{}
	case server.BALANCE:
		packet = &server.Balance{}
	case server.CONTEST_ITEMS:
		packet = &server.ContestItems{}
	case server.BONUS:
		packet = &server.Bonus{}
	case server.REWARDS:
		packet = &server.Rewards{}
	case server.BALANCE_ITEMS:
		packet = &server.BalanceItems{}
	case server.COLLECTIONS_POINTS:
		packet = &server.CollectionsPoints{}
	case server.REWARD_GOT:
		packet = &server.RewardGot{}
	case server.BOTTLE_PLAY_DENIED:
		packet = &server.BottlePlayDenied{}
	case server.BOTTLE_ROOM:
		packet = &server.BottleRoom{}
	case server.BOTTLE_JOIN:
		packet = &server.BottleJoin{}
	case server.BOTTLE_LEAVE:
		packet = &server.BottleLeave{}
	case server.BOTTLE_LEADER:
		packet = &server.BottleLeader{}
	case server.BOTTLE_ROLL:
		packet = &server.BottleRoll{}
	case server.BOTTLE_KISS:
		packet = &server.BottleKiss{}
	case server.BOTTLE_ENTER:
		packet = &server.BottleEnter{}
	case server.CHAT_MESSAGE:
		packet = &server.ChatMessage{}
	case server.CHAT_WHISPER:
		packet = &server.ChatWhisper{}
	}

	if packet != nil {
		if pack, ok := packet.(IServerPacket); ok {
			err = pack.Unmarshal(reader)
			if err != nil {
				if socket.errorHandle != nil {
					socket.errorHandle(socket, err)
				}
				return
			}
		}

		socket.Log(fmt.Sprintf("[read] %#v", packet))
	}
}

func (socket *Socket) close_connection() {
	if socket.client != nil {
		socket.client.Close()
		socket.client = nil
		close(socket.done)
	}
}

func (socket *Socket) setClosedRule(rule byte) {
	if socket.rule_close == 255 {
		socket.rule_close = rule
	}
}

func (socket *Socket) getCloseRuleMsg() string {
	return closed_rules[socket.rule_close]
}

func (socket *Socket) timeoutToGame() {
	<-time.After(time.Duration(socket.config.TimeInTheGame) * time.Second)
	socket.setClosedRule(ERROR_TIMEOUT_CLOSE)
	if socket.errorHandle != nil {
		socket.errorHandle(socket, ErrTimeoutTheGame)
	}
}
