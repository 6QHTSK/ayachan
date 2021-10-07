package Controllers

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Services"
	"ayachanV2/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MapInfoOutput struct {
	Models.MapInfo
	Result bool `json:"result"`
}

// MapInfoFromBestdori 从Bestdori获取谱面并分析数据
//@description 根据谱面ID、谱面难度从Bestdori获取谱面并分析得到数据
//@Summary 从Bestdori获取谱面并分析数据
//@tags mapInfo
//@Param chartID path int true "谱面ID"
//@Param diff query string false "谱面难度差分指定，默认Expert"
//@Param map_format_out query string false "输出谱面格式，默认BestdoriV2"
//@Produce json
//@Success 200 {object} MapInfoOutput "解析的谱面信息"
//@failed 400	{object} utils.ErrorObject "传入参数有误"
//@failed 404	{object} utils.ErrorObject	"找不到谱面或服务器无连接"
//@Router /map-info/bestdori/:chartID [get]
func MapInfoFromBestdori(c *gin.Context) {
	chartID, suc := utils.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	diff, suc := utils.ConvertQueryInt(c, "diff", strconv.Itoa(int(chartFormat.Diff_Expert)))
	if !suc {
		return
	}
	BestdoriV2Map, errCode, err := Services.GetMapData(chartID, diff)
	if err != nil {
		utils.ErrorHandle(c, errCode, err)
		return
	}

	Map := BestdoriV2Map.Decode()
	MapInfo := Services.MapInfoGetter(Map, chartFormat.DiffType(diff))

	c.JSON(http.StatusOK, MapInfoOutput{
		MapInfo: MapInfo,
		Result:  true,
	})
}

// MapInfo 从用户上传谱面获取谱面信息
//@description 根据用户上传的谱面获取指定格式谱面信息
//@Summary 从用户上传谱面获取谱面信息
//@tags mapInfo
//@Accept json
//@Param map body utils.MapInput true "需转换的谱面"
//@Produce json
//@Success 200 {object} MapInfoOutput "对应的谱面"
//@failed 400	{object} utils.ErrorObject "传入参数有误"
//@Router /map-info [post]
func MapInfo(c *gin.Context) {
	Map, Options, suc := utils.ReadMap(c)
	if !suc {
		return
	}

	diff, err := strconv.Atoi(Options["diff"])
	if err != nil {
		utils.ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法解析body参数：diff"))
	}
	MapInfo := Services.MapInfoGetter(Map, chartFormat.DiffType(diff))

	c.JSON(http.StatusOK, MapInfoOutput{
		MapInfo: MapInfo,
		Result:  true,
	})
}
