package task

import (
	"flag"
	"fmt"
	"gameserver/config"
	"gameserver/utils"
	"gameserver/utils/db"
	"gameserver/utils/log"
	"github.com/spf13/viper"
	"testing"
)

func Test_Task(t *testing.T) {
	var fileName string
	flag.StringVar(&fileName, "conf", "cfg.toml", "Configuration file to start game")
	flag.Parse()
	if !utils.PathExists(fileName) {
		fmt.Println("conf file not exist !!", fileName)
		return
	}
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: %s \n", err)
	}
	log.Info("conf:", fileName)

	if err := db.Get().Error; err != nil {
		log.Error(err)
	}

	//TEST
	config.InitJson("config")


	t.Error("TaskData:", config.Task().Get("35001"))

	Constant := config.GetConstantCfg()
	t.Error(Constant.TASKTYPE.MAIN)

	t.Error(Constant.TASKJUMPTYPE.DUNGEONALL)

	//task := Task{}
	//task.InitByConfig()
	//tasks := task.GetTasksList()
	//for _, v := range tasks {
	//	if v.Count > 0 {
	//		t.Error(v)
	//	}
	//}
	//t.Error( tasks )

	//var tasks []TaskInfo
	//task := &TaskInfo{}
	//task.State = 0
	//task.Count =1
	//task.TaskType = 2
	//task.Id = "3"
	//
	//tasks = append(tasks, *task)
	//
	//task.State = 2
	//
	//t.Error(tasks)
	//
	//t.Error(task)
	//
	//tasks[len(tasks)-1] = *task   //TODO
	//
	//t.Error(tasks)
	//myTasks := &Task{PlayerId: 2}
	//tasks := myTasks.UpdateDungeon(41004)
	////t.Error(tasks)
	//for _, v := range tasks {
	//	if v.State == 2 {
	//		t.Error(v)
	//	}
	//}
	//t.Error(int(utils.InetToaton("192.168.1.190")))
	//
	//t.Error(config.Task().Get(strconv.Itoa(233)))
}

func Test_DB(t *testing.T) {
	var fileName string
	flag.StringVar(&fileName, "conf", "cfg.toml", "Configuration file to start game")
	flag.Parse()
	if !utils.PathExists(fileName) {
		fmt.Println("conf file not exist !!", fileName)
		return
	}
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: %s \n", err)
	}
	log.Info("conf:", fileName)

	if err := db.Get().Error; err != nil {
		log.Error(err)
	}



}