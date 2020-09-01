package http

import (
	"errors"
	"gameserver/http/model"
	"gameserver/http/protoc"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		return []byte(viper.GetString("general.token_key")), nil
	}
}

//cookie验证
func CookieAuth(cookies string) (err error, ok bool) {
	if len(cookies)==0 && cookies == "" {
		err = errors.New("token is null")
		return
	}
	var token *jwt.Token
	token, err = jwt.Parse(cookies, secret())
	if err != nil {
		log.Error("parse token fail")
		return
	}
	data, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("token.Claims fail")
		log.Error(err)
		return
	}
	if !token.Valid {
		err = errors.New("token.Validl")
		log.Error(err)
		return
	}

	id := data["id"].(float64)
	player := &model.User{PlayerId: uint(id)}
	if err = player.GetByuID(); err != nil {
		log.Error(ErrGetPlayer)
		return
	}
	ok = true
	return
}

func Servers(ctx iris.Context) {
	var err error
	if err, _ = CookieAuth(ctx.GetHeader("Cookie")); err != nil {
		log.Error(err)
		return
	}

	var gameServers protoc.GameServers
	if err = ctx.ReadForm(&gameServers); err != nil {
		log.Error(ErrGetParam, err)
		return
	}
	gameServersed := protoc.GameServersed{Code: server.CODE_USE_ERROR}

	if gameServersed.Servers, err = GetGameServers(); err != nil {
		ctx.JSON(gameServersed)
		return
	}
	gameServersed.Code = server.CODE_OK
	ctx.JSON(gameServersed)
}

func GetGameServers() (value []protoc.Server, err error) {
	var serversInfo []protoc.Server
	servers := &model.Servers{}
	if err := servers.Get(); err != nil {
		log.Error(ErrGetServers, err)
		return serversInfo, err
	}
	for _, v := range *servers {
		server := protoc.Server{ServerID: v.ID, Host: v.Addr, Port: v.Port, State: v.Status, Name: v.Name}
		serversInfo = append(serversInfo, server)
	}
	if len(serversInfo) == 0 { //没有服务器列表
		err = ErrNoServerList
		log.Warning(ErrNoServerList)
	}
	return serversInfo, err
}
