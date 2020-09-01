package http

import "errors"

var (
	ErrPlayerExist  = errors.New("player already exist, ")
	ErrGetParam   = errors.New("Fatal error get param, ")
	ErrGetPlayer   = errors.New("get player in login error, ")
	ErrGenerateToken   = errors.New("generate token error, ")
	ErrAddPlayer   = errors.New("add player error, ")
	ErrGetUsername   = errors.New("get username error, ")
	ErrHashPass   = errors.New("hash pass error, ")
	ErrGetServers   = errors.New("get servers error, ")
	ErrNoServerList   = errors.New("no server list error, ")
)
