package pva

import (
	"gameserver/config"
	"gameserver/game/award"
	"gameserver/game/robot"
	task2 "gameserver/game/task"
	"gameserver/model"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/push"
	"gameserver/protocol/pva"
	"gameserver/protocol/route"
	"gameserver/utils"
	"gameserver/utils/log"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"math/rand"
	"strconv"
	"time"
)

// 竞技塔

const (
	WIN  = 1
	LOSS = 2
)

func init() {
	msg.GetMsg().Reg(route.GetWeeklyBoardInfo, &pva.GetWeeklyBoardInfo{}, &pva.GetWeeklyBoardInfoed{}, setWeeklyBoardInfo)
	msg.GetMsg().Reg(route.ChallengeWeeklyOpponent, &pva.ChallengeWeeklyOpponent{}, &pva.ChallengeWeeklyOpponented{}, challengeWeeklyOpponent)
	msg.GetMsg().Reg(route.FinishChallengeWeeklyOpponent, &pva.FinishChallengeWeeklyOpponent{}, &pva.FinishChallengeWeeklyOpponented{}, finishChallengeWeeklyOpponent)
	msg.GetMsg().Reg(route.BuyWeeklyChallengeTimes, &pva.BuyWeeklyChallengeTimes{}, &pva.BuyWeeklyChallengeTimesed{}, buyWeeklyChallengeTimes)
	msg.GetMsg().Reg(route.RecvWeeklyAwards, &pva.RecvWeeklyAwards{}, &pva.RecvWeeklyAwardsed{}, recvWeeklyAwards)
}

//获取PVA榜信息
func setWeeklyBoardInfo(sess server.Session, req *pva.GetWeeklyBoardInfo, resp *pva.GetWeeklyBoardInfoed) {
	Max := config.Param().PvaShowMax
	Head := 3
	Tail := Max - Head - 1
	resp.Code = proto.OK
	self := &model.Player{PlayerId: sess.UId()}
	if err := self.GetWeeklyPlayer(); err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}

	var totalPlayers model.Players
	if self.WeeklyRank > 0 && self.WeeklyRank <= Max {
		players := model.Players{}
		players.Rank(self.WeeklyRank, Max-1)

		totalPlayers = append(totalPlayers, self)
		for _, v := range totalPlayers {
			if v.WeeklyRank > 0 && v.WeeklyRank <= Max {
				totalPlayers = append(totalPlayers, players...)
			}
		}

		for i := 1; i <= Max; i++ {
			has := false
			for _, v := range totalPlayers {
				if v.WeeklyRank == i {
					has = true
					break
				}
			}
			if !has {
				robot := robot.Pool().Rank(i)
				totalPlayers = append(totalPlayers, robot)
			}
		}
	} else {
		head := model.Players{}
		head.Rank(self.WeeklyRank, 3)
		for _, v := range head {
			if v.WeeklyRank > 0 && v.WeeklyRank <= 3 {
				totalPlayers = append(totalPlayers, v)
			}
		}
		for i := 1; i <= 3; i++ {
			has := false
			for _, v := range head {
				if v.WeeklyRank == i {
					has = true
					break
				}
			}
			if !has {
				robot := robot.Pool().Rank(i)
				totalPlayers = append(totalPlayers, robot)
			}
		}

		near := model.Players{}
		var lower, upper int
		lower = (self.WeeklyRank) - rand.Intn(6) - 5
		upper = (self.WeeklyRank) + rand.Intn(6) + 8
		near.RankNear(self.WeeklyRank, Tail, lower, upper)

		totalPlayers = append(totalPlayers, near...)

		ex := make([]int, 0, len(totalPlayers))
		for _, v := range totalPlayers {
			ex = append(ex, v.WeeklyRank)
		}
		ex = append(ex, self.WeeklyRank)

		robots := robot.Pool().Match(ex, self.Power, Max-len(totalPlayers)-1)
		for _, v := range robots {
			p := v
			totalPlayers = append(totalPlayers, p)
		}

		totalPlayers = append(totalPlayers, self)
	}

	for _, v := range totalPlayers {
		opp := pva.Opponent{}
		utils.StructAtoB(&opp, v)
		resp.Opponents = append(resp.Opponents, opp)
	}

	self.UpdateDayWeeklyPKTimes(0)
	resp.WeeklyRank = self.WeeklyRank
	resp.WeeklyPKTimes = self.DayWeeklyPKTimes
	resp.WeeklyLeftTime = config.Param().PvaRecover - (int(time.Now().Unix()) - self.WeeklyLeftTime) // 获取pK次数倒计时(2小时恢复1个)
}

//挑战PVA场对手
func challengeWeeklyOpponent(sess server.Session, req *pva.ChallengeWeeklyOpponent, resp *pva.ChallengeWeeklyOpponented) {
	resp.Code = proto.OK
	if sess.UId() == req.OpponentId {
		resp.Code = proto.FAIL
		return
	}

	self := model.Player{PlayerId: sess.UId()}
	if err := self.GetWeeklyPlayer(); err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}
	if self.DayWeeklyPKTimes <= 0 {
		resp.Code = proto.WEEKLY_TIMES_NOT_ENOUGH
		log.Error("WEEKLY_TIMES_NOT_ENOUGH")
		return
	}

	if req.OpponentId< robot.ROBOT{
		opp := robot.Pool().Get(req.OpponentId)
		u := opp.UserData().(*online.UserData)
		oppPlayer := u.GetPlayer()
		utils.StructAtoB(&resp.Opponent, oppPlayer)
		partners := u.GetPartners()
		resp.Opponent.Partners.CopyTo(&oppPlayer, &partners)
	} else {
		opp := &model.Player{PlayerId: req.OpponentId}
		if err := opp.GetWeeklyOpp(); err != nil {
			log.Error(err)
			resp.Code = proto.FAIL
			return
		}

		utils.StructAtoB(&resp.Opponent, opp)
		par := &model.Partners{}
		if err := par.Get(req.OpponentId); err != nil {
			resp.Code = proto.FAIL
			return
		}
		resp.Opponent.Partners.CopyTo(opp, par)
	}

	err := self.UpdateDayWeeklyPKTimes(-1)
	if err != nil {
		log.Error(err)
		resp.Code = proto.WEEKLY_TIMES_NOT_ENOUGH
	}
}

