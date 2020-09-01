package model

import (
	"errors"
	"gameserver/config"
	"gameserver/utils/db"
	"github.com/jinzhu/gorm"
	"time"
)

// 玩家数据
type Player struct {
	// 玩家基础数据
	PlayerId        uint   `gorm:"primary_key;column:playerId;COMMENT:'用户游戏唯一标识';type:int(11) unsigned " json:"playerId"`
	Nickname        string `gorm:"type:varchar(30);column:nickname;index:index_nickname" json:"nickname"`
	ResType         int    `gorm:"type:int(4);column:resType;default:0;COMMENT:'男女，1为男，2为女 默认为男'" json:"resType"`
	Diamond         int    `gorm:"type:int(11) unsigned;column:diamond;default:0" json:"diamond"`
	Gold            int    `gorm:"type:int(11) unsigned;column:gold;default:0" json:"gold"`
	Honor           int    `gorm:"type:int(11) unsigned;column:honor;default:0" json:"honor"`
	Stamina         int    `gorm:"type:int(11) unsigned;column:stamina;default:0;COMMENT:''" json:"stamina"`
	Power           int    `gorm:"type:int(8) unsigned ;column:power;default:0;COMMENT:''" json:"power"`
	Exp             int    `gorm:"type:int(11) unsigned;column:exp;default:0" json:"exp"`
	Level           int    `gorm:"type:int(4) unsigned ;column:level;default:0;COMMENT:''" json:"level"`
	Vip             int    `gorm:"type:int(4) unsigned;column:vip;default:0;COMMENT:''" json:"vip"`
	Avatar          int    `gorm:"type:int(4);column:avatar;default:0;COMMENT:'头像，0为默认头像'" json:"avatar"`
	Robot           int    `gorm:"type:int(4);column:robot;default:0;COMMENT:'0真是玩家,1是机器人'" json:"robot"`
	LoginTime       int    `gorm:"type:int(11);column:loginTime;COMMENT:''" json:"loginTime"`
	LeaveTime       int    `gorm:"type:int(11);column:leaveTime;COMMENT:''" json:"leaveTime"`
	UsedEnergySkill int    `gorm:"type:int(11);column:usedEnergySkill;default:0;COMMENT:'当前使用的能量技，在initByConfig或者initByDB中初始化'" json:"usedEnergySkill"`
	// 通关数据
	CurThroughId int    `gorm:"type:int(11);column:curThroughId;default:0;COMMENT:'当前挑战的关卡索引'" json:"curThroughId"`
	IsGuide      int    `gorm:"type:int(4);column:isGuide;default:0;COMMENT:'0开 1关闭 新手引导的开关'" json:"isGuide"`
	NewStep      int    `gorm:"type:int(4);column:newStep;default:0;COMMENT:'主线引导步骤'" json:"newStep"`
	GuideSteps   string `gorm:"type:varchar(30);column:guideSteps;COMMENT:'异步引导步骤Battle1@0&Battle2@0  0代表没执行过1代表执行'" json:"guideSteps"`

	ClearDebuffCount      int `gorm:"type:int(11);column:clearDebuffCount;default:0" json:"clearDebuffCount"`
	ClearIconCount        int `gorm:"type:int(11);column:clearIconCount;default:0" json:"clearIconCount"`
	DayOnlineTotalTime    int `gorm:"type:int(11);column:dayOnlineTotalTime;default:0;COMMENT:''" json:"dayOnlineTotalTime"`
	StaminaRecvTimes      int `gorm:"type:int(11);column:staminaRecvTimes;default:0;COMMENT:'每天领取好友赠送的体力的次数限定最大好友数量默认30'" json:"staminaRecvTimes"`
	LastRecvStaminaTime   int `gorm:"type:int(11);column:lastRecvStaminaTime" json:"lastRecvStaminaTime"`
	RecvLevelAwardIndex   int `gorm:"type:int(11);column:recvLevelAwardIndex;default:0;COMMENT:''" json:"recvLevelAwardIndex"`
	RecvOnlineAwardIndex  int `gorm:"type:int(11);column:recvOnlineAwardIndex;default:0;COMMENT:''" json:"recvOnlineAwardIndex"`
	RecvOpenSvrAwardIndex int `gorm:"type:int(11);column:recvOpenSvrAwardIndex;default:0;COMMENT:''" json:"recvOpenSvrAwardIndex"`
	RecvStarIndex         int `gorm:"type:int(11);column:recvStarIndex;default:0;COMMENT:''" json:"recvStarIndex"`
	LoginTotalDays        int `gorm:"type:int(11);column:loginTotalDays;default:0" json:"loginTotalDays"`
	TodaySignTime         int `gorm:"type:int(11);column:todaySignTime;COMMENT:''" json:"todaySignTime"`
	TotalHarmToMonster    int `gorm:"type:int(11);column:totalHarmToMonster;default:0;COMMENT:''" json:"totalHarmToMonster"`
	TotalOnlineTime       int `gorm:"type:int(11);column:totalOnlineTime;COMMENT:''" json:"totalOnlineTime"`
	MonthSignDays         int `gorm:"type:int(11);column:monthSignDays;default:0" json:"monthSignDays"`

	//无尽挑战数据
	ChallengeTimes    int `gorm:"type:int(11);column:challengeTimes;default:0" json:"challengeTimes"`
	MaxChallengeId    int `gorm:"type:int(11);column:maxChallengeId" json:"maxChallengeId"`
	LastChallengeTime int `gorm:"type:int(11);column:lastChallengeTime;default:0" json:"lastChallengeTime"`

	// 周赛场数据
	PvpRank             int `gorm:"type:int(11);column:pvpRank;default:0;COMMENT:''" json:"pvpRank"`
	PvpScore            int `gorm:"type:int(11);column:pvpScore;default:0;COMMENT:''" json:"pvpScore"`
	ArenaPKLoses        int `gorm:"type:int(11);column:arenaPKLoses;default:0;COMMENT:''" json:"arenaPKLoses"`
	ArenaPKWins         int `gorm:"type:int(11);column:arenaPKWins;default:0;COMMENT:''" json:"arenaPKWins"`
	LastArenaAwardIndex int `gorm:"type:int(11);column:lastArenaAwardIndex;default:0" json:"lastArenaAwardIndex"`
	LastArenaAwardTime  int `gorm:"type:int(11);column:lastArenaAwardTime" json:"lastArenaAwardTime"`
	ContinuousArenaPKWins int `gorm:"type:int(11);column:continuousArenaPKWins;default:0;COMMENT:'pvp连赢次数'" json:"continuousArenaPKWins"`
	DayArenaPKTimes       int `gorm:"type:int(11);column:dayArenaPKTimes;default:0;COMMENT:'pvp pk次数'" json:"dayArenaPKTimes"`

	//竞技塔数据
	WeeklyPKLoses          int `gorm:"type:int(11);column:weeklyPKLoses;default:0;COMMENT:''" json:"weeklyPKLoses"`
	WeeklyPKWins           int `gorm:"type:int(11);column:weeklyPKWins;default:0;COMMENT:''" json:"weeklyPKWins"`
	WeeklyRank             int `gorm:"type:int(11);column:weeklyRank;index:weeklyRank;default:0;COMMENT:'竞技塔,这个值只是镜像'" json:"weeklyRank"`
	ContinuousWeeklyPKWins int `gorm:"type:int(11);column:continuousWeeklyPKWins;default:0;COMMENT:'竞技塔连赢次数'" json:"continuousWeeklyPKWins"`
	BuyWeeklyPKTimes       int `gorm:"type:int(11);column:buyWeeklyPKTimes;default:0;COMMENT:''" json:"buyWeeklyPKTimes"`
	DayWeeklyPKTimes       int `gorm:"type:int(11);column:dayWeeklyPKTimes;default:0;COMMENT:'玩家竞技塔剩余pk次数'" json:"dayWeeklyPKTimes"`
	WeeklyLeftTime         int `gorm:"type:int(11);column:weeklyLeftTime;COMMENT:'2小时自增竞技塔pk次数领取时间'" json:"weeklyLeftTime"`
	LastWeeklyRank         int `gorm:"type:int(11);column:lastWeeklyRank;default:0" json:"lastWeeklyRank"`
	LeftWeeklyAwardTime int `gorm:"type:int(11);column:leftWeeklyAwardTime;default:0;COMMENT:'竞技塔领奖剩余时间'" json:"leftWeeklyAwardTime"` //领奖剩余时间
}

