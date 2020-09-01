package push

import (
	"gameserver/protocol/route"
	msg "gameserver/utils/socket/message"
)

// 推送消息注册
func init() {
	msg.GetMsg().RegPush(route.OnUpdatePlayer, &OnUpdatePlayer{})
	msg.GetMsg().RegPush(route.OnUpdateThroughRecords, &OnUpdateThroughRecords{})
	msg.GetMsg().RegPush(route.OnUpdatePlayerTasks, &OnUpdatePlayerTasks{})
	msg.GetMsg().RegPush(route.OnUpdatePlayerFriends, &OnUpdatePlayerFriends{})
	msg.GetMsg().RegPush(route.OnUpdateEmails, &OnUpdateEmails{})
	msg.GetMsg().RegPush(route.OnTipsReceiveMessage, &OnTipsReceiveMessage{})
	msg.GetMsg().RegPush(route.OnChatReceiveMessage, &OnChatReceiveMessage{})
	msg.GetMsg().RegPush(route.OnUpdateBag, &OnUpdateBag{})
	msg.GetMsg().RegPush(route.OnUpdatePartners, &OnUpdatePartners{})
	msg.GetMsg().RegPush(route.OnUpdateAgreeMakeFriend, &OnUpdateAgreeMakeFriend{}) //同意对方加好友,好友上下线更新好友在线状态复用这个推送
	msg.GetMsg().RegPush(route.OnUpdateMakeFriend, &OnUpdateMakeFriend{})
	msg.GetMsg().RegPush(route.OnUpdateDelFriend, &OnUpdateDelFriend{})

	msg.GetMsg().RegPush(route.OnUpdateBuyDiamond, &OnUpdateBuyDiamond{})
	msg.GetMsg().RegPush(route.OnUpdatePvfFigthAlert, &OnUpdatePvfFigthAlert{})     //更新好友挑战提示
	msg.GetMsg().RegPush(route.OnUpdateReadyToFight, &OnUpdateReadyToFight{})       //更新备战状态PUSH
	msg.GetMsg().RegPush(route.OnUpdateBattleStart, &OnUpdateBattleStart{})         //更新战斗开始
	msg.GetMsg().RegPush(route.OnUpdateFightAction, &OnUpdateFightAction{})         //更新战斗动作
	msg.GetMsg().RegPush(route.OnUpdateBattleEnd, &OnUpdateBattleEnd{})             //更新战斗结束
	msg.GetMsg().RegPush(route.OnUpdateOpponentHadGone, &OnUpdateOpponentHadGone{}) //更新对手断线
	msg.GetMsg().RegPush(route.OnUpdateFindOverTime, &OnUpdateFindOverTime{})       //更新匹配对手超时
	msg.GetMsg().RegPush(route.OnUpdateUseItem, &OnUpdateUseItem{})

}
