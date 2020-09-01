package push

import (
	"gameserver/model"
	"gameserver/online"
	"gameserver/utils"
)

//更新玩家信息PUSH
type OnUpdatePlayer struct {
	Level               int    `json:"level"`               //等级
	Nickname            string `json:"nickname"`            //昵称
	PlayerId            uint   `json:"playerId"`            //Id
	PvpScore            int    `json:"pvpScore"`            //arena积分
	Diamond             int    `json:"diamond"`             //钻石
	TotalDungeonStar    int    `json:"totalDungeonStar"`    //推图副本星星总和
	Honor               int    `json:"honor"`               //荣耀
	PvpRank             int    `json:"pvpRank"`             //arena段位
	Power               int    `json:"power"`               //战斗力
	Exp                 int    `json:"exp"`                 //
	Gold                int    `json:"gold"`                //金币
	Vip                 int    `json:"vip"`                 //vip等级
	Stamina             int    `json:"stamina"`             //体力值
	LastRecvStaminaTime int    `json:"lastRecvStaminaTime"` //  最后一次领取体力时间
	ResType             int    `json:"resType"`             // 男女，1为男，2为女 默认为男
}

func (this *OnUpdatePlayer) Push(player *model.Player,after ...bool) error {
	player.UpdateStamina(0)		// 根据最后一次领取体力时间戳校正体力值
	utils.StructAtoB(this, player)
	//this.PvpRank
	totalStar := (&model.DungeonRecords{}).GetTotalStar(player.PlayerId)
	this.TotalDungeonStar = totalStar
	online.Get().Push(player.PlayerId, this,after...)
	return nil
}
