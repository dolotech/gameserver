package task

import (
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/protocol/task"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"reflect"
	"strconv"
)

// 消息注册
func init() {
	msg.GetMsg().Reg(route.RecvTaskAward, &task.RecvTaskAward{}, &task.RecvTaskAwarded{}, recvTaskAwardCb)
	msg.GetMsg().Reg(route.RecvCollectStarAward, &task.StarTaskAward{}, &task.StarTaskAwarded{}, recvStarAwardCb)
}

// 请求完成任务
func recvTaskAwardCb(sess server.Session, req *task.RecvTaskAward, resp *task.RecvTaskAwarded) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
	}
	tasks := &Task{PlayerId: sess.UId()}
	taskConfig := config.Task().Get(strconv.Itoa(req.TaskId))
	constantCfg := config.GetConstantCfg()
	if reflect.DeepEqual(taskConfig, config.TaskData{}) {
		resp.Code = proto.FA_INVALID_TASK_ID
		return
	}
	userTask := &model.Task{
		PlayerId: sess.UId(),
		TaskId:   req.TaskId,
	}
	if err := userTask.GetDone(); err != nil {
		log.Error(err)
		resp.Code = proto.FA_INVALID_TASK_ID
		return
	}

	if userTask.State != constantCfg.TASKSTATE.DONE {
		resp.Code = proto.FA_TASK_NOT_DONE
		return
	}

	saveAward := award.Awards{}
	saveAward.Drop(taskConfig.DropId)
	saveAward.SaveAwards(player)

	if saveAward.OnBag() {
		(&push.OnUpdateBag{}).Push(player.PlayerId)
	}
	if saveAward.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(player)
	}

	userTask.State = constantCfg.TASKSTATE.CLOSED
	tasks.UpdateTaskState(userTask.TaskId, userTask.State)
}

func recvStarAwardCb(sess server.Session, req *task.StarTaskAward, resp *task.StarTaskAwarded) {
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
	}
	player.RecvStarIndex = player.RecvStarIndex + 1
	starCfg := config.Star().Get(strconv.Itoa(player.RecvStarIndex))
	if reflect.DeepEqual(starCfg, config.StarData{}) {
		resp.Code = proto.FA_INVALID_STAR_INDEX
		return
	}
	totalStar := (&model.DungeonRecords{}).GetTotalStar(sess.UId())

	if totalStar < starCfg.Star {
		resp.Code = proto.FA_STAR_NOT_ENOUGH
		return
	}
	player.AddStarAward()
	saveAward := award.Awards{}
	saveAward.Drop(starCfg.DropID)
	saveAward.SaveAwards(player)
	if saveAward.OnBag() {
		(&push.OnUpdateBag{}).Push(player.PlayerId)
	}
	if saveAward.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(player)
	}
}

//任务info
type TaskInfo struct {
	Id       int //任务标识
	State    int    //任务状态
	Count    int    //任务进度
	TaskType int    //任务类型	(1 主线任务 2 日常任务 3 成就任务)
	TaskJump int    //跳转类型
}

type Task struct {
	PlayerId uint
	Tasks    []TaskInfo //用于存放从数据库拿出来的task或者新增的task
	TasksOn	 []TaskInfo //用于存放当前更新或者新增的task，用于推送
}

/**
 * @description 首次初始化，通过配置文件来进行
 */
func (this *Task) InitByConfig() {
	constant := config.GetConstantCfg()
	for taskId, taskCfg := range *config.Task() {
		if taskCfg.IsFrist != 1 {
			continue
		}
		newTask := this.AddNewTask(taskId, &taskCfg, constant.TASKSTATE.DOING)
		if taskCfg.TaskType == constant.TASKTYPE.ACHIEVE {
			//成就任务初始设置
			switch taskCfg.TaskJump {
			case constant.TASKJUMPTYPE.PLAYERLEVEL: //主角等级初始化1一级
				newTask.Count = 1
			case constant.TASKJUMPTYPE.PVPACH: //pvp 段位初始化1一级
				newTask.Count = 1
			case constant.TASKJUMPTYPE.ROLEADVCE: //角色 阶级初始化1一级
				newTask.Count = 1
			case constant.TASKJUMPTYPE.EQUIMENTADVCE: //武器 初始化1阶
				newTask.Count = 1
			}
		}
	}
	this.AddTasksToDB()
}

