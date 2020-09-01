package model

import (
	"fmt"
	"gameserver/config"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"strings"
)

// 任务   gameserver
type Task struct {
	Id       uint `gorm:"primary_key;column:id;AUTO_INCREMENT;COMMENT:'任务唯一标识'"`
	TaskId   int  `gorm:"type:int(11);column:taskId;unique_index:index_playerId" json:"taskId" `
	PlayerId uint `gorm:"type:int(11);column:playerId;unique_index:index_playerId" json:"playerId" `
	State    int  `gorm:"type:int(8);column:state" json:"state" `
	Count    int  `gorm:"type:int(8);column:count;COMMENT:'用于进度条'" json:"count" `
	Type     int  `gorm:"type:int(8);column:type" json:"type" `
	TaskJump int  `gorm:"type:int(8);column:taskJump" json:"taskJump" `
}
type Tasks []*Task

func (this *Task) AddTask() error {
	return db.Get().Create(this).Error
}

func (this *Task) UpdateTaskCount(playerId uint, taskId int, count int) error {
	return db.Get().Model(this).Where("playerId = ? and taskId= ?", playerId, taskId).Update("count", count).Error
}

func (this *Task) UpdateTaskState(playerId uint, taskId int, state int) error {
	return db.Get().Model(this).Where("playerId = ? and taskId= ?", playerId, taskId).Update("state", state).Error
}

func (this *Task) UpdateTaskStateCount(playerId uint, taskId int, state int, count int) error {
	return db.Get().Model(this).Where("playerId = ? and taskId= ?", playerId, taskId).Updates(map[string]interface{}{
		"state": state,
		"count": count,
	}).Error
}

func (this *Task) GetDone() error {
	ConstantCfg := config.GetConstantCfg()
	return db.Get().Model(this).Where("playerId = ? and taskId = ? and state = ?", this.PlayerId, this.TaskId, ConstantCfg.TASKSTATE.DONE).Find(this).Error
}

func (this *Tasks) GetDoing(playerID uint) error {
	ConstantCfg := config.GetConstantCfg()
	return db.Get().Model(this).Where("playerId = ? and state = ?", playerID, ConstantCfg.TASKSTATE.DOING).Find(this).Error
}

func (this *Tasks) GetDoingDone(playerID uint) error {
	ConstantCfg := config.GetConstantCfg()
	return db.Get().Model(this).Where("playerId = ? and (state = ? or state = ?)", playerID, ConstantCfg.TASKSTATE.DOING, ConstantCfg.TASKSTATE.DONE).Find(this).Error
}

func (this *Tasks) GetByType(playerID uint, typeId int) error {
	return db.Get().Model(this).Where("playerId = ? and type = ?", playerID, typeId).Find(this).Error
}




//一次添加多条记录
func (this *Tasks) AddTasks() {
	taskStrings := make([]string, 0, len(*this))
	taskArgs := make([]interface{}, 0, len(*this)*5)
	for _, t := range *this {
		taskStrings = append(taskStrings, "(?, ?, ?, ?, ?, ?)")
		taskArgs = append(taskArgs, t.TaskId)
		taskArgs = append(taskArgs, t.PlayerId)
		taskArgs = append(taskArgs, t.Count)
		taskArgs = append(taskArgs, t.State)
		taskArgs = append(taskArgs, t.Type)
		taskArgs = append(taskArgs, t.TaskJump)
	}
	stmt := fmt.Sprintf("INSERT INTO tasks (taskId, playerId, count, state, type, taskJump) VALUES %s",
		strings.Join(taskStrings, ","))
	if err := db.Get().Exec(stmt, taskArgs...).Error; err != nil {
		log.Error(err)
	}
	return
}
