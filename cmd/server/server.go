package main

import (
	"log"

	"github.com/suvrick/go-kiss-core/server"
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	server := server.NewServer(logger)
	if err := server.Run(); err != nil {
		logger.Fatalln(err.Error())
	}
}
