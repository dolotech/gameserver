package pva

import (
	"gameserver/protocol"
)

type GetWeeklyBoardInfo struct {
}

type GetWeeklyBoardInfoed struct {
	Code           int        `json:"code"`
	Opponents      []Opponent `json:"opponents"`
	WeeklyRank     int        `json:"weeklyRank"`
	WeeklyPKTimes  int        `json:"weeklyPKTimes"`
	WeeklyLeftTime int        `json:"weeklyLeftTime"`
}

type GetWeeklyOpponents struct {
}

type GetWeeklyOpponentsed struct {
	Code int `json:"code"`
}

type ChallengeWeeklyOpponent struct {
	OpponentId uint `json:"opponentId"`
}

type ChallengeWeeklyOpponented struct {
	Code     int      `json:"code"`
	Opponent Opponent `json:"opponent"`
}

type FinishChallengeWeeklyOpponent struct {
	Result     int  `json:"result"` //result 0为打输了，1为打赢了，2为打平了
	OpponentId uint `json:"opponentId"`
}

type FinishChallengeWeeklyOpponented struct {
	Code     int             `json:"code"`
	NewRank  int             `json:"newRank"`
	Items    []protocol.Item `json:"items"`
	Opponent Opponent        `json:"opponent"`
}

type BuyWeeklyChallengeTimes struct {
	Times int `json:"times"`
}

type BuyWeeklyChallengeTimesed struct {
	Code           int `json:"code"`
	ChallengeTimes int `json:"challengeTimes"`
}

type RecvWeeklyAwards struct {
}

type RecvWeeklyAwardsed struct {
	Code                int `json:"code"`
	LeftWeeklyAwardTime int `json:"leftWeeklyAwardTime"` //领奖剩余时间
}

type Opponent struct {
	Level                  int               `json:"level"`
	ContinuousWeeklyPKWins int               `json:"continuousWeeklyPKWins"`
	Power                  int               `json:"power"`
	PlayerId               uint              `json:"playerId"`
	Avatar                 int               `json:"avatar"`
	UsedEnergySkill        int               `json:"usedEnergySkill"`
	WeeklyPKWins           int               `json:"weeklyPKWins"`
	Vip                    int               `json:"vip"`
	Nickname               string            `json:"nickname"`
	WeeklyRank             int               `json:"weeklyRank"`
	WeeklyPKLoses          int               `json:"weeklyPKLoses"`
	ResType                int               `json:"resType"`
	Partners               protocol.Partners `json:"partners"` //所有的角色
}
