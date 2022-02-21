package Databases

import "github.com/6QHTSK/ayachan/Models/DatabaseModel"

func GetBestdoriOverallInfo() (info DatabaseModel.BestdoriOverAllInfo, err error) {
	err = SqlDB.Get(&info, "SELECT t1.chartCount,t1.latest,t2.totalNote,t2.totalTime,t3.authorCount FROM (SELECT count(*) as chartCount, max(postTime) as latest from BestdoriFanMade) as t1, (SELECT sum(totalTime) as totalTime,sum(totalNote) as totalNote from BestdoriFanMadeMetrics) as t2, (SELECT count(*) as authorCount from BestdoriAuthor) as t3")
	return info, err
}