func (this *Task) PushTasks() {
	if len(this.TasksOn) <= 0 {
		return
	}
	pushTasks := &push.OnUpdatePlayerTasks{PlayerId: this.PlayerId}
	for _, gt := range this.TasksOn {
		pushTasks.Tasks = append(pushTasks.Tasks, push.Task{
			State: gt.State,
			Count: gt.Count,
			Id:    gt.Id,
		})
	}
	pushTasks.PushTasks()
}

func (this *Task) AddTasksToDB() {
	tasks := model.Tasks{}
	for _, t := range this.Tasks {
		tasks = append(tasks, &model.Task{
			PlayerId: this.PlayerId,
			TaskId:   t.Id,
			State:    t.State,
			Count:    t.Count,
			Type:     t.TaskType,
			TaskJump: t.TaskJump,
		})
	}
	tasks.AddTasks()
}

func (this *Task) AddTaskToDB(newTask *TaskInfo) {
	taskM := &model.Task{
		TaskId:   newTask.Id,
		PlayerId: this.PlayerId,
		State:    newTask.State,
		Count:    newTask.Count,
		Type:     newTask.TaskType,
		TaskJump: newTask.TaskJump,
	}
	if err := taskM.AddTask(); err != nil {
		log.Error(err)
	}
}

func (this *Task) UpdateTaskState(id int, state int) {
	if err := (&model.Task{}).UpdateTaskState(this.PlayerId, id, state); err != nil {
		log.Error(err)
	}
}

func (this *Task) UpdateTaskStateCount(id int, state int, count int) {
	if err := (&model.Task{}).UpdateTaskStateCount(this.PlayerId, id, state, count); err != nil {
		log.Error(err)
	}
}

/**
 * @description 添加一个任务
 * @param {string} taskId - 任务标识
 * @param {number} type - 任务类型，1为主线，2为日常，3为成就
 * @param {Constant.TASK_STATE} state - 任务状态
 * @returns {task}
 */
func (this *Task) AddNewTask(taskId string, taskCfg *config.TaskData, state int) *TaskInfo {
	id, _ := strconv.Atoi(taskId)
	this.Tasks = append(this.Tasks, TaskInfo{Id: id, TaskType: taskCfg.TaskType, Count: 0, State: state, TaskJump: taskCfg.TaskJump})
	return &this.Tasks[len(this.Tasks)-1]
}

/**
 * @description 通过任务标识来查找任务
 * @param {string} id - 任务标识
 * @returns {task} task - 任务
 */
func (this *Task) GetTaskById(id string) (value TaskInfo, has bool) {
	task := TaskInfo{}
	has = false
	taskId,_ := strconv.Atoi(id)

	if len(this.Tasks) <= 0 {
		this.getTasksFromDb()
	}
	for _, v := range this.Tasks {
		if taskId == v.Id {
			task = v
			has = true
			break
		}
	}
	return task, has
}

//初始化，从数据库获取所有在做的任务(状态为1，即DOING)
func (this *Task) init(){
	if len(this.Tasks) <= 0 {
		this.getTasksFromDb()
	}
}

func (this *Task) getTasksFromDb() {
	tasks := &model.Tasks{}
	tasks.GetDoing(this.PlayerId)
	for _, t := range *tasks {
		this.Tasks = append(this.Tasks, TaskInfo{
			Id:       t.TaskId,
			State:    t.State,
			Count:    t.Count,
			TaskType: t.Type,
			TaskJump: t.TaskJump,
		})
	}
}

func (this *Task) dealTask(taskOn TaskInfo, cfg config.TaskData, nextCount int){
	Constant := config.GetConstantCfg()
	if taskOn.Count >= cfg.Condition[1] {
		taskOn.State = Constant.TASKSTATE.DONE
		if cfg.NextTaskId > 0 {
			newTask := this.AddNewTask(strconv.Itoa(cfg.NextTaskId), &cfg, Constant.TASKSTATE.DOING)
			newTask.Count = nextCount
			this.AddTaskToDB(newTask)
			this.TasksOn = append(this.TasksOn, *newTask)
		}
	}
	this.TasksOn = append(this.TasksOn, taskOn)
	err := (&model.Task{}).UpdateTaskStateCount(this.PlayerId, taskOn.Id, taskOn.State, taskOn.Count)
	if err != nil{
		log.Error(err)
	}
}

