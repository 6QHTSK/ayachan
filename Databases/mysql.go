package Databases

import (
	"database/sql"
	"github.com/6QHTSK/ayachan/Config"
	"github.com/6QHTSK/ayachan/Models/ChartFormat"
	"github.com/6QHTSK/ayachan/Models/DatabaseModel"
	"time"
)

// GetLastUpdate 获取上一次的更新时间
func GetLastUpdate() (lastUpdate time.Time, err error) {
	err = SqlDB.QueryRow("SELECT MAX(lastUpdate) from BestdoriFanMade").Scan(&lastUpdate)
	return lastUpdate, err
}

// BestdoriFanMadeSongCount 获取BestdoriFanMade表中某个ID的数量？
//func BestdoriFanMadeSongCount(chartID int) (count int, err error) {
//	err = SqlDB.QueryRow("SELECT COUNT(chartID) from BestdoriFanMade where chartID = ?", chartID).Scan(&count)
//	return count, err
//}

func CheckBestdoriSongVersion(ChartID int) (bool, error) {
	var version int
	err := SqlDB.Get(&version, "SELECT version from BestdoriFanMadeMetrics where chartID=?", ChartID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return version == Config.BestdoriFanMadeVersion, err
}

func QueryBestdoriFanMadeByLastUpdate(lastUpdate time.Time) (items []DatabaseModel.BestdoriFanMadeView, err error) {
	err = SqlDB.Select(&items, "SELECT * from BestdoriFanMadeView where lastUpdate > ? order by lastUpdate", lastUpdate)
	return items, err
}

func InsertBestdori(item ChartFormat.BestdoriChartItem) (err error) {
	//使用Trigger来解决黑名单问题
	//下面的部分使用函数 返回值为BOOL即是否改变了nickname
	var isNicknameChanged int
	err = SqlDB.Get(&isNicknameChanged, "SELECT insertFanMadeF(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.ChartID, item.Title, item.Artists, item.Author.Username, item.Author.Nickname, item.Diff, item.Level, item.SongUrl.Cover, item.SongUrl.Audio, item.Likes, item.PostTime, item.TotalNote, item.TotalTime, item.SPRhythm, item.IrregularInfo.Irregular, Config.BestdoriFanMadeVersion, item.Content)
	return err
}

func UpdateBestdori(item ChartFormat.BestdoriChartUpdateItem) (err error) {
	var statusCode int8
	err = SqlDB.Get(&statusCode, "SELECT updateFanMadeF(?,?,?,?,?,?)", item.ChartID, item.Username, item.Nickname, item.Diff, item.Level, item.Likes)
	return err
}
