package Services

import (
	"ayachan/Config"
	"ayachan/Models/MapFormat"
	"ayachan/utils"
	"fmt"
	"net/http"
	"net/url"
)

type mapDataRequest struct {
	Result bool                      `json:"result"`
	Map    MapFormat.BestdoriV2Chart `json:"map"`
}

// GetMapData 从Bestdori拉取指定ChartID的谱面
func GetMapData(chartID int, diff int) (Map MapFormat.BestdoriV2Chart, errorCode int, err error) {
	mapDataParam, err := url.Parse(fmt.Sprintf("%d/map?diff=%d", chartID, diff))
	mapDataUrl := Config.BestdoriAPIUrl.ResolveReference(mapDataParam)
	var request mapDataRequest
	errorCode, err = utils.HttpGet(mapDataUrl.String(), &request)
	if err != nil {
		return nil, errorCode, err
	}
	if request.Result {
		return request.Map, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, fmt.Errorf("解析错误出现故障")
	}
}