/**
 * @description 更新副本情况主线 关卡 日常任务
 * @param {string} dungeonId - 副本标识
 */
func (this *Task) UpdateDungeon(dungeonId int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		switch t.TaskType {
		case Constant.TASKTYPE.MAIN: //主线任务 有后置任务
			if t.TaskJump == Constant.TASKJUMPTYPE.DUNGEON &&
				dungeonId == cfg.Condition[0] {
				t.Count++
				this.dealTask(t, cfg, 0)
			}
		case Constant.TASKTYPE.DAILY: //日常任务
			if (t.TaskJump == Constant.TASKJUMPTYPE.DUNGEON && dungeonId == cfg.Condition[0]) ||
				t.TaskJump == Constant.TASKJUMPTYPE.DUNGEONALL{
				t.Count++
				this.dealTask(t, cfg, 0) //无后置任务
			}
		case Constant.TASKTYPE.ACHIEVE:
			if t.TaskJump == Constant.TASKJUMPTYPE.DUNGEONALL{
				t.Count++
				this.dealTask(t, cfg, cfg.Condition[1])
			}
		}
	}
}

/**
 * @description 更新购买体力情况任务（日常任务）
 */
func (this *Task) UpdateBuyStamina() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump == Constant.TASKJUMPTYPE.STAMINA &&
			t.TaskType == Constant.TASKTYPE.DAILY{
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}


/**
 * @description 更新pvp情况匹配任务（日常任务）
 */
func (this *Task) UpdatePvpDaily() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump == Constant.TASKJUMPTYPE.PVPITMES &&
			t.TaskType == Constant.TASKTYPE.DAILY{
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新pvp段位任务（成就任务）
 */
func (this *Task) UpdatePvpAch(rank int){
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.PVPACH {
			if t.Count < rank {
				t.Count = rank //需要注意的是，这种做法仅在rank不会连跳两级的前提下成立，否则会少算完成的任务
				this.dealTask(t, cfg, rank)
			}
		}
	}
}

/**
 * @description 更新PVB 战斗成就任务
 */
func (this *Task) UpdatePvbAch(point int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.PVB {
			if t.Count < point {
				t.Count = point
				this.dealTask(t, cfg, cfg.Condition[1])
			}
		}
	}
}

/**
 * @description 更新PVB 战斗日常任务
 */
func (this *Task) UpdatePvbDaily() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.DAILY &&
			t.TaskJump == Constant.TASKJUMPTYPE.PVB {
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新pva情况匹配任务（日常任务）
 */
func (this *Task) UpdatePvaDaily(){
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.DAILY &&
			t.TaskJump == Constant.TASKJUMPTYPE.PVAITMES {
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新强化次数任务（成就+日常任务）
 */
func (this *Task) UpdateWeaponStrength() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump != Constant.TASKJUMPTYPE.EQUIMENTSTRENTH {
			continue
		}
		if t.TaskType == Constant.TASKTYPE.DAILY  {
			t.Count++
			this.dealTask(t, cfg, 0)
		}else if t.TaskType == Constant.TASKTYPE.ACHIEVE {
			t.Count++
			this.dealTask(t, cfg, cfg.Condition[1])
		}
	}
}

/**
 * @description 更新训练英雄任务（日常任务）
 */
func (this *Task) UpdateRoleTrainDaily() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.DAILY &&
			t.TaskJump == Constant.TASKJUMPTYPE.ROLETRIAN {
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新能量技次数任务（成就+日常任务）
 */
func (this *Task) UpdateSkillPower() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump != Constant.TASKJUMPTYPE.POWERSKILL {
			continue
		}
		if t.TaskType == Constant.TASKTYPE.DAILY  {
			t.Count++
			this.dealTask(t, cfg, 0)
		}else if t.TaskType == Constant.TASKTYPE.ACHIEVE {
			t.Count++
			this.dealTask(t, cfg, cfg.Condition[1])
		}
	}
}

/**
 * @description 更新玩家花费钻石任务（日常任务）
 */
