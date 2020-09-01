package protocol

import (
	"gameserver/config"
	"gameserver/model"
	"gameserver/utils"
	"gameserver/utils/log"
	"time"
)

//训练
type Train struct {
	Level int `json:"level"`
	Point int `json:"point"`
	Star  int `json:"star"`
}

//强化
type Equipment struct {
	Strength int `json:"strength"`
	Star     int `json:"star"`
	WeaponId int `json:"weaponId"` //武器id
}

//物品单元
type Item struct {
	Count int `json:"count"` //物品数量
	ID    int `json:"id"`    //物品Id
}

//角色
type Partner struct {
	Level     int       `json:"level"` //等级从1开始
	Exp       int       `json:"exp"`
	PartnerId int       `json:"partnerId"` //角色ID
	Train     Train     `json:"train"`     //角色训练
	Equipment Equipment `json:"equipment"` //装备的实例
}

type Partners []Partner

func (this *Partners) CopyTo(player *model.Player, partners *model.Partners) {
	for _, par := range *partners {
		p := Partner{
			PartnerId: par.PartnerId,
			Level:     player.Level,
		}
		// 武器属性
		p.Equipment.Star = par.WStar
		p.Equipment.WeaponId = par.WeaponId
		p.Equipment.Strength = par.Strength
		// 角色属性
		p.Train.Star = par.PStar
		p.Train.Point = par.Point
		p.Train.Level = par.Level
		(*this) = append((*this), p)
	}
}

type Player struct {
	PlayerId               uint     `json:"playerId"`
	Nickname               string   `json:"nickname"` //昵称
	ResType                int      `json:"resType"`  // 男女，1为男，2为女 默认为男
	Exp                    int      `json:"exp"`      //经练
	Diamond                int      `json:"diamond"`  //钻石
	Gold                   int      `json:"gold"`     //金币
	Honor                  int      `json:"honor"`    //荣誉值
	Stamina                int      `json:"stamina"`  //体力值
	Avatar                 int      `json:"avatar"`   //头像，0为默认头像
	Vip                    int      `json:"vip"`      //vip等级
	Level                  int      `json:"level"`
	Power                  int      `json:"power"`                  //战斗力
	Items                  []Item   `json:"items"`                  //背包物品
	Partners               Partners `json:"partners"`               //所有的角色
	UsedEnergySkill        int      `json:"usedEnergySkill"`        // 当前使用的能量技，在initByConfig或者initByDB中初始化
	IsGuide                int      `json:"isGuide"`                //0开 1关闭	(新手引导的开关 ANJUN)
	NewStep                int      `json:"newStep"`                //主线引导步骤
	GuideSteps             string   `json:"guideSteps"`             //异步引导步骤Battle1@0&Battle2@0  0代表没执行过1代表执行
	RecvOpenSvrAwardIndex  int      `json:"recvOpenSvrAwardIndex"`  //开服奖励领取到哪个索引
	DayOnlineTotalTime     int      `json:"dayOnlineTotalTime"`     //在线奖励累计时间
	TodaySignTime          int      `json:"todaySignTime"`          //今天签到时间
	ArenaPKLoses           int      `json:"arenaPKLoses"`           //竞技场输次数
	WeeklyRank             int      `json:"weeklyRank"`             //周赛排名,这个值只是镜像
	RecvLevelAwardIndex    int      `json:"recvLevelAwardIndex"`    //领取登记奖励的索引
	LastArenaAwardIndex    int      `json:"lastArenaAwardIndex"`    //arena周末段位奖励,要发给客户端
	DayArenaPKTimes        int      `json:"dayArenaPKTimes"`        //每天竞技场pk次数
	ContinuousArenaPKWins  int      `json:"continuousArenaPKWins"`  //竞技场连赢次数
	ChallengeTimes         int      `json:"challengeTimes"`         //挑战关卡的次数
	StaminaCountDownTime   int      `json:"staminaCountDownTime"`   //离上次恢复体力过去了多久，时间为毫秒
	MonthSignDays          int      `json:"monthSignDays"`          //每个月累计签到天数
	ArenaPKWins            int      `json:"arenaPKWins"`            //竞技场赢次数
	StaminaRecvTimes       int      `json:"staminaRecvTimes"`       //当前领取好友体力的次数
	OwnedEnergySkills      []int    `json:"ownedEnergySkills"`      // 拥有的能量技，在initByConfig或者initByDB中初始化
	ChallengeCountDownTime int      `json:"challengeCountDownTime"` //返回从上次恢复体力后，过去了多少时间，单位为毫秒。要注意：这个函数最好是紧接着getChallengeTimes后调用，不然数值会有偏差  (原来是由getChallengeCountDownTime方法获取 ANJUN )
	RecvStarIndex          int      `json:"recvStarIndex"`          //领取集星奖励的索引
	LoginTotalDays         int      `json:"loginTotalDays"`         //索引从1有效
	LeftWeeklyAwardTime    int      `json:"leftWeeklyAwardTime"`    //领奖剩余时间
	OsTime                 int64    `json:"osTime"`                 //
	RecvOnlineAwardIndex   int      `json:"recvOnlineAwardIndex"`   //索引从1有效
	TotalDungeonStar       int      `json:"totalDungeonStar"`       //推图副本星星总和
	MaxChallengeId         int      `json:"maxChallengeId"`         // 获取最远打到的挑战关卡
	CurThroughId           int      `json:"curThroughId"`           //当前挑战的关卡索引
}

func (this *Player) P2P(player *model.Player) (partners model.Partners, err error) {
	utils.StructAtoB(this, player)

	//--------------------计算竞技塔排行榜奖励剩余时间--------------------------
	this.LeftWeeklyAwardTime = config.Param().PvaAwardsTime - (int(time.Now().Unix()) - player.LeftWeeklyAwardTime)
	//--------------------以下属性是数据库和网络协议不对称的--------------------------
	this.OsTime = time.Now().Unix()
	this.StaminaCountDownTime = player.StaminaRecvTimes
	this.ChallengeCountDownTime = player.ChallengeTimes
	//---------------------OwnedEnergySkills-----------------------
	ownedSkills := model.OwnedEnergySkills{}
	if err = ownedSkills.Get(player.PlayerId); err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < len(ownedSkills); i++ {
		skill := ownedSkills[i]
		this.OwnedEnergySkills = append(this.OwnedEnergySkills, skill.SkillId)
		if skill.Status > 0 {
			this.UsedEnergySkill = skill.SkillId
		}
	}
	//---------------------Items-----------------------
	items := model.Items{}
	if err = items.Get(player.PlayerId); err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < len(items); i++ {
		it := items[i]
		item := Item{ID: it.ItemId, Count: it.Count}
		this.Items = append(this.Items, item)
	}
	//----------------------Partners--------------------------------
	if err = partners.Get(player.PlayerId); err != nil {
		log.Error(err)
		return
	}
	this.Partners.CopyTo(player, &partners)
	return
}
