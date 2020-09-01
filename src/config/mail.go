package config

type Mail struct {
	Head   string `json:"head"`
	Text   string `json:"text"`
	DropId int    `json:"dropId"`
}

var mailPool **MailPool

type MailPool map[string]Mail

func MailPoll() *MailPool {
	return *mailPool
}

func (this *MailPool) Get(typeId string) Mail {
	return (*this)[typeId]
}
