package online

import (
	"gameserver/model"
	"sync"
	"sync/atomic"
)

func NewUserData(playerId uint) *UserData {
	return &UserData{
		UId:  playerId,
		lock: &sync.RWMutex{},
	}
}

func (this *UserData) SetPartners(partners model.Partners) {
	defer this.lock.Unlock()
	this.lock.Lock()
	this.partners = partners.Clone()
}

func (this *UserData) SetPlayer(p model.Player) {
	defer this.lock.Unlock()
	this.lock.Lock()
	this.player = p
}

func (this *UserData) GetPartners() model.Partners {
	defer this.lock.RUnlock()
	this.lock.RLock()
	return this.partners.Clone()
}

func (this *UserData) GetPlayer() model.Player {
	defer this.lock.RUnlock()
	this.lock.RLock()
	return this.player
}

type UserData struct {
	UId                  uint
	lock                 *sync.RWMutex
	player               model.Player
	RefreshRandomFriends int64
	ChatSendMessage      int64
	partners             model.Partners
	room                 uint32
}

func (this *UserData) SetRoom(id uint32) {
	atomic.StoreUint32(&this.room, id)
}

func (this *UserData) GetRoom() uint32 {
	return atomic.LoadUint32(&this.room)
}
