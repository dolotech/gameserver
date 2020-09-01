package http

import (
	"gameserver/http/model"
	"gameserver/http/protoc"
	"gameserver/protocol"
	"gameserver/utils"
	"gameserver/utils/filter"
	"gameserver/utils/log"
	"gameserver/utils/socket/server"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

func Login(ctx iris.Context) {
	var login protoc.Login
	resp := protoc.Logined{Code: server.CODE_USE_ERROR}

	if err := ctx.ReadForm(&login); err != nil {
		log.Error(ErrGetParam, err)
		ctx.JSON(resp)
		return
	}

	player := &model.User{Username: login.Username}
	if err := player.GetByUsername(); err != nil {
		log.Error(ErrGetPlayer, err)
		resp.Code = protocol.FA_USER_NOT_EXIST
		ctx.JSON(resp)
		return
	}

	if match := CheckPasswordHash(login.Password, player.Password); match {
		token, err := GenerateToken(uint(player.PlayerId))
		if err != nil {
			log.Error(ErrGenerateToken, err)
			ctx.JSON(resp)
			return
		}
		resp := &protoc.Logined{
			Code:         server.CODE_OK,
			Token:        token,
			UID:          strconv.Itoa(int(player.PlayerId)),
			LastServerID: 0,
			Servers:      nil,
		}
		if resp.Servers, err = GetGameServers(); err == nil {
			resp.LastServerID = getLastServer(resp, player)
		}
		ctx.JSON(resp)
		return
	}
	ctx.JSON(resp)
}

func Register(ctx iris.Context) {
	var reg protoc.Register
	if err := ctx.ReadForm(&reg); err != nil {
		log.Error(ErrGetParam, err)
		return
	}
	str := []rune(reg.Username)
	filter.FilterText("dic.txt", str, []rune{}, '*')
	reg.Username = string(str)

	password := HashPassword(reg.Password)
	player := &model.User{
		Username:   reg.Username,
		Password:   password,
		CreateTime: (int)(time.Now().Unix()),
		LoginTime:  (int)(time.Now().Unix()),
		RegIP:      int(utils.InetToaton(ctx.RemoteAddr())),
	}
	res := doRegist(player)
	ctx.Write(res)
}

func doRegist(player *model.User) []byte {
	var res []byte
	registerd := &protoc.Logined{}
	if player.Username != "" {
		err := player.GetByUsername()
		if err == nil {
			log.Error(ErrPlayerExist, player.Username)
			registerd.Code = protocol.FA_USERNAME_EXIST
			res, _ = registerd.Encode()
			return res
		}
	}

	if err := player.CreateUser(); err == nil {
		token, _ := GenerateToken(player.PlayerId)
		registerd.Code = server.CODE_OK
		registerd.Token = token
		registerd.UID = (string)(player.PlayerId)
		if registerd.Servers, err = GetGameServers(); err == nil {
			registerd.LastServerID = getLastServer(registerd, player)
		}

	} else {
		log.Error(ErrAddPlayer, err)
		registerd.Code = server.CODE_USE_ERROR
	}

	res, _ = registerd.Encode()
	return res
}

func GuestLogin(ctx iris.Context) {
	guest := &protoc.GuestLogin{}
	if err := ctx.ReadForm(guest); err != nil {
		log.Error(ErrGetParam, err)
	}

	player := &model.User{MacSn: guest.MacSn}
	var res []byte
	if err := player.GetByMacSn(); err != nil {
		res = doRegist(player)
	} else {
		token, err := GenerateToken(player.PlayerId)
		if err != nil {
			log.Error(ErrGenerateToken, err)
			return
		}
		logined := &protoc.Logined{
			Code:         server.CODE_OK,
			Token:        token,
			UID:          strconv.Itoa(int(player.PlayerId)),
			LastServerID: 0,
			Servers:      nil,
		}
		if logined.Servers, err = GetGameServers(); err == nil {
			logined.LastServerID = getLastServer(logined, player)
		}
		res, _ = logined.Encode()
	}
	ctx.Write(res)
}

func getLastServer(logined *protoc.Logined, user *model.User) uint {
	var serverID uint
	if user.ServerId > 0 {
		serverID = user.ServerId
	} else {
		if len(logined.Servers) > 0 {
			serverID = logined.Servers[len(logined.Servers)-1].ServerID
		} else {
			log.Error("没有找到可用服务器")
		}
	}
	return serverID
}

func BindGuestUser(ctx iris.Context) {
	bind := &protoc.BindGuestUser{}
	if err := ctx.ReadForm(bind); err != nil {
		log.Error(ErrGetParam, err)
		return
	}

	if !utils.LegalName(bind.Username, 12) {
		log.Error(ErrGetParam)
		return
	}

	if !filter.Valid("dic.txt", []rune(bind.Username), []rune{}, '*') {
		log.Error(ErrGetParam)
		return
	}

	binded := &protoc.BindGuestUsered{Code: server.CODE_OK}

	player := &model.User{
		MacSn: bind.MacSn,
	}
	if err := player.GetByMacSn(); err == nil {
		player.Username = bind.Username
		player.Password = HashPassword(bind.Password)
		if err = player.BindGuestUser(); err != nil {
			binded.Code = server.CODE_OK
		}
	}
	res, _ := binded.Encode()
	ctx.Write(res)
}

func UpdateServerId(ctx iris.Context) {
	req := &protoc.UpdateServerId{}
	if err := ctx.ReadForm(req); err != nil {
		log.Error(ErrGetParam, err)
		return
	}
	player := &model.User{Username: req.Username}
	if err := player.GetByUsername(); err != nil {
		log.Error(ErrGetUsername, err)
		return
	}
	if err := player.UpdateServerId(req.ServerId); err != nil {
		resp := &protoc.UpdateServerIded{Code: server.CODE_USE_ERROR}
		ctx.JSON(resp)
		return
	}
	reseted := &protoc.UpdateServerIded{Code: server.CODE_OK}
	ctx.JSON(reseted)
	return
}
func ResetPassword(ctx iris.Context) {
	reset := &protoc.ResetPassword{}
	if err := ctx.ReadForm(reset); err != nil {
		log.Error(ErrGetParam, err)
		return
	}
	player := &model.User{Username: reset.Username}
	if err := player.GetByUsername(); err != nil {
		log.Error(ErrGetUsername, err)
		return
	}
	if ok := CheckPasswordHash(reset.OldPassword, player.Password); ok {
		player.UpdatePassword(HashPassword(reset.NewPassword))
		reseted := &protoc.ResetPassworded{Code: server.CODE_OK}
		ctx.JSON(reseted)
		return
	}
	reseted := &protoc.ResetPassworded{Code: server.CODE_USE_ERROR}
	ctx.JSON(reseted)
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error(ErrHashPass, err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(playerId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   playerId,
		"rand": rand.Intn(1000),
		"exp":  time.Now().Add(time.Second * viper.GetDuration("general.token_expire")).Unix(),
	})
	return token.SignedString([]byte(viper.GetString("general.token_key")))
}
