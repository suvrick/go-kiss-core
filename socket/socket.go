package socket

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/types"
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
	config      *SocketConfig
	packetIndex uint64
	rule_close  byte
	client      *websocket.Conn
	logger      *log.Logger
	proxy       *url.URL
	done        chan struct{}
	Role        byte
	Self        *models.Bot

	//openHandle  func(sender *Socket)
	updateSelfHandle func(sender *Socket, self *models.Bot)

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
		Self: &models.Bot{
			Room: &models.Room{},
			Info: &models.Player{},
		},
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

func (socket *Socket) SetUpdateSelfHandler(handler func(sender *Socket, self *models.Bot)) {
	socket.updateSelfHandle = handler
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

	pack, err := leb128.Marshal(packet)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	data := make([]byte, 0)
	data = leb128.AppendInt(data, int64(socket.packetIndex)) // messageID
	data = leb128.AppendUint(data, uint64(packetID))         // packetID
	data = leb128.AppendUint(data, uint64(5))                //device
	data = append(data, pack...)

	data_len := make([]byte, 0)
	data_len = leb128.AppendInt(data_len, int64(len(data))) // packet len
	data_len = append(data_len, data...)

	err = socket.client.WriteMessage(websocket.BinaryMessage, data_len)
	if err != nil {
		socket.setClosedRule(ERROR_SEND_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
	}

	socket.packetIndex++
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

func (socket *Socket) read(reader io.Reader) {

	//read packetLen
	_, err := leb128.ReadUint(reader, 32)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	//read packetIndex
	_, err = leb128.ReadUint(reader, 32)
	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
		}
		return
	}

	// read packetID
	ID, err := leb128.ReadUint(reader, 32)
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
		len, err := leb128.ReadUint(reader, 16)
		if err == nil {
			msg := make([]byte, len)
			reader.Read(msg)
			mask, err := leb128.ReadInt(reader, 64)
			if err == nil && types.I(mask) == server.INFOMASK {
				reader = bytes.NewReader(msg)
				packet = &server.Info{}
			}
		}
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

		err = leb128.Unmarshal(reader, packet)

		if err != nil {
			if socket.errorHandle != nil {
				socket.errorHandle(socket, err)
			}
			return
		}

		socket.Log(fmt.Sprintf("[read] %#v", packet))

		if pack, ok := packet.(interfaces.IServerPacket); ok {
			err = pack.Use(socket.Self, socket)
			if err != nil {
				if socket.errorHandle != nil {
					socket.errorHandle(socket, err)
				}
				return
			}
		}
	}
}

func (socket *Socket) close_connection() {
	if socket.client != nil {
		socket.client.Close()
		socket.client = nil
	}

	close(socket.done)
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
	socket.close_connection()
}

func (socket *Socket) UpdateSelfEmit() {
	if socket.updateSelfHandle != nil {
		socket.updateSelfHandle(socket, socket.Self)
	}
}
