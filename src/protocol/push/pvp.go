package push

import (
	"gameserver/protocol"
	"gameserver/protocol/friend"
)

//更新匹配对手超时 PUSH
type OnUpdateFindOverTime struct {
	PlayerId uint `json:"playerId"` //
}

//更新战斗结束PUSH
type OnUpdateBattleEnd struct {
	WinId  uint `json:"winId"`  //赢家ID
	IsGone int  `json:"isGone"` //角斗士是否离开（逃走） 1：离开，0：未离开
}

type BattlePlayer struct {
	PlayerId uint `json:"playerId"`
	Hp       int  `json:"hp"`
	Dp       int  `json:"dp"`
}

//更新战斗开始PUSH
type OnUpdateBattleStart struct {
	Left     BattlePlayer `json:"left"`     // 左边玩家ID
	Right    BattlePlayer `json:"right"`    // 右边玩家ID
	BattleId int          `json:"battleId"` //战斗标识
}

//更新战斗动作 PUSH
type OnUpdateFightAction struct {
	SkillId    int        `json:"skillId"`    //技能ID
	FightIndex int        `json:"fightIndex"` //攻击者角色，取值为1－10
	Defenders  []Defender `json:"defenders"`  //被攻击者们
}

type Defender struct {
	Dp          int `json:"dp"`          // 防御的变化值，不管是加还减都是正数
	Lethal      int `json:"lethal"`      // 致死，0为没死，1为死了
	AttackRatio int `json:"attackRatio"` // 攻击类型，0为普攻，1为爆击，2为致命
	Hp          int `json:"hp"`          // 血量的变化值，不管是加还减都是正数
	FightIndex  int `json:"fightIndex"`  //被攻击者角色，取值为1－10
}

//更新好友挑战提示PUSH
type OnUpdatePvfFigthAlert struct {
	friend.Friend
}

//对手信息
type Opponent struct {
	Level                 int               `json:"level"`                 //等级
	Power                 int               `json:"power"`                 //战斗力
	PlayerId              uint              `json:"playerId"`              //玩家ID
	Avatar                int               `json:"avatar"`                //头像
	ContinuousArenaPKWins int               `json:"continuousArenaPKWins"` //竞技场连赢次数
	PvpRank               int               `json:"pvpRank"`               //arena段位
	Partners              protocol.Partners `json:"partners"`              //所有的角色
	UsedEnergySkill       int               `json:"usedEnergySkill"`       //当前使用的能量技
	ArenaPKWins           int               `json:"arenaPKWins"`           //竞技场赢次数
	ArenaPKLoses          int               `json:"arenaPKLoses"`          //竞技场输次数
	Vip                   int               `json:"vip"`                   //vip等级
	Nickname              string            `json:"nickname"`              //昵称
	ResType               int               `json:"resType"`               // 男女，1为男，2为女 默认为男
}

//更新备战状态PUSH
type OnUpdateReadyToFight struct {
	PlayerId uint     `json:"playerId"` //玩家ID
	MapId    int      `json:"mapId"`    //产生一个随机消除板，取值为1-100的整数  var mapId = Random.randRange(1, 100); (ANJUN)
	Opponent Opponent `json:"opponent"` //对手信息
}

//更新对手断线 PUSH
type OnUpdateOpponentHadGone struct {
	PlayerId uint `json:"playerId"` //
	FriendId uint `json:"friendId"` //
}

type OnUpdateUseItem struct {
	ItemId int `json:"itemId"` //物品标识
}
