package role

import (
	"gameserver/config"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/role"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"reflect"
	"strconv"
)

// 消息注册
func init() {
	msg.GetMsg().Reg(route.Train, &role.Train{}, &role.Trained{}, train)
	msg.GetMsg().Reg(route.UpgradeStar, &role.UpgradeStar{}, &role.UpgradeStared{}, upgradeStar)
}

func train(sess server.Session, req *role.Train, resp *role.Trained){
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners := model.Partners{}
	index := -1
	partners.Get(sess.UId())
	for i, p := range partners {
		if p.PartnerId == req.PartnerId{
			index = i
		}
	}
	if index == -1 {
		resp.Code = proto.FA_INVALID_PARTNER
		return
	}
	nextTrain := config.RoleTrain().Get(strconv.Itoa(partners[index].Level + 1))
	if reflect.DeepEqual(nextTrain, config.RoleTrainData{}){
		resp.Code = proto.FA_MAX_TRAIN_LEVEL
		return
	}
	trainCfg := config.RoleTrain().Get(strconv.Itoa(partners[index].Level))
	trainNeed := config.RoleTrainNeed()
	if player.Gold < (*trainNeed)["2"].Num{
		resp.Code = proto.FA_NOT_ENOUGH_GOLD
		return
	}

	//数据库操作
	item := model.Item{PlayerId: sess.UId(), ItemId: (*trainNeed)["1"].ItemId}
	if err := item.Reduce((*trainNeed)["1"].Num);err != nil{
		resp.Code = proto.FA_NOT_ENOUGH_ITEMS
		return
	}
	partners[index].Point += 10 //1本训练书加10经验
	if partners[index].Point >= trainCfg.TrainExp {
		//升级了
		partners[index].Point = 0
		partners[index].Level += 1
		player.Power = partners.CalcuPower(player.Level)
	}

	if err := partners[index].UpdateTrain();err != nil {
		log.Error(err)
	}
	player.Gold -= (*trainNeed)["2"].Num
	if err := player.UpdatePowerAndGold();err != nil{
		log.Error(err)
	}

	//任务处理
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdatePlayerGold((*trainNeed)["2"].Num)
	myTasks.UpdateRoleTrainDaily()
	//消息推送
	myTasks.PushTasks()
	(&push.OnUpdateBag{PlayerId:player.PlayerId}).Push(player.PlayerId)
	//返回
	resp.Power = player.Power
}

func upgradeStar(sess server.Session, req *role.UpgradeStar, resp *role.UpgradeStared){
	resp.Code = proto.OK
	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	partners := model.Partners{}
	index := -1
	partners.Get(sess.UId())
	for i, p := range partners {
		if p.PartnerId == req.PartnerId{
			index = i
		}
	}
	if index == -1 {
		resp.Code = proto.FA_INVALID_PARTNER
		return
	}
	starCfg := config.RoleStar().Get(strconv.Itoa(partners[index].PStar + 1))
	if reflect.DeepEqual(starCfg, config.RoleStarData{}) {
		resp.Code = proto.FA_MAX_TRAIN_STAR
		return
	}

	if partners[index].Level < starCfg.NeedTrain {
		resp.Code = proto.FA_NOT_ENOUGH_TRAIN_LEVEL
		return
	}
	if player.Honor < starCfg.NeedHonor {
		resp.Code = proto.FA_NOT_ENOUGH_HONOR
		return
	}
	if player.Gold < starCfg.NeedGold {
		resp.Code = proto.FA_NOT_ENOUGH_GOLD
		return
	}
	//数据库操作
	player.Gold -= starCfg.NeedGold
	player.Honor -= starCfg.NeedHonor
	partners[index].PStar +=   1
	player.Power = partners.CalcuPower(player.Level)
	if err := player.UpdateRisingStar();err != nil{
		log.Error(err)
	}
	if err := partners[index].UpdatePStar();err != nil {
		log.Error(err)
	}

	//任务处理
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdatePlayerGold(starCfg.NeedGold)
	myTasks.UpdatePartnerAdvce(partners[index].PStar)
	//消息推送
	myTasks.PushTasks()

	//返回
	resp.Power = player.Power
}