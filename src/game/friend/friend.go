package friend

import (
	"gameserver/config"
	task2 "gameserver/game/task"
	"gameserver/model"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/friend"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/utils"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"time"
)

// 消息注册 todo 不同意加好友处理
func init() {
	msg.GetMsg().Reg(route.RefreshRandomFriends, &friend.RefreshRandomFriends{}, &friend.RefreshRandomFriendsed{}, refreshRandomFriendsCb)
	msg.GetMsg().Reg(route.FindFriends, &friend.FindFriends{}, &friend.FindFriendsed{}, findFriendsCb)
	msg.GetMsg().Reg(route.MakeFriend, &friend.MakeFriend{}, &friend.MakeFriended{}, makeFriendCb)
	msg.GetMsg().Reg(route.DelFriend, &friend.DelFriend{}, &friend.DelFriended{}, delFriendCb)
	msg.GetMsg().Reg(route.AgreeMakeFriend, &friend.AgreeMakeFriend{}, &friend.AgreeMakeFriended{}, agreeMakeFriendCb)
	msg.GetMsg().Reg(route.GiveStamina, &friend.GiveStamina{}, &friend.GiveStaminaed{}, giveStaminaCb)
	msg.GetMsg().Reg(route.RecvStamina, &friend.RecvStamina{}, &friend.RecvStaminaed{}, recvStaminaCb)
}

// 领取体力
func recvStaminaCb(sess server.Session, req *friend.RecvStamina, resp *friend.RecvStaminaed) {
	resp.Code = proto.OK
	f := model.Friend{}
	err := f.Get(sess.UId(), req.FriendId)
	if err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}

	if f.Status != model.ING {
		log.Error("f.Status != model.ING", f.Status)
		resp.Code = proto.FAIL
		return
	}

	p := model.Player{PlayerId: sess.UId()}

	p.UpdateStamina(config.Param().GiftStamina)

	f.Status = model.CLOSE
	f.GiftStatus(time.Now().Unix())
	f.UpateStatus(f.Status)

	(&push.OnUpdatePlayerFriends{}).Push(sess.UId())

	p.FrienRecvTimes()
}

// 请求请求赠送体力
func giveStaminaCb(sess server.Session, req *friend.GiveStamina, resp *friend.GiveStaminaed) {
	resp.Code = proto.OK
	f := model.Friend{}
	f.Get(sess.UId(), req.FriendId)
	now := time.Now().Unix()
	giftt := int64(f.GiftTime)
	if utils.SameDay(giftt, now) {
		resp.Code = proto.FAIL
		return
	}
	err := f.UpateGiftTime()
	if err != nil {
		resp.Code = proto.FAIL
		return
	}

	p := model.Player{PlayerId: sess.UId()}
	p.GetStamina()

	if err := p.UpdateStamina(-config.Param().GiftStamina); err != nil {
		resp.Code = proto.FA_NOT_ENOUGH_STAMINA
		return
	}
	p.UpdateHonor(config.Param().GiftHoner)

	// ----------------给对方加状态---------------------
	f = model.Friend{}
	err = f.Get(req.FriendId, sess.UId())
	if err != nil {
		resp.Code = proto.FAIL
		return
	}
	f.UpateStatus(model.ING)

	if online.Get().Online(req.FriendId) {
		msg := &push.OnUpdateAgreeMakeFriend{}
		msg.IsOnline = 1
		msg.Statue = model.ING
		utils.StructAtoB(&msg.Friend, sess.UserData().(*online.UserData).GetPlayer())
		online.Get().Push(req.FriendId, msg)
	}
	//任务
	task := &task2.Task{PlayerId:p.PlayerId}
	task.UpdateFriendGiveDaily()
	task.PushTasks()
}

// 请求同意添加好友
func agreeMakeFriendCb(sess server.Session, req *friend.AgreeMakeFriend, resp *friend.AgreeMakeFriended) {
	resp.Code = proto.OK
	if req.IsAgree > 0 {
		f := model.Friend{}
		f.AgreeFriend(sess.UId(), req.FriendId)

		if online.Get().Online(req.FriendId) {
			msg := &push.OnUpdateAgreeMakeFriend{}
			utils.StructAtoB(&msg.Friend, sess.UserData().(*online.UserData).GetPlayer())
			msg.Statue = f.Status
			msg.IsOnline = 1
			online.Get().Push(req.FriendId, msg)
		}
		//任务
		tasks := task2.Task{PlayerId:req.FriendId}
		tasks.UpdatePlayerFriends()
		tasks.PushTasks()
	} else {
		f := model.Friend{}
		f.Delete(sess.UId(), req.FriendId)
	}
}

