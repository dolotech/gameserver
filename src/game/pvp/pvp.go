package pvp

// 好友对战
import (
	"gameserver/config"
	"gameserver/online"
	proto "gameserver/protocol"
	"gameserver/protocol/pvp"
	"gameserver/protocol/route"
	msg "gameserver/utils/socket/message"
	"gameserver/utils/socket/server"
	"time"
)

// 游戏消息注册
func init() {
	msg.GetMsg().Reg(route.EnterArena, &pvp.EnterArena{}, &pvp.EnterArenaed{}, enterArena)
	msg.GetMsg().Reg(route.SetFighterBuff, &pvp.SetFighterBuff{}, &pvp.SetFighterBuffed{}, setFighterBuff)
	msg.GetMsg().Reg(route.RecvArenaRankAward, &pvp.RecvArenaRankAward{}, &pvp.RecvArenaRankAwarded{}, recvArenaRankAward)
	msg.GetMsg().Reg(route.BuyArenaPKTimes, &pvp.BuyPvpChallengeTimes{}, &pvp.BuyPvpChallengeTimesed{}, buyArenaPKTimes)
}

//进入PVP
func enterArena(sess server.Session, req *pvp.EnterArena, resp *pvp.EnterArenaed) {
	resp.Code = proto.OK
	self := sess.UserData().(*online.UserData).GetPlayer()
	self.GetPvpOpp()

	resp.PvpRank = self.PvpRank + 1
	if resp.PvpRank == 0 {
		resp.PvpRank += 1
	}
	resp.PvpScore = self.PvpScore
	resp.PkTimes = self.DayArenaPKTimes
	resp.WinNum = self.ContinuousArenaPKWins //pvp连赢次数
	leftArenaAwardTime := config.Param().PvpRecover - (int(time.Now().Unix()) - self.LastArenaAwardTime)
	if leftArenaAwardTime < 0 {
		leftArenaAwardTime = 0
	}
	resp.LeftArenaAwardTime = leftArenaAwardTime
}

//PVP设置自己的BUFF添加跟删除
func setFighterBuff(sess server.Session, req *pvp.EnterArena, resp *pvp.EnterArenaed) {
	resp.Code = proto.OK

}

//PVP奖励领取
func recvArenaRankAward(sess server.Session, req *pvp.RecvArenaRankAward, resp *pvp.RecvArenaRankAwarded) {
	resp.Code = proto.OK

	self := sess.UserData().(*online.UserData).GetPlayer()
	self.RecvPvpAward(int(time.Now().Unix()))

	resp.LeftArenaAwardTime = config.Param().PvpRecover
}

//PVP购买次数
func buyArenaPKTimes(sess server.Session, req *pvp.BuyPvpChallengeTimes, resp *pvp.BuyPvpChallengeTimesed) {
	resp.Code = proto.OK
	self := sess.UserData().(*online.UserData).GetPlayer()
	self.GetDiamondAndArenaPK()

	if self.DayArenaPKTimes >= config.Param().PvpMaxNum {
		resp.PkTimes = self.DayArenaPKTimes
		resp.Diamond = self.Diamond
		resp.Code = proto.FAIL
		return
	}

	price := req.Times * config.Param().PvpPrice
	if self.Diamond < price {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		resp.PkTimes = self.DayArenaPKTimes
		resp.Diamond = self.Diamond
		return
	}

	if self.UpdateDiamond(-price) != nil {
		resp.Code = proto.FA_NOT_ENOUGH_DIAMOND
		resp.PkTimes = self.DayArenaPKTimes
		resp.Diamond = self.Diamond
		return
	}
	self.AddDayPvpPKTimes(req.Times)
	self.DayArenaPKTimes = self.DayArenaPKTimes + req.Times
	self.Diamond = self.Diamond - price

	resp.PkTimes = self.DayArenaPKTimes
	resp.Diamond = self.Diamond

}
