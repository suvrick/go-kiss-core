package game

import (
	"fmt"
	"io"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/suvrick/go-kiss-core/bot"
	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/socket"
)

type Game struct {
	socket *socket.Socket
	msgID  int64
	Done   chan bot.Bot

	/* packets */
	login *client.Login
	buy   *client.Buy

	bot *bot.Bot
}

func NewGame(config *GameConfig) *Game {

	game := Game{
		Done:  make(chan bot.Bot),
		msgID: 0,
		bot: &bot.Bot{
			Time:           time.Now(),
			RewardGot:      make([]int, 0),
			BalanceHistory: make([]uint, 0),
			ErrorHistory:   make([]string, 0),
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

func (g *Game) Connection() error {
	return g.socket.Connection()
}

func (g *Game) ConnectionWithProxy(proxy *url.URL) error {
	return g.socket.Connection()
}

func (game *Game) LoginSend(login *client.Login) {

	if game.login == nil {
		game.login = login

		game.bot.ID = GetBotID(login)
	}

	game.Send(client.LOGIN, game.login)
}

func (game *Game) SetBuyPacket(buy *client.Buy) {
	game.buy = buy
}

func (game *Game) BuySend() {
	if game.buy != nil {
		game.Send(client.BUY, game.buy)
	}
}

func (game *Game) OpenHandler() {
	game.Log("open connection.")
}

func (game *Game) CloseHandler(rule byte, msg string) {
	game.Logf("game over. %s", msg)
	game.Done <- *game.bot
}

func (game *Game) ErrorHandler(err error) {
	err_str := fmt.Sprintf("catch error. %s", err.Error())
	game.Logf(err_str)
	game.bot.ErrorHistory = append(game.bot.ErrorHistory, err_str)
}

func (game *Game) ReadHandler(reader io.Reader) {

	//read packetLen
	_, err := leb128.ReadUint(reader, 32)
	if err != nil {
		game.Logf(err.Error())
		return
	}

	//read packetIndex
	_, err = leb128.ReadUint(reader, 32)
	if err != nil {
		game.Logf(err.Error())
		return
	}

	// read packetID
	packetID, err := leb128.ReadUint(reader, 32)
	if err != nil {
		game.Logf(err.Error())
		return
	}

	packetType := server.PacketServerType(packetID)
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
	data = leb128.AppendInt(data, int64(game.msgID)) // messageID
	data = leb128.AppendUint(data, uint64(packType)) // packetID
	data = leb128.AppendUint(data, uint64(5))        //device
	data = append(data, pack...)

	data_len := make([]byte, 0)
	data_len = leb128.AppendInt(data_len, int64(len(data))) // packet len
	data_len = append(data_len, data...)

	game.LogSendPacket(packet)
	game.socket.Send(data_len)
	atomic.AddInt64(&game.msgID, 1)
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
