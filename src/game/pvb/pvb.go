package pvb

import (
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/pvb"
	"gameserver/protocol/route"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"reflect"
	"strconv"
	"time"
)

const (
	WIN = 1
	LOSE = 0
)

func init() {
	msg.GetMsg().Reg(route.StartFightChallenge, &pvb.StartFightChallenge{}, &pvb.StartFightThroughed{}, startFightChallenge)
	msg.GetMsg().Reg(route.FinishFightChallenge, &pvb.FinishFightChallenge{}, &pvb.FinishFightChallenged{}, FinishFightChallenge)
	msg.GetMsg().Reg(route.BuyChallenge, &pvb.BuyChallenge{}, &pvb.BuyChallenged{}, buyChallenge)
	msg.GetMsg().Reg(route.ResetChallenge, &pvb.ResetChallenge{}, &pvb.ResetChallenged{}, resetChallenge) 
	msg.GetMsg().Reg(route.GetPvbInfo, &pvb.GetPvbInfo{}, &pvb.GetPvbInfoed{}, getPvbInfocb)

}
func getPvbInfocb(sess server.Session, req *pvb.GetPvbInfo, resp *pvb.GetPvbInfoed){
	resp.Code = proto.OK
	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	player.UpdateChallengeTimes(0)
	resp.ChallengeCountDownTime = config.Param().PvbRecover - (int(time.Now().Unix()) - player.LastChallengeTime)
	resp.Times = player.ChallengeTimes
}

func startFightChallenge(sess server.Session, req *pvb.StartFightChallenge, resp *pvb.StartFightThroughed){
	resp.Code = proto.OK
	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	toll := config.Tollgate().Get(strconv.Itoa(req.DungeonId))
	if reflect.DeepEqual(toll, config.TollgateData{}){
		resp.Code = proto.FA_INVALID_DUNGEON_ID
		return
	}
	if player.MaxChallengeId < req.DungeonId {
		resp.Code = proto.FA_CHALLENGE_WRONG_DUNGEON_ID
		return
	}
	if player.Level < config.Param().OpenPvb {
		resp.Code = proto.FA_LEVEL_LESS_THAN
		return
	}
	if err := player.UpdateChallengeTimes(-1);err != nil{
		resp.Code = proto.FA_CHALLENGE_NO_CHANCE
		log.Error(err)
		return
	}

	tasks := &task.Task{PlayerId: player.PlayerId}
	tasks.UpdatePvbDaily()
	tasks.PushTasks()
}


func FinishFightChallenge(sess server.Session, req *pvb.FinishFightChallenge, resp *pvb.FinishFightChallenged){
	resp.Code = proto.OK
	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	toll := config.Tollgate().Get(strconv.Itoa(req.DungeonId))
	if reflect.DeepEqual(toll, config.TollgateData{}){
		resp.Code = proto.FA_INVALID_DUNGEON_ID
		return
	}
	if player.MaxChallengeId < req.DungeonId {
		resp.Code = proto.FA_CHALLENGE_WRONG_DUNGEON_ID
		return
	}
	awards := &award.Awards{}
	if req.Win == WIN {
		//获胜
		nextChallengeId := req.DungeonId + 1
		if reflect.DeepEqual(config.Tollgate().Get(strconv.Itoa(req.DungeonId + 1)), config.TollgateData{}){
			nextChallengeId = req.DungeonId
		}
		player.MaxChallengeId = nextChallengeId
		awards.Drop(toll.FristDropId)
		//任务
		tasks := &task.Task{PlayerId:player.PlayerId}
		tasks.UpdatePvbAch(toll.Point)
		tasks.PushTasks()

	} else {
		//失败
		awards.Drop(toll.DropId)
	}
	awards.SaveAwards(player)
	resp.Items = awards.ConvertItem()
	player.UpdateChallengeId()
	if awards.OnBag(){
		(&push.OnUpdateBag{}).Push(sess.UId())
	}
	if awards.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(player)
	}
}

func buyChallenge(sess server.Session, req *pvb.BuyChallenge, resp *pvb.BuyChallenged){
	resp.Code = proto.OK
	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	price := req.Times * config.Param().PvbPrice
	if player.Diamond < price{
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}

	player.Diamond -= price
	if err := player.UpdateBuyChallenge(req.Times);err != nil{
		resp.Code = proto.FA_CHALLENGE_MAX_TIMES
		log.Error(err)
		return
	}
	resp.Diamond = player.Diamond
	resp.Times = player.ChallengeTimes
}

func resetChallenge(sess server.Session, req *pvb.ResetChallenge, resp *pvb.ResetChallenged) {
	resp.Code = proto.OK
	player := &model.Player{
		PlayerId: sess.UId(),
	}
	if err := player.Get(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST
		return
	}
	player.MaxChallengeId = config.GetPlayerCfg().Dungeon.ChallengeDungeonId
	if err := player.UpdateChallengeId();err != nil {
		resp.Code = proto.FAIL
		log.Error(err)
		return
	}
}





