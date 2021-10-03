package Controllers

import (
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Services"
	"ayachanV2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// MapDataFromBestdori 从Bestdori获取谱面
//@description 根据谱面ID、谱面难度从Bestdori获取谱面
//@Summary 从Bestdori获取谱面
//@tags mapData
//@Param chartID path int true "谱面ID"
//@Param diff query string false "谱面难度差分指定，默认Expert"
//@Param map_format_out query string false "输出谱面格式，默认BestdoriV2"
//@Produce json
//@Success 200 {object} utils.MapOutput "对应的谱面"
//@failed 400	{object} utils.ErrorObject "传入参数有误"
//@failed 404	{object} utils.ErrorObject	"找不到谱面或服务器无连接"
//@Router /map-data/bestdori/:chartID [get]
func MapDataFromBestdori(c *gin.Context) {

	format := c.DefaultQuery("map_format_out", "BestdoriV2")

	chartID, suc := utils.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}

	diff, suc := utils.ConvertQueryInt(c, "diff", strconv.Itoa(int(chartFormat.Diff_Expert)))
	if !suc {
		return
	}

	BestdoriV2map, err := Services.GetMapData(chartID, diff)
	if utils.ErrorHandle(c, http.StatusNotFound, err) {
		return
	}

	Map := BestdoriV2map.Decode()

	suc = utils.ReturnMap(c, Map, format)
	if !suc {
		return
	}
}

// MapData 从用户上传谱面获取谱面
//@description 根据用户上传的谱面获取指定格式谱面
//@Summary 从用户上传谱面获取谱面
//@tags mapData
//@Accept json
//@Param map body utils.MapInput true "需转换的谱面"
//@Produce json
//@Success 200 {object} utils.MapOutput "对应的谱面"
//@failed 400	{object} utils.ErrorObject "传入参数有误"
//@Router /map-data [post]
func MapData(c *gin.Context) {
	Map, Options, suc := utils.ReadMap(c)
	if !suc {
		return
	}

	suc = utils.ReturnMap(c, Map, Options["map_format_out"])
	if !suc {
		return
	}
}
