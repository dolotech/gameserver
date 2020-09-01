package config

//竞技场开放时间12:00——14:00
type time struct {
	StHour int `json:"stHour"` //开始时
	StMin  int `json:"stMin"`  //开始分
	EdHour int `json:"edHour"` //结束时
	EdMin  int `json:"edMin"`  //结束分
}

//竞技场开放时间等配置
type ArenaConfig struct {
	Time []time `json:"time"` // 0:竞技场开放时间12:00——14:00; 	1:竞技场开放时间18:00——20:00
}

var arenaConfig **ArenaConfig

func GetArenaCfg() *ArenaConfig {
	return *arenaConfig
}
