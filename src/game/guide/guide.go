package guide

import (
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
)

// 游戏消息注册
func init() {
	msg.GetMsg().Reg(route.Setnewstep, &proto.SetNewStep{}, &proto.SetNewSteped{}, SetNewStepCb)
}

// 请求主线引导记录
func SetNewStepCb(sess server.Session, req *proto.SetNewStep, resp *proto.SetNewSteped) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.GetById();err != nil{
		log.Error(err)
		sess.Close()
		return
	}
	if err := player.UpdateNewStep(req.NewStep);err != nil{
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}
}