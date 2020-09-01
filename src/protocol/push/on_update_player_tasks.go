package push

import (
	"gameserver/model"
	"gameserver/online"
)

//更新任务PUSH
type OnUpdatePlayerTasks struct {
	PlayerId uint   `json:"playerId"`
	Tasks    []Task `json:"tasks"` //所有任务信息
}

//任务信息
type Task struct {
	State int `json:"state"` //任务状态
	Count int `json:"count"` //任务进度
	Id    int `json:"id"`    //任务标识
}

func (this *OnUpdatePlayerTasks) Push(player *model.Player,after ...bool) error {
	tasks := &model.Tasks{}
	tasks.GetDoingDone(player.PlayerId)
	this.PlayerId = player.PlayerId
	for _, t := range *tasks {
		this.Tasks = append(this.Tasks, Task{
			State: t.State,
			Count: t.Count,
			Id:    t.TaskId,
		})
	}
	online.Get().Push(player.PlayerId, this,after...)
	return nil
}

func (this *OnUpdatePlayerTasks) PushTasks() error {
	online.Get().Push(this.PlayerId, this)
	return nil
}


