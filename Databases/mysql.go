package Databases

import (
	"ayachanV2/Config"
	"ayachanV2/Models/DatabaseModel"
	"ayachanV2/Models/chartFormat"
	"database/sql"
	"fmt"
	"time"
)

// GetLastUpdate 获取上一次的更新时间
func GetLastUpdate() (lastUpdate time.Time, err error) {
	err = SqlDB.QueryRow("SELECT MAX(lastUpdate) from BestdoriFanMade").Scan(&lastUpdate)
	return lastUpdate, err
}

// BestdoriFanMadeSongCount 获取BestdoriFanMade表中某个ID的数量？
func BestdoriFanMadeSongCount(chartID int) (count int, err error) {
	err = SqlDB.QueryRow("SELECT COUNT(chartID) from BestdoriFanMade where chartID = ?", chartID).Scan(&count)
	return count, err
}

func CheckBestdoriSongVersion(ChartID int) (bool, error) {
	var version int
	err := SqlDB.Get(&version, "SELECT version from BestdoriFanMadeMetrics where chartID=?", ChartID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return version == Config.BestdoriFanMadeVersion, err
}

func convertToFanMadeView(item chartFormat.BestdoriChartItem) (dItem DatabaseModel.BestdoriFanMadeView) {
	return DatabaseModel.BestdoriFanMadeView{
		ChartID:    item.ChartID,
		Title:      item.Title,
		Artists:    item.Artists,
		Username:   item.Author.Username,
		Nickname:   item.Author.Nickname,
		Diff:       int(item.Diff),
		ChartLevel: item.Level,
		CoverURL:   item.SongUrl.Cover,
		SongURL:    item.SongUrl.Audio,
		Likes:      item.Likes,
		PostTime:   item.PostTime,
		LastUpdate: time.Now(),
		TotalNote:  item.TotalNote,
		TotalTime:  item.TotalTime,
		TotalNPS:   item.TotalNPS,
		SPRhythm:   item.SPRhythm,
		Irregular:  int(item.Irregular),
		Content:    item.Content,
	}
}

//func QueryBestdoriSongByAuthor(authorName string) (items []chartFormat.BestdoriChartItem, err error) {
//	var databaseItems []DatabaseModel.BestdoriFanMadeView
//	err = SqlDB.Select(&databaseItems, "SELECT * FROM BestdoriFanMadeView WHERE author = ?", authorName)
//	if err != nil {
//		return items, err
//	}
//	for _, di := range databaseItems {
//		items = append(items, di.ToBestdoriChart())
//	}
//	return items, err
//}

func queryBestdoriSongByAuthor(authorName string) (items []DatabaseModel.BestdoriFanMadeView, err error) {
	err = SqlDB.Select(&items, "SELECT * FROM BestdoriFanMadeView WHERE author = ?", authorName)
	if err != nil {
		return items, err
	}
	return items, err
}

func InsertBestdori(item chartFormat.BestdoriChartItem) (err error) {
	//使用Trigger来解决黑名单问题
	//下面的部分使用函数 返回值为BOOL即是否改变了nickname
	var isNicknameChanged int
	err = SqlDB.Get(&isNicknameChanged, "SELECT insertFanMadeF(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.ChartID, item.Title, item.Artists, item.Author.Username, item.Author.Nickname, item.Diff, item.Level, item.SongUrl.Cover, item.SongUrl.Audio, item.Likes, item.PostTime, item.TotalNote, item.TotalTime, item.SPRhythm, item.IrregularInfo.Irregular, Config.BestdoriFanMadeVersion, item.Content)
	if err != nil {
		return err
	}
	if isNicknameChanged == 1 {
		err := UpdateNickname(item.Author.Username)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if isNicknameChanged >= 0 {
		InsertChart(convertToFanMadeView(item))
	}
	return nil
}

func UpdateBestdori(item chartFormat.BestdoriChartUpdateItem) (err error) {
	var statusCode int8
	err = SqlDB.Get(&statusCode, "SELECT updateFanMadeF(?,?,?,?,?,?)", item.ChartID, item.Username, item.Nickname, item.Diff, item.Level, item.Likes)
	if err != nil {
		return err
	}
	if statusCode == -1 {
		return fmt.Errorf("not found")
	} else if statusCode == 1 {
		UpdateChart(item)
	} else if statusCode == 2 {
		err := UpdateNickname(item.Username)
		if err != nil {
			return err
		}
	}
	return nil
}
