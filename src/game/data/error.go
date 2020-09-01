package data

import "errors"

var (
	ErrUpdateStamina = errors.New("db update stamina error, ")
	ErrSaveTollgateData = errors.New("save player tollgate data error,")
	ErrSaveDungeonRecord = errors.New("save Dungeon Record  error,")
	ErrSavePlayerData = errors.New("save player data error,")
)
