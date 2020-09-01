package robot

import (
	"gameserver/model"
)

type SortPlayerList []*model.Player

func (this *SortPlayerList) Len() int {
	return len(*this)
}

func (this *SortPlayerList) Less(i, j int) bool {
	if (*this)[i].Power > (*this)[j].Power {
		return true
	}
	return false
}

func (this *SortPlayerList) Swap(i, j int) {
	(*this)[i], (*this)[j] = (*this)[j], (*this)[i]
}
