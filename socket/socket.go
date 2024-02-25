package socket

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/proxy"
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
	HiroID     uint64
	Name       string

	openHandle  func(sender *Socket)
	closeHandle func(sender *Socket, rule byte, caption string)
	errorHandle func(sender *Socket, err error)
	recvHandle  func(sender *Socket, packetID types.PacketServerType, packet any)
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
	sock := &Socket{
		config:     config,
		logger:     config.Logger,
		done:       make(chan struct{}),
		rule_close: 255,
	}

	sock.logger.SetPrefix(fmt.Sprintf("[%s] ", sock.getUID()))

	return sock
}

func (s *Socket) Log(msg string) {
	s.logger.Println(msg)
}

func (s *Socket) Logf(msg string, param ...any) {
	s.logger.Printf(msg, param...)
}

func (s *Socket) getUID() string {
	return fmt.Sprintf("%p", s)
}

func (socket *Socket) Connection() error {

	dialer := websocket.Dialer{
		HandshakeTimeout: (socket.config.ConnectTimeout),
	}

	// uid := socket.getUID()
	p := proxy.GetNetProxy(socket.Name)
	if p == nil {
		socket.setClosedRule(ERROR_CONNECT_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(socket, fmt.Errorf("get proxy fialed"))
			return ErrConnectionFail
		}
	}

	dialer.Proxy = p

	client, resp, err := dialer.Dial(socket.config.Host, socket.config.Head)

	if err != nil {
		socket.setClosedRule(ERROR_CONNECT_CLOSE)
		if socket.errorHandle != nil {
			socket.errorHandle(socket, err)
			return ErrConnectionFail
		}
	}

	socket.client = client

	socket.Logf("proxy set remote addr %s", socket.client.RemoteAddr())
	socket.Logf("proxy set local addr %s", socket.client.LocalAddr())

	// brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-158.46.166.29:hcx7fnqnph27 +
	// brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-103.241.53.114:hcx7fnqnph27
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

func (socket *Socket) SetRecvHandler(handler func(sender *Socket, packetID types.PacketServerType, packet any)) {
	socket.recvHandle = handler
}

func (socket *Socket) Send(packetID types.PacketClientType, packet interface{}) {

	if socket.client == nil {
		if socket.errorHandle != nil {
			socket.errorHandle(socket, ErrConnectionFail)
		}
		return
	}

	pckt, ok := packet.(IClientPacket)

	if !ok {
		socket.Logf("[error] %s -> %v", packet, fmt.Errorf("not impliment ICleintPacket"))
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
	data, err := leb128.WriteUInt64(nil, socket.messageID)
	if err != nil {
		socket.Logf("[error] %s -> %v", packet, err)
		return
	}

	socket.messageID++

	// packetID
	data, err = leb128.WriteUInt64(data, uint64(packetID))
	if err != nil {
		socket.Logf("[error] %s -> %v", packet, err)
		return
	}

	//device
	data, err = leb128.WriteByte(data, 5)
	if err != nil {
		socket.Logf("[error] %s -> %v", packet, err)
		return
	}

	data = append(data, pack...)

	// packet len
	data_len, err := leb128.WriteUInt64(nil, uint64(len(data)))
	if err != nil {
		socket.Logf("[error] %s -> %v", packet, err)
		return
	}

	data_len = append(data_len, data...)

	socket.Logf("%s[send] %s -> %#v", string("\033[33m"), packet, packet)

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

	//read messageID
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

	var packet IServerPacket

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
		err = packet.Unmarshal(reader)
		if err != nil {
			socket.Logf("[read] %s -> %#v", packet, err)
			// if socket.errorHandle != nil {
			// 	socket.errorHandle(socket, err)
			// }
			return
		}

		if socket.recvHandle != nil {
			socket.recvHandle(socket, packetID, packet)
		}

		socket.Logf("%s[read] %s -> %#v", string("\033[34m"), packet, packet)
	}
}

func getHex(s string) uint32 {

	hex := fnv.New32a()

	_, err := hex.Write([]byte(s))
	if err != nil {
		return 0
	}

	return hex.Sum32()
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
