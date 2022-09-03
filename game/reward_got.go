package game

import (
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/packets/server"
)

func (game *Game) RewardGot(reader io.Reader) (interface{}, error) {
	rewardGot := &server.RewardGot{}

	err := leb128.Unmarshal(reader, rewardGot)
	if err != nil {
		return rewardGot, err
	}

	return rewardGot, nil
}
