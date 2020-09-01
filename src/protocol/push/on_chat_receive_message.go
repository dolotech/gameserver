package push

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/utils"
	"gameserver/utils/filter"
)

//更新聊天信息PUSH,全服聊天广播
type OnChatReceiveMessage struct {
	Level    int    `json:"level"`    //等级
	Nickname string `json:"nickname"` //昵称
	PlayerId uint   `json:"playerId"` //Id
	Avatar   int    `json:"avatar"`   //头像
	Message  string `json:"message"`  //信息
	ResType  int    `json:"resType"`  //性别
}

// player:聊天信息发送者
// msg发送的:聊天内容
func (this *OnChatReceiveMessage) Push(player *model.Player, msg string) error {
	utils.StructAtoB(this, player)
	str := []rune(msg)
	filter.FilterText("dic.txt", str, []rune{}, '*')
	this.Message = string(str)
	online.Get().BCAll(this)
	return nil
}
