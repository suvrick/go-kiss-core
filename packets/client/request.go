package client

import (
	"sync"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const REQUEST types.PacketClientType = 8

var poolRequest = sync.Pool{
	New: func() interface{} { return Request{} },
}

// REQUEST (8) ""
type Request struct {
	Players []uint64
	Mask    int64
}

func NewRequest() *Request {
	return poolRequest.Get().(*Request)
}

func (request *Request) Down() {
	poolRequest.Put(request)
}

func (request Request) String() string {
	return "REQUEST(8)"
}

func (request *Request) Marshal() ([]byte, error) {

	data, err := leb128.WriteUInt64(nil, uint64(len(request.Players)))
	if err != nil {
		return nil, err
	}

	for _, v := range request.Players {
		data, err = leb128.WriteUInt64(data, v)
		if err != nil {
			return nil, err
		}
	}

	data, err = leb128.WriteInt64(data, request.Mask)
	if err != nil {
		return nil, err
	}

	return data, nil
}
