package cronJob

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/6QHTSK/ayachan/internal/cronJob/config"
	"github.com/6QHTSK/ayachan/internal/pkg/generalModels"
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
	"time"
)

var client = meilisearch.NewClient(meilisearch.ClientConfig{Host: config.Config.MeiliSearch, APIKey: config.Config.MeiliSearchKey})
var index = client.Index("BestdoriFanMade")
var location, _ = time.LoadLocation("Asia/Shanghai")
var SqlDB *sqlx.DB

func init() {
	var err error
	SqlDB, err = sqlx.Open("mysql", config.Config.Mysql)
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

func GetMeiliLastUpdate() (lastUpdate time.Time, err error) {
	res, err := index.Search("", &meilisearch.SearchRequest{
		Limit: 1,
		Sort:  []string{"last_update:desc"},
	})
	if err != nil {
		return lastUpdate, err
	}
	if res.NbHits > 0 {
		item := res.Hits[0].(map[string]interface{})
		lastUpdate, err = time.ParseInLocation(time.RFC3339, item["last_update"].(string), location)
		if err != nil {
			return lastUpdate, err
		}
	}
	return lastUpdate, nil
}

func AddDocument(docs interface{}) (err error) {
	updateTask, err := index.AddDocuments(docs)
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	task, err := client.WaitForTask(updateTask, struct {
		Context  context.Context
		Interval time.Duration
	}{Context: ctx, Interval: time.Second})
	if err != nil {
		return err
	}
	if task.Status == meilisearch.TaskStatusSucceeded {
		return nil
	} else if task.Status == meilisearch.TaskStatusFailed {
		errorStr, _ := json.Marshal(task.Error)
		return fmt.Errorf("%s", errorStr)
	} else if task.Status == meilisearch.TaskStatusUnknown {
		return fmt.Errorf("meiliSearch update status unknown")
	}
	return nil
}

func CheckBestdoriSongVersion(ChartID int) (bool, error) {
	var version int
	err := SqlDB.Get(&version, "SELECT version from BestdoriFanMadeMetrics where chartID=?", ChartID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return version == generalModels.BestdoriFanMadeVersion, err
}

func QueryBestdoriFanMadeByLastUpdate(lastUpdate time.Time) (items []generalModels.BestdoriFanMadeView, err error) {
	err = SqlDB.Select(&items, "SELECT * from BestdoriFanMadeView where lastUpdate > ? order by lastUpdate", lastUpdate)
	return items, err
}

func InsertBestdori(item generalModels.BestdoriChartItem) (err error) {
	//使用Trigger来解决黑名单问题
	//下面的部分使用函数 返回值为BOOL即是否改变了nickname
	var isNicknameChanged int
	err = SqlDB.Get(&isNicknameChanged, "SELECT insertFanMadeF(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.ChartID, item.Title, item.Artists, item.Author.Username, item.Author.Nickname, item.Diff, item.Level, item.SongUrl.Cover, item.SongUrl.Audio, item.Likes, item.PostTime, item.TotalNote, item.TotalTime, item.SPRhythm, item.IrregularInfo.Irregular, generalModels.BestdoriFanMadeVersion, item.Content)
	return err
}

func UpdateBestdori(item generalModels.BestdoriChartUpdateItem) (err error) {
	var statusCode int8
	err = SqlDB.Get(&statusCode, "SELECT updateFanMadeF(?,?,?,?,?,?)", item.ChartID, item.Username, item.Nickname, item.Diff, item.Level, item.Likes)
	return err
}
