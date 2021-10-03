package Databases

import "ayachanV2/Models/chartFormat"

func GetCharterPostRank(page int, limit int) (list []chartFormat.CharterRankItem, err error) {
	return list, nil
}
func GetCharterLikeRank(page int, limit int) (list []chartFormat.CharterRankItem, err error) {
	return list, nil
}
func SongLikeRank(page int, limit int) (list []chartFormat.SongRankItem, err error) {
	return list, nil
}
func GetCharterList() (list []chartFormat.Author, err error) {
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
