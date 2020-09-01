package pvb

import proto "gameserver/protocol"

//请求开始通关
type StartFightChallenge struct {
	DungeonId int `json:"dungeonId"` //关卡ID
}

//返回开始通关
type StartFightThroughed struct {
	Code int `json:"code"`
}

//请求结束通关
type FinishFightChallenge struct {
	DungeonId int `json:"dungeonId"` //关卡ID
	Win       int `json:"win"`       //是否获胜 0 1
}

//请求结束通关
type FinishFightChallenged struct {
	Code  int          `json:"code"`
	Items []proto.Item `json:"items"` //玩家获取掉落奖励
}

//购买挑战次数
type BuyChallenge struct {
	Times int `json:"times"` //次数
}

//购买挑战次数返回
type BuyChallenged struct {
	Code    int `json:"code"`    //次数
	Diamond int `json:"diamond"` //当前钻石数
	Times   int `json:"times"`   //当前次数
}

//重置无限挑战
type ResetChallenge struct {
}

//重置无限挑战返回
type ResetChallenged struct {
	Code int `json:"code"` //次数
}

//请求PVB信息
type GetPvbInfo struct {
}

//返回PVB信息
type GetPvbInfoed struct {
	Code                   int `json:"code"`
	Times                  int `json:"times"`
	ChallengeCountDownTime int `json:"challengeCountDownTime"` //挑战次数倒计时
}
