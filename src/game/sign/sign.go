package sign

import (
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/game/task"
	"gameserver/model"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/route"
	"gameserver/protocol/sign"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"strconv"
	"time"
)

// 消息注册
func init() {
	msg.GetMsg().Reg(route.RecvSignAward, &sign.RecvSignAward{}, &sign.RecvSignAwarded{}, recvSignAwardCb)
}

// 请求签到奖励
func recvSignAwardCb(sess server.Session, req *sign.RecvSignAward, resp *sign.RecvSignAwarded) {
	resp.Code = proto.OK
	p := &model.Player{PlayerId: sess.UId()}
	if err := p.GetById(); err != nil {
		resp.Code = proto.FA_USER_NOT_EXIST //FA_INVALID_PLAYER
		return
	}
	if p.GetMonthTodaySign() != nil {
		resp.Code = proto.FAIL
		return
	}

	if p.TodaySignTime > 0 {
		resp.Code = proto.FAIL //FA_TODAY_HAD_SIGN
		return
	}

	monthSignDays := p.MonthSignDays
	monthSignDays++
	config := config.SignAward().Get(strconv.Itoa(monthSignDays))
	if config.DropId < 1 {
		log.Warning("recvSignAward, playerId: %d, monthSignDays: %d", p.PlayerId, monthSignDays)
		resp.Code = proto.FAIL
		return
	}
	signAward := &award.Awards{}
	signAward.Drop(config.DropId)
	signAward.SaveAwards(p)
	if signAward.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(p)
	}
	if signAward.OnBag() {
		(&push.OnUpdateBag{}).Push(sess.UId())
	}
	//log.Debug("recvSignAward, playerId: %d, monthSignDays: %d, DropId:%d,award: %d", p.PlayerId, monthSignDays, config.DropId, signAward)
	p.UpdateMonthTodaySign(monthSignDays, time.Now().Unix())
	//更新签到成就任务
	signTask:= &task.Task{PlayerId: p.PlayerId}
	signTask.UpdatePlayerSign()
	signTask.PushTasks()
}