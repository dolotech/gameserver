package push

import (
	"gameserver/online"
	"gameserver/protocol"
)

//更新邮件信息PUSH
type OnUpdateEmail struct {
	EmailId uint   `json:"emailId"`
	Head    string `json:"head"`
	Text    string `json:"text"`
	Items   []protocol.Item
}

type OnUpdateEmails []OnUpdateEmail

func (this *OnUpdateEmails) Push(playerId uint,after ...bool) error {
	online.Get().Push(playerId, this,after...)
	return nil
}
