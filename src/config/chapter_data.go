package config

import (
	"strconv"
)

//章节表
type ChapterData struct {
	StartId int `json:"startId"` //起始ID
	Num     int `json:"num"`     //关卡数
	DropId  int `json:"dropId"`  //章节奖励
}

var chapterPool **ChapterDataPool

type ChapterDataPool map[string]ChapterData

func Chapter() *ChapterDataPool {
	return *chapterPool
}

func (this *ChapterDataPool) Get(id string) ChapterData {
	return (*this)[id]
}

func (this *ChapterDataPool) GetMaxId() int {
	maxId := 0
	for id, _ := range *this {
		i, _ := strconv.Atoi(id)
		if maxId < i {
			maxId = i
		}
	}
	return maxId
}

func (this *ChapterDataPool) GetAllId() []string {
	ids := make([]string, 0, len(*this))
	for id, _ := range *this {
		ids = append(ids, id)
	}
	return ids
}
