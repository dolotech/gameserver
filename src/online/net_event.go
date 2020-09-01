package online

import (
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"strconv"
	"sync"
)

type CloseCall interface {
	OnClose(server.Session)

}

func GetEvent() *netevent {
	eventOne.Do(func() {
		event = &netevent{onClose: sync.Map{}}

	})
	return event
}

var event *netevent
var eventOne sync.Once

type netevent struct {
	onClose sync.Map
}

// 握手时登陆验证
func (this *netevent) OnAuth(sess server.Session, to string) (uint, bool) {
	var id float64
	i, err := strconv.Atoi(to)
	if err == nil && viper.GetString("general.token_key") == "" {
		id = float64(i)
	} else {
		token, err := jwt.Parse(to, secret())
		if err != nil {
			log.Error(err)
			return 0, false
		}
		data, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error(" token.Claims fail")
			return 0, false
		}
		if !token.Valid {
			log.Error(" token.Validl")
			return 0, false
		}
		id = data["id"].(float64)
	}

	// 登陆成功绑定用户id和推送服务
	// todo 如果在前端登陆entrygame 接口请求返回前，服务端下发推送前端会不会报错？
	sess.SetUid(uint(id))
	userdata := NewUserData(sess.UId())
	sess.SetUserData(userdata)
	Get().Set(sess)
	return uint(id), true
}

func (this *netevent) RemoveClose(f CloseCall) {
	this.onClose.Delete(f)
}

func (this *netevent) AddClose(f CloseCall) {
	this.onClose.Store(f, struct{}{})
}


//玩家离线了
func (this *netevent) OnClose(sess server.Session) {
	Get().Del(sess.UId())
	this.onClose.Range(func(key, value interface{}) bool {
		c, ok := key.(CloseCall)
		if ok {
			c.OnClose(sess)
		}
		return true
	})
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		return []byte(viper.GetString("general.token_key")), nil
	}
}
