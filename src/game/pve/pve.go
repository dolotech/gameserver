package pve

import (
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/game/data"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"strconv"
)

// 游戏消息注册
func init() {
	msg.GetMsg().Reg(route.StartFightThrough, &proto.StartFightThrough{}, &proto.StartFightThroughed{}, startFightThroughCb)
	msg.GetMsg().Reg(route.FinishFightThrough, &proto.FinishFightThrough{}, &proto.FinishFightThroughed{}, FinishFightThroughCb)
	msg.GetMsg().Reg(route.UseItem, &proto.UseItem{}, &proto.UseItemed{}, useItemCb)
}

// 战斗外使用道具
func useItemCb(sess server.Session, req *proto.UseItem, resp *proto.UseItemed) {
	resp.Code = proto.OK

	//如果是机器人直接使用  TODO

	bag := &model.Item{
		PlayerId: sess.UId(),
		ItemId:   req.ItemId,
	}
	itmeCount, _ := bag.GetCount()
	if itmeCount < 1 {
		log.Warning("has not enough itemId! itemId: %d", req.ItemId)
		resp.Code = proto.FAIL //FA_NOT_ENOUGH_ITEMS
		return
	}
	bag.Reduce(1)
	resp.ItemId = req.ItemId
	resp.Key = req.Key
}

// 请求开始通关
func startFightThroughCb(sess server.Session, req *proto.StartFightThrough, resp *proto.StartFightThroughed) {
	resp.Code = proto.OK

	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	toll := config.Tollgate().Get(strconv.Itoa(req.DungeonId))
	if err := player.UpdateStamina(-toll.Stamina); err != nil {
		log.Error(data.ErrUpdateStamina, err)
		resp.Code = proto.FA_NOT_ENOUGH_STAMINA
		return
	}
	resp.Stamina = player.Stamina
}

// 请求完成通关
func FinishFightThroughCb(sess server.Session, req *proto.FinishFightThrough, resp *proto.FinishFightThroughed) {
	//	客户端Req {"route":"logic.playerHandler.finishFightThrough","id":9,"type":2,"body":{"star":3,"code":200,"items":[{"count":1,"id":22001},{"count":1000,"id":2},{"count":10,"id":5}]},"compressRoute":0}
	resp.Code = proto.OK

	player := &model.Player{PlayerId: sess.UId()}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	if player.CurThroughId < req.DungeonId {
		resp.Code = proto.FA_INVALID_DUNGEON_ID
		return
	}
	//todo 数据分析是否作弊

	toll := config.Tollgate().Get(strconv.Itoa(req.DungeonId))
	//---------------------------奖励掉落------------------
	dropAwards := &award.Awards{}
	if req.DungeonId == player.CurThroughId {
		//初次通关
		nextDungeonId := config.GetNextDungeonId(req.DungeonId)
		player.CurThroughId = nextDungeonId
		dropAwards.Drop(toll.FristDropId)
	} else {
		//重刷关卡
		dropAwards.Drop(toll.DropId)
	}
	dropAwards.Exp += toll.PlayExp
	upgrade := dropAwards.OnLevelUp(player.Level, player.Exp)
	dropAwards.SaveAwards(player)

	resp.Items = dropAwards.ConvertItem()

	//-------------------------战斗评星----------------------
	resp.Star = getDungeonStar(toll.DefSteps, req)
	//TODO 处理error
	record := &model.DungeonRecord{}
	if err := record.SaveRecord(req.DungeonId, sess.UId(), resp.Star); err != nil {
		log.Error(data.ErrSaveDungeonRecord, err)
	}

	if upgrade {
		partners := &model.Partners{}
		partners.Get(player.PlayerId)
		player.Power = partners.CalcuPower(player.Level)
		(&push.OnUpdatePartners{}).Push(player, partners)
	}

	//----------player保存---------------------------
	if err := player.UpdateDungeonData(); err != nil {
		resp.Code = proto.FAIL
		log.Error(data.ErrSaveTollgateData, err)
		return
	}

	//--------------------任务管理---------------------
	myTasks := &task.Task{PlayerId: player.PlayerId}
	myTasks.UpdateDungeon(req.DungeonId)
	myTasks.UpdatePlayersStar((&model.DungeonRecords{}).GetTotalStar(player.PlayerId))
	if upgrade {
		myTasks.UpdatePlayerLevel(player.Level)
	}

	myTasks.PushTasks()
	(&push.OnUpdatePlayer{}).Push(player)

	if dropAwards.OnBag() {
		(&push.OnUpdateBag{}).Push(sess.UId())
	}
}

func getDungeonStar(defSteps int, req *proto.FinishFightThrough) int {
	//1,2,3 为星数
	constCfg := config.GetConstantCfg()
	if req.DungeonId < constCfg.ROOKIEGUIDE {
		//新手引导
		return 3
	}
	statusStar := 1
	statusRatio := (req.LeftHP + req.LeftDF) / (req.MaxHP + req.MaxDF)
	if statusRatio >= constCfg.STARRATE.STATUSTHREE {
		statusStar = 3
	} else if statusRatio >= constCfg.STARRATE.STATUSTWO {
		statusStar = 2
	}
	if defSteps <= 0 {
		defSteps = 1
	}
	stepStar := 1
	stepRatio := float64(req.UseSteps / defSteps)
	if stepRatio <= constCfg.STARRATE.STEPTHREE {
		stepStar = 3
	} else if stepRatio <= constCfg.STARRATE.STEPTWO {
		stepStar = 2
	}
	star := (stepStar + statusStar + 1) / 2 //加1作为四舍五入
	if star <= 0 {
		star = 1
	}
	return star
}
