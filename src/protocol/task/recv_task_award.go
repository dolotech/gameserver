package task


//完成任务
type RecvTaskAward struct {
	TaskId int `json:"taskId"`
}


//完成任务返回
type RecvTaskAwarded struct {
	Code   int `json:"code"`
}

//完成任务
type StarTaskAward struct {
}


//完成任务返回
type StarTaskAwarded struct {
	Code   int `json:"code"`
}