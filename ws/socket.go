package ws

import (
	"errors"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PCLogin struct {
	ID      uint64
	SocType uint16
	Device2 byte
	Token   string
}

type PSBalance struct {
}

type PSBonus struct {
	CanCollect byte
	Day        byte
}

type PSLogin struct {
	Result  byte
	UserID  uint64
	Balance int32
}

type PSUserInfo struct {
	Result  byte
	UserID  uint64
	Balance int32
}

// GameSock ...
type Socket struct {
	*SocketConfig

	rule_close byte
	client     *websocket.Conn
	msgID      int
	wg         sync.WaitGroup

	send_chan  chan []byte
	read_chan  chan []byte
	error_chan chan error

	openHandle  func()
	closeHandle func(byte, string)
	readHandle  func(PacketServerType, interface{})
	errorHandle func(error)
}

const (
	NORMAL_CLOSE        = 0x00
	ERROR_CONNECT_CLOSE = 0x01
	ERROR_READ_CLOSE    = 0x02
	ERROR_SEND_CLOSE    = 0x03
	ERROR_TIMEOUT_CLOSE = 0x04
)

type PacketClientType uint16

const (
	LOGIN     PacketClientType = 4
	GET_BONUS PacketClientType = 61
	BUY       PacketClientType = 6
)

type PacketServerType uint16

const (
	LOGIN_SERVER PacketServerType = 4
	INFO         PacketServerType = 5
	BALANCE      PacketServerType = 7
	BONUS        PacketServerType = 17
)

var (
	ErrTimeoutTheGame           = errors.New("time in the game success")
	ErrConnectionNot101         = errors.New("websocket connection fail")
	ErrConnectionFail           = errors.New("connection fail")
	ErrProxyConnectionFail      = errors.New("proxy connection fail")
	ErrCloseConnectionByTimeout = errors.New("close connection by timeout")
	ErrMarshalClientPacket      = errors.New("error marshal client packet")
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

func (socket *Socket) SetReadHandler(handler func(packetType PacketServerType, structure interface{})) {
	socket.readHandle = handler
}

func (socket *Socket) SetErrorHandler(handler func(err error)) {
	socket.errorHandle = handler
}

func (socket *Socket) Go() {
	socket.connect()
	socket.wg.Add(2)
	go socket.read()
	go socket.send()
	go socket.done()
}

func (socket *Socket) Send(packet []byte) {

	socket.wg.Add(1)

	go func() {
		defer socket.wg.Done()
		if socket.client == nil {
			return
		}

		socket.send_chan <- packet
	}()
}

func (socket *Socket) SendPacket(packType uint16, payload interface{}) {

	pack, err := Marshal(payload)

	if err != nil {
		if socket.errorHandle != nil {
			socket.errorHandle(ErrMarshalClientPacket)
		}
		return
	}

	data := make([]byte, 0)

	data = AppendInt(data, int64(socket.msgID)) // message ID

	data = AppendUint(data, uint64(packType)) // packet type

	data = AppendUint(data, uint64(6)) //device

	data = append(data, pack...)

	data_len := make([]byte, 0)

	data_len = AppendInt(data_len, int64(len(data))) // packet len

	data_len = append(data_len, data...)

	// socket.Logger.Printf("[SendPacket] >> %x\n", pack)

	// socket.Logger.Printf("[SendPacket2] >> %x\n", data_len)

	socket.Send(data_len)

	socket.msgID = socket.msgID + 1
}

func (socket *Socket) Close() {
	socket.setClosedRule(NORMAL_CLOSE)
	socket.close_connection()
}

func (socket *Socket) getRuleMsg() string {
	return closed_rules[socket.rule_close]
}

func (socket *Socket) connect() {

	dialer := websocket.Dialer{
		HandshakeTimeout: (socket.Timeout),
	}

	// if socket.proxy != nil {
	// 	dialer.Proxy = http.ProxyURL(&url.URL{
	// 		Scheme: socket.proxy.Scheme,
	// 		Host:   socket.proxy.Host,
	// 		User:   url.UserPassword(socket.proxy.User, socket.proxy.Password),
	// 	})
	// }

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
	if socket.errorHandle != nil {
		socket.errorHandle(ErrTimeoutTheGame)
	}
}

func (socket *Socket) done() {
	socket.wg.Wait()
	if socket.closeHandle != nil {
		socket.closeHandle(socket.rule_close, socket.getRuleMsg())
	}
	socket.close_chan()
}

func (socket *Socket) send() {
	defer func() {
		socket.wg.Done()
	}()

	var packet []byte

	for socket.client != nil {

		select {
		case packet = <-socket.send_chan:
		default:
			continue
		}

		err := socket.client.WriteMessage(websocket.BinaryMessage, packet)

		if err != nil {
			socket.setClosedRule(ERROR_SEND_CLOSE)
			if socket.errorHandle != nil {
				socket.errorHandle(err)
			}
			break
		}

	}
}

func (socket *Socket) read() {

	defer func() {
		socket.wg.Done()
	}()

	for socket.client != nil {

		_, reader, err := socket.client.NextReader()

		if err != nil {
			socket.setClosedRule(ERROR_READ_CLOSE)
			if socket.errorHandle != nil {
				socket.errorHandle(err)
			}
			break
		}

		if socket.readHandle != nil {

			ReadInt(reader, 32)
			ReadInt(reader, 32)
			t, _ := ReadUint(reader, 16)

			packetType := PacketServerType(t)

			socket.Logger.Printf("server type: %d\n", packetType)

			switch packetType {
			case LOGIN_SERVER:
				login := PSLogin{}
				if err := Unmarshal(reader, &login); err != nil {
					if socket.errorHandle != nil {
						socket.errorHandle(err)
					}
				} else {
					socket.readHandle(packetType, login)
				}

			case BALANCE:
				balance := PSBalance{}
				if err := Unmarshal(reader, &balance); err != nil {
					if socket.errorHandle != nil {
						socket.errorHandle(err)
					}
				} else {
					socket.readHandle(packetType, balance)
				}
			case BONUS:
				bonus := PSBonus{}
				if err := Unmarshal(reader, &bonus); err != nil {
					if socket.errorHandle != nil {
						socket.errorHandle(err)
					}
				} else {
					socket.readHandle(packetType, bonus)
				}
			}

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

func Unmarshal(reader io.Reader, s interface{}) error {

	var _err error = nil
	var _uint uint64 = 0
	var _int int64 = 0

	structure := reflect.ValueOf(s)

	numfield := reflect.ValueOf(s).Elem().NumField()

	for x := 0; x < numfield; x++ {

		field := structure.Elem().Field(x)

		switch reflect.ValueOf(s).Elem().Field(x).Kind() {
		case reflect.Int8:
			_int, _err = ReadInt(reader, 8)
			if _err != nil {
				return _err
			}

			field.SetInt(_int)
		case reflect.Uint8:
			_uint, _err = ReadUint(reader, 8)
			if _err != nil {
				return _err
			}

			field.SetUint(_uint)
		case reflect.Int16:
			_int, _err = ReadInt(reader, 16)
			if _err != nil {
				return _err
			}

			field.SetInt(_int)
		case reflect.Uint16:
			_uint, _err = ReadUint(reader, 32)
			if _err != nil {
				return _err
			}

			field.SetUint(_uint)
		case reflect.Int:
			_int, _err = ReadInt(reader, 32)
			if _err != nil {
				return _err
			}

			field.SetInt(_int)
		case reflect.Uint:
			_uint, _err = ReadUint(reader, 32)
			if _err != nil {
				return _err
			}

			field.SetUint(_uint)
		case reflect.Int32:
			_int, _err = ReadInt(reader, 32)
			if _err != nil {
				return _err
			}

			field.SetInt(_int)
		case reflect.Uint32:
			_uint, _err = ReadUint(reader, 32)
			if _err != nil {
				return _err
			}

			field.SetUint(_uint)
		case reflect.Int64:
			_int, _err = ReadInt(reader, 64)
			if _err != nil {
				return _err
			}

			field.SetInt(_int)
		case reflect.Uint64:
			_uint, _err = ReadUint(reader, 64)
			if _err != nil {
				return _err
			}

			field.SetUint(_uint)
		case reflect.String:
			_int, _err = ReadInt(reader, 32)
			if _err != nil {
				return _err
			}

			str := make([]byte, _int)

			reader.Read(str)

			field.SetString(string(str))
		default:
			_err = ErrMarshalClientPacket
			return _err
		}
	}

	return nil
}

func Marshal(s interface{}) ([]byte, error) {

	result := make([]byte, 0)

	var err error = nil

	structure := reflect.ValueOf(s)

	numfield := reflect.ValueOf(s).Elem().NumField()

	for x := 0; x < numfield; x++ {

		if err != nil {
			return nil, err
		}

		field := structure.Elem().Field(x)

		switch reflect.ValueOf(s).Elem().Field(x).Kind() {
		case reflect.Int8:
			if i, ok := field.Interface().(int8); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint8:
			if i, ok := field.Interface().(uint8); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int16:
			if i, ok := field.Interface().(int16); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint16:
			if i, ok := field.Interface().(uint16); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int:
			if i, ok := field.Interface().(int); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint:
			if i, ok := field.Interface().(uint); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int32:
			if i, ok := field.Interface().(int32); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint32:
			if i, ok := field.Interface().(uint32); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int64:
			if i, ok := field.Interface().(int64); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint64:
			if i, ok := field.Interface().(uint64); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.String:
			if i, ok := field.Interface().(string); ok {
				str := []byte(i)
				result = AppendInt(result, int64(len(str)))
				result = append(result, str...)
			} else {
				err = ErrMarshalClientPacket
			}
		default:
			err = ErrMarshalClientPacket
		}
	}

	return result, nil
}

// AppendUleb128 appends v to b using unsigned LEB128 encoding.
func AppendUint(b []byte, v uint64) []byte {
	for {
		//берём 7 бит
		// 	13 -> 0 0 0 0  1 1 0 1
		// 127 -> 0 1 1 1  1 1 1 1
		c := uint8(v & 0x7f)
		//сдвигаем на 7 бит
		v >>= 7
		if v != 0 {
			// 1 0 0 0 0 0 0 0
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// AppendSleb128 appends v to b using signed LEB128 encoding.
func AppendInt(b []byte, v int64) []byte {
	for {
		c := uint8(v & 0x7f) // берем первых 7 бит
		s := uint8(v & 0x40)
		v >>= 7 // сдвигайем на 7 бит вправо
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			// если вошли сюда то
			c |= 0x80 // дописываем 8 бит
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

func ReadUint(r io.Reader, n uint) (uint64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res uint64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := uint64(p[0])
		switch {
		case b < 1<<7 && b < 1<<n:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid uint")
		}
	}
}

func ReadInt(r io.Reader, n uint) (int64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res int64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := int64(p[0])
		switch {
		case b < 1<<6 && uint64(b) < uint64(1<<(n-1)):
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<6 && b < 1<<7 && uint64(b)+1<<(n-1) >= 1<<7:
			res += (1 << shift) * (b - 1<<7)
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid int")
		}
	}
}
