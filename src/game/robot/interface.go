package robot

import (
	"gameserver/model"
	"gameserver/utils/socket/server"
)

type Room interface {
	Ready()
	GetPlayer(int) BPlayer
	AddPlayer(int, BPlayer)
	Start()
	GameOver(uint, int)
	AllLost()
	Lost(lostPlayerId uint)
	Close()
	RId() uint32
	Get(playerId uint) BPlayer
	PushByUId(uid uint, msg interface{})
	Push(msg interface{})
	PushOpp(uid uint, msg interface{})
	Opp(uid uint) BPlayer
	OppId(uid uint) uint
	CheckEnd() uint
	AllTimeOut()bool
}

type BPlayer interface {
	UId() uint
	Start(uint32)
	OnClose()
	Push(interface{})
	Power() int
	Alive() bool
	HarmHp(int32) int32
	HarmDp(int32) int32
	GetPlayer() *model.Player
	GetPartners() *model.Partners
	GetSession() server.Session
	TimeOut() bool
	OnGameOver(uint)
	Hp() int32
	Dp() int32
	AllReady()
}
