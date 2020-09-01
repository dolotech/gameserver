package push

import (
	"gameserver/model"
	"gameserver/online"
)

//更新关卡信息PUSH
type OnUpdateThroughRecords struct {
	Records []Record `json:"records"` //所有通关记录
}

type Record struct {
	Star int `json:"star"` //关卡星级
	ID   int `json:"id"`   //关卡ID
}

func (this *OnUpdateThroughRecords) Push(player *model.Player,after ...bool) error {
	dungeonRecords := &model.DungeonRecords{}
	if dungeonRecords.Get(player.PlayerId) == nil {
		for _, record := range *dungeonRecords {
			info := Record{Star: record.Star, ID: record.DungeonId}
			this.Records = append(this.Records, info)
		}
		online.Get().Push(player.PlayerId, this,after...)
	}
	return nil
}
