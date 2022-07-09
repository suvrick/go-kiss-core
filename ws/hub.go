package ws

import (
	"log"

	"github.com/suvrick/go-kiss-core/frame"
)

var hub map[string]*Socket

func NewHub() {
	hub = make(map[string]*Socket)
}

type bot struct {
	bot_id  string
	packets [][]interface{}
}

func UpdateFromFrames(frames []string) []bot {

	f := frame.NewFrameDefault()
	result := make([]bot, 0)

	for _, v := range frames {

		bot_id, params, err := f.GetValue(v)
		if err != nil {
			log.Println(err)
			continue
		}

		b := bot{
			bot_id:  bot_id,
			packets: make([][]interface{}, 1),
		}

		b.packets = append(b.packets, params)

		result = append(result, result...)
	}

	return result
}
