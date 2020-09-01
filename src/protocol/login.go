package protocol

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//请求进入游戏
type EntryGame struct {
	Username  string `json:"username"`  //用户名称
	Operator  string `json:"operator"`  //平台
	Token     string `json:"token"`     //token标志 包含uid等信息
	ServerId  int    `json:"serverId"`  //进入的服务器id
	IsGetInfo int    `json:"isGetInfo"` //	（取用户信息的相关，前后端一直都是1 ANJUN）
}

//返回进入游戏
type EntryGamed struct {
	Player        Player `json:"player"`   //用户信息
	IsNew         int    `json:"isNew"`    //是否新用户
	PlayerId      uint   `json:"playerId"` //用户ID
	IsNewDay      int    `json:"isNewDay"` //是否当天第一次登录
	Code          int    `json:"code"`
	PvpTime int    `json:"pvpTime"` //周赛场距离结束时间
}
