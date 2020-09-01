package protocol

type ChatSendMessage struct {
	Message string `json:"message"`
}

type ChatSendMessaged struct {
	Code int `json:"code"`
}

//跑马灯
type TipSendMessage struct {
	Message string `json:"message"`
}

//跑马灯响应
type TipSendMessaged struct {
	Code int `json:"code"`
}
