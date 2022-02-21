package Controllers

import (
	"github.com/6QHTSK/ayachan/Services"
	"github.com/6QHTSK/ayachan/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BestdoriFanMadeSearch(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	page, suc := utils.ConvertQueryInt(c, "page", "0")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "20")
	if !suc {
		return
	}
	search := Services.NewSearch(query, int64(page), int64(limit))

	minLevel, suc := utils.ConvertQueryInt(c, "level_min", "5")
	if !suc {
		return
	}
	maxLevel, suc := utils.ConvertQueryInt(c, "level_max", "30")
	if !suc {
		return
	}
	if minLevel > maxLevel {
		temp := maxLevel
		maxLevel = minLevel
		minLevel = temp
	}
	search.FilterLevel(minLevel, maxLevel)

	minDiff, suc := utils.ConvertQueryInt(c, "diff_min", "0")
	if !suc {
		return
	}
	maxDiff, suc := utils.ConvertQueryInt(c, "diff_max", "4")
	if !suc {
		return
	}
	if minDiff > maxDiff {
		temp := maxDiff
		maxDiff = minDiff
		minDiff = temp
	}
	search.FilterDiff(minDiff, maxDiff)

	minTimeInt, suc := utils.ConvertQueryInt(c, "time_min", "0")
	if !suc {
		return
	}
	maxTimeEn, suc := utils.ConvertQueryInt(c, "time_max_en", "0")
	if !suc {
		return
	}
	if maxTimeEn == 1 {
		maxTimeInt, suc := utils.ConvertQueryInt(c, "time_max", "23333")
		if !suc {
			return
		}
		search.FilterTime(float64(minTimeInt), float64(maxTimeInt))
	} else {
		search.FilterTimeLow(float64(minTimeInt))
	}

	minNPSInt, suc := utils.ConvertQueryInt(c, "nps_min", "0")
	if !suc {
		return
	}
	maxNPSEn, suc := utils.ConvertQueryInt(c, "nps_max_en", "0")
	if !suc {
		return
	}
	if maxNPSEn == 1 {
		maxNPS, suc := utils.ConvertQueryInt(c, "nps_max", "23333")
		if !suc {
			return
		}
		search.FilterNPS(float64(minNPSInt), float64(maxNPS))
	} else {
		search.FilterNPSLow(float64(minNPSInt))
	}

	SPRhythm, suc := utils.ConvertQueryInt(c, "sp_rhythm", "0")
	if !suc {
		return
	}
	if SPRhythm > 0 {
		search.FilterSP(SPRhythm == 1)
	}

	regular, suc := utils.ConvertQueryInt(c, "regular", "0")
	if !suc {
		return
	}
	if regular > 0 {
		search.FilterIrregular(regular == 1)
	}

	documents, totalPage, err := search.Search()
	if utils.ErrorHandle(c, http.StatusInternalServerError, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":    true,
		"docs":      documents,
		"totalPage": totalPage,
	})
}

func BestdoriFanMadeGet(c *gin.Context) {
	chartID, suc := utils.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	chart, err := Services.BestdoriFanMadeGet(chartID)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"doc":    chart,
	})
}
