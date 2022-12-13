package main

import (
	"log"

	"github.com/suvrick/go-kiss-core/client"
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	client := client.NewClient(logger)
	if err := client.Run(); err != nil {
		logger.Fatalln(err.Error())
	}

	client.LoginSend()

	<-client.Done()
}
