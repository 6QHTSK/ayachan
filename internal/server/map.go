package server

import (
	"fmt"
	"github.com/6QHTSK/ayachan/internal/pkg/ginx"
	"github.com/6QHTSK/ayachan/internal/pkg/httpx"
	"github.com/6QHTSK/ayachan/internal/server/config"
	"github.com/6QHTSK/ayachan/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"net/url"
	"strconv"
)

type mapInput struct {
	MapFormatIn string            `json:"map_format_in"`
	Options     map[string]string `json:"options"`
	Map         interface{}       `json:"map"`
}

type bestdoriV2Input struct {
	Map pkg.BestdoriV2Chart `json:"map"`
}

type mapOutput struct {
	Result bool        `json:"result"`
	Map    interface{} `json:"map"`
}

type mapInfoOutput struct {
	pkg.MapInfo
	Result bool `json:"result"`
}

type mapDataRequest struct {
	Result bool                `json:"result"`
	Map    pkg.BestdoriV2Chart `json:"map"`
}

// readMap 读取内容字段内容
func readMap(c *gin.Context) (Map pkg.InputMap, Options map[string]string, suc bool) {
	var inputOptions mapInput
	err := c.ShouldBindBodyWith(&inputOptions, binding.JSON)
	if err != nil {
		ginx.ErrorHandle(c, http.StatusBadRequest, err)
		return Map, Options, false
	}

	Options = inputOptions.Options

	if inputOptions.MapFormatIn == "BestdoriV2" {
		var input bestdoriV2Input
		err := c.ShouldBindBodyWith(&input, binding.JSON)
		if err != nil {
			ginx.ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法处理传入谱面"))
			return Map, Options, false
		}
		Map = &input.Map
	} else {
		ginx.ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("不支持的传入格式：%s", inputOptions.MapFormatIn))
		return Map, Options, false
	}

	return Map, Options, true
}

// returnMap 返回指定格式的谱面 会处理错误
func returnMap(c *gin.Context, Map pkg.Chart, formatOut string) (suc bool) {
	if formatOut == "BestdoriV2" {
		var BestdoriV2Map pkg.BestdoriV2Chart
		err := BestdoriV2Map.Encode(Map)
		if err != nil {
			ginx.ErrorHandle(c, http.StatusInternalServerError, fmt.Errorf("转换内部格式错误：%s", err.Error()))
			return false
		}
		c.JSON(http.StatusOK, mapOutput{
			Result: true,
			Map:    BestdoriV2Map,
		})
	} else {
		ginx.ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("不支持的传出格式：%s", formatOut))
		return false
	}
	return true
}

// getMapData 从Bestdori拉取指定ChartID的谱面
func getMapData(chartID int, diff int) (Map pkg.BestdoriV2Chart, errorCode int, err error) {
	mapDataParam, err := url.Parse(fmt.Sprintf("%d/map?diff=%d", chartID, diff))
	mapDataUrl := config.BestdoriAPIUrl.ResolveReference(mapDataParam)
	var request mapDataRequest
	errorCode, err = httpx.HttpGet(mapDataUrl.String(), &request)
	if err != nil {
		return nil, errorCode, err
	}
	if request.Result {
		return request.Map, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, fmt.Errorf("解析错误出现故障")
	}
}

func MapInfoFromBestdori(c *gin.Context) {
	chartID, suc := ginx.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	diff, suc := ginx.ConvertQueryInt(c, "diff", strconv.Itoa(3)) // Expert
	if !suc {
		return
	}
	BestdoriV2Map, errCode, err := getMapData(chartID, diff)
	if err != nil {
		ginx.ErrorHandle(c, errCode, err)
		return
	}
	MapInfo, err := pkg.MapInfoGetter(&BestdoriV2Map, diff)
	if err != nil {
		ginx.ErrorHandle(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, mapInfoOutput{
		MapInfo: MapInfo,
		Result:  true,
	})
}

func MapInfo(c *gin.Context) {
	inputMap, Options, suc := readMap(c)
	if !suc {
		return
	}

	diff, err := strconv.Atoi(Options["diff"])
	if err != nil {
		ginx.ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法解析body参数：diff"))
		return
	}
	MapInfo, err := pkg.MapInfoGetter(inputMap, diff)
	if err != nil {
		ginx.ErrorHandle(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, mapInfoOutput{
		MapInfo: MapInfo,
		Result:  true,
	})
}

func MapDataFromBestdori(c *gin.Context) {

	format := c.DefaultQuery("map_format_out", "BestdoriV2")

	chartID, suc := ginx.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}

	diff, suc := ginx.ConvertQueryInt(c, "diff", strconv.Itoa(3)) // EXPERT
	if !suc {
		return
	}

	BestdoriV2map, errCode, err := getMapData(chartID, diff)
	if err != nil {
		ginx.ErrorHandle(c, errCode, err)
		return
	}

	Map := BestdoriV2map.Decode()

	suc = returnMap(c, Map, format)
	if !suc {
		return
	}
}

func MapData(c *gin.Context) {
	Map, Options, suc := readMap(c)
	if !suc {
		return
	}
	check, err := Map.MapCheck()
	if !check {
		ginx.ErrorHandle(c, http.StatusBadRequest, err)
		return
	}
	suc = returnMap(c, Map.Decode(), Options["map_format_out"])
	if !suc {
		return
	}
}
