package Databases

import "github.com/6QHTSK/ayachan/Models/ChartFormat"

func GetChartDisplay(page int, limit int) (ChartSet []ChartFormat.Chart, suc bool) {
	return ChartSet, true
}

func GetChartDisplayID(chartID int) (Chart ChartFormat.Chart, suc bool) {
	return Chart, true
}
