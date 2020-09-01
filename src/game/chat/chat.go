package chat

import (
	"gameserver/config"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"time"
)

func init() {
	msg.GetMsg().Reg(route.ChatSendMessage, &proto.ChatSendMessage{}, &proto.ChatSendMessaged{}, chat)
	msg.GetMsg().Reg(route.TipsSendMessage, &proto.TipSendMessage{}, &proto.TipSendMessaged{}, tip)
}

func chat(sess server.Session, req *proto.ChatSendMessage, resp *proto.ChatSendMessaged) {
	userdata := sess.UserData().(*online.UserData)
	now := time.Now().Unix()
	if now-userdata.ChatSendMessage < int64(config.Param().CharRest){
		resp.Code = proto.TOO_FAST_CHAT
		return
	}
	resp.Code = proto.OK
	p := userdata.GetPlayer()
	(&push.OnChatReceiveMessage{}).Push(&p, req.Message)
	userdata.ChatSendMessage = time.Now().Unix()
}

func tip(sess server.Session, req *proto.TipSendMessage, resp *proto.TipSendMessaged) {
	resp.Code = proto.OK
	(&push.OnTipsReceiveMessage{}).Push(sess.UId(), req.Message)
}
