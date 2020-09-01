package protocol

import (
	jsoniter "github.com/json-iterator/go"
)

//完成通关
type FinishFightThrough struct {
	DungeonId int `json:"dungeonId"` //关卡ID
	//(计算星星用的 见getDungeonStar() ANJUN)
	LeftHP   float64 `json:"leftHP"`   //玩家当前血量
	MaxHP    float64 `json:"maxHP"`    //玩家满血量
	LeftDF   float64 `json:"leftDF"`   //玩家当前防御力
	MaxDF    float64 `json:"maxDF"`    //玩家满防御力
	UseSteps int     `json:"useSteps"` //消除次数
	//以下不用了服务器
	NpcId            int `json:"npcId"`            //没用的参数 当前挑战的npc索引
	ClearIconCount   int `json:"clearIconCount"`   //消除海马、海星、鱼、水母的个数
	ClearDebuffCount int `json:"clearDebuffCount"` //清除负面状态的个数
}

func (this *FinishFightThrough) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *FinishFightThrough) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//完成通关返回
type FinishFightThroughed struct {
	Star  int    `json:"star"` //玩家得星
	Code  int    `json:"code"`
	Items []Item `json:"items"` //玩家获取掉落奖励
}

func (this *FinishFightThroughed) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *FinishFightThroughed) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}
