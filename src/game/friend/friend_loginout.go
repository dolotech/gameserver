package friend

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/protocol/push"
	"gameserver/utils"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
)

func init() {
	online.GetEvent().AddClose(&FriendLoginAndLogout{})
}

type FriendLoginAndLogout struct{}

func (f *FriendLoginAndLogout) OnLogin(sess server.Session) {
	f.push(sess, 1)
}
func (f *FriendLoginAndLogout) OnClose(sess server.Session) {
	f.push(sess, 0)
}

func (f *FriendLoginAndLogout) push(sess server.Session, isonline int) {
	fs := model.Friends{}
	if err := fs.GetFriend(sess.UId()); err != nil {
		log.Error(err)
		return
	}
	for _, v := range fs {
		if online.Get().Online(v.FriendId) {
			f:=model.Friend{}
			f.Get(v.FriendId,sess.UId())
			msg := &push.OnUpdateAgreeMakeFriend{}
			msg.IsOnline = isonline
			msg.Statue = f.Status
			utils.StructAtoB(&msg.Friend, sess.UserData().(*online.UserData).GetPlayer())
			online.Get().Push(v.FriendId, msg)
		}
	}
}
