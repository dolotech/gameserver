package protocol

import jsoniter "github.com/json-iterator/go"

//请求开始通关
type StartFightThrough struct {
	NpcId     int `json:"npcId"`     //当前挑战的npc索引 (好像没有用到，一直是0 ANJUN)
	DungeonId int `json:"dungeonId"` //关卡ID
}

func (this *StartFightThrough) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *StartFightThrough) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

//返回开始通关
type StartFightThroughed struct {
	Code    int `json:"code"`
	Stamina int `json:"stamina"`		//当前用户体力值

}
func (this *StartFightThroughed) Encode() ([]byte, error) {
	return jsoniter.Marshal(this)
}

func (this *StartFightThroughed) Decode(data []byte) error {
	return jsoniter.Unmarshal(data, this)
}

