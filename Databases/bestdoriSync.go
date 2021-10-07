package Databases

import (
	"ayachanV2/Config"
	"ayachanV2/Models/chartFormat"
	"fmt"
	"time"
)

func GetLastUpdate() (lastUpdate time.Time, err error) {
	err = SqlDB.QueryRow("SELECT MAX(lastUpdate) from bestdoriSongList").Scan(&lastUpdate)
	return lastUpdate, err
}

func BestdoriSongCount(chartID int, diff int) (count int, err error) {
	err = SqlDB.QueryRow("SELECT COUNT(id) from bestdoriSongList where chartID = ? and (not isOfficial or (isOfficial and diff = diff))").Scan(&count)
	return count, err
}

// QueryBestdoriSong 需提前检查版本是否适配
func QueryBestdoriSong(chartID int, diff int) (item chartFormat.BestdoriChartItem, err error) {
	count, err := BestdoriSongCount(chartID, diff)
	if err != nil {
		return item, nil
	}
	if count != 1 {
		return item, fmt.Errorf("not Found / too many")
	}
	rows, err := SqlDB.Query("SELECT chartID,title,artists,bestdoriAuthorList.username,bestdoriAuthorList.nickname,diff,chartLevel,coverURL,audioURL,isOfficial,likes,postTime,lastUpdate,totalNote,totalTime,totalNPS,SPRhythm,BPMLow,BPMHigh,BPMMain FROM bestdoriSongList,bestdoriAuthorList,bestdoriSongInfo WHERE bestdoriSongList.username = bestdoriAuthorList.username and bestdoriSongList.id = bestdoriSongInfo.id and chartID = ? and (not isOfficial or (isOfficial and diff = ?))", chartID, diff)
	defer rows.Close()
	if err != nil {
		return item, err
	}
	err = rows.Scan(&item.ChartID, &item.Title, &item.Artists, &item.Author.Username, &item.Author.Nickname, &item.Diff, &item.Level, &item.SongUrl.Cover, &item.SongUrl.Audio, &item.Official, &item.Likes, &item.PostTime, &item.LastUpdateTime, &item.TotalNote, &item.TotalTime, &item.TotalNPS, &item.SPRhythm, &item.BPMLow, &item.BPMHigh, &item.MainBPM)
	if err != nil {
		return item, err
	}
	return item, err
}

func CheckBestdoriSongVersion(chartID int, diff int) (result bool, err error) {
	var songListVer, songInfoVer, authorVer int
	count, err := BestdoriSongCount(chartID, diff)
	if err != nil {
		return result, nil
	}
	if count != 1 {
		return result, fmt.Errorf("not found/too many")
	}
	err = SqlDB.QueryRow("SELECT bestdoriSongList.version, bestdoriSongInfo.version, bestdoriAuthorList.version from bestdoriSongInfo,bestdoriAuthorList,bestdoriSongList where bestdoriSongInfo.id = bestdoriSongList.id and bestdoriSongList.username = bestdoriAuthorList.username").Scan(&songListVer, &songInfoVer, &authorVer)
	return songListVer == Config.BestdoriSongListVersion && songInfoVer == Config.BestdoriSongInfoVersion && authorVer == Config.BestdoriAuthorListVersion, err
}

func UpdateBestdori(item chartFormat.BestdoriChartItem) (err error) {
	var id int
	if item.Official {
		id = item.ChartID*10 + int(item.Diff)
	} else {
		id = item.ChartID*10 + 9
	}
	tx, err := SqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	resultSongList, err := tx.Exec("REPLACE INTO bestdoriSongList (id,chartID,title,artists,username,diff,chartLevel,coverURL,audioURL,isOfficial,likes,postTime,lastUpdate,version) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),?)",
		id, item.ChartID, item.Title, item.Artists, item.Author.Username, item.Diff, item.Level, item.SongUrl.Cover, item.SongUrl.Audio, item.Official, item.Likes, item.PostTime, Config.BestdoriSongListVersion)
	if err != nil {
		return err
	}
	resultAuthorList, err := tx.Exec("REPLACE INTO bestdoriAuthorList (username, nickname, version) VALUES (?,?,?)", item.Author.Username, item.Author.Nickname, Config.BestdoriAuthorListVersion)
	if err != nil {
		return err
	}
	resultSongInfo, err := tx.Exec("REPLACE INTO  bestdoriSongInfo (id, totalNote, totalTime, totalNPS, SPRhythm, BPMLow, BPMHigh, BPMMain, version) values (?,?,?,?,?,?,?,?,?)",
		id, item.TotalNote, item.TotalTime, item.TotalNPS, item.SPRhythm, item.BPMLow, item.BPMHigh, item.MainBPM, Config.BestdoriSongInfoVersion)
	if err != nil {
		return err
	}
	countSongList, _ := resultSongList.RowsAffected()
	countAuthorList, _ := resultAuthorList.RowsAffected()
	countSongInfo, _ := resultSongInfo.RowsAffected()
	if countSongList > 0 && countAuthorList > 0 && countSongInfo > 0 {
		err = tx.Commit()
	}
	return err
}
