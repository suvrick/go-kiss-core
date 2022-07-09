package ws_test

import (
	"encoding/binary"
	"fmt"
	"log"
	"testing"

	"github.com/suvrick/go-kiss-core/packets/net"
	"github.com/suvrick/go-kiss-core/ws"
)

func TestGameSocket(t *testing.T) {

	parser := net.NewParser()
	parser.Initialize()
	if parser.Error != nil {
		log.Fatalln(parser.Error.Error())
	}

	config := ws.GetDefaultGameSocketConfig()

	gs := ws.NewGameSocket(config)
	gs.Run()
	login_params := []interface{}{1000015, 32, 4, "200514254f3678c2f79cb18760ba048d", 0, ""}
	gs.Send(4, login_params)
	<-gs.Done
}

/*
func TestConnection(t *testing.T) {

	logger := log.Logger{}

	file, err := os.Create("log.txt")
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	//b := bufio.NewWriter(file)

	logger.SetOutput(io.MultiWriter(file, os.Stdout))

	parser := net.NewParser()
	parser.Initialize()
	if parser.Error != nil {
		log.Fatalln(parser.Error.Error())
	}

	login_params := []interface{}{113594657, 32, 4, "7a2b140e7b42935768c040a54ade4cfc", 0, "8c9991f3e49ef7d20d33432d1534e378"}

	pack := packets.CreateClientPacket(4, login_params...)

	if pack.Error != nil {
		t.Fatal(pack.Error.Error())
	}

	config := &ws.SocketConfig{
		Host:          "wss://bottlews.itsrealgames.com",
		Timeout:       time.Second * 30,
		TimeInTheGame: 600,
		Head: http.Header{
			"Origin": {
				"https://inspin.me",
			},
		},
		Logger: &logger,
	}

	//proxy_str := "zproxy.lum-superproxy.io:22225:lum-customer-c_07f044e7-zone-static:hcx7fnqnph27"

	//proxy := proxy.Parse(proxy_str, ":")

	socket := ws.NewSocket(config, nil)

	//log.Printf("%#v\n", proxy)

	socket.Go()
	socket.Send(&pack)

END:
	for {
		select {
		case p, ok := <-socket.NextRead():
			if !ok {
				logger.Printf("[read] -> try read from closed chanel")
				continue
			}

			err_str := ""

			if p.Error != nil {
				err_str = p.Error.Error()
			}

			logger.Printf("[read] -> %s(%d) %#v PARAMS: %v ERROR: %#v\n", p.Name, p.Type, p.Format, p.Params, err_str)

			if p.Type == 4 {
				//socket.GameOver()
				socket.Send(&pack)
				socket.Send(&pack)
				socket.Send(&pack)
				//socket.Send(&packets.Packet{})
				// socket.Send(&packets.Packet{})
			}

		case <-socket.NextError():
			logger.Println("Catch error")
			socket.GameOver()
		case rule := <-socket.Done:
			logger.Printf("Done! %s\n", socket.GetRuleClosed())
			switch rule {
			case 0x00:
				//....
			}
			socket.Close()
			break END
		}
	}

	time.Sleep(60 * time.Second)
}
*/

type Number int64
type UNumber uint64

func TestMain(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 24234)
	fmt.Print(b)
}
