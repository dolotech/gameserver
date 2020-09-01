package main

import (
	"gameserver/model"
	"gameserver/protocol"
	"gameserver/utils"
	"strings"
	"testing"
	"time"
)

func Test_RobotPool(t *testing.T) {

}
func Test_timer(t *testing.T) {
	heartbeat := time.Now().Unix()
	timer := time.NewTicker(time.Second)
	timeoutChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-timer.C:
				if heartbeat+5 < time.Now().Unix() {
					timer.Stop()
					return
				}
			case <-timeoutChan:
				timer.Stop()
				return
			}
		}
	}()

	time.AfterFunc(time.Second*3, func() {
		close(timeoutChan)
	})

	select {}

}
func TestIPAndPortReg(t *testing.T) {
	str := "Ddddfff"
	first := str[0:1]
	first = strings.ToLower(first)
	//fmt.Errorf("%s   %s",str, first)
	t.Error(first + str[1:])
}

func Test_Struct2Struct(t *testing.T) {

	aa := model.Player{Nickname:"michael"}

	bb := protocol.Player{}

	list := utils.StructAtoB(&bb, aa)


	t.Error(aa)
	t.Error(bb)
	t.Error(list)

}

