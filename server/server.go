package server

import (
	"context"
	"log"
	"net"

	pb "github.com/suvrick/go-kiss-core/protos"
	"google.golang.org/grpc"
)

type server struct {
	log *log.Logger
	pb.UnimplementedGameServer
}

func NewServer(l *log.Logger) *server {
	return &server{
		log: l,
	}
}

func (s *server) Run() error {

	s.log.Println("start game server")

	rpc_server := grpc.NewServer()

	pb.RegisterGameServer(rpc_server, s)

	l, err := net.Listen("tcp", ":55555")

	if err != nil {
		s.log.Fatalln(err.Error())
	}

	return rpc_server.Serve(l)
}

func (s *server) Login(ctx context.Context,p *pb.ClientLogin) (*pb.ServerResponse, error)
	s.LoginS(context.Background(), &pb.ServerLogin{
		Result:  0,
		GameId:  123,
		Balance: 555,
	})
	return nil
}
