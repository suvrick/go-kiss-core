package main

import (
	"fmt"
	"log"

	"github.com/suvrick/go-kiss-core/frame"
)

func main() {

	f := frame.NewFrameDefault()
	l, err := f.Parse("")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(l)
}
