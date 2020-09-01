package room

import (
	"gameserver/game/robot"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"math/rand"
	"runtime/debug"
	"time"
)

func NewRobotPlayer(sess server.Session) robot.BPlayer {
	r := &RobotPlayer{}
	r.BattlePlayer.session = sess
	r.close = make(chan struct{})
	return r
}

type RobotPlayer struct {
	BattlePlayer
	close chan struct{}
}

func (this *RobotPlayer) OnGameOver(win uint) {

}
func (this *RobotPlayer) AllReady() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error(string(debug.Stack()), e)
			}
		}()
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(5000)+2000))
			select {
			case <-this.close:
				return
			default:
				if rand.Intn(100) < 20 {
					this.attack()
					this.attack()
				} else {
					this.attack()
				}
			}
		}
	}()
}

func (this *RobotPlayer) attack() {
	//todo
	//  AttackRatio 攻击类型，0为普攻，1为爆击，2为致命
	// FighterIndex  1-10
	// skillId为0是普攻，1是武器技能
	roomid := this.roomId
	if roomid > 0 {
		r := robot.Get().Get(roomid)
		if r == nil {
			return
		}
		index := rand.Intn(5) + 6

		attackRatio := 0
		ratio := rand.Intn(100)
		if ratio <= 10 {
			attackRatio = 1
		} else if ratio > 10 && ratio <= 20 {
			attackRatio = 2
		}

		skillId := rand.Intn(2)

		Fire(this, r, index, skillId, attackRatio)
	}
}
func (this *RobotPlayer) OnClose() {
	this.BattlePlayer.OnClose()
	select {
	case <-this.close:
	default:
		robot.Get().Cancel(this.UId()) // 归还机器人需要通过统一的通道
		close(this.close)
	}
}
