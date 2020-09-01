package config

import (
	"flag"
	"fmt"
	"gameserver/utils"
	"gameserver/utils/log"
	"github.com/spf13/viper"
	"math/rand"
	"testing"
)

func Test_robot(t *testing.T) {

}
func Test_Shop(t *testing.T) {

	//for k,v:=range *Skill(){
		//t.Error(k,v)
	//}
	t.Error(rand.Int31n(5))
	t.Error(Skill().Get("35001"))
}
func Test_json(t *testing.T) {
	// using
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
	// 初始化json表
	Load(viper.GetString("general.configDir"))
	//TollgateData := utils.DataPool{}
	//utils.DecodeJson("TollgateData", &TollgateData)
	//t.Error("TollgateData:", TollgateData)
	//dataMap := utils.DataMap{}
	//utils.DecodeJson("ParamData", &dataMap)
	//t.Error("ParamData:", dataMap)
	//
	//WeaponStarUpData := WeaponStarUpPool{}
	//utils.DecodeJson("WeaponStarUpData", &WeaponStarUpData)
	//t.Error("WeaponStarUpData:", WeaponStarUpData)
	//
	//ParamData := ParamData{}
	//utils.DecodeJson("ParamData", &ParamData)
	//t.Error("ParamData==========:", ParamData)
	//
	//t.Error("weapon:", Weapon().Get("130001"))
	//
	//t.Error("RoleData:",Role().Get("10001"))
	//
	//t.Error("SignAwardData:",SignAward().Get("18").DropId)
	//t.Error("TaskData:", Task().Get("36104").NextTaskId)
	//t.Error("GetConstantCfg:",GetConstantCfg().TASKSTATE.DOING)
	t.Error("TaskData:", GetArenaCfg().Time)

	for k, v := range GetArenaCfg().Time{
		t.Error(k,v)
	}


	//keys := Task()
	//for taskId, _ := range *keys{
	//	taskinfo := Task().Get(taskId)
	//	//if (taskinfo.IsFrist==1) {
	//		fmt.Println("taskinfo:",taskinfo.Condition)
	//	//}
	//
	//}
	//t.Error("chapter data:", Chapter().Get(strconv.Itoa(1)))
	//t.Error("param", Param().MinThroughDungeonId)
	//t.Error("TaskData:", Goods().Get("6").Sell)
}
