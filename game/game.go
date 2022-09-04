package game

import (
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/suvrick/go-kiss-core/bot"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/socket"
)

type Game struct {
	socket *socket.Socket
	msgID  int64
	Done   chan bot.Bot

	bot *bot.Bot
}

func NewGame(config *GameConfig) *Game {

	game := Game{
		Done:  make(chan bot.Bot),
		msgID: 0,
		bot: &bot.Bot{
			Time: time.Now(),
		},
	}

	ws := socket.NewSocket(config.SocketConfig)

	ws.SetErrorHandler(game.ErrorHandler)

	ws.SetOpenHandler(game.OpenHandler)

	ws.SetCloseHandler(game.CloseHandler)

	ws.SetReadHandler(game.ReadHandler)

	game.socket = ws

	return &game
}

func (game *Game) OpenHandler() {
	game.socket.Logger.Println("open")
}

func (game *Game) CloseHandler(rule byte, msg string) {
	game.socket.Logger.Printf("game over. %s\n", msg)
	game.Done <- *game.bot
}

func (game *Game) ErrorHandler(err error) {
	game.socket.Logger.Println(err.Error())
	game.GameOver()
}

func (game *Game) ReadHandler(reader io.Reader) {

	leb128.ReadInt(reader, 32)

	leb128.ReadInt(reader, 32)

	t, _ := leb128.ReadUint(reader, 16)

	packetType := server.PacketServerType(t)

	var packet interface{}
	var err error

	switch packetType {
	case server.LOGIN:
		packet, err = game.Login(reader)
	case server.INFO:
		packet, err = game.Info(reader)
	case server.BALANCE:
		packet, err = game.Balance(reader)
	case server.BONUS:
		packet, err = game.Bonus(reader)
	case server.REWARDS:
		packet, err = game.Rewards(reader)
	case server.BALANCE_ITEMS:
		packet, err = game.BalanceItems(reader)
	case server.COLLECTIONS_POINTS:
		packet, err = game.CollectionsPoints(reader)
	case server.REWARD_GOT:
		packet, err = game.RewardGot(reader)
	default:
		return
	}

	if err != nil {
		game.socket.Logger.Printf("Read [%T] %s\n", packet, err.Error())
		return
	}
}

func (game *Game) GameOver() {
	game.socket.Close()
}

func (game *Game) Send(packType client.PacketClientType, packet interface{}) {

	pack, err := leb128.Marshal(packet)
	if err != nil {
		log.Fatalln(err)
	}

	data := make([]byte, 0)

	data = leb128.AppendInt(data, int64(game.msgID)) // message ID

	data = leb128.AppendUint(data, uint64(packType)) // packet type

	data = leb128.AppendUint(data, uint64(6)) //device

	data = append(data, pack...)

	data_len := make([]byte, 0)

	data_len = leb128.AppendInt(data_len, int64(len(data))) // packet len

	data_len = append(data_len, data...)

	game.socket.Logger.Printf("Send [%T] %+v\n", packet, packet)

	game.socket.Send(data_len)

	atomic.AddInt64(&game.msgID, 1)
}

func (game *Game) Run() {
	game.socket.Go()
}
