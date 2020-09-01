package protocol

import jsoniter "github.com/json-iterator/go"

//主线引导记录
type SetNewStep struct {
	NewStep  int    `json:"newStep"`			//主线引导步骤
}

func (this *SetNewStep) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *SetNewStep) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//主线引导记录返回
type SetNewSteped struct {
	Code int `json:"code"`
}
func (this *SetNewSteped) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *SetNewSteped) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}