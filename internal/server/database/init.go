package database

import (
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	"github.com/6QHTSK/ayachan/internal/server/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var SqlDB *sqlx.DB

func init() {
	var err error
	SqlDB, err = sqlx.Open("mysql", config.Config.Database.Mysql)
	if err != nil {
		logrus.Log.Fatal(err.Error())
	}
	SqlDB.SetMaxOpenConns(20)
	SqlDB.SetMaxIdleConns(20)
	err = SqlDB.Ping()
	if err != nil {
		logrus.Log.Fatal(err.Error())
	}
}
