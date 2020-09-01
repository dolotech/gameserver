package message

import (
	"sync"
)

var poolmsg sync.Pool

func init() {
	poolmsg.New = func() interface{} {
		return &Message{}
	}
}

func PoolGet() *Message {
	return poolmsg.Get().(*Message)
}

func PoolPut(p *Message) {
	p.Reset()
	p.Data = p.Data[:]
	p.Route = ""
	p.Type = 0
	p.ID = 0
	poolmsg.Put(p)
}
