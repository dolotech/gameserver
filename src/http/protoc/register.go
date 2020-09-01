package protoc

import jsoniter "github.com/json-iterator/go"

//注册请求
type Register struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (this *Register) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *Register) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//注册返回
type Registered struct {
	Logined
}

func (this *Registered) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *Registered) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//绑定匿名用户请求
type BindGuestUser struct {
	Password string `json:"password"`
	Username string `json:"username"`
	MacSn    string `json:"macSn"`
}

func (this *BindGuestUser) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *BindGuestUser) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//绑定匿名用户返回
type BindGuestUsered struct {
	Code int `json:"code"`
}

func (this *BindGuestUsered) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *BindGuestUsered) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//重置密码请求
type ResetPassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

//重置密码返回
type ResetPassworded struct {
	Code int `json:"code"`
}

//记录玩家最后玩游戏服
type UpdateServerId struct {
	Username string `json:"username"`
	ServerId int `json:"serverId"`
}

type UpdateServerIded struct {
	Code int `json:"code"`
}
