package room

import (
	"gameserver/game/robot"
	"gameserver/protocol/push"
	"gameserver/utils"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"math/rand"
	"sync/atomic"
)

type BattleRoom struct {
	roomId uint32
	players [2]robot.BPlayer
	ready   uint32
	Mode    string
	closed  uint32
}



func (this *BattleRoom) RId() uint32{
	return this.roomId
}
func (this *BattleRoom) Get(playerId uint) robot.BPlayer {
	for _, v := range this.players {
		if v.UId() == playerId {
			return v
		}
	}
	return nil
}
func (this *BattleRoom) PushByUId(uid uint, msg interface{}) {
	for _, v := range this.players {
		if v.UId() == uid {
			v.Push(msg)
			return
		}
	}
}
func (this *BattleRoom) Push(msg interface{}) {
	for _, v := range this.players {
		v.Push(msg)
	}
}

func (this *BattleRoom) PushOpp(uid uint, msg interface{}) {
	for _, v := range this.players {
		if v.UId() != uid {
			v.Push(msg)
			return
		}
	}
}

func (this *BattleRoom) Opp(uid uint) robot.BPlayer {
	for _, v := range this.players {
		if v.UId() != uid {
			return v
		}
	}
	return nil
}
func (this *BattleRoom) OppId(uid uint) uint {
	for _, v := range this.players {
		if v.UId() != uid {
			return v.UId()
		}
	}
	return 0
}

func (this *BattleRoom) CheckEnd() uint {
	for _, v := range this.players {
		if !v.Alive() {
			return this.OppId(v.UId())
		}
	}
	return 0
}

func (this *BattleRoom) sendReadyToFight(playerId uint, mapid int) {
	msg := &push.OnUpdateReadyToFight{MapId: mapid, PlayerId: playerId}
	p := this.Opp(playerId).GetPlayer()
	utils.StructAtoB(&msg.Opponent, p)
	par := this.Opp(playerId).GetPartners()
	msg.Opponent.Partners.CopyTo(p, par)
	this.PushByUId(playerId, msg)
}

func (this *BattleRoom) Start() {
	mapid := int(rand.Int31n(100) + 1)
	for _, v := range this.players {
		v.Start(this.roomId)
		this.sendReadyToFight(v.UId(), mapid)
	}
	log.Info("房间创建成功", this)
}

func (this *BattleRoom) genBattlePlayer(b robot.BPlayer) *push.BattlePlayer {
	return &push.BattlePlayer{
		PlayerId: b.UId(),
		Hp:       int(b.Hp()),
		Dp:       int(b.Dp()),
	}
}

func (this *BattleRoom) Ready() {
	if atomic.AddUint32(&this.ready, 1) == 2 {
		// todo 都准备好了，开始游戏
		log.Info("都准备好了，开始游戏", this)
		msg := &push.OnUpdateBattleStart{}
		msg.Left = *this.genBattlePlayer(this.GetPlayer(0))
		msg.Right = *this.genBattlePlayer(this.GetPlayer(1))
		msg.BattleId = int(this.roomId)
		for _, v := range this.players {
			v.AllReady()
		}
		this.Push(msg)
	}
}

func (this *BattleRoom) GameOver(winner uint, isGone int) {
	for _, v := range this.players {
		v.OnGameOver(winner)
	}
	end := &push.OnUpdateBattleEnd{
		WinId:  winner,
		IsGone: isGone,
	}
	this.Push(end)
}


func (this *BattleRoom) AllLost() {
	for _, v := range this.players {
		this.Lost(v.UId())
	}
}
//更新对手断线 PUSH
func (this *BattleRoom) Lost(lostPlayerId uint) {
	log.Info("更新对手断线", lostPlayerId)
	var pushPlayer = this.OppId(lostPlayerId)
	lost := &push.OnUpdateOpponentHadGone{
		PlayerId: pushPlayer,
		FriendId: lostPlayerId,
	}
	this.PushByUId(pushPlayer, lost)
}

// 房间中只要有一人掉线即被调用
func (this *BattleRoom) Close() {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		for _, v := range this.players {
			v.OnClose()
		}
	}
}

func (this *BattleRoom) AllTimeOut() bool {
	count:=0
	for _, v := range this.players {
		if v.TimeOut(){
			count +=1
		}
	}
	return count ==2
}
func (this *BattleRoom) GetPlayer(index int) robot.BPlayer {
	if index < 0 || index > 1 {
		return nil
	}
	return this.players[index]
}
func (this *BattleRoom) AddPlayer(index int, bp robot.BPlayer) {
	if index < 0 || index > 1 {
		return
	}
	this.players[index] = bp
}

func randRobot(robot robot.BPlayer)  {
	hp:=int32(robot.GetPartners().TotalHP(robot.GetPlayer().Level) * 6)
	dp:=int32(robot.GetPartners().TotalDP(robot.GetPlayer().Level) * 6)
	dp = dp + int32(rand.Intn(int(dp/5))) - int32(dp/10)
	robot.HarmDp(dp)
	hp = hp + int32(rand.Intn(int(hp/5))) - int32(hp/10)
	robot.HarmHp(hp)
}

func NewRoom(mode string, battleID uint32, organiger, opponent server.Session)robot.Room {
	// 机器人一定要在右边
	r := &BattleRoom{Mode: mode, roomId: battleID}
	if organiger.UId() < robot.ROBOT {
		robot:=NewRobotPlayer(organiger)
		randRobot(robot)
		r.AddPlayer(1, robot)
		r.AddPlayer(0, NewBattlePlayer(battleID, opponent))

	} else if opponent.UId() < robot.ROBOT{
		robot:=NewRobotPlayer(opponent)
		r.AddPlayer(1, robot)
		randRobot(robot)
		r.AddPlayer(0, NewBattlePlayer(battleID, organiger))
	} else {
		r.AddPlayer(0, NewBattlePlayer(battleID, organiger))
		r.AddPlayer(1, NewBattlePlayer(battleID, opponent))
	}
	return r
}
