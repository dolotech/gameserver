package utils

import "math"

type AutoInc struct {
	queue chan uint32
	i     uint32
}

func NewID() (ai *AutoInc) {
	ai = &AutoInc{
		queue: make(chan uint32, 1000),
	}
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for {
		ai.i += 1
		if ai.i >= math.MaxUint32 {
			ai.i = 1
		}
		ai.queue <- ai.i
	}
}

func (ai *AutoInc) Id() uint32 {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	close(ai.queue)
}
