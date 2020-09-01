package config

type RobotData struct {
	Nickname         string  `json:"nickname"`         //机器人名字
	FightTime        []int   `json:"fightTime"`        //发招时间
	SkillProbability int     `json:"skillProbability"` //能量技能概率
	ItemProbability  int     `json:"itemProbability"`  //道具概率
	AttProbability   [][]int `json:"attProbability"`   //攻击概率
	Level            int     `json:"level"`            //玩家等级
	ResType          int     `json:"resType"`          //男女类型
	Avatar           int     `json:"avatar"`           //头像
	UsedEnergySkill  int     `json:"usedEnergySkill"`  //能量技能
	Items            []int   `json:"items"`            //有的道具
	Partner1         []int   `json:"partner1"`         //骑士   角色属性顺序[trainLevel,trainStar,eqId,eqStrength,eqStar]
	Partner2         []int   `json:"partner2"`         //战士
	Partner3         []int   `json:"partner3"`         //法师
	Partner4         []int   `json:"partner4"`         //牧师
	Partner5         []int   `json:"partner5"`         //射手
}

var robotPool **RobotDataPool

type RobotDataPool map[string]RobotData

func Robots() *RobotDataPool {
	return *robotPool
}
