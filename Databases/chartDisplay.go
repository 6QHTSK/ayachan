package Databases

import "ayachanV2/Models/chartFormat"

func GetChartDisplay(page int, limit int) (ChartSet []chartFormat.Chart, suc bool) {
	return ChartSet, true
}

func GetChartDisplayID(chartID int) (Chart chartFormat.Chart, suc bool) {
	return Chart, true
}
