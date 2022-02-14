package Services

import (
	"ayachan/Log"
	"github.com/robfig/cron/v3"
	"sync"
)

var MysqlSyncRandMutex sync.Mutex
var MysqlSyncFirstMutex sync.Mutex
var MeiliSyncMutex sync.Mutex
var MysqlSyncRand bool
var MysqlSyncFirst bool
var MeiliSync bool

func cronMysqlSyncRand() {
	MysqlSyncRandMutex.Lock()
	if MysqlSyncRand {
		Log.Log.Warning("Hourly Sync Last Job did not finish!")
		return
	}
	MysqlSyncRand = true
	MysqlSyncRandMutex.Unlock()

	Log.Log.Info("Start Sync Mysql hourly")
	_, err := BestdoriFanMadeSyncRand()
	if err != nil {
		Log.Log.Warningf("Failed sync: Error %s", err)
	}

	MysqlSyncRandMutex.Lock()
	MysqlSyncRand = false
	MysqlSyncRandMutex.Unlock()
}

func cronMysqlSyncFirst() {
	MysqlSyncFirstMutex.Lock()
	if MysqlSyncFirst {
		Log.Log.Warning("Minutely Sync Mysql Last Job did not finish!")
		return
	}
	MysqlSyncFirst = true
	MysqlSyncFirstMutex.Unlock()

	Log.Log.Info("Start Sync Mysql minutely")
	_, _, err := BestdoriFanMadeSyncPage(0)
	if err != nil {
		Log.Log.Warningf("Failed sync minute : Error %s", err)
	} else {
		Log.Log.Info("Sync Mysql Success")
	}

	MysqlSyncFirstMutex.Lock()
	MysqlSyncFirst = false
	MysqlSyncFirstMutex.Unlock()
}

func cronMeiliSync() {
	MeiliSyncMutex.Lock()
	if MeiliSync {
		Log.Log.Warning("Minutely Sync Meili Last Job did not finish!")
		return
	}
	MeiliSync = true
	MeiliSyncMutex.Unlock()

	Log.Log.Info("Start Sync MeiliSearch minutely")
	err := MysqlSyncToMeiliSearch()
	if err != nil {
		Log.Log.Warningf("Failed sync minute : Error %s", err)
	} else {
		Log.Log.Info("Sync MeiliSearch Success")
	}
	MeiliSyncMutex.Lock()
	MeiliSync = false
	MeiliSyncMutex.Unlock()
}

func CronSync() {
	c := cron.New(cron.WithSeconds())
	// 每小时的Bestdori随机更新任务
	_, err := c.AddFunc("@hourly", cronMysqlSyncRand)
	if err != nil {
		Log.Log.Fatalf("Cannot add hourly job:%s", err)
	}
	// 每分钟的Bestdori拉取第一页任务 除整点
	_, err = c.AddFunc("0 1-59 * * * *", cronMysqlSyncFirst)
	if err != nil {
		Log.Log.Fatalf("Cannot add minutely Mysql job:%s", err)
	}
	// 每分钟的MeiliSearch同步任务
	_, err = c.AddFunc("30 * * * * *", cronMeiliSync)
	if err != nil {
		Log.Log.Fatalf("Cannot add minutely Meilisearch job:%s", err)
	}
	c.Start()
}
