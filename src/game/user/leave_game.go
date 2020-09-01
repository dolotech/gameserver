package user

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
)

// 玩家离线

func init()  {
	online.GetEvent().AddClose(&userLogout{})
}

type userLogout struct {}

func (u *userLogout)OnClose(sess server.Session) {
	//log.Info("userLogout)OnLogin(",sess.UId())
	p := &model.Player{}
	p.PlayerId = sess.UId()
	err := p.LogoutUpdate()

	if err != nil {
		log.Error(err)
	}
}
