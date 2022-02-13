package Databases

import (
	"ayachanV2/Models/DatabaseModel"
	"ayachanV2/Models/chartFormat"
	"encoding/json"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"io/ioutil"
	"log"
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
	task, err := client.WaitForTask(updateTask)
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

func HandleErr(document DatabaseModel.BestdoriFanMadeView, err error) {
	// 不再处理failList
	currentTimeStamp := time.Now().Unix()
	logName := fmt.Sprintf("log/%d.%d.log", document.ChartID, currentTimeStamp)
	fErr := ioutil.WriteFile(logName, []byte(fmt.Sprintf("%d,%s", document.ChartID, err)), 0666)
	if fErr != nil {
		log.Printf("Fail To Write File %s\n", logName)
		log.Printf("Sync Fail , chartID = %d\n", document.ChartID)
	}
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
