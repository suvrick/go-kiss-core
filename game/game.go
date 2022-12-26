package game

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-core/packets/client"
	"github.com/suvrick/go-kiss-core/packets/server"
	"github.com/suvrick/go-kiss-core/socket"
	"github.com/suvrick/go-kiss-core/types"
)

const Tototo93 types.I = 22132982

type Game struct {
	ws        *socket.Socket
	selfID    types.I
	gameOver  chan struct{}
	closeRole CloseRole
	b         uint8
}

func NewGame() *Game {
	g := Game{}
	return &g
}

func (g *Game) SetCloseRule(rule CloseRole) {
	g.closeRole = rule
}

func (g *Game) Connection() error {

	g.gameOver = make(chan struct{})
	ws := socket.NewSocket(socket.GetDefaultSocketConfig())
	ws.SetOpenHandler(g.openHandler)
	ws.SetErrorHandler(g.errorHandler)
	ws.SetReadHandler(g.readHandler)
	ws.SetCloseHandler(g.closeHandler)

	if err := ws.Connection(); err != nil {
		return err
	}

	g.ws = ws
	g.b = 7

	return nil
}

func (g *Game) Send(packetID client.PacketClientType, packet interface{}) {
	if g.ws != nil {
		g.ws.Send(packetID, packet)
	}
}

func (g *Game) GoRoom(roomID types.B) {
	if g.ws != nil {
		g.ws.Send(client.BOTTLE_PLAY, &client.BottlePlay{
			RoomID: roomID,
		})
	}
}

func (g *Game) GoRoomToPlayer(playerID types.I) {
	if g.ws != nil {
		g.ws.Send(client.MOVE, &client.Move{
			PlayerID:  playerID,
			ByteField: 0,
		})
	}
}

func (g *Game) GameOver() chan struct{} {
	return g.gameOver
}

func (g *Game) errorHandler(game *socket.Socket, err error) {
	game.Log(fmt.Sprintf("[error] %v", err.Error()))
}

func (g *Game) readHandler(ws *socket.Socket, ID server.PacketServerType, packet interface{}) {

	ws.Log(fmt.Sprintf("[read] %T %+v", packet, packet))

	switch ID {
	case server.LOGIN:
		p := packet.(*server.Login)
		switch p.Result {
		case server.Success:
			g.selfID = p.GameID
			g.b &^= 1 // off
		default:
			g.ws.Close()
		}
	case server.INFO:
		p := packet.(*server.Info)
		if len(p.Players) > 0 && p.Players[0].GameID == g.selfID {
			g.b &^= 4 //off
			// Fill info
		}
	case server.BONUS:
		p := packet.(*server.Bonus)
		if p.CanCollect == 1 {
			ws.Send(client.BONUS, &client.Bonus{})
		}
	case server.REWARDS:
		p := packet.(*server.Rewards)
		for _, reward := range g.getRewards(p.Rewards) {
			if reward.Count > 0 {
				g.b |= 2 // on
				ws.Send(client.GAME_REWARDS_GET, &client.GameRewardsGet{
					RewardID: reward.ID,
				})
				return
			}
		}

		if (g.b & 4) == 4 {
			ws.Send(client.REQUEST, &client.Request{
				Players: []types.I{g.selfID},
				Mask:    server.INFOMASK,
			})
		}

		g.b &^= 2 // off

	case server.BOTTLE_ROOM:
		p := packet.(*server.BottleRoom)
		for _, v := range p.Players {
			if v != 0 {
				ws.Send(client.REQUEST, &client.Request{
					Players: []types.I{v},
					Mask:    server.INFOMASK,
				})
			}
		}
	case server.BOTTLE_JOIN:
		p := packet.(*server.BottleJoin)
		ws.Send(client.REQUEST, &client.Request{
			Players: []types.I{p.PlayerID},
			Mask:    server.INFOMASK,
		})
	case server.BOTTLE_LEADER:
		p := packet.(*server.BottleLeader)
		if p.LeaderID == g.selfID {
			ws.Log("I am leader!")
			<-time.After(time.Second * 5)
			ws.Log("I am rolled bottle!")
			ws.Send(client.BOTTLE_ROLL, &client.BottleRoll{
				IntField: 10,
			})
		}
	case server.BOTTLE_ROLL:
		p := packet.(*server.BottleRoll)
		if p.LeaderID == g.selfID || p.RollerID == g.selfID {

			if p.RollerID == g.selfID {
				ws.Log("I am kissed as roller!")
			} else {
				ws.Log("I am kissed as leader!")
			}

			<-time.After(time.Second * 5)

			ws.Send(client.BOTTLE_KISS, &client.BottleKiss{
				Answer: 1,
			})
		}
	case server.BOTTLE_LEAVE:
		//ws.Log("I am close game!")
		//ws.Close()
	}

	if g.b == 0 && g.closeRole == FAST {
		go func() {
			<-time.After(time.Millisecond * 500)
			g.ws.Close()
		}()
	}
}

func (g *Game) getRewards(r []server.Reward) (rewards []server.Reward) {
	for _, v := range r {
		if v.ID > 0 && v.Count > 0 {
			rewards = append(rewards, v)
		}
	}
	return
}

func (g *Game) openHandler(game *socket.Socket) {
	game.Log("[open] game socket open")
}

func (g *Game) closeHandler(game *socket.Socket, rule byte, caption string) {
	game.Log(fmt.Sprintf("[close] game socket close by %s", caption))
	close(g.gameOver)
}

type CloseRole uint8

const (
	FAST  CloseRole = 0
	NEVER CloseRole = 1
)
