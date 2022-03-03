package server

import (
	"github.com/6QHTSK/ayachan/internal/pkg/ginx"
	"github.com/6QHTSK/ayachan/internal/server/bestdori"
	"github.com/6QHTSK/ayachan/internal/server/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoOutput struct {
	Result bool        `json:"result"`
	List   interface{} `json:"list,omitempty"`
}

func BestdoriOverAllInfo(c *gin.Context) {
	info, err := database.GetBestdoriOverallInfo()
	if err != nil {
		ginx.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   info,
	})
}

func BestdoriFanMadeSearch(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	page, suc := ginx.ConvertQueryInt(c, "page", "0")
	if !suc {
		return
	}
	limit, suc := ginx.ConvertQueryInt(c, "limit", "20")
	if !suc {
		return
	}
	search := bestdori.NewSearch(query, int64(page), int64(limit))

	minLevel, suc := ginx.ConvertQueryInt(c, "level_min", "5")
	if !suc {
		return
	}
	maxLevel, suc := ginx.ConvertQueryInt(c, "level_max", "30")
	if !suc {
		return
	}
	if minLevel > maxLevel {
		temp := maxLevel
		maxLevel = minLevel
		minLevel = temp
	}
	search.FilterLevel(minLevel, maxLevel)

	minDiff, suc := ginx.ConvertQueryInt(c, "diff_min", "0")
	if !suc {
		return
	}
	maxDiff, suc := ginx.ConvertQueryInt(c, "diff_max", "4")
	if !suc {
		return
	}
	if minDiff > maxDiff {
		temp := maxDiff
		maxDiff = minDiff
		minDiff = temp
	}
	search.FilterDiff(minDiff, maxDiff)

	minTimeInt, suc := ginx.ConvertQueryInt(c, "time_min", "0")
	if !suc {
		return
	}
	maxTimeEn, suc := ginx.ConvertQueryInt(c, "time_max_en", "0")
	if !suc {
		return
	}
	if maxTimeEn == 1 {
		maxTimeInt, suc := ginx.ConvertQueryInt(c, "time_max", "23333")
		if !suc {
			return
		}
		search.FilterTime(float64(minTimeInt), float64(maxTimeInt))
	} else {
		search.FilterTimeLow(float64(minTimeInt))
	}

	minNPSInt, suc := ginx.ConvertQueryInt(c, "nps_min", "0")
	if !suc {
		return
	}
	maxNPSEn, suc := ginx.ConvertQueryInt(c, "nps_max_en", "0")
	if !suc {
		return
	}
	if maxNPSEn == 1 {
		maxNPS, suc := ginx.ConvertQueryInt(c, "nps_max", "23333")
		if !suc {
			return
		}
		search.FilterNPS(float64(minNPSInt), float64(maxNPS))
	} else {
		search.FilterNPSLow(float64(minNPSInt))
	}

	SPRhythm, suc := ginx.ConvertQueryInt(c, "sp_rhythm", "0")
	if !suc {
		return
	}
	if SPRhythm > 0 {
		search.FilterSP(SPRhythm == 1)
	}

	regular, suc := ginx.ConvertQueryInt(c, "regular", "0")
	if !suc {
		return
	}
	if regular > 0 {
		search.FilterIrregular(regular == 1)
	}

	documents, totalCount, totalPage, err := search.Search()
	if ginx.ErrorHandle(c, http.StatusInternalServerError, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":     true,
		"docs":       documents,
		"totalCount": totalCount,
		"totalPage":  totalPage,
	})
}

func BestdoriFanMadeGet(c *gin.Context) {
	chartID, suc := ginx.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	chart, err := database.Get(chartID)
	if err != nil {
		ginx.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"doc":    chart,
	})
}
