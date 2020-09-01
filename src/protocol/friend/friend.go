package friend

//同意添加好友
type AgreeMakeFriend struct {
	IsAgree  int  `json:"isAgree"`
	FriendId uint `json:"friendId"` //
}

//同意添加好友返回
type AgreeMakeFriended struct {
	Code int `json:"code"`
}

//删除好友
type DelFriend struct {
	FriendId uint `json:"friendId"` //
}

//删除好友返回
type DelFriended struct {
	Code int `json:"code"`
}

//查找好友
type FindFriends struct {
	Nickname string `json:"nickname"`
}

//查找好友返回
type FindFriendsed struct {
	Code    int      `json:"code"`
	Friends []Friend `json:"friends"`
}

//请求赠送体力
type GiveStamina struct {
	FriendId uint `json:"friendId"` //
}

//请求赠送体力返回
type GiveStaminaed struct {
	Code int `json:"code"`
}

//添加好友
type MakeFriend struct {
	FriendId uint `json:"friendId"`
}

//添加好友返回
type MakeFriended struct {
	Code int `json:"code"`
}

//请求回复赠送体力
type RecvStamina struct {
	FriendId uint `json:"friendId"` //
}

//请求回复赠送体力返回
type RecvStaminaed struct {
	Code int `json:"code"`
} //推荐好友
type RefreshRandomFriends struct{}

//好友信息
type Friend struct {
	Level    int    `json:"level"`    //等级
	Power    int    `json:"power"`    //战力
	PlayerId uint   `json:"playerId"` //ID
	Avatar   int    `json:"avatar"`   //头像
	Statue   int    `json:"statue"`
	Nickname string `json:"nickname"` //昵称
	Vip      int    `json:"vip"`      //VIP登级
	IsOnline int    `json:"isOnline"` //是否在线
	ResType  int    `json:"resType"`  //性别
}

//推荐好友返回
type RefreshRandomFriendsed struct {
	Code    int      `json:"code"`
	Friends []Friend `json:"friends"`
}
