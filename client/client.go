package client

import (
	"context"
	"log"

	pb "github.com/suvrick/go-kiss-core/protos"
	"google.golang.org/grpc"
)

type client struct {
	log *log.Logger
	pb.GameClient
	game pb.GameClient
	done chan struct{}
}

func NewClient(logger *log.Logger) *client {
	return &client{
		log:  logger,
		done: make(chan struct{}),
	}
}

func (c *client) Done() chan struct{} {
	return c.done
}

func (c *client) Run() error {
	conn, err := grpc.Dial(":55555", grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.game = pb.NewGameClient(conn)
	return nil
}

func (c *client) LoginS(ctx context.Context, in *pb.ServerLogin) error {
	return nil
}

func (c *client) LoginSend() {
	c.game.Login(context.Background(), &pb.ClientLogin{
		LoginId:    103786258,
		NetType:    32,
		DeviceType: 5,
		Key:        "dc93c8e0c365ca792cf1198ab71c73e7",
	})
}

// func (c *client) Login(p *pb.Server_Login) {
// 	c.log.Printf("%d, %d, %d\n", p.Result, p.GameId, p.Balance)
// }

// func (c *client) LoginSend() {
// 	p := &pb.Client{
// 		Event: &pb.Client_ClientLogin{
// 			ClientLogin: &pb.Client_Login{
// 				LoginId:    103786258,
// 				NetType:    32,
// 				DeviceType: 5,
// 				Key:        "dc93c8e0c365ca792cf1198ab71c73e7",
// 			},
// 		},
// 	}
// 	c.stream.Send(p)
// }

// func (c *client) ShutdownSend() {
// 	p := &pb.Client{
// 		Event: &pb.Client_ClientShutdown{
// 			ClientShutdown: &pb.Client_Shutdown{},
// 		},
// 	}
// 	c.stream.Send(p)
// }
