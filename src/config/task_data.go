package config

//任务数据
type TaskData struct {
	TaskType   int   `json:"taskType"`   //任务类型	(1 主线任务 2 日常任务 3 成就任务)
	IsFrist    int   `json:"isFrist"`    //是否是初始化任务
	NextTaskId int   `json:"nextTaskId"` //后置任务编号
	DropId     int   `json:"dropId"`     //任务掉落奖励
	TaskJump   int   `json:"taskJump"`   //任务跳转	(1 跳转关卡 2 跳转角色 3 跳转物品 4 跳转商城 5 周赛场 6 竞技场)
	Condition  []int `json:"condition"`  //完成条件参数
}

var taskPool **TaskPool

type TaskPool map[string]TaskData

func Task() *TaskPool {
	return *taskPool
}

func (this *TaskPool) Get(id string) TaskData { //id :任务ID  21开头的是用来给支线  25开头的用来给日常
	return (*this)[id]
}
