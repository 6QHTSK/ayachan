package Databases

import (
	"ayachanV2/Models/DatabaseModel"
	"ayachanV2/Models/chartFormat"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"strconv"
)

type failItem struct {
	document interface{}
	err      error
}

var FailList = list.New()

var index = meilisearch.NewClient(meilisearch.ClientConfig{Host: meiliSearchURL, APIKey: meiliSearchKey}).Index("BestdoriFanMade")

func handleErr(document interface{}, err error) {
	FailList.PushBack(failItem{document, err})
}

func handleFailList() {
	count := FailList.Len()
	for i := 0; i < count; i++ {
		elem := FailList.Front()
		item := elem.Value.(failItem)
		FailList.Remove(elem)
		insertDocuments(item.document)
	}
}

func insertDocuments(documents interface{}) {
	updateTask, err := index.UpdateDocuments(documents)
	if err != nil {
		handleErr(documents, err)
		return
	}
	status, err := index.WaitForPendingUpdate(context.TODO(), 1, updateTask)
	if err != nil {
		handleErr(documents, err)
		return
	}
	if status == meilisearch.UpdateStatusFailed {
		updateResp, err := index.GetUpdateStatus(updateTask.UpdateID)
		if err != nil {
			handleErr(documents, err)
			return
		}
		err = fmt.Errorf("meilisearch Error: %s", updateResp.Error)
		handleErr(documents, err)
		return
	}
}

func UpdateNickname(username string) error {
	documents, err := queryBestdoriSongByAuthor(username)
	if err != nil {
		return err
	}
	handleFailList()
	insertDocuments(documents)
	return nil
}

func InsertChart(document DatabaseModel.BestdoriFanMadeView) {
	handleFailList()
	insertDocuments(document)
}

func UpdateChart(document chartFormat.BestdoriChartUpdateItem) {
	handleFailList()
	insertDocuments(document)
}

func Query(q string, page int64, limit int64, filter []string) (charts []chartFormat.BestdoriChartItem, totalChart int64, err error) {
	res, err := index.Search(q, &meilisearch.SearchRequest{
		Offset: page * limit,
		Limit:  limit,
		Filter: filter,
		Sort:   []string{"chart_id:desc"},
	})
	if err != nil {
		return charts, totalChart, err
	}
	jsonStr, err := json.Marshal(res.Hits)
	if err != nil {
		return charts, totalChart, err
	}
	var queriedChart []DatabaseModel.BestdoriFanMadeView
	err = json.Unmarshal(jsonStr, &queriedChart)
	if err != nil {
		return charts, totalChart, err
	}
	for _, chart := range queriedChart {
		charts = append(charts, chart.ToBestdoriChart())
	}
	totalChart = res.NbHits
	return charts, totalChart, err
}

func Get(chartID int) (chart chartFormat.BestdoriChartItem, err error) {
	var dbChart DatabaseModel.BestdoriFanMadeView
	err = index.GetDocument(strconv.Itoa(chartID), &dbChart)
	return dbChart.ToBestdoriChart(), err
}
