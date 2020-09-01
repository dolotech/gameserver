package room

import (
	"gameserver/game/robot"
	"gameserver/model"
	"gameserver/online"
	"gameserver/utils/socket/server"
	"sync/atomic"
	"time"
)

type BattlePlayer struct {
	roomId  uint32
	hp      int32
	dp      int32
	session server.Session
	timeout int64 // 记录战斗操作超时
}

func (this *BattlePlayer) Push(msg interface{}) {
	this.session.Push(msg)
}
func (this *BattlePlayer) UId() uint {
	return this.session.UId()
}

func (this *BattlePlayer) GetPartners() *model.Partners {
	ps := this.session.UserData().(*online.UserData).GetPartners()
	return &ps
}
func (this *BattlePlayer) GetPlayer() *model.Player {
	p := this.session.UserData().(*online.UserData).GetPlayer()
	return &p
}
func (this *BattlePlayer) GetSession() server.Session {
	return this.session
}
func (this *BattlePlayer) TimeOut() bool {
	return time.Now().Unix()-atomic.LoadInt64(&this.timeout) > robot.RoomTimeOut
}

func (this *BattlePlayer) Start(roomid uint32) {
	atomic.StoreInt64(&this.timeout, time.Now().Unix())
	this.roomId = roomid
	this.session.UserData().(*online.UserData).SetRoom(roomid)
}

func (this *BattlePlayer) OnGameOver(win uint) {
	this.GetPlayer().AddDayPvpPKTimes(-1) // 扣除Pk次数
	if win == this.UId() {
		this.GetPlayer().AddContinuousArenaPKWins(1) // 连胜次数加1
	} else {
		this.GetPlayer().UpdateContinuousArenaPKWins(0) // 重置连胜次数
	}
}

func (this *BattlePlayer) OnClose() {
	this.session.UserData().(*online.UserData).SetRoom(0)
}
func (this *BattlePlayer) Power() int {
	return this.GetPartners().CalcuPower(this.GetPlayer().Level)
}

func (this *BattlePlayer) Alive() bool {
	return atomic.LoadInt32(&this.hp) > 0
}
func (this *BattlePlayer) Hp() int32 {
	return atomic.LoadInt32(&this.hp)
}
func (this *BattlePlayer) Dp() int32 {
	return atomic.LoadInt32(&this.hp)
}

func (this *BattlePlayer) HarmHp(hp int32) int32 {
	atomic.StoreInt64(&this.timeout, time.Now().Unix())
	return atomic.AddInt32(&this.hp, hp)
}

func (this *BattlePlayer) HarmDp(dp int32) int32 {
	return atomic.AddInt32(&this.dp, dp)
}
func (this *BattlePlayer) AllReady() {

}
func NewBattlePlayer(rid uint32, sess server.Session) (battlePlayer robot.BPlayer) {
	u := sess.UserData().(*online.UserData) // FindArenaOpponent 匹配协议已经获取到最新的玩家数据
	player := u.GetPlayer()
	partners := u.GetPartners()
	battlePlayer = &BattlePlayer{
		roomId:  rid,
		session: sess,
	}
	battlePlayer.HarmHp(int32(partners.TotalHP(player.Level) * 6))
	battlePlayer.HarmDp(int32(partners.TotalDP(player.Level) * 6))
	return battlePlayer
}
