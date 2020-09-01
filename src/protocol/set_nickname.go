package protocol

//设置昵称
type SetNickname struct {
	Nickname   string `json:"nickname"`   //用户昵称
	ResType    int    `json:"resType"`    //性别 1:男; 2:女
	IsModified int    `json:"isModified"` //是否修改昵称 1:修改昵称;0:不修改
}

//设置昵称返回
type SetNicknamed struct {
	Code   int    `json:"code"`
	Player Player `json:"player"` //玩家信息
}

type SetAvatar struct {
	Avatar int `json:"avatar"`
}

type SetAvatared struct {
	Code   int `json:"code"`
	Avatar int `json:"avatar"`
}
