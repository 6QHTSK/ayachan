package main

import (
	"flag"
	"fmt"
	"github.com/6QHTSK/ayachan/Config"
	"github.com/6QHTSK/ayachan/Databases"
	"github.com/6QHTSK/ayachan/Log"
	"github.com/6QHTSK/ayachan/Router"
	"github.com/6QHTSK/ayachan/Services"
	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
)

// @title ayachan API
// @version 2.0
// @description api 计算Bestdori谱面难度，获得Bestdori数据，常见Bandori谱面格式转换等

// @contact.name 6QHTSK

// @license.name MIT
// @license.url https://mit-license.org/

// @host 127.0.0.1:8080
// @BasePath /v2

var syncAll bool
var showVer bool
var runAddr string

func init() {
	flag.BoolVar(&syncAll, "s", false, "更新全部内容（耗时约3小时）")
	flag.BoolVar(&showVer, "v", false, "查看版本号")
	flag.StringVar(&runAddr, "a", Config.Config.RunAddr, "运行地址")
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "开始更新全部内容(耗时约3小时），是否继续?",
		Items: []string{"是", "否"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		Log.Log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "是"
}

func main() {
	defer func(SqlDB *sqlx.DB) {
		err := SqlDB.Close()
		if err != nil {
			Log.Log.Fatal(err)
		}
	}(Databases.SqlDB)

	flag.Parse()
	if syncAll {
		if yesNo() {
			_, err := Services.BestdoriFanMadeSyncAll()
			if err != nil {
				Log.Log.Fatal(err)
			}
			err = Services.MysqlSyncToMeiliSearch()
			if err != nil {
				Log.Log.Fatal(err)
			}
		}
	} else if showVer {
		fmt.Println(Config.Version)
	} else {
		lastUpdate, err := Databases.GetLastUpdate()
		if err != nil {
			Log.Log.Warning("读表失败，表为空，最后更新设为0")
		}
		Config.SetLastUpdate(lastUpdate)
		if Config.Config.StartCron {
			Services.CronSync()
		} else {
			Log.Log.Warning("同步任务被config文件禁用")
		}

		router := Router.InitRouter()

		Router.InitAPIV2(router)
		_ = router.Run(runAddr)
	}
}
