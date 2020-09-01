package config

import (
	"strconv"
)

//关卡数据
type TollgateData struct {
	Point       int `json:"point"`       //关卡的关数
	Stamina     int `json:"stamina"`     //体力消耗
	DropId      int `json:"dropId"`      //有可能掉落
	FristDropId int `json:"fristDropId"` //首杀掉落
	PlayExp     int `json:"playExp"`     //关卡提供角色经验
	DefSteps    int `json:"defSteps"`    //默认步数
}

var tollgatePool **TollgateDataPool

type TollgateDataPool map[string]TollgateData

func Tollgate() *TollgateDataPool {
	return *tollgatePool
}

func (this *TollgateDataPool) Get(id string) TollgateData {
	return (*this)[id]
}

func (this *TollgateDataPool) GetAllId() []string {
	ids := []string{}
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}

func GetNextDungeonId(DungeonId int) int {
	param := Param()
	chapterId := (DungeonId-param.MinThroughDungeonId)/param.ChapterCap + 1
	chapter := Chapter().Get(strconv.Itoa(chapterId))
	lastId := chapter.StartId + chapter.Num - 1
	//log.Error("chapter id:", chapterId, "last id:", lastId, "chapter:", chapter)
	if DungeonId < lastId {
		return DungeonId + 1
	}
	if Chapter().GetMaxId() == chapterId {
		//最后一关
		return DungeonId
	}
	return Chapter().Get(strconv.Itoa(chapterId + 1)).StartId
}
