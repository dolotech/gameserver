package main

import (
	"flag"
	"fmt"
	"gameserver/config"
	_ "gameserver/config"
	_ "gameserver/game/chat"
	_ "gameserver/game/friend"
	_ "gameserver/game/guide"
	_ "gameserver/game/mail"
	_ "gameserver/game/pva"
	_ "gameserver/game/pvb"
	_ "gameserver/game/pve"
	_ "gameserver/game/pvp"
	"gameserver/game/pvp/room"
	"gameserver/game/robot"
	_ "gameserver/game/role"
	_ "gameserver/game/shop"
	_ "gameserver/game/sign"
	_ "gameserver/game/skill"
	_ "gameserver/game/task"
	_ "gameserver/game/user"
	_ "gameserver/game/weapon"
	"gameserver/http"
	lmodel "gameserver/http/model"
	"gameserver/model"
	"gameserver/online"
	_ "gameserver/protocol"
	"gameserver/utils"
	"gameserver/utils/filter"
	"gameserver/utils/log"
	"gameserver/utils/signal"
	"gameserver/utils/socket/server"
	"gameserver/utils/socket/tcp"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/viper"
	"os"
)


func NewRouter() *router.APIBuilder {
	api := router.NewAPIBuilder()
	v1 := api.Party("/", ).AllowMethods(iris.MethodOptions)
	p := pprof.New()
	v1.Post("/login", http.Login)
	v1.Post("/register", http.Register)
	v1.Post("/guestLogin", http.GuestLogin)
	v1.Post("/bindGuestUser", http.BindGuestUser)
	v1.Post("/resetPassword", http.ResetPassword)
	v1.Post("/updateServerId", http.UpdateServerId)
	v1.Post("/servers", http.Servers)
	v1.Any("/debug/pprof", p)
	v1.Any("/debug/pprof/{action:path}", p)
	return api
}
func main() {
	var fileName string
	flag.StringVar(&fileName, "conf", "cfg.toml", "Configuration file to start game")
	flag.Parse()
	if !utils.PathExists(fileName) {
		fmt.Println("conf file not exist !!", fileName)
		return
	}
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: %s \n", err)
	}
	log.Info("conf:", fileName)

	signal.HandleSignals(func(o os.Signal) {
		config.Reload(viper.GetString("general.configDir"))
	})
	// 初始化数据库表格
	model.AutoMigrate()
	lmodel.AutoMigrate()
	// 初始化json表
	config.Load(viper.GetString("general.configDir"))
	// 脏字库，支持多个库
	filter.LoadDicFiles([]string{viper.GetString("general.dic")})

	utils.LoadNames("name.txt")
	app := iris.New()
	app.APIBuilder = NewRouter()
	app.Use(recover.New())
	go app.Run(iris.Addr(":" + viper.GetString("login.port")))

	log.Info("Now listening on: localhost:", viper.GetString("server.port"))
	tcpServer := new(tcp.TCPServer)
	tcpServer.Addr = ":" + viper.GetString("server.port")
	tcpServer.MaxConnNum = 10000
	tcpServer.PendingWriteNum = 300
	tcpServer.MaxMsgLen = 4096
	tcpServer.LittleEndian = false
	tcpServer.NewAgent = func(conn tcp.Conn) tcp.Agent {
		return server.NewAgent(conn, online.GetEvent())
	}
	robot.Robot()    //初始化机器人数据
	robot.InitRoomPool(room.NewRoom)    //初始化pvp房间列表

	tcpServer.Start()
}
