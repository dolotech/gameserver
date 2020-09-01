package robot

import (
	"gameserver/protocol/push"
	"gameserver/utils/socket/server"
	"math/rand"
	"time"
)

type Match struct {
	Mode       string //对战模式 pvf || pvp || no
	Organiger  server.Session
	Opponent   uint  // 指定对手匹配一定是好友实时对战了
	CreateTime int64 // 匹配开始时间
	Rand       int   // 匹配随机时间
	Power      int
}

// 匹配超时被调用
func (this *Match) Timeout() {
	to := &push.OnUpdateFindOverTime{PlayerId: this.Organiger.UId()}
	this.Organiger.Push(to)
}
func NewMatch(organiger server.Session, mode string, opponent uint, power int) *Match {
	r := &Match{
		Power:      power,
		Mode:       mode,
		Opponent:   opponent,
		Organiger:  organiger,
		CreateTime: time.Now().Unix(),
		Rand:       rand.Intn(7)+3,
	}
	return r
}
