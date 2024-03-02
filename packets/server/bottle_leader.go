package server

import (
	"bytes"
	"sync"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const BOTTLE_LEADER types.PacketServerType = 28

var poolBottleLeader = sync.Pool{
	New: func() any { return &BottleLeader{} },
}

func NewBottleLeader() *BottleLeader {
	return poolBottleLeader.Get().(*BottleLeader)
}

func (p *BottleLeader) Reset() {
	poolBottleLeader.Put(p)
}

// "I"
type BottleLeader struct {
	LeaderID uint64
}

func (p BottleLeader) String() string {
	return "BOTTLE_LEADER(28)"
}

func (bottleLeader *BottleLeader) Unmarshal(r *bytes.Reader) error {
	var err error
	bottleLeader.LeaderID, err = leb128.ReadUInt64(r)
	return err
}
