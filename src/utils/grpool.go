package utils

import (
	"gameserver/utils/log"
	"runtime/debug"
	"sync"
)

var p *pool
var gronce sync.Once

type job func()
type pool struct {
	jobs chan job
}

func Go(j job) {
	gronce.Do(func() {
		p = &pool{jobs: make(chan job, 1000)}
		p.start(16)
	})
	p.jobs <- j
}
func (p *pool) start(workerNum int) {
	for i := 0; i < workerNum; i++ {
		go func() {
			for job := range p.jobs {
				func() {
					defer func() {
						if e := recover(); e != nil {
							log.Error(string(debug.Stack()),e)
						}
					}()
					job()
				}()
			}
		}()
	}
}