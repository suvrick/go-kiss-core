package client

import "github.com/suvrick/go-kiss-core/types"

const BONUS types.PacketClientType = 61

// BONUS(61) ""
type Bonus struct{}

func (bonus *Bonus) Marshal() ([]byte, error) {
	return nil, nil
}
