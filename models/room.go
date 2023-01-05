package models

import (
	"github.com/suvrick/go-kiss-core/types"
)

type Room struct {
	RoomID           types.I
	LeaderID         types.I
	RollerID         types.I
	KissAnswerLeader KissAnswer
	KissAnswerRoller KissAnswer
	Players          map[types.I]*Player
}
