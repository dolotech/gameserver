package online

import (
	"gameserver/utils/socket"
	"gameserver/utils/socket/server"
	"sync"
)

var onlineOne sync.Once
var online *onlines

// 玩家数据列表
// 列表内玩家的数据改变需要自行维护,目前只保证在线的状态是准确的
func Get() *onlines {
	onlineOne.Do(func() {
		online = &onlines{}
	})
	return online
}

type onlines struct {
	socket.SessionMap
}

func (this *onlines) Online(playerId uint) bool {
	_, ok := this.Load(playerId)
	return ok
}

// after 为true 时本消息会在当前的玩家请求后推送，否则立马推送
func (this *onlines) BCAll(msg interface{}, after ...bool) {
	this.Range(func(key uint, value server.Session) bool {
		value.Push(msg,  after...)
		return true
	})
}

// after 为true 时本消息会在当前的玩家请求后推送，否则立马推送
func (this *onlines) BC(playerIds []uint, msg interface{}, after ...bool) {
	for _, v := range playerIds {
		if p, ok := this.Load(v); ok {
			p.Push(msg,  after...)
		}
	}
}

// after 为true 时本消息会在当前的玩家请求后推送，否则立马推送
func (this *onlines) Push(playerId uint, msg interface{}, after ...bool) {
	if p, ok := this.Load(playerId); ok {
		p.Push(msg, after...)
	}
}
func (this *onlines) Del(playerId uint) {
	if playerId > 0 {
		this.Delete(playerId)
	}
}
func (this *onlines) Set(session server.Session) {
	if session.UId() > 0 {
		if oldSess, ok := this.Load(session.UId()); ok {
			oldSess.Close() // 挤号
		}
		this.Store(session.UId(), session)
	}
}
func (this *onlines) Get(playerId uint) server.Session {
	if p, ok := this.Load(playerId); ok {
		return p
	}
	return nil
}
