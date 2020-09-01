package push

import "gameserver/online"

//更新走马灯信息PUSH
type OnTipsReceiveMessage struct {
	PlayerId uint   `json:"playerId"` //Id
	Message  string `json:"message"`  //信息
}

func (this *OnTipsReceiveMessage) Push(playerid uint, msg string) error {
	this.PlayerId = playerid
	this.Message = msg
	online.Get().BCAll(this)
	return nil
}
