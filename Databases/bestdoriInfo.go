package Databases

import (
	"ayachanV2/Models/chartFormat"
)

func GetCharterPostRank(page int, limit int) (list []chartFormat.CharterRankItem, err error) {
	rows, err := SqlDB.Query("select @rank := @rank + 1 as 'rank', unpC.* from (select bestdoriSongList.username,bestdoriAuthorList.nickname,count(id) postCount from bestdoriSongList,bestdoriAuthorList where bestdoriSongList.username = bestdoriAuthorList.username and chartLevel >= 21 and diff >= 3 group by bestdoriSongList.username having count(id) > 5 order by count(id) desc LIMIT ? OFFSET ?) as unpC, (select @rank:=?) as r", limit, page*limit, page*limit)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var item chartFormat.CharterRankItem
		err := rows.Scan(&item.Rank, &item.Author.Username, &item.Author.Nickname, &item.Count)
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}
	return list, nil
}
func GetCharterLikeRank(page int, limit int) (list []chartFormat.CharterRankItem, err error) {
	rows, err := SqlDB.Query("select @rank := @rank + 1 as 'rank', unpC.* from (select bestdoriSongList.username,bestdoriAuthorList.nickname,SUM(likes) likeCount from bestdoriSongList,bestdoriAuthorList where bestdoriSongList.username = bestdoriAuthorList.username and chartLevel >= 21 and diff >= 3 group by bestdoriSongList.username having count(id) > 5 order by SUM(likes) desc LIMIT ? OFFSET ?) as unpC, (select @rank:= ?) as r", limit, page*limit, page*limit)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var item chartFormat.CharterRankItem
		err := rows.Scan(&item.Rank, &item.Author.Username, &item.Author.Nickname, &item.Count)
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}
	return list, nil
}
func SongLikeRank(page int, limit int) (list []chartFormat.SongRankItem, err error) {
	rows, err := SqlDB.Query("select @rank := @rank + 1 as 'rank', unpC.* from (select chartID, title, artists ,bestdoriSongList.username,bestdoriAuthorList.nickname, diff , chartLevel , likes from bestdoriSongList,bestdoriAuthorList where bestdoriSongList.username = bestdoriAuthorList.username order by likes desc LIMIT ? OFFSET ?) as unpC, (select @rank:=?) as r", limit, page*limit, page*limit)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var item chartFormat.SongRankItem
		err := rows.Scan(&item.Rank, &item.ChartID, &item.Title, &item.Artists, &item.Author.Username, &item.Author.Nickname, &item.Diff, &item.Level, &item.Likes)
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}
	return list, nil
}
func GetCharterList() (list []chartFormat.Author, err error) {
	//rows, err := SqlDB.Query("select bestdoriSongList.username,bestdoriAuthorList.nickname from bestdoriSongList,bestdoriAuthorList where bestdoriSongList.username = bestdoriAuthorList.username and chartLevel >= 21 and diff >= 3 group by username having COUNT(id) > 5")
	//if err != nil{
	//	return list,err
	//}
	//defer rows.Close()
	//for rows.Next(){
	//	var item chartFormat.Author
	//	err := rows.Scan(&item.Username, &item.Nickname)
	//	if err != nil{
	//		return list,err
	//	}
	//	list = append(list, item)
	//}
	return list, nil
}
func GetCharterSelfBasic(charter string) (info chartFormat.CharterSelfInfoBasic, err error) {
	return info, err
}
func GetCharterSelfPost(charter string, page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfLikeRank(charter string, page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfNoteRank(charter string, page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfTimeRank(charter string, page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfNPSRank(charter string, page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, err
}
