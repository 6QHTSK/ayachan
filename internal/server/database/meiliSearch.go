package database

import (
	"encoding/json"
	"github.com/6QHTSK/ayachan/internal/pkg/generalModels"
	"github.com/6QHTSK/ayachan/internal/server/config"
	"github.com/meilisearch/meilisearch-go"
	"strconv"
)

var client = meilisearch.NewClient(meilisearch.ClientConfig{Host: config.Config.Database.MeiliSearch, APIKey: config.Config.Database.MeiliSearchKey})
var index = client.Index("BestdoriFanMade")

func Query(q string, page int64, limit int64, filter []string) (charts []generalModels.BestdoriChartItem, totalChart int64, err error) {
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
	var queriedChart []generalModels.BestdoriFanMadeView
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

func Get(chartID int) (chart generalModels.BestdoriChartItem, err error) {
	var dbChart generalModels.BestdoriFanMadeView
	err = index.GetDocument(strconv.Itoa(chartID), &dbChart)
	return dbChart.ToBestdoriChart(), err
}
