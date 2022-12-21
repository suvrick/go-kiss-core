package server

import "github.com/suvrick/go-kiss-core/types"

const BONUS PacketServerType = 17

// BONUS(17) "BB"
type Bonus struct {
	CanCollect types.B
	Day        types.B
}