type Players []*Player

func (this *Player) GetById() error {
	return db.Get().Model(this).Limit(1).Find(this).Error
}

func (this *Player) AddPlayer() error {
	return db.Get().Create(this).Error
}

func (this *Player) Get() error {
	return db.Get().Model(this).Limit(1).Find(this).Error
}

func (this *Player) Select(fields string) error {
	return db.Get().Model(this).Select(fields).Limit(1).Find(this).Error
}

func (this *Player) Exist() bool {
	return db.Get().Model(this).Select("playerId").Find(this).Error == nil
}

// 获取指定用户的钻石和金币
func (this *Player) GetDiamondGoldHonor() error {
	return db.Get().Model(this).Select("diamond,gold,honor").Find(this).Error
}

func (this *Player) GetDiamond() error {
	return db.Get().Model(this).Select("diamond").Find(this).Error
}

func (this *Player) GetGold() error {
	return db.Get().Model(this).Select("gold").Find(this).Error
}

// 更新指定用户的昵称和性别
func (this *Player) UpdateNicknameSex(nick string, sex int) error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"nickname": nick,
		"resType":  sex,
	}).Error
}

func (this *Player) UpdateAvatar(avatar int) error {
	return db.Get().Model(this).Update("avatar", avatar).Error
}

// 登陆更新玩家数据
func (this *Player) LoginUpdate() error {
	return db.Get().Model(this).Update("loginTime", time.Now().Unix()).Error
}

