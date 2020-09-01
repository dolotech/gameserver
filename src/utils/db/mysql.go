package db

import (
	"gameserver/utils/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

var (
	ldb     *gorm.DB
	ldbOnce sync.Once
)

func LDB() *gorm.DB {
	ldbOnce.Do(func() {
		var err error
		ldb, err = gorm.Open("mysql", viper.GetString("ldb.addr"))
		if err != nil {
			log.Error(err)
		}
		ldb.DB().SetMaxIdleConns(viper.GetInt("ldb.max_idle_conns"))
		ldb.DB().SetMaxOpenConns(viper.GetInt("ldb.max_open_conns"))
		ldb.DB().SetConnMaxLifetime(time.Second *time.Duration(viper.GetInt("db.set_conn_maxLifetime")))
		ldb.LogMode(false)
	})
	return ldb
}

func Get() *gorm.DB {
	dbOnce.Do(func() {
		var err error
		db, err = gorm.Open("mysql", viper.GetString("db.addr"))
		if err != nil {
			log.Error(err)
		}
		db.DB().SetMaxIdleConns(viper.GetInt("db.max_idle_conns"))
		db.DB().SetMaxOpenConns(viper.GetInt("db.max_open_conns"))
		db.DB().SetConnMaxLifetime(time.Second *time.Duration(viper.GetInt("db.set_conn_maxLifetime")))
		db.LogMode(true)
	})
	return db
}
