package main

import (
	"flag"
	"github.com/6QHTSK/ayachan/internal/cronJob"
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	"github.com/manifoldco/promptui"
)

var syncAll, syncRand bool
var syncID int

func init() {
	flag.BoolVar(&syncAll, "a", false, "全部更新")
	flag.BoolVar(&syncRand, "r", false, "部分更新")
	flag.IntVar(&syncID, "i", 0, "更新某谱面")
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "开始更新全部内容(耗时约3小时），是否继续?",
		Items: []string{"是", "否"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		logrus.Log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "是"
}

func main() {
	flag.Parse()
	if syncAll {
		if yesNo() {
			errCode, err := cronJob.BestdoriFanMadeSyncAll()
			if err != nil {
				logrus.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
			} else {
				logrus.Log.Info("Success")
			}
		}
		logrus.Log.Info("Canceled")
	} else if syncRand {
		errCode, err := cronJob.BestdoriFanMadeSyncRand()
		if err != nil {
			logrus.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
		} else {
			logrus.Log.Info("Success")
		}
	} else if syncID > 0 {
		errCode, err := cronJob.BestdoriFanMadeInsertID(syncID)
		if err != nil {
			logrus.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
		} else {
			logrus.Log.Info("Success")
		}
	} else {
		logrus.Log.Info("Start Cron Job Now...")
		cronJob.CronSync()
	}
}
