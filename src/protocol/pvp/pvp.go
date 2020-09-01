package pvp

import "gameserver/protocol"

//取消PVP对手查找
type CancelFindArenaOpponent struct {
}

//取消PVP对手查找返回
type CancelFindArenaOpponented struct {
	Code int `json:"code"`
}

//PVP使用能量技
type ConjureSkills struct {
	SkillId  int `json:"skillId"`  //能量技ID
	BattleId int `json:"battleId"` //战斗标识
}

//PVP使用能量技返回
type ConjureSkillsed struct {
	Code int `json:"code"`
}

//PVP对手查找
type FindArenaOpponent struct {
	Fightmodel string `json:"fightmodel"` //对战模式 pvf || pvp || no
	OpponentId uint   `json:"opponentId"` //对手ID
}

//PVP对手查找返回
type FindArenaOpponented struct {
	Code int `json:"code"`
}

//PVP资源都加载完毕
type FinishLoadBattleRes struct {
}

//PVP资源都加载完毕返回
type FinishLoadBattleResed struct {
	Code int `json:"code"`
}

//PVP设置自己的属性更改  (发起一次攻击)
type SetFighterAttribute struct {
	BattleId     int `json:"battleId"`     //战役标识
	SkillId      int `json:"skillId"`      //技能ID
	AttackRatio  int `json:"attackRatio"`  //攻击类型，0为普攻，1为爆击，2为致命
	FighterIndex int `json:"fighterIndex"` //角色，取值为1－10
}

//PVP设置自己的属性更改 (发起一次攻击)返回
type SetFighterAttributed struct {
	Code int `json:"code"`
}

//战斗中使用物品
type UseItemOnBattle struct {
	protocol.UseItem
}

//战斗中使用物品返回
type UseItemOnBattled struct {
	protocol.UseItemed
}

//START 周赛场

type EnterArena struct {
}

type EnterArenaed struct {
	PvpRank            int `json:"pvpRank"` //排位
	PkTimes            int `json:"pkTimes"` //挑战次数
	Code               int `json:"code"`
	PvpScore           int `json:"pvpScore"`           //得分
	WinNum             int `json:"winNum"`             //连胜次数
	LeftArenaAwardTime int `json:"leftArenaAwardTime"` //周赛场领奖剩余时间
}

//END
type RecvArenaRankAward struct {
}

type RecvArenaRankAwarded struct {
	Code               int `json:"code"`
	LeftArenaAwardTime int `json:"leftArenaAwardTime"` //领奖剩余时间
}

type SetFighterBuff struct {
}

type SetFighterBuffed struct {
	Code int `json:"code"`
}

type BuyPvpChallengeTimes struct {
	Times int `json:"times"`
}

type BuyPvpChallengeTimesed struct {
	Code    int `json:"code"`
	PkTimes int `json:"pkTimes"`
	Diamond int `json:"diamond"`
}
