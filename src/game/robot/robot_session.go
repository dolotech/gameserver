package robot

import (
	"gameserver/online"
	"gameserver/utils/log"
	"gameserver/utils/socket/message"
	"net"
)

type RobotSession struct {
	userData interface{}
	uid      uint
}

func (this *RobotSession) SetUid(uid uint) {
	this.uid = uid
}

func (this *RobotSession) UId() uint {
	return this.uid
}

func (this *RobotSession) Push(msg interface{}, after ...bool) {
	route := message.GetMsg().GetPushMsg(msg)
	if route == "onUpdateReadyToFight" {
		roomid := this.UserData().(*online.UserData).GetRoom()
		if roomid > 0 {
			r := Get().Get(roomid)
			r.Ready()
		} else {
			log.Error("roomid is 0")
		}
	}
}
func (this *RobotSession) WriteMsg(*message.Message) {

}
func (this *RobotSession) LocalAddr() net.Addr {
	return nil
}
func (this *RobotSession) RemoteAddr() net.Addr {
	return nil
}
func (this *RobotSession) Close() {

}
func (this *RobotSession) UserData() interface{} {
	return this.userData
}
func (this *RobotSession) SetUserData(data interface{}) {
	this.userData = data
}
