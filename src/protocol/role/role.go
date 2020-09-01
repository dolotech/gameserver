package role


//角色训练
type Train struct {
	PartnerId int `json:"partnerId"`
}

//角色训练返回
type Trained struct {
	Power        int          `json:"power"`
	Code         int          `json:"code"`
}


//角色升星
type UpgradeStar struct {
	PartnerId int `json:"partnerId"`
}

//角色升星返回
type UpgradeStared struct {
	Power        int          `json:"power"`
	Code         int          `json:"code"`
}
