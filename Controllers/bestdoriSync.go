package Controllers

import (
	"ayachanV2/Config"
	"ayachanV2/Services"
	"ayachanV2/utils"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"time"
)

func SyncAll(c *gin.Context) {
	go func() {
		errCode, err := Services.BestdoriFanMadeSyncAll()
		if err != nil {
			log.Printf("SyncRand %d,%s\n", errCode, err.Error())
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
			log.Printf("SyncRand %d,%s\n", errCode, err.Error())
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
			log.Printf("SyncRand %d,%s\n", errCode, err.Error())
			return
		}
	}()
	c.JSON(http.StatusAccepted, InfoOutput{Result: true})
}

func CronSync() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		log.Print("Start Sync hourly")
		_, err := Services.BestdoriFanMadeSyncRand()
		if err != nil {
			log.Printf("Failed sync: Error %s", err)
		}
	})
	if err != nil {
		log.Fatalf("Cannot add hourly job:%s", err)
	}
	_, err = c.AddFunc("1-59 * * * *", func() {
		log.Print("Start Sync minutely")
		_, _, err := Services.BestdoriFanMadeSyncPage(0)
		if err != nil {
			log.Printf("Failed sync minute : Error %s", err)
		} else {
			log.Print("Sync minutely Success")
		}
	})
	if err != nil {
		log.Fatalf("Cannot add minutely job:%s", err)
	}
	c.Start()
}

//func SyncBlackList(c *gin.Context){
//	err := Services.RefreshBlackList()
//	if err != nil{
//		utils.ErrorHandle(c,http.StatusInternalServerError,err)
//		return
//	}
//	c.JSON(http.StatusOK,InfoOutput{Result: true})
//}