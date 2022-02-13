package Databases

import (
	"ayachanV2/Models/ChartFormat"
)

func GetCharterPostRank(page int, limit int) (list []ChartFormat.CharterRankItem, err error) {
	err = SqlDB.Select(&list, "select author,nickname,count(DISTINCT (title,artists)) postCount from BestdoriFanMadeView where chartLevel >= 21 and diff >= 3 group by author having count(*) > 5 order by count(*) desc LIMIT ? OFFSET ?", limit, page*limit)
	return list, err
}
func GetCharterLikeRank(page int, limit int) (list []ChartFormat.CharterRankItem, err error) {
	err = SqlDB.Select(&list, "select author,nickname,SUM(likes) likeCount from BestdoriFanMadeView group by author order by SUM(likes) desc LIMIT ? OFFSET ?", limit, page*limit)
	return list, err
}
func SongLikeRank(page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	err = SqlDB.Select(&list, "select chartID, title, artists ,author ,nickname, diff , chartLevel , likes from BestdoriFanMadeView order by likes desc LIMIT ? OFFSET ?", limit, page*limit)
	return list, err
}
func GetCharterList() (list []ChartFormat.Author, err error) {
	// err = SqlDB.Select(&list,"select username,nickname from BestdoriFanMadeView where chartLevel >= 21 and diff >= 3 group by username having COUNT(*) > 5")
	return list, nil
}
func GetCharterSelfBasic(charter string) (info ChartFormat.CharterSelfInfoBasic, err error) {
	// err = SqlDB.Get(&info,"select count(DISTINCT (title,artists)) as totalPost, SUM(likes) as totalLike From BestdoriFanMadeView  ")
	return info, err
}
func GetCharterSelfPost(charter string, page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfLikeRank(charter string, page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfNoteRank(charter string, page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfTimeRank(charter string, page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	return list, err
}
func GetCharterSelfNPSRank(charter string, page int, limit int) (list []ChartFormat.SongRankItem, err error) {
	return list, err
}
