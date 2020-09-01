package config

type PvpEveryAwardData struct {
	Level     int `json:"level"`
	WinNum    int `json:"winNum"`
	WinAward  int `json:"winAward"`
	LostAward int `json:"lostAward"`
}

var pvpEveryAward **PvpEveryAwardPool

type PvpEveryAwardPool map[string]PvpEveryAwardData

func PvpEveryAward() *PvpEveryAwardPool {
	return *pvpEveryAward
}

type PvpLevelAwardData struct {
	Award int `json:"award"`
}

var pvpLevelAward **PvpLevelAwardDataPool

type PvpLevelAwardDataPool map[string]PvpLevelAwardData

func PvpLevelAward() *PvpLevelAwardDataPool {
	return *pvpLevelAward
}

type PvpLevel struct {
	Level int `json:"level"`
	Exp   int `json:"exp"`
}

var pvpLevelData **PvpLevelDataPool

type PvpLevelDataPool map[string]PvpLevel

func PvpLevelData() *PvpLevelDataPool {
	return *pvpLevelData
}
