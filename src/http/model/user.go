package model

import "gameserver/utils/db"

// 用户数据
type User struct {
	PlayerId   uint   `gorm:"primary_key;column:playerId;AUTO_INCREMENT;COMMENT:'用户游戏唯一标识'" json:"playerId"`
	Password   string `gorm:"type:varchar(200);column:password" json:"password"`
	MacSn      string `gorm:"type:varchar(50);column:macSn" json:"macSn"`
	CreateTime int    `gorm:"type:int(11);column:createTime;default:0" json:"createTime"`
	RegIP      int    `gorm:"type:int(11);column:regIP;default:0" json:"regIP"`
	LoginTime  int    `gorm:"type:int(11);column:loginTime;default:0" json:"loginTime"`
	Username   string `gorm:"type:varchar(30);column:username" json:"username"`
	ServerId   uint   `gorm:"type:int(4);column:serverId;default:0" json:"serverId"` // 用户最后一次选择的服务器id
}

func (this *User) GetByUsername() error {
	return db.LDB().Model(this).Where("username = ?", this.Username).Limit(1).Find(this).Error
}
func (this *User) GetByuID() error {
	return db.LDB().Model(this).Where("playerId = ?", this.PlayerId).Limit(1).Find(this).Error
}
func (this *User) GetByMacSn() error {
	return db.LDB().Model(this).Where("macSn = ?", this.MacSn).Find(this).Error
}

func (this *User) BindGuestUser() error {
	m := map[string]interface{}{
		"username": this.Username,
		"password": this.Password,
	}
	return db.LDB().Model(this).Updates(m).Error
}

func (this *User) UpdateServerId(serverId int) error {
	return db.LDB().Model(this).Updates(map[string]interface{}{"serverId": serverId}).Error
}

//更新密码
func (this *User) UpdatePassword(newPass string) error {
	return db.LDB().Model(this).Updates(map[string]interface{}{"password": newPass}).Error
}

func (this *User) CreateUser() error {
	return db.LDB().Create(this).Error
}