// 请求删除好友
func delFriendCb(sess server.Session, req *friend.DelFriend, resp *friend.DelFriended) {
	resp.Code = proto.OK
	f := model.Friend{}
	f.Delete(sess.UId(), req.FriendId)
	if online.Get().Online(req.FriendId) {
		online.Get().Push(req.FriendId, &push.OnUpdateDelFriend{FriendId: sess.UId(), PlayerId: req.FriendId})
	}
}

// 请求添加好友
func makeFriendCb(sess server.Session, req *friend.MakeFriend, resp *friend.MakeFriended) {
	resp.Code = proto.OK
	if req.FriendId == 0 {
		resp.Code = proto.FAIL
		return
	}

	if (&model.Friend{}).Count(sess.UId()) >= config.Param().FriendMax {
		resp.Code = proto.FA_OVER_FRIENDS_COUNT
		return
	}

	if (&model.Friend{}).Count(req.FriendId) >= config.Param().FriendMax {
		resp.Code = proto.FA_OVER_FRIENDS_COUNT
		return
	}

	// 判断对方玩家是否存在
	p := &model.Player{PlayerId: req.FriendId}
	if err := p.GetById(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		log.Error("请求查找好友，没有找到")
		return
	}
	f := &model.Friend{}
	if err := f.ApplyFriend(sess.UId(), p.PlayerId); err == nil {
		// 如果对方在线，推送加要有通知
		if online.Get().Online(req.FriendId) {
			msg := &push.OnUpdateMakeFriend{}
			utils.StructAtoB(&msg.Friend, sess.UserData().(*online.UserData).GetPlayer())
			online.Get().Push(req.FriendId, msg)
			//log.Error("req.FriendId",req.FriendId)
			(&push.OnUpdatePlayerFriends{}).Push(req.FriendId)
		}
	}
}

// 请求推荐好友
func refreshRandomFriendsCb(sess server.Session, req *friend.RefreshRandomFriends, resp *friend.RefreshRandomFriendsed) {

	userdata := sess.UserData().(*online.UserData)
	now := time.Now().Unix()
	if now-userdata.RefreshRandomFriends < 6 {
		resp.Code = proto.FAIL
		return
	}

	if (&model.Friend{}).Count(sess.UId()) >= config.Param().FriendMax {
		resp.Code = proto.FAIL
		return
	}

	userdata.RefreshRandomFriends = time.Now().Unix()

	resp.Code = proto.OK
	selfFriend := model.Friends{}
	err := selfFriend.Get(sess.UId())
	if err != nil {
		log.Error(err)
	}
	//log.Error("请求推荐好友")
	ids := make([]uint, 0, len(selfFriend)+1)
	ids = append(ids, sess.UId())
	if len(selfFriend) > 0 {
		for _, v := range selfFriend {
			ids = append(ids, v.FriendId)
		}
	}
	ps := model.Players{}
	//log.Error("请求推荐好友", len(ps), ids)
	err = ps.GetNotIn(ids)
	if err != nil {
		log.Error(err)
		return
	}
	//log.Error("请求推荐好友1", len(ps))
	resp.Friends = make([]friend.Friend, len(ps))
	for k, v := range ps {
		utils.StructAtoB(&resp.Friends[k], v)
		isOn := online.Get().Online(v.PlayerId)
		if isOn {
			resp.Friends[k].IsOnline = 1
		}
	}
}

// 请求查找好友
func findFriendsCb(sess server.Session, req *friend.FindFriends, resp *friend.FindFriendsed) {
	resp.Code = proto.OK
	ps := model.Players{}
	if err := ps.GetByNick(sess.UId(), req.Nickname); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		log.Error("请求查找好友，没有找到")
		return
	}
	resp.Friends = make([]friend.Friend, len(ps))
	for k, v := range ps {
		utils.StructAtoB(&resp.Friends[k], v)
		if online.Get().Online(v.PlayerId) {
			resp.Friends[k].IsOnline = 1
		}
	}
}
