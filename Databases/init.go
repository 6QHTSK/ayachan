package Databases

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var SqlDB *sqlx.DB

func init() {
	initMysql()
}

func initMysql() {
	var err error
	SqlDB, err = sqlx.Open("mysql", databaseInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	SqlDB.SetMaxOpenConns(20)
	SqlDB.SetMaxIdleConns(20)
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