func (this *Task) UpdateSpendDiamondDaily()  {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump == Constant.TASKJUMPTYPE.SPENDDIAMOND &&
			t.TaskType == Constant.TASKTYPE.DAILY{
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新玩家赠送体力任务（日常任务）
 */
func (this *Task) UpdateFriendGiveDaily() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump == Constant.TASKJUMPTYPE.FRIENDGIVE &&
			t.TaskType == Constant.TASKTYPE.DAILY{
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新好友对战（日常任务）
 */
func (this *Task) UpdateFriendPVP() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskJump == Constant.TASKJUMPTYPE.FRIENDPVP &&
			t.TaskType == Constant.TASKTYPE.DAILY{
			t.Count++
			this.dealTask(t, cfg, 0)
		}
	}
}

/**
 * @description 更新pva累计胜利任务（成就任务）
 */
func (this *Task) UpdatePvaAch() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.PVAACH {
			t.Count++
			this.dealTask(t, cfg, cfg.Condition[1])
		}
	}
}

/**
 * @description 更新主角等级任务（成就任务）
 */
func (this *Task) UpdatePlayerLevel(level int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.PLAYERLEVEL {
			if t.Count < level {
				t.Count = level
				this.dealTask(t, cfg, level)
			}
		}
	}
}

/**
 * @description 更新任意角色进阶任务（成就任务）
 */
func (this *Task) UpdatePartnerAdvce(level int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.ROLEADVCE {
			if t.Count < level{
				t.Count = level
				this.dealTask(t, cfg, level)
			}
		}
	}
}

/**
 * @description 更新任意武器进阶任务（成就任务）
 */
func (this *Task) UpdateWeaponAdvce(level int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.EQUIMENTADVCE {
			if t.Count < level{
				t.Count = level
				this.dealTask(t, cfg, level)
			}
		}
	}
}

/**
 * @description 更新使用金币任务（成就任务）
 */
func (this *Task) UpdatePlayerGold(amount int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.GOLD {
			t.Count += amount
			nextCount := t.Count - cfg.Condition[1] //之前进度没有超出的部分归零
			this.dealTask(t, cfg, nextCount)
		}
	}
}

/**
 * @description 更新签到币任务（成就任务）
 */
func (this *Task) UpdatePlayerSign() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.SIGN {
			t.Count++
			this.dealTask(t, cfg, t.Count)
		}
	}
}

/**
 * @description 更新加好友任务（成就任务）
 */
func (this *Task) UpdatePlayerFriends() {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.FRIEND {
			t.Count++
			this.dealTask(t, cfg, t.Count)
		}
	}
}

/**
 * @description 更新加集星任务（成就任务）
 *@star 总星数
 */
func (this *Task) UpdatePlayersStar(star int) {
	this.init()
	Constant := config.GetConstantCfg()
	for _ , t := range this.Tasks {
		cfg := config.Task().Get(strconv.Itoa(t.Id))
		if t.TaskType == Constant.TASKTYPE.ACHIEVE &&
			t.TaskJump == Constant.TASKJUMPTYPE.STAR {
			if t.Count < star {
				t.Count = star
				this.dealTask(t, cfg, star)
			}
		}
	}
}

/**
 * @description 每天重置任务
 */
func (this *Task) RestTask() []config.TaskData {
	var tasksConfig []config.TaskData
	Constant := config.GetConstantCfg()
	tasks := &model.Tasks{}
	tasks.GetByType(this.PlayerId, Constant.TASKTYPE.DAILY)
	for _, t := range *tasks {
		this.Tasks = append(this.Tasks, TaskInfo{
			Id:       t.TaskId,
			State:    t.State,
			Count:    t.Count,
			TaskType: t.Type,
		})
	}
	for key, taskInfo := range this.Tasks {
		config := config.Task().Get(strconv.Itoa(taskInfo.Id))
		if taskInfo.State == Constant.TASKSTATE.DONE { //完成没领取要发邮件
			tasksConfig = append(tasksConfig, config)
		}
		taskInfo.Count = 0
		taskInfo.State = Constant.TASKSTATE.DOING
		this.UpdateTaskStateCount(taskInfo.Id, taskInfo.State, taskInfo.Count)
		this.Tasks[key] = taskInfo
	}
	return tasksConfig
}
