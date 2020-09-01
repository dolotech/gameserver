package protoc

import jsoniter "github.com/json-iterator/go"

//登录请求
type Login struct {
	Password string `json:"password"`
	Username string `json:"username"`
	//Cookie   string `json:"cookie"`
}

func (this *Login) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *Login) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//登录返回
type Logined struct {
	Code         int      `json:"code"`
	Token        string   `json:"token"`
	UID          string   `json:"uid"`
	LastServerID uint      `json:"lastServerId"`
	Servers      []Server `json:"servers"`
}
type Server struct {
	State    int    `json:"state"`
	ServerID uint    `json:"serverId"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (this *Logined) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *Logined) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//匿名登录请求
type GuestLogin struct {
	MacSn string `json:"macSn"`
}

func (this *GuestLogin) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *GuestLogin) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//匿名登录返回
type GuestLogined struct {
	Logined
}

func (this *GuestLogined) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *GuestLogined) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//游戏服务器请求
type GameServers struct {
	Token         string      `json:"token"`
}

//游戏服务器返回
type GameServersed struct {
	Code         int      `json:"code"`
	Servers      []Server `json:"servers"`
}