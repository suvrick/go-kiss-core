package game

import (
	"fmt"
	"io"
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

	login *client.Login
	bot   *bot.Bot
}

func NewGame(config *GameConfig) *Game {

	game := Game{
		Done:  make(chan bot.Bot),
		msgID: 0,
		bot: &bot.Bot{
			Time:           time.Now(),
			RewardGot:      make([]int, 0),
			BalanceHistory: make([]uint, 0),
			Log:            make([]string, 0),
			Rewards:        make([]server.Reward, 0),
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

// return game instance with deault configure
func NewGameDefault() *Game {
	config := GetDefaultGameConfig()
	return NewGame(config)
}

// func (game *Game) LiveUp() {
// 	game.bot.Live++
// 	game.bot.LiveHistory = append(game.bot.LiveHistory, game.bot.Live)
// }

// func (game *Game) LiveDown() {
// 	game.bot.Live--
// 	game.bot.LiveHistory = append(game.bot.LiveHistory, game.bot.Live)
// }

// func (game *Game) Loop() {
// 	for {
// 		<-time.After(time.Microsecond * 500)
// 		if game.bot.Live < 0 {
// 			break
// 		}
// 	}

// 	game.Log("close loop")
// 	game.GameOver()
// }

func (game *Game) LoginSend(login *client.Login) {

	if game.login == nil {
		game.login = login

		game.bot.ID = GetBotID(login)
	}

	game.Send(client.LOGIN, game.login)
}

func (game *Game) OpenHandler() {
	game.Log("open")

	//go game.Loop()
}

func (game *Game) CloseHandler(rule byte, msg string) {
	game.Logf("game over. %s", msg)
	game.Done <- *game.bot
}

func (game *Game) ErrorHandler(err error) {
	game.Logf("catch error. %s", err.Error())
	//game.GameOver()
}

func (game *Game) ReadHandler(reader io.Reader) {

	leb128.ReadInt(reader, 32)

	leb128.ReadInt(reader, 32)

	t, _ := leb128.ReadUint(reader, 16)

	packetType := server.PacketServerType(t)
	switch packetType {
	case server.LOGIN:
		game.Login(reader)
	case server.INFO:
		game.Info(reader)
	case server.BALANCE:
		game.Balance(reader)
	case server.BONUS:
		game.Bonus(reader)
	case server.REWARDS:
		game.Rewards(reader)
	case server.BALANCE_ITEMS:
		game.BalanceItems(reader)
	case server.COLLECTIONS_POINTS:
		game.CollectionsPoints(reader)
	case server.REWARD_GOT:
		game.RewardGot(reader)
	default:
		return
	}
}

func (game *Game) GameOver() {
	game.socket.Close()
}

func (game *Game) Send(packType client.PacketClientType, packet interface{}) {

	pack, err := leb128.Marshal(packet)
	if err != nil {
		game.LogErrorPacket(packet, err)
		return
	}

	data := make([]byte, 0)

	data = leb128.AppendInt(data, int64(game.msgID)) // message ID

	data = leb128.AppendUint(data, uint64(packType)) // packet type

	data = leb128.AppendUint(data, uint64(6)) //device

	data = append(data, pack...)

	data_len := make([]byte, 0)

	data_len = leb128.AppendInt(data_len, int64(len(data))) // packet len

	data_len = append(data_len, data...)

	game.LogSendPacket(packet)

	game.socket.Send(data_len)

	atomic.AddInt64(&game.msgID, 1)
}

func (game *Game) Run() {
	game.socket.Go()
}

func (game *Game) LogErrorPacket(p any, err error) {
	game.Logf("Error packet [%T] %s", p, err)
}

func (game *Game) LogReadPacket(p any) {
	s := fmt.Sprintf("Read [%T] %+v", p, p)
	game.Log(s)
}

func (game *Game) LogSendPacket(p any) {
	s := fmt.Sprintf("Send [%T] %+v", p, p)
	game.Log(s)
}

func (game *Game) Logf(s string, param ...any) {
	s = fmt.Sprintf(s, param...)
	game.bot.Log = append(game.bot.Log, s)
}

func (game *Game) Log(s string) {
	game.bot.Log = append(game.bot.Log, s)
}

func GetBotID(l *client.Login) string {
	return fmt.Sprintf("%d%d", l.NetType, l.ID)
}