// 离线更新玩家数据
func (this *Player) LogoutUpdate() error {
	return db.Get().Model(this).Update("leaveTime", time.Now().Unix()).Error
}

var enoughStamina = errors.New("not enough stamina")

// stamina  如果扣除传负值
func (this *Player) UpdateStamina(stamina int) error {
	if this.Stamina == 0 {
		this.GetStamina()
	}
	staminaAdd, staminaTime := this.GetAddStamina(this.Stamina, this.LastRecvStaminaTime)

	if stamina < 0 && this.Stamina+staminaAdd+stamina < 0 {
		return enoughStamina
	}

	this.Stamina += (stamina + staminaAdd)
	this.LastRecvStaminaTime = staminaTime
	return db.Get().Model(this).Updates(map[string]interface{}{
		"stamina":             this.Stamina,
		"lastRecvStaminaTime": this.LastRecvStaminaTime,
	}).Error
}

func (this *Player) GetAddStamina(Stamina, LastRecvStaminaTime int) (int, int) {
	now := int(time.Now().Unix())
	staminaTime := now
	staminaAdd := 0
	if LastRecvStaminaTime != 0 {
		staminaAdd = (now - LastRecvStaminaTime) / config.Param().StaminaCd
		staminaTime -= (now - LastRecvStaminaTime) % config.Param().StaminaCd
	}
	LastRecvStaminaTime = staminaTime
	if staminaAdd+Stamina > config.Param().StaminaMax {
		//只加到满体力为止
		staminaAdd = config.Param().StaminaMax - Stamina
		LastRecvStaminaTime = now
	}
	if Stamina >= config.Param().StaminaMax {
		//原本就处于体力最大值，要更新一下领取时间
		staminaAdd = 0
		LastRecvStaminaTime = now
	}
	return staminaAdd, staminaTime
}

//更新玩家离线时间
func (this *Player) UpdateLeaveTime() error {
	return db.Get().Model(this).Update("leaveTime", time.Now().Unix()).Error
}

//更新玩家引导步骤
func (this *Player) UpdateNewStep(step int) error {
	return db.Get().Model(this).Update("newStep", step).Error
}

func (this *Player) UpdateGold(gold int) error {
	return db.Get().Model(this).UpdateColumn("gold", gorm.Expr("gold + ?", gold)).Error
}

func (this *Player) UpdateDiamond(diamond int) error {
	return db.Get().Model(this).UpdateColumn("diamond", gorm.Expr("diamond + ?", diamond)).Error
}

func (this *Player) UpdateHonor(honor int) error {
	return db.Get().Model(this).UpdateColumn("honor", gorm.Expr("honor + ?", honor)).Error
}

func (this *Player) UpdateDungeonData() error {
	return db.Get().Model(this).Update("curThroughId", this.CurThroughId).Error
}

func (this *Player) AddStarAward() error {
	return db.Get().Model(this).UpdateColumn("recvStarIndex", gorm.Expr("recvStarIndex + ?", 1)).Error
}

//包括武器和角色升星
func (this *Player) UpdateRisingStar() error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"gold":  this.Gold,
		"power": this.Power,
		"honor": this.Honor,
	}).Error
}

func (this *Player) UpdatePowerAndGold() error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"gold":  this.Gold,
		"power": this.Power,
	}).Error
}

func (this *Player) GetStamina() error {
	return db.Get().Model(this).Select("playerId,stamina,lastRecvStaminaTime").Limit(1).Find(this).Error
}

//获取指定用户的每个月累计签到天数和今天签到时间
func (this *Player) GetMonthTodaySign() error {
	return db.Get().Model(this).Select("monthSignDays,todaySignTime").Find(this).Error
}

//更新指定用户的每个月累计签到天数和今天签到时间
func (this *Player) UpdateMonthTodaySign(monthSignDays int, todaySignTime int64) error {
	return db.Get().Model(this).Updates(map[string]interface{}{
		"monthSignDays": monthSignDays,
		"todaySignTime": todaySignTime,
	}).Error
}

// 更新指定用户每天计签到天数
func (this *Player) UpdateTodaySign(todaySignTime int64) error {
	return db.Get().Model(this).Update("todaySignTime", todaySignTime).Error
}

// 更新指定用户每个月累计签到天数
func (this *Player) UpdateMonthSign(monthSignDays int) error {
	return db.Get().Model(this).Update("monthSignDays", monthSignDays).Error
}
