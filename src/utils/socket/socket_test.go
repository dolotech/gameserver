package socket

import (
	"gameserver/utils/socket/server"
	"gameserver/utils/socket/tcp"
	"testing"
)

/*
type Agent struct {
	conn *tcp.tcpConn
}
func (this *Agent)OnClose() ()  {
	fmt.Println("Agent)OnClose")
}

func (this *Agent)Run ()  {
	fmt.Println("Agent)Run () start ")
	for {
		data, err := this.conn.ReadMsg()
		if err != nil {
			fmt.Println("read message: ", err)
			break
		}

		pkgType, body := pkg.Decode(data)
		fmt.Println("pkgType",pkgType, string(body))
	}
	fmt.Println("Agent)Run () end ")
}*/
//
//MaxConnNum:      20000,
//PendingWriteNum: 200,
//MaxMsgLen:       4096,
//HTTPTimeout:     10 * time.Second,
//HeartbeatTimeout: 10 * time.Second,
//LenMsgLen:       2,
//LittleEndian:    false,


func TestSocket(t *testing.T)  {
	tcpServer := new(tcp.TCPServer)
	tcpServer.Addr = ":443"
	tcpServer.MaxConnNum = 20000
	tcpServer.PendingWriteNum = 200
	tcpServer.MaxMsgLen = 4096
	tcpServer.LittleEndian = false
	tcpServer.NewAgent = func(conn tcp.Conn) tcp.Agent {
		return  server.NewAgent(conn)
	}
	tcpServer.Start()
}