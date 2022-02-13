package Controllers

import (
	"ayachanV2/Config"
	"ayachanV2/Log"
	"ayachanV2/Services"
	"ayachanV2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SyncAll(c *gin.Context) {
	go func() {
		errCode, err := Services.BestdoriFanMadeSyncAll()
		if err != nil {
			Log.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
		}
	}()
	c.JSON(http.StatusAccepted, InfoOutput{Result: true})
}

func SyncRand(c *gin.Context) {
	if time.Since(Config.LastUpdate) < time.Hour {
		c.JSON(http.StatusForbidden, InfoOutput{Result: false})
		return
	}
	go func() {
		errCode, err := Services.BestdoriFanMadeSyncRand()
		if err != nil {
			Log.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
		}
	}()
	c.JSON(http.StatusAccepted, InfoOutput{Result: true})
}

func SyncChartID(c *gin.Context) {
	chartID, suc := utils.ConvertParamInt(c, "chartID")
	if !suc {
		return
	}
	go func() {
		errCode, err := Services.BestdoriFanMadeInsertID(chartID)
		if err != nil {
			Log.Log.Warningf("SyncRand %d,%s\n", errCode, err.Error())
			return
		}
	}()
	c.JSON(http.StatusAccepted, InfoOutput{Result: true})
}

//func SyncBlackList(c *gin.Context){
//	err := Services.RefreshBlackList()
//	if err != nil{
//		utils.ErrorHandle(c,http.StatusInternalServerError,err)
//		return
//	}
//	c.JSON(http.StatusOK,InfoOutput{Result: true})
//}
