package config

//关卡
type dungeon struct {
	ThroughDungeonId   int `json:"throughDungeonId"`   //推图第一关  从41001开始
	ChallengeDungeonId int `json:"challengeDungeonId"` //挑战第一关（即无限boss） 从51001开始
	ChallengeTimes     int `json:"challengeTimes"`     //挑战关卡的次数
}

//角色
type partner struct {
	ID    int `json:"id"`
	Level int `json:"level"`
	Exp   int `json:"exp"`
}

//背包
type bagItem struct {
	Count int `json:"count"` //背包里物品数量
	ID    int `json:"id"`    //背包里物品Id
}

//玩家初始信息配置
type PlayerConfig struct {
	Gold      int       `json:"gold"`    // 初始金币
	Diamond   int       `json:"diamond"` // 初始钻石
	Level     int       `json:"level"`
	Exp       int       `json:"exp"`
	Vip       int       `json:"vip"`
	Stamina   int       `json:"stamina"`   // 初始体力值
	Medal     int       `json:"medal"`     //初始功勋
	Power     int       `json:"power"`     // 初始战力
	TrainBook int       `json:"trainBook"` // 初始一千本训练书
	Honor     int       `json:"honor"`     // 10000荣誉测试
	ResType   []int     `json:"resType"`
	Partners  []partner `json:"partners"` // 初始伙伴
	Bag       []bagItem `json:"bag"`      // 初始背包
	Dungeon   dungeon   `json:"dungeon"`  //推图和挑战副本的初始ID
	//初始能量技
	UsedEnergySkill     int   `json:"usedEnergySkill"`     //当前使用的能量技
	OwnedEnergySkills   []int `json:"ownedEnergySkills"`   //拥有的能量技
	ChallengeLevelLimit int   `json:"challengeLevelLimit"` // 挑战副本开放等级
	MinFightTime        int   `json:"minFightTime"`        // 最短副本战斗时间
}

var playerConfig **PlayerConfig

func GetPlayerCfg() *PlayerConfig {
	return *playerConfig
}
