package Databases

import (
	"ayachanV2/Config"
	"ayachanV2/Models/chartFormat"
	"fmt"
	"time"
)

func GetLastUpdate() (lastUpdate time.Time, err error) {
	err = SqlDB.QueryRow("SELECT MAX(lastUpdate) from BestdoriFanMade").Scan(&lastUpdate)
	return lastUpdate, err
}

func BestdoriFanMadeSongCount(chartID int) (count int, err error) {
	err = SqlDB.QueryRow("SELECT COUNT(chartID) from BestdoriFanMade where chartID = ?").Scan(&count)
	return count, err
}

// QueryBestdoriSong 需提前检查版本是否适配
func QueryBestdoriSong(chartID int) (item chartFormat.BestdoriChartItem, err error) {
	count, err := BestdoriFanMadeSongCount(chartID)
	if err != nil {
		return item, nil
	}
	if count != 1 {
		return item, fmt.Errorf("not Found / too many")
	}
	err = SqlDB.QueryRow("SELECT chartID,title,artists,BestdoriAuthorList.username,BestdoriAuthorList.nickname,diff,chartLevel,coverURL,songURL,likes,postTime,lastUpdate,totalNote,totalTime,totalNPS,SPRhythm,irregular FROM BestdoriFanMade,BestdoriAuthorList WHERE BestdoriFanMade.author = BestdoriAuthorList.username and chartID = ?", chartID).Scan(
		&item.ChartID, &item.Title, &item.Artists, &item.Author.Username, &item.Author.Nickname, &item.Diff, &item.Level, &item.SongUrl.Cover, &item.SongUrl.Audio, &item.Likes, &item.PostTime, &item.LastUpdateTime, &item.TotalNote, &item.TotalTime, &item.TotalNPS, &item.SPRhythm, &item.Irregular)
	if err != nil {
		return item, err
	}
	return item, err
}

func CheckBestdoriSongVersion(chartID int, diff int) (result bool, err error) {
	var songListVer, authorVer int
	count, err := BestdoriFanMadeSongCount(chartID)
	if err != nil {
		return result, nil
	}
	if count != 1 {
		return result, fmt.Errorf("not found/too many")
	}
	err = SqlDB.QueryRow("SELECT BestdoriFanMade.version, BestdoriAuthorList.version from BestdoriAuthorList,BestdoriFanMade where BestdoriFanMade.author = BestdoriAuthorList.username").Scan(&songListVer, &authorVer)
	return songListVer == Config.BestdoriFanMadeVersion && authorVer == Config.BestdoriAuthorListVersion, err
}

func UpdateBestdori(item chartFormat.BestdoriChartItem) (err error) {
	tx, err := SqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	resultSongList, err := tx.Exec("REPLACE INTO BestdoriFanMade (chartID,title,artists,author,diff,chartLevel,coverURL,songURL,likes,postTime,totalTime,totalNote,totalNPS,spRhythm,irregular,version) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.ChartID, item.Title, item.Artists, item.Author.Username, item.Diff, item.Level, item.SongUrl.Cover, item.SongUrl.Audio, item.Likes, item.PostTime, item.TotalTime, item.TotalNote, item.TotalNPS, item.SPRhythm, item.IrregularInfo.Irregular, Config.BestdoriFanMadeVersion)
	if err != nil {
		return err
	}
	resultAuthorList, err := tx.Exec("REPLACE INTO BestdoriAuthorList (username, nickname, version) VALUES (?,?,?)", item.Author.Username, item.Author.Nickname, Config.BestdoriAuthorListVersion)
	if err != nil {
		return err
	}
	countSongList, _ := resultSongList.RowsAffected()
	countAuthorList, _ := resultAuthorList.RowsAffected()
	if countSongList > 0 && countAuthorList > 0 {
		err = tx.Commit()
	}
	return err
}
