package client

import "github.com/suvrick/go-kiss-core/types"

const BONUS types.PacketClientType = 61

// BONUS(61) ""
type Bonus struct{}

func (bonus Bonus) String() string {
	return "BONUS(61)"
}

func (bonus *Bonus) Marshal() ([]byte, error) {
	return nil, nil
}
