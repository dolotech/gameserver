package robot

import (
	"gameserver/online"
	"gameserver/protocol/push"
	"gameserver/utils"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"math"
	"runtime/debug"
	"sync"
	"time"
)

var pvpRoomPool *roomPool
const RoomTimeOut = 300 // 不进行操作的房间超时
const ROBOT = 50000     //

func InitRoomPool(f func(string, uint32, server.Session, server.Session) Room) {
	pvpRoomPool = &roomPool{
		matchPool:   make(map[uint]*Match, 1000),
		remove:      make(chan uint, 100),
		match:       make(chan *Match, 1000),
		ticker:      time.NewTicker(time.Millisecond * 1000),
		cleanTicker: time.NewTicker(time.Minute * 2),
		creator:     f,
	}
	online.GetEvent().AddClose(pvpRoomPool)
	go pvpRoomPool.start()
	go pvpRoomPool.cleanRoom()
}

func Get() *roomPool {
	return pvpRoomPool
}

type roomPool struct {
	matchPool   map[uint]*Match
	roomPool    sync.Map
	remove      chan uint
	match       chan *Match
	id          uint32
	ticker      *time.Ticker
	cleanTicker *time.Ticker
	creator     func(string, uint32, server.Session, server.Session) Room
}

func (this *roomPool) cleanRoom() {
	defer func() {
		if e := recover(); e != nil {
			log.Error(string(debug.Stack()), e)
			time.Sleep(time.Second)
			go this.cleanRoom()
		}
	}()
	for {
		select {
		case <-this.cleanTicker.C:
			this.roomPool.Range(func(key, value interface{}) bool {
				if room, ok := value.(Room); ok {
					if room.AllTimeOut() {
						room.AllLost()
						this.Del(room.RId())
					}
				}
				return true
			})
		}
	}
}
func (this *roomPool) start() {
	defer func() {
		if e := recover(); e != nil {
			log.Error(string(debug.Stack()), e)
		}
	}()
	for {
		select {
		case <-this.ticker.C:
			for _, v := range this.matchPool {
				if v.CreateTime+int64(v.Rand )<= time.Now().Unix() {
					if v.Mode == "pvf" {
						log.Info("pvf 超时")
						// 匹配超时从匹配移除
						v.Timeout()
						delete(this.matchPool, v.Organiger.UId())
					} else if v.Mode == "pvp" {
						if v.Organiger.UId() < ROBOT { // 机器人超时就退出
							delete(this.matchPool, v.Organiger.UId())
							continue
						}

						//  给玩家匹配机器人
						sess := Pool().Search(v.Organiger.UId(), v.Power)
						log.Info("给玩家匹配机器人", sess)
						if sess != nil {
							// 机器人不能匹配机器人
							this.matchSucess(NewMatch(sess, v.Mode, 0, 0), v)
							continue
						} else { // 没有匹配到适合的机器人

							log.Info("没有匹配到适合的机器人")
							v.Timeout()
							delete(this.matchPool, v.Organiger.UId())
						}
					}
				}
			}
			this.matching()
		case m := <-this.match:
			if m.Mode == "pvf" {
				if _, ok := this.matchPool[m.Organiger.UId()]; !ok {
					if _, ok := this.matchPool[m.Opponent]; !ok { //判断是同意接受挑战方，就不再发推送给对手
						msg := &push.OnUpdatePvfFigthAlert{}
						utils.StructAtoB(&msg.Friend, m.Organiger.UserData().(*online.UserData).GetPlayer()) // todo UserData里面的Player数据不实时
						msg.IsOnline = 1
						online.Get().Push(m.Opponent, msg)
					}
				}
			}

			this.matchPool[m.Organiger.UId()] = m
			//this.matching()
		case id := <-this.remove:
			log.Info("取消", id)
			if _, ok := this.matchPool[id]; ok {
				delete(this.matchPool, id)
			} else {
				/*for _, v := range this.matchPool {
					// 对手不同意pk
					if id > 0 && v.Opponent == id {
						delete(this.matchPool, v.Organiger.UId())
						v.Timeout()
					}
				}*/
			}
		}
	}
}

func (this *roomPool) matchByPower(first uint, power int) *Match {
	var score float32 = 1.1
	for ; score < 1.6; {
		for uidkey, vo := range this.matchPool {
			if first == uidkey {
				continue
			}
			if power > vo.Power && int(float32(vo.Power)*score) >= power {
				return vo
			} else if power < vo.Power && int(float32(power)*score) >= vo.Power {
				return vo
			} else if power == vo.Power {
				return vo
			}
		}
		score += 0.1
	}
	return nil
}

func (this *roomPool) matching() {
	for uid, v := range this.matchPool {
		if v.Opponent > 0 {
			if vo, ok := this.matchPool[v.Opponent]; ok && v != vo {
				this.matchSucess(v, vo)
			}
		} else if v.Mode == "pvp" {
			if target := this.matchByPower(uid, v.Power); target != nil {
				this.matchSucess(v, target)
			}
		}
	}
}

func (this *roomPool) matchSucess(first, two *Match) {
	if first != nil && two != nil {
		this.id += 1
		if this.id >= math.MaxUint32 {
			this.id = 1
		}
		r := this.creator(first.Mode, this.id, first.Organiger, two.Organiger)
		this.roomPool.Store(r.RId(), r)

		// 匹配成功从匹配移除
		delete(this.matchPool, first.Organiger.UId())
		delete(this.matchPool, two.Organiger.UId())

		r.Start()
	}
}
func (this *roomPool) Match(match *Match) {
	if match == nil {
		return
	}
	this.match <- match
}

// 掉线和或者玩家主动取消匹配
func (this *roomPool) Cancel(organiger uint) {
	this.remove <- organiger
}

func (this *roomPool) Del(battleId uint32) {
	if room := this.Get(battleId); room != nil {
		room.Close()
	}
	this.roomPool.Delete(battleId)
}
func (this *roomPool) Get(battleId uint32) Room {
	if r, ok := this.roomPool.Load(battleId); ok {
		return r.(Room)
	}
	return nil
}

// 有玩家离线了
func (this *roomPool) OnClose(sess server.Session) {
	this.Cancel(sess.UId())
	if u, ok := sess.UserData().(*online.UserData); ok {
		if roomid := u.GetRoom(); roomid > 0 {
			if room := this.Get(roomid); room != nil {
				log.Info("有玩家离线了roomid ", roomid, sess.UId())
				room.GameOver(room.OppId(sess.UId()), 1)
			}
			this.Del(roomid)
		}
	}
}
