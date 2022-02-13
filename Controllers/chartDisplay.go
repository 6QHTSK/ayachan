package Controllers

import (
	"ayachanV2/Databases"
	"ayachanV2/Models/ChartFormat"
	"ayachanV2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChartDisplayOutput struct {
	ChartSet []ChartFormat.Chart `json:"chart_set"`
	Result   bool                `json:"result"`
}

type ChartDisplayOutputID struct {
	Chart  ChartFormat.Chart `json:"chart_set"`
	Result bool              `json:"result"`
}

// ChartDisplay 获取谱面展示数据库的谱面
//@description 获取谱面展示数据库的谱面
//@Summary 获取谱面展示数据库的谱面
//@tags chartDisplay
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} ChartDisplayOutput "对应的谱面"
//@Router /chart-display [get]
func ChartDisplay(c *gin.Context) {
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "20")
	if !suc {
		return
	}
	ChartSet, suc := Databases.GetChartDisplay(page, limit)
	if !suc {
		return
	}
	c.JSON(http.StatusOK, ChartDisplayOutput{
		ChartSet: ChartSet,
		Result:   true,
	})
}

// ChartDisplayID 获取谱面展示数据库的指定ID的谱面
//@description 获取谱面展示数据库的指定ID的谱面
//@Summary 获取谱面展示数据库的谱面
//@tags chartDisplay
//@Param chartID path int true "谱面ID"
//@Produce json
//@Success 200 {object} ChartDisplayOutputID "对应的谱面"
//@Router /chart-display/:chartID [get]
func ChartDisplayID(c *gin.Context) {
	chartID, suc := utils.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	Chart, suc := Databases.GetChartDisplayID(chartID)
	if !suc {
		return
	}
	c.JSON(http.StatusOK, ChartDisplayOutputID{
		Chart:  Chart,
		Result: true,
	})
}
