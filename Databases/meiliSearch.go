package Databases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/6QHTSK/ayachan/Models/ChartFormat"
	"github.com/6QHTSK/ayachan/Models/DatabaseModel"
	"github.com/meilisearch/meilisearch-go"
	"strconv"
	"time"
)

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
		lastUpdate, err = time.ParseInLocation(time.RFC3339, item["last_update"].(string), MysqlLocation)
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

func Query(q string, page int64, limit int64, filter []string) (charts []ChartFormat.BestdoriChartItem, totalChart int64, err error) {
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

func Get(chartID int) (chart ChartFormat.BestdoriChartItem, err error) {
	var dbChart DatabaseModel.BestdoriFanMadeView
	err = index.GetDocument(strconv.Itoa(chartID), &dbChart)
	return dbChart.ToBestdoriChart(), err
}