//完成PVA场对手的挑战
func finishChallengeWeeklyOpponent(sess server.Session, req *pva.FinishChallengeWeeklyOpponent, resp *pva.FinishChallengeWeeklyOpponented) {
	resp.Code = proto.OK

	self := model.Player{PlayerId: sess.UId()}
	if err := self.GetWeeklyOpp(); err != nil {
		log.Error(err)
		resp.Code = proto.FAIL
		return
	}

	awards, ok := config.LdtEveryAward().Get(strconv.Itoa(req.Result))
	if !ok {
		log.Error("no award")
		resp.Code = proto.FAIL
		return
	}

	saveAward := award.Awards{}
	saveAward.Drop(awards.Award)
	saveAward.SaveAwardsById(sess.UId())
	resp.Items = saveAward.ConvertItem()

	if req.Result == LOSS {
		self.AddWeeklyPKLoses()
	} else if req.Result == WIN {
		self.AddWeeklyPKWins()
		//累计胜利任务
		tasks := &task2.Task{PlayerId:self.PlayerId}
		tasks.UpdatePvaAch()
		tasks.PushTasks()
	}
	var opp *model.Player
	if req.OpponentId < robot.ROBOT{
		sess := robot.Pool().Get(req.OpponentId)
		p := sess.UserData().(*online.UserData).GetPlayer()
		opp = &p
	} else {
		opp = &model.Player{PlayerId: req.OpponentId}
		if err := opp.GetWeeklyOpp(); err != nil {
			log.Error(err)
			resp.Code = proto.FAIL
			return
		}
	}

	if req.Result == WIN {
		//调换排行
		if opp.WeeklyRank < self.WeeklyRank {
			opp.WeeklyRank, self.WeeklyRank = self.WeeklyRank, opp.WeeklyRank
			self.UpdateRank()
			opp.UpdateRank()
		}
	} else if req.Result == LOSS {
		//调换排行
		if opp.WeeklyRank > self.WeeklyRank {
			opp.WeeklyRank, self.WeeklyRank = self.WeeklyRank, opp.WeeklyRank
			self.UpdateRank()
			opp.UpdateRank()
		}
	}
	resp.NewRank = self.WeeklyRank

	utils.StructAtoB(&resp.Opponent, opp)
	//任务
	task := &task2.Task{PlayerId:self.PlayerId}
	task.UpdatePvaDaily()
	task.PushTasks()
}

// 领取获取挑战PVA的奖励
func recvWeeklyAwards(sess server.Session, req *pva.RecvWeeklyAwards, resp *pva.RecvWeeklyAwardsed) {
	resp.Code = proto.OK

	p := &model.Player{PlayerId: sess.UId()}
	if err := p.Get(); err != nil {
		resp.Code = proto.FAIL
		return
	}
	if int(time.Now().Unix())-p.LeftWeeklyAwardTime < config.Param().PvaAwardsTime {
		resp.Code = proto.FAIL
		return
	}
	a := config.LdtLevelAward().Get(p.LastWeeklyRank)
	if a == 0 {
		resp.Code = proto.FAIL
		return
	}
	saveAward := award.Awards{}
	saveAward.Drop(a)
	saveAward.SaveAwards(p)
	t := int(time.Now().Unix())
	p.RecvWeeklyAward(t)

	resp.LeftWeeklyAwardTime = config.Param().PvaAwardsTime

	if saveAward.OnBag() {
		(&push.OnUpdateBag{}).Push(p.PlayerId)
	}
	if saveAward.OnPlayer() {
		(&push.OnUpdatePlayer{}).Push(p)
	}
}

// 购买挑战PVA场次数
func buyWeeklyChallengeTimes(sess server.Session, req *pva.BuyWeeklyChallengeTimes, resp *pva.BuyWeeklyChallengeTimesed) {
	resp.Code = proto.OK

	p := &model.Player{PlayerId: sess.UId()}
	if err := p.Get(); err != nil {
		resp.Code = proto.FAIL
		return
	}
	if p.DayWeeklyPKTimes >= config.Param().PvaMaxNum {
		log.Error("PvaMaxNum:", p.DayWeeklyPKTimes, config.Param().PvaMaxNum)
		resp.Code = proto.FAIL
		return
	}

	price := req.Times * config.Param().PvaPrice
	if p.Diamond < price {
		log.Error("FA_NOT_ENOUGH_DIAMOND:", p.Diamond, price)
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}

	if p.UpdateDiamond(-price) != nil {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		return
	}
	p.DayWeeklyPKTimes = p.DayWeeklyPKTimes + req.Times
	p.Diamond = p.Diamond - price

	p.UpdateDayWeeklyPKTimes(req.Times)

	resp.ChallengeTimes = p.DayWeeklyPKTimes

	(&push.OnUpdatePlayer{}).Push(p)
}
