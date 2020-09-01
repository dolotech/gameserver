package config

import (
	"strconv"
)

//角色的升级配置文件，每等级升级所需要的经验值
type ExpData struct {
	Exp int `json:"exp"` //升级经验
}

var expPool **ExpDataPool

type ExpDataPool map[string]ExpData

func Exp() *ExpDataPool {
	return *expPool
}

func (this *ExpDataPool) Get(id string) ExpData {
	return (*this)[id]
}

func (this *ExpDataPool) GetMaxLevel() int {
	max := 0
	for id, _ := range *this {
		i, _ := strconv.Atoi(id)
		if max < i {
			max = i
		}
	}
	return max
}
