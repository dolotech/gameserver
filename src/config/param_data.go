package config

type ParamData struct {
	StaminaNum            int `json:"staminaNum"`            //每次购买体力数量
	StaminaMax            int `json:"staminaMax"`            //体力上限
	StaminaPrice          int `json:"staminaPrice"`          //每次购买体力花费钻石
	StaminaCd             int `json:"staminaCd"`             //体力恢复的时间间隔（秒）
	NickNameGold          int `json:"nickNameGold"`          //购买昵称金币
	GiftHoner             int `json:"giftHoner"`             //赠送体力获得荣誉数
	GiftStamina           int `json:"giftStamina"`           //赠送体力数
	RecevStamina          int `json:"recevStamina"`          //领取体力数
	FriendMax             int `json:"friendMax"`             //好友最大数
	PvaPrice              int `json:"pvaPrice"`              //PVA购买次数钻石
	PvaMaxNum             int `json:"pvaMaxNum"`             //PVA最大次数
	PvbPrice              int `json:"pvbPrice"`              //PVB购买次数钻石(无尽boss购买次数花费钻石)
	PvbMaxNum             int `json:"pvbMaxNum"`             //PVB最大次数
	PvpPrice              int `json:"pvpPrice"`              //PVP购买次数钻石(周赛场购买花费)
	PvpMaxNum             int `json:"pvpMaxNum"`             //PVP最大次数
	OpenPowerSkill        int `json:"openPowerSkill"`        //能量技能开启关卡
	OpenDayTask           int `json:"openDayTask"`           //日常任务开启关卡
	OpenPva               int `json:"openPva"`               //竞技塔开启关卡
	OpenPvp               int `json:"openPvp"`               //pvp开启等级
	OpenPvb               int `json:"openPvb"`               //pvb开启等级
	OpenRole              int `json:"openRole"`              //开启人物等级
	IngotsId              int `json:"IngotsId"`              //铁锭ID
	TrainId               int `json:"TrainId"`               //训练书ID
	ScrollId              int `json:"ScrollId"`              //卷轴ID
	MinThroughDungeonId   int `json:"minThroughDungeonId"`   //初始关卡id
	MinChallengeDungeonId int `json:"minChallengeDungeonId"` //初始无尽挑战关卡id
	ChapterCap            int `json:"chapterCap"`            //每章节最大关卡数
	PvbRecover            int `json:"pvbRecover"`            //无尽挑战次数恢复时间
	PvpRecover            int `json:"pvpRecover"`            //pvp 领奖间隔时间 周赛场
	PvaRecover            int `json:"pvaRecover"`            //pva恢复时间 竞技塔

	PvaShowMax   int `json:"pvaShowMax"` // 排行榜显示人数
	PvaAwardsTime int `json:"pvaAwardsTime"` // 竞技塔排行领奖时间间隔
	CharRest int `json:"charRest"` // 聊天发送间隔(单位秒)
}

var paramData **ParamData

func Param() *ParamData {
	return *paramData
}
