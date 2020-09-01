package socket

import (
	"fmt"
	"gameserver/utils/log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)


func TestChannel_bbb(t *testing.T)  {
	c:=make(chan struct{},1)



		select {
		case  c<- struct{}{}:
			t.Error("少到")
		}



	select {

	}
}
func TestChannel_aaa(t *testing.T)  {
	c:=make(chan struct{})

	w:=sync.WaitGroup{}
	go func() {
		w.Add(1)
		select{
		case	<-c:
			t.Error("channel return")

		}
		w.Done()
		t.Error("complete goroutine1")
	}()
/*
	go func() {
		w.Add(1)
		//c <- struct{}{}
		w.Done()
		t.Error("complete goroutine2")
	}()


	go func() {
		w.Add(1)
		//c <- struct{}{}
		w.Done()
		t.Error("complete goroutine3")
	}()
*/
	time.Sleep(time.Second)


	time.AfterFunc(time.Second,func() {
		close(c)
	})
	w.Wait()

	t.Error("all complete")
}



func TestChannel(t *testing.T)  {
	testChan1()
}
var writeChan chan []byte

var closeFlag bool
var lock sync.Mutex

var closeF int32
var closeChan chan struct{}

func Destroy1() {
	if atomic.CompareAndSwapInt32(&closeF, 0, 1) {
		close(closeChan)
	}
}

func Write1(b []byte) {
	if b == nil {
		return
	}

	if atomic.LoadInt32(&closeF) == 0 {
		writeChan <- b
	}
}
func testChan1() {
	writeChan = make(chan []byte, 100)
	closeChan = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case b, ok := <-writeChan:
				if !ok {
					log.Error("!ok")
					goto BR
				}
				if b == nil {
					log.Error("b is nil")
					goto BR
				}
			case <-closeChan:
				log.Error("closeChan")
				goto BR
			}
			time.Sleep(time.Millisecond * 20)
		}
	BR:
		log.Error("break:")
	}()

	wg := sync.WaitGroup{}
	var atom int32
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {

			Write1([]byte(fmt.Sprintf("helle%d", atomic.AddInt32(&atom, 1))))

			wg.Done()
		}()
	}

	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 1000; i++ {
		go func() {
			Destroy1()
		}()
	}
	log.Error("Wait:")
	wg.Wait()
	log.Info("===全部完成===")
}
func testChan() {
	//=============================================================================================================
	writeChan = make(chan []byte, 1)

	go func() {
		for b := range writeChan {
			if b == nil {
				log.Info("b:%s", b)
				break
			}
			time.Sleep(time.Millisecond * 20)
		}

		log.Info("break:")
		lock.Lock()
		closeFlag = true
		lock.Unlock()
	}()

	wg := sync.WaitGroup{}
	var atom int32
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {

			Write([]byte(fmt.Sprintf("helle%d", atomic.AddInt32(&atom, 1))))

			wg.Done()
		}()
	}

	for i := 0; i < 10000; i++ {
		go func() {
			Destroy()
		}()
	}

	wg.Wait()

	log.Info("===全部完成===")
	//=============================================================================================================
}

func Destroy() {
	lock.Lock()
	defer lock.Unlock()

	if !closeFlag {
		close(writeChan)
		closeFlag = true
	}
}

func Write(b []byte) {
	lock.Lock()
	defer lock.Unlock()
	if closeFlag || b == nil {
		return
	}

	writeChan <- b
}
