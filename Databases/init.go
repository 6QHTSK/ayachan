package Databases

import (
	"github.com/6QHTSK/ayachan/Config"
	"github.com/6QHTSK/ayachan/Log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
	"time"
)

var SqlDB *sqlx.DB
var MysqlLocation *time.Location
var client *meilisearch.Client
var index *meilisearch.Index

func init() {
	initMysql()
	initMeili()
}

func initMysql() {
	var err error
	SqlDB, err = sqlx.Open("mysql", Config.Config.Database.Mysql)
	if err != nil {
		Log.Log.Fatal(err.Error())
	}
	SqlDB.SetMaxOpenConns(20)
	SqlDB.SetMaxIdleConns(20)
	err = SqlDB.Ping()
	if err != nil {
		Log.Log.Fatal(err.Error())
	}
	MysqlLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		Log.Log.Fatal(err.Error())
	}
}

func initMeili() {
	client = meilisearch.NewClient(meilisearch.ClientConfig{Host: Config.Config.Database.MeiliSearch, APIKey: Config.Config.Database.MeiliSearchKey})
	index = client.Index("BestdoriFanMade")
}
