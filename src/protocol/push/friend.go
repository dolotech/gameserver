package push

import "gameserver/protocol/friend"

//通知谁删掉了好友 PUSH
type OnUpdateDelFriend struct {
	PlayerId uint `json:"playerId"` //
	FriendId uint `json:"friendId"` //
}
//同意对方加好友PUSH
type OnUpdateAgreeMakeFriend struct {
	friend.Friend
}
type OnUpdateMakeFriend struct {
	friend.Friend
}
